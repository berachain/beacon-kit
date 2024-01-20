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

	"github.com/itsdevbear/bolaris/beacon/execution"
	consensusv1 "github.com/itsdevbear/bolaris/types/consensus/v1"
	"github.com/itsdevbear/bolaris/types/consensus/v1/interfaces"
	prsymexecution "github.com/prysmaticlabs/prysm/v4/beacon-chain/execution"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	"golang.org/x/sync/errgroup"
)

// BuildNextBlock constructs the next block in the blockchain.
func (s *Service) BuildNextBlock(
	ctx context.Context, slot primitives.Slot, time uint64,
) (interfaces.BeaconKitBlock, error) {
	// The goal here is to build a payload whose parent is the previously
	// finalized block, such that, if this payload is accepted, it will be
	// the next finalized block in the chain. A byproduct of this design
	// is that we get the nice property of lazily propogating the finalized
	// and safe block hashes to the execution client.
	lastFinalizedBlock := s.fcsp.ForkChoiceStore(ctx).GetFinalizedBlockHash()
	executionData, err := s.buildNewPayloadAtSlotWithParent(ctx, slot, lastFinalizedBlock[:])
	if err != nil {
		return nil, err
	}

	// Create a new block with the payload.
	return consensusv1.NewBaseBeaconKitBlock(
		slot, time, executionData,
		s.beaconCfg.ActiveForkVersion(primitives.Epoch(slot)),
	)
}

// buildNewBlockOnTopOf constructs a new block on top of an existing head of the execution client.
func (s *Service) buildNewPayloadAtSlotWithParent(
	ctx context.Context, slot primitives.Slot, headHash []byte,
) (interfaces.ExecutionData, error) {
	finalHash := s.fcsp.ForkChoiceStore(ctx).GetFinalizedBlockHash()
	safeHash := s.fcsp.ForkChoiceStore(ctx).GetSafeBlockHash()
	payloadIDBytes, err := s.en.NotifyForkchoiceUpdate(
		ctx, slot,
		execution.NewNotifyForkchoiceUpdateArg(
			headHash, safeHash[:], finalHash[:],
		),
		true,
		true,
	)

	if err != nil {
		s.logger.Error("Failed to notify forkchoice update",
			"finalized_hash", finalHash,
			"safe_hash", safeHash,
			"head_hash", headHash,
			"error", err)
		return nil, err
	}

	// TODO: Do we need to wait for the forkchoice to update?
	time.Sleep(payloadBuildDelay * time.Second)

	payload, _, _, err := s.en.GetBuiltPayload(
		ctx, slot,
	)
	if err != nil {
		s.logger.Error("Failed to get built payload", "error", err, "payload_id", payloadIDBytes)
		return nil, err
	}
	return payload, err
}

func (s *Service) ReceiveBlock(ctx context.Context, block interfaces.BeaconKitBlock) error {
	// If we get any sort of error from the execution client, we bubble
	// it up and reject the proposal, as we do not want to write a block
	// finalization to the consensus layer that is invalid.
	var (
		eg, groupCtx   = errgroup.WithContext(ctx)
		isValidPayload bool
	)

	// TODO: Do we need to wait for the forkchoice to update?

	// var postState state.BeaconState
	eg.Go(func() error {
		err := s.validateStateTransition(groupCtx, block)
		if err != nil {
			s.logger.Error("failed to validate state transition", "error", err)
			return err
			// return errors.Wrapf(err, "failed to validate consensus state transition function")
		}
		return nil
	})

	eg.Go(func() error {
		var err error
		if isValidPayload, err = s.validateExecutionOnBlock(
			groupCtx, block.ExecutionData(),
		); err != nil {
			s.logger.Error("failed to notify engine of new payload", "error", err)
			return err
		}
		return nil
	})

	// Wait for the goroutines to finish.
	if err := eg.Wait(); err != nil {
		return err
	}

	s.logger.Info("Validation complete", "isValidPayload", isValidPayload)

	if err := s.postBlockProcess(
		ctx, block /*blockCopy, blockRoot, postState,*/, isValidPayload,
	); err != nil {
		// err := errors.Wrap(err, "could not process block")
		// tracing.AnnotateError(span, err)
		return err
	}

	return nil
}

// validateStateTransition validates the state transition of a given block.
// TODO: add more rules and modularize, I am unsure if this is the best / correct place for this.
// It's also not very modular, its just hardcoded to single slot finality for now, which is fine,
// but maybe not the most extensible.
func (s *Service) validateStateTransition(
	ctx context.Context, block interfaces.BeaconKitBlock,
) error {
	parentHash := block.ExecutionData().ParentHash()
	finalizedHash := s.fcsp.ForkChoiceStore(ctx).GetFinalizedBlockHash()
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
	isValidPayload, err := s.en.NotifyNewPayload(ctx, 0, header /*, nil, [32]byte{}*/)
	if err != nil && errors.Is(err, prsymexecution.ErrAcceptedSyncingPayloadStatus) {
		s.logger.Error("Failed to validate execution on block", "error", err)
		return isValidPayload, err
	} else if err != nil || !isValidPayload {
		return isValidPayload, prsymexecution.ErrInvalidPayloadStatus
	}
	return isValidPayload, nil
}

func (s *Service) postBlockProcess(
	ctx context.Context, block interfaces.BeaconKitBlock, isValidPayload bool,
) error {
	// todo: don't get slot off the execution data incase it's incorrect somehow.
	_ = isValidPayload
	slot := primitives.Slot(block.ExecutionData().BlockNumber())
	headRoot := []byte{} // todo: get from forkchoice store. (stops proposer from double choicing)
	if !bytes.Equal(block.ExecutionData().StateRoot(), headRoot) {
		finalized := s.fcsp.ForkChoiceStore(ctx).GetFinalizedBlockHash()
		safe := s.fcsp.ForkChoiceStore(ctx).GetSafeBlockHash()
		_, err := s.en.NotifyForkchoiceUpdate(
			ctx, slot,
			execution.NewNotifyForkchoiceUpdateArg(
				block.ExecutionData().BlockHash(), safe[:], finalized[:],
			), true, true)
		if err != nil {
			return err
		}
	}
	return nil
}
