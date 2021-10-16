package snapshot

import (
	"errors"
	"fmt"
	"io"
)

func readByteArray(buf io.Reader) ([]byte, error) {
	valueSize, err := readUvarint(buf)
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

func readUvarint(r io.Reader) (uint64, error) {
	buf := make([]byte, 1)
	var x uint64
	var s uint
	for i := 0; ; i++ {
		reads, err := r.Read(buf)
		if err != nil {
			return x, err
		}
		if reads != 1 {
			return x, fmt.Errorf("uvarint failed reading next byte")
		}
		b := buf[0]
		if b < 0x80 {
			if i > 9 || i == 9 && b > 1 {
				return x, errors.New("binary: varint overflows a 64-bit interger")
			}
			return x | uint64(b)<<s, nil
		}
		x |= uint64(b&0x7f) << s
		s += 7
	}
}

func readZeroTerminatedString(r io.Reader) (out string, err error) {
	b := make([]byte, 1)
	for i := 0; ; i++ {
		if i > 10000 {
			return out, fmt.Errorf("unsupported string over 10k")
		}
		readz, err := r.Read(b)
		if err != nil {
			return out, err
		}
		if readz != 1 {
			return out, fmt.Errorf("read 0, expected 1")
		}
		if b[0] == 0x00 {
			return out, nil
		}
		out = out + string(b[0])
	}
}
