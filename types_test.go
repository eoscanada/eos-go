package eos

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/eoscanada/eos-go/ecc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChecksum256String(t *testing.T) {
	s := Checksum256{0x01, 0x02, 0x03, 0x04}
	assert.Equal(t, "01020304", s.String())
}

func TestSafeString(t *testing.T) {
	const nonUTF8 = "\xca\xc0\x20\xbd\xe7\x0a"
	filtered := strings.Map(fixUtf, nonUTF8)

	require.NotEqual(t, filtered, nonUTF8)

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	enc.writeString(nonUTF8)

	d := NewDecoder(buf.Bytes())
	var ss SafeString

	err := d.Decode(&ss)
	require.NoError(t, err)
	assert.Equal(t, SafeString(filtered), ss, "SafeString should contain filtered data")

	d = NewDecoder(buf.Bytes())
	var s string
	err = d.Decode(&s)
	require.NoError(t, err)
	assert.Equal(t, nonUTF8, s, "string should return unfiltered data")

}

func TestFloat64JSON_MarshalUnmarshal(t *testing.T) {
	f := Float64(math.Inf(1))

	var out Float64

	v, err := f.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, []byte("\"inf\""), v)
	err = out.UnmarshalJSON(v)
	require.NoError(t, err)
	assert.Equal(t, out, Float64(math.Inf(1)))

	f = Float64(12.3)
	v, err = f.MarshalJSON()
	require.NoError(t, err)
	err = out.UnmarshalJSON(v)
	require.NoError(t, err)
	require.Equal(t, f, out)

}

func TestUint128JSONUnmarshal(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedLo      uint64
		expectedHi      uint64
		expectedError   string
		expectedDecimal string
	}{
		{
			name:            "zero",
			input:           `"0x00000000000000000000000000000000"`,
			expectedLo:      0,
			expectedHi:      0,
			expectedDecimal: "0",
		},
		{
			name:            "one",
			input:           `"0x01000000000000000000000000000000"`,
			expectedLo:      1,
			expectedHi:      0,
			expectedDecimal: "1",
		},
		{
			name:            "value",
			input:           `"0x9ea6ce00000000000000000000000000"`,
			expectedLo:      13543070,
			expectedHi:      0,
			expectedDecimal: "13543070",
		},
		{
			name:            "max uint64",
			input:           `"0xffffffffffffffff0000000000000000"`,
			expectedLo:      math.MaxUint64,
			expectedHi:      0,
			expectedDecimal: "18446744073709551615",
		},
		{
			name:            "one more than uint64",
			input:           `"0x00000000000000000100000000000000"`,
			expectedLo:      0,
			expectedHi:      1,
			expectedDecimal: "18446744073709551616",
		},
		{
			name:            "value from nodeos serialization",
			input:           `"0x9d030000000000007d00000000000000"`,
			expectedLo:      925,
			expectedHi:      125,
			expectedDecimal: "2305843009213693952925",
		},
		{
			name:            "one less then largest ever",
			input:           `"0xfeffffffffffffffffffffffffffffff"`,
			expectedLo:      0xFFFFFFFFFFFFFFFE, // 18446744073709551614
			expectedHi:      math.MaxUint64,
			expectedDecimal: "340282366920938463463374607431768211454",
		},
		{
			name:            "largest ever",
			input:           `"0xffffffffffffffffffffffffffffffff"`,
			expectedLo:      math.MaxUint64,
			expectedHi:      math.MaxUint64,
			expectedDecimal: "340282366920938463463374607431768211455",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var i Uint128
			err := json.Unmarshal([]byte(test.input), &i)

			if test.expectedError != "" {
				require.Error(t, err)
				assert.Equal(t, test.expectedError, err.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expectedLo, i.Lo, "lo")
				assert.Equal(t, test.expectedHi, i.Hi, "hi")

				res, err := json.Marshal(i)
				require.NoError(t, err)
				assert.Equal(t, test.input, string(res))

				assert.Equal(t, test.expectedDecimal, i.DecimalString(), "numerical")
			}
		})
	}
}

func Test_twosComplement(t *testing.T) {
	tests := []struct {
		name         string
		input        []byte
		expectOutput []byte
	}{
		{
			name: "-1",
			// 0xffffffffffffffffffffffffffffffff
			input: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			// 0x00000000000000000000000000000001
			// the current algorithm will simply omit MSB 0's
			expectOutput: []byte{0x01},
		},
		{
			name: "-18446744073709551615",
			// 0xffffffffffffffff0000000000000001
			input: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
			// 0x0000000000000000ffffffffffffffff
			// the current algorithm will simply omit MSB 0's
			expectOutput: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		},
		{
			name: "-170141183460469231731687303715884105727",
			// 0x80000000000000000000000000000001
			input: []byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
			// 0x7fffffffffffffffffffffffffffffff
			expectOutput: []byte{0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		},
		{
			name: "-170141183460469231731687303715884105728",
			// 0x80000000000000000000000000000000
			input: []byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			// 0x80000000000000000000000000000000
			expectOutput: []byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expectOutput, twosComplement(test.input))
		})
	}
}

func TestInt128JSONUnmarshal(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedLo      uint64
		expectedHi      uint64
		expectedError   string
		expectedDecimal string
	}{
		{
			name:          "broken prefix",
			input:         `"mama"`,
			expectedError: "int128 expects 0x prefix",
		},
		{
			name:          "broken length",
			input:         `"0xmama"`,
			expectedError: "int128 expects 32 characters after 0x, had 4",
		},
		{
			name:          "broken hex",
			input:         `"0xmamamamamamamamamamamamamamamama"`,
			expectedError: "encoding/hex: invalid byte: U+006D 'm'",
		},
		{
			name:            "zero",
			input:           `"0x00000000000000000000000000000000"`,
			expectedLo:      0,
			expectedHi:      0,
			expectedDecimal: "0",
		},
		{
			name:            "one",
			input:           `"0x01000000000000000000000000000000"`,
			expectedLo:      1,
			expectedHi:      0,
			expectedDecimal: "1",
		},
		{
			name:            "negative one",
			input:           `"0xffffffffffffffffffffffffffffffff"`,
			expectedLo:      math.MaxUint64,
			expectedHi:      math.MaxUint64,
			expectedDecimal: "-1",
		},
		{
			name:            "max uint64",
			input:           `"0xffffffffffffffff0000000000000000"`,
			expectedLo:      math.MaxUint64,
			expectedHi:      0,
			expectedDecimal: "18446744073709551615",
		},
		{
			name:            "negative max uint64",
			input:           `"0x0100000000000000ffffffffffffffff"`,
			expectedLo:      1,
			expectedHi:      math.MaxUint64,
			expectedDecimal: "-18446744073709551615",
		},
		{
			name:            "largest positive number",
			input:           `"0xffffffffffffffffffffffffffffff7f"`,
			expectedLo:      math.MaxUint64,
			expectedHi:      0x7fffffffffffffff, //9223372036854775807
			expectedDecimal: "170141183460469231731687303715884105727",
		},
		{
			name:            "before smallest negative number",
			input:           `"0x01000000000000000000000000000080"`,
			expectedLo:      1,
			expectedHi:      0x8000000000000000, //9223372036854775808
			expectedDecimal: "-170141183460469231731687303715884105727",
		},
		{
			name:            "smallest negative number",
			input:           `"0x00000000000000000000000000000080"`,
			expectedLo:      0,
			expectedHi:      0x8000000000000000,
			expectedDecimal: "-170141183460469231731687303715884105728",
		},
		{
			name:            "value from nodeos serialization",
			input:           `"0x9d030000000000007d00000000000000"`,
			expectedLo:      925,
			expectedHi:      125,
			expectedDecimal: "2305843009213693952925",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var i Int128
			err := json.Unmarshal([]byte(test.input), &i)

			if test.expectedError != "" {
				require.Error(t, err)
				assert.Equal(t, test.expectedError, err.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expectedLo, i.Lo, "lo")
				assert.Equal(t, test.expectedHi, i.Hi, "hi")

				res, err := json.Marshal(i)
				require.NoError(t, err)
				assert.Equal(t, test.input, string(res))

				assert.Equal(t, test.expectedDecimal, i.DecimalString(), "decimal")
			}
		})
	}
}

func TestAssetMarshalUnmarshal(t *testing.T) {
	tests := []struct {
		in  string
		out Asset
	}{
		// Haven't seen such a thing yet though..
		{"808d5b000000000004454f5300000000",
			Asset{6000000, Symbol{Precision: 4, Symbol: "EOS"}}},
	}

	for _, test := range tests {
		bin, err := hex.DecodeString(test.in)
		require.NoError(t, err)

		var a Asset
		require.NoError(t, UnmarshalBinary(bin, &a))

		assert.Equal(t, test.out, a)

		marshaled, err := MarshalBinary(test.out)
		require.NoError(t, err)

		assert.Equal(t, test.in, hex.EncodeToString(marshaled))
	}

	/*
		18 = len
		808d5b0000000000 quantity = 6000000
		04  precision = 4
		454f5300000000  "EOS"
	*/
}

func TestAssetToString(t *testing.T) {
	tests := []struct {
		in  Asset
		out string
	}{
		// Haven't seen such a thing yet though..
		{
			Asset{6000000, Symbol{Precision: 4, Symbol: "EOS"}},
			"600.0000 EOS",
		},
		{
			Asset{-6000000, Symbol{Precision: 4, Symbol: "EOS"}},
			"-600.0000 EOS",
		},
		{
			Asset{10, Symbol{Precision: 5, Symbol: "SYS"}},
			"0.00010 SYS",
		},
		{
			Asset{-10, Symbol{Precision: 5, Symbol: "SYS"}},
			"-0.00010 SYS",
		},
		{
			Asset{6000, Symbol{Precision: 0, Symbol: "MAMA"}},
			"6000 MAMA",
		},
		{
			Asset{-6000, Symbol{Precision: 0, Symbol: "MAMA"}},
			"-6000 MAMA",
		},
		{
			Asset{0, Symbol{Precision: 255, Symbol: "EOS"}},
			"0.000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000 EOS",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.out, test.in.String())
	}
}

func TestLegacyAssetToJSON(t *testing.T) {
	LegacyJSON4Asset = true
	tests := []struct {
		in  Asset
		out string
	}{
		// Haven't seen such a thing yet though..
		{
			Asset{6000000, Symbol{Precision: 4, Symbol: "EOS"}},
			"\"600.0000 EOS\"",
		},
		{
			Asset{-6000000, Symbol{Precision: 4, Symbol: "EOS"}},
			"\"-600.0000 EOS\"",
		},
		{
			Asset{10, Symbol{Precision: 5, Symbol: "SYS"}},
			"\"0.00010 SYS\"",
		},
		{
			Asset{-10, Symbol{Precision: 5, Symbol: "SYS"}},
			"\"-0.00010 SYS\"",
		},
		{
			Asset{6000, Symbol{Precision: 0, Symbol: "MAMA"}},
			"\"6000 MAMA\"",
		},
		{
			Asset{-6000, Symbol{Precision: 0, Symbol: "MAMA"}},
			"\"-6000 MAMA\"",
		},
		{
			Asset{0, Symbol{Precision: 255, Symbol: "EOS"}},
			"\"0.000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000 EOS\"",
		},
	}

	for _, test := range tests {
		bytes, err := test.in.MarshalJSON()
		if err != nil {
			t.Error("MarshalJSON() error: %w", err)
		}
		assert.Equal(t, test.out, string(bytes))
	}
}

func TestAssetToJSON(t *testing.T) {
	LegacyJSON4Asset = false
	tests := []struct {
		in  Asset
		out string
	}{
		// Haven't seen such a thing yet though..
		{
			Asset{6000000, Symbol{Precision: 4, Symbol: "EOS"}},
			assertJson(t, 600, "EOS", 4),
		},
		{
			Asset{-6000000, Symbol{Precision: 4, Symbol: "EOS"}},
			assertJson(t, -600, "EOS", 4),
		},
		{
			Asset{10, Symbol{Precision: 5, Symbol: "SYS"}},
			assertJson(t, 0.0001, "SYS", 5),
			// "\"0.00010 SYS\"",
		},
		{
			Asset{-10, Symbol{Precision: 5, Symbol: "SYS"}},
			assertJson(t, -0.0001, "SYS", 5),
			// "\"-0.00010 SYS\"",
		},
		{
			Asset{6000, Symbol{Precision: 0, Symbol: "MAMA"}},
			assertJson(t, 6000, "MAMA", 0),
			// "\"6000 MAMA\"",
		},
		{
			Asset{-6000, Symbol{Precision: 0, Symbol: "MAMA"}},
			assertJson(t, -6000, "MAMA", 0),
			// "\"-6000 MAMA\"",
		},
		{
			Asset{0, Symbol{Precision: 255, Symbol: "EOS"}},
			assertJson(t, 0, "EOS", 255),
		},
	}

	for _, test := range tests {
		bytes, err := test.in.MarshalJSON()
		if err != nil {
			t.Error("MarshalJSON() error: %w", err)
		}
		assert.Equal(t, test.out, string(bytes))
	}
}

func assertJson(t testing.TB, amount float64, symbol string, precision uint8) string {
	t.Helper()
	data, err := json.Marshal(map[string]interface{}{
		"amount":    amount,
		"symbol":    symbol,
		"precision": precision,
	})
	if err != nil {
		t.Fatal("json.Marshal error: %w", err)
	}
	return string(data)
}

func TestSimplePacking(t *testing.T) {
	type S struct {
		P string
	}
	type M struct {
		Acct AccountName
		A    []*S
	}
	cnt, err := MarshalBinary(&M{
		Acct: AccountName("bob"),
		A:    []*S{{"hello"}, {"world"}},
	})

	require.NoError(t, err)
	assert.Equal(t, "0000000000000e3d020568656c6c6f05776f726c64", hex.EncodeToString(cnt))
}

func TestSetRefBlock(t *testing.T) {
	tx := &Transaction{}
	blockID, err := hex.DecodeString("0012cf6247be7e2050090bd83b473369b705ba1d280cd55d3aef79998c784b9b")
	//                                    ^^^^        ^^....^^
	require.NoError(t, err)
	tx.setRefBlock(blockID)
	assert.Equal(t, uint16(0xcf62), tx.RefBlockNum) // 53090
	assert.Equal(t, uint32(0xd80b0950), tx.RefBlockPrefix)
}

func FixmeTestPackTransaction(t *testing.T) {
	stamp := time.Date(2018, time.March, 22, 1, 1, 1, 1, time.UTC)
	blockID, _ := hex.DecodeString("00106438d58d4fcab54cf89ca8308e5971cff735979d6050c6c1b45d8aadcad6")
	tx := &Transaction{
		TransactionHeader: TransactionHeader{
			Expiration: JSONTime{stamp},
		},
		Actions: []*Action{
			{
				Account: AccountName("eosio"),
				Name:    ActionName("transfer"),
				Authorization: []PermissionLevel{
					{AccountName("eosio"), PermissionName("active")},
				},
			},
		},
	}
	tx.setRefBlock(blockID)

	buf, err := MarshalBinary(tx)
	assert.NoError(t, err)
	// "cd6a4001  0000  0010  b54cf89c  0000  0000  00  01  0000000000ea3055  000000572d3ccdcd
	// 01
	// permission level: 0000000000ea3055  00000000a8ed3232
	// data: 00"

	// A transaction:
	// expiration: 82a3ac5a
	// region: 0000
	// refblocknum: 2a2a
	// refblockprefix: 06e90b85
	// packedbandwidth: 0000
	// contextfreecpubandwidth: 0000
	// contextfreeactions: 00
	// actions: 01
	// - account: 0000000000ea3055 (eosio)
	//   action: 000000572d3ccdcd (transfer)
	//   []auths: 01
	//     - acct: 0000000000ea3055 (eosio)
	//       perm: 00000000a8ed3232 (active)
	// ... missing Transfer !

	assert.Equal(t, `4d00b35a00003864b54cf89c00000000010000000000ea3055000000572d3ccdcd010000000000ea305500000000a8ed323200`, hex.EncodeToString(buf))
}

func TestUnpackBinaryTableRows(t *testing.T) {
	resp := &GetTableRowsResp{
		Rows: json.RawMessage(`["044355520000000004435552000000000000000000000000"]`),
	}
	var out []*MyStruct
	require.NoError(t, resp.BinaryToStructs(&out))
	assert.Equal(t, "CUR", string(out[0].Currency.Name))
}

func TestStringToName(t *testing.T) {
	i, err := StringToName("tbcox2.3")
	require.NoError(t, err)
	assert.Equal(t, uint64(0xc9d14e8803000000), i)

	//h, _ := hex.DecodeString("0000001e4d75af46")
	//fmt.Println("NAMETOSTRING", NameToString(binary.LittleEndian.Uint64(h)))
}

func TestVaruint32MarshalUnmarshal(t *testing.T) {
	tests := []struct {
		in  uint32
		out string
	}{
		{0, "00"},
		{1, "01"},
		{128, "8001"},
		{127, "7f"},
	}

	for _, test := range tests {
		v := Varuint32(test.in)
		res, err := MarshalBinary(v)
		require.NoError(t, err)
		assert.Equal(t, test.out, hex.EncodeToString(res))

		require.NoError(t, UnmarshalBinary(res, &v))
		assert.Equal(t, test.in, uint32(v))
	}
}

func TestNameToString(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"0000000000000000", ""},
		{"0000000000003055", "eos"},
		{"0000001e4d75af46", "currency"},
		{"0000000000ea3055", "eosio"},
		{"00409e9a2264b89a", "newaccount"},
		{"00000003884ed1c9", "tbcox2.3"},
		{"00000000a8ed3232", "active"},
		{"000000572d3ccdcd", "transfer"},
		{"00000059b1abe931", "abourget"},
		{"c08fca86a9a8d2d4", "undelegatebw"},
		{"0040cbdaa86c52d5", "updateauth"},
		{"0000000080ab26a7", "owner"},
		{"00000040258ab2c2", "setcode"},
		{"00000000b863b2c2", "setabi"},
	}

	for _, test := range tests {
		h, err := hex.DecodeString(test.in)
		require.NoError(t, err)
		res := NameToString(binary.LittleEndian.Uint64(h))
		assert.Equal(t, test.out, res)
	}
}

func BenchmarkNameToString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NameToString(5093418677655568384)
		NameToString(6138663577826885632)
		NameToString(11148770977341390848)
		NameToString(14542491017828892672)
		NameToString(3617214756542218240)
		NameToString(14829575313431724032)
		NameToString(3596594555622785024)
		NameToString(15335505127214321600)
		NameToString(15371467950649982976)
		NameToString(12044502819693133824)
		NameToString(14029427681804681216)
		NameToString(14029385431137648640)
	}
}

func TestPackAccountName(t *testing.T) {
	// SHOULD IMPLEMENT: string_to_name from contracts/eosiolib/types.hpp
	tests := []struct {
		in  string
		out string
	}{
		{"eosio",
			"0000000000ea3055"},
		{"eosio.system",
			"2055c61e03ea3055"},
		{"tbcox2.3",
			"00000003884ed1c9"},
		{"tbcox2.",
			"00000000884ed1c9"},
		{"quantity",
			"0000003ebb3c8db6"},
		{"genesis.1",
			"000008003baca662"},
		{"genesis.z",
			"0000f8003baca662"},
		{"genesis.zzzz",
			"f0ffff003baca662"},
	}

	for idx, test := range tests {
		acct := AccountName(test.in)
		out, err := hex.DecodeString(test.out)
		require.NoError(t, err)
		buf, err := MarshalBinary(acct)
		assert.NoError(t, err)
		assert.Equal(t, out, buf, fmt.Sprintf("index %d, was %q", idx, hex.EncodeToString(buf)))
	}
}

