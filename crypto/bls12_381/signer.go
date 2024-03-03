package bls12381

import "github.com/prysmaticlabs/prysm/v5/crypto/bls"

type BlsSigner interface {
	Sign([]byte) ([SignatureLength]byte, error)
}

type blsSigner struct {
	secretKey [SecretKeyLength]byte
}

func NewBlsSigner(secretKey [SecretKeyLength]byte) BlsSigner {
	return &blsSigner{secretKey: secretKey}
}

func (b blsSigner) Sign(msg []byte) ([SignatureLength]byte, error) {
	privKey, err := bls.SecretKeyFromBytes(b.secretKey[:])
	if err != nil {
		return [SignatureLength]byte{}, err
	}

	sign := privKey.Sign(msg)

	var signature [SignatureLength]byte
	copy(signature[:], sign.Marshal())

	return signature, nil
}
