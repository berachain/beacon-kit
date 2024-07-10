package types

import (
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/merkleizer"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

type ExecutionPayload[
	T any,
	ExecutionPayloadHeaderT ExecutionPayloadHeader[ExecutionPayloadHeaderT],
	WithdrawalT any,
] interface {
	InnerExecutionPayload[WithdrawalT]
	constraints.EmptyWithVersion[T]
	// ToHeader converts the execution payload to a header.
	ToHeader(
		txsMerkleizer *merkleizer.Merkleizer[[32]byte, common.Root],
		maxWithdrawalsPerPayload uint64,
	) (ExecutionPayloadHeaderT, error)
}

// InnerExecutionPayload represents the inner execution payload.
type InnerExecutionPayload[
	WithdrawalT any,
] interface {
	ExecutionPayloadBody
	// GetTransactions returns the transactions of the execution payload.
	GetTransactions() [][]byte
	// GetWithdrawals returns the withdrawals of the execution payload.
	GetWithdrawals() []WithdrawalT
}

// ExecutionPayloadBody is the interface for the execution data of a block.
type ExecutionPayloadBody interface {
	constraints.SSZMarshallable
	constraints.JSONMarshallable
	constraints.Nillable
	constraints.Versionable
	// GetPrevRandao returns the previous randao of the execution payload.
	GetPrevRandao() common.Bytes32
	// GetBlockHash returns the block hash of the execution payload.
	GetBlockHash() gethprimitives.ExecutionHash
	// GetParentHash returns the parent hash of the execution payload.
	GetParentHash() gethprimitives.ExecutionHash
	// GetNumber returns the number of the execution payload.
	GetNumber() math.U64
	// GetGasLimit returns the gas limit of the execution payload.
	GetGasLimit() math.U64
	// GetGasUsed returns the gas used of the execution payload.
	GetGasUsed() math.U64
	// GetTimestamp returns the timestamp of the execution payload.
	GetTimestamp() math.U64
	// GetExtraData returns the extra data of the execution payload.
	GetExtraData() []byte
	// GetBaseFeePerGas returns the base fee per gas of the execution payload.
	GetBaseFeePerGas() math.Wei
	// GetFeeRecipient returns the fee recipient of the execution payload.
	GetFeeRecipient() gethprimitives.ExecutionAddress
	// GetStateRoot returns the state root of the execution payload.
	GetStateRoot() common.Bytes32
	// GetReceiptsRoot returns the receipts root of the execution payload.
	GetReceiptsRoot() common.Bytes32
	// GetLogsBloom returns the logs bloom of the execution payload.
	GetLogsBloom() []byte
	// GetBlobGasUsed returns the blob gas used of the execution payload.
	GetBlobGasUsed() math.U64
	// GetExcessBlobGas returns the excess blob gas of the execution payload.
	GetExcessBlobGas() math.U64
}
