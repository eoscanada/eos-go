package blog

// Block Index Documentation
//  (Copied from )
//
// The block log is an external append only log of the blocks with a header. Blocks should only
// be written to the log after they irreverisble as the log is append only. The log is a doubly
// linked list of blocks. There is a secondary index file of only block positions that enables
// O(1) random access lookup by block number.
//
// +---------+----------------+---------+----------------+-----+------------+-------------------+
// | Block 1 | Pos of Block 1 | Block 2 | Pos of Block 2 | ... | Head Block | Pos of Head Block |
// +---------+----------------+---------+----------------+-----+------------+-------------------+
//
// +----------------+----------------+-----+-------------------+
// | Pos of Block 1 | Pos of Block 2 | ... | Pos of Head Block |
// +----------------+----------------+-----+-------------------+
//
// The block log can be walked in order by deserializing a block, skipping 8 bytes, deserializing a
// block, repeat... The head block of the file can be found by seeking to the position contained
// in the last 8 bytes the file. The block log can be read backwards by jumping back 8 bytes, following
// the position, reading the block, jumping back 8 bytes, etc.
//
// Blocks can be accessed at random via block number through the index file. Seek to 8 * (block_num - 1)
// to find the position of the block in the main file.
//
// The main file is the only file that needs to persist. The index file can be reconstructed during a
// linear scan of the main file.
//

import (
	"encoding/binary"
	"io"
	"os"

	"go.uber.org/zap"

	"github.com/eoscanada/eos-go"

	"github.com/pkg/errors"
)

var ErrBlockLogNotOpened = errors.New("block log must be open prior doing this operation")

type BlockLog interface {
	Open() error
	Close() error

	ReadFirst() (*eos.SignedBlock, error)
	ReadHead() (*eos.SignedBlock, error)
	ReadBlock(blockNum uint32) (*eos.SignedBlock, error)
}

type MemoryBlockLog struct {
	buffer []byte
	index  BlockIndex
}

func NewMemoryBlockLog(logBuffer []byte, indexBuffer []byte) (*MemoryBlockLog, error) {
	index, err := NewMemoryBlockIndex(indexBuffer)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create memory index")
	}

	return &MemoryBlockLog{
		buffer: logBuffer,
		index:  index,
	}, nil
}

func (i *MemoryBlockLog) Open() (err error) {
	return i.index.Open()
}

func (i *MemoryBlockLog) Close() (err error) {
	return i.index.Close()
}

func (i *MemoryBlockLog) ReadFirst() (*eos.SignedBlock, error) {
	return i.ReadBlock(1)
}

func (i *MemoryBlockLog) ReadHead() (*eos.SignedBlock, error) {
	panic(errors.New("not implemented yet"))
}

func (i *MemoryBlockLog) ReadBlock(blockNum uint32) (*eos.SignedBlock, error) {
	blockByteOffset, err := readBlockPosition(i.index, blockNum)
	if err != nil {
		return nil, err
	}

	return i.readBlockAtOffset(int64(blockByteOffset))
}

func (i *MemoryBlockLog) readBlockAtOffset(byteOffset int64) (*eos.SignedBlock, error) {
	var block *eos.SignedBlock
	err := eos.UnmarshalBinary(i.buffer[byteOffset:], block)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to deserialize block from binary at offset %d", byteOffset)
	}

	return block, nil
}

type FileBlockLog struct {
	filename string
	index    BlockIndex

	handle *os.File
}

func NewFileBlockLog(logFilename string, indexFilename string) *FileBlockLog {
	return &FileBlockLog{
		filename: logFilename,
		index:    NewFileBlockIndex(indexFilename),
	}
}

func (i *FileBlockLog) Open() (err error) {
	err = i.index.Open()
	if err != nil {
		return errors.Wrap(err, "unable to open block index")
	}

	i.handle, err = os.Open(i.filename)
	if err != nil {
		return errors.Wrapf(err, "unable to open file block index %s", i.filename)
	}

	return nil
}

func (i *FileBlockLog) Close() (err error) {
	defer func() {
		err = i.index.Close()
		if err != nil {
			err = errors.Wrap(err, "unable to close block index")
		}
	}()

	if i.handle != nil {
		currentHanddle := i.handle
		i.handle = nil

		err = currentHanddle.Close()
		if err != nil {
			return errors.Wrapf(err, "unable to close file block log %s", i.filename)
		}
	}

	return nil
}

func (i *FileBlockLog) ReadFirst() (*eos.SignedBlock, error) {
	return i.ReadBlock(1)
}

func (i *FileBlockLog) ReadHead() (*eos.SignedBlock, error) {
	panic(errors.New("not implemented yet"))
}

func (i *FileBlockLog) ReadBlock(blockNum uint32) (*eos.SignedBlock, error) {
	blockByteOffset, err := i.index.ReadBlockByteOffset(blockNum)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get block position for block num %d", blockNum)
	}

	return i.readBlockAtOffset(int64(blockByteOffset))
}

