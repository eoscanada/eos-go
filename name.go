package eos

import (
	"strings"
)

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

func NameToString(in uint64) string {
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

	return strings.TrimRight(string(a), ".")
}
