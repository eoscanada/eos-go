package eos

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecoder_P2PMessageEnvelope(t *testing.T) {

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)

	msg := &P2PMessageEnvelope{
		Length:  4,
		Type:    SignedTransactionMessageType,
		Payload: []byte{1, 2, 3},
	}

	err := enc.Encode(msg)
	assert.NoError(t, err)
	d := NewDecoder(enc.data)

	var decoded P2PMessageEnvelope

	err = d.Decode(&decoded)
	assert.NoError(t, err)
	assert.Equal(t, uint32(4), decoded.Length)
	assert.Equal(t, SignedTransactionMessageType, decoded.Type)
	assert.Equal(t, []byte{1, 2, 3}, decoded.Payload)
}
