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

package logs

import (
	"context"
	"math/big"

	"cosmossdk.io/errors"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"

	eth "github.com/itsdevbear/bolaris/beacon/execution/engine/ethclient"
	"github.com/itsdevbear/bolaris/beacon/execution/logs/callback"
)

// type Contract

// Processor is responsible for processing logs fr.
type Processor struct {
	eth1Client *eth.Eth1Client
	handlers   map[common.Address]callback.LogHandler
}

// NewProcessor creates a new instance of Processor with the provided options.
// It applies each option to the Processor and returns an error if any of the options fail.
func NewProcessor(opts ...Option) (*Processor, error) {
	s := &Processor{
		handlers: make(map[common.Address]callback.LogHandler),
	}
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}
	return s, nil
}

// ProcessETH1Block processes logs from the provided eth1 block.
func (s *Processor) ProcessETH1Block(ctx context.Context, blkNum *big.Int) error {
	// TODO: add safety check here? Check if Safe/Finalized before processing?

	// Gather all the addresses we have handlers for.
	addresses := make([]common.Address, 0)
	for addr := range s.handlers {
		addresses = append(addresses, addr)
	}

	// Create a filter query for the block, to acquire all logs from contracts
	// that we care about.
	query := ethereum.FilterQuery{
		Addresses: addresses,
		FromBlock: blkNum,
		ToBlock:   blkNum,
	}

	// Gather all the logs from this block.
	logs, err := s.eth1Client.FilterLogs(ctx, query)
	if err != nil {
		return err
	}

	// Process each log.
	for i, filterLog := range logs {
		// Skip logs that are not from the block we are processing.
		// This should never happen, but defensively check anyway.
		if filterLog.BlockNumber != blkNum.Uint64() {
			continue
		}

		// Skip logs that are not from the addresses we care about.
		// This should never happen, but defensively check anyway.
		handler, found := s.handlers[filterLog.Address]
		if !found {
			continue
		}

		// Process the log with the handler.
		if err = handler.HandleLog(ctx, logs[i]); err != nil {
			return errors.Wrap(err, "could not process log")
		}
	}

	// if !s.chainStartData.Chainstarted {
	// 	if err := s.processChainStartFromBlockNum(ctx, blkNum); err != nil {
	// 		return err
	// 	}
	// }
	return nil
}
