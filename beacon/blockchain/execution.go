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

	"github.com/itsdevbear/bolaris/beacon/execution"
	consensusv1 "github.com/itsdevbear/bolaris/types/consensus/v1"
	"github.com/itsdevbear/bolaris/types/consensus/v1/interfaces"
	prsymexecution "github.com/prysmaticlabs/prysm/v4/beacon-chain/execution"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"

	"cosmossdk.io/core/header"
)

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
	_, err := s.en.NotifyForkchoiceUpdate(
		ctx, slot,
		execution.NewNotifyForkchoiceUpdateArg(
			headHash, safeHash[:], finalHash[:],
		),
		true,
		true,
	)

	if err != nil {
		return nil, err
	}

	// todo we need to wait for the forkchoice to update?
	time.Sleep(payloadBuildDelay * time.Second)

	payload, _, _, err := s.en.GetBuiltPayload(
		ctx, slot,
	)
	return payload, err
}

// ValidateProposedBeaconBlock validates a proposed beacon block.
func (s *Service) ValidateProposedBeaconBlock(ctx context.Context,
	block header.Info, header interfaces.ExecutionData,
) (*enginev1.PayloadIDBytes, error) {
	// We must first notify the execution client we have received a new payload. We ask
	// the execution client to try to insert this payload to check it's validity.
	isValidPayload, err := s.en.NotifyNewPayload(ctx, 0, header /*, nil, [32]byte{}*/)
	if err != nil {
		if !errors.Is(err, prsymexecution.ErrAcceptedSyncingPayloadStatus) {
			s.logger.Error("failed to validate execution on block", "err", err)
			return nil, err
		}
	} else if !isValidPayload {
		return nil, prsymexecution.ErrInvalidPayloadStatus
	}

	// Forkchoice our execution client's head to be the block that we validated as correct
	// above. We also lazily update our finalized and safe block hashes to be the same as
	// what is currently on the beacon chain.
	finalized := s.fcsp.ForkChoiceStore(ctx).GetFinalizedBlockHash()
	safe := s.fcsp.ForkChoiceStore(ctx).GetSafeBlockHash()
	return s.en.NotifyForkchoiceUpdate(
		ctx, primitives.Slot(block.Height),
		execution.NewNotifyForkchoiceUpdateArg(
			header.BlockHash(), safe[:], finalized[:],
		), true, true)
}