func TestAuthorityBinaryMarshal(t *testing.T) {
	key, err := ecc.NewPublicKey(ecc.PublicKeyPrefixCompat + "6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV")
	require.NoError(t, err)
	a := Authority{
		Threshold: 2,
		Keys: []KeyWeight{
			{
				PublicKey: key,
				Weight:    5,
			},
		},
	}
	cnt, err := MarshalBinary(a)
	require.NoError(t, err)

	// threshold: 02000000
	// []keys: 01
	// - pubkey: 0002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf
	//   weight: 0500
	// []accounts: 00
	assert.Equal(t, `02000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf05000000`, hex.EncodeToString(cnt))
}

func TestActionNoData(t *testing.T) {
	a := &Action{
		Account: AccountName("eosio"),
		Name:    ActionName("transfer"),
	}

	cnt, err := json.Marshal(a)
	require.NoError(t, err)
	// account + name + emptylist + emptylist
	assert.Equal(t,
		`{"account":"eosio","name":"transfer","authorization":[]}`,
		string(cnt),
	)

	var newAction *Action
	require.NoError(t, json.Unmarshal(cnt, &newAction))
}

func TestHexBytes(t *testing.T) {
	a := HexBytes("hello world")
	cnt, err := json.Marshal(a)
	require.NoError(t, err)
	assert.Equal(t,
		`"68656c6c6f20776f726c64"`,
		string(cnt),
	)

	var b HexBytes
	require.NoError(t, json.Unmarshal(cnt, &b))
	assert.Equal(t, a, b)
}

func TestNewAssetFromString(t *testing.T) {
	longDecimals := "." + strings.Repeat("1", math.MaxUint8+1)

	tests := []struct {
		in          string
		amount      int64
		precision   uint8
		symbolCode  string
		expectedErr error
	}{
		{"1000.0000000 TEST", 10000000000, 7, "TEST", nil},
		{"1000.0000 TEST", 10000000, 4, "TEST", nil},
		{"1000 TEST", 1000, 0, "TEST", nil},
		{"1000.1 TEST", 10001, 1, "TEST", nil},
		{"1000.001 TEST", 1000001, 3, "TEST", nil},
		{"1.0001 TEST", 10001, 4, "TEST", nil},
		{"0.1 TEST", 1, 1, "TEST", nil},
		{".1 TEST", 1, 1, "TEST", nil},

		{"", 0, 0, "", errors.New("input cannot be empty")},
		{".00.001", 0, 0, "", errors.New(`invalid asset amount ".00.001", expected amount to have at most a single dot`)},
		{"1 ABCDEFGH", 0, 0, "", errors.New(`invalid asset "1 ABCDEFGH", symbol should have less than 7 characters`)},
		{"1 A AND B", 0, 0, "", errors.New(`invalid asset "1 A AND B", expecting an amount alone or an amount and a currency symbol`)},
		{"1 A AND B", 0, 0, "", errors.New(`invalid asset "1 A AND B", expecting an amount alone or an amount and a currency symbol`)},
		{longDecimals, 0, 0, "", fmt.Errorf(`invalid asset amount precision "%s", should have less than 255 characters`, longDecimals)},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			asset, err := NewAssetFromString(test.in)
			require.Equal(t, test.expectedErr, err)

			if test.expectedErr == nil {
				assert.Equal(t, Int64(test.amount), asset.Amount)
				assert.Equal(t, int(test.precision), int(asset.Symbol.Precision))
				assert.Equal(t, test.symbolCode, asset.Symbol.Symbol)
			}
		})
	}
}

func TestNewFixedSymbolAssetFromString(t *testing.T) {
	tests := []struct {
		in          string
		amount      int64
		expectedErr error
	}{
		{"1000.0000 SYS", 10000000, nil},
		{"1000", 10000000, nil},
		{"1000 SYS", 10000000, nil},
		{"1000.1 SYS", 10001000, nil},
		{"1000.1", 10001000, nil},
		{"1000.01", 10000100, nil},
		{"1000.001", 10000010, nil},
		{"1.0001", 10001, nil},
		{"0.1", 1000, nil},
		{"0.0001", 1, nil},
		{".0001", 1, nil},

		{".00001", 1000, errors.New("symbol 4,SYS precision mismatch: expected 4, got 5")},
		{".0001 BOS", 1000, errors.New("symbol 4,SYS code mismatch: expected SYS, got BOS")},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			asset, err := NewFixedSymbolAssetFromString(Symbol{Precision: 4, Symbol: "SYS"}, test.in)
			require.Equal(t, test.expectedErr, err)

			if test.expectedErr == nil {
				assert.Equal(t, Int64(test.amount), asset.Amount)
				assert.Equal(t, "SYS", asset.Symbol.Symbol)
				assert.Equal(t, 4, int(asset.Symbol.Precision))
			}
		})
	}
}

func TestNameToSymbol(t *testing.T) {
	tests := []struct {
		in          string
		expected    Symbol
		expectedErr error
	}{
		{".....l2nep1k4", Symbol{Precision: 4, Symbol: "CUSD", symbolCode: uint64(1146312003)}, nil},
		{"......2ndx2k4", Symbol{Precision: 4, Symbol: "EOS", symbolCode: uint64(5459781)}, nil},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			actual, err := NameToSymbol(Name(test.in))
			if test.expectedErr == nil {
				require.NoError(t, err)
				assert.Equal(t, test.expected, actual)
			} else {
				assert.Equal(t, test.expectedErr, err)
			}
		})
	}
}

func TestNewSymbolFromUint64(t *testing.T) {
	tests := []struct {
		in       uint64
		expected Symbol
	}{
		{293455872769, Symbol{Precision: 1, Symbol: "CUSD", symbolCode: uint64(1146312003)}},
		{5327108, Symbol{Precision: 4, Symbol: "IQ", symbolCode: uint64(20809)}},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			actual := NewSymbolFromUint64(test.in)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestStringToSymbol(t *testing.T) {
	tests := []struct {
		in           string
		expected     Symbol
		expectedName string
		expectedErr  error
	}{
		{"1,CUSD", Symbol{Precision: 1, Symbol: "CUSD"}, ".....l2nep1k1", nil},
		{"2,CUSD", Symbol{Precision: 2, Symbol: "CUSD"}, ".....l2nep1k2", nil},
		{"2,KARMA", Symbol{Precision: 2, Symbol: "KARMA"}, "...42nemc55k2", nil},
		{"4,IQ", Symbol{Precision: 4, Symbol: "IQ"}, "........e54k4", nil},
		{"4,EOS", Symbol{Precision: 4, Symbol: "EOS"}, "......2ndx2k4", nil},
		{"9,EOSEOSA", Symbol{Precision: 9, Symbol: "EOSEOSA"}, "c5doylendx2kd", nil},

		{"EOS", Symbol{}, "", errors.New("EOS is not a valid symbol")},
		{",EOS", Symbol{}, "", errors.New(",EOS is not a valid symbol")},
		{"10,EOS", Symbol{}, "", errors.New("10,EOS is not a valid symbol")},
		{"10,EOS", Symbol{}, "", errors.New("10,EOS is not a valid symbol")},
		{"1,EOSEOSEO", Symbol{}, "", errors.New("1,EOSEOSEO is not a valid symbol")},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			actual, err := StringToSymbol(test.in)
			if test.expectedErr == nil {
				require.NoError(t, err)
				assert.Equal(t, test.expected, actual)

				asName, err := actual.ToName()
				if test.expectedName != "" {
					require.NoError(t, err)
					assert.Equal(t, test.expectedName, asName)
				}
			} else {
				assert.Equal(t, test.expectedErr, err)
			}
		})
	}
}

func TestStringToSymbolCode(t *testing.T) {
	tests := []struct {
		in            string
		expectedValue uint64
		expectedName  string
		expectedErr   error
	}{
		{"CUSD", 1146312003, "......24eheo3", nil},
		{"KARMA", 280470110539, ".....kehed.of", nil},
		{"IQ", 20809, ".........1cod", nil},
		{"EOS", 5459781, "........ehbo5", nil},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			actual, err := StringToSymbolCode(test.in)
			if test.expectedErr == nil {
				require.NoError(t, err)
				assert.Equal(t, test.expectedValue, uint64(actual))
				assert.Equal(t, test.expectedName, actual.ToName())
			} else {
				assert.Equal(t, test.expectedErr, err)
			}
		})
	}
}

func TestSymbolCode_String(t *testing.T) {
	tests := []struct {
		in          uint64
		expected    string
		expectedErr error
	}{
		{1146312003, "CUSD", nil},
		{280470110539, "KARMA", nil},
		{20809, "IQ", nil},
		{5459781, "EOS", nil},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			actual := SymbolCode(test.in).String()
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestNameToSymbolCode(t *testing.T) {
	tests := []struct {
		in             string
		expected       SymbolCode
		expectedString string
		expectedErr    error
	}{
		{"......24eheo3", SymbolCode(1146312003), "CUSD", nil},
		{".....kehed.of", SymbolCode(280470110539), "KARMA", nil},
		{".........1cod", SymbolCode(20809), "IQ", nil},
		{"........ehbo5", SymbolCode(5459781), "EOS", nil},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			actual, err := NameToSymbolCode(Name(test.in))
			if test.expectedErr == nil {
				require.NoError(t, err)
				assert.Equal(t, test.expected, actual)
				assert.Equal(t, test.expectedString, actual.String())
			} else {
				assert.Equal(t, test.expectedErr, err)
			}
		})
	}
}