func (i *FileBlockLog) ForEach(callback func(block *eos.SignedBlock) error) error {
	err := i.checkAndSeek(0)
	if err != nil {
		return err
	}

	// Should be performed once at initialization time, when initializing the block log!
	offset, err := i.peekOffset()
	if err != nil {
		return errors.Wrap(err, "unable to peek offset")
	}

	zlog.Debug("Offset prior reading version", zap.Int64("offset", offset))
	version, err := readVersion(i.handle)
	if err != nil {
		return errors.Wrap(err, "unable to read version")
	}

	offset, err = i.peekOffset()
	if err != nil {
		return errors.Wrap(err, "unable to peek offset")
	}

	zlog.Debug("Offset prior reading first block num", zap.Int64("offset", offset))
	firstBlockNum, err := readFirstBlockNum(i.handle)
	if err != nil {
		return errors.Wrap(err, "unable to read first block num")
	}

	offset, err = i.peekOffset()
	if err != nil {
		return errors.Wrap(err, "unable to peek offset")
	}

	zlog.Debug("Offset prior reading genesis state", zap.Int64("offset", offset))
	genesisState, err := readGenesisState(i.handle)
	if err != nil {
		return errors.Wrap(err, "unable to read genesis state")
	}

	// FIXME: We hit a road block here where the offset is at the end of the block log
	//        completely instead of being right after the end of the genesis state struct.
	//
	//        The problem lies down to the `eos.UnmarshalBinaryReader` method that we
	//        use to decode the struct, the `GenesisState` struct in this case. The
	//        underlying decoder implementation relies on having the full buffer of
	//        the struct available. So, the `eos.UnmarshalBinaryReader` when it starts
	//        do an `ioutil.ReadAll` call to read the full reader and obtain the buffer.
	//
	//        This is highly problematic in this case because we would read the full
	//        blocks log in memory which in this case, it's not viable. Even when the
	//        full blocks log fit in memory (like in our playground files), it's still
	//        a big problem since the offset is all wrong and we need to start at the
	//        right offset for correctly reading the next data (a block or a block
	//        position)
	offset, err = i.peekOffset()
	if err != nil {
		return errors.Wrap(err, "unable to peek offset")
	}

	zlog.Debug("Offset prior reading totem", zap.Int64("offset", offset))
	totem, err := readTotem(i.handle)
	if err != nil {
		return errors.Wrap(err, "unable to read totem")
	}

	zlog.Info("Read block log header section",
		zap.Uint32("version", version),
		zap.Uint32("first_block_num", firstBlockNum),
		zap.Reflect("genesis_state", genesisState),
		zap.Uint64("totem", totem),
	)

	offset, err = i.peekOffset()
	if err != nil {
		return errors.Wrap(err, "unable to peek offset")
	}

	zlog.Debug("Offset prior reading first block", zap.Int64("offset", offset))

	for err == nil {
		// FIXME: We should if there is 0 bytes available, if it's the case, it's the end
		block, err := readBlock(i.handle)
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return errors.Wrap(err, "unable to read block")
		}

		callbackErr := callback(block)
		if callbackErr != nil {
			return callbackErr
		}

		_, err = readBlockByteOffset(i.handle)
		if err != nil {
			return errors.Wrap(err, "unable to read block byte offset that should follow block")
		}
	}

	return nil
}

func (i *FileBlockLog) readBlockAtOffset(byteOffset int64) (*eos.SignedBlock, error) {
	err := i.checkAndSeek(byteOffset)
	if err != nil {
		return nil, err
	}

	block, err := readBlock(i.handle)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to deserialize block from file at offset %d", byteOffset)
	}

	return block, nil
}

func (i *FileBlockLog) checkAndSeek(byteOffset int64) error {
	if i.handle == nil {
		return ErrBlockLogNotOpened
	}

	newByteOffset, err := i.handle.Seek(byteOffset, 0)
	if err != nil {
		return errors.Wrapf(err, "unable to seek file to byte offset %d", byteOffset)
	}

	if newByteOffset != byteOffset {
		return errors.Errorf("expected to have been at offset %d, but now at offset %d", byteOffset, newByteOffset)
	}

	return nil
}

func (i *FileBlockLog) peekOffset() (int64, error) {
	currentOffset, err := i.handle.Seek(0, 1)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to peek offset of file")
	}

	return currentOffset, nil
}

func readBlockPosition(index BlockIndex, blockNum uint32) (uint64, error) {
	blockByteOffset, err := index.ReadBlockByteOffset(blockNum)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to get block position for block num %d", blockNum)
	}

	return blockByteOffset, nil
}

func readVersion(reader io.Reader) (uint32, error) {
	buffer := make([]byte, 4)
	err := readBuffer(reader, buffer)
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint32(buffer), nil
}

func readFirstBlockNum(reader io.Reader) (uint32, error) {
	buffer := make([]byte, 4)
	err := readBuffer(reader, buffer)
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint32(buffer), nil
}

func readGenesisState(reader io.Reader) (*eos.GenesisState, error) {
	var state *eos.GenesisState
	err := eos.UnmarshalBinaryReader(reader, &state)
	if err != nil {
		return nil, err
	}

	return state, nil
}

func readTotem(reader io.Reader) (uint64, error) {
	buffer := make([]byte, 8)
	err := readBuffer(reader, buffer)
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint64(buffer), nil
}

func readBlock(reader io.Reader) (*eos.SignedBlock, error) {
	var block *eos.SignedBlock
	err := eos.UnmarshalBinaryReader(reader, &block)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func readBlockByteOffset(reader io.Reader) (uint64, error) {
	buffer := make([]byte, 8)
	err := readBuffer(reader, buffer)
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint64(buffer), nil
}

func readBuffer(reader io.Reader, buffer []byte) error {
	byteReadCount, err := reader.Read(buffer)
	if err != nil {
		return errors.Wrap(err, "unable to read buffer")
	}

	if byteReadCount != len(buffer) {
		return errors.Errorf("expected to have read %d bytes from buffer, only read %d", len(buffer), byteReadCount)
	}

	return nil
}
