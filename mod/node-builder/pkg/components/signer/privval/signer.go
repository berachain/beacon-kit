package privval

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/cometbft/cometbft/privval"
	"github.com/itsdevbear/comet-bls12-381/bls/blst"
)

// FileBLSSigner holds the secret key used for signing operations.
type FileBLSigner struct {
	*privval.FilePV
}

// NewFileBLSigner creates a new FileBLSigner instance given a key file path.
func NewFileBLSigner(keyFilePath, stateFilePath string) (*FileBLSigner, error) {
	return &FileBLSigner{
		FilePV: privval.LoadOrGenFilePV(keyFilePath, stateFilePath),
	}, nil
}

// PublicKey returns the public key of the signer.
func (b *FileBLSigner) PublicKey() crypto.BLSPubkey {
	key, err := b.FilePV.GetPubKey()
	if err != nil {
		return crypto.BLSPubkey{}
	}
	return crypto.BLSPubkey(key.Bytes())
}

// Sign generates a signature for a given message using the signer's secret key.
// It returns the signature and any error encountered during the signing
// process.
func (b *FileBLSigner) Sign(msg []byte) (crypto.BLSSignature, error) {
	sig, err := b.FilePV.SignBytes(msg)
	if err != nil {
		return crypto.BLSSignature{}, err
	}
	return crypto.BLSSignature(sig), nil
}

// VerifySignature verifies a signature for a given message using the signer's
// public key.
func (FileBLSigner) VerifySignature(
	pubKey crypto.BLSPubkey,
	msg []byte,
	signature crypto.BLSSignature,
) error {
	pubkey, err := blst.PublicKeyFromBytes(pubKey[:])
	if err != nil {
		return err
	}

	sig, err := blst.SignatureFromBytes(signature[:])
	if err != nil {
		return err
	}

	if !sig.Verify(pubkey, msg) {
		return ErrInvalidSignature
	}
	return nil
}
