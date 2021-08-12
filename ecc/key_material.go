package ecc

import (
	"bytes"
	"fmt"

	"github.com/eoscanada/eos-go/btcsuite/btcutil/base58"
)

type keyMaterialDecoder func(input string) []byte

var keyMaterialDecoders = map[CurveID]keyMaterialDecoder{
	CurveR1: base58.Decode,
	CurveK1: base58.Decode,
	CurveWA: base58.DecodeVarSize,
}

// decodePublicKeyMaterial decodes a public key checksumed string and verifies its checksum.
func decodePublicKeyMaterial(in string, curve CurveID) (result []byte, err error) {
	return decodeKeyMaterial("public key", in, curve, ripemd160checksum)
}

// decodeSignatureMaterial decodes a signature checksumed string and verifies its checksum.
func decodeSignatureMaterial(in string, curve CurveID) (result []byte, err error) {
	return decodeKeyMaterial("signature", in, curve, ripemd160checksumHashCurve)
}

// decodeKeyMaterial decodes a string that was encoded with CheckEncode and verifies the checksum.
func decodeKeyMaterial(tag string, input string, curve CurveID, hasher func([]byte, CurveID) []byte) (result []byte, err error) {
	decoder := keyMaterialDecoders[curve]
	if decoder == nil {
		return nil, fmt.Errorf("%s invalid curve %s (%d): no decoder", tag, curve, uint8(curve))
	}

	decoded := decoder(input)
	if len(decoded) < 5 {
		return nil, fmt.Errorf("%s invalid format, expected at least 5 bytes, got %d", tag, len(decoded))
	}

	checksumOffset := len(decoded) - 4
	payload := decoded[:checksumOffset]
	checksum := decoded[checksumOffset:]
	verifyChecksum := hasher(payload, curve)
	if !bytes.Equal(verifyChecksum, checksum) {
		return nil, fmt.Errorf("%s checksum failed, found %x but expected %x", tag, verifyChecksum, checksum)
	}

	return payload, nil
}
