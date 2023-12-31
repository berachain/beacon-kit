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

package eth1

import (
	"context"

	"cosmossdk.io/log"

	"github.com/ethereum/go-ethereum"
)

// ethClient is an interface that wraps the ChainSyncReader from the go-ethereum package.
type ethClient interface {
	ethereum.ChainSyncReader
}

// SyncStatus is a struct that holds the logger and the eth1client.
type SyncStatus struct {
	logger     log.Logger
	eth1client ethClient
}

// NewSyncStatus is a constructor function for SyncStatus. It takes an ethClient as an argument.
func NewSyncStatus(opts ...Option) *SyncStatus {
	ss := &SyncStatus{}
	for _, opt := range opts {
		if err := opt(ss); err != nil {
			panic(err)
		}
	}
	return ss
}

// Syncing is a method of SyncStatus that returns the sync status of the execution client.
func (s *SyncStatus) Syncing(ctx context.Context) (bool, error) {
	sp, err := s.eth1client.SyncProgress(ctx)
	if err != nil {
		return false, err
	}

	// A nil response indicates that the client is not syncing.
	if sp == nil {
		return false, nil
	}

	s.logger.Error("Syncing Execution Client",
		"currentBlock", sp.CurrentBlock,
		"highestBlock", sp.HighestBlock)

	return sp == nil, nil
}
