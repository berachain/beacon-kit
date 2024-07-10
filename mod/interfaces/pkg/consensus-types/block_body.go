package types

import (
	// gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
)

// BeaconBlockBody is the interface for a beacon block body.
type BeaconBlockBody[
	T any,
	DepositT any,
	Eth1DataT any,
	ExecutionPayloadT any,
] interface {
	RawBeaconBlockBody[DepositT, Eth1DataT, ExecutionPayloadT]
	constraints.EmptyWithVersion[T]
}

// RawBeaconBlockBody is the interface for a raw beacon block body.
type RawBeaconBlockBody[
	DepositT any,
	Eth1DataT any,
	ExecutionPayloadT any,
] interface {
	WriteOnlyBeaconBlockBody[DepositT, Eth1DataT, ExecutionPayloadT]
	ReadOnlyBeaconBlockBody[DepositT, Eth1DataT, ExecutionPayloadT]
	// Length returns the length of the block body.
	Length() uint64
}

// WriteOnlyBeaconBlockBody is the interface for a write-only beacon block body.
type WriteOnlyBeaconBlockBody[
	DepositT any,
	Eth1DataT any,
	ExecutionPayloadT any,
] interface {
	// SetDeposits sets the deposits of the block.
	SetDeposits([]DepositT)
	// SetEth1Data sets the eth1 data of the block.
	SetEth1Data(Eth1DataT)
	// SetExecutionData sets the execution data of the block.
	SetExecutionData(ExecutionPayloadT) error
	// SetBlobKzgCommitments sets the KZG commitments of the block.
	SetBlobKzgCommitments(eip4844.KZGCommitments[gethprimitives.ExecutionHash])
	// SetRandaoReveal sets the RANDAO reveal of the block.
	SetRandaoReveal(crypto.BLSSignature)
	// SetGraffiti sets the graffiti of the block.
	SetGraffiti(common.Bytes32)
}

// ReadOnlyBeaconBlockBody is the interface for
// a read-only beacon block body.
type ReadOnlyBeaconBlockBody[
	DepositT any,
	Eth1DataT any,
	ExecutionPayloadT any,
] interface {
	constraints.SSZMarshallable
	constraints.Nillable
	// GetDeposits returns the deposits of the block.
	GetDeposits() []DepositT
	// GetEth1Data returns the eth1 data of the block.
	GetEth1Data() Eth1DataT
	// GetGraffiti returns the graffiti of the block.
	GetGraffiti() common.Bytes32
	// GetRandaoReveal returns the RANDAO reveal of the block.
	GetRandaoReveal() crypto.BLSSignature
	// GetExecutionPayload returns the execution data of the block.
	GetExecutionPayload() ExecutionPayloadT
	// GetBlobKzgCommitments returns the KZG commitments of the block.
	GetBlobKzgCommitments() eip4844.KZGCommitments[gethprimitives.ExecutionHash]
	// GetTopLevelRoots returns the top-level roots of the block.
	GetTopLevelRoots() ([][32]byte, error)
}
