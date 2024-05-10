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

type BeaconStorageBackend[
	BeaconStateT any, BlobSidecarsT BlobSidecars,
] interface {
	AvailabilityStore(
		context.Context,
	) core.AvailabilityStore[
		types.BeaconBlockBody, BlobSidecarsT,
	]
	BeaconState(context.Context) BeaconStateT
}

// BlobsSidecars is the interface for blobs sidecars.
type BlobSidecars interface {
	ssz.Marshaler
	ssz.Unmarshaler
	Len() int
}

// BlockVerifier is the interface for the block verifier.
type BlockVerifier[BeaconStateT any] interface {
	ValidateBlock(
		st BeaconStateT,
		blk types.ReadOnlyBeaconBlock[types.BeaconBlockBody],
	) error
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
type LocalBuilder[BeaconStateT any] interface {
	RequestPayload(
		ctx context.Context,
		st BeaconStateT,
		slot math.Slot,
		timestamp uint64,
		parentBlockRoot primitives.Root,
		parentEth1Hash common.ExecutionHash,
	) (*engineprimitives.PayloadID, error)
}

// PayloadVerifier is the interface for the payload verifier.
type PayloadVerifier[BeaconStateT any] interface {
	VerifyPayload(
		st BeaconStateT,
		payload engineprimitives.ExecutionPayload,
	) error
}

// RandaoProcessor is the interface for the randao processor.
type RandaoProcessor[BeaconStateT any] interface {
	BuildReveal(
		st BeaconStateT,
	) (crypto.BLSSignature, error)
	MixinNewReveal(
		st BeaconStateT,
		reveal crypto.BLSSignature,
	) error
	VerifyReveal(
		st BeaconStateT,
		proposerPubkey crypto.BLSPubkey,
		reveal crypto.BLSSignature,
	) error
}

// StakingService is the interface for the staking service.
type StakingService interface {
	// ProcessLogsInETH1Block processes logs in an eth1 block.
	ProcessLogsInETH1Block(
		ctx context.Context,
		blockHash common.ExecutionHash,
	) error

	// PruneDepositEvents prunes deposit events.
	// TODO: decouple.
	PruneDepositEvents(idx uint64) error
}
