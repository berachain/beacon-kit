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
	"reflect"
	"sort"

	"cosmossdk.io/errors"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	engineclient "github.com/itsdevbear/bolaris/engine/client"
	enginetypes "github.com/itsdevbear/bolaris/engine/types"
	enginev1 "github.com/itsdevbear/bolaris/engine/types/v1"
	"github.com/itsdevbear/bolaris/runtime/service"
	"github.com/itsdevbear/bolaris/types/consensus/primitives"
	"golang.org/x/sync/errgroup"
)

// Service is responsible for delivering beacon chain notifications to
// the execution client and processing logs received from the execution client.
type Service struct {
	service.BaseService
	// engine gives the notifier access to the engine api of the execution
	// client.
	engine engineclient.Caller
	// logFactory is the factory for creating objects from Ethereum logs.
	logFactory LogFactory
}

// Start spawns any goroutines required by the service.
func (s *Service) Start(ctx context.Context) {
	go s.engine.Start(ctx)
}

// Status returns error if the service is not considered healthy.
func (s *Service) Status() error {
	return s.engine.Status()
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
	slot primitives.Slot,
	payload enginetypes.ExecutionPayload,
	versionedHashes []common.Hash,
	parentBlockRoot [32]byte,
) (bool, error) {
	return s.notifyNewPayload(
		ctx,
		slot,
		payload,
		versionedHashes,
		parentBlockRoot,
	)
}

// GetLogsInETH1Block gets logs in the Eth1 block
// received from the execution client and uses LogFactory to
// convert them into appropriate objects that can be consumed
// by other services.
func (s *Service) GetLogsInETH1Block(
	ctx context.Context,
	blkNum uint64,
) ([]*reflect.Value, error) {
	// Gather all the logs corresponding to the handlers from this block.
	registeredAddrs := s.logFactory.GetRegisteredAddresses()
	logsInBlock, err := s.engine.GetLogs(ctx, blkNum, blkNum, registeredAddrs)
	if err != nil {
		return nil, err
	}
	return s.ProcessLogs(logsInBlock, blkNum)
}

// ProcessLogs processes the logs received from the execution client
// in parallel but returns the values in the same order of the received logs.
func (s *Service) ProcessLogs(
	logs []ethtypes.Log,
	blkNum uint64,
) ([]*reflect.Value, error) {
	type logResult struct {
		val   *reflect.Value
		index int
	}

	eg := new(errgroup.Group)
	logResults := make([]*logResult, len(logs))
	// Unmarshal the logs into objects.
	for i := range logs {
		i := i
		log := &logs[i]
		eg.Go(func() error {
			// Skip logs that are not from the block we are processing
			// This should never happen, but defensively check anyway.
			if log.BlockNumber != blkNum {
				return nil
			}

			// Skip logs that are not registered with the factory.
			// They may be from unregistered contracts (defensive check)
			// or emitted from unregistered events in the registered contracts.
			if !s.logFactory.IsRegisteredLog(log) {
				return nil
			}

			val, err := s.logFactory.UnmarshalEthLog(log)
			if err != nil {
				return errors.Wrap(err, "could not unmarshal log")
			}
			logResults[i] = &logResult{val: &val, index: i}
			return nil
		})
	}

	// Wait for all goroutines to complete
	// or return the first error encountered.
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	// Sort the result by the index of the log in the block.
	sort.Slice(logResults, func(i, j int) bool {
		if logResults[i] == nil {
			return false
		}
		if logResults[j] == nil {
			return true
		}
		return logResults[i].index < logResults[j].index
	})

	vals := make([]*reflect.Value, 0, len(logResults))
	for _, res := range logResults {
		if res != nil {
			// At this point, we should not have any nil values
			// or non-nil errors in the result.
			vals = append(vals, res.val)
		}
	}
	return vals, nil
}
