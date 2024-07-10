package types

import "github.com/berachain/beacon-kit/mod/primitives/pkg/common"

type Genesis[
	DepositT any,
	ExecutionPayloadHeaderT ExecutionPayloadHeader[ExecutionPayloadHeaderT],
] interface {
	// GetForkVersion returns the fork version of the genesis.
	GetForkVersion() common.Version
	// GetDeposits returns the deposits of the genesis.
	GetDeposits() []DepositT
	// GetExecutionPayloadHeader returns the execution payload header of the genesis.
	GetExecutionPayloadHeader() ExecutionPayloadHeaderT
	// UnmarshalJSON unmarshals the genesis from JSON.
	UnmarshalJSON([]byte) error
}
