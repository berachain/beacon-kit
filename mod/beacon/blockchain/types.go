// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package blockchain

import (
	"context"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	engineprimitives "github.com/berachain/beacon-kit/mod/primitives-engine"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
	ssz "github.com/ferranbt/fastssz"
)

// The AvailabilityStore interface is responsible for validating and storing
// sidecars for specific blocks, as well as verifying sidecars that have already
// been stored.
type AvailabilityStore[BeaconBlockBodyT any, BlobSidecarsT any] interface {
	// IsDataAvailable ensures that all blobs referenced in the block are
	// securely stored before it returns without an error.
	IsDataAvailable(
		context.Context, math.Slot, BeaconBlockBodyT,
	) bool
	// Persist makes sure that the sidecar remains accessible for data
	// availability checks throughout the beacon node's operation.
	Persist(math.Slot, BlobSidecarsT) error
}

// ReadOnlyBeaconState defines the interface for accessing various components of
// the
// beacon state.
type ReadOnlyBeaconState[T any] interface {
	// GetSlot retrieves the current slot of the beacon state.
	GetSlot() (math.Slot, error)

	// GetLatestExecutionPayloadHeader returns the most recent execution payload
	// header.
	GetLatestExecutionPayloadHeader() (
		engineprimitives.ExecutionPayloadHeader,
		error,
	)

	// GetEth1DepositIndex returns the index of the most recent eth1 deposit.
	GetEth1DepositIndex() (uint64, error)

	// GetLatestBlockHeader returns the most recent block header.
	GetLatestBlockHeader() (
		*types.BeaconBlockHeader,
		error,
	)

	// HashTreeRoot returns the hash tree root of the beacon state.
	HashTreeRoot() ([32]byte, error)

	// Copy creates a copy of the beacon state.
	Copy() T

	// ValidatorIndexByPubkey finds the index of a validator based on their
	// public key.
	ValidatorIndexByPubkey(crypto.BLSPubkey) (math.ValidatorIndex, error)
}

type StorageBackend[
	BeaconStateT any,
	BlobSidecarsT BlobSidecars,
	DepositStoreT DepositStore,
] interface {
	AvailabilityStore(
		context.Context,
	) core.AvailabilityStore[
		types.BeaconBlockBody, BlobSidecarsT,
	]
	StateFromContext(context.Context) BeaconStateT
	DepositStore(context.Context) DepositStoreT
}

// BlobVerifier is the interface for the blobs processor.
type BlobProcessor[
	AvailabilityStoreT AvailabilityStore[types.BeaconBlockBody, BlobSidecarsT],
	BlobSidecarsT any,
] interface {
	// ProcessBlobs processes the blobs and ensures they match the local state.
	ProcessBlobs(
		slot math.Slot,
		avs AvailabilityStoreT,
		sidecars BlobSidecarsT,
	) error
}

// BlobsSidecars is the interface for blobs sidecars.
type BlobSidecars interface {
	ssz.Marshaler
	ssz.Unmarshaler
	Len() int
}

// DepositContract is the ABI for the deposit contract.
type DepositContract interface {
	GetDeposits(
		ctx context.Context,
		blockNumber uint64,
	) ([]*types.Deposit, error)
}

// DepositStore defines the interface for managing deposit operations.
type DepositStore interface {
	// PruneToIndex prunes the deposit store up to the specified index.
	PruneToIndex(index uint64) error

	// EnqueueDeposits adds a list of deposits to the deposit store.
	EnqueueDeposits(deposits []*types.Deposit) error
}

type ExecutionEngine interface {
	// GetPayload returns the payload and blobs bundle for the given slot.
	GetPayload(
		ctx context.Context,
		req *engineprimitives.GetPayloadRequest,
	) (engineprimitives.BuiltExecutionPayloadEnv, error)

	// NotifyForkchoiceUpdate notifies the execution client of a forkchoice
	// update.
	NotifyForkchoiceUpdate(
		ctx context.Context,
		req *engineprimitives.ForkchoiceUpdateRequest,
	) (*engineprimitives.PayloadID, *common.ExecutionHash, error)

	// VerifyAndNotifyNewPayload verifies the new payload and notifies the
	// execution client.
	VerifyAndNotifyNewPayload(
		ctx context.Context,
		req *engineprimitives.NewPayloadRequest[types.ExecutionPayload],
	) error
}

// LocalBuilder is the interface for the builder service.
type LocalBuilder[ReadOnlyBeaconStateT any] interface {
	// RequestPayload requests a new payload for the given slot.
	RequestPayload(
		ctx context.Context,
		st ReadOnlyBeaconStateT,
		slot math.Slot,
		timestamp uint64,
		parentBlockRoot primitives.Root,
		parentEth1Hash common.ExecutionHash,
	) (*engineprimitives.PayloadID, error)
}

// PayloadVerifier is the interface for the payload verifier.
type PayloadVerifier[ReadOnlyBeaconStateT any] interface {
	VerifyPayload(
		st ReadOnlyBeaconStateT,
		payload engineprimitives.ExecutionPayload,
	) error
}

// RandaoProcessor is the interface for the randao processor.
type RandaoProcessor[ReadOnlyBeaconStateT any] interface {
	// BuildReveal generates a new RANDAO reveal for the given state.
	BuildReveal(
		st ReadOnlyBeaconStateT,
	) (crypto.BLSSignature, error)

	// MixinNewReveal mixes a new RANDAO reveal into the given state.
	MixinNewReveal(
		st ReadOnlyBeaconStateT,
		reveal crypto.BLSSignature,
	) error

	// VerifyReveal verifies the provided RANDAO reveal against the
	// given state and proposer public key.
	VerifyReveal(
		st ReadOnlyBeaconStateT,
		proposerPubkey crypto.BLSPubkey,
		reveal crypto.BLSSignature,
	) error
}

// StateProcessor defines the interface for processing various state transitions
// in the beacon chain.
type StateProcessor[ReadOnlyBeaconStateT, BlobSidecarsT any] interface {
	// ProcessSlot processes the state transition for a single slot.
	ProcessSlot(
		st ReadOnlyBeaconStateT,
	) error

	// Transition processes the state transition for a given block.
	Transition(
		ctx core.Context,
		st ReadOnlyBeaconStateT,
		blk types.BeaconBlock,
	) error
}
