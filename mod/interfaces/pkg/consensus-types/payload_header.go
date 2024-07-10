package types

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
)

type ExecutionPayloadHeader[T any] interface {
	InnerExecutionPayloadHeader
	constraints.EmptyWithVersion[T]
	constraints.NewFromSSZable[T]
	// NewFromJSON creates a new execution payload header from the given JSON bytes.)
	NewFromJSON(bz []byte, forkVersion uint32) (T, error)
}

type InnerExecutionPayloadHeader interface {
	ExecutionPayloadBody
	// GetTransactionsRoot returns the transactions root of the execution payload header.
	GetTransactionsRoot() common.Root
	// GetWithdrawalsRoot returns the withdrawals root of the execution payload header.
	GetWithdrawalsRoot() common.Root
}
