package ecc

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"errors"
	"math/big"

	"github.com/eoscanada/eos-go/btcsuite/btcd/btcec"
	"github.com/eoscanada/eos-go/btcsuite/btcutil/base58"
)

type innerWASignature struct {
}

func newInnerWASignature() innerSignature {
	return &innerWASignature{}
}

func (s innerWASignature) verify(content []byte, hash []byte, pubKey PublicKey) bool {
	key, err := pubKey.Key()
	if err != nil {
		// FIXME: How to log that to the other world?
		return false
	}

	R := new(big.Int).SetBytes(content[1:33])
	S := new(big.Int).SetBytes(content[33:65])

	return ecdsa.Verify(
		key.ToECDSA(),
		hash,
		R,
		S,
	)
}

func (s *innerWASignature) publicKey(content []byte, hash []byte) (out PublicKey, err error) {
	// This is the start of the algorithm to extract the PublicKey from a SIG_WA signature.
	// This was extracted from the `btcec` package which works on a different curve settings.
	// This is untested, but from my (small) comprehension of Elliptic Curves, this should
	// still work for ECDSA P-256 curves.
	//
	// However, one remaining step is re-construction of the full WA public key material
	// which contains more than just the X and Y points in the WA case on EOSIO. Indeed, the
	// key material also contain the flag presence and the relay party ID in string format.
	// This information are presented in the signature `content` after the R and S values of
	// the actual signature. But they are encoded in `cbor` encoding. To finish this, we need
	// to:
	//   1- Extract the attestation object from the signature
	//   2- Extract the client data JSON object from the signature
	//   3- Extract the user presence flag from the attestation object, need to be decoded against cbor algorithm first
	//   4- Extract the relay party ID from the client data JSON
	//   5- Embed everything into the proper WA public key encoding scheme.

	// R := new(big.Int).SetBytes(content[1:33])
	// S := new(big.Int).SetBytes(content[33:65])

	// iteration := int((content[0] - 27) & ^byte(4))

	// key, err := recoverKeyFromSignature(elliptic.P256(), R, S, hash, iteration, false)

	return out, errors.New("WA not supported")
}

func (s innerWASignature) string(content []byte) string {
	checksum := ripemd160checksumHashCurve(content, CurveWA)
	buf := append(content[:], checksum...)
	return "SIG_WA_" + base58.Encode(buf)
}

func (s innerWASignature) signatureMaterialSize() *int {
	return nil
}

// recoverKeyFromSignature recovers a public key from the signature "sig" on the
// given message hash "msg". Based on the algorithm found in section 5.1.5 of
// SEC 1 Ver 2.0, page 47-48 (53 and 54 in the pdf). This performs the details
// in the inner loop in Step 1. The counter provided is actually the j parameter
// of the loop * 2 - on the first iteration of j we do the R case, else the -R
// case in step 1.6. This counter is used in the bitcoin compressed signature
// format and thus we match bitcoind's behaviour here.
//
// This was copied straight (with slight modifications) from `btcec.recoverKeySignature`
// seems it looks the algorithm is valid for all Elliptic curves. Not tested yet.
func recoverKeyFromSignature(curve elliptic.Curve, R, S *big.Int, msg []byte, iter int, doChecks bool) (*btcec.PublicKey, error) {
	// 1.1 x = (n * i) + r
	Rx := new(big.Int).Mul(curve.Params().N,
		new(big.Int).SetInt64(int64(iter/2)))
	Rx.Add(Rx, R)
	if Rx.Cmp(curve.Params().P) != -1 {
		return nil, errors.New("calculated Rx is larger than curve P")
	}

	// convert 02<Rx> to point R. (step 1.2 and 1.3). If we are on an odd
	// iteration then 1.6 will be done with -R, so we calculate the other
	// term when uncompressing the point.
	// FIXME: Shall we have the curve? I have hard-coded it as P256
	Ry, err := decompressPoint(Rx, iter%2 == 1)
	if err != nil {
		return nil, err
	}

	// 1.4 Check n*R is point at infinity
	if doChecks {
		nRx, nRy := curve.ScalarMult(Rx, Ry, curve.Params().N.Bytes())
		if nRx.Sign() != 0 || nRy.Sign() != 0 {
			return nil, errors.New("n*R does not equal the point at infinity")
		}
	}

	// 1.5 calculate e from message using the same algorithm as ecdsa
	// signature calculation.
	e := hashToInt(msg, curve)

	// Step 1.6.1:
	// We calculate the two terms sR and eG separately multiplied by the
	// inverse of r (from the signature). We then add them to calculate
	// Q = r^-1(sR-eG)
	invr := new(big.Int).ModInverse(R, curve.Params().N)

	// first term.
	invrS := new(big.Int).Mul(invr, S)
	invrS.Mod(invrS, curve.Params().N)
	sRx, sRy := curve.ScalarMult(Rx, Ry, invrS.Bytes())

	// second term.
	e.Neg(e)
	e.Mod(e, curve.Params().N)
	e.Mul(e, invr)
	e.Mod(e, curve.Params().N)
	minuseGx, minuseGy := curve.ScalarBaseMult(e.Bytes())

	// TODO: this would be faster if we did a mult and add in one
	// step to prevent the jacobian conversion back and forth.
	Qx, Qy := curve.Add(sRx, sRy, minuseGx, minuseGy)

	return &btcec.PublicKey{
		Curve: curve,
		X:     Qx,
		Y:     Qy,
	}, nil
}

// hashToInt converts a hash value to an integer. There is some disagreement
// about how this is done. [NSA] suggests that this is done in the obvious
// manner, but [SECG] truncates the hash to the bit-length of the curve order
// first. We follow [SECG] because that's what OpenSSL does. Additionally,
// OpenSSL right shifts excess bits from the number if the hash is too large
// and we mirror that too.
// This is borrowed from crypto/ecdsa.
func hashToInt(hash []byte, c elliptic.Curve) *big.Int {
	orderBits := c.Params().N.BitLen()
	orderBytes := (orderBits + 7) / 8
	if len(hash) > orderBytes {
		hash = hash[:orderBytes]
	}

	ret := new(big.Int).SetBytes(hash)
	excess := len(hash)*8 - orderBits
	if excess > 0 {
		ret.Rsh(ret, uint(excess))
	}
	return ret
}
