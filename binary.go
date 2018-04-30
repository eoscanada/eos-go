package eos

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
)

// From github.com/alecthomas/binary

var print = func(s string) {
	//	fmt.Print(s)
}
var println = func(s string) {
	//	print(fmt.Sprintf("%s\n", s))
}

var (
	LittleEndian  = binary.LittleEndian
	BigEndian     = binary.BigEndian
	DefaultEndian = LittleEndian
)

func MarshalBinary(v interface{}) ([]byte, error) {
	b := &bytes.Buffer{}
	if err := NewEncoder(b).Encode(v); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func UnmarshalBinary(b []byte, v interface{}) error {
	return NewDecoder(bytes.NewReader(b)).Decode(v)
}

func UnmarshalBinaryWithAction(b []byte, v interface{}, act Action) error {
	d := NewDecoder(bytes.NewReader(b))
	d.lastAction = act
	return d.Decode(v)
}

type Encoder struct {
	Order  binary.ByteOrder
	w      io.Writer
	buf    []byte
	strict bool
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		Order: DefaultEndian,
		w:     w,
		buf:   make([]byte, 8),
	}
}

// NewStrictEncoder creates an encoder similar to NewEncoder, however
// if this encoder attempts to encode a struct and the struct has no encodable
// fields an error is returned whereas the encoder returned from NewEncoder
// will simply not write anything to `w`.
func NewStrictEncoder(w io.Writer) *Encoder {
	e := NewEncoder(w)
	e.strict = true
	return e
}

func (e *Encoder) writeVarint(v int) error {
	l := binary.PutUvarint(e.buf, uint64(v))
	_, err := e.w.Write(e.buf[:l])
	return err
}

func (b *Encoder) Encode(v interface{}) (err error) {
	if i, ok := v.(OptionalBinaryMarshaler); ok {
		if i.OptionalBinaryMarshalerPresent() {
			b.w.Write([]byte{0x01})
		} else {
			b.w.Write([]byte{0x01})
			return nil
		}
	}

	switch cv := v.(type) {
	case encoding.BinaryMarshaler:
		buf, err := cv.MarshalBinary()
		if err != nil {
			return err
		}
		// let the Marshallers write their own prefix.. we already support
		// handling unmarshallers..
		// if err = b.writeVarint(len(buf)); err != nil {
		// 	return err
		// }
		_, err = b.w.Write(buf)

	case []byte: // fast-path byte arrays
		if err = b.writeVarint(len(cv)); err != nil {
			return
		}
		_, err = b.w.Write(cv)

	default:
		rv := reflect.Indirect(reflect.ValueOf(v))
		t := rv.Type()
		switch t.Kind() {
		case reflect.Ptr:
			return b.Encode(rv.Elem().Interface())

		case reflect.Array:
			l := t.Len()
			for i := 0; i < l; i++ {
				if err = b.Encode(rv.Index(i).Interface()); err != nil {
					return
				}
			}

		case reflect.Slice:
			l := rv.Len()
			if err = b.writeVarint(l); err != nil {
				return
			}
			for i := 0; i < l; i++ {
				if err = b.Encode(rv.Index(i).Interface()); err != nil {
					return
				}
			}

		case reflect.Struct:
			l := rv.NumField()
			n := 0
			for i := 0; i < l; i++ {
				//fmt.Println("Bin -> ", t.Field(i).Name)
				if v := rv.Field(i); t.Field(i).Name != "_" {
					if v.CanInterface() {
						iface := v.Interface()
						if iface != nil {
							if err = b.Encode(iface); err != nil {
								return
							}
						}
					}
					n++
				}
			}
			if b.strict && n == 0 {
				return fmt.Errorf("binary: struct had no encodable fields")
			}

		case reflect.Map:
			l := rv.Len()
			if err = b.writeVarint(l); err != nil {
				return
			}
			for _, key := range rv.MapKeys() {
				value := rv.MapIndex(key)
				if err = b.Encode(key.Interface()); err != nil {
					return err
				}
				if err = b.Encode(value.Interface()); err != nil {
					return err
				}
			}

		case reflect.String:
			if err = b.writeVarint(rv.Len()); err != nil {
				return
			}
			_, err = b.w.Write([]byte(rv.String()))

		case reflect.Bool:
			var out byte
			if rv.Bool() {
				out = 1
			}
			err = binary.Write(b.w, b.Order, out)

		case reflect.Int:
			err = binary.Write(b.w, b.Order, int64(rv.Int()))

		case reflect.Uint:
			err = binary.Write(b.w, b.Order, int64(rv.Uint()))

		case reflect.Int8, reflect.Uint8, reflect.Int16, reflect.Uint16,
			reflect.Int32, reflect.Uint32, reflect.Int64, reflect.Uint64,
			reflect.Float32, reflect.Float64,
			reflect.Complex64, reflect.Complex128:
			err = binary.Write(b.w, b.Order, v)

		default:
			return errors.New("binary: unsupported type " + t.String())
		}
	}
	return
}

type ByteReader struct {
	io.Reader
}

func (b *ByteReader) ReadByte() (byte, error) {
	var buf [1]byte
	if _, err := io.ReadFull(b, buf[:]); err != nil {
		return 0, err
	}
	return buf[0], nil
}

func (b *ByteReader) Read(p []byte) (int, error) {
	if cap(p) == 0 {
		return 0, nil
	}
	return b.Reader.Read(p)
}

type Decoder struct {
	Order      binary.ByteOrder
	r          *ByteReader
	lastAction Action
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		Order: DefaultEndian,
		r:     &ByteReader{r},
	}
}

// UnmarshalBinarySizer determines the size of what we need to read to
// unmarshal this object. If not specified, and not a native type,
// falls back to a uvarint prefix.
type UnmarshalBinarySizer interface {
	UnmarshalBinarySize() int
}

