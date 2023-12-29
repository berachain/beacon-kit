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
	"context"
	"errors"
	"time"

	consensusv1 "github.com/itsdevbear/bolaris/types/consensus/v1"
	"github.com/itsdevbear/bolaris/types/consensus/v1/interfaces"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/execution"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"

	"cosmossdk.io/core/header"
)

const payloadBuildDelay = 2

func (s *Service) BuildNextBlock(
	ctx context.Context, slot primitives.Slot, time uint64,
) (interfaces.BeaconKitBlock, error) {
	// The goal here is to build a payload whose parent is the previously
	// finalized block, such that, if this payload is accepted, it will be
	// the next finalized block in the chain. A byproduct of this design
	// is that we get the nice property of lazily propogate the finalized
	// and safe block hashes to the execution client.
	lastFinalizedBlock := s.fcsp.ForkChoiceStore(ctx).GetFinalizedBlockHash()
	executionData, err := s.buildNewBlockOnTopOf(ctx, slot, lastFinalizedBlock[:])
	if err != nil {
		return nil, err
	}

	// Create a new block with the payload.
	return consensusv1.NewBaseBeaconKitBlock(
		slot, time, executionData,
		s.beaconCfg.ActiveForkVersion(primitives.Epoch(slot)),
	)
}

// buildNewBlockOnTopOf builds a new block on top of an existing head of the execution client.
func (s *Service) buildNewBlockOnTopOf(ctx context.Context,
	slot primitives.Slot, headHash []byte) (interfaces.ExecutionData, error) {
	finalHash := s.fcsp.ForkChoiceStore(ctx).GetFinalizedBlockHash()
	safeHash := s.fcsp.ForkChoiceStore(ctx).GetSafeBlockHash()
	payloadIDNew, err := s.notifyForkchoiceUpdateWithSyncingRetry(
		ctx, slot,
		&notifyForkchoiceUpdateArg{
			headHash:  headHash,
			safeHash:  safeHash[:],
			finalHash: finalHash[:],
		},
		true,
	)

	if err != nil {
		return nil, err
	}

	// todo we need to wait for the forkchoice to update?
	time.Sleep(payloadBuildDelay * time.Second)

	payload, _, _, err := s.engine.GetPayload(
		ctx, [8]byte(payloadIDNew[:]), slot,
	)
	return payload, err
}

// ValidateProposedBeaconBlock validates a proposed beacon block.
func (s *Service) ValidateProposedBeaconBlock(ctx context.Context,
	block header.Info, header interfaces.ExecutionData,
) (*enginev1.PayloadIDBytes, error) {
	isValidPayload, err := s.notifyNewPayload(ctx, 0, header /*, nil, [32]byte{}*/)
	if err != nil {
		if !errors.Is(err, execution.ErrAcceptedSyncingPayloadStatus) {
			s.logger.Error("failed to validate execution on block", "err", err)
			return nil, err
		}
	} else if !isValidPayload {
		return nil, execution.ErrInvalidPayloadStatus
	}

	// Forkchoice our execution client's head to be the block that we validated as correct
	// above. NOTE: we do not touch our safe or finalized hashes here.
	finalized := s.fcsp.ForkChoiceStore(ctx).GetFinalizedBlockHash()
	safe := s.fcsp.ForkChoiceStore(ctx).GetSafeBlockHash()
	return s.notifyForkchoiceUpdateWithSyncingRetry(
		ctx, primitives.Slot(block.Height), &notifyForkchoiceUpdateArg{
			headHash:  header.BlockHash(),
			safeHash:  safe[:],
			finalHash: finalized[:],
		}, true)
}

func (s *Service) FinalizeBlockAsync(
	_ context.Context, beaconBlock header.Info, toFinalize []byte,
) error {
	s.finalizer.RequestFinalization(toFinalize, beaconBlock)
	return nil
}

// FinalizeBlock marks the block as finalized on the execution layer.
func (s *Service) FinalizeBlock(
	ctx context.Context, slot uint64, toFinalize []byte,
) error {
	_, err := s.notifyForkchoiceUpdateWithSyncingRetry(
		ctx, primitives.Slot(slot),
		&notifyForkchoiceUpdateArg{
			headHash:  toFinalize,
			safeHash:  toFinalize,
			finalHash: toFinalize,
		},
		false, // todo: maybe we can store a cache of payloadIDs and
		// beginb uilding a new payoad here
		// OR we abstract away payload building into its own async thingy.
	)
	return err
}
