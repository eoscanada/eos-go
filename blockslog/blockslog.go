package blockslog

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/eoscanada/eos-go"
)

type Reader struct {
	filename string
	fl       *os.File

	Version       uint8
	FirstBlockNum uint32
	ChainID       string

	firstOffset int64
	nextOffset  int64
	prevOffset  int64
}

func NewReader(filename string) *Reader {
	return &Reader{filename: filename}
}

func (r *Reader) Close() error {
	if r.fl != nil {
		return r.fl.Close()
	}
	return nil
}

func (r *Reader) ReadHeader() error {
	fl, err := os.Open(r.filename)
	if err != nil {
		return err
	}
	r.fl = fl

	versionData := make([]byte, 4)
	_, err = fl.Read(versionData)
	if err != nil {
		return err
	}

	version := versionData[0]
	r.Version = version

	// fmt.Println("Version", version)

	firstBlockData := []byte{1, 0, 0, 0}

	if version > 1 {
		// fmt.Println("Reading first block")
		_, err = fl.Read(firstBlockData)
		if err != nil {
			return err
		}
	}

	firstBlockNum := binary.LittleEndian.Uint32(firstBlockData)

	r.FirstBlockNum = firstBlockNum
	// fmt.Println("First block", firstBlockNum)

	// Certain conditions where the genesis state is written:
	// bool block_log::contains_genesis_state(uint32_t version, uint32_t first_block_num) {
	//    return version <= 2 || first_block_num == 1;
	// }

	chainID := make([]byte, 32)
	if version >= 3 && firstBlockNum > 1 {
		_, err = fl.Read(chainID)
		if err != nil {
			return err
		}
	}

	r.ChainID = hex.EncodeToString(chainID)

	if version != 1 {
		totem := make([]byte, 8)
		_, err := fl.Read(totem)
		if err != nil {
			return err
		}
	}

	startPos, err := fl.Seek(0, os.SEEK_CUR)
	if err != nil {
		return err
	}

	r.firstOffset = startPos

	r.First()

	return nil
}

func (r *Reader) Next() (*eos.SignedBlock, []byte, error) {
	if r.nextOffset == -1 {
		return nil, nil, io.EOF
	}

	_, err := r.fl.Seek(r.nextOffset, os.SEEK_SET)
	if err != nil {
		return nil, nil, err
	}

	blk, bytesRead, err := r.readSignedBlock()
	if err != nil {
		return nil, nil, err
	}

	r.prevOffset = r.nextOffset
	r.nextOffset = r.nextOffset + int64(len(bytesRead)) + 8

	return blk, bytesRead, nil
}

func (r *Reader) First() {
	r.nextOffset = r.firstOffset
	r.prevOffset = -1
}

func (r *Reader) Last() error {
	_, err := r.fl.Seek(-8, os.SEEK_END)
	if err != nil {
		return err
	}

	cnt := make([]byte, 8)

	_, err = r.fl.Read(cnt)
	if err != nil {
		return err
	}

	lastBlockOffset := binary.LittleEndian.Uint64(cnt)

	// TODO: perhaps check if the blocks log is empty, and the whole offset
	// thing is less than a single block's size..

	r.prevOffset = int64(lastBlockOffset)
	r.nextOffset = -1

	return nil
}

func (r *Reader) Prev() (*eos.SignedBlock, []byte, error) {
	if r.prevOffset == -1 {
		return nil, nil, io.EOF
	}

	_, err := r.fl.Seek(r.prevOffset-8, os.SEEK_SET)
	if err != nil {
		return nil, nil, fmt.Errorf("seek -8: %w", err)
	}

	prevOffsetBin := make([]byte, 8)
	_, err = r.fl.Read(prevOffsetBin)
	if err != nil {
		return nil, nil, fmt.Errorf("read offset: %w", err)
	}
	prevOffset := binary.LittleEndian.Uint64(prevOffsetBin)

	blk, bytesRead, err := r.readSignedBlock()
	if err != nil {
		return nil, nil, fmt.Errorf("read signed block: %w", err)
	}

	r.nextOffset = r.prevOffset
	r.prevOffset = int64(prevOffset)

	return blk, bytesRead, nil
}

func (r *Reader) readSignedBlock() (block *eos.SignedBlock, bytesRead []byte, err error) {
	cnt := make([]byte, 1000000)

	// prevPos, err := r.fl.Seek(offset, os.SEEK_SET)
	// if err != nil {
	// 	return
	// }

	_, err = r.fl.Read(cnt)
	if err != nil {
		return
	}

	d := eos.NewDecoder(cnt)

	if err = d.Decode(&block); err != nil {
		err = fmt.Errorf("decoding signed block: %w", err)
		return
	}

	// jsonStr, err := json.Marshal(block)
	// if err != nil {
	// 	return err
	// }

	// fmt.Println(string(jsonStr))

	bytesRead = cnt[:d.LastPos()]

	// r.nextOffset = prevPos + int64(d.LastPos()) + 8
	// r.prevOffset = prevPos
	return
}
