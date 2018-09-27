package ecc

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/eoscanada/eos-go/btcsuite/btcutil/base58"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func FixmeTestSignatureSerialization(t *testing.T) {
	privkey, err := NewPrivateKey("5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3")
	require.NoError(t, err)

	payload := []byte("payload")
	sig, err := privkey.Sign(sigDigest(make([]byte, 32, 32), payload))
	require.NoError(t, err)
	assert.Equal(t, `SIG_K1_K2JjfxmYpoVwCKkohDiQPcepeyetSWMgQPjx3zqagzao5NeQhnW4JQ2qwxd4txU7dR5TdS6PnP75vmMs5qSXzjphqUZz6N`, sig.String()) // not checked after..
	assert.True(t, isCanonical([]byte(sig.Content)))
}

func FixmeTestSignatureCanonical(t *testing.T) {
	privkey, err := NewPrivateKey("5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3")
	require.NoError(t, err)

	fmt.Println("Start")
	payload := []byte("payload1") // doesn't fail
	sig, err := privkey.Sign(sigDigest(make([]byte, 32, 32), payload))
	fmt.Println("Signed")
	require.NoError(t, err)
	fmt.Println("MAM", sig.String())
	assert.True(t, isCanonical([]byte(sig.Content)))
	fmt.Println("End")

	fmt.Println("Start")
	payload = []byte("payload6") // fails
	sig, err = privkey.Sign(sigDigest(make([]byte, 32, 32), payload))
	fmt.Println("Signed")
	require.NoError(t, err)
	fmt.Println("MAM1", sig.String())
	assert.True(t, isCanonical([]byte(sig.Content)))
	fmt.Println("End")
}

func FixmeTestSignatureMarshalUnmarshal(t *testing.T) {
	fromEOSIOC := "SIG_K1_K5yY5ehsnDMc6xcRhsLYzFuZGUaKwb4hc8oLmP5HA1EhU42NRo3ygx3zvLRJ1nkw1NA5nCSegwcYkSfkZBQBzqMDsCGnNK"
	sig, err := NewSignature(fromEOSIOC)
	require.NoError(t, err)
	assert.Equal(t, fromEOSIOC, sig.String())
	assert.True(t, isCanonical([]byte(sig.Content)))
}
func FixmeTestSignatureMarshalUnmarshal_bilc(t *testing.T) {
	fromEOSIOC := "SIG_K1_Jy9G6SgmGSjAbu7n82veUiqV8LFFL6wqr9G26H37dy1WExUj9kYwS17X3ffT5W9M51HkpKF4xQ6MoFCCMxBEHbk64dgbMg"
	sig, err := NewSignature(fromEOSIOC)
	require.NoError(t, err)
	assert.Equal(t, fromEOSIOC, sig.String())
	assert.True(t, isCanonical([]byte(sig.Content)))
}

func isCanonical(compactSig []byte) bool {
	// !(c.data[1] & 0x80)
	// && !(c.data[1] == 0 && !(c.data[2] & 0x80))
	// && !(c.data[33] & 0x80)
	// && !(c.data[33] == 0 && !(c.data[34] & 0x80));

	d := compactSig
	t1 := (d[1] & 0x80) == 0
	t2 := !(d[1] == 0 && ((d[2] & 0x80) == 0))
	t3 := (d[33] & 0x80) == 0
	t4 := !(d[33] == 0 && ((d[34] & 0x80) == 0))
	return t1 && t2 && t3 && t4
}

func FixmeTestSignaturePublicKeyExtraction(t *testing.T) {
	// was signed with EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV
	fromEOSIOC := "SIG_K1_K5yY5ehsnDMc6xcRhsLYzFuZGUaKwb4hc8oLmP5HA1EhU42NRo3ygx3zvLRJ1nkw1NA5nCSegwcYkSfkZBQBzqMDsCGnNK"
	sig, err := NewSignature(fromEOSIOC)
	require.NoError(t, err)

	payload, err := hex.DecodeString("20d8af5a0000b32bcc0e37eb0000000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea305500001059b1abe93101000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed32320100")
	require.NoError(t, err)

	pubKey, err := sig.PublicKey(sigDigest(make([]byte, 32, 32), payload))
	require.NoError(t, err)

	// Ok, we'd need to find values where we know the signature is valid, and comes from the given key.
	assert.Equal(t, "PUB_K1_6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV", pubKey.String())
}

