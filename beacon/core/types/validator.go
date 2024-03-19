package types

// Validator is a struct that represents a validator in the beacon chain.
// Validator represents a participant in the beacon chain consensus mechanism.
// It holds the validator's public key, withdrawal credentials, effective balance, and slashing status.
type Validator struct {
	// Pubkey is the validator's 48-byte BLS public key.
	Pubkey [48]byte `json:"pubkey" ssz-size:"48"`
	// Credentials are address that controls the validator.
	Credentials [32]byte `json:"withdrawal_credentials" ssz-size:"32"`
	// EffectiveBalance is the validator's current effective balance in gwei.
	EffectiveBalance uint64 `json:"effective_balance"`
	// Slashed indicates whether the validator has been slashed.
	Slashed bool `json:"slashed"`
}
