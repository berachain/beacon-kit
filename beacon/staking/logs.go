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

package staking

import (
	"context"
	"math/big"

	"cosmossdk.io/errors"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	coretypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"

	stakingabi "github.com/itsdevbear/bolaris/beacon/staking/abi"
)

const (
	depositEventName    = "Delegate"
	withdrawalEventName = "Undelegate"
)

// ProcessFinalizedETH1Block processes logs from an eth1 block, but before doing so
// it checks if the block is safe to process.
func (s *Service) ProcessFinalizedETH1Block(ctx context.Context, blkNum *big.Int) error {
	// Get the safe block number from the eth1 client.
	// TODO do we want to come up with a heuristic around when we check the execution client,
	// vs when we check the forkchoice store.
	finalBlock, err := s.eth1Client.BlockByNumber(ctx, big.NewInt(int64(rpc.FinalizedBlockNumber)))
	if err != nil {
		return err
	}

	// Ensure we don't start processing the logs of a block that is ahead of the safe block.
	if finalBlock.Number().Cmp(blkNum) < 0 {
		return errors.Wrapf(
			ErrProcessingUnfinalizedBlock, "safe block %d is behind block %d", finalBlock.Number(), blkNum,
		)
	}

	return s.ProcessETH1Block(ctx, finalBlock.Number())
}

// GatherLogsFromEth1Block gathers all the logs from the provided eth1 block.
func (s *Service) GatherLogsFromEth1Block(
	ctx context.Context, blkNum *big.Int,
) ([]coretypes.Log, error) {
	// Create a filter query for the block, to acquire all logs from contracts
	// that we care about.
	query := ethereum.FilterQuery{
		Addresses: []common.Address{s.depositContractAddress},
		FromBlock: blkNum,
		ToBlock:   blkNum,
	}

	// Gather all the logs from this block.
	return s.eth1Client.FilterLogs(ctx, query)
}

// ProcessETH1Block processes logs from the provided eth1 block.
func (s *Service) ProcessETH1Block(ctx context.Context, blkNum *big.Int) error {
	// Gather all the logs from this block.
	logs, err := s.GatherLogsFromEth1Block(ctx, blkNum)
	if err != nil {
		return err
	}

	stakingAbi, err := stakingabi.StakingMetaData.GetAbi()
	if err != nil {
		return errors.Wrap(err, "could not get staking ABI")
	}
	depositEvent := stakingAbi.Events[depositEventName]
	withdrawalEvent := stakingAbi.Events[withdrawalEventName]

	// Process each log.
	for _, filterLog := range logs {
		// Skip logs that are not from the block we are processing.
		// This should never happen, but defensively check anyway.
		if filterLog.BlockNumber != blkNum.Uint64() {
			continue
		}

		// Skip logs that are not from the addresses we care about.
		// This should never happen, but defensively check anyway.
		if filterLog.Address != s.depositContractAddress {
			continue
		}

		var unpackedArgs []any
		switch filterLog.Topics[0] {
		case depositEvent.ID:
			// Process the deposit log.
			unpackedArgs, err = depositEvent.Inputs.Unpack(filterLog.Data)
			if err != nil {
				return errors.Wrap(err, "could not unpack deposit log")
			}
			err = s.ProcessDeposit(ctx, unpackedArgs)
			if err != nil {
				return errors.Wrap(err, "could not process deposit")
			}
		case withdrawalEvent.ID:
			// Process the withdrawal log.
			unpackedArgs, err = withdrawalEvent.Inputs.Unpack(filterLog.Data)
			if err != nil {
				return errors.Wrap(err, "could not unpack withdrawal log")
			}
			err = s.ProcessWithdrawal(ctx, unpackedArgs)
			if err != nil {
				return errors.Wrap(err, "could not process withdrawal")
			}
		default:
			// Skip logs that are not from the events we care about.
			continue
		}
	}

	// if !s.chainStartData.Chainstarted {
	// 	if err := s.processChainStartFromBlockNum(ctx, blkNum); err != nil {
	// 		return err
	// 	}
	// }
	return nil
}
