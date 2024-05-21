package privval

import (
	"errors"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/cometbft/cometbft/privval"
	"github.com/itsdevbear/comet-bls12-381/bls/blst"
)

// FileBLSSigner implements Privval using data persisted to disk to prevent
// double signing
type FileBLSSigner struct {
	*privval.FilePV
}

func NewFileBLSSigner(keyFilePath string, stateFilePath string) *FileBLSSigner {
	filePV := privval.LoadOrGenFilePV(keyFilePath, stateFilePath)
	return &FileBLSSigner{FilePV: filePV}
}

// ========================== Implements BLS Signer ==========================

// PublicKey returns the public key of the signer.
func (f *FileBLSSigner) PublicKey() crypto.BLSPubkey {
	key, err := f.FilePV.GetPubKey()
	if err != nil {
		return crypto.BLSPubkey{}
	}
	return crypto.BLSPubkey(key.Bytes())
}

// Sign generates a signature for a given message using the signer's secret key.
func (f *FileBLSSigner) Sign(msg []byte) (crypto.BLSSignature, error) {
	sig, err := f.FilePV.SignBytes(msg)
	if err != nil {
		return crypto.BLSSignature{}, err
	}
	return crypto.BLSSignature(sig), nil
}

// VerifySignature verifies a signature against a message and a public key.
func (f *FileBLSSigner) VerifySignature(
	blsPk crypto.BLSPubkey,
	msg []byte,
	signature crypto.BLSSignature) error {
	pk, err := blst.PublicKeyFromBytes(blsPk[:])
	if err != nil {
		return err
	}

	sig, err := blst.SignatureFromBytes(signature[:])
	if err != nil {
		return err
	}

	if !sig.Verify(pk, msg) {
		return errors.New("signature verification failed")
	}
	return nil
}
