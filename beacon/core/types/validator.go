package types

// Validator is a struct that represents a validator in the beacon chain.
type Validator struct {
	Pubkey           [48]byte `json:"pubkey" ssz-size:"48"`
	Credentials      [32]byte `json:"withdrawal_credentials" ssz-size:"32"`
	EffectiveBalance uint64   `json:"effective_balance"`
	Slashed          bool     `json:"slashed"`
}
