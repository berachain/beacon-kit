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

	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives"
)

// Service is the blockchain service.
type Service[
	BeaconStateT BeaconState[BeaconStateT],
	BlobSidecarsT BlobSidecars,
	DepositStoreT DepositStore,
] struct {
	// bsb represents the backend storage for beacon states and associated
	// sidecars.
	bsb BeaconStorageBackend[
		BeaconStateT, BlobSidecarsT, DepositStoreT,
	]

	// logger is used for logging messages in the service.
	logger log.Logger[any]

	// cs holds the chain specifications.
	cs primitives.ChainSpec

	// ee is the execution engine responsible for processing execution payloads.
	ee ExecutionEngine

	// bdc is a connection to the deposit contract.
	bdc DepositContract

	// lb is a local builder for constructing new beacon states.
	lb LocalBuilder[BeaconStateT]

	// bv is responsible for verifying beacon blocks.
	bv BlockVerifier[BeaconStateT]

	// sp is the state processor for beacon blocks and states.
	sp StateProcessor[BeaconStateT, BlobSidecarsT]

	// pv verifies the payload of beacon blocks.
	pv PayloadVerifier[BeaconStateT]
}

// NewService creates a new validator service.
func NewService[
	BeaconStateT BeaconState[BeaconStateT],
	BlobSidecarsT BlobSidecars,
	DepositStoreT DepositStore,
](
	bsb BeaconStorageBackend[
		BeaconStateT, BlobSidecarsT, DepositStoreT],
	logger log.Logger[any],
	cs primitives.ChainSpec,
	ee ExecutionEngine,
	lb LocalBuilder[BeaconStateT],
	bv BlockVerifier[BeaconStateT],
	sp StateProcessor[BeaconStateT, BlobSidecarsT],
	pv PayloadVerifier[BeaconStateT],
	bdc DepositContract,
) *Service[BeaconStateT, BlobSidecarsT, DepositStoreT] {
	return &Service[BeaconStateT, BlobSidecarsT, DepositStoreT]{
		bsb:    bsb,
		logger: logger,
		cs:     cs,
		ee:     ee,
		lb:     lb,
		bv:     bv,
		sp:     sp,
		pv:     pv,
		bdc:    bdc,
	}
}

// Name returns the name of the service.
func (s *Service[
	BeaconStateT, BlobSidecarsT, DepositStoreT,
]) Name() string {
	return "blockchain"
}

func (s *Service[
	BeaconStateT, BlobSidecarsT, DepositStoreT,
]) Start(
	context.Context,
) {
}

func (s *Service[
	BeaconStateT, BlobSidecarsT, DepositStoreT,
]) Status() error {
	return nil
}

func (s *Service[
	BeaconStateT, BlobSidecarsT, DepositStoreT,
]) WaitForHealthy(
	context.Context,
) {
}

// TODO: Remove
func (s Service[
	BeaconStateT, BlobSidecarsT, DepositStoreT,
]) StateFromContext(
	ctx context.Context,
) BeaconStateT {
	return s.bsb.StateFromContext(ctx)
}
