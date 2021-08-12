package eos

import (
	"strings"
)

// ExtendedStringToName acts similar to StringToName with the big differences
// that it will automtically try to infer from which format to convert to a name.
// Current rules are:
// - If the `s` contains a `,` character, assumes it's a `Symbol`
// - If the `s` contains only upper-case characters and length is <= 7, assumes it's a `SymbolCode`
// - Otherwise, forwards `s` to `StringToName` directly
func ExtendedStringToName(s string) (val uint64, err error) {
	if strings.Contains(s, ",") {
		symbol, err := StringToSymbol(s)
		if err != nil {
			return 0, err
		}

		return symbol.ToUint64()
	}

	if symbolCodeRegex.MatchString(s) {
		symbolCode, err := StringToSymbolCode(s)
		if err != nil {
			return 0, err
		}

		return uint64(symbolCode), nil
	}

	return StringToName(s)
}

func MustStringToName(s string) (val uint64) {
	var err error
	val, err = StringToName(s)
	if err != nil {
		panic(err)
	}

	return
}

func StringToName(s string) (val uint64, err error) {
	// ported from the eosio codebase, libraries/chain/include/eosio/chain/name.hpp
	var i uint32
	sLen := uint32(len(s))
	for ; i <= 12; i++ {
		var c uint64
		if i < sLen {
			c = uint64(charToSymbol(s[i]))
		}

		if i < 12 {
			c &= 0x1f
			c <<= 64 - 5*(i+1)
		} else {
			c &= 0x0f
		}

		val |= c
	}

	return
}

func charToSymbol(c byte) byte {
	if c >= 'a' && c <= 'z' {
		return c - 'a' + 6
	}
	if c >= '1' && c <= '5' {
		return c - '1' + 1
	}
	return 0
}

var base32Alphabet = []byte(".12345abcdefghijklmnopqrstuvwxyz")

var eosioNameUint64 = uint64(6138663577826885632)
var eosioTokenNameUint64 = uint64(6138663591592764928)
var cachedNames = map[uint64]string{
	6138663577826885632: "eosio",
	6138663591592764928: "eosio.token",
}

func NameToString(in uint64) string {
	// Some particularly used name are pre-cached, so we avoid the transformation altogether, and reduce memory usage
	if name, found := cachedNames[in]; found {
		return name
	}

	// ported from libraries/chain/name.cpp in eosio
	a := []byte{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.', '.', '.', '.'}

	tmp := in
	i := uint32(0)
	for ; i <= 12; i++ {
		bit := 0x1f
		if i == 0 {
			bit = 0x0f
		}
		c := base32Alphabet[tmp&uint64(bit)]
		a[12-i] = c

		shift := uint(5)
		if i == 0 {
			shift = 4
		}

		tmp >>= shift
	}

	// We had a call to `strings.TrimRight` before, but that was causing lots of
	// allocation and lost CPU cycles. We now have our own cutting method that
	// improves performance a lot.
	return trimRightDots(a)
}

func trimRightDots(bytes []byte) string {
	trimUpTo := -1
	for i := 12; i >= 0; i-- {
		if bytes[i] == '.' {
			trimUpTo = i
		} else {
			break
		}
	}

	if trimUpTo == -1 {
		return string(bytes)
	}

	return string(bytes[0:trimUpTo])
}
