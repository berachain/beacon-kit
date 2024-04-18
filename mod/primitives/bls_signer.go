package primitives

// BLSSigner defines an interface for cryptographic signing operations.
// It uses generic type parameters Signature and Pubkey, both of which are
// slices of bytes.
type BLSSigner interface {
	// PublicKey returns the public key of the signer.
	PublicKey() BLSPubkey

	// Sign takes a message as a slice of bytes and returns a signature as a
	// slice of bytes and an error.
	Sign([]byte) (BLSSignature, error)
}
