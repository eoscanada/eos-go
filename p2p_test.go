package eos

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {

	hexString := `09000000050100000019000000`
	decoded, err := hex.DecodeString(hexString)
	if err != nil {
		t.Error(err)
	}

	var s P2PMessage

	assert.NoError(t, UnmarshalBinary(decoded, &s))
	assert.Equal(t, uint32(9), s.Length)
	assert.Equal(t, byte(5), s.Type)
	assert.Equal(t, []byte{0x1, 0x0, 0x0, 0x0, 0x19, 0x0, 0x0, 0x0}, s.Payload)

}