// UnmarshalBinaryReader will read by-per-by (implementing a varint of
// some sort), and consume what it needs out of the stream.
type UnmarshalBinaryReader interface {
	UnmarshalBinaryRead(io.Reader) error
}

type UnmarshalBinaryWithCurrentAction interface {
	UnmarshalBinaryWithCurrentAction(data []byte, act Action) error
}

type OptionalBinaryMarshaler interface {
	OptionalBinaryMarshalerPresent() bool
}

func (d *Decoder) Decode(v interface{}) (err error) {

	//fmt.Printf("MAMA!!!: %#v %T\n", v, v)
	if i, ok := v.(UnmarshalBinaryReader); ok {
		return i.UnmarshalBinaryRead(d.r)
	}

	if _, ok := v.(OptionalBinaryMarshaler); ok {
		isPresent := make([]byte, 1, 1)
		_, err := d.r.Read(isPresent)
		if err != nil {
			return err
		}

		if isPresent[0] == 0 {
			return nil
		}
	}

	if i, ok := v.(UnmarshalBinaryWithCurrentAction); ok {
		var l uint64
		if l, err = binary.ReadUvarint(d.r); err != nil {
			return
		}
		buf := make([]byte, l)
		_, err = d.r.Read(buf)
		return i.UnmarshalBinaryWithCurrentAction(buf, d.lastAction)
	}

	if i, ok := v.(encoding.BinaryUnmarshaler); ok {
		//fmt.Println("BinaryUnmarshaler")
		// if we need, we'll implement an UnmarshalBinaryRead() that'll take precedence over this
		// and that will read byte-per-byte on its own..
		var l uint64
		if sizer, ok := v.(UnmarshalBinarySizer); ok {
			l = uint64(sizer.UnmarshalBinarySize())
		} else {
			if l, err = binary.ReadUvarint(d.r); err != nil {
				return
			}
		}
		buf := make([]byte, l)
		_, err = d.r.Read(buf)
		return i.UnmarshalBinary(buf)
	}

	rv := reflect.Indirect(reflect.ValueOf(v))
	if !rv.CanAddr() {
		return errors.New("binary: can only Decode to pointer type")
	}
	t := rv.Type()

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		newRV := reflect.New(t)
		rv.Set(newRV)
		rv = reflect.Indirect(newRV)
	}

	switch t.Kind() {
	case reflect.Array:
		//fmt.Println("Array")
		len := t.Len()
		for i := 0; i < int(len); i++ {
			if err = d.Decode(rv.Index(i).Addr().Interface()); err != nil {
				return
			}
		}

	case reflect.Slice:
		print("Reading Slice length ")
		var l uint64
		if l, err = binary.ReadUvarint(d.r); err != nil {
			return
		}
		println(fmt.Sprintf("Slice [%T] of length: %d", v, l))
		if t.Kind() == reflect.Slice {
			rv.Set(reflect.MakeSlice(t, int(l), int(l)))
		} else if int(l) != t.Len() {
			return fmt.Errorf("binary: encoded size %d != real size %d", l, t.Len())
		}
		for i := 0; i < int(l); i++ {
			if err = d.Decode(rv.Index(i).Addr().Interface()); err != nil {
				return
			}
		}

	case reflect.Struct:
		l := rv.NumField()
		for i := 0; i < l; i++ {
			if v := rv.Field(i); v.CanSet() && t.Field(i).Name != "_" {
				iface := v.Addr().Interface()
				println(fmt.Sprintf("Struct Field name: %s", t.Field(i).Name))
				if err = d.Decode(iface); err != nil {
					return
				}
			}
		}

	case reflect.Map:
		//fmt.Println("Map")
		var l uint64
		if l, err = binary.ReadUvarint(d.r); err != nil {
			return
		}
		kt := t.Key()
		vt := t.Elem()
		rv.Set(reflect.MakeMap(t))
		for i := 0; i < int(l); i++ {
			kv := reflect.Indirect(reflect.New(kt))
			if err = d.Decode(kv.Addr().Interface()); err != nil {
				return
			}
			vv := reflect.Indirect(reflect.New(vt))
			if err = d.Decode(vv.Addr().Interface()); err != nil {
				return
			}
			rv.SetMapIndex(kv, vv)
		}

	case reflect.String:
		//fmt.Println("String")
		var l uint64
		if l, err = binary.ReadUvarint(d.r); err != nil {
			return
		}
		buf := make([]byte, l)
		_, err = d.r.Read(buf)
		rv.SetString(string(buf))

	case reflect.Bool:
		//fmt.Println("Bool")
		var out byte
		err = binary.Read(d.r, d.Order, &out)
		rv.SetBool(out != 0)

	case reflect.Int:
		//fmt.Println("Int")
		var out int64
		err = binary.Read(d.r, d.Order, &out)
		rv.SetInt(out)

	case reflect.Uint:
		//fmt.Println("uInt")
		var out uint64
		err = binary.Read(d.r, d.Order, &out)
		rv.SetUint(out)

	case reflect.Int8, reflect.Uint8, reflect.Int16, reflect.Uint16,
		reflect.Int32, reflect.Uint32, reflect.Int64, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		//fmt.Println("Some funky ints")
		err = binary.Read(d.r, d.Order, v)

	default:
		return errors.New("binary: unsupported type " + t.String())
	}
	return
}

type LoggerReader struct {
	Reader io.Reader
}

func (l *LoggerReader) Read(p []byte) (n int, err error) {

	length := len(p)
	n, err = l.Reader.Read(p)
	if err == nil {
		fmt.Printf("\t\t[%d] data [%s]\n", length, hex.EncodeToString(p[:int(math.Min(float64(2000), float64(len(p))))]))
	}

	return
}
