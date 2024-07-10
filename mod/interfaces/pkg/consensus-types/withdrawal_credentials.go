package types

import (
	"encoding"
	"encoding/json"

	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
)

type WithdrawalCredentials[T any] interface {
	json.Unmarshaler
	encoding.TextMarshaler
	encoding.TextUnmarshaler
	// NewFromExecutionAddress creates a new withdrawal credentials from an execution address.
	NewFromExecutionAddress(address gethprimitives.ExecutionAddress) T
	// ToExecutionAddress converts the withdrawal credentials to an execution address.
	ToExecutionAddress() (gethprimitives.ExecutionAddress, error)
}
