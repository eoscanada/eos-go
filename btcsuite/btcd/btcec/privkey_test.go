// Copyright (c) 2013-2016 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package btcec

import (
	"bytes"
	"crypto/elliptic"
	"encoding/hex"
	"github.com/eoscanada/eos-go/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
	"testing"
)

func TestPrivKeys(t *testing.T) {
	tests := []struct {
		name string
		key  []byte
	}{
		{
			name: "check curve",
			key: []byte{
				0xea, 0xf0, 0x2c, 0xa3, 0x48, 0xc5, 0x24, 0xe6,
				0x39, 0x26, 0x55, 0xba, 0x4d, 0x29, 0x60, 0x3c,
				0xd1, 0xa7, 0x34, 0x7d, 0x9d, 0x65, 0xcf, 0xe9,
				0x3c, 0xe1, 0xeb, 0xff, 0xdc, 0xa2, 0x26, 0x94,
			},
		},
	}

	for _, test := range tests {
		priv, pub := PrivKeyFromBytes(S256(), test.key)

		_, err := ParsePubKey(pub.SerializeUncompressed(), S256())
		if err != nil {
			t.Errorf("%s privkey: %v", test.name, err)
			continue
		}

		hash := []byte{0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9}
		sig, err := priv.Sign(hash)
		if err != nil {
			t.Errorf("%s could not sign: %v", test.name, err)
			continue
		}

		if !sig.Verify(hash, pub) {
			t.Errorf("%s could not verify: %v", test.name, err)
			continue
		}

		serializedKey := priv.Serialize()
		if !bytes.Equal(serializedKey, test.key) {
			t.Errorf("%s unexpected serialized bytes - got: %x, "+
				"want: %x", test.name, serializedKey, test.key)
		}
	}
}

func sign2String(signature []byte, curve elliptic.Curve) string {
	h := ripemd160.New()
	_, _ = h.Write(signature)
	_, _ = h.Write([]byte("K1"))
	sum := h.Sum(nil)

	buf := append(signature[:], sum[:4]...)
	return "SIG_K1_" + base58.Encode(buf)
}

/*
   Test cases from eosjs (src/signature.test.js)

   const digest = '6cb75bc5a46a7fdb64b92efefca01ed7b060ab5e0d625226e8efbc0980c3ddc1'
   const wif = '5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3'
   const expectedSig = 'SIG_K1_Kk1yUXAG2Cfo2qvWuJiyvaGdwZBQ1HzSf4EZ9arUTWBL4kTngLM1GSUU59bJUVAqwJ886CNQMcR7mmx323gjQGvhEU8WpX'

   const digest = Buffer.alloc(32, 0)
   const wif = '5HxQKWDznancXZXm7Gr2guadK7BhK9Zs8ejDhfA9oEBM89ZaAru'
   const expectedSig = 'SIG_K1_Jz9d1rKmMV51EY6dnU3pNaDiLvGTeVdxDZGvJEfAkdcwzs97gNg5yYPhPSdEg33Jyp5736Tnnzccf1p6h6vedXpHSUBio1'
*/

func TestEosSignature(t *testing.T) {

	cases := []struct {
		name string
		exec func(t *testing.T)
	}{
		{
			name: "SUCCESS: Deterministic Sign",
			exec: func(t *testing.T) {
				//Given
				keyValue := []byte{
					0xd2, 0x65, 0x3f, 0xf7, 0xcb, 0xb2, 0xd8, 0xff, 0x12, 0x9a, 0xc2, 0x7e, 0xf5, 0x78, 0x1c, 0xe6,
					0x8b, 0x25, 0x58, 0xc4, 0x1a, 0x74, 0xaf, 0x1f, 0x2d, 0xdc, 0xa6, 0x35, 0xcb, 0xee, 0xf0, 0x7d,
				}
				key, _ := PrivKeyFromBytes(S256(), keyValue)

				digest, err := hex.DecodeString("6cb75bc5a46a7fdb64b92efefca01ed7b060ab5e0d625226e8efbc0980c3ddc1")
				if err != nil {
					t.Fatalf("DecodeString: %s", err.Error())
				}
				//When
				signature, err := key.SignCanonical(S256(), digest)
				if err != nil {
					t.Fatalf("Sign: %s", err.Error())
				}

				signatureExpected := "SIG_K1_Kk1yUXAG2Cfo2qvWuJiyvaGdwZBQ1HzSf4EZ9arUTWBL4kTngLM1GSUU59bJUVAqwJ886CNQMcR7mmx323gjQGvhEU8WpX"
				//Then
				if signatureExpected != sign2String(signature, S256()) {
					t.Fatalf("unexpected signature case1")
				}
			},
		},
		{
			name: "SUCCESS: Deterministic canonical2",
			exec: func(t *testing.T) {
				//Given
				keyValue := []byte{
					0x12, 0x6d, 0xc4, 0x3d, 0x71, 0x94, 0x51, 0xc5, 0x8c, 0xae, 0x72, 0x8e, 0xb7, 0x0c, 0x6a, 0xdf,
					0x16, 0x17, 0xb6, 0x23, 0xcf, 0xbf, 0x88, 0x75, 0xba, 0x6d, 0x13, 0x33, 0x2e, 0x6d, 0x92, 0x0d,
				}
				key, _ := PrivKeyFromBytes(S256(), keyValue)

				digest, err := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
				if err != nil {
					t.Fatalf("DecodeString: %s", err.Error())
				}
				//When
				signature, err := key.SignCanonical(S256(), digest)
				if err != nil {
					t.Fatalf("Sign: %s", err.Error())
				}

				signatureExpected := "SIG_K1_Jz9d1rKmMV51EY6dnU3pNaDiLvGTeVdxDZGvJEfAkdcwzs97gNg5yYPhPSdEg33Jyp5736Tnnzccf1p6h6vedXpHSUBio1"
				//Then
				if signatureExpected != sign2String(signature, S256()) {
					t.Fatalf("unexpected signature case2")
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.exec(t)
		})
	}

}
