package eos

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncoder_MapStringString(t *testing.T) {
	buffer := bytes.NewBuffer(nil)
	encoder := NewEncoder(buffer)
	data := map[string]string{
		"a": "1",
		"b": "2",
	}

	err := encoder.Encode(data)
	require.NoError(t, err)

	out := hex.EncodeToString(buffer.Bytes())

	// Sadly, we cannot do much for map not ordered here, so let's check that it's either one or the other
	expected1 := "020162013201610131"
	expected2 := "020161013101620132"

	if out != expected1 && out != expected2 {
		require.Fail(t, "encoded map is invalid", "must be either %q or %q, got %q", expected1, expected2, out)
	}
}
