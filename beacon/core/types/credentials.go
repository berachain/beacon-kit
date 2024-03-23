package types

import "github.com/berachain/beacon-kit/primitives"

// DepositCredentials is a staking credential that is used to identify a
// validator.
type DepositCredentials [32]byte

// NewCredentialsFromExecutionAddress creates a new DepositCredentials from an.
func NewCredentialsFromExecutionAddress(
	address primitives.ExecutionAddress,
) DepositCredentials {
	credentials := DepositCredentials{}
	credentials[0] = 0x01
	copy(credentials[12:], address[:])
	return credentials
}

// ToExecutionAddress converts the DepositCredentials to an ExecutionAddress.
func (c DepositCredentials) ToExecutionAddress() (
	primitives.ExecutionAddress,
	error,
) {
	if c[0] != byte(EthSecp256k1CredentialPrefix) {
		return primitives.ExecutionAddress{}, ErrInvalidDepositCredentials
	}
	return primitives.ExecutionAddress(c[12:]), nil
}
