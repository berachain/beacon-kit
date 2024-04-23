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

	"github.com/berachain/beacon-kit/mod/core/state"
	"github.com/berachain/beacon-kit/mod/primitives"
	engineprimitives "github.com/berachain/beacon-kit/mod/primitives-engine"
	"github.com/berachain/beacon-kit/mod/primitives/math"
)

type ExecutionEngine interface {
	// GetPayload returns the payload and blobs bundle for the given slot.
	GetPayload(
		ctx context.Context,
		req *engineprimitives.GetPayloadRequest,
	) (engineprimitives.BuiltExecutionPayload, error)

	// NotifyForkchoiceUpdate notifies the execution client of a forkchoice
	// update.
	NotifyForkchoiceUpdate(
		ctx context.Context,
		req *engineprimitives.ForkchoiceUpdateRequest,
	) (*engineprimitives.PayloadID, *primitives.ExecutionHash, error)

	// VerifyAndNotifyNewPayload verifies the new payload and notifies the
	// execution
	VerifyAndNotifyNewPayload(
		ctx context.Context,
		req *engineprimitives.NewPayloadRequest,
	) (bool, error)
}

// LocalBuilder is the interface for the builder service.
type LocalBuilder interface {
	RequestPayload(
		ctx context.Context,
		st state.BeaconState,
		parentEth1Hash primitives.ExecutionHash,
		slot math.Slot,
		timestamp uint64,
		parentBlockRoot primitives.Root,
	) (*engineprimitives.PayloadID, error)
}

// RandaoProcessor is the interface for the randao processor.
type RandaoProcessor interface {
	BuildReveal(
		st state.BeaconState,
	) (primitives.BLSSignature, error)
	MixinNewReveal(
		st state.BeaconState,
		reveal primitives.BLSSignature,
	) error
	VerifyReveal(
		st state.BeaconState,
		proposerPubkey primitives.BLSPubkey,
		reveal primitives.BLSSignature,
	) error
}

// StakingService is the interface for the staking service.
type StakingService interface {
	// ProcessLogsInETH1Block processes logs in an eth1 block.
	ProcessLogsInETH1Block(
		ctx context.Context,
		st state.BeaconState,
		blockHash primitives.ExecutionHash,
	) error
}
