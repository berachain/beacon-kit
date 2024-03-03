package bls12381

import "github.com/prysmaticlabs/prysm/v5/crypto/bls"

func VerifySignature(
	pubKey [PubKeyLength]byte,
	msg []byte,
	signature [SignatureLength]byte,
) bool {
	pub, err := bls.PublicKeyFromBytes(pubKey[:])
	if err != nil {
		panic(err)
	}

	sig, err := bls.SignatureFromBytes(signature[:])
	if err != nil {
		panic(err)
	}

	return sig.Verify(pub, msg)
}
