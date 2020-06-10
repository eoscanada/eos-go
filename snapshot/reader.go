package snapshot

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
)

type Reader struct {
	Header *Header

	filename   string
	fl         *os.File
	buf        bufio.Reader
	nextOffset uint64
}

func NewReader(filename string) (r *Reader, err error) {
	r = &Reader{filename: filename}
	r.fl, err = os.Open(filename)
	if err != nil {
		return nil, err
	}

	h, err := r.readHeader()
	if err != nil {
		return nil, err
	}

	beginOffset, err := r.fl.Seek(0, os.SEEK_CUR)
	if err != nil {
		return nil, err
	}

	r.Header = h
	r.nextOffset = uint64(beginOffset)
	//fmt.Println("Beign offset", r.nextOffset)

	return
}

var magicNumber = []byte{0x50, 0x05, 0x51, 0x30}

func (r *Reader) readHeader() (*Header, error) {
	buf := make([]byte, 8)
	if _, err := r.fl.Read(buf); err != nil {
		return nil, err
	}

	if bytes.Compare(buf[:4], magicNumber) != 0 {
		return nil, fmt.Errorf("invalid magic number (first 4 bytes): %v, expected %v", buf[:4], magicNumber)
	}

	h := &Header{
		Version: binary.LittleEndian.Uint32(buf[4:8]),
	}

	return h, nil
}

// Next retrieves the next section.
func (r *Reader) Next() (*Section, error) {
	beginOffset, err := r.fl.Seek(int64(r.nextOffset), os.SEEK_SET)
	if err != nil {
		return nil, err
	}

	vals := make([]byte, 16)
	bytesRead, err := r.fl.Read(vals)
	if err != nil {
		return nil, err
	}

	// end marker
	if bytesRead == 8 && bytes.Compare(vals[:8], []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}) == 0 {
		return nil, io.EOF
	}

	sectionSize := binary.LittleEndian.Uint64(vals[:8])
	rowCount := binary.LittleEndian.Uint64(vals[8:16])

	buf := bufio.NewReaderSize(r.fl, int(sectionSize))
	str, err := buf.ReadString(0x00)
	if err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("EOF while reading string section (partial: %s)", str)
		}
		return nil, err
	}

	r.nextOffset = uint64(beginOffset) + sectionSize + 8 // well well, sectionSize includes the rowCount I guess?

	return &Section{
		Name:       strings.TrimRight(str, string([]byte{0x00})),
		Size:       sectionSize,
		RowCount:   rowCount,
		BufferSize: sectionSize - uint64(len(str)) - 8,
		Buffer:     buf,
	}, nil
}

func (r *Reader) Close() error {
	return r.fl.Close()
}
