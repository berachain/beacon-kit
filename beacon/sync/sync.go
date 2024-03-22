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
	"context"

	"github.com/cometbft/cometbft/rpc/client"
)

// CheckELSync checks if the execution layer is syncing.
func (s *Service) CheckCLSync(ctx context.Context) {
	// Call the CometBFT Client to get the sync progress.
	resultStatus, err := s.clientCtx.Client.Status(ctx)
	s.isSyncedCond.L.Lock()
	defer s.isSyncedCond.L.Unlock()
	if err != nil {
		s.isCLSynced = false
		return
	}

	// If we are not catchup, then say we are synced.
	s.isCLSynced = !resultStatus.SyncInfo.CatchingUp

	// Add a log if syncing.
	if !s.isCLSynced {
		s.Logger().Warn(
			"beacon client is attemping to sync.... ",
			"current_beacon", resultStatus.SyncInfo.LatestBlockHeight,
			"highest_beacon", resultStatus.SyncInfo.CatchingUp,
			"starting_beacon", resultStatus.SyncInfo.EarliestBlockHeight,
		)
	}
}

// CheckELSync checks if the execution layer is syncing.
func (s *Service) CheckELSync(ctx context.Context) {
	// Call the ethClient to get the sync progress
	progress, err := s.engineClient.SyncProgress(ctx)
	s.isSyncedCond.L.Lock()
	defer s.isSyncedCond.L.Unlock()
	if err != nil {
		s.isELSynced = false
		return
	}

	// If progress == nil, then we say the execution client is synced.
	s.isELSynced = progress == nil

	// Add a log if syncing.
	if !s.isELSynced && progress != nil {
		s.Logger().Warn(
			"your execution client is out of sync please investigate.... ",
			"current_eth1", progress.CurrentBlock,
			"highest_eth1", progress.HighestBlock,
			"starting_eth1", progress.StartingBlock,
		)
	}
}

// UpdateNumCLPeers updates the number of peers connected at the consensus
// layer.
func (s *Service) UpdateNumCLPeers(ctx context.Context) {
	// Call the CometBFT Client to get the sync progress.
	netInfo, err := s.clientCtx.Client.(client.NetworkClient).NetInfo(ctx)
	if err != nil {
		s.clNumPeers = 0
		return
	}

	//#nosec:G701 // if our number of peers overflows a int64
	// we have bigger problems.
	s.clNumPeers = uint64(netInfo.NPeers)
}

// UpdateNumELPeers updates the number of peers connected at the execution
// layer.
func (s *Service) UpdateNumELPeers(_ context.Context) {
	// TODO: Net is not avail over the 8551 port?
	// // Call the ethClient to get the sync progress
	// numPeers, err := s.engineClient.PeerCount(ctx)
	// if err != nil {
	// 	s.elNumPeers = 0
	// 	return
	// }
	s.elNumPeers = 0
}
