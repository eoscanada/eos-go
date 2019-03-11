package blog

import (
	"encoding/binary"
	"os"

	"github.com/pkg/errors"
)

var ErrBlockIndexNotOpened = errors.New("block index must be open prior doing this operation")

type BlockIndex interface {
	Open() error
	Close() error

	ReadBlockByteOffset(blockNum uint32) (uint64, error)
}

type MemoryBlockIndex struct {
	buffer []byte
}

func NewMemoryBlockIndex(buffer []byte) (*MemoryBlockIndex, error) {
	if len(buffer)%8 != 0 {
		return nil, errors.Errorf("buffer lenght should be a multiple of 8")
	}

	return &MemoryBlockIndex{
		buffer: buffer,
	}, nil
}

func (i *MemoryBlockIndex) Open() (err error)  { return nil }
func (i *MemoryBlockIndex) Close() (err error) { return nil }

func (i *MemoryBlockIndex) ReadBlockByteOffset(blockNum uint32) (uint64, error) {
	if blockNum == 0 {
		return 0, errors.Errorf("blockNum should be greater than 0, got %d", blockNum)
	}

	byteOffset := (blockNum - 1) * 8

	return binary.LittleEndian.Uint64(i.buffer[byteOffset:]), nil
}

type FileBlockIndex struct {
	filename string

	handle *os.File
}

func NewFileBlockIndex(filename string) *FileBlockIndex {
	return &FileBlockIndex{
		filename: filename,
	}
}

func (i *FileBlockIndex) Open() (err error) {
	i.handle, err = os.Open(i.filename)
	if err != nil {
		return errors.Wrapf(err, "unable to open file block index %s", i.filename)
	}

	return nil
}

func (i *FileBlockIndex) Close() (err error) {
	if i.handle != nil {
		currentHanddle := i.handle
		i.handle = nil

		err = currentHanddle.Close()
		if err != nil {
			return errors.Wrapf(err, "unable to close file block index %s", i.filename)
		}
	}

	return nil
}

func (i *FileBlockIndex) ReadBlockByteOffset(blockNum uint32) (uint64, error) {
	if i.handle == nil {
		return 0, ErrBlockIndexNotOpened
	}

	if blockNum == 0 {
		return 0, errors.Errorf("blockNum should be greater than 0, got %d", blockNum)
	}

	buffer := make([]byte, 8)
	byteOffset := int64((blockNum - 1) * 8)

	byteReadCount, err := i.handle.ReadAt(buffer, byteOffset)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to seek file to byte offset %d", byteOffset)
	}

	if byteReadCount != len(buffer) {
		return 0, errors.Errorf("expected to have read %d bytes, only read %d", len(buffer), byteReadCount)
	}

	return binary.LittleEndian.Uint64(buffer), nil
}
