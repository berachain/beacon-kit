package blst

import (
	"github.com/berachain/comet-bls12-381/bls"
	"github.com/berachain/comet-bls12-381/bls/blst"

	"github.com/itsdevbear/bolaris/beacon/core/randao"
)

// Reveal represents the reveal of the validator for the current epoch.
type Reveal [randao.BLSSignatureLength]byte

func (r Reveal) Verify(pubKey []byte, signingData randao.SigningData) bool {
	bytes, err := blst.SignatureFromBytes(r[:])
	if err != nil {
		panic(err)
	}

	p, err := blst.PublicKeyFromBytes(pubKey)
	if err != nil {
		panic(err)
	}

	return bytes.Verify(p, signingData.Marshal())
}

// Marshal returns the reveal as a byte slice.
func (r Reveal) Marshal() []byte {
	return r[:]
}

// NewRandaoReveal creates the randao reveal for a given signing data and private key.
func NewRandaoReveal(signingData randao.SigningData, privKey bls.SecretKey) (randao.Reveal, error) {
	sig := privKey.Sign(signingData.Marshal())

	var reveal Reveal
	copy(reveal[:], sig.Marshal())

	return reveal, nil
}
