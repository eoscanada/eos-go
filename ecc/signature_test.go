package ecc

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignatureSerialization(t *testing.T) {
	privkey, err := NewPrivateKey("5JFhynQnFBYNTPDA9TiKeE7TmujNYaExcbZi9bsRUjhVxwZF4Mt")
	require.NoError(t, err)

	payload, err := hex.DecodeString("89529cb031c69eccc92f3e8492393a8688bd3d071d7346677b6ff59d314d5121")
	require.NoError(t, err)
	digest := sigDigest(make([]byte, 32, 32), payload, nil)
	sig, err := privkey.Sign(digest)
	require.NoError(t, err)
	pubKey, err := sig.PublicKey(digest)
	require.NoError(t, err)
	assert.Equal(t, `EOS5jSQLpKBHLaMtuzkftnYE6bCMA5Jxso8f22uZyKj6cDEp32eSj`, pubKey.String()) // not checked after..
	assert.True(t, isCanonical([]byte(sig.Content)))
}

func TestSignatureCanonical(t *testing.T) {
	privkey, err := NewPrivateKey("5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3")
	require.NoError(t, err)

	//fmt.Println("Start")
	payload := []byte("payload1") // doesn't fail
	sig, err := privkey.Sign(sigDigest(make([]byte, 32, 32), payload, nil))
	//fmt.Println("Signed")
	require.NoError(t, err)
	//fmt.Println("MAM", sig.String())
	assert.True(t, isCanonical([]byte(sig.Content)))
	//fmt.Println("End")

	//fmt.Println("Start")
	payload = []byte("payload6") // fails
	sig, err = privkey.Sign(sigDigest(make([]byte, 32, 32), payload, nil))
	//fmt.Println("Signed")
	require.NoError(t, err)
	//fmt.Println("MAM1", sig.String())
	assert.True(t, isCanonical([]byte(sig.Content)))
	//fmt.Println("End")
}

func TestSignatureMarshalUnmarshal(t *testing.T) {
	cases := []struct {
		name          string
		signature     string
		testCanonical bool
	}{
		{
			name:          "K1",
			signature:     "SIG_K1_KVp1bPmzswSvbcZCMENXbawKFVXPyYrUeJNZ9ChgWdhxLd5K8WtRmCtFY5cqVFgxjCZH8CwdNkxM3HBZ7EXeJmzcK78mHA",
			testCanonical: true,
		},
		{
			name:      "R1",
			signature: "SIG_R1_KE33Ucjr5N3GR4ZosFh8KtGMytHHNtnmdUaSoMLJVXpVXoC8B9zfoXYrLiQJZqroe3LKciaP2uJT7Myqqoo4PZH7iSnso8",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			sig, err := NewSignature(c.signature)
			require.NoError(t, err)
			assert.Equal(t, c.signature, sig.String())
			if c.testCanonical {
				assert.True(t, isCanonical([]byte(sig.Content)))
			}
		})
	}
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

