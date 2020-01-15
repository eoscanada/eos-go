package ecc

import (
	"bytes"
	"fmt"

	"github.com/eoscanada/eos-go/btcsuite/btcutil/base58"
)

var keyMaterialDecoders = map[CurveID]keyMaterialDecoder{
	CurveR1: keyMaterialDecoderFunc(base58.Decode),
	CurveK1: keyMaterialDecoderFunc(base58.Decode),
	CurveWA: keyMaterialDecoderFunc(base58.DecodeVarSize),
}

// decodeKeyMaterial decodes a string that was encoded with CheckEncode and verifies the checksum.
func decodeKeyMaterial(input string, curve CurveID) (result []byte, err error) {
	decoder := keyMaterialDecoders[curve]
	if decoder == nil {
		decoder = keyMaterialDecoderFunc(base58.Decode)
	}

	decoded := decoder.Decode(input)
	if len(decoded) < 5 {
		return nil, fmt.Errorf("invalid format")
	}
	var cksum [4]byte
	copy(cksum[:], decoded[len(decoded)-4:])
	///// WARN: ok the ripemd160checksum should include the prefix in CERTAIN situations,
	// like when we imported the PubKey without a prefix ?! tied to the string representation
	// or something ? weird.. checksum shouldn't change based on the string reprsentation.
	if bytes.Compare(ripemd160checksum(decoded[:len(decoded)-4], curve), cksum[:]) != 0 {
		return nil, fmt.Errorf("invalid checksum")
	}
	// perhaps bitcoin has a leading net ID / version, but EOS doesn't
	payload := decoded[:len(decoded)-4]
	result = append(result, payload...)
	return
}

type keyMaterialDecoder interface {
	Decode(input string) []byte
}

type keyMaterialDecoderFunc func(input string) []byte

func (f keyMaterialDecoderFunc) Decode(input string) []byte {
	return f(input)
}
