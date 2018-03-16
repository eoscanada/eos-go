package eosapi

import (
	"encoding/binary"
	"io"

	"github.com/lunixbochs/struc"
)

// from: https://github.com/lunixbochs/struc/issues/52
type byteWrapper struct {
	io.Reader
}

func (w *byteWrapper) ReadByte() (byte, error) {
	var b [1]byte
	_, err := w.Read(b[:])
	return b[0], err
}

// Varint used to pack and unpack structs with a Varint length field.
type Varint uint64

func (v *Varint) Pack(p []byte, order binary.ByteOrder) (int, error) {
	return binary.PutUvarint(p, uint64(*v)), nil
}
func (v *Varint) Unpack(r io.Reader, length int, opt *struc.Options) error {
	n, err := binary.ReadUvarint(&byteWrapper{r})
	*v = Varint(n)
	return err
}
func (v *Varint) SizeOf() int {
	var buf [8]byte
	return binary.PutUvarint(buf[:], uint64(*v))
}

// Varint used to pack and unpack structs with a Varint length field.
type Varint32 uint32

func (v Varint32) Pack(p []byte, order binary.ByteOrder) (int, error) {
	return binary.PutUvarint(p, uint64(v)), nil
}
func (v *Varint32) Unpack(r io.Reader, length int, opt *struc.Options) error {
	n, err := binary.ReadUvarint(&byteWrapper{r})
	*v = Varint32(n)
	return err
}
func (v Varint32) SizeOf() int {
	var buf [4]byte
	return binary.PutUvarint(buf[:], uint64(v))
}

// // Varint used to pack and unpack structs with a Varint length field.
// type Varint uint64

// func (v *Varint) Pack(p []byte, opt *struc.Options) (int, error) {
// 	return binary.PutUvarint(p, uint64(*v)), nil
// }
// func (v *Varint) Unpack(r io.Reader, length int, opt *struc.Options) error {
// 	n, err := binary.ReadUvarint(&byteWrapper{r})
// 	*v = Varint(n)
// 	return err
// }
// func (v *Varint) Size(opt *struc.Options) int {
// 	var buf [8]byte
// 	return binary.PutUvarint(buf[:], uint64(*v))
// }
// func (v *Varint) String() string {
// 	return strconv.FormatUint(uint64(*v), 10)
// }

// // Varint used to pack and unpack structs with a Varint length field.
// type Varint32 uint32

// func (v Varint32) Pack(p []byte, opt *struc.Options) (int, error) {
// 	return binary.PutUvarint(p, uint64(v)), nil
// }
// func (v *Varint32) Unpack(r io.Reader, length int, opt *struc.Options) error {
// 	n, err := binary.ReadUvarint(&byteWrapper{r})
// 	*v = Varint32(n)
// 	return err
// }
// func (v Varint32) Size(opt *struc.Options) int {
// 	var buf [4]byte
// 	return binary.PutUvarint(buf[:], uint64(v))
// }
// func (v Varint32) String() string {
// 	return strconv.FormatUint(uint64(v), 10)
// }
