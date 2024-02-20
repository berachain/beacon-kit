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

package sync

import (
	"bytes"
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/itsdevbear/bolaris/beacon/execution"
	"github.com/itsdevbear/bolaris/runtime/service"
)

// Service is responsible for tracking the synchornization status
// of both the beacon and execution chains.
type Service struct {
	service.BaseService
	ethClient ethClient
	es        executionService
}

// CheckSyncStatus returns the current synchronization status of the beacon and
// execution chains.
//
// TODO, We need to add a handler than does the following after this function
// returns `StatusBeaconAhead`.
// 1. Fire off event to the dispatcher to trigger a fork choice
// 2. Block here until it is sync'd.
// 3. Return we are blessed.
func (s *Service) CheckSyncStatus(ctx context.Context) *BeaconSyncProgress {
	// First lets grab the beacon chains view of the last finalized execution
	// layer block.
	finalHash := s.BeaconState(ctx).GetFinalizedEth1BlockHash()

	// If the chain hasn't been started met, we are at genesis, and we can't
	// really do anything. This is to handle calling this function before
	// InitGenesis has been called. If InitGenesis has previously been called,
	// we will continue on. We return StatuSynced here even if it is not totally
	// true. This is because we don't want to block the beacon chain from
	// starting up.
	isBeaconGenesis := bytes.Equal(finalHash[:], common.Hash{}.Bytes())
	if isBeaconGenesis {
		return &BeaconSyncProgress{Status: StatusSynced}
	}

	// The only other thing we can do before ABCI starts is to handle the case
	// where the beacon chain is AHEAD of the execution chain. We can't check
	// the converse, since we don't know what blocks we are missing, so there at
	// this point in time, we cannot tell the execution chain where to jump to
	// anyways.

	// We previously grabbed the beacon chain's view of what is finalized. We
	// first ensure it exists. If it exists on the chain, this is bullish. If it
	// doesn't we need to forkchoice.
	clFinalized, err := s.ethClient.HeaderByHash(ctx,
		common.BytesToHash(finalHash[:]))
	if clFinalized == nil || err != nil {
		// We need to fork choice to find the latest finalized block. This is
		// trigger the execution chain to start asking it's peers to help it sync
		// and build the chain required for the following forkchoice.
		return &BeaconSyncProgress{Status: StatusBeaconAhead, clFinalized: finalHash}
	}

	// If clFinalized != nil, then we know that the beacon chain is at or behind
	// the execution chain. So let's figure out whats going on by getting the
	// last block that the execution chain believes is finalized.
	elFinalized, err := s.ethClient.HeaderByNumber(
		ctx, big.NewInt(int64(rpc.FinalizedBlockNumber)))
	if err != nil || elFinalized == nil {
		// If we have something confirmed as finalized on the beacon chain, but
		// we don't have anything confirmed as finalized on the execution chain,
		// we need to forkchoice to find the latest finalized block.
		return &BeaconSyncProgress{Status: StatusBeaconAhead, clFinalized: finalHash}
	}

	// Once we reach here, we can confirm that the consensus layer and the
	// execution layer have their own view of the world, and we now need to
	// configure whether or not these views align. We will define "things being
	// in sync" when the latest finalized beacon chain block, is either equal to
	// the execution chain block, or AT MOST 1 block ahead. This 1 block ahead
	// provision is due to the one block delay in finalization.
	clBlockNum := clFinalized.Number
	elBlockNum := elFinalized.Number

	// Check if the beacon chain block is either equal to the execution chain
	// block or at most 1 block ahead.
	if clBlockNum.Cmp(elBlockNum) == 0 || clBlockNum.Cmp(
		new(big.Int).Add(elBlockNum, big.NewInt(1)),
	) == 0 {
		// The beacon chain and the execution chain are at the same number || The
		// beacon chain is at most 1 block ahead of the execution chain.
		s.Logger().Info(
			"beacon and execution chains are synced ✅",
			"finalized_hash", common.BytesToHash(finalHash[:]),
		)
		return &BeaconSyncProgress{
			Status:      StatusSynced,
			elFinalized: elFinalized.Hash(),
			clFinalized: finalHash,
		}
	} else if clBlockNum.Cmp(elBlockNum) > 0 {
		// The beacon chain is ahead of the execution chain.
		return &BeaconSyncProgress{
			Status:      StatusBeaconAhead,
			elFinalized: elFinalized.Hash(),
			clFinalized: finalHash,
		}
	}

	// By ruling out everything else, we can say the execution chain is ahead of
	// the beacon chain. There is nothing really actionable to do here, as we
	// need to just let the beacon chain keep syncing until it passes the
	// execution chain head. Only then can we issue a forkchoice update to start
	// syncing the execution chain again.
	return &BeaconSyncProgress{
		Status:      StatusExecutionAhead,
		elFinalized: elFinalized.Hash(),
		clFinalized: finalHash,
	}
}

// WaitForExecutionClientSync is a helper function that calls CheckSyncStatus
// and then triggers a forkchoice update if the beacon chain is ahead of the
// execution chain.
func (s *Service) WaitForExecutionClientSync(ctx context.Context) error {
	// First start by checking the sync status.
	bss := s.CheckSyncStatus(ctx)
	clFinalizedHash := bss.clFinalized
	rpcFinalizedBlockNumber := big.NewInt(int64(rpc.FinalizedBlockNumber))

	s.Logger().Info("verifying execution client is synced with consensus client 🔎")

	// If the beacon chain is ahead of the execution chain, we need to trigger a
	// forkchoice update to get the execution chain to start syncing, otherwise we
	// can just return.
	if bss.Status == StatusExecutionAhead {
		s.Logger().Info(
			"skipping startup forkchoice update, execution client is ahead",
			"status", bss.Status,
		)
		return nil
	}

	// Only forkchoice update if the beacon chain has a valid finalized block.
	// This is a hack at genesis to prevent the forkchoice update from being
	// called before InitGenesis, this should be fixed later on.
	if bytes.Equal(bss.clFinalized.Bytes(), (common.Hash{}).Bytes()) {
		return nil
	}
	// We will spin here forever until the consensus client and the
	// execution client, have the same view of the world. This is
	// important as we don't want to start the beacon chain until
	// the execution chain is ready to start processing blocks.
	for {
		// We forkchoice here to trigger the execution client to start syncing.
		if _, err := s.es.NotifyForkchoiceUpdate(
			ctx,
			&execution.FCUConfig{
				HeadEth1Hash: clFinalizedHash,
			},
		); err != nil {
			s.Logger().Warn(
				"failed to send forkchoice update",
				"error", err,
			)
		}

		// Update the elLatestFinalizedHeader to the latest block to update the for loop
		// exit variable.
		elLatestFinalizedHeader, err := s.ethClient.HeaderByNumber(
			ctx, rpcFinalizedBlockNumber,
		)
		if err != nil {
			return err
		}

		s.Logger().Info(
			"waiting for execution client to sync with consensus client 🔎",
			"cl_finalized", clFinalizedHash.Hex()[:8]+"...",
			"el_finalized", elLatestFinalizedHeader.Hash().Hex()[:8]+"...",
		)

		// If the execution client has caught up to the consensus client, we can
		// break out of the loop.
		if clFinalizedHash == elLatestFinalizedHeader.Hash() {
			break
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(3 * time.Second): //nolint:gomnd // 3 seconds is fine.
			continue
		}
	}

	s.Logger().Info("execution client is synced with consensus client 🎉")
	return nil
}
