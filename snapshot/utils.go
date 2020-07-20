package snapshot

import (
	"bufio"
	"encoding/binary"
	"fmt"
)

func readByteArray(buf *bufio.Reader) ([]byte, error) {
	valueSize, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, err
	}

	val := make([]byte, valueSize)

	written, err := buf.Read(val)
	if err != nil {
		return nil, err
	}
	if uint64(written) != valueSize {
		return nil, fmt.Errorf("inconsistent read, expected %d bytes, read %d", valueSize, written)
	}

	return val, nil
}
