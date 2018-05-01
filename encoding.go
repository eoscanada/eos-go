package eos

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Decoder implements the EOS unpacking, similar to FC_BUFFER
type Decoder struct {
	data               []byte
	pos                int
	decodeP2PMessage   bool
	decodeTransactions bool
	decodeActions      bool

	actionMap    map[AccountName]map[ActionName]interface{}
	actionABIMap map[AccountName]map[ActionName]ABIDef

	//lastSeenAction ActionName
}

func NewDecoder(data []byte) *Decoder {
	return &Decoder{
		data:               data,
		decodeP2PMessage:   true,
		decodeTransactions: true,
		decodeActions:      true,
	}
}

func (d *Decoder) SetActionMap(actions map[AccountName]map[ActionName]interface{}) {
	d.actionMap = actions
}

func (d *Decoder) SetDecodeTransactions(enabled bool) {
	d.decodeTransactions = enabled
}

func (d *Decoder) SetDecodeP2PMessage(enabled bool) {
	d.decodeP2PMessage = enabled
}

func (d *Decoder) SetDecodeActions(enabled bool) {
	d.decodeActions = enabled
}

func (d *Decoder) Decode(v interface{}) (err error) {
	// TODO: reflect, get a new instance of the type or someth'n.. similar to  json.Unmarshal
	//based on the type of v

	switch dst := v.(type) {
	case *P2PMessageEnvelope:
		if d.remaining() < 6 {
			return 0, fmt.Errorf("%d bytes remaining, reading P2PMEssageEnvelope requires at least 6", d.remaining())
		}
		dst.Length = d.readUint32()
		dst.Type = P2PMessageType(d.readByte())
		dst.Payload, err = d.readByteArray()
		if err != nil {
			return
		}
		if d.decodeP2PMessage {
			subD := NewDecoder(dst.Payload)
			switch dst.Type {
				dst.Message = subD.readThisType(&SignedBlockTransaction{})
			}
		}

	case *P2PMessage
	case *PackedTransaction:
	case *Transaction:
	case *Action:
	case string:
	case HexBytes:
	case interface{}:4
		if d.p2pType
	default:
		return fmt.Errorf("unsupported type %T", v)
	}
	return nil
}

func (d *Decoder) readByteArray() (out []byte, err error) {
	varlen, read := binary.Varint(d.data[d.pos])
	if read == 0 {
		return nil, fmt.Errorf("varint: buffer too small")
	}
	if read < 0 {
		return nil, fmt.Errorf("varint: 64 bits overflow")
	}

	d.pos += read

	// TODO: do bounds check on the varint itself.. was it successful? did you miss bytes ?

	if len(d.data) < d.pos+varlen {
		return nil, fmt.Errorf("byte array: varlen=%d, missing %d bytes", varlen, d.pos+varlen-len(d.data))
	}

	out = d.data[d.pos : d.pos+varlen]
	d.pos += int(varlen)

	return
}

func (d *Decoder) readUint16() (out uint16) {
	out = binary.LittleEndian.Uint16(d.data[d.pos:])
	d.pos += 2
	return
}

func (d *Decoder) readByte() (out byte) {
	out = binary.LittleEndian.Byte(d.data[d.pos:])
	d.pos++
	return
}

func (d *Decoder) readUint32() (out uint32) {
	out = binary.LittleEndian.Uint32(d.data[d.pos:])
	d.pos += 4
	return
}

func (d *Decoder) remaining() int {
	return len(d.data) - d.pos
}

// Encoder implements the EOS packing, similar to FC_BUFFER
type Encoder struct {
	output io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		output: w,
	}
}
