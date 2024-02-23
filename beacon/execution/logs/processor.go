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
	"math/big"

	"cosmossdk.io/errors"
	"cosmossdk.io/log"
	"github.com/ethereum/go-ethereum/common"
	coretypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

//nolint:gochecknoglobals // no sense reallocating over and over.
var rpcFinalizedBlockQuery = big.NewInt(int64(rpc.FinalizedBlockNumber))

// Processor is responsible for processing logs fr.
type Processor struct {
	logger   log.Logger
	handlers map[common.Address]Handler
}

// NewProcessor creates a new instance of Processor with the provided options.
// It applies each option to the Processor and returns an error if any of the
// options fail.
func NewProcessor(opts ...Option) (*Processor, error) {
	s := &Processor{
		handlers: make(map[common.Address]Handler),
	}
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}
	return s, nil
}

// ProcessLogs processes logs from the provided Eth1 block.
func (s *Processor) ProcessLogs(
	ctx context.Context,
	logs []coretypes.Log,
	blkNum uint64,
) error {
	// Process each log.
	for i, log := range logs {
		// Skip logs that are not from the block we are processing.
		// This should never happen, but defensively check anyway.
		if log.BlockNumber != blkNum {
			continue
		}

		// Skip logs that are not from the addresses we care about.
		// This should never happen, but defensively check anyway.
		handler, found := s.handlers[log.Address]
		if !found {
			continue
		}

		// Process the log with the handler.
		if err := handler.HandleLog(ctx, &logs[i]); err != nil {
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
