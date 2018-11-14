package ecc

type innerR1Signature struct {
}

func (s innerR1Signature) verify(content []byte, hash []byte, pubKey PublicKey) bool {
	//recoveredKey, _, err := btcec.RecoverCompact(btcec.S256(), content, hash)
	//if err != nil {
	//	return false
	//}
	//key, err := pubKey.key()
	//if err != nil {
	//	return false
	//}
	//if recoveredKey.IsEqual(key) {
	//	return true
	//}
	return false
}

func (s *innerR1Signature) publicKey(content []byte, hash []byte) (out PublicKey, err error) {

	//var recoveredKey *btcec.publicKey
	//switch s.Curve {
	//case CurveK1:
	//	recoveredKey, _, err = btcec.RecoverCompact(btcec.S256(), s.Content, hash)
	//case CurveR1:
	//	curve := btcec.S256R1()
	//
	//	recoveredKey, _, err = btcec.RecoverCompact(curve, s.Content, hash)
	//default:
	//	return publicKey{}, fmt.Errorf("invalid curve: %s", s.Curve)
	//}
	//
	//if err != nil {
	//	return out, err
	//}
	//
	//return publicKey{
	//	Curve:   s.Curve,
	//	Content: recoveredKey.SerializeCompressed(),
	//}, nil

	return PublicKey{
		Curve:   CurveR1,
		Content: nil,
		inner:   nil,
	}, nil
}
