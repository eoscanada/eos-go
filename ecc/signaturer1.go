package ecc

type InnerR1Signature struct {
}

func (s InnerR1Signature) Verify(content []byte, hash []byte, pubKey PublicKey) bool {
	//recoveredKey, _, err := btcec.RecoverCompact(btcec.S256(), content, hash)
	//if err != nil {
	//	return false
	//}
	//key, err := pubKey.Key()
	//if err != nil {
	//	return false
	//}
	//if recoveredKey.IsEqual(key) {
	//	return true
	//}
	return false
}

func (s *InnerR1Signature) PublicKey(content []byte, hash []byte) (out PublicKey, err error) {

	//var recoveredKey *btcec.PublicKey
	//switch s.Curve {
	//case CurveK1:
	//	recoveredKey, _, err = btcec.RecoverCompact(btcec.S256(), s.Content, hash)
	//case CurveR1:
	//	curve := btcec.S256R1()
	//
	//	recoveredKey, _, err = btcec.RecoverCompact(curve, s.Content, hash)
	//default:
	//	return PublicKey{}, fmt.Errorf("invalid curve: %s", s.Curve)
	//}
	//
	//if err != nil {
	//	return out, err
	//}
	//
	//return PublicKey{
	//	Curve:   s.Curve,
	//	Content: recoveredKey.SerializeCompressed(),
	//}, nil

	return PublicKey{}, nil
}
