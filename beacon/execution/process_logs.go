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
	"runtime"
	"sort"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

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
	// Use a semaphore to limit the number
	// of concurrent goroutines.
	sem := make(chan bool, runtime.NumCPU())
	logResults := make([]*logResult, len(logs))
	// Unmarshal the logs into objects.
	for i := range logs {
		i := i
		log := &logs[i]
		// Acquire a semaphore slot.
		sem <- true
		eg.Go(func() error {
			defer func() {
				// Release the semaphore slot when done.
				<-sem
			}()
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
				return errors.Wrapf(err, "%s", "could not unmarshal log")
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
