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

package execution

import (
	"context"

	"github.com/berachain/beacon-kit/engine"
	"github.com/berachain/beacon-kit/primitives"
	"github.com/berachain/beacon-kit/runtime/service"
)

// Service is responsible for delivering beacon chain notifications to
// the execution client and processing logs received from the execution client.
type Service struct {
	service.BaseService
	// engine gives the notifier access to the engine api of the execution
	// client.
	engine *engine.ExecutionEngine
	sks    StakingService
}

// Start spawns any goroutines required by the service.
func (s *Service) Start(ctx context.Context) {
	go s.engine.Start(ctx)
}

// Status returns error if the service is not considered healthy.
func (s *Service) Status() error {
	return s.engine.Status()
}

// ProcessLogsInETH1Block gets logs in the Eth1 block
// received from the execution client and processes them to
// convert them into appropriate objects that can be consumed
// by other services.
func (s *Service) ProcessLogsInETH1Block(
	ctx context.Context,
	blockHash primitives.ExecutionHash,
) error {
	// Gather all the logs corresponding to
	// the addresses of interest from this block.
	logsInBlock, err := s.engine.GetLogs(
		ctx,
		blockHash,
		[]primitives.ExecutionAddress{
			s.BeaconCfg().DepositContractAddress,
		},
	)
	if err != nil {
		return err
	}

	return s.sks.ProcessBlockEvents(ctx, logsInBlock)
}
