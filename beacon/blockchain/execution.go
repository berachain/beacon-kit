// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
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
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/ethereum/go-ethereum/common"
	"github.com/itsdevbear/bolaris/beacon/execution"
	consensusv1 "github.com/itsdevbear/bolaris/types/consensus/v1"
	"github.com/itsdevbear/bolaris/types/consensus/v1/interfaces"
	prsymexecution "github.com/prysmaticlabs/prysm/v4/beacon-chain/execution"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	"golang.org/x/sync/errgroup"
)

// GetOrBuildBlock constructs the next block in the execution chain.
func (s *Service) GetOrBuildBlock(
	ctx context.Context, slot primitives.Slot,
) (interfaces.BeaconKitBlock, error) {
	// The goal here is to acquire a payload whose parent is the previously
	// finalized block, such that, if this payload is accepted, it will be
	// the next finalized block in the chain. A byproduct of this design
	// is that we get the nice property of lazily propogating the finalized
	// and safe block hashes to the execution client.
	var (
		beaconState        = s.bsp.BeaconState(ctx)
		err                error
		executionData      interfaces.ExecutionData
		lastFinalizedBlock = beaconState.GetFinalizedEth1BlockHash()
	)

	// Attempt to get a previously built payload, otherwise we will trigger a payload to be build.
	// Building a payload at this point is not ideal, as it will block the consensus client
	// from proposing a block until the payload is built. However, this is a rare case, and
	// we can optimize this later.
	// TODO: 4844.
	if executionData, _, _, err = s.en.GetBuiltPayload(
		ctx, slot, lastFinalizedBlock,
	); err != nil || executionData == nil {
		// This branch represents a cache miss.
		telemetry.IncrCounter(1, MetricGetBuiltPayloadMiss)
		executionData, err = s.buildNewPayloadAtSlotWithParent(ctx, slot, lastFinalizedBlock)
		if err != nil {
			return nil, err
		}
	} else {
		// This branch represents a cache hit.
		telemetry.IncrCounter(1, MetricGetBuiltPayloadHit)
	}

	// Create a new block with the payload.
	return consensusv1.NewBaseBeaconKitBlockFromState(
		beaconState, executionData,
		s.BeaconCfg().ActiveForkVersion(primitives.Epoch(slot)),
	)
}

// buildNewBlockOnTopOf constructs a new block on top of an existing head of the execution client.
func (s *Service) buildNewPayloadAtSlotWithParent(
	ctx context.Context, slot primitives.Slot, headHash common.Hash,
) (interfaces.ExecutionData, error) {
	finalHash := s.bsp.BeaconState(ctx).GetFinalizedEth1BlockHash()
	safeHash := s.bsp.BeaconState(ctx).GetSafeEth1BlockHash()

	// Notify the execution client of the new finalized and safe hashes, while also
	// triggering a block to be built. We wait for the payload ID to be returned.
	err := s.en.NotifyForkchoiceUpdate(
		ctx, slot,
		execution.NewNotifyForkchoiceUpdateArg(
			headHash, safeHash, finalHash,
		),
		true,
		true,
		false,
	)

	if err != nil {
		s.Logger().Error("Failed to notify forkchoice update",
			"finalized_hash", finalHash,
			"safe_hash", safeHash,
			"head_hash", headHash,
			"error", err)
		return nil, err
	}

	return s.waitForPayload(ctx, payloadBuildDelay*time.Second, slot, headHash)
}

// waitForPayload waits for a payload to be built by the execution client.
func (s *Service) waitForPayload(
	ctx context.Context, delay time.Duration, slot primitives.Slot, headHash common.Hash,
) (interfaces.ExecutionData, error) {
	// Create a timer for the delay duration
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		// The context was cancelled before the timer expired
		s.Logger().Error(
			"Context cancelled while waiting for payload", "slot", slot, "head_hash", headHash,
		)
		return nil, ctx.Err()
	case <-timer.C:
		// Hitting this case is basically the desired behaviour.
		break
	}

	// If neither context cancellation nor timeout occurred, proceed to get the payload
	payload, _, _, err := s.en.GetBuiltPayload(ctx, slot, headHash)
	if err != nil {
		s.Logger().Error(
			"Failed to get built payload", "error", err, "slot", slot, "head_hash", headHash,
		)
		return nil, err
	}
	return payload, nil
}

