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
	"github.com/ethereum/go-ethereum/core/types"

	eth "github.com/itsdevbear/bolaris/beacon/execution/engine/ethclient"
)

// type Contract

// Processor is responsible for processing logs fr.
type Processor struct {
	eth1Client   *eth.Eth1Client
	contractAddr common.Address
}

// NewProcessor creates a new instance of Processor with the provided options.
// It applies each option to the Processor and returns an error if any of the options fail.
func NewProcessor(opts ...Option) (*Processor, error) {
	s := &Processor{}
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}
	return s, nil
}

// ProcessETH1Block processes logs from the provided eth1 block.
func (s *Processor) ProcessETH1Block(ctx context.Context, blkNum *big.Int) error {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{
			s.contractAddr,
		},
		FromBlock: blkNum,
		ToBlock:   blkNum,
	}

	logs, err := s.eth1Client.FilterLogs(ctx, query)
	if err != nil {
		return err
	}

	for i, filterLog := range logs {
		// ignore logs that are not of the required block number
		if filterLog.BlockNumber != blkNum.Uint64() {
			continue
		}

		if err = s.ProcessLog(ctx, &logs[i]); err != nil {
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

func (s *Processor) ProcessLog(ctx context.Context, log *types.Log) error {
	_ = ctx
	// if log.Address == s.depositContractAddr {
	// 	return s.processDepositLog(ctx, log)
	// }
	return nil
}
