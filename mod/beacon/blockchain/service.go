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

	"github.com/berachain/beacon-kit/mod/core"
	"github.com/berachain/beacon-kit/mod/core/state"
	datypes "github.com/berachain/beacon-kit/mod/da/pkg/types"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives"
)

// Service is the blockchain service.
type Service struct {
	// service.BaseService
	bsb    BeaconStorageBackend
	logger log.Logger[any]
	cs     primitives.ChainSpec
	ee     ExecutionEngine
	lb     LocalBuilder
	sks    StakingService
	bv     *core.BlockVerifier
	sp     *core.StateProcessor[*datypes.BlobSidecars]
	pv     *core.PayloadVerifier
}

// NewService creates a new validator service.
func NewService(
	opts ...Option,
) *Service {
	s := &Service{}
	for _, opt := range opts {
		if err := opt(s); err != nil {
			panic(err)
		}
	}

	return s
}

// Name returns the name of the service.
func (s *Service) Name() string {
	return "blockchain"
}

func (s *Service) Start(context.Context) {}

func (s *Service) Status() error { return nil }

func (s *Service) WaitForHealthy(context.Context) {}

// TODO: Remove
func (s Service) BeaconState(ctx context.Context) state.BeaconState {
	return s.bsb.BeaconState(ctx)
}

// TODO: Remove
func (s Service) ChainSpec() primitives.ChainSpec {
	return s.cs
}
