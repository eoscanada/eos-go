package ecc

import (
	"crypto/elliptic"
	"fmt"
	"math/big"

	"github.com/eoscanada/eos-go/btcsuite/btcd/btcec"
)

type innerWAPublicKey struct {
}

func newInnerWAPublicKey() innerPublicKey {
	return &innerWAPublicKey{}
}

func (p *innerWAPublicKey) key(content []byte) (*btcec.PublicKey, error) {
	// First byte can be either 0x02 or 0x03
	if content[0] != 0x02 && content[0] != 0x03 {
		return nil, fmt.Errorf("expected compressed public key format, expecting 0x02 or 0x03, got %d", content[0])
	}

	ySign := content[0] & 0x01
	x := content[1:33]

	X := new(big.Int).SetBytes(x)
	Y, err := decompressPoint(X, ySign == 1)
	if err != nil {
		return nil, fmt.Errorf("unable to decompress compressed publick key material: %w", err)
	}

	return &btcec.PublicKey{
		Curve: elliptic.P256(),
		X:     X,
		Y:     Y,
	}, nil
}

func (p *innerWAPublicKey) keyMaterialSize() *int {
	return nil
}

func (p *innerWAPublicKey) prefix() string {
	return PublicKeyWAPrefix
}

// We use 3 a couple of times in the decompress algorithm below
var B3 = big.NewInt(3)

// decompressPoint Y derived from the point X and curve information. Adapted
// from an algorithm found on StackOverflow.
//
// @see https://stackoverflow.com/a/46289709/697930
// @see https://crypto.stackexchange.com/a/8916
func decompressPoint(x *big.Int, ybit bool) (*big.Int, error) {
	// The equation we need to solve is `y^2 = x^3 + aX + b`, in for elliptic curves
	// here, `a` is always `-3` so we need to solve `y^2 = x^3 - 3x + b`. We have `x`,
	// let's compute `y` from that information. We need to remember that all computation
	// must be done modulo the field prime of the group.
	//
	// This is probably not an optimal algorithm, just a working one.
	c := elliptic.P256().Params()

	// First, x^3, mod P
	xCubed := new(big.Int).Exp(x, B3, c.P)

	// Next, 3x, mod P
	threeX := new(big.Int).Mul(x, B3)
	threeX.Mod(threeX, c.P)

	// x^3 - 3x ...
	ySquared := new(big.Int).Sub(xCubed, threeX)

	// ... + b mod P
	ySquared.Add(ySquared, c.B)
	ySquared.Mod(ySquared, c.P)

	// Now we need to find the square root mod P, hopefully, Go's big int library has it!
	y := new(big.Int).ModSqrt(ySquared, c.P)
	if y == nil {
		return nil, fmt.Errorf("unable to compute square root of Y, we seem to have and invalid curve point X")
	}

	// Finally, check if you have the correct root by comparing
	// the low bit with the low bit of the sign byte. If it’s not
	// the same you want -y mod P instead of y.
	if ybit != isOdd(y) {
		y.Sub(c.P, y)
	}

	if ybit != isOdd(y) {
		return nil, fmt.Errorf("ybit doesn't match oddness")
	}

	return y, nil
}

func isOdd(a *big.Int) bool {
	return a.Bit(0) == 1
}
