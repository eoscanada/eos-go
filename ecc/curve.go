package ecc

type CurveID uint8

const (
	CurveK1 = CurveID(iota)
	CurveR1
	CurveWA
)

func (c CurveID) String() string {
	switch c {
	case CurveK1:
		return "K1"
	case CurveR1:
		return "R1"
	case CurveWA:
		return "WA"
	default:
		return "UN" // unknown
	}
}

func (c CurveID) StringPrefix() string {
	return c.String() + "_"
}