func (s *Service) ProcessReceivedBlock(
	ctx context.Context,
	block interfaces.ReadOnlyBeaconKitBlock,
) error {
	// If we get any sort of error from the execution client, we bubble
	// it up and reject the proposal, as we do not want to write a block
	// finalization to the consensus layer that is invalid.
	var (
		eg, groupCtx   = errgroup.WithContext(ctx)
		isValidPayload bool
	)

	// TODO: Do we need to wait for the forkchoice to update?
	// TODO: move the error group to use GCD.

	// var postState state.BeaconState
	eg.Go(func() error {
		err := s.validateStateTransition(groupCtx, block)
		if err != nil {
			s.Logger().Error("failed to validate state transition", "error", err)
			return err
		}
		return nil
	})

	eg.Go(func() error {
		var err error
		if isValidPayload, err = s.validateExecutionOnBlock(
			groupCtx, block.ExecutionData(),
		); err != nil {
			s.Logger().Error("failed to notify engine of new payload", "error", err)
			return err
		}
		return nil
	})

	// Wait for the goroutines to finish.
	if err := eg.Wait(); err != nil {
		return err
	}

	if err := s.postBlockProcess(
		ctx, block /*blockCopy, blockRoot, postState,*/, isValidPayload,
	); err != nil {
		return err
	}

	return nil
}

// validateStateTransition validates the state transition of a given block.
// TODO: add more rules and modularize, I am unsure if this is the best / correct place for this.
// It's also not very modular, its just hardcoded to single slot finality for now, which is fine,
// but maybe not the most extensible.
func (s *Service) validateStateTransition(
	ctx context.Context, block interfaces.ReadOnlyBeaconKitBlock,
) error {
	parentHash := block.ExecutionData().ParentHash()
	finalizedHash := s.bsp.BeaconState(ctx).GetFinalizedEth1BlockHash()
	if !bytes.Equal(finalizedHash[:], parentHash) {
		return fmt.Errorf(
			"parent block with hash %x is not finalized, expected finalized hash %x",
			parentHash, finalizedHash,
		)
	}
	return nil
}

// ValidateProposedBeaconBlock checks the validity of a proposed beacon block.
func (s *Service) validateExecutionOnBlock(ctx context.Context, header interfaces.ExecutionData,
) (bool, error) {
	isValidPayload, err := s.en.NotifyNewPayload(ctx, 0, header)
	if err != nil && errors.Is(err, prsymexecution.ErrAcceptedSyncingPayloadStatus) {
		s.Logger().Error("Failed to validate execution on block", "error", err)
		return isValidPayload, err
	} else if err != nil || !isValidPayload {
		return isValidPayload, prsymexecution.ErrInvalidPayloadStatus
	}
	return isValidPayload, nil
}

func (s *Service) postBlockProcess(
	ctx context.Context, block interfaces.ReadOnlyBeaconKitBlock, isValidPayload bool,
) error {
	if !isValidPayload {
		telemetry.IncrCounter(1, MetricReceivedInvalidPayload)
		return ErrInvalidPayload
	}

	// TODO: don't get slot off the execution data incase it's incorrect somehow.
	slot := primitives.Slot(block.ExecutionData().BlockNumber())

	// We notify the execution client of the new block, and wait for it to return
	// a payload ID. If the payload ID is nil, we return an error. One thing to notice here however
	// is that we pass in `slot+1` to the execution client. We do this so that we can begin building
	// the next block in the background while we are finalizing this block.
	// We are okay pushing this asynchonous work to the execution client, as it is
	return s.en.NotifyForkchoiceUpdate(
		ctx, slot+1,
		execution.NewNotifyForkchoiceUpdateArg(
			common.BytesToHash(block.ExecutionData().BlockHash()),
			s.bsp.BeaconState(ctx).GetSafeEth1BlockHash(),
			s.bsp.BeaconState(ctx).GetFinalizedEth1BlockHash(),
		), true, true, false)
}
