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
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core/state"
)

// Service is the blockchain service.
type Service[
	BlobSidecarsT BlobSidecars,
] struct {
	// service.BaseService
	bsb    BeaconStorageBackend[BlobSidecarsT]
	logger log.Logger[any]
	cs     primitives.ChainSpec
	ee     ExecutionEngine
	lb     LocalBuilder
	sks    StakingService
	bv     BlockVerifier
	sp     *core.StateProcessor[BlobSidecarsT]
	pv     PayloadVerifier
}

// NewService creates a new validator service.
func NewService[BlobSidecarsT BlobSidecars](
	bsb BeaconStorageBackend[BlobSidecarsT],
	logger log.Logger[any],
	cs primitives.ChainSpec,
	ee ExecutionEngine,
	lb LocalBuilder,
	sks StakingService,
	bv BlockVerifier,
	sp *core.StateProcessor[BlobSidecarsT],
	pv PayloadVerifier,
) *Service[BlobSidecarsT] {
	return &Service[BlobSidecarsT]{
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
func (s *Service[BlobSidecarsT]) Name() string {
	return "blockchain"
}

func (s *Service[BlobSidecarsT]) Start(context.Context) {}

func (s *Service[BlobSidecarsT]) Status() error { return nil }

func (s *Service[BlobSidecarsT]) WaitForHealthy(context.Context) {}

// TODO: Remove
func (s Service[BlobSidecarsT]) BeaconState(
	ctx context.Context,
) state.BeaconState {
	return s.bsb.BeaconState(ctx)
}

// TODO: Remove
func (s Service[BlobSidecarsT]) ChainSpec() primitives.ChainSpec {
	return s.cs
}
