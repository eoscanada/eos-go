package snapshot

import "bufio"

type Section struct {
	Name       string
	Size       uint64 // This includes the section name and row count
	BufferSize uint64 // This represents the bytes that are following the section header
	RowCount   uint64 // This is a count of rows packed in `Buffer`
	Buffer     *bufio.Reader
}

// Next reads the next row
func (s *Section) Next() ([]byte, error) {
	return nil, nil
}
