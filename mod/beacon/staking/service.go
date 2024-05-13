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

package staking

import (
	"context"

	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// Service represents the staking service.
type Service struct {
	logger log.Logger[any]
	bsb    BeaconStorageBackend

	// depositContract represents the deposit contract.
	depositContract DepositContract

	// deposit represents the deposit store.
	ds DepositStore
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
	return "staking"
}

func (s *Service) Start(context.Context) {}

func (s *Service) Status() error { return nil }

func (s *Service) WaitForHealthy(context.Context) {}

// ProcessLogsInETH1Block gets logs in the Eth1 block
// received from the execution client and processes them to
// convert them into appropriate objects that can be consumed
// by other services.
func (s *Service) ProcessLogsInETH1Block(
	ctx context.Context,
	blockNumber math.U64,
) error {
	deposits, err := s.depositContract.
		GetDeposits(ctx, blockNumber.Unwrap())
	if err != nil {
		return err
	}

	return s.ds.EnqueueDeposits(deposits)
}

// PruneDepositEvents prunes deposit events.
func (s *Service) PruneDepositEvents(idx uint64) error {
	return s.ds.PruneToIndex(idx)
}
