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

	"github.com/berachain/beacon-kit/mod/core/state"
	"github.com/berachain/beacon-kit/mod/execution"
	"github.com/berachain/beacon-kit/mod/node-builder/service"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/runtime/services/staking/abi"
	"github.com/berachain/beacon-kit/mod/storage/deposit"
)

// Service represents the staking service.
type Service struct {
	// BaseService is the base service.
	service.BaseService

	// ee represents the execution engine.
	ee *execution.Engine

	// abi represents the configured deposit contract's
	// abi.
	abi *abi.WrappedABI

	// deposit represents the deposit store.
	ds *deposit.KVStore
}

// ProcessLogsInETH1Block gets logs in the Eth1 block
// received from the execution client and processes them to
// convert them into appropriate objects that can be consumed
// by other services.
func (s *Service) ProcessLogsInETH1Block(
	ctx context.Context,
	st state.BeaconState,
	blockHash primitives.ExecutionHash,
) error {
	// Gather all the logs corresponding to
	// the addresses of interest from this block.
	logsInBlock, err := s.ee.GetLogs(
		ctx,
		blockHash,
		[]primitives.ExecutionAddress{
			s.ChainSpec().DepositContractAddress(),
		},
	)
	if err != nil {
		return err
	}

	return s.ProcessBlockEvents(ctx, st, logsInBlock)
}

func (s *Service) PruneDepositEvents(st state.BeaconState) error {
	idx, err := st.GetEth1DepositIndex()
	if err != nil {
		return err
	}
	s.Logger().Info("🍇 pruning deposit events", "index", idx)
	return s.ds.PruneToIndex(idx)
}