func FixmeTestJsonWithZLIBcompression(t *testing.T) {
	jsonString := `{"timestamp":"2018-06-09T11:57:57.500","producer":"eosio","confirmed":0,"previous":"000000b0090f3beb6593010f6d0623f7367af83544afabe1715e1a43fe94c699","transaction_mroot":"7c8e49e232e0ef11b997246fbde51a596e95ce46a65f5ca21ab0b8769c34ec8a","action_mroot":"d1b6c6a9c016ad1c130b161f6041e598ba170bcee7d7dd5ef8b9d7b489f7ff75","schedule_version":0,"new_producers":null,"header_extensions":[],"producer_signature":"SIG_K1_K2wGMewyNR4sgLeKgzJdQk6sXSkwvpBweNvSPMv3wE8gDUYDvvMUtwDpsswUz17X33zmk9gEGcAuDGdkbKaK7NQeR64y8M","transactions":[{"status":"executed","cpu_usage_us":3270,"net_usage_words":1108,"trx":{"id":"c9f28cf4fa099d0c64b3bc1fe4b0a389754b181e46d23b6640d323f418beea4f","signatures":["SIG_K1_JzFpJxiVsU1vtvaW2QUBf65f1XjFmus2UVqpJjz5RyH2ippF34g79NkFmZJhTHvskXn19MQYHq4R8KWt7GMHauqVsNTKSK"],"compression":"zlib","packed_context_free_data":"","context_free_data":[],"packed_trx":"78daed7c7d8c5cd775dfbdef7be6cd709fa48dbc22d9e8ce331daf5cc9597eadd6b46cf1b15a511b39faae23d80886c3dda1b833fbc599598a7454eeba5163226902fe61a029dac40ce254ae032244820446fe0898404555c02e14c44e82043088360894a20154204583366d7a7ee7dcf7e6cdecae48a78ae1b65989f3debbef7e9c7bbeee39e79e77ffe3ad039ffd15f5d8fff8f0b715fd39f851ff69e61fd2efc90ffde4affe8e56aaffb22305ea2b7f7ee4c87ffdd7a5a277e8a1d55fa512a57f567fe08cde5267bcad2b57ece58a3aa3e87f7de58c83fb608bfff4d619678b6ab8721f5275fabb42b74a5ed1afa6ea1e7580d200572ea53ed1d0c70557145fe1f22de9f08ad4da02002e8fafa9c4f96fcebd6e7beda2df3abbde1b2817f7496b71b0bcbed65c6a0d5acdfef2e7daaa8ae2ef5f3cdf5eec3637dabdd5e57e1f155a9b83f3ebbde5cfb5505d5550e941a934e8b5d6fab69bd15a018fb0b8d9ebb5d706cd5e7bb1bd7cb1dd531e8a6b79f16079b52d45f1d2d9e6b9e5b5a5e6f2ec31554749954a5e690fb82042419d0a7aedd5f58b6d2e53dc1395f507eb3d29da9757dbdca03949d904576bafd3449aad7ebf4d930f5114acb65717372e4bd7213da063794a7aedd652b3841c157327bdf685cd651a0a1355351eab5c7444f95cd66fd33496dae7da34c72595b8b3c7e338aeaa280c63fb778fa29ff0def0de38d4f74d86b19e74ee0b43e7fb507abf52d1fdf1fdcaf57cbda183c077b5d2e17ff1ef01c0ebbdcb8ebabff999f685179e3e7c9849d0df5c3d727cf6c5c34d75ef8e17876767f0e23e7ab1d6dee5c5a4bbb6feaafabe99e6679e39ce183a7ca43ca1f9179e7ee6c5e6e1d91227acb42fb657e6d5fd3775d1e6e8e13bb082ede670a902170d8e369b878ff6db8385dd4779e6c5a3cd632bed7e7fe1c5e3cd797efc586b65657db1450497b279eee85873e1f0ccc6e6d995e5c566b77df99917679b0b2f3ec14de6ecddfcbcfac0974a30cfbc378fcf5fbe3c0ae168ff33395c4778909911b88e14701d6ebec7cc165eccb8f1717b373f7f59a5070b08e756375706cbfde5571edde8ad6facf7dbf317d5073fbecbebd60655b8d89e27b09ac7d65aabedf93d68f6c82776b6fed8e6da5db73ff2a19ded67175b6b8bf472d8f8b23af6c19df58eb52fb517476acdf934f0ca65f53196c3d50d9505ab8c4475ca3bd76bb7d57ca5a655a6e2e0c9b9238f1cfb60f5e63fd3c7636594d19949b2f9b87c5f2fee33759a442df14fdd7afd13db6fc4fbe8c5432eff444912ffa6ef3a5bfa4ab89551cfd38e97cd74d3daace3995a96748c4aa632d588e95f85fe559d93c69976a2d49daf53ed464095d4ac1399da29f52875e9de68f8ee491364ba43effc53e1d5d43fa55ea8abd83989ffd058a5a17132af937a4fd6952bc5611ad1c8a99a4789a11aa8a6eaaae6c4b18930204a4c442d4c94459d4644d522aec44f345395d73a4da3c504420a1042e3e1d904f3751d9b9a0992074c8d866a84b8a886aa69cc2da47f0aef693287e462e812ce3ad33443d3314eb27fbf34ac48c32a9a29025d635c677774387742871674e8ef023a92f7131d3a3960d111db867c4d52ee67d2845d53cd881d0e988aa976e58e10649c032636aacb7789db0869ae982471d9e3841253c57d55eeb93b7ec9e34f721d6a58cb663ae04e82f187e29fd7da2df836e9a61ebd21d85d0c063a28a20321cac3135d08ee83fb19c31e90e5f0d044db07a8bed7c13c5c20aec003cdcac3e09e5cd0838b6ed1d1dfdb8f4b36d549be9f2e87e432d3491edc9f61c00440c6371ced6f09ad7958a68a838b97065da39f2500021e32f5bba96b34b141b6f5a9ba93fd35fd858d002263e89593fd2f2a709fae030627d38334a086dec2819448e91041934f369c9a1bf3441a8a994e9d504a38411103a60126a0b3ad4eaa8993bc1813a44ad48c5841252769b28e093a0deadfb89dd46f04c35e2c07e7ad63e3f340a90b3664f035c1d86df88037d34f1198bed184d9da7e412886303e20700d030246351e83042ae8e4f18938fe0f0eab22bf20e9dbaa9b56807f3c35a8cbe34e44003297bb32d790a0f4de20e02ac683807924606994850302faafd585de09356d184c9fd8ee873252788685a7924d77320d00432e008a18c088c553443310d1d410341ed491415d1ab492051d13d050f54e96a21ff4c8a5f2e8d8411b32a85b0cca6f5d2e2025c5833273b0cc8722f3611ab06a28cb7cb087cc8722f33c0da96565de68bec930eedb8aa546c57feeb87a04c724365181624f50ec8ecc168ce0be41e88f8c0b14bb40710814fb16c513325baf40b1c09a4d5814075c001433c9c3128afd1114d3680e541e8d2ef3a52e43e952f049e3269d2c620473116a32dd028bef72f590e473bcb297e33b6429147c07a9cff4b6f8664cfadc03e0620c83d1581107fc2473925a826f7f88ef48e49fd0fd73be1310baad1c52db6a8781cd6216a8e34e8d0969052c027a235922a88988ec7167f28e75b2831dd668c0e9541a76b3e9d51b7768040c18e74999be5d1eb249eec6c54493d4eba67ee66deed911f4a87181203fdbb800dd1180eb425a21b8970eba3948dd503b9728b9374421c905cc82d343f0629ad4f4705287fe9f99d4b18e65fb696706937a0f58ca93f2eda498190592cce90813df15300980f1478051f1975cd723f6d4656d105a6d10b14ae866ae08141e20abacd3f989ad96a232d5d4a33509c0526567acf268b75167bc671f95594078217a837413091de61440ff38d03f9ed53ffb58638c6a795ea03cd102fbac1ef278e12f54bdc3c8805270f76aaba16e766bcb626921734620d380cc2934e358d716a2893d7a258832c5d412fdf15b8eeb1281dc1dea5a4fcb8f12b1f0457f07236ba347e0f96f90ae8a86c6a7202e1c052f18052fda05bc2047d8a8d9ca9c58e86e2fd7dd8eed351c997454ea32cabbdcbd6667ac32db0c8c97a44197839dc4eccf258cc6d75caf11d4c6f5ef55c19fb383c14bf883048ee14f10078e1b439cb717e2caec7617d8f30aec3985718151f35e75a95759084598ec4aea588440528009189cd4a7034c04654c0c39e98fb4eb8f8bfa1827dd205629a3c113aac28709c1e56464001b21b0e18e72b9378a0d7f948d7c2e2810e1e672278d725408661c52d0cf8c4d9e6ce0e1e48d8fc98f50f9cb6238ba77a27261488a7923a6a30bb7a8301d03cccb1db56bc614833b4a656bcae5221cecb01fd93e2e64c41d9191a1c1229417b6b75de6a6cd1e353bbb56c6b297966cbf21075c134f3f2ab064bae984784cfbe81ff9ae8d7aeee0fbe2d14666423ca909f168fd1b6c074558614adc30e6e057c4c10f861e6d258db15055c5daaaa61554ab5a6b2bc680283131b014833631558bb9123f99aaa9e6b5c4da8ae0e04764f507bca445ecd14e98889cb909f1d327d84f174f8e8317785f852bc997842ee45b4cd10c13ebe073c39a34acc736dee1961cfc7174387742872be870bf0be848de4f74b870f0b9e13edb10568336756228ff00b9e4f5aedc55c5abdf276e3edd7989d3a88cf9f275dcd7c57f9f20430aec5881fffecff7f90e6be7b25e7a17be24aa3869986830f55b9ab1e6679ff716609825a7d86bcd5831400e372e489318d6dfacf3a6668db18f5608d0e9a6a6e9fa8818cc3a5fd152765d273e1b8754f6458deb09f5b35a9effa5e6f59ca365d96304ed552dded59b1aa3bc86dbb77443cb60977021cf7dd6d9d624eee1d399818f997d4db375f69a086968ef3b450f54dbe90eabdfdcbd3a57cefe52d1cbecbaee2487f277f8b9d49d7506c5ebabf4fa43f4fa21f76b3a9990e94d03763b9567b3198c84d71c0a421030751f759f228c6c3deabe2c933f83b93fea3e474f7fa23a123809b377ec2de205606bea737b7bfbea4f7ffdf7ffe2c616fabca993802c1ff514bb4b6f29c4887e002a6e521e9f240f34213e75e6b36b1a608008733224fb555f26dce27a5d2396e0b1ba8679e5675bb4020dcdab42b6a063d9d4a7c56fe3023526db16847640ec9a0beef0bb74a7630ea1347c6079ae437cf100c679acc1c3cd7139421393724984aa53ac372dd21bfb85faf81da47e47188097052858ea85ec79f44945898bb5e2961e417598bd090699e4db9bb84df8f62fe0954ff1ed5b283d28c87f47a3ed5389f73c09033ab24afc9b5aaaa2d5db72ff0e9afd21c2bc9f8495c23a1a412fb2f880dac4f8790310e26d8df97f58668a90810f1ea011b9371fbea39e756e8bb83103a41229934ec40742d0ca019811f89ea7224ad191080f02545358a281fa30bbcdcc1975ecf05440189d161f3b97ca08d24d4a02975cddd2046a1d900211bac7c9a4c00f73aa9146d3bbf072220804f51e2651cbd977a603dee47b21e9ac531bd2222ab0ee8ea03c1aa23cd989f1999d18c72591cb14231ecd4c81f78704efb7f52e883f04e7e8eb7ab8868a35f42c472819d70ae1328b6beac286d310a229217bda22b840f6d709d91f6164dfd639b641df12ba2725ccea533f12e5b4d158a6842a2851bc624a88d04ea79abb616e79d8320b412ebeaacace77520efff912905563e33ac5b80599efc80134f3d3b977cd5d63e21a52ec81195509ba39b66204e98f95cacf94a12e78dccfc1f601b6cf6160112265c3d4450097e3c6a3b02981cdb7b011d3df0138d27343e8beacf90d5b6bef6a0e1dbf9d3b2e62b222244feaaf305d87b664c8fcb290bded58d3300d1f518a4c5db26b6194862c9330681111d517d26070c3f8af43a76621e2045e165e1c5a2cfe1ba998cd12775112a248831f469496a611749387c71d4087342f992e1c3a09c0b27801ab3af306bce0a958c232bae427b1761f3adfc1702db43edf2ecee2784d0e22ef55d9863c1c319054fc79b784d122b8213d691b76180d413845bc82566a6d6b626a12d0180d6e8c54cebb0580bbf56c372a98a2f042726fc8929191c4c1460fb474189779f4d4790f5a3a424b70e22749c50dd1eb9db600614f2f27c16e5e2ab67f2c5401e3db82c4321208341e8798ac53cbd0840c4d984313029a70149a4922c14f794c026fabd869ac3422c4a7091d902e1bc6642f748a3d63d8785db8682494764722115f915d5dc46121ae705989e33bd979a6b6c8b6580768291e3164112b76129297f7e0c54cf7b39b3435ae1c58db3220dbd213db926ce9dcb6f4d8b664f64fc2fde4077bd818e38bc1c5042c67d3dcd1b60759f1927932ca45478508d41f6a78312fd3812c1a812c1a812cd3816c72851cf016e68164325f47b8e9f088d87c83e797fc7d78c31ac64384601e5fb0a4f8580323f8c95380a826db4e84d7d4861dd47011f78c8f85c5c3c61e35022053980b07a268cc625dc11243eff927c14f8d57709ad414e999bca168488f371ee7d92c230dc9da278472f444c14576ff2e88ffc473c3f180846f85085ca6ed3e193dd1cc1865f4ccf86457cc2767881725be781c8e274ee12782a9e10961a71dec508a6674bb1cf027146e7d8ad807ce14b6ce522660c47dc9a614223f2ede11faab06619238dfb08da40b8dd931aa82e47196ab34107020f4bcef280c97d8087224efa85594d792b122be430837e2106e6518c2adc0eaad0cadde11360d854dfd219b86cca6411e2780826606c2a8d3764c89f3544cb16b14d8ed2c84bbaac311537f2cdc65a2220a1211e6f2e887f57d7dc839ea8c0fc601f7676c334fe23864b73fe446650f027a1339091a3c5c0154099cd1d28cc3641e0b4c054d9397608212664f5d2180f56bf4f38f3e7d803804f8a84abc63829c646c6afe1b577b5bf9e6ab677761bad9d46aca1ad5cf7ef95b5ffd16c421fb2a6eb0276dd77923d582ec9b5ffdd6479eae0ba6f580b413add90b07c0e46c734dad249f24ede2e43be7c17043558bf1636d4f08f4d40a49a2986e8c050509d4cc9da206d01f226fd85ac67e2a8baf2cb376a738ccb737ce4b0c8978dec939ef7c073bd4b8bb044571a9c3fde27980e7c1f07903cf1bc36752af04101bd376c7c26e8987561ed9c06c381c07901d5c5a08fd91bddcf2ee0bd42e2f610deb5573b284e8883da428c63630bcb5d02a71c39616ecabf8aac7d1cb7ce520551137107b88772a8c68a830b404c3b1b3da888602ece74313abf9d62c1311167bcfb8a3f2eb73e5468539dc62acc2bd7263de9e654996748348dca18a8dfd582108d8a525b1939ddb9d21231a6492494e62527992b5d4c60503298d991ba43b7fb4bbbda24f3bfb903d662b4d11e98f8ae88f58e2b4d01f15d61f552cc7d0843445f68ea25ca154728542ef784582c670e11e5709858db260ffcc37bef1ed5f3f3e2ad7217a845c1f2a415211b90ec7e53a845c873be53abefd7702fdbd28d093b9404f9504fa6f26ebffd8b19b92365a6f4dd26087ebe3ef707dfcdcf5f18b78bc3766a006af73dc8f0dd492eb6382379080c33918768ba69b3ae2f5103dbaa524213f777c1a52930d46542a5c1f98a2ec78c5a58d4936ee52ff5936eebca18d2f6ba264a1c8de0659cabfece860afcc21f6c82479282cf89f7c349b3f44d2b74d2f258748494e0eb17eea501da4113d9e6711b99be4f4fa9c45441a11a154db9494804fae1b728ac6b926a879c2496c5d907b47e04e40507c4e2572aca0248ad3895c7600b1c142e3875da60cb76a708496ef7c76c33967c883b8143943412c310fcb99bae04cc9983162c98161febbe338bc34c8e2702bcf1a72a175f8423a3f494869cdd988266932c4912a08274dd385a476860a65dd80d951e1a84476db493e9a729e0754291b366964e3f89245e9bd60b7eff0e03f8f642cf22db196489c22421659c5b03220bd7b3adb76599372225985c6493818c5e39145cf8a35f9414080fc80200f5454b263d62324352602194120238953443be314939ca429b63bab10570452e5f99b12c2285e4908838c5006bb42667a9e09476e163bf0411e01e28ca25b9c5114ffbae614b771dfcdcd5d2fe49819b7cb5bc3d69437faf9b2ef86403571e198efe6967c3750c018253164d9c1088639aca5e50641685f82d010abc3230d6119abd306ba2288af7b1eac09d933b41645adc1098676df90dbdd02b1301207065dab661a555ed0ac87019527a4f1f3441c4e1f31556b85f860daaab0089e6dc25992bb5758f7c9589d617509aba46a1d16f6b95db64b64c117c7cac5c295d4f763a986672c528810130df07c76d34be20c9e105b3f3e768a545ab58ab89b466c38c0c9706138b8234e86d820c8cb0c4c5552cdc87e70a9dad006a9da34a83b74b55b1fb22766577e976c10d98f32352ab43648956d90185b52b04d68866c83b8b90d52cd6d9098510a62c5627c108bb20d42ae3629171826d5d26855b1333c5093ed0c257686073bc32bd919b1b533be34d4c0d611cd1335410807eb8b9b2b4c6575ad37b43588626159e18250a470914aea92c2256f192eb66d4fd8f544e18685c295e1ec564603d1a27054bd721a35ab576d19cb05a3e921a3b9440182437725f597d52d0d0f93d6655e63ad0be9d385d6f58ba432d6baeeb83dc0b6338295f16ffeff1ecc99d919cc911d8040760002499cbeeb98cea1fdf9c3243f7080e7d09e019e437fbb019e8fec1de099cc033c53a500cf7714fbf9db5fb12779c59e7a3f56ec77bd3bafd847ffef59b17ff1fd5cfa6aa5a52ff99b2c7d5169e99b7abf973e1b8bfabbd5ef3b5efd7ecbd3c314f1dbcaa68ab9c904eb415231229511a472c64490ca398e7f80cdad544e735d91c9e38c5617a61fb4389305f975bc099cf808d2a0c71a2e27945556939c713eccc3c07864fcbf2497e74467bf4c75df56a4b4920fed978e9885a4d94276cd4f383d6604e248204e04e2a90262bc1615c2128addd464166f934725417f2af572f98ea073254bdf78e32eb21abac865f9f60bf9f64bbb81aa8849b17c17af44be3d7191593d17f2ed41be03eb1d492f04ddc932746f2a0b9ecac11beec786c3fdd850c00b0b18c25118447fdb4dd7083a644f2058c3df1625f31bdf25b760f23ddd823cbc3ce617fcc048cb11bf60db75f495608b231aa870d32f69c1db7e12e7abad3f3e2b568b8ecc2a9039e69f68002407938b603d88892386869d1c26763abb1ef09e261adad94ab04f2c270f9653249f2140534662244636afc1ca539e305ede8d098a1579187df198b692e2329d078d6c9a86c9b535835b1862255a46d8ca1d61786fc8f077357464033fdca9dd1121542651fc1b430a98f7a480f75e14f04b46669902057b45bcdfe2df8902b2ed291408840262b8e614f0730a042258411e9e180ad6642e586372ffde3207db987e0ae204bc5955264e70b7c40986c411f3c4040554c13854410155300a55b037dd0e11ddde0d38423fccb4beaebb696c570ddeb08b91a85628e2188a785b9b189af82aae34ad6bdac456175f92faf9f251a5c7b7c79bbf2cadcf48e3f345dbe7a472a1cb695da85011357c029713ea2979fa14b5a894169839acabccb7482f8ac1b40ff3525b132acc50335a1462b24a65b5a9202b732eef4332ecb86723204db1b62f62bd93a9750bb103a7606c786c166e71eee6d8669fa463b0ed48f4dbb890b2a5a7ed677b91a4b8455dbad36c0576b15110b3119f3c00b80e12e16398f028e74d4cb928991911846411c17b6409440224a62a5f3e0592dd860e52eee720723f86d96db12ca4b171c08df9d60819bb92fb4dad937b2c09abb4880a416421e53e23b018c373900066d6a3f2f343033acefeb3b27ab69a5bd055f0346cc446241bbd60e648783a2a5299a2f154a6a848658a46539922b1a0ab2c4231d21cb96f31a1237cd3c57a6b08df355d06f08bda02b88b895f1d9af85501b05a40511d85a25ab6e3e3ec9ade1b0a36d1ae4b8ad05f05aeb3e59536d06f5b370ac1f2531f6ff8d9eda08187d03d590a2507a7824f432b3fac60617f69ab83d078f6c0d3a41dddec9760c9d61c62de06dcd9f8e53a7f56ea229d0c9a357b56aafb0b076015bb374ecdfdc4a9b92f00cae8d4833f6efc5388844cfcf84fa0041b32b432046464b2b20a5f6fe0838a1f7b2df54f7dce5ac9e4d8d35d01f1bbbb43ec7daf408cb0f7f3d97698588b441bf57c09fcef75840bba35fbb7a43e0b101deee5f71044fd116cc3c078dafe261eb168cae3adaffcc2effd36cdd6b1cfdbdbd7bffa07a4ad38977bd63946170d15492ef243eecca3a4402bd91faa0e7ff49a3cb61f1b07d2f08b6c94e51dddfaa33ffbe3dfffd3975ea8bbd2cf1374716679ab9ffa3989404adecf5427f9f87e60a722cabb92912ffc94547c0a8197bc62c2151d5bf131ba7850ee5c71ee5152ce7945c390695bf125bab063c1159f434427af482efe89fde2465b0bf79f6a6d435f64aa4aca5c4352e61ada399968f6d8d905931812397ce425bbc84af793f93c77a3f0fc38f345de2302e7647a53a2747477815d5f0e9475ecd7d4e449dda3d82c4426998e7fc1255de0162e3d6b82221a37025c9813dd82687d44cd1ea9965d4fb920e69183a76dd04ce0cbdd265c108dd2e2946971ca347bade4978967069e88ecc73a27c585e04850f2899283aac941e53a122b60b3ea31d47b38b77fd0040b95ce23715c48ee2c2254d0f0261fe890bc3d2663714cb4f0ee22b0267c218b615dc2b0e6ef7a24c8170bec689cf76233766dc24e50f47732e57cf9488a1057c497d2e13d5a951c2222d1afe1cb6b97bfbcb61fdf449cc2108d53a852a2907c1d6f29a489429150482e40764e214e52420153a82aa4e14b2497a9e1d7095599974d0eba7401df349f2ede00cb21cf526a90255a7ac95f1745f957d81cb7c9730ae6879cab395a09339ca34508084523ecccc8ae30b22bcc588cb7ca3dfc6557553e1221a47d4b6b49d854163b48678c38425404551d7e2c7d0ce5d9d0889ce5b057338e6aefd172cf26c9ae4d583cf84b21a34b95a7f6ea3f9ea359b11f43ebc78194e88ae33ddc93387302dc712b44923129de7bc051783b6f9358f7a9e19919f371ec2b6c052cf02918b2d3cb4a070cfe08561ba371f14e8379b036680cd969804a5b9dd4b15f1ce1d36104fefdb8a2b237c9a54afe41fc0d4fd7b6863bfc002a4bebe407258d3a7572ebb7492612fb4933ef3a90a71e5e489d2ed2200ee02b917c2fe4cd8cecb7ba64b293e1b3d2c9b649a058f2aa1cf3789d1ed33a619c599d5e42b2eb6c1fc5fc2559a3ea16b93b7553ed30799e1626f24c9d1a1b34e7dd944b178cd3cb13ec4859309fc947f7e55a5ca5ce270f60f745254ff0c11671f657387803db1488a4bf45b36cc476ba8d1a66f35696ba8d3a6b903a2693eaecaa0445047631089fccaefd3bc5f9b6afcb4bc4303bdc2e600319d907f5ec356b7757e4bc8f0b8c940a072ae4f39667d9a8ac6363bf6ac26791b3285f8151159a1be1531a767826c659908fc1d014ff160e506f82959e4dd8309cd350370eb214c8c490f642488c473a04c3a1bfba0917f87b64629b585c014e13263cf2a92f4f1e48eb8257577c423d118b2a1e4132d203e22ff83ada126e68f0570e1931e89f2ef06775d3ce3b0b9ccdf838e7db64fa04bd912f1707d8e3786701892a6e61ac840434f9bdc9c5d4799c4cef1fe6ef51a2ee49b59fdd4c2a643789bbe5cd33ac832efa811ffc943d18896887c46a8214db7091f47a218db2cfffa5bb9006073a46cac20b2638d03561172ba7405d5704a95b82d405a4ae40ea087b3b321629363ec505092750b43e3cf61cfc88c067db00899f04be8900ba63410f80f3803342087237dbc22699853c2d31967c2f025df6641d1f67614db89a33a27c781a8081eca7c3c20ed96ba4a94cd4e949069b930719a44e8929482110efa7f6fc0c2de9fcace60186086fbe4bc8ec1f19497c8926e2fc5b96a800db15c8749ef644d3b2f3a08537b55fb21610860584ba80d0d6294398c3981fe342d0dcca10d0e174f6092c100414ab0027fe31d295a2992cefbdbe90ca9135d99bff5652af6a2bb8250dc5a17dfb8903a9635ebd90fdff34731ce0e00736d76bbc11cbf23701acbd564c6708aa4bba3afef7e194cabcd87bf7538a0cb8240e5eed2d0fda746b621f67ecd1dd4c3c593a17ceb42f6d2ce3d43c4516497c74fd6ca7bd38301b38b76fc90cd60db5eee194b526eefb666d7d6096d70c9f32d65c5e5b6a5fa286d774fc413933adb5625e5d1e9c3783f36dd36fadb60dce1ea32196fb833e55bca5e34679ec9113e0ccb9d6f20a03f2ae8e8f2eb6d630d622013d681b01ab8fa107adb32b5470ceb4e8fdf976cf2caeaf519f8b036a78d3893fd0eef5d67b06735d5e7ba5001fd37362ef9536aadd76e2fb0a7831cab9f5cd350cbcedc647e574367a61674b8061362b34058c8ad3fadafd012127afc81373e38fee40ddeafad2f2b9cb7b62ed5db798a5ad79b7b3f4e28fcb718b546cd174beb5f64adb6cf496575bbdcba6dbbe6c5e3ddf5eb31d0311ad35db3d46f6e207d7d68b0950b3f6c5e5f5cdfeca65f30a5167c044b8e6c70fe6bdf3917366736db0bc22fc22074e12287efcc88e7913c6fbed3da77ddb8f8fd87ea5e25dcefa7a101f6f0d06edd58d01069163296d63a2506b605e6ded35641047eb6bcc18987d10077c561e081ec627ed1b233cf9e1be69f5a8df354206e16679c99cebadaf320ba4dc2a35fdcb44ff55aabf486c03d068358e3f2d67e93551bb89e3f496189ea5767ff99535c10bf779b64de85c5901ff9c03019be7dbad0daebab8bebab1d21eb4a9923ddb1074181e8caabeb6b8cbc1a85fa8e625ffb3b28f013c71a27576f9070f7f74c6ad59109b90420f3f13a543068bc2d81e215914d492f1b30895e353a5f5de488fd561adf17e03e95279a16db01394f2a8f5114570eff8e89ffd510f6782fa672f0fdafd7bcba75e12f296da3d1554875cb90f279c3637d69769b07e7bb1de6b9f6b9e2574769b6b9babc126151f9e4d8685c4fbe7962f71f9d123f7adb62e35d7da83e666bff54abbf9ea7a6fa95fb9d8eac9db046f173736eddbd5be8ff2b9ca527ba57519630dab56da9706ed35cc4039dee0f246db0e5c9a475c9ac72e737227c1fbd40d33933d1ab51fc9f5b33f1a8e17dc5feea218bd1f17b7d4c69e6ba9bcc8de8dd2b39eebc421912a85badb8528eea077a93c8bd01e6ca9dcbbeedfe7ae76b05ba53824f3ffbcaf40b4d777d051242dc62a7a3850f33be905f537076315a3bca272766973ef466bb1db5e2a1fec2accb2af586c48a99d5b57ee2e8def2b68d52c6aef42b57bf16a79e90eb57c9cc4dcfacc9fdd28980605ab9feb1e2fa8acd473ffe2db9ff9c6374bc4c2df13ffe45f6539cef96ffba51f11dce194e7df3d437dbacbb3c7f4e80cb4378a1d75eb775768b03d6b8e2384fefe37767b3f41","transaction":{"expiration":"2018-06-09T11:58:27","ref_block_num":175,"ref_block_prefix":3727162172,"max_net_usage_words":0,"max_cpu_usage_ms":0,"delay_sec":0,"context_free_actions":[],"actions":[{"account":"eosio","name":"setcode","authorization":[{"actor":"eosio.msig","permission":"active"}],"data":{"account":"eosio.msig","vmtype":0,"vmversion":0,"code":"0061736d010000000198011760017f0060047f7e7e7f0060047f7e7e7e006000006000017e60027e7e0060067f7f7f7f7f7f017f60027f7f0060037f7f7f017f60077e7e7f7f7f7f7e017f6000017f60027f7f017f60017e0060047e7e7e7e017f60067e7e7e7e7f7f017f60047f7e7f7f0060057f7e7f7f7f0060057e7e7f7f7e017f60047f7f7f7f0060037f7e7f017f60047f7f7e7f0060037e7e7e0060017f017f02f7021203656e760561626f7274000303656e7610616374696f6e5f646174615f73697a65000a03656e761e636865636b5f7065726d697373696f6e5f617574686f72697a6174696f6e000903656e761f636865636b5f7472616e73616374696f6e5f617574686f72697a6174696f6e000603656e761063757272656e745f7265636569766572000403656e760c63757272656e745f74696d65000403656e760b64625f66696e645f693634000d03656e760a64625f6765745f693634000803656e760d64625f72656d6f76655f693634000003656e760c64625f73746f72655f693634000e03656e760d64625f7570646174655f693634000f03656e760c656f73696f5f617373657274000703656e76066d656d637079000803656e76076d656d6d6f7665000803656e7610726561645f616374696f6e5f64617461000b03656e760c726571756972655f61757468000c03656e760d726571756972655f6175746832000503656e760d73656e645f646566657272656400100336350b0b0b0a000807070b0b0b0b0b0b0b0b11000b0b0b071207120b07011314070b01140213070702150b0b0b0716000008160b16000304050170010606050301000107f10511066d656d6f72790200165f5a6571524b3131636865636b73756d32353653315f0012165f5a6571524b3131636865636b73756d31363053315f0013165f5a6e65524b3131636865636b73756d31363053315f0014036e6f770015305f5a4e35656f73696f3132726571756972655f6175746845524b4e535f31367065726d697373696f6e5f6c6576656c450016b0015f5a4e35656f73696f3331636865636b5f7472616e73616374696f6e5f617574686f72697a6174696f6e45524b4e535f31317472616e73616374696f6e45524b4e5374335f5f3133736574494e535f31367065726d697373696f6e5f6c6576656c454e53335f346c6573734953355f45454e53335f39616c6c6f6361746f724953355f45454545524b4e53345f4931307075626c69635f6b65794e53365f4953445f45454e53385f4953445f4545454500179f015f5a4e35656f73696f3330636865636b5f7065726d697373696f6e5f617574686f72697a6174696f6e457979524b4e5374335f5f31337365744931307075626c69635f6b65794e53305f346c6573734953325f45454e53305f39616c6c6f6361746f724953325f45454545524b4e53315f494e535f31367065726d697373696f6e5f6c6576656c454e53335f4953415f45454e53355f4953415f454545457900221c5f5a4e35656f73696f386d756c74697369673770726f706f7365457600233b5f5a4e35656f73696f386d756c746973696737617070726f766545794e535f346e616d65454e535f31367065726d697373696f6e5f6c6576656c45002d3d5f5a4e35656f73696f386d756c746973696739756e617070726f766545794e535f346e616d65454e535f31367065726d697373696f6e5f6c6576656c450032255f5a4e35656f73696f386d756c74697369673663616e63656c45794e535f346e616d6545790034235f5a4e35656f73696f386d756c7469736967346578656345794e535f346e616d6545790038056170706c790039066d656d636d700041066d616c6c6f63004204667265650045090c010041000b064638322d34230ab09701350b002000200141201041450b0b002000200141201041450b0d0020002001412010414100470b0a00100542c0843d80a70b0e002000290300200029030810100bba0503027f017e077f4100410028020441306b220c360204200c41106a200010184100210b410021094100210a024020022802082203450d0041002106200c4100360208200c42003703002003ad21050340200641016a2106200542078822054200520d000b02400240024020022802002207200241046a2204460d00034002400240200722082802042200450d0003402000220728020022000d000c020b0b200828020822072802002008460d00200841086a210803402008280200220041086a2108200020002802082207280200470d000b0b200641226a210620072004470d000b2006450d010b200c20061019200c2802042107200c28020021000c010b41002107410021000b200c2000360224200c2000360220200c2007360228200c41206a2002101a1a200c2802042109200c280200210a0b41002100024020012802082202450d0041002106200c4100360208200c42003703002002ad21050340200641016a2106200542078822054200520d000b02400240024020012802002207200141046a2204460d00034002400240200722082802042200450d0003402000220728020022000d000c020b0b200828020822072802002008460d00200841086a210803402008280200220041086a2108200020002802082207280200470d000b0b200641106a210620072004470d000b2006450d010b200c20061019200c2802042107200c28020021000c010b41002107410021000b200c2000360224200c2000360220200c2007360228200c41206a2001101b1a200c280204210b200c28020021000b200c2802102207200c28021420076b200a410020031b2009200a6b410020031b2000410020021b200b20006b410020021b1003210702402000450d002000103f0b0240200a450d00200a103f0b0240200c2802102200450d00200c20003602142000103f0b4100200c41306a360204200741004a0b9e0101037f4100410028020441106b22043602044100210320004100360208200042003702002004410036020020042001101c1a0240024020042802002202450d00200020021019200041046a2802002103200028020021000c010b410021000b20042000360204200420003602002004200336020820042001101d1a2004200141186a101e200141246a101e200141306a101f1a4100200441106a3602040bad0201057f0240024002400240024020002802082202200028020422066b20014f0d002006200028020022056b220320016a2204417f4c0d0241ffffffff0721060240200220056b220241feffffff034b0d0020042002410174220620062004491b2206450d020b2006103e21020c030b200041046a21000340200641003a00002000200028020041016a22063602002001417f6a22010d000c040b0b41002106410021020c010b20001040000b200220066a2104200220036a220521060340200641003a0000200641016a21062001417f6a22010d000b2005200041046a2203280200200028020022016b22026b2105024020024101480d00200520012002100c1a200028020021010b2000200536020020032006360200200041086a20043602002001450d002001103f0f0b0be20203027f017e057f4100410028020441d0006b22093602042000280204210520013502082104200041086a2103200041046a210703402004a721062009200442078822044200522208410774200641ff0071723a0028200328020020056b41004a4110100b2007280200200941286a4101100c1a2007200728020041016a220536020020080d000b024020012802002206200141046a2201460d00200041086a2102200041046a21030340200941066a20062208410d6a4122100c1a200941286a200941066a4122100c1a200228020020056b41214a4110100b2003280200200941286a4122100c1a2003200328020041226a22053602000240024020082802042207450d0003402007220628020022070d000c020b0b200828020822062802002008460d00200841086a210803402008280200220741086a2108200720072802082206280200470d000b0b20062001470d000b0b4100200941d0006a36020420000bed0203017f017e057f4100410028020441106b22083602042000280204210420013502082103200041086a2102200041046a210603402003a721052008200342078822034200522207410774200541ff0071723a000f200228020020046b41004a4110100b20062802002008410f6a4101100c1a2006200628020041016a220436020020070d000b024020012802002205200141046a2201460d00200041046a21020340200041086a220628020020046b41074a4110100b20022802002005220741106a4108100c1a2002200228020041086a2205360200200628020020056b41074a4110100b2002280200200741186a4108100c1a2002200228020041086a22043602000240024020072802042206450d0003402006220528020022060d000c020b0b200728020822052802002007460d00200741086a210703402007280200220641086a2107200620062802082205280200470d000b0b20052001470d000b0b4100200841106a36020420000b9d0502067f017e200020002802002206410a6a3602002006410b6a2106200135020c21080340200641016a2106200842078822084200520d000b20002006360200200135021421080340200641016a2106200842078822084200520d000b200020063602002001411c6a2802002202200128021822076b41286dad21080340200641016a2106200842078822084200520d000b20002006360200024020072002460d000340200641106a2106200741146a2802002203200728021022046b2205410475ad21080340200641016a2106200842078822084200520d000b024020042003460d00200541707120066a21060b2006200741206a28020022036a200728021c22046b2106200320046bad21080340200641016a2106200842078822084200520d000b200741286a22072002470d000b200020063602000b200141286a2802002202200128022422076b41286dad21080340200641016a2106200842078822084200520d000b20002006360200024020072002460d000340200641106a2106200741146a2802002203200728021022046b2205410475ad21080340200641016a2106200842078822084200520d000b024020042003460d00200541707120066a21060b2006200741206a28020022036a200728021c22046b2106200320046bad21080340200641016a2106200842078822084200520d000b200741286a22072002470d000b200020063602000b200141346a2802002205200128023022076b410475ad21080340200641016a2106200842078822084200520d000b20002006360200024020072005460d0003402006200741086a28020022036a41026a200728020422046b2106200320046bad21080340200641016a2106200842078822084200520d000b200741106a22072005470d000b200020063602000b20000b9f0303047f017e017f4100410028020441106b2207360204200028020820002802046b41034a4110100b200028020420014104100c1a2000200028020441046a2204360204200028020820046b41014a4110100b2000280204200141046a4102100c1a2000200028020441026a2204360204200028020820046b41034a4110100b2000280204200141086a4104100c1a2000200028020441046a2205360204200135020c210603402006a721042007200642078822064200522202410774200441ff0071723a000e200041086a28020020056b41004a4110100b200041046a22042802002007410e6a4101100c1a2004200428020041016a220536020020020d000b200041086a220328020020056b41004a4110100b200041046a2204280200200141106a4101100c1a2004200428020041016a22053602002001350214210603402006a721022007200642078822064200522201410774200241ff0071723a000f200328020020056b41004a4110100b20042802002007410f6a4101100c1a2004200428020041016a220536020020010d000b4100200741106a36020420000bbe0203037f017e037f4100410028020441106b2208360204200128020420012802006b41286dad210520002802042106200041086a2103200041046a210403402005a721072008200542078822054200522202410774200741ff0071723a000f200328020020066b41004a4110100b20042802002008410f6a4101100c1a2004200428020041016a220636020020020d000b024020012802002207200141046a2802002203460d00200041046a21040340200041086a220228020020066b41074a4110100b200428020020074108100c1a2004200428020041086a2206360200200228020020066b41074a4110100b2004280200200741086a4108100c1a2004200428020041086a3602002000200741106a10212007411c6a10201a200741286a22072003460d01200428020021060c000b0b4100200841106a36020420000b880203037f017e027f4100410028020441106b2207360204200128020420012802006b410475ad210520002802042106200041086a210303402005a721042007200542078822054200522202410774200441ff0071723a000f200328020020066b41004a4110100b200041046a22042802002007410f6a4101100c1a2004200428020041016a220636020020020d000b024020012802002204200141046a2802002202460d00200041086a21030340200328020020066b41014a4110100b200041046a220628020020044102100c1a2006200628020041026a3602002000200441046a10201a200441106a22042002460d01200628020021060c000b0b4100200741106a36020420000bda0103057f017e017f4100410028020441106b2208360204200128020420012802006bad210720002802042106200041086a2104200041046a210503402007a721022008200742078822074200522203410774200241ff0071723a000f200428020020066b41004a4110100b20052802002008410f6a4101100c1a2005200528020041016a220636020020030d000b200041086a28020020066b200141046a280200200128020022026b22054e4110100b200041046a220628020020022005100c1a2006200628020020056a3602004100200841106a36020420000ba40203027f017e037f4100410028020441106b2207360204200128020420012802006b410475ad210420002802042105200041086a210203402004a721032007200442078822044200522206410774200341ff0071723a000f200228020020056b41004a4110100b200041046a22032802002007410f6a4101100c1a2003200328020041016a220536020020060d000b024020012802002206200141046a2802002201460d00200041046a21030340200041086a220228020020056b41074a4110100b200328020020064108100c1a2003200328020041086a2205360200200228020020056b41074a4110100b2003280200200641086a4108100c1a2003200328020041086a2205360200200641106a22062001470d000b0b4100200741106a36020420000b900503027f017e087f4100410028020441206b220f3602044100210e4100210c4100210d024020022802082205450d0041002108200f4100360208200f42003703002005ad21070340200841016a2108200742078822074200520d000b02400240024020022802002209200241046a2206460d000340024002402009220b280204220a450d000340200a2209280200220a0d000c020b0b200b2802082209280200200b460d00200b41086a210b0340200b280200220a41086a210b200a200a2802082209280200470d000b0b200841226a210820092006470d000b2008450d010b200f20081019200f2802042109200f280200210a0c010b410021094100210a0b200f200a360214200f200a360210200f2009360218200f41106a2002101a1a200f280204210c200f280200210d0b4100210a024020032802082202450d0041002108200f4100360208200f42003703002002ad21070340200841016a2108200742078822074200520d000b02400240024020032802002209200341046a2206460d000340024002402009220b280204220a450d000340200a2209280200220a0d000c020b0b200b2802082209280200200b460d00200b41086a210b0340200b280200220a41086a210b200a200a2802082209280200470d000b0b200841106a210820092006470d000b2008450d010b200f20081019200f2802042109200f280200210a0c010b410021094100210a0b200f200a360214200f200a360210200f2009360218200f41106a2003101b1a200f280204210e200f280200210a0b20002001200d410020051b200c200d6b410020051b200a410020021b200e200a6b410020021b2004100221090240200a450d00200a103f0b0240200d450d00200d103f0b4100200f41206a360204200941004a0b990e05027f017e027f017e017f4100410028020441f0016b220236020420022207100122053602c801024002402005418104490d002005104221020c010b410020022005410f6a4170716b22023602040b200720023602c40120022005100e1a200742003703b00141002105200741003602a801200742003703a001100521062007410036029401200741003a0098012007410036029c012007200642c0843d80a7413c6a36028801200720072802c401220236027c20072802c8012101200720023602782007200220016a36028001200141074b4120100b200741b8016a200728027c4108100c1a2007200728027c41086a220236027c20072802800120026b41074b4120100b200741b0016a200728027c4108100c1a2007200728027c41086a36027c200741f8006a200741a0016a10241a2007200728027c20072802786b360274200741f8006a20074188016a10251a20072903b801100f10052106200728028801200642c0843d80a74f4130100b20072903b80121062007200029030022033703482007427f3703582007410036026020072006370350200741e4006a4100360200200741e8006a41003602004100210202402003200642808080888dccd6f4ad7f20072903b001100622014100480d00200741c8006a200110262202280214200741c8006a4641d000100b0b200245419001100b200742003703382007410036024020072802a401220220072802a00122016b2204410475ad210603402005417f6a2105200642078822064200520d000b024002400240024020012002460d00200441707122022005470d0141002102410021050c030b410020056b21050c010b200220056b21050b200741386a20051019200728023c2102200728023821050b200720053602142007200536021020072002360218200741106a200741a0016a10211a20072802c401200728027422056a20072802c80120056b4100410020072802382205200728023c20056b100341004a41c001100b20072903b80121062007200741c4016a3602142007200741b0016a3602102007200741f4006a3602182007200741c8016a36021c200720063703e801200729034810045141f001100b2007200741106a3602d4012007200741c8006a3602d0012007200741e8016a3602d8014120103e2205420037030020054200370208200541003602102005200741c8006a360214200741d0016a20051027200720053602082007200529030022063703d0012007200528021822013602e00102400240200741e4006a22042802002202200741c8006a41206a2802004f0d00200220063703082002200136021020074100360208200220053602002004200241186a3602000c010b200741e0006a200741086a200741d0016a200741e0016a10280b200728020821052007410036020802402005450d00024020052802082202450d002005410c6a20023602002002103f0b2005103f0b2007427f3703202007410036022820072903b8012106200720002903002203370310200720063703182007412c6a22024100360200200741306a220141003602002007200741a0016a36020c2007200741b0016a360208200720063703e801200310045141f001100b2007200741086a3602d4012007200741106a3602d0012007200741e8016a3602d8014130103e220542003703002005420037020820054200370210200542003702182005200741106a360220200741d0016a20051029200720053602e0012007200529030022063703d0012007200528022422043602cc01024002402002280200220020012802004f0d002000200637030820002004360210200741003602e001200020053602002002200041186a3602000c010b200741286a200741e0016a200741d0016a200741cc016a102a0b20072802e0012105200741003602e00102402005450d00024020052802142200450d00200541186a20003602002000103f0b024020052802082200450d002005410c6a20003602002000103f0b2005103f0b024020072802282201450d00024002402007412c6a220428020022002001460d000340200041686a220028020021052000410036020002402005450d00024020052802142202450d00200541186a20023602002002103f0b024020052802082202450d002005410c6a20023602002002103f0b2005103f0b20012000470d000b200741286a28020021050c010b200121050b200420013602002005103f0b024020072802382205450d002007200536023c2005103f0b024020072802602201450d0002400240200741e4006a220428020022052001460d000340200541686a220528020021002005410036020002402000450d00024020002802082202450d002000410c6a20023602002002103f0b2000103f0b20012005470d000b200741e0006a28020021050c010b200121050b200420013602002005103f0b024020072802a0012205450d00200720053602a4012005103f0b4100200741f0016a3602040bd00203037f017e027f200028020421074100210642002105200041086a2102200041046a21030340200720022802004941d002100b200328020022072d000021042003200741016a2207360200200441ff0071200641ff0171220674ad2005842105200641076a210620044107760d000b0240024002402005a7220420012802042202200128020022076b41047522064d0d002001200420066b102c20012802002207200141046a2802002202470d010c020b0240200420064f0d00200141046a200720044104746a22023602000b20072002460d010b200041046a220428020021060340200041086a220328020020066b41074b4120100b200720042802004108100c1a2004200428020041086a2206360200200328020020066b41074b4120100b200741086a20042802004108100c1a2004200428020041086a2206360200200741106a22072002470d000b0b20000b810303037f017e027f200028020820002802046b41034b4120100b200120002802044104100c1a2000200028020441046a2202360204200028020820026b41014b4120100b200141046a20002802044102100c1a2000200028020441026a2202360204200028020820026b41034b4120100b200141086a20002802044104100c1a2000200028020441046a2204360204410021064200210503402004200041086a2802004941d002100b200041046a220728020022042d000021022007200441016a2204360200200241ff0071200641ff0171220674ad2005842105200641076a210620024107760d000b200120053e020c200041086a22032802002004474120100b200141106a200041046a22042802004101100c1a2004200428020041016a220636020041002107420021050340200620032802004941d002100b200428020022062d000021022004200641016a2206360200200241ff0071200741ff0171220774ad2005842105200741076a210720024107760d000b200120053e021420000b8b0403037f017e047f410028020441306b220921084100200936020402402000411c6a280200220720002802182202460d00410020026b2103200741686a21060340200641106a2802002001460d0120062107200641686a22042106200420036a4168470d000b0b0240024020072002460d00200741686a28020021060c010b20014100410010072206411f7641017341b002100b024002402006418104490d002006104221040c010b410020092006410f6a4170716b22043602040b20012004200610071a20082004360224200820043602202008200420066a2207360228024020064180044d0d0020041045200841286a2802002107200828022421040b4120103e22064200370300200642003702082006410036021020062000360214200720046b41074b4120100b200620044108100c1a2008200441086a360224200841206a200641086a102b1a200620013602182008200636021820082006290300220537031020082006280218220736020c024002402000411c6a22012802002204200041206a2802004f0d00200420053703082004200736021020084100360218200420063602002001200441186a3602000c010b200041186a200841186a200841106a2008410c6a10280b20082802182104200841003602182004450d00024020042802082207450d002004410c6a20073602002007103f0b2004103f0b4100200841306a36020420060be40403077f017e017f4100410028020441106b22053602042000280200210220012000280204220628020029030037030020062802042802002107410021082005220a4100360200200a410036020420062802082103200a41003602084100210402400240200628020c280200200328020022036b2206450d002006417f4c0d01200a41086a2006103e220420066a2208360200200a20043602002004200720036a2006100c1a200a20083602040b0240024020012802082206450d002001410c6a20063602002006103f200141106a22064100360200200141086a42003702000c010b200141106a21060b200620083602002001410c6a2008360200200141086a2004360200200841086a20046b2106200820046bad21090340200641016a2106200942078822094200520d000b024002402006418104490d002006104221070c010b410020052006410f6a4170716b22073602040b200641074a4110100b200720014108100c1a2001410c6a280200200141086a2802006bad2109200741086a2108200720066a210303402009a72104200a200942078822094200522205410774200441ff0071723a000f200320086b41004a4110100b2008200a410f6a4101100c1a200841016a210820050d000b200320086b2001410c6a280200200141086a28020022046b22054e4110100b200820042005100c1a2001200229030842808080888dccd6f4ad7f20002802082903002001290300220920072006100936021802402006418104490d00200710450b024020092002290310540d00200241106a427e200942017c2009427d561b3703000b4100200a41106a3602040f0b200a1040000bc50301047f024002402000280204200028020022066b41186d220441016a220541abd5aad5004f0d0041aad5aad500210702400240200028020820066b41186d220641d4aad52a4b0d0020052006410174220720072005491b2207450d010b200741186c103e21060c020b41002107410021060c010b20001040000b20012802002105200141003602002006200441186c6a2201200536020020012002290300370308200120032802003602102006200741186c6a2104200141186a210502400240200041046a280200220620002802002207460d000340200641686a2202280200210320024100360200200141686a2003360200200141786a200641786a280200360200200141746a200641746a280200360200200141706a200641706a280200360200200141686a21012002210620072002470d000b200041046a2802002107200028020021020c010b200721020b20002001360200200041046a2005360200200041086a2004360200024020072002460d000340200741686a220728020021012007410036020002402001450d00024020012802082206450d002001410c6a20063602002006103f0b2001103f0b20022007470d000b0b02402002450d002002103f0b0b880403057f017e047f410028020441106b220b210a4100200b36020420002802002102200120002802042208280200290300370300200141086a2103200828020421080240024020012802082205450d002001410c6a220920053602002005103f200141106a22054100360200200342003702000c010b200141106a21052001410c6a21090b200320082802003602002009200828020436020020052008280208360200200842003702002008410036020820092802002209200328020022046b2206410475ad2107410821080340200841016a2108200742078822074200520d000b200141146a2105024020042009460d00200641707120086a21080b200141186a2802002209200528020022046b2206410475ad21070340200841016a2108200742078822074200520d000b024020042009460d00200641707120086a21080b024002402008418104490d002008104221090c010b4100200b2008410f6a4170716b22093602040b200a2009360200200a200920086a360208200841074a4110100b200920014108100c1a200a200941086a360204200a200310211a200a200510211a20012002290308428080808ecdcddeb53520002802082903002001290300220720092008100936022402402008418104490d00200910450b024020072002290310540d00200241106a427e200742017c2007427d561b3703000b4100200a41106a3602040be00301047f024002402000280204200028020022066b41186d220441016a220541abd5aad5004f0d0041aad5aad500210702400240200028020820066b41186d220641d4aad52a4b0d0020052006410174220720072005491b2207450d010b200741186c103e21060c020b41002107410021060c010b20001040000b20012802002105200141003602002006200441186c6a2201200536020020012002290300370308200120032802003602102006200741186c6a2104200141186a210502400240200041046a280200220620002802002207460d000340200641686a2202280200210320024100360200200141686a2003360200200141786a200641786a280200360200200141746a200641746a280200360200200141706a200641706a280200360200200141686a21012002210620072002470d000b200041046a2802002107200028020021020c010b200721020b20002001360200200041046a2005360200200041086a2004360200024020072002460d000340200741686a220728020021012007410036020002402001450d00024020012802142206450d00200141186a20063602002006103f0b024020012802082206450d002001410c6a20063602002006103f0b2001103f0b20022007470d000b0b02402002450d002002103f0b0b820203047f017e017f200028020421054100210742002106200041086a2102200041046a21030340200520022802004941d002100b200328020022052d000021042003200541016a2205360200200441ff0071200741ff0171220774ad2006842106200741076a210720044107760d000b024002402006a7220320012802042207200128020022046b22024d0d002001200320026b1019200041046a2802002105200141046a2802002107200128020021040c010b200320024f0d00200141046a200420036a22073602000b200041086a28020020056b200720046b22054f4120100b2004200041046a22072802002005100c1a2007200728020020056a36020020000bab0201067f0240024002400240024020002802082202200028020422076b41047520014f0d002007200028020022066b410475220320016a22044180808080014f0d0241ffffffff0021050240200220066b220241047541feffff3f4b0d0020042002410375220520052004491b2205450d0220054180808080014f0d040b2005410474103e2102200041046a2802002107200028020021060c040b200041046a200720014104746a3602000f0b41002105410021020c020b20001040000b1000000b200220034104746a2203200720066b22076b2104200320014104746a2101200220054104746a2105024020074101480d00200420062007100c1a200028020021060b20002004360200200041046a2001360200200041086a200536020002402006450d002006103f0b0bfa0202027e047f4100410028020441c0006b220936020420032903002204200329030822051010200941386a4100360200200920013703202009427f37032820094200370330200920002903003703180240200941186a200241e002102e220628020822002006410c6a2802002208460d000340024020002903002004520d00200041086a2903002005510d020b2008200041106a2200470d000b200821000b200920003602102000200847418003100b200920033602082009200941106a36020c200941186a20062001200941086a102f024020092802302206450d0002400240200941346a220728020022082006460d000340200841686a220828020021002008410036020002402000450d00024020002802142203450d00200041186a20033602002003103f0b024020002802082203450d002000410c6a20033602002003103f0b2000103f0b20062008470d000b200941306a28020021000c010b200621000b200720063602002000103f0b4100200941c0006a3602040bb50101057f02402000411c6a280200220720002802182203460d00200741686a2106410020036b2104034020062802002903002001510d0120062107200641686a22052106200520046a4168470d000b0b0240024020072003460d00200741686a280200220628022020004641d000100b0c010b4100210620002903002000290308428080808ecdcddeb5352001100622054100480d00200020051031220628022020004641d000100b0b20064100472002100b20060ba00404017e047f017e037f410028020441106b220c210b4100200c360204200128022020004641c003100b200029030010045141f003100b2003280200210a2001290300210402400240200141186a220728020022052001411c6a280200460d002005200a290300370300200541086a200a41086a2903003703002007200728020041106a3602000c010b200141146a200a10300b02402001410c6a220a2802002003280204280200220341106a22056b22064104752207450d00200320052006100d1a0b200a200320074104746a360200200420012903005141b004100b41082103200141086a2105200a280200220a200128020822066b2208410475ad21090340200341016a2103200942078822094200520d000b200141146a210702402006200a460d00200841707120036a21030b200141186a280200220a200728020022066b2208410475ad21090340200341016a2103200942078822094200520d000b02402006200a460d00200841707120036a21030b024002402003418104490d0020031042210a0c010b4100200c2003410f6a4170716b220a3602040b200b200a360200200b200a20036a360208200341074a4110100b200a20014108100c1a200b200a41086a360204200b200510211a200b200710211a20012802242002200a2003100a02402003418104490d00200a10450b024020042000290310540d00200041106a427e200442017c2004427d561b3703000b4100200b41106a3602040b9f0201067f02400240024020002802042206200028020022056b410475220241016a22034180808080014f0d0041ffffffff00210402400240200028020820056b220741047541feffff3f4b0d0020032007410375220420042003491b2204450d0120044180808080014f0d030b2004410474103e2107200041046a2802002106200028020021050c030b41002104410021070c020b20001040000b1000000b200720024104746a22032001290300370300200341086a200141086a2903003703002003200620056b22016b2106200720044104746a2104200341106a2103024020014101480d00200620052001100c1a200028020021050b20002006360200200041046a2003360200200041086a200436020002402005450d002005103f0b0bba0403037f017e047f410028020441306b220921084100200936020402402000411c6a280200220720002802182202460d00410020026b2103200741686a21060340200641106a2802002001460d0120062107200641686a22042106200420036a4168470d000b0b0240024020072002460d00200741686a28020021060c010b20014100410010072206411f7641017341b002100b024002402006418104490d002006104221040c010b410020092006410f6a4170716b22043602040b20012004200610071a20082004360224200820043602202008200420066a2207360228024020064180044d0d0020041045200841286a2802002107200828022421040b4130103e2206420037030020064200370208200642003702102006420037021820062000360220200720046b41074b4120100b200620044108100c1a2008200441086a360224200841206a200641086a10241a200841206a200641146a10241a200620013602242008200636021820082006290300220537031020082006280224220736020c024002402000411c6a22012802002204200041206a2802004f0d00200420053703082004200736021020084100360218200420063602002001200441186a3602000c010b200041186a200841186a200841106a2008410c6a102a0b20082802182104200841003602182004450d00024020042802142207450d00200441186a20073602002007103f0b024020042802082207450d002004410c6a20073602002007103f0b2004103f0b4100200841306a36020420060bfa0202027e047f4100410028020441c0006b220936020420032903002204200329030822051010200941386a4100360200200920013703202009427f37032820094200370330200920002903003703180240200941186a200241e002102e22062802142200200641186a2802002208460d000340024020002903002004520d00200041086a2903002005510d020b2008200041106a2200470d000b200821000b20092000360210200020084741f004100b200920033602082009200941106a36020c200941186a20062001200941086a1033024020092802302206450d0002400240200941346a220728020022082006460d000340200841686a220828020021002008410036020002402000450d00024020002802142203450d00200041186a20033602002003103f0b024020002802082203450d002000410c6a20033602002003103f0b2000103f0b20062008470d000b200941306a28020021000c010b200621000b200720063602002000103f0b4100200941c0006a3602040ba30404017e047f017e037f410028020441106b220c210b4100200c360204200128022020004641c003100b200029030010045141f003100b2003280200210a20012903002104024002402001410c6a22072802002205200141106a280200460d002005200a290300370300200541086a200a41086a2903003703002007200728020041106a3602000c010b200141086a200a10300b0240200141186a220a2802002003280204280200220341106a22056b22064104752207450d00200320052006100d1a0b200a200320074104746a360200200420012903005141b004100b41082103200141086a21052001410c6a280200220a200128020822066b2208410475ad21090340200341016a2103200942078822094200520d000b200141146a210702402006200a460d00200841707120036a21030b200141186a280200220a200728020022066b2208410475ad21090340200341016a2103200942078822094200520d000b02402006200a460d00200841707120036a21030b024002402003418104490d0020031042210a0c010b4100200c2003410f6a4170716b220a3602040b200b200a360200200b200a20036a360208200341074a4110100b200a20014108100c1a200b200a41086a360204200b200510211a200b200710211a20012802242002200a2003100a02402003418104490d00200a10450b024020042000290310540d00200041106a427e200442017c2004427d561b3703000b4100200b41106a3602040bbe0401057f4100410028020441e0006b22083602042003100f200841286a41206a4100360200200820013703302008427f3703382008420037034020082000290300370328200841286a200241e00210352107024020032001510d002007410c6a280200210420072802082105100521032008410036020c200841003a0010200841003602142008200342c0843d80a7413c6a360200200820053602542008200536025020082004360258200841d0006a200810251a100521032008280200200342c0843d80a749419005100b0b200841206a4100360200200820013703082008427f37031020084200370318200820002903003703002008200241e002102e2100200841286a20071036200820001037024020082802182204450d00024002402008411c6a220628020022072004460d000340200741686a220728020021002007410036020002402000450d00024020002802142205450d00200041186a20053602002005103f0b024020002802082205450d002000410c6a20053602002005103f0b2000103f0b20042007470d000b200841186a28020021000c010b200421000b200620043602002000103f0b024020082802402204450d0002400240200841c4006a220628020022002004460d000340200041686a220028020021072000410036020002402007450d00024020072802082205450d002007410c6a20053602002005103f0b2007103f0b20042000470d000b200841c0006a28020021000c010b200421000b200620043602002000103f0b4100200841e0006a3602040bb60101057f02402000411c6a280200220720002802182203460d00200741686a2106410020036b2104034020062802002903002001510d0120062107200641686a22052106200520046a4168470d000b0b0240024020072003460d00200741686a280200220628021420004641d000100b0c010b410021062000290300200029030842808080888dccd6f4ad7f2001100622054100480d00200020051026220628021420004641d000100b0b20064100472002100b20060b800302017e067f200128021420004641b005100b200029030010045141e005100b02402000411c6a2205280200220720002802182203460d0020012903002102410020036b2106200741686a2108034020082802002903002002510d0120082107200841686a22042108200420066a4168470d000b0b200720034741a006100b200741686a210802400240200720052802002204460d00410020046b2103200821070340200741186a2208280200210620084100360200200728020021042007200636020002402004450d00024020042802082206450d002004410c6a20063602002006103f0b2004103f0b200741106a200741286a280200360200200741086a200741206a29030037030020082107200820036a4168470d000b2000411c6a28020022072008460d010b0340200741686a220728020021042007410036020002402004450d00024020042802082206450d002004410c6a20063602002006103f0b2004103f0b20082007470d000b0b2000411c6a2008360200200128021810080bb60302017e067f200128022020004641b005100b200029030010045141e005100b02402000411c6a2204280200220720002802182203460d0020012903002102410020036b2105200741686a2106034020062802002903002002510d0120062107200641686a22082106200820056a4168470d000b0b200720034741a006100b200741686a210802400240200720042802002206460d00410020066b2103200821060340200641186a2208280200210520084100360200200628020021072006200536020002402007450d00024020072802142205450d00200741186a20053602002005103f0b024020072802082205450d002007410c6a20053602002005103f0b2007103f0b200641106a200641286a280200360200200641086a200641206a29030037030020082106200820036a4168470d000b2000411c6a28020022072008460d010b0340200741686a220728020021062007410036020002402006450d00024020062802142205450d00200641186a20053602002005103f0b024020062802082205450d002006410c6a20053602002005103f0b2006103f0b20082007470d000b0b2000411c6a2008360200200128022410080bf00603057f017e027f4100410028020441a0016b220b3602042003100f41002108200b41f8006a41206a4100360200200b200137038001200b427f37038801200b420037039001200b2000290300370378200b41f8006a200241e0021035210a200b41d0006a41206a4100360200200b2001370358200b427f370360200b4200370368200b2000290300370350200b41d0006a200241e002102e210010052109200b4100360244200b41003a0048200b410036024c200b200942c0843d80a7413c6a360238200b200a2802082206360228200b200636022c200b200a410c6a280200360230200b41286a200b41386a10251a10052109200b280238200942c0843d80a74f4130100b200b4100360220200b4200370318200041186a2802002206200028021422056b2207410475ad2109200041146a210403402008417f6a2108200942078822094200520d000b024002400240024020052006460d00200741707122062008470d0141002106410021080c030b410020086b21080c010b200620086b21080b200b41186a20081019200b28021c2106200b28021821080b200b2008360204200b2008360200200b2006360208200b200410211a200a41086a22082802002206200a410c6a220528020020066b41004100200b2802182206200b28021c20066b100341004a41c001100b200b2001370308200b2002370300200b200320082802002208200528020020086b41001011200b41f8006a200a1036200b41d0006a200010370240200b2802182208450d00200b200836021c2008103f0b0240200b2802682206450d0002400240200b41ec006a2205280200220a2006460d000340200a41686a220a2802002108200a410036020002402008450d00024020082802142200450d00200841186a20003602002000103f0b024020082802082200450d002008410c6a20003602002000103f0b2008103f0b2006200a470d000b200b41e8006a28020021080c010b200621080b200520063602002008103f0b0240200b280290012206450d0002400240200b4194016a220528020022082006460d000340200841686a2208280200210a200841003602000240200a450d000240200a2802082200450d00200a410c6a20003602002000103f0b200a103f0b20062008470d000b200b4190016a28020021080c010b200621080b200520063602002008103f0b4100200b41a0016a3602040bfc0603027f047e017f4100410028020441e0006b220936020442002106423b210541e00621044200210703400240024002400240024020064206560d0020042c00002203419f7f6a41ff017141194b0d01200341a5016a21030c020b420021082006420b580d020c030b200341d0016a41002003414f6a41ff01714105491b21030b2003ad42388642388721080b2008421f83200542ffffffff0f838621080b200441016a2104200642017c2106200820078421072005427b7c2205427a520d000b024020072002520d0042002106423b210541f00621044200210703400240024002400240024020064204560d0020042c00002203419f7f6a41ff017141194b0d01200341a5016a21030c020b420021082006420b580d020c030b200341d0016a41002003414f6a41ff01714105491b21030b2003ad42388642388721080b2008421f83200542ffffffff0f838621080b200441016a2104200642017c2106200820078421072005427b7c2205427a520d000b2007200151418007100b0b0240024020012000510d0042002106423b210541e00621044200210703400240024002400240024020064206560d0020042c00002203419f7f6a41ff017141194b0d01200341a5016a21030c020b420021082006420b580d020c030b200341d0016a41002003414f6a41ff01714105491b21030b2003ad42388642388721080b2008421f83200542ffffffff0f838621080b200441016a2104200642017c2106200820078421072005427b7c2205427a520d000b20072002520d010b200920003703580240024002400240200242ffffffffd3cddeb535570d0020024280808080d4cddeb535510d0120024280808080c0a8a1d3c100510d02200242808080808080a0aad700520d04200941003602342009410136023020092009290330370228200941d8006a200941286a103c1a0c040b2002428080808094ccd6f4ad7f510d022002428080c0dae9dbd6e654520d03200941003602442009410236024020092009290340370218200941d8006a200941186a103b1a0c030b2009410036024c2009410336024820092009290348370210200941d8006a200941106a103b1a0c020b2009410036023c2009410436023820092009290338370220200941d8006a200941206a103c1a0c010b200941003602542009410536025020092009290350370208200941d8006a200941086a103a1a0b4100200941e0006a3602040b8c0101047f4100280204220521042001280204210220012802002101024010012203450d00024020034180044d0d002003104222052003100e1a200510450c010b410020052003410f6a4170716b220536020420052003100e1a0b200020024101756a210302402002410171450d00200328020020016a28020021010b200320011100004100200436020441010ba10303027f037e037f410028020441e0006b22092108410020093602042001280204210220012802002107024002400240024010012203450d002003418104490d012003104221010c020b410021010c020b410020092003410f6a4170716b22013602040b20012003100e1a0b200842003703182008420037031020082001360254200820013602502008200120036a3602582008200841d0006a3602302008200841106a360240200841c0006a200841306a103d02402003418104490d00200110450b200841106a41086a29030021052008413c6a2008412c6a280200360200200841306a41086a2201200841286a28020036020020082903102104200820082802203602302008200841246a280200360234200841c0006a41086a200129030037030020082008290330370340200020024101756a210102402002410171450d00200128020020076a28020021070b200841d0006a41086a200841c0006a41086a2903002206370300200841086a200637030020082008290340220637035020082006370300200120042005200820071101004100200841e0006a36020441010bb30203017f037e057f410028020441206b2208210a410020083602042001280204210220012802002109024002400240024010012201450d002001418104490d012001104221080c020b410021080c020b410020082001410f6a4170716b22083602040b20082001100e1a0b200a4200370310200a4200370308200a4200370318200141074b4120100b200a41086a20084108100c1a200141787122064108474120100b200a41086a41086a2207200841086a4108100c1a20064110474120100b200a41086a41106a2206200841106a4108100c1a02402001418104490d00200810450b200020024101756a21012006290300210520072903002104200a290308210302402002410171450d00200128020020096a28020021090b200120032004200520091102004100200a41206a36020441010bd50101027f200028020021022001280200220328020820032802046b41074b4120100b200220032802044108100c1a2003200328020441086a360204200028020021002001280200220328020820032802046b41074b4120100b200041086a20032802044108100c1a2003200328020441086a3602042001280200220328020820032802046b41074b4120100b200041106a20032802044108100c1a2003200328020441086a2201360204200328020820016b41074b4120100b200041186a20032802044108100c1a2003200328020441086a3602040b3801027f02402000410120001b2201104222000d0003404100210041002802c0072202450d012002110300200110422200450d000b0b20000b0e0002402000450d00200010450b0b05001000000b4901037f4100210502402002450d000240034020002d0000220320012d00002204470d01200141016a2101200041016a21002002417f6a22020d000c020b0b200320046b21050b20050b090041c407200010430bcd04010c7f02402001450d00024020002802c041220d0d004110210d200041c0c1006a41103602000b200141086a200141046a41077122026b200120021b210202400240024020002802c441220a200d4f0d002000200a410c6c6a4180c0006a21010240200a0d0020004184c0006a220d2802000d0020014180c000360200200d20003602000b200241046a210a034002402001280208220d200a6a20012802004b0d002001280204200d6a220d200d28020041808080807871200272360200200141086a22012001280200200a6a360200200d200d28020041808080807872360200200d41046a22010d030b2000104422010d000b0b41fcffffff0720026b2104200041c8c1006a210b200041c0c1006a210c20002802c8412203210d03402000200d410c6c6a22014188c0006a28020020014180c0006a2205280200464190c900100b20014184c0006a280200220641046a210d0340200620052802006a2107200d417c6a2208280200220941ffffffff07712101024020094100480d000240200120024f0d000340200d20016a220a20074f0d01200a280200220a4100480d012001200a41ffffffff07716a41046a22012002490d000b0b20082001200220012002491b200941808080807871723602000240200120024d0d00200d20026a200420016a41ffffffff07713602000b200120024f0d040b200d20016a41046a220d2007490d000b41002101200b4100200b28020041016a220d200d200c280200461b220d360200200d2003470d000b0b20010f0b2008200828020041808080807872360200200d0f0b41000b870501087f20002802c44121010240024041002d00e649450d0041002802e84921070c010b3f002107410041013a00e6494100200741107422073602e8490b200721030240024002400240200741ffff036a41107622023f0022084d0d00200220086b40001a4100210820023f00470d0141002802e84921030b41002108410020033602e84920074100480d0020002001410c6c6a210220074180800441808008200741ffff037122084181f8034922061b6a2008200741ffff077120061b6b20076b2107024041002d00e6490d003f002103410041013a00e6494100200341107422033602e8490b20024180c0006a210220074100480d01200321060240200741076a417871220520036a41ffff036a41107622083f0022044d0d00200820046b40001a20083f00470d0241002802e84921060b4100200620056a3602e8492003417f460d0120002001410c6c6a22014184c0006a2802002206200228020022086a2003460d020240200820014188c0006a22052802002201460d00200620016a2206200628020041808080807871417c20016b20086a72360200200520022802003602002006200628020041ffffffff07713602000b200041c4c1006a2202200228020041016a220236020020002002410c6c6a22004184c0006a200336020020004180c0006a220820073602000b20080f0b02402002280200220820002001410c6c6a22034188c0006a22012802002207460d0020034184c0006a28020020076a2203200328020041808080807871417c20076b20086a72360200200120022802003602002003200328020041ffffffff07713602000b2000200041c4c1006a220728020041016a22033602c0412007200336020041000f0b2002200820076a36020020020b7b01037f024002402000450d0041002802844922024101480d0041c4c70021032002410c6c41c4c7006a21010340200341046a2802002202450d010240200241046a20004b0d00200220032802006a20004b0d030b2003410c6a22032001490d000b0b0f0b2000417c6a2203200328020041ffffffff07713602000b0300000b0bce07180041040b04f04c00000041100b067772697465000041200b0572656164000041300b147472616e73616374696f6e2065787069726564000041d0000b336f626a6563742070617373656420746f206974657261746f725f746f206973206e6f7420696e206d756c74695f696e64657800004190010b2370726f706f73616c2077697468207468652073616d65206e616d6520657869737473000041c0010b217472616e73616374696f6e20617574686f72697a6174696f6e206661696c6564000041f0010b3363616e6e6f7420637265617465206f626a6563747320696e207461626c65206f6620616e6f7468657220636f6e7472616374000041b0020b176572726f722072656164696e67206974657261746f72000041d0020b04676574000041e0020b1370726f706f73616c206e6f7420666f756e6400004180030b33617070726f76616c206973206e6f74206f6e20746865206c697374206f662072657175657374656420617070726f76616c73000041c0030b2e6f626a6563742070617373656420746f206d6f64696679206973206e6f7420696e206d756c74695f696e646578000041f0030b3363616e6e6f74206d6f64696679206f626a6563747320696e207461626c65206f6620616e6f7468657220636f6e7472616374000041b0040b3b757064617465722063616e6e6f74206368616e6765207072696d617279206b6579207768656e206d6f64696679696e6720616e206f626a656374000041f0040b1f6e6f20617070726f76616c2070726576696f75736c79206772616e74656400004190050b1f63616e6e6f742063616e63656c20756e74696c2065787069726174696f6e000041b0050b2d6f626a6563742070617373656420746f206572617365206973206e6f7420696e206d756c74695f696e646578000041e0050b3263616e6e6f74206572617365206f626a6563747320696e207461626c65206f6620616e6f7468657220636f6e7472616374000041a0060b35617474656d707420746f2072656d6f7665206f626a656374207468617420776173206e6f7420696e206d756c74695f696e646578000041e0060b086f6e6572726f72000041f0060b06656f73696f00004180070b406f6e6572726f7220616374696f6e277320617265206f6e6c792076616c69642066726f6d207468652022656f73696f222073797374656d206163636f756e7400004190c9000b566d616c6c6f635f66726f6d5f6672656564207761732064657369676e656420746f206f6e6c792062652063616c6c6564206166746572205f686561702077617320636f6d706c6574656c7920616c6c6f636174656400"},"hex_data":"0000735802ea30550000e8a9010061736d010000000198011760017f0060047f7e7e7f0060047f7e7e7e006000006000017e60027e7e0060067f7f7f7f7f7f017f60027f7f0060037f7f7f017f60077e7e7f7f7f7f7e017f6000017f60027f7f017f60017e0060047e7e7e7e017f60067e7e7e7e7f7f017f60047f7e7f7f0060057f7e7f7f7f0060057e7e7f7f7e017f60047f7f7f7f0060037f7e7f017f60047f7f7e7f0060037e7e7e0060017f017f02f7021203656e760561626f7274000303656e7610616374696f6e5f646174615f73697a65000a03656e761e636865636b5f7065726d697373696f6e5f617574686f72697a6174696f6e000903656e761f636865636b5f7472616e73616374696f6e5f617574686f72697a6174696f6e000603656e761063757272656e745f7265636569766572000403656e760c63757272656e745f74696d65000403656e760b64625f66696e645f693634000d03656e760a64625f6765745f693634000803656e760d64625f72656d6f76655f693634000003656e760c64625f73746f72655f693634000e03656e760d64625f7570646174655f693634000f03656e760c656f73696f5f617373657274000703656e76066d656d637079000803656e76076d656d6d6f7665000803656e7610726561645f616374696f6e5f64617461000b03656e760c726571756972655f61757468000c03656e760d726571756972655f6175746832000503656e760d73656e645f646566657272656400100336350b0b0b0a000807070b0b0b0b0b0b0b0b11000b0b0b071207120b07011314070b01140213070702150b0b0b0716000008160b16000304050170010606050301000107f10511066d656d6f72790200165f5a6571524b3131636865636b73756d32353653315f0012165f5a6571524b3131636865636b73756d31363053315f0013165f5a6e65524b3131636865636b73756d31363053315f0014036e6f770015305f5a4e35656f73696f3132726571756972655f6175746845524b4e535f31367065726d697373696f6e5f6c6576656c450016b0015f5a4e35656f73696f3331636865636b5f7472616e73616374696f6e5f617574686f72697a6174696f6e45524b4e535f31317472616e73616374696f6e45524b4e5374335f5f3133736574494e535f31367065726d697373696f6e5f6c6576656c454e53335f346c6573734953355f45454e53335f39616c6c6f6361746f724953355f45454545524b4e53345f4931307075626c69635f6b65794e53365f4953445f45454e53385f4953445f4545454500179f015f5a4e35656f73696f3330636865636b5f7065726d697373696f6e5f617574686f72697a6174696f6e457979524b4e5374335f5f31337365744931307075626c69635f6b65794e53305f346c6573734953325f45454e53305f39616c6c6f6361746f724953325f45454545524b4e53315f494e535f31367065726d697373696f6e5f6c6576656c454e53335f4953415f45454e53355f4953415f454545457900221c5f5a4e35656f73696f386d756c74697369673770726f706f7365457600233b5f5a4e35656f73696f386d756c746973696737617070726f766545794e535f346e616d65454e535f31367065726d697373696f6e5f6c6576656c45002d3d5f5a4e35656f73696f386d756c746973696739756e617070726f766545794e535f346e616d65454e535f31367065726d697373696f6e5f6c6576656c450032255f5a4e35656f73696f386d756c74697369673663616e63656c45794e535f346e616d6545790034235f5a4e35656f73696f386d756c7469736967346578656345794e535f346e616d6545790038056170706c790039066d656d636d700041066d616c6c6f63004204667265650045090c010041000b064638322d34230ab09701350b002000200141201041450b0b002000200141201041450b0d0020002001412010414100470b0a00100542c0843d80a70b0e002000290300200029030810100bba0503027f017e077f4100410028020441306b220c360204200c41106a200010184100210b410021094100210a024020022802082203450d0041002106200c4100360208200c42003703002003ad21050340200641016a2106200542078822054200520d000b02400240024020022802002207200241046a2204460d00034002400240200722082802042200450d0003402000220728020022000d000c020b0b200828020822072802002008460d00200841086a210803402008280200220041086a2108200020002802082207280200470d000b0b200641226a210620072004470d000b2006450d010b200c20061019200c2802042107200c28020021000c010b41002107410021000b200c2000360224200c2000360220200c2007360228200c41206a2002101a1a200c2802042109200c280200210a0b41002100024020012802082202450d0041002106200c4100360208200c42003703002002ad21050340200641016a2106200542078822054200520d000b02400240024020012802002207200141046a2204460d00034002400240200722082802042200450d0003402000220728020022000d000c020b0b200828020822072802002008460d00200841086a210803402008280200220041086a2108200020002802082207280200470d000b0b200641106a210620072004470d000b2006450d010b200c20061019200c2802042107200c28020021000c010b41002107410021000b200c2000360224200c2000360220200c2007360228200c41206a2001101b1a200c280204210b200c28020021000b200c2802102207200c28021420076b200a410020031b2009200a6b410020031b2000410020021b200b20006b410020021b1003210702402000450d002000103f0b0240200a450d00200a103f0b0240200c2802102200450d00200c20003602142000103f0b4100200c41306a360204200741004a0b9e0101037f4100410028020441106b22043602044100210320004100360208200042003702002004410036020020042001101c1a0240024020042802002202450d00200020021019200041046a2802002103200028020021000c010b410021000b20042000360204200420003602002004200336020820042001101d1a2004200141186a101e200141246a101e200141306a101f1a4100200441106a3602040bad0201057f0240024002400240024020002802082202200028020422066b20014f0d002006200028020022056b220320016a2204417f4c0d0241ffffffff0721060240200220056b220241feffffff034b0d0020042002410174220620062004491b2206450d020b2006103e21020c030b200041046a21000340200641003a00002000200028020041016a22063602002001417f6a22010d000c040b0b41002106410021020c010b20001040000b200220066a2104200220036a220521060340200641003a0000200641016a21062001417f6a22010d000b2005200041046a2203280200200028020022016b22026b2105024020024101480d00200520012002100c1a200028020021010b2000200536020020032006360200200041086a20043602002001450d002001103f0f0b0be20203027f017e057f4100410028020441d0006b22093602042000280204210520013502082104200041086a2103200041046a210703402004a721062009200442078822044200522208410774200641ff0071723a0028200328020020056b41004a4110100b2007280200200941286a4101100c1a2007200728020041016a220536020020080d000b024020012802002206200141046a2201460d00200041086a2102200041046a21030340200941066a20062208410d6a4122100c1a200941286a200941066a4122100c1a200228020020056b41214a4110100b2003280200200941286a4122100c1a2003200328020041226a22053602000240024020082802042207450d0003402007220628020022070d000c020b0b200828020822062802002008460d00200841086a210803402008280200220741086a2108200720072802082206280200470d000b0b20062001470d000b0b4100200941d0006a36020420000bed0203017f017e057f4100410028020441106b22083602042000280204210420013502082103200041086a2102200041046a210603402003a721052008200342078822034200522207410774200541ff0071723a000f200228020020046b41004a4110100b20062802002008410f6a4101100c1a2006200628020041016a220436020020070d000b024020012802002205200141046a2201460d00200041046a21020340200041086a220628020020046b41074a4110100b20022802002005220741106a4108100c1a2002200228020041086a2205360200200628020020056b41074a4110100b2002280200200741186a4108100c1a2002200228020041086a22043602000240024020072802042206450d0003402006220528020022060d000c020b0b200728020822052802002007460d00200741086a210703402007280200220641086a2107200620062802082205280200470d000b0b20052001470d000b0b4100200841106a36020420000b9d0502067f017e200020002802002206410a6a3602002006410b6a2106200135020c21080340200641016a2106200842078822084200520d000b20002006360200200135021421080340200641016a2106200842078822084200520d000b200020063602002001411c6a2802002202200128021822076b41286dad21080340200641016a2106200842078822084200520d000b20002006360200024020072002460d000340200641106a2106200741146a2802002203200728021022046b2205410475ad21080340200641016a2106200842078822084200520d000b024020042003460d00200541707120066a21060b2006200741206a28020022036a200728021c22046b2106200320046bad21080340200641016a2106200842078822084200520d000b200741286a22072002470d000b200020063602000b200141286a2802002202200128022422076b41286dad21080340200641016a2106200842078822084200520d000b20002006360200024020072002460d000340200641106a2106200741146a2802002203200728021022046b2205410475ad21080340200641016a2106200842078822084200520d000b024020042003460d00200541707120066a21060b2006200741206a28020022036a200728021c22046b2106200320046bad21080340200641016a2106200842078822084200520d000b200741286a22072002470d000b200020063602000b200141346a2802002205200128023022076b410475ad21080340200641016a2106200842078822084200520d000b20002006360200024020072005460d0003402006200741086a28020022036a41026a200728020422046b2106200320046bad21080340200641016a2106200842078822084200520d000b200741106a22072005470d000b200020063602000b20000b9f0303047f017e017f4100410028020441106b2207360204200028020820002802046b41034a4110100b200028020420014104100c1a2000200028020441046a2204360204200028020820046b41014a4110100b2000280204200141046a4102100c1a2000200028020441026a2204360204200028020820046b41034a4110100b2000280204200141086a4104100c1a2000200028020441046a2205360204200135020c210603402006a721042007200642078822064200522202410774200441ff0071723a000e200041086a28020020056b41004a4110100b200041046a22042802002007410e6a4101100c1a2004200428020041016a220536020020020d000b200041086a220328020020056b41004a4110100b200041046a2204280200200141106a4101100c1a2004200428020041016a22053602002001350214210603402006a721022007200642078822064200522201410774200241ff0071723a000f200328020020056b41004a4110100b20042802002007410f6a4101100c1a2004200428020041016a220536020020010d000b4100200741106a36020420000bbe0203037f017e037f4100410028020441106b2208360204200128020420012802006b41286dad210520002802042106200041086a2103200041046a210403402005a721072008200542078822054200522202410774200741ff0071723a000f200328020020066b41004a4110100b20042802002008410f6a4101100c1a2004200428020041016a220636020020020d000b024020012802002207200141046a2802002203460d00200041046a21040340200041086a220228020020066b41074a4110100b200428020020074108100c1a2004200428020041086a2206360200200228020020066b41074a4110100b2004280200200741086a4108100c1a2004200428020041086a3602002000200741106a10212007411c6a10201a200741286a22072003460d01200428020021060c000b0b4100200841106a36020420000b880203037f017e027f4100410028020441106b2207360204200128020420012802006b410475ad210520002802042106200041086a210303402005a721042007200542078822054200522202410774200441ff0071723a000f200328020020066b41004a4110100b200041046a22042802002007410f6a4101100c1a2004200428020041016a220636020020020d000b024020012802002204200141046a2802002202460d00200041086a21030340200328020020066b41014a4110100b200041046a220628020020044102100c1a2006200628020041026a3602002000200441046a10201a200441106a22042002460d01200628020021060c000b0b4100200741106a36020420000bda0103057f017e017f4100410028020441106b2208360204200128020420012802006bad210720002802042106200041086a2104200041046a210503402007a721022008200742078822074200522203410774200241ff0071723a000f200428020020066b41004a4110100b20052802002008410f6a4101100c1a2005200528020041016a220636020020030d000b200041086a28020020066b200141046a280200200128020022026b22054e4110100b200041046a220628020020022005100c1a2006200628020020056a3602004100200841106a36020420000ba40203027f017e037f4100410028020441106b2207360204200128020420012802006b410475ad210420002802042105200041086a210203402004a721032007200442078822044200522206410774200341ff0071723a000f200228020020056b41004a4110100b200041046a22032802002007410f6a4101100c1a2003200328020041016a220536020020060d000b024020012802002206200141046a2802002201460d00200041046a21030340200041086a220228020020056b41074a4110100b200328020020064108100c1a2003200328020041086a2205360200200228020020056b41074a4110100b2003280200200641086a4108100c1a2003200328020041086a2205360200200641106a22062001470d000b0b4100200741106a36020420000b900503027f017e087f4100410028020441206b220f3602044100210e4100210c4100210d024020022802082205450d0041002108200f4100360208200f42003703002005ad21070340200841016a2108200742078822074200520d000b02400240024020022802002209200241046a2206460d000340024002402009220b280204220a450d000340200a2209280200220a0d000c020b0b200b2802082209280200200b460d00200b41086a210b0340200b280200220a41086a210b200a200a2802082209280200470d000b0b200841226a210820092006470d000b2008450d010b200f20081019200f2802042109200f280200210a0c010b410021094100210a0b200f200a360214200f200a360210200f2009360218200f41106a2002101a1a200f280204210c200f280200210d0b4100210a024020032802082202450d0041002108200f4100360208200f42003703002002ad21070340200841016a2108200742078822074200520d000b02400240024020032802002209200341046a2206460d000340024002402009220b280204220a450d000340200a2209280200220a0d000c020b0b200b2802082209280200200b460d00200b41086a210b0340200b280200220a41086a210b200a200a2802082209280200470d000b0b200841106a210820092006470d000b2008450d010b200f20081019200f2802042109200f280200210a0c010b410021094100210a0b200f200a360214200f200a360210200f2009360218200f41106a2003101b1a200f280204210e200f280200210a0b20002001200d410020051b200c200d6b410020051b200a410020021b200e200a6b410020021b2004100221090240200a450d00200a103f0b0240200d450d00200d103f0b4100200f41206a360204200941004a0b990e05027f017e027f017e017f4100410028020441f0016b220236020420022207100122053602c801024002402005418104490d002005104221020c010b410020022005410f6a4170716b22023602040b200720023602c40120022005100e1a200742003703b00141002105200741003602a801200742003703a001100521062007410036029401200741003a0098012007410036029c012007200642c0843d80a7413c6a36028801200720072802c401220236027c20072802c8012101200720023602782007200220016a36028001200141074b4120100b200741b8016a200728027c4108100c1a2007200728027c41086a220236027c20072802800120026b41074b4120100b200741b0016a200728027c4108100c1a2007200728027c41086a36027c200741f8006a200741a0016a10241a2007200728027c20072802786b360274200741f8006a20074188016a10251a20072903b801100f10052106200728028801200642c0843d80a74f4130100b20072903b80121062007200029030022033703482007427f3703582007410036026020072006370350200741e4006a4100360200200741e8006a41003602004100210202402003200642808080888dccd6f4ad7f20072903b001100622014100480d00200741c8006a200110262202280214200741c8006a4641d000100b0b200245419001100b200742003703382007410036024020072802a401220220072802a00122016b2204410475ad210603402005417f6a2105200642078822064200520d000b024002400240024020012002460d00200441707122022005470d0141002102410021050c030b410020056b21050c010b200220056b21050b200741386a20051019200728023c2102200728023821050b200720053602142007200536021020072002360218200741106a200741a0016a10211a20072802c401200728027422056a20072802c80120056b4100410020072802382205200728023c20056b100341004a41c001100b20072903b80121062007200741c4016a3602142007200741b0016a3602102007200741f4006a3602182007200741c8016a36021c200720063703e801200729034810045141f001100b2007200741106a3602d4012007200741c8006a3602d0012007200741e8016a3602d8014120103e2205420037030020054200370208200541003602102005200741c8006a360214200741d0016a20051027200720053602082007200529030022063703d0012007200528021822013602e00102400240200741e4006a22042802002202200741c8006a41206a2802004f0d00200220063703082002200136021020074100360208200220053602002004200241186a3602000c010b200741e0006a200741086a200741d0016a200741e0016a10280b200728020821052007410036020802402005450d00024020052802082202450d002005410c6a20023602002002103f0b2005103f0b2007427f3703202007410036022820072903b8012106200720002903002203370310200720063703182007412c6a22024100360200200741306a220141003602002007200741a0016a36020c2007200741b0016a360208200720063703e801200310045141f001100b2007200741086a3602d4012007200741106a3602d0012007200741e8016a3602d8014130103e220542003703002005420037020820054200370210200542003702182005200741106a360220200741d0016a20051029200720053602e0012007200529030022063703d0012007200528022422043602cc01024002402002280200220020012802004f0d002000200637030820002004360210200741003602e001200020053602002002200041186a3602000c010b200741286a200741e0016a200741d0016a200741cc016a102a0b20072802e0012105200741003602e00102402005450d00024020052802142200450d00200541186a20003602002000103f0b024020052802082200450d002005410c6a20003602002000103f0b2005103f0b024020072802282201450d00024002402007412c6a220428020022002001460d000340200041686a220028020021052000410036020002402005450d00024020052802142202450d00200541186a20023602002002103f0b024020052802082202450d002005410c6a20023602002002103f0b2005103f0b20012000470d000b200741286a28020021050c010b200121050b200420013602002005103f0b024020072802382205450d002007200536023c2005103f0b024020072802602201450d0002400240200741e4006a220428020022052001460d000340200541686a220528020021002005410036020002402000450d00024020002802082202450d002000410c6a20023602002002103f0b2000103f0b20012005470d000b200741e0006a28020021050c010b200121050b200420013602002005103f0b024020072802a0012205450d00200720053602a4012005103f0b4100200741f0016a3602040bd00203037f017e027f200028020421074100210642002105200041086a2102200041046a21030340200720022802004941d002100b200328020022072d000021042003200741016a2207360200200441ff0071200641ff0171220674ad2005842105200641076a210620044107760d000b0240024002402005a7220420012802042202200128020022076b41047522064d0d002001200420066b102c20012802002207200141046a2802002202470d010c020b0240200420064f0d00200141046a200720044104746a22023602000b20072002460d010b200041046a220428020021060340200041086a220328020020066b41074b4120100b200720042802004108100c1a2004200428020041086a2206360200200328020020066b41074b4120100b200741086a20042802004108100c1a2004200428020041086a2206360200200741106a22072002470d000b0b20000b810303037f017e027f200028020820002802046b41034b4120100b200120002802044104100c1a2000200028020441046a2202360204200028020820026b41014b4120100b200141046a20002802044102100c1a2000200028020441026a2202360204200028020820026b41034b4120100b200141086a20002802044104100c1a2000200028020441046a2204360204410021064200210503402004200041086a2802004941d002100b200041046a220728020022042d000021022007200441016a2204360200200241ff0071200641ff0171220674ad2005842105200641076a210620024107760d000b200120053e020c200041086a22032802002004474120100b200141106a200041046a22042802004101100c1a2004200428020041016a220636020041002107420021050340200620032802004941d002100b200428020022062d000021022004200641016a2206360200200241ff0071200741ff0171220774ad2005842105200741076a210720024107760d000b200120053e021420000b8b0403037f017e047f410028020441306b220921084100200936020402402000411c6a280200220720002802182202460d00410020026b2103200741686a21060340200641106a2802002001460d0120062107200641686a22042106200420036a4168470d000b0b0240024020072002460d00200741686a28020021060c010b20014100410010072206411f7641017341b002100b024002402006418104490d002006104221040c010b410020092006410f6a4170716b22043602040b20012004200610071a20082004360224200820043602202008200420066a2207360228024020064180044d0d0020041045200841286a2802002107200828022421040b4120103e22064200370300200642003702082006410036021020062000360214200720046b41074b4120100b200620044108100c1a2008200441086a360224200841206a200641086a102b1a200620013602182008200636021820082006290300220537031020082006280218220736020c024002402000411c6a22012802002204200041206a2802004f0d00200420053703082004200736021020084100360218200420063602002001200441186a3602000c010b200041186a200841186a200841106a2008410c6a10280b20082802182104200841003602182004450d00024020042802082207450d002004410c6a20073602002007103f0b2004103f0b4100200841306a36020420060be40403077f017e017f4100410028020441106b22053602042000280200210220012000280204220628020029030037030020062802042802002107410021082005220a4100360200200a410036020420062802082103200a41003602084100210402400240200628020c280200200328020022036b2206450d002006417f4c0d01200a41086a2006103e220420066a2208360200200a20043602002004200720036a2006100c1a200a20083602040b0240024020012802082206450d002001410c6a20063602002006103f200141106a22064100360200200141086a42003702000c010b200141106a21060b200620083602002001410c6a2008360200200141086a2004360200200841086a20046b2106200820046bad21090340200641016a2106200942078822094200520d000b024002402006418104490d002006104221070c010b410020052006410f6a4170716b22073602040b200641074a4110100b200720014108100c1a2001410c6a280200200141086a2802006bad2109200741086a2108200720066a210303402009a72104200a200942078822094200522205410774200441ff0071723a000f200320086b41004a4110100b2008200a410f6a4101100c1a200841016a210820050d000b200320086b2001410c6a280200200141086a28020022046b22054e4110100b200820042005100c1a2001200229030842808080888dccd6f4ad7f20002802082903002001290300220920072006100936021802402006418104490d00200710450b024020092002290310540d00200241106a427e200942017c2009427d561b3703000b4100200a41106a3602040f0b200a1040000bc50301047f024002402000280204200028020022066b41186d220441016a220541abd5aad5004f0d0041aad5aad500210702400240200028020820066b41186d220641d4aad52a4b0d0020052006410174220720072005491b2207450d010b200741186c103e21060c020b41002107410021060c010b20001040000b20012802002105200141003602002006200441186c6a2201200536020020012002290300370308200120032802003602102006200741186c6a2104200141186a210502400240200041046a280200220620002802002207460d000340200641686a2202280200210320024100360200200141686a2003360200200141786a200641786a280200360200200141746a200641746a280200360200200141706a200641706a280200360200200141686a21012002210620072002470d000b200041046a2802002107200028020021020c010b200721020b20002001360200200041046a2005360200200041086a2004360200024020072002460d000340200741686a220728020021012007410036020002402001450d00024020012802082206450d002001410c6a20063602002006103f0b2001103f0b20022007470d000b0b02402002450d002002103f0b0b880403057f017e047f410028020441106b220b210a4100200b36020420002802002102200120002802042208280200290300370300200141086a2103200828020421080240024020012802082205450d002001410c6a220920053602002005103f200141106a22054100360200200342003702000c010b200141106a21052001410c6a21090b200320082802003602002009200828020436020020052008280208360200200842003702002008410036020820092802002209200328020022046b2206410475ad2107410821080340200841016a2108200742078822074200520d000b200141146a2105024020042009460d00200641707120086a21080b200141186a2802002209200528020022046b2206410475ad21070340200841016a2108200742078822074200520d000b024020042009460d00200641707120086a21080b024002402008418104490d002008104221090c010b4100200b2008410f6a4170716b22093602040b200a2009360200200a200920086a360208200841074a4110100b200920014108100c1a200a200941086a360204200a200310211a200a200510211a20012002290308428080808ecdcddeb53520002802082903002001290300220720092008100936022402402008418104490d00200910450b024020072002290310540d00200241106a427e200742017c2007427d561b3703000b4100200a41106a3602040be00301047f024002402000280204200028020022066b41186d220441016a220541abd5aad5004f0d0041aad5aad500210702400240200028020820066b41186d220641d4aad52a4b0d0020052006410174220720072005491b2207450d010b200741186c103e21060c020b41002107410021060c010b20001040000b20012802002105200141003602002006200441186c6a2201200536020020012002290300370308200120032802003602102006200741186c6a2104200141186a210502400240200041046a280200220620002802002207460d000340200641686a2202280200210320024100360200200141686a2003360200200141786a200641786a280200360200200141746a200641746a280200360200200141706a200641706a280200360200200141686a21012002210620072002470d000b200041046a2802002107200028020021020c010b200721020b20002001360200200041046a2005360200200041086a2004360200024020072002460d000340200741686a220728020021012007410036020002402001450d00024020012802142206450d00200141186a20063602002006103f0b024020012802082206450d002001410c6a20063602002006103f0b2001103f0b20022007470d000b0b02402002450d002002103f0b0b820203047f017e017f200028020421054100210742002106200041086a2102200041046a21030340200520022802004941d002100b200328020022052d000021042003200541016a2205360200200441ff0071200741ff0171220774ad2006842106200741076a210720044107760d000b024002402006a7220320012802042207200128020022046b22024d0d002001200320026b1019200041046a2802002105200141046a2802002107200128020021040c010b200320024f0d00200141046a200420036a22073602000b200041086a28020020056b200720046b22054f4120100b2004200041046a22072802002005100c1a2007200728020020056a36020020000bab0201067f0240024002400240024020002802082202200028020422076b41047520014f0d002007200028020022066b410475220320016a22044180808080014f0d0241ffffffff0021050240200220066b220241047541feffff3f4b0d0020042002410375220520052004491b2205450d0220054180808080014f0d040b2005410474103e2102200041046a2802002107200028020021060c040b200041046a200720014104746a3602000f0b41002105410021020c020b20001040000b1000000b200220034104746a2203200720066b22076b2104200320014104746a2101200220054104746a2105024020074101480d00200420062007100c1a200028020021060b20002004360200200041046a2001360200200041086a200536020002402006450d002006103f0b0bfa0202027e047f4100410028020441c0006b220936020420032903002204200329030822051010200941386a4100360200200920013703202009427f37032820094200370330200920002903003703180240200941186a200241e002102e220628020822002006410c6a2802002208460d000340024020002903002004520d00200041086a2903002005510d020b2008200041106a2200470d000b200821000b200920003602102000200847418003100b200920033602082009200941106a36020c200941186a20062001200941086a102f024020092802302206450d0002400240200941346a220728020022082006460d000340200841686a220828020021002008410036020002402000450d00024020002802142203450d00200041186a20033602002003103f0b024020002802082203450d002000410c6a20033602002003103f0b2000103f0b20062008470d000b200941306a28020021000c010b200621000b200720063602002000103f0b4100200941c0006a3602040bb50101057f02402000411c6a280200220720002802182203460d00200741686a2106410020036b2104034020062802002903002001510d0120062107200641686a22052106200520046a4168470d000b0b0240024020072003460d00200741686a280200220628022020004641d000100b0c010b4100210620002903002000290308428080808ecdcddeb5352001100622054100480d00200020051031220628022020004641d000100b0b20064100472002100b20060ba00404017e047f017e037f410028020441106b220c210b4100200c360204200128022020004641c003100b200029030010045141f003100b2003280200210a2001290300210402400240200141186a220728020022052001411c6a280200460d002005200a290300370300200541086a200a41086a2903003703002007200728020041106a3602000c010b200141146a200a10300b02402001410c6a220a2802002003280204280200220341106a22056b22064104752207450d00200320052006100d1a0b200a200320074104746a360200200420012903005141b004100b41082103200141086a2105200a280200220a200128020822066b2208410475ad21090340200341016a2103200942078822094200520d000b200141146a210702402006200a460d00200841707120036a21030b200141186a280200220a200728020022066b2208410475ad21090340200341016a2103200942078822094200520d000b02402006200a460d00200841707120036a21030b024002402003418104490d0020031042210a0c010b4100200c2003410f6a4170716b220a3602040b200b200a360200200b200a20036a360208200341074a4110100b200a20014108100c1a200b200a41086a360204200b200510211a200b200710211a20012802242002200a2003100a02402003418104490d00200a10450b024020042000290310540d00200041106a427e200442017c2004427d561b3703000b4100200b41106a3602040b9f0201067f02400240024020002802042206200028020022056b410475220241016a22034180808080014f0d0041ffffffff00210402400240200028020820056b220741047541feffff3f4b0d0020032007410375220420042003491b2204450d0120044180808080014f0d030b2004410474103e2107200041046a2802002106200028020021050c030b41002104410021070c020b20001040000b1000000b200720024104746a22032001290300370300200341086a200141086a2903003703002003200620056b22016b2106200720044104746a2104200341106a2103024020014101480d00200620052001100c1a200028020021050b20002006360200200041046a2003360200200041086a200436020002402005450d002005103f0b0bba0403037f017e047f410028020441306b220921084100200936020402402000411c6a280200220720002802182202460d00410020026b2103200741686a21060340200641106a2802002001460d0120062107200641686a22042106200420036a4168470d000b0b0240024020072002460d00200741686a28020021060c010b20014100410010072206411f7641017341b002100b024002402006418104490d002006104221040c010b410020092006410f6a4170716b22043602040b20012004200610071a20082004360224200820043602202008200420066a2207360228024020064180044d0d0020041045200841286a2802002107200828022421040b4130103e2206420037030020064200370208200642003702102006420037021820062000360220200720046b41074b4120100b200620044108100c1a2008200441086a360224200841206a200641086a10241a200841206a200641146a10241a200620013602242008200636021820082006290300220537031020082006280224220736020c024002402000411c6a22012802002204200041206a2802004f0d00200420053703082004200736021020084100360218200420063602002001200441186a3602000c010b200041186a200841186a200841106a2008410c6a102a0b20082802182104200841003602182004450d00024020042802142207450d00200441186a20073602002007103f0b024020042802082207450d002004410c6a20073602002007103f0b2004103f0b4100200841306a36020420060bfa0202027e047f4100410028020441c0006b220936020420032903002204200329030822051010200941386a4100360200200920013703202009427f37032820094200370330200920002903003703180240200941186a200241e002102e22062802142200200641186a2802002208460d000340024020002903002004520d00200041086a2903002005510d020b2008200041106a2200470d000b200821000b20092000360210200020084741f004100b200920033602082009200941106a36020c200941186a20062001200941086a1033024020092802302206450d0002400240200941346a220728020022082006460d000340200841686a220828020021002008410036020002402000450d00024020002802142203450d00200041186a20033602002003103f0b024020002802082203450d002000410c6a20033602002003103f0b2000103f0b20062008470d000b200941306a28020021000c010b200621000b200720063602002000103f0b4100200941c0006a3602040ba30404017e047f017e037f410028020441106b220c210b4100200c360204200128022020004641c003100b200029030010045141f003100b2003280200210a20012903002104024002402001410c6a22072802002205200141106a280200460d002005200a290300370300200541086a200a41086a2903003703002007200728020041106a3602000c010b200141086a200a10300b0240200141186a220a2802002003280204280200220341106a22056b22064104752207450d00200320052006100d1a0b200a200320074104746a360200200420012903005141b004100b41082103200141086a21052001410c6a280200220a200128020822066b2208410475ad21090340200341016a2103200942078822094200520d000b200141146a210702402006200a460d00200841707120036a21030b200141186a280200220a200728020022066b2208410475ad21090340200341016a2103200942078822094200520d000b02402006200a460d00200841707120036a21030b024002402003418104490d0020031042210a0c010b4100200c2003410f6a4170716b220a3602040b200b200a360200200b200a20036a360208200341074a4110100b200a20014108100c1a200b200a41086a360204200b200510211a200b200710211a20012802242002200a2003100a02402003418104490d00200a10450b024020042000290310540d00200041106a427e200442017c2004427d561b3703000b4100200b41106a3602040bbe0401057f4100410028020441e0006b22083602042003100f200841286a41206a4100360200200820013703302008427f3703382008420037034020082000290300370328200841286a200241e00210352107024020032001510d002007410c6a280200210420072802082105100521032008410036020c200841003a0010200841003602142008200342c0843d80a7413c6a360200200820053602542008200536025020082004360258200841d0006a200810251a100521032008280200200342c0843d80a749419005100b0b200841206a4100360200200820013703082008427f37031020084200370318200820002903003703002008200241e002102e2100200841286a20071036200820001037024020082802182204450d00024002402008411c6a220628020022072004460d000340200741686a220728020021002007410036020002402000450d00024020002802142205450d00200041186a20053602002005103f0b024020002802082205450d002000410c6a20053602002005103f0b2000103f0b20042007470d000b200841186a28020021000c010b200421000b200620043602002000103f0b024020082802402204450d0002400240200841c4006a220628020022002004460d000340200041686a220028020021072000410036020002402007450d00024020072802082205450d002007410c6a20053602002005103f0b2007103f0b20042000470d000b200841c0006a28020021000c010b200421000b200620043602002000103f0b4100200841e0006a3602040bb60101057f02402000411c6a280200220720002802182203460d00200741686a2106410020036b2104034020062802002903002001510d0120062107200641686a22052106200520046a4168470d000b0b0240024020072003460d00200741686a280200220628021420004641d000100b0c010b410021062000290300200029030842808080888dccd6f4ad7f2001100622054100480d00200020051026220628021420004641d000100b0b20064100472002100b20060b800302017e067f200128021420004641b005100b200029030010045141e005100b02402000411c6a2205280200220720002802182203460d0020012903002102410020036b2106200741686a2108034020082802002903002002510d0120082107200841686a22042108200420066a4168470d000b0b200720034741a006100b200741686a210802400240200720052802002204460d00410020046b2103200821070340200741186a2208280200210620084100360200200728020021042007200636020002402004450d00024020042802082206450d002004410c6a20063602002006103f0b2004103f0b200741106a200741286a280200360200200741086a200741206a29030037030020082107200820036a4168470d000b2000411c6a28020022072008460d010b0340200741686a220728020021042007410036020002402004450d00024020042802082206450d002004410c6a20063602002006103f0b2004103f0b20082007470d000b0b2000411c6a2008360200200128021810080bb60302017e067f200128022020004641b005100b200029030010045141e005100b02402000411c6a2204280200220720002802182203460d0020012903002102410020036b2105200741686a2106034020062802002903002002510d0120062107200641686a22082106200820056a4168470d000b0b200720034741a006100b200741686a210802400240200720042802002206460d00410020066b2103200821060340200641186a2208280200210520084100360200200628020021072006200536020002402007450d00024020072802142205450d00200741186a20053602002005103f0b024020072802082205450d002007410c6a20053602002005103f0b2007103f0b200641106a200641286a280200360200200641086a200641206a29030037030020082106200820036a4168470d000b2000411c6a28020022072008460d010b0340200741686a220728020021062007410036020002402006450d00024020062802142205450d00200641186a20053602002005103f0b024020062802082205450d002006410c6a20053602002005103f0b2006103f0b20082007470d000b0b2000411c6a2008360200200128022410080bf00603057f017e027f4100410028020441a0016b220b3602042003100f41002108200b41f8006a41206a4100360200200b200137038001200b427f37038801200b420037039001200b2000290300370378200b41f8006a200241e0021035210a200b41d0006a41206a4100360200200b2001370358200b427f370360200b4200370368200b2000290300370350200b41d0006a200241e002102e210010052109200b4100360244200b41003a0048200b410036024c200b200942c0843d80a7413c6a360238200b200a2802082206360228200b200636022c200b200a410c6a280200360230200b41286a200b41386a10251a10052109200b280238200942c0843d80a74f4130100b200b4100360220200b4200370318200041186a2802002206200028021422056b2207410475ad2109200041146a210403402008417f6a2108200942078822094200520d000b024002400240024020052006460d00200741707122062008470d0141002106410021080c030b410020086b21080c010b200620086b21080b200b41186a20081019200b28021c2106200b28021821080b200b2008360204200b2008360200200b2006360208200b200410211a200a41086a22082802002206200a410c6a220528020020066b41004100200b2802182206200b28021c20066b100341004a41c001100b200b2001370308200b2002370300200b200320082802002208200528020020086b41001011200b41f8006a200a1036200b41d0006a200010370240200b2802182208450d00200b200836021c2008103f0b0240200b2802682206450d0002400240200b41ec006a2205280200220a2006460d000340200a41686a220a2802002108200a410036020002402008450d00024020082802142200450d00200841186a20003602002000103f0b024020082802082200450d002008410c6a20003602002000103f0b2008103f0b2006200a470d000b200b41e8006a28020021080c010b200621080b200520063602002008103f0b0240200b280290012206450d0002400240200b4194016a220528020022082006460d000340200841686a2208280200210a200841003602000240200a450d000240200a2802082200450d00200a410c6a20003602002000103f0b200a103f0b20062008470d000b200b4190016a28020021080c010b200621080b200520063602002008103f0b4100200b41a0016a3602040bfc0603027f047e017f4100410028020441e0006b220936020442002106423b210541e00621044200210703400240024002400240024020064206560d0020042c00002203419f7f6a41ff017141194b0d01200341a5016a21030c020b420021082006420b580d020c030b200341d0016a41002003414f6a41ff01714105491b21030b2003ad42388642388721080b2008421f83200542ffffffff0f838621080b200441016a2104200642017c2106200820078421072005427b7c2205427a520d000b024020072002520d0042002106423b210541f00621044200210703400240024002400240024020064204560d0020042c00002203419f7f6a41ff017141194b0d01200341a5016a21030c020b420021082006420b580d020c030b200341d0016a41002003414f6a41ff01714105491b21030b2003ad42388642388721080b2008421f83200542ffffffff0f838621080b200441016a2104200642017c2106200820078421072005427b7c2205427a520d000b2007200151418007100b0b0240024020012000510d0042002106423b210541e00621044200210703400240024002400240024020064206560d0020042c00002203419f7f6a41ff017141194b0d01200341a5016a21030c020b420021082006420b580d020c030b200341d0016a41002003414f6a41ff01714105491b21030b2003ad42388642388721080b2008421f83200542ffffffff0f838621080b200441016a2104200642017c2106200820078421072005427b7c2205427a520d000b20072002520d010b200920003703580240024002400240200242ffffffffd3cddeb535570d0020024280808080d4cddeb535510d0120024280808080c0a8a1d3c100510d02200242808080808080a0aad700520d04200941003602342009410136023020092009290330370228200941d8006a200941286a103c1a0c040b2002428080808094ccd6f4ad7f510d022002428080c0dae9dbd6e654520d03200941003602442009410236024020092009290340370218200941d8006a200941186a103b1a0c030b2009410036024c2009410336024820092009290348370210200941d8006a200941106a103b1a0c020b2009410036023c2009410436023820092009290338370220200941d8006a200941206a103c1a0c010b200941003602542009410536025020092009290350370208200941d8006a200941086a103a1a0b4100200941e0006a3602040b8c0101047f4100280204220521042001280204210220012802002101024010012203450d00024020034180044d0d002003104222052003100e1a200510450c010b410020052003410f6a4170716b220536020420052003100e1a0b200020024101756a210302402002410171450d00200328020020016a28020021010b200320011100004100200436020441010ba10303027f037e037f410028020441e0006b22092108410020093602042001280204210220012802002107024002400240024010012203450d002003418104490d012003104221010c020b410021010c020b410020092003410f6a4170716b22013602040b20012003100e1a0b200842003703182008420037031020082001360254200820013602502008200120036a3602582008200841d0006a3602302008200841106a360240200841c0006a200841306a103d02402003418104490d00200110450b200841106a41086a29030021052008413c6a2008412c6a280200360200200841306a41086a2201200841286a28020036020020082903102104200820082802203602302008200841246a280200360234200841c0006a41086a200129030037030020082008290330370340200020024101756a210102402002410171450d00200128020020076a28020021070b200841d0006a41086a200841c0006a41086a2903002206370300200841086a200637030020082008290340220637035020082006370300200120042005200820071101004100200841e0006a36020441010bb30203017f037e057f410028020441206b2208210a410020083602042001280204210220012802002109024002400240024010012201450d002001418104490d012001104221080c020b410021080c020b410020082001410f6a4170716b22083602040b20082001100e1a0b200a4200370310200a4200370308200a4200370318200141074b4120100b200a41086a20084108100c1a200141787122064108474120100b200a41086a41086a2207200841086a4108100c1a20064110474120100b200a41086a41106a2206200841106a4108100c1a02402001418104490d00200810450b200020024101756a21012006290300210520072903002104200a290308210302402002410171450d00200128020020096a28020021090b200120032004200520091102004100200a41206a36020441010bd50101027f200028020021022001280200220328020820032802046b41074b4120100b200220032802044108100c1a2003200328020441086a360204200028020021002001280200220328020820032802046b41074b4120100b200041086a20032802044108100c1a2003200328020441086a3602042001280200220328020820032802046b41074b4120100b200041106a20032802044108100c1a2003200328020441086a2201360204200328020820016b41074b4120100b200041186a20032802044108100c1a2003200328020441086a3602040b3801027f02402000410120001b2201104222000d0003404100210041002802c0072202450d012002110300200110422200450d000b0b20000b0e0002402000450d00200010450b0b05001000000b4901037f4100210502402002450d000240034020002d0000220320012d00002204470d01200141016a2101200041016a21002002417f6a22020d000c020b0b200320046b21050b20050b090041c407200010430bcd04010c7f02402001450d00024020002802c041220d0d004110210d200041c0c1006a41103602000b200141086a200141046a41077122026b200120021b210202400240024020002802c441220a200d4f0d002000200a410c6c6a4180c0006a21010240200a0d0020004184c0006a220d2802000d0020014180c000360200200d20003602000b200241046a210a034002402001280208220d200a6a20012802004b0d002001280204200d6a220d200d28020041808080807871200272360200200141086a22012001280200200a6a360200200d200d28020041808080807872360200200d41046a22010d030b2000104422010d000b0b41fcffffff0720026b2104200041c8c1006a210b200041c0c1006a210c20002802c8412203210d03402000200d410c6c6a22014188c0006a28020020014180c0006a2205280200464190c900100b20014184c0006a280200220641046a210d0340200620052802006a2107200d417c6a2208280200220941ffffffff07712101024020094100480d000240200120024f0d000340200d20016a220a20074f0d01200a280200220a4100480d012001200a41ffffffff07716a41046a22012002490d000b0b20082001200220012002491b200941808080807871723602000240200120024d0d00200d20026a200420016a41ffffffff07713602000b200120024f0d040b200d20016a41046a220d2007490d000b41002101200b4100200b28020041016a220d200d200c280200461b220d360200200d2003470d000b0b20010f0b2008200828020041808080807872360200200d0f0b41000b870501087f20002802c44121010240024041002d00e649450d0041002802e84921070c010b3f002107410041013a00e6494100200741107422073602e8490b200721030240024002400240200741ffff036a41107622023f0022084d0d00200220086b40001a4100210820023f00470d0141002802e84921030b41002108410020033602e84920074100480d0020002001410c6c6a210220074180800441808008200741ffff037122084181f8034922061b6a2008200741ffff077120061b6b20076b2107024041002d00e6490d003f002103410041013a00e6494100200341107422033602e8490b20024180c0006a210220074100480d01200321060240200741076a417871220520036a41ffff036a41107622083f0022044d0d00200820046b40001a20083f00470d0241002802e84921060b4100200620056a3602e8492003417f460d0120002001410c6c6a22014184c0006a2802002206200228020022086a2003460d020240200820014188c0006a22052802002201460d00200620016a2206200628020041808080807871417c20016b20086a72360200200520022802003602002006200628020041ffffffff07713602000b200041c4c1006a2202200228020041016a220236020020002002410c6c6a22004184c0006a200336020020004180c0006a220820073602000b20080f0b02402002280200220820002001410c6c6a22034188c0006a22012802002207460d0020034184c0006a28020020076a2203200328020041808080807871417c20076b20086a72360200200120022802003602002003200328020041ffffffff07713602000b2000200041c4c1006a220728020041016a22033602c0412007200336020041000f0b2002200820076a36020020020b7b01037f024002402000450d0041002802844922024101480d0041c4c70021032002410c6c41c4c7006a21010340200341046a2802002202450d010240200241046a20004b0d00200220032802006a20004b0d030b2003410c6a22032001490d000b0b0f0b2000417c6a2203200328020041ffffffff07713602000b0300000b0bce07180041040b04f04c00000041100b067772697465000041200b0572656164000041300b147472616e73616374696f6e2065787069726564000041d0000b336f626a6563742070617373656420746f206974657261746f725f746f206973206e6f7420696e206d756c74695f696e64657800004190010b2370726f706f73616c2077697468207468652073616d65206e616d6520657869737473000041c0010b217472616e73616374696f6e20617574686f72697a6174696f6e206661696c6564000041f0010b3363616e6e6f7420637265617465206f626a6563747320696e207461626c65206f6620616e6f7468657220636f6e7472616374000041b0020b176572726f722072656164696e67206974657261746f72000041d0020b04676574000041e0020b1370726f706f73616c206e6f7420666f756e6400004180030b33617070726f76616c206973206e6f74206f6e20746865206c697374206f662072657175657374656420617070726f76616c73000041c0030b2e6f626a6563742070617373656420746f206d6f64696679206973206e6f7420696e206d756c74695f696e646578000041f0030b3363616e6e6f74206d6f64696679206f626a6563747320696e207461626c65206f6620616e6f7468657220636f6e7472616374000041b0040b3b757064617465722063616e6e6f74206368616e6765207072696d617279206b6579207768656e206d6f64696679696e6720616e206f626a656374000041f0040b1f6e6f20617070726f76616c2070726576696f75736c79206772616e74656400004190050b1f63616e6e6f742063616e63656c20756e74696c2065787069726174696f6e000041b0050b2d6f626a6563742070617373656420746f206572617365206973206e6f7420696e206d756c74695f696e646578000041e0050b3263616e6e6f74206572617365206f626a6563747320696e207461626c65206f6620616e6f7468657220636f6e7472616374000041a0060b35617474656d707420746f2072656d6f7665206f626a656374207468617420776173206e6f7420696e206d756c74695f696e646578000041e0060b086f6e6572726f72000041f0060b06656f73696f00004180070b406f6e6572726f7220616374696f6e277320617265206f6e6c792076616c69642066726f6d207468652022656f73696f222073797374656d206163636f756e7400004190c9000b566d616c6c6f635f66726f6d5f6672656564207761732064657369676e656420746f206f6e6c792062652063616c6c6564206166746572205f686561702077617320636f6d706c6574656c7920616c6c6f636174656400"},{"account":"eosio","name":"setabi","authorization":[{"actor":"eosio.msig","permission":"active"}],"data":{"account":"eosio.msig","abi":"0e656f73696f3a3a6162692f312e30030c6163636f756e745f6e616d65046e616d650f7065726d697373696f6e5f6e616d65046e616d650b616374696f6e5f6e616d65046e616d650c107065726d697373696f6e5f6c6576656c0002056163746f720c6163636f756e745f6e616d650a7065726d697373696f6e0f7065726d697373696f6e5f6e616d6506616374696f6e0004076163636f756e740c6163636f756e745f6e616d65046e616d650b616374696f6e5f6e616d650d617574686f72697a6174696f6e127065726d697373696f6e5f6c6576656c5b5d0464617461056279746573127472616e73616374696f6e5f68656164657200060a65787069726174696f6e0e74696d655f706f696e745f7365630d7265665f626c6f636b5f6e756d0675696e743136107265665f626c6f636b5f7072656669780675696e743332136d61785f6e65745f75736167655f776f7264730976617275696e743332106d61785f6370755f75736167655f6d730575696e74380964656c61795f7365630976617275696e74333209657874656e73696f6e000204747970650675696e74313604646174610562797465730b7472616e73616374696f6e127472616e73616374696f6e5f6865616465720314636f6e746578745f667265655f616374696f6e7308616374696f6e5b5d07616374696f6e7308616374696f6e5b5d167472616e73616374696f6e5f657874656e73696f6e730b657874656e73696f6e5b5d0770726f706f736500040870726f706f7365720c6163636f756e745f6e616d650d70726f706f73616c5f6e616d65046e616d6509726571756573746564127065726d697373696f6e5f6c6576656c5b5d037472780b7472616e73616374696f6e07617070726f766500030870726f706f7365720c6163636f756e745f6e616d650d70726f706f73616c5f6e616d65046e616d65056c6576656c107065726d697373696f6e5f6c6576656c09756e617070726f766500030870726f706f7365720c6163636f756e745f6e616d650d70726f706f73616c5f6e616d65046e616d65056c6576656c107065726d697373696f6e5f6c6576656c0663616e63656c00030870726f706f7365720c6163636f756e745f6e616d650d70726f706f73616c5f6e616d65046e616d650863616e63656c65720c6163636f756e745f6e616d65046578656300030870726f706f7365720c6163636f756e745f6e616d650d70726f706f73616c5f6e616d65046e616d650865786563757465720c6163636f756e745f6e616d650870726f706f73616c00020d70726f706f73616c5f6e616d65046e616d65127061636b65645f7472616e73616374696f6e0562797465730e617070726f76616c735f696e666f00030d70726f706f73616c5f6e616d65046e616d65137265717565737465645f617070726f76616c73127065726d697373696f6e5f6c6576656c5b5d1270726f76696465645f617070726f76616c73127065726d697373696f6e5f6c6576656c5b5d0500000040615ae9ad0770726f706f736500000000406d7a6b3507617070726f7665000000509bde5acdd409756e617070726f766500000000004485a6410663616e63656c00000000000080545704657865630002000000d1605ae9ad03693634010d70726f706f73616c5f6e616d6501046e616d650870726f706f73616c0000c0d16c7a6b3503693634010d70726f706f73616c5f6e616d6501046e616d650e617070726f76616c735f696e666f000000"},"hex_data":"0000735802ea3055fd090e656f73696f3a3a6162692f312e30030c6163636f756e745f6e616d65046e616d650f7065726d697373696f6e5f6e616d65046e616d650b616374696f6e5f6e616d65046e616d650c107065726d697373696f6e5f6c6576656c0002056163746f720c6163636f756e745f6e616d650a7065726d697373696f6e0f7065726d697373696f6e5f6e616d6506616374696f6e0004076163636f756e740c6163636f756e745f6e616d65046e616d650b616374696f6e5f6e616d650d617574686f72697a6174696f6e127065726d697373696f6e5f6c6576656c5b5d0464617461056279746573127472616e73616374696f6e5f68656164657200060a65787069726174696f6e0e74696d655f706f696e745f7365630d7265665f626c6f636b5f6e756d0675696e743136107265665f626c6f636b5f7072656669780675696e743332136d61785f6e65745f75736167655f776f7264730976617275696e743332106d61785f6370755f75736167655f6d730575696e74380964656c61795f7365630976617275696e74333209657874656e73696f6e000204747970650675696e74313604646174610562797465730b7472616e73616374696f6e127472616e73616374696f6e5f6865616465720314636f6e746578745f667265655f616374696f6e7308616374696f6e5b5d07616374696f6e7308616374696f6e5b5d167472616e73616374696f6e5f657874656e73696f6e730b657874656e73696f6e5b5d0770726f706f736500040870726f706f7365720c6163636f756e745f6e616d650d70726f706f73616c5f6e616d65046e616d6509726571756573746564127065726d697373696f6e5f6c6576656c5b5d037472780b7472616e73616374696f6e07617070726f766500030870726f706f7365720c6163636f756e745f6e616d650d70726f706f73616c5f6e616d65046e616d65056c6576656c107065726d697373696f6e5f6c6576656c09756e617070726f766500030870726f706f7365720c6163636f756e745f6e616d650d70726f706f73616c5f6e616d65046e616d65056c6576656c107065726d697373696f6e5f6c6576656c0663616e63656c00030870726f706f7365720c6163636f756e745f6e616d650d70726f706f73616c5f6e616d65046e616d650863616e63656c65720c6163636f756e745f6e616d65046578656300030870726f706f7365720c6163636f756e745f6e616d650d70726f706f73616c5f6e616d65046e616d650865786563757465720c6163636f756e745f6e616d650870726f706f73616c00020d70726f706f73616c5f6e616d65046e616d65127061636b65645f7472616e73616374696f6e0562797465730e617070726f76616c735f696e666f00030d70726f706f73616c5f6e616d65046e616d65137265717565737465645f617070726f76616c73127065726d697373696f6e5f6c6576656c5b5d1270726f76696465645f617070726f76616c73127065726d697373696f6e5f6c6576656c5b5d0500000040615ae9ad0770726f706f736500000000406d7a6b3507617070726f7665000000509bde5acdd409756e617070726f766500000000004485a6410663616e63656c00000000000080545704657865630002000000d1605ae9ad03693634010d70726f706f73616c5f6e616d6501046e616d650870726f706f73616c0000c0d16c7a6b3503693634010d70726f706f73616c5f6e616d6501046e616d650e617070726f76616c735f696e666f000000"}],"transaction_extensions":[]}}}],"block_extensions":[],"id":"000000b177f6f2e5a0d29837ede5dff616178aece8f4c9f891f376f186f1cb72","block_num":177,"ref_block_prefix":932762272}`
	jsonDataSignedBlock := []byte(jsonString)

	signedBlock := &SignedBlock{}

	err := json.Unmarshal(jsonDataSignedBlock, signedBlock)
	assert.NoError(t, err)

	tx := signedBlock.Transactions[0]
	assert.NoError(t, err)

	unpacked, err := tx.Transaction.Packed.Unpack()
	assert.NoError(t, err)

	assert.Equal(t, 2, len(unpacked.Actions))

	packed, err := unpacked.Pack(CompressionZlib)
	assert.NoError(t, err)

	EqualNoDiff(t,
		unpacked.packed.PackedTransaction,
		packed.PackedTransaction,
		fmt.Sprintf("\nActual: Unpack -> Pack (len=%d)\nExpected: Pack (len=%d)",
			len(unpacked.packed.PackedTransaction),
			len(packed.PackedTransaction)))
}

func EqualNoDiff(t *testing.T, expected interface{}, actual interface{}, message string) bool {
	if !assert.ObjectsAreEqual(expected, actual) {
		return assert.Fail(t, fmt.Sprintf("Not equal: %s\n", message))
	}

	return true
}

func TestBlob(t *testing.T) {
	b := Blob("RU9TIEdv")

	t.Run("String", func(tt *testing.T) {
		assert.Equal(tt, "RU9TIEdv", b.String())
	})

	t.Run("Data", func(tt *testing.T) {
		data, err := b.Data()
		require.Nil(tt, err)
		assert.Equal(tt, []byte("EOS Go"), data)
	})

	t.Run("malformed data", func(tt *testing.T) {
		b := Blob("not base64")
		data, err := b.Data()
		require.Equal(tt, "illegal base64 data at input byte 3", err.Error())
		assert.Empty(tt, data)
	})
}
