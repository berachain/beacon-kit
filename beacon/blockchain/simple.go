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
	"time"

	consensusv1 "github.com/itsdevbear/bolaris/types/consensus/v1"
	"github.com/itsdevbear/bolaris/types/consensus/v1/interfaces"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"

	"cosmossdk.io/core/header"
)

const payloadBuildDelay = 2

func (s *Service) BuildNextBlock(
	ctx context.Context, slot primitives.Slot, time uint64,
) (interfaces.BeaconKitBlock, error) {
	// The goal here is to build a payload whose parent is the previously
	// finalized block, such that, if this payload is accepted, it will be
	// the next finalized block in the chain.
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
