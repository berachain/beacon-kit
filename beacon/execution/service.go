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

package execution

import (
	"context"
	"math/big"
	"reflect"

	"cosmossdk.io/errors"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/itsdevbear/bolaris/engine"
	"github.com/itsdevbear/bolaris/runtime/service"
	"github.com/itsdevbear/bolaris/types/consensus/primitives"
	enginetypes "github.com/itsdevbear/bolaris/types/engine"
	enginev1 "github.com/itsdevbear/bolaris/types/engine/v1"
)

//nolint:gochecknoglobals // TODO
var rpcFinalizedBlockNumber = big.NewInt(int64(rpc.FinalizedBlockNumber))

// Service is responsible for delivering beacon chain notifications to
// the execution client and processing logs received from the execution client.
type Service struct {
	service.BaseService
	// engine gives the notifier access to the engine api of the execution
	// client.
	engine engine.Caller
	// logFactory is the factory for creating objects from Ethereum logs.
	logFactory LogFactory
}

// Start spawns any goroutines required by the service.
func (s *Service) Start(ctx context.Context) {
	s.engine.Start(ctx)
}

// Status returns error if the service is not considered healthy.
func (s *Service) Status() error {
	if !s.engine.IsConnected() {
		return ErrExecutionClientDisconnected
	}
	return nil
}

// NotifyForkchoiceUpdate notifies the execution client of a forkchoice update.
// TODO: handle the bools better i.e attrs, retry, async.
func (s *Service) NotifyForkchoiceUpdate(
	ctx context.Context, fcuConfig *FCUConfig,
) (*enginev1.PayloadIDBytes, error) {
	var (
		err       error
		payloadID *enginev1.PayloadIDBytes
	)
	// Push the forkchoice request to the forkchoice dispatcher, we want to
	// block until
	if e := s.GCD().GetQueue(forkchoiceDispatchQueue).Sync(func() {
		payloadID, err = s.notifyForkchoiceUpdate(ctx, fcuConfig)
	}); e != nil {
		return nil, e
	}

	return payloadID, err
}

// GetPayload returns the payload and blobs bundle for the given slot.
func (s *Service) GetPayload(
	ctx context.Context, payloadID primitives.PayloadID, slot primitives.Slot,
) (enginetypes.ExecutionPayload, *enginev1.BlobsBundle, bool, error) {
	return s.engine.GetPayload(ctx, payloadID, slot)
}

// NotifyNewPayload notifies the execution client of a new payload.
// It returns true if the EL has returned VALID for the block.
func (s *Service) NotifyNewPayload(
	ctx context.Context,
	payload enginetypes.ExecutionPayload,
	slot primitives.Slot,
) (bool, error) {
	return s.notifyNewPayload(ctx, payload, slot)
}

// GetLogsInFinalizedETH1Block gets logs in the finalized block
// received from the execution client and uses LogFactory to
// convert them into appropriate objects that can be consumed
// by other services.
func (s *Service) GetLogsInFinalizedETH1Block(
	ctx context.Context,
	blkNum uint64,
) ([]reflect.Value, error) {
	// Get the block from the eth1 client and
	// check if the block is safe to process.
	// TODO: Do we want to come up with a heuristic around
	// when we check the execution client,
	// vs when we check the forkchoice store.
	finalBlock, err := s.engine.BlockByNumber(
		ctx, rpcFinalizedBlockNumber,
	)
	if err != nil {
		return nil, err
	}

	// Ensure we don't start processing the logs
	// of a block that is ahead of the safe block.
	if finalBlock.Number().Uint64() < blkNum {
		return nil, errors.Wrapf(
			ErrProcessingUnfinalizedBlock,
			"safe block %d is behind block %d", finalBlock.Number(), blkNum,
		)
	}

	// Gather all the logs corresponding to the handlers from this block.
	registeredAddrs := s.logFactory.GetRegisteredAddresses()
	logsInBlock, err := s.engine.GetLogs(ctx, blkNum, blkNum, registeredAddrs)
	if err != nil {
		return nil, err
	}

	// Unmarshal the logs into objects.
	vals := make([]reflect.Value, 0, len(logsInBlock))
	for i, log := range logsInBlock {
		// Skip logs that are not from the block we are processing
		// This should never happen, but defensively check anyway.
		if log.BlockNumber != blkNum {
			continue
		}

		// Skip logs that are not registered with the factory.
		// They may be from unregistered contracts (defensive check)
		// or emitted from unregistered events in the registered contracts.
		if !s.logFactory.IsRegisteredLog(&logsInBlock[i]) {
			continue
		}

		var val reflect.Value
		val, err = s.logFactory.UnmarshalEthLog(&logsInBlock[i])
		if err != nil {
			return nil, errors.Wrap(err, "could not unmarshal log")
		}
		vals = append(vals, val)
	}
	return vals, nil
}
