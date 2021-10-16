package ecc

import "golang.org/x/crypto/ripemd160"

func ripemd160checksum(in []byte, curve CurveID) []byte {
	h := ripemd160.New()
	_, _ = h.Write(in) // this implementation has no error path

	if curve != CurveK1 {
		_, _ = h.Write([]byte(curve.String()))
	}

	sum := h.Sum(nil)
	return sum[:4]
}

func ripemd160checksumHashCurve(in []byte, curve CurveID) []byte {
	h := ripemd160.New()
	_, _ = h.Write(in) // this implementation has no error path

	// FIXME: this seems to be only rolled out to the `SIG_` things..
	// proper support for importing `EOS` keys isn't rolled out into `dawn4`.
	_, _ = h.Write([]byte(curve.String())) // conditionally ?
	sum := h.Sum(nil)
	return sum[:4]
}