func TestSignaturePublicKeyExtraction(t *testing.T) {

	//SIG_R1_KE33Ucjr5N3GR4ZosFh8KtGMytHHNtnmdUaSoMLJVXpVXoC8B9zfoXYrLiQJZqroe3LKciaP2uJT7Myqqoo4PZH7iSnso8
	//PUB_R1_78rbUHSk87e7eCBoccgWUkhNTCZLYdvJzerDRHg6fxj2SQy6Xm

	cases := []struct {
		name                   string
		signature              string
		payload                string
		chainID                string
		expectedPubKey         string
		expectedSignatureError string
		expectedPubKeyError    string
	}{
		{
			name:           "K1",
			signature:      "SIG_K1_KW4qcHDh6ziqWELRAsFx42sgPuP3VfCpTKX4D5A3uZhFb3fzojTeGohja19g4EJa9Zv7SrGZ47H8apo1sNa2bwPvGwW2ba",
			payload:        "45e2ea5b22f87c6f74430000000001a0904b1822f330550040346aabab904b01a0904b1822f3305500000000a8ed32329d01fb5f27000000000027e2ea5b0000000082b4c2a389d911f1cef87b3f10dc38e8f5118ce5b83e160c5813447db849ea89c1d910841a3662747dd0e6e0040b1317be571384054a30f7e6851ebda9adab9c0a9394a5bb26479b697937fbe8b4a9d2780bee68334b2800000000000004454f5300000000000000000000000004454f53000000000000000000000000000000000000000004454f530000000000",
			chainID:        "aca376f206b8fc25a6ed44dbdc66547c36c6c33e3a119ffbeaef943642f0e906",
			expectedPubKey: "EOS7KtnQUSGVf4vbFE2eQsWmDp4iV93jVcSmdQXtRdRRnWj2ubbFW",
		},
		{
			name:                "R1",
			signature:           "SIG_R1_KE33Ucjr5N3GR4ZosFh8KtGMytHHNtnmdUaSoMLJVXpVXoC8B9zfoXYrLiQJZqroe3LKciaP2uJT7Myqqoo4PZH7iSnso8",
			payload:             "45e2ea5b22f87c6f74430000000001a0904b1822f330550040346aabab904b01a0904b1822f3305500000000a8ed32329d01fb5f27000000000027e2ea5b0000000082b4c2a389d911f1cef87b3f10dc38e8f5118ce5b83e160c5813447db849ea89c1d910841a3662747dd0e6e0040b1317be571384054a30f7e6851ebda9adab9c0a9394a5bb26479b697937fbe8b4a9d2780bee68334b2800000000000004454f5300000000000000000000000004454f53000000000000000000000000000000000000000004454f530000000000",
			chainID:             "aca376f206b8fc25a6ed44dbdc66547c36c6c33e3a119ffbeaef943642f0e906",
			expectedPubKeyError: "R1 not supported",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			sig, err := NewSignature(c.signature)
			if c.expectedSignatureError != "" {
				require.Equal(t, fmt.Errorf(c.expectedSignatureError), err)
				return
			}

			require.NoError(t, err)

			payload, err := hex.DecodeString(c.payload)
			require.NoError(t, err)

			chainID, err := hex.DecodeString(c.chainID)
			require.NoError(t, err)

			pubKey, err := sig.PublicKey(sigDigest(chainID, payload, nil))
			if c.expectedPubKeyError != "" {
				require.Equal(t, fmt.Errorf(c.expectedPubKeyError), err)
				return
			}
			assert.Equal(t, c.expectedPubKey, pubKey.String())
		})
	}

}

func TestEOSIOCSigningComparison(t *testing.T) {
	// try with: ec sign -k 5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3 '{"expiration":"2018-03-21T23:02:32","region":0,"ref_block_num":2156,"ref_block_prefix":1532582828,"packed_bandwidth_words":0,"context_free_cpu_bandwidth":0,"context_free_actions":[],"actions":[],"signatures":[],"context_free_data":[]}'
	wif := "5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3" // corresponds to: EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV
	privKey, err := NewPrivateKey(wif)
	require.NoError(t, err)

	chainID, err := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
	require.NoError(t, err)

	payload, err := hex.DecodeString("88e4b25a00006c08ac5b595b000000000000") // without signed transaction bytes
	require.NoError(t, err)

	digest := sigDigest(chainID, payload, nil)

	sig, err := privKey.Sign(digest)
	require.NoError(t, err)

	pubKey, err := sig.PublicKey(digest)
	require.NoError(t, err)

	assert.Equal(t, "EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV", pubKey.String())
}

func TestNodeosSignatureComparison(t *testing.T) {
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

	pubKey, err := sig.PublicKey(digest)

	assert.Equal(t, "EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV", pubKey.String())
}

func TestSignatureUnmarshalChecksum(t *testing.T) {
	fromEOSIOC := "SIG_K1_KW4qcHDh6ziqWELRAsFx42sgPuP3VfCpTKX4D5A3uZhFb3fzojTeGohja19g4EJa9Zv7SrGZ47H8apo1sNa2bwPvGwW2bb" // simply checked the last 2 bytes
	_, err := NewSignature(fromEOSIOC)
	require.Equal(t, "signature checksum failed, found a9c72981 expected a9c72982", err.Error())
}

//to do this here because of a import cycle when use eos.SigDigest
func sigDigest(chainID, payload, contextFreeData []byte) []byte {
	h := sha256.New()
	if len(chainID) == 0 {
		_, _ = h.Write(make([]byte, 32, 32))
	} else {
		_, _ = h.Write(chainID)
	}
	_, _ = h.Write(payload)

	if len(contextFreeData) > 0 {
		h2 := sha256.New()
		_, _ = h2.Write(contextFreeData)
		_, _ = h.Write(h2.Sum(nil)) // add the hash of CFD to the payload
	} else {
		_, _ = h.Write(make([]byte, 32, 32))
	}
	return h.Sum(nil)
}
