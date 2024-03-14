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

	randaotypes "github.com/berachain/beacon-kit/beacon/core/randao/types"
	beacontypes "github.com/berachain/beacon-kit/beacon/core/types"
	"github.com/berachain/beacon-kit/beacon/execution"
	bls12381 "github.com/berachain/beacon-kit/crypto/bls12-381"
	enginetypes "github.com/berachain/beacon-kit/engine/types"
	"github.com/berachain/beacon-kit/primitives"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

// ExecutionService is the interface for the execution service.
type ExecutionService interface {
	// NotifyForkchoiceUpdate notifies the execution client of a forkchoice
	// update.
	NotifyForkchoiceUpdate(
		ctx context.Context,
		fcuConfig *execution.FCUConfig,
	) (*enginetypes.PayloadID, error)

	// NotifyNewPayload notifies the execution client of a new payload.
	NotifyNewPayload(
		ctx context.Context,
		slot primitives.Slot,
		payload enginetypes.ExecutionPayload,
		versionedHashes []primitives.ExecutionHash,
		parentBlockRoot [32]byte,
	) (bool, error)

	// ProcessLogsInETH1Block processes logs in an eth1 block.
	ProcessLogsInETH1Block(
		ctx context.Context,
		blockHash primitives.ExecutionHash,
	) error
}

// LocalBuilder is the interface for the builder service.
type LocalBuilder interface {
	BuildLocalPayload(
		ctx context.Context,
		parentEth1Hash primitives.ExecutionHash,
		slot primitives.Slot,
		timestamp uint64,
		parentBlockRoot [32]byte,
	) (*enginetypes.PayloadID, error)
}

// RandaoProcessor is the interface for the randao processor.
type RandaoProcessor interface {
	BuildReveal(
		st context.Context,
		epoch primitives.Epoch,
	) (randaotypes.Reveal, error)
	MixinNewReveal(
		ctx context.Context,
		blk beacontypes.BeaconBlock,
	) error
	VerifyReveal(
		proposerPubkey [bls12381.PubKeyLength]byte,
		epoch primitives.Epoch,
		reveal randaotypes.Reveal,
	) error
}

// StakingService is the interface for the staking service.
type StakingService interface {
	ProcessBlockEvents(
		ctx context.Context,
		logs []ethtypes.Log,
	) error
}
type SyncService interface {
	IsInitSync() bool
	Status() error
}
