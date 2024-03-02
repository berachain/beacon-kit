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

package logs

import (
	"context"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/itsdevbear/bolaris/beacon/core/state"
)

// LogRequest is a request for logs sent from a service.
type LogRequest struct {
	ContractAddress ethcommon.Address
	Allocator       *TypeAllocator
}

// LogCache is the interface for the cache of logs.
type LogCache interface {
	// ShouldProcess returns true if the cache
	// determines that the log should be processed.
	ShouldProcess(log *ethtypes.Log) bool
	// Push pushes the log value container into the cache.
	Push(container LogValueContainer) error
	// LastFinalizedBlock returns the block number of
	// the last finalized block in cache.
	LastFinalizedBlock() uint64
	// SetLastFinalizedBlock sets the block number of
	// the last finalized block in cache.
	SetLastFinalizedBlock(blockNumber uint64)
}

// LogValueContainer is the interface for the container of
// the unmarsheled value of a log and other related information.
type LogValueContainer interface {
	// BlockHash returns the block hash of the log.
	BlockHash() ethcommon.Hash
	// BlockNumber returns the block number of the log.
	BlockNumber() uint64
	// LogIndex returns the index of the log in the block.
	LogIndex() uint
}

type LogFactory interface {
	GetRegisteredAddresses() []ethcommon.Address
	ProcessLog(log *ethtypes.Log) (LogValueContainer, error)
}

// ReadOnlyForkChoiceProvider provides the read-only fork choice store.
type ReadOnlyForkChoiceProvider interface {
	ForkchoiceStore(ctx context.Context) state.ReadOnlyForkChoice
}
