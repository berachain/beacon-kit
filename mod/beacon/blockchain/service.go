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
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core/state"
)

// Service is the blockchain service.
type Service[
	BeaconStateT state.BeaconState,
	BlobSidecarsT BlobSidecars,
] struct {
	// bsb represents the backend storage for beacon states and associated
	// sidecars.
	bsb BeaconStorageBackend[
		BeaconStateT, BlobSidecarsT,
	]

	// logger is used for logging messages in the service.
	logger log.Logger[any]

	// cs holds the chain specifications.
	cs primitives.ChainSpec

	// ee is the execution engine responsible for processing execution payloads.
	ee ExecutionEngine

	// lb is a local builder for constructing new beacon states.
	lb LocalBuilder[BeaconStateT]

	// sks is the staking service managing staking logic.
	sks StakingService

	// bv is responsible for verifying beacon blocks.
	bv BlockVerifier[BeaconStateT]

	// sp is the state processor for beacon blocks and states.
	sp *core.StateProcessor[types.BeaconBlock, BeaconStateT, BlobSidecarsT]

	// pv verifies the payload of beacon blocks.
	pv PayloadVerifier[BeaconStateT]
}

// NewService creates a new validator service.
func NewService[
	BeaconStateT state.BeaconState, BlobSidecarsT BlobSidecars,
](
	bsb BeaconStorageBackend[BeaconStateT, BlobSidecarsT],
	logger log.Logger[any],
	cs primitives.ChainSpec,
	ee ExecutionEngine,
	lb LocalBuilder[BeaconStateT],
	sks StakingService,
	bv BlockVerifier[BeaconStateT],
	sp *core.StateProcessor[
		types.BeaconBlock,
		BeaconStateT,
		BlobSidecarsT,
	],
	pv PayloadVerifier[BeaconStateT],
) *Service[BeaconStateT, BlobSidecarsT] {
	return &Service[BeaconStateT, BlobSidecarsT]{
		bsb:    bsb,
		logger: logger,
		cs:     cs,
		ee:     ee,
		lb:     lb,
		sks:    sks,
		bv:     bv,
		sp:     sp,
		pv:     pv,
	}
}

// Name returns the name of the service.
func (s *Service[BeaconStateT, BlobSidecarsT]) Name() string {
	return "blockchain"
}

func (s *Service[BeaconStateT, BlobSidecarsT]) Start(context.Context) {}

func (s *Service[BeaconStateT, BlobSidecarsT]) Status() error { return nil }

func (s *Service[BeaconStateT, BlobSidecarsT]) WaitForHealthy(
	context.Context,
) {
}

// TODO: Remove
func (s Service[BeaconStateT, BlobSidecarsT]) BeaconState(
	ctx context.Context,
) BeaconStateT {
	return s.bsb.BeaconState(ctx)
}

// TODO: Remove
func (s Service[BeaconStateT, BlobSidecarsT]) ChainSpec() primitives.ChainSpec {
	return s.cs
}
