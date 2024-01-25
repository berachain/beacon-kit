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
	StatusUnknown = Status(-1)   //nolint:errname // initial status of the service.
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
func (s *Service) Start() {}

// Stop terminates all goroutines belonging to the service,
// blocking until they are all terminated.
func (s *Service) Stop() error { return nil }

// Status returns error if the service is not considered healthy.
func (s *Service) Status() error { return nil }

// CheckSyncStatus returns the current synchronization status of the beacon and execution chains.
//
// TODO, We need to add a handler than does the following after this function returns
// `StatusBeaconAhead`.
// 1. Fire off event to the dispatcher to trigger a fork choice
// 2. Block here until it is sync'd.
// 3. Return we are blessed.
func (s *Service) CheckSyncStatus(ctx context.Context) Status {
	// First lets grab the beacon chains view of the last finalized execution layer block.
	finalHash := s.fcsp.ForkChoiceStore(ctx).GetFinalizedBlockHash()

	// If the chain hasn't been started yet, we are at genesis, and we can't really do anything.
	// This is to handle calling this function before InitGenesis has been called. If InitGenesis
	// has previously been called, we will continue on. We return StatuSynced here even if it is
	// not totally true. This is because we don't want to block the beacon chain from
	// starting up.
	isBeaconGenesis := bytes.Equal(finalHash[:], common.Hash{}.Bytes())
	if isBeaconGenesis {
		return StatusSynced
	}

	// The only other thing we can do before ABCI starts is to handle the case where the beacon
	// chain is AHEAD of the execution chain. We can't check the converse, since we don't know
	// what blocks we are missing, so there at this point in time, we cannot tell the execution
	// chain where to jump to anyways.

	// We previously grabbed the beacon chain's view of what is finalized. We first ensure it
	// exists. If it exists on the chain, this is bullish. If it doesn't we need to forkchoice.
	clFinalized, _ := s.ethClient.HeaderByHash(ctx, common.BytesToHash(finalHash[:]))
	if clFinalized == nil {
		// We need to fork choice to find the latest finalized block. This is trigger the execution
		// chain to start asking it's peers to help it sync and build the chain required for
		// the following forkchoice.
		return StatusBeaconAhead
	}

	// If clFinalized != nil, then we know that the beacon chain is at or behind the execution chain.
	// So let's figure out whats going on by getting the last block that the execution chain believes
	// is finalized.
	elFinalized, err := s.ethClient.HeaderByNumber(
		ctx, big.NewInt(int64(rpc.FinalizedBlockNumber)))
	if err != nil || elFinalized == nil {
		// If the execution chain doesn't have a finalized block, then we are in a weird state, and
		// your client is like kinda fucked up rn, because if this is the case then the above
		// clFinalized call should've failed and a forkchoice should've been triggered.
		s.logger.Error("Error getting latest finalized block from execution chain", "error", err)
		return StatusUnknown
	}

	// Once we reach here, we can confirm that the consensus layer and the execution
	// layer have their own view of the world, and we now need to configure whether or not these
	// views align. We will define "things being in sync" when the latest finalized beacon chain
	// block, is either equal to the execution chain block, or AT MOST 1 block ahead. This 1 block
	// ahead provision is due to the one block delay in finalization.
	clBlockNum := clFinalized.Number
	elBlockNum := elFinalized.Number

	// Check if the beacon chain block is either equal to the execution chain block or at most
	// 1 block ahead.
	if clBlockNum.Cmp(elBlockNum) == 0 || clBlockNum.Cmp(
		new(big.Int).Add(elBlockNum, big.NewInt(1)),
	) == 0 {
		// The beacon chain and the execution chain are at the same number || The beacon chain is at
		// most 1 block ahead of the execution chain.
		s.logger.Info(
			"beacon and execution chains are synced ✅",
			"finalized_hash", common.BytesToHash(finalHash[:]),
		)
		return StatusSynced
	} else if clBlockNum.Cmp(elBlockNum) > 0 {
		// The beacon chain is ahead of the execution chain.
		return StatusBeaconAhead
	}

	// By ruling out everything else, we can say the execution chain is ahead of the beacon chain.
	// There is nothing really actionable to do here, as we need to just let the beacon chain
	// keep syncing, until it passes the execution chain head. Only then can we issue a forkchoice
	// update to start syncing the execution chain again.
	return StatusExecutionAhead
}

func (*Service) RequestForkchoiceUpdate() error {
	// TODO: this function needs to queue up a request for the execution service to
	// notify the execution layer a forkchoice is required.
	// TODO: wire up dispatcher.
	return nil
}