// OKAY I THINK THERE's a problem with the internal representation of
// a signature.. it works in an of itself, but it is not compatible
// with the main `nodeos` software.
//
// FIXME: We need to fix that for this library to be able to sign
// transactions and push them to the network without relying on an
// external wallet, or eosjs-ecc or something..
func TestSignaturePublicKeyExtractionSecond(t *testing.T) {
	// this was transaction be72ed8f391277c7792caec781b70f3e97766920c1f3844fdbb82b7db5f0381e
	// was signed with EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV
	fromEOSIOC := "SIG_K1_KkwLhwDoRF8gpGFbcUKiaPdeeKo6U7eDuXQw9szMiNE4K4cFe17sffk6hmy3mWf1ogtzd5J5kvnvFD3Lq5cF6VyYb3KsGy"
	sig, err := NewSignature(fromEOSIOC)
	require.NoError(t, err)

	payload, err := hex.DecodeString("30d3b35a0000be0194c22fe70000000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea305500000059b1abe93101000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed32320100")
	require.NoError(t, err)

	pubKey, err := sig.PublicKey(sigDigest(make([]byte, 32, 32), payload))
	require.NoError(t, err)

	// Ok, we'd need to find values where we know the signature is valid, and comes from the given key.
	assert.Equal(t, "EOS5DguRMaGh72NvbVX5LKHTb5cvbRmAxgrm9i2NNPKv5TC7FadXs", pubKey.String())
}

func FixmeTestEOSIOCSigningComparison(t *testing.T) {
	// try with: ec sign -k 5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3 '{"expiration":"2018-03-21T23:02:32","region":0,"ref_block_num":2156,"ref_block_prefix":1532582828,"packed_bandwidth_words":0,"context_free_cpu_bandwidth":0,"context_free_actions":[],"actions":[],"signatures":[],"context_free_data":[]}'
	wif := "5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3" // corresponds to: EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV
	privKey, err := NewPrivateKey(wif)
	require.NoError(t, err)

	chainID, err := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
	require.NoError(t, err)

	payload, err := hex.DecodeString("88e4b25a00006c08ac5b595b000000000000") // without signed transaction bytes
	require.NoError(t, err)

	digest := sigDigest(chainID, payload)

	sig, err := privKey.Sign(digest)
	require.NoError(t, err)

	fromEOSIOC := "SIG_K1_K2WBNtiTY8o4mqFSz7HPnjkiT9JhUYGFa81RrzaXr3aWRF1F8qwVfutJXroqiL35ZiHTcvn8gPWGYJDwnKZTCcbAGJy4i9"
	assert.Equal(t, fromEOSIOC, sig.String())
}

func FixmeTestNodeosSignatureComparison(t *testing.T) {
	wif := "5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3" // corresponds to: EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV
	privKey, err := NewPrivateKey(wif)
	require.NoError(t, err)

	// produce with `cleos create account eosio abourget EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV
	// transaction:
	// chainID + 30d3b35a0000be0194c22fe70000000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea305500000059b1abe93101000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed323201000000
	// hashes to:
	digest, _ := hex.DecodeString("a744a49dd60badd5e7073e7287d53e184914242e94ef309d2694e954077dcb27")

	sig, err := privKey.Sign(digest)
	require.NoError(t, err)

	// from that tx:
	fromEOSIOCTx := "SIG_K1_K9JDNpqcgUin9i2PtsV6QbLG8QGzYPN8kqVicJ63CgHBiwq9q27qykaerbNh8kD6baLFWcuKyTmVUwFRF6myjqFQcHvRXf"
	assert.Equal(t, fromEOSIOCTx, sig.String())

	// decode
	fmt.Println("From EOSIO sig:", hex.EncodeToString(base58.Decode(fromEOSIOCTx[3:])))
	fmt.Println("From GO sig:", hex.EncodeToString(base58.Decode(sig.String()[3:])))
}

func FixmeTestSignatureUnmarshalChecksum(t *testing.T) {
	fromEOSIOC := "SIG_K1_K5yY5ehsnDMc6xcRhsLYzFuZGUaKwb4hc8oLmP5HA1EhU42NRo3ygx3zvLRJ1nkw1NA5nCSegwcYkSfkZBQBzqMDsCGnZZ" // simply checked the last 2 bytes
	_, err := NewSignature(fromEOSIOC)
	require.Equal(t, "signature checksum failed, found 02c9bc70 expected 02c9befc", err.Error())
}

func sigDigest(chainID, payload []byte) []byte {
	h := sha256.New()
	_, _ = h.Write(chainID)
	_, _ = h.Write(payload)
	return h.Sum(nil)
}
