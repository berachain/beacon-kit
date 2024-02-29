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

package builder

import (
	"context"
	"errors"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	stakinglogs "github.com/itsdevbear/bolaris/beacon/staking/logs"
	"github.com/itsdevbear/bolaris/config"
	"github.com/itsdevbear/bolaris/runtime/service"
	"github.com/itsdevbear/bolaris/types/consensus"
	"github.com/itsdevbear/bolaris/types/consensus/primitives"
	consensusv1 "github.com/itsdevbear/bolaris/types/consensus/v1"
	"github.com/sourcegraph/conc/pool"
)

// Service is responsible for building beacon blocks.
type Service struct {
	service.BaseService
	cfg *config.Builder

	// es is the execution service that is responsible for
	// processing logs in the eth1 block.
	es ExecutionService

	// ss is the staking service that is responsible for
	// accepting deposits into the deposit queue.
	ss StakingService

	// localBuilder represents the local block builder, this builder
	// is connected to this nodes execution client via the EngineAPI.
	// Building blocks is done by submitting forkchoice updates through.
	// The local Builder.
	localBuilder   PayloadBuilder
	remoteBuilders []PayloadBuilder
}

// LocalBuilder returns the local builder.
func (s *Service) LocalBuilder() PayloadBuilder {
	return s.localBuilder
}

// RequestBestBlock builds a new beacon block.
func (s *Service) RequestBestBlock(
	ctx context.Context, slot primitives.Slot,
) (consensus.BeaconKitBlock, error) {
	s.Logger().Info("our turn to propose a block 🙈", "slot", slot)
	// The goal here is to acquire a payload whose parent is the previously
	// finalized block, such that, if this payload is accepted, it will be
	// the next finalized block in the chain. A byproduct of this design
	// is that we get the nice property of lazily propogating the finalized
	// and safe block hashes to the execution client.

	// // // TODO: SIGN UR RANDAO THINGY HERE OR SOMETHING.
	// _ = s.beaconKitValKey
	// // _, err := s.beaconKitValKey.Key.PrivKey.Sign([]byte("hello world"))
	// // if err != nil {
	// // 	return nil, err
	// // }

	parentBlockRoot := s.BeaconState(ctx).GetParentBlockRoot()

	// Create a new empty block from the current state.
	beaconBlock, err := consensus.EmptyBeaconKitBlock(
		slot, parentBlockRoot, s.ActiveForkVersionForSlot(slot),
	)
	if err != nil {
		return nil, err
	}

	// TODO: right now parent is ALWAYS the previously finalized. But
	// maybe not forever?
	parentEth1Hash := s.BeaconState(ctx).GetFinalizedEth1BlockHash()

	// Get the payload for the block.
	payload, blobsBundle, overrideBuilder, err := s.localBuilder.GetBestPayload(
		ctx, slot, parentBlockRoot, parentEth1Hash,
	)
	if err != nil {
		return nil, err
	}

	p := pool.New().WithErrors()
	// Using goroutines here with only one task
	// is unnecessary, but it makes sense if we
	// want to add more tasks in the future.
	p.Go(func() error {
		return s.handleLogs(ctx, common.BytesToHash(payload.GetBlockHash()))
	})
	err = p.Wait()
	if err != nil {
		return nil, err
	}

	// TODO: Dencun
	_ = blobsBundle

	// TODO: allow external block builders to override the payload.
	_ = overrideBuilder

	// Assemble a new block with the payload.
	if err = beaconBlock.AttachExecution(payload); err != nil {
		return nil, err
	}

	// Dequeue deposits, up to MaxDepositsPerBlock, from the deposit queue.
	expectedDeposits, err := s.ss.DequeueDeposits(ctx)
	if err != nil {
		return nil, err
	}

	// Attach the deposits to the block.
	if err = beaconBlock.AttachDeposits(expectedDeposits); err != nil {
		return nil, err
	}

	// Return the block.
	return beaconBlock, nil
}

// handleLogs processes logs into values and
// does the appropriate action based on the log type.
func (s *Service) handleLogs(ctx context.Context, blkHash common.Hash) error {
	// Process logs in the eth1 block into values.
	var logValues []*reflect.Value
	logValues, err := s.es.ProcessLogsInETH1Block(ctx, blkHash)
	if err != nil {
		return err
	}
	// Process the log values based on their types.
	// Deposits are accepted into the deposit queue.
	deposits := make([]*consensusv1.Deposit, 0, len(logValues))
	for _, logValue := range logValues {
		logType := reflect.TypeOf(logValue.Interface())
		if logType.Kind() == reflect.Ptr {
			logType = logType.Elem()
		}
		switch logType {
		case stakinglogs.DepositType:
			deposit, ok := logValue.Interface().(*consensusv1.Deposit)
			if !ok {
				return errors.New("could not cast log value to deposit")
			}
			deposits = append(deposits, deposit)
		case stakinglogs.WithdrawalType:
		}
	}
	if err = s.ss.AcceptDepositsIntoQueue(ctx, deposits); err != nil {
		return err
	}
	return nil
}
