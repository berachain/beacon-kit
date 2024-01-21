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

package initialsync

import (
	"bytes"
	"context"
	"math/big"

	"cosmossdk.io/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/itsdevbear/bolaris/types"
)

type Status int

const (
	StatusWaiting = Status(iota) //nolint:errname // initial status of the service.
	StatusBeaconAhead
	StatusExecutionAhead
	StatusSynced
)

type ForkChoiceStoreProvider interface {
	ForkChoiceStore(ctx context.Context) types.ForkChoiceStore
}

// Service is responsible for tracking the synchornization status
// of both the beacon and execution chains.
type Service struct {
	logger    log.Logger
	ethClient ethClient
	fcsp      ForkChoiceStoreProvider
}

func NewService(opts ...Option) *Service {
	s := &Service{}
	for _, opt := range opts {
		if err := opt(s); err != nil {
			panic(err)
		}
	}
	return s
}

// Start spawns any goroutines required by the service.
func (s *Service) Start() {
	// go func() {
	// 	ticker := time.NewTicker(8 * time.Second)
	// 	defer ticker.Stop()
	// 	for {
	// 		select {
	// 		case <-ticker.C:
	// 			ctx := context.Background()
	// 			status, err := s.CheckSyncStatus(ctx)
	// 			if err != nil {
	// 				s.logger.Error("Error checking sync status", "error", err)
	// 				continue
	// 			}

	// 			switch status {
	// 			case StatusBeaconAhead:
	// 				s.logger.Info("Beacon chain is ahead of execution chain")
	// 			case StatusExecutionAhead:
	// 				s.logger.Info("Execution chain is ahead of beacon chain")
	// 			case StatusSynced:
	// 				s.logger.Info("Beacon and execution chains are synced")
	// 			}
	// 		}
	// 	}
	// }()
}

// Stop terminates all goroutines belonging to the service,
// blocking until they are all terminated.
func (s *Service) Stop() error { return nil }

// Status returns error if the service is not considered healthy.
func (s *Service) Status() error { return nil }

// CheckSyncStatus returns the current relative sync status of the beacon and execution.
func (s *Service) CheckSyncStatus(ctx context.Context) Status {
	finalHash := s.fcsp.ForkChoiceStore(ctx).GetFinalizedBlockHash()

	// Get the latest finalized block from the execution chain.
	elFinalized, err := s.ethClient.HeaderByNumber(
		ctx, big.NewInt(int64(rpc.FinalizedBlockNumber)))
	if err != nil {
		s.logger.Error("Error getting latest finalized block from execution chain", "error", err)
	}
	if elFinalized == nil {
		s.logger.Info("execution chain is waiting for a finalized block 😚")
		return StatusWaiting
	} else if bytes.Equal(elFinalized.Hash().Bytes(), finalHash[:]) {
		// If the beacon chain and the execution chain have the same finalized block,
		// then they are synced.
		s.logger.Info(
			"beacon and execution chains are synced ✅",
			"finalized_hash", common.BytesToHash(finalHash[:]),
		)
		return StatusSynced
	}

	fields := []any{
		"finalized_execution", elFinalized.Hash(),
		"finalized_execution_num", elFinalized.Number.Uint64(),
	}

	// Otherwise we need to check if the beacon chain is ahead of the execution chain.
	// Get the latest finalized block from the beacon chain.
	clFinalized, _ := s.ethClient.HeaderByHash(ctx, common.BytesToHash(finalHash[:]))
	if clFinalized == nil || clFinalized.Number.Uint64() > elFinalized.Number.Uint64() {
		// Prevent nil pointer dereference.
		if clFinalized != nil {
			fields = append([]any{
				"finalized_beacon", clFinalized.Hash(),
				"finalized_beacon_num", clFinalized.Number.Uint64(),
			}, fields...)
		}

		s.logger.Info(
			"block finalization on the beacon chain is ahead of the execution chain",
			fields...,
		)
		return StatusBeaconAhead
	}

	s.logger.Info(
		"block finalization on the execution chain is ahead of the beacon chain",
		append([]any{
			"finalized_beacon", clFinalized.Hash(),
			"finalized_beacon_num", clFinalized.Number.Uint64(),
		}, fields...)...,
	)
	return StatusExecutionAhead
}
