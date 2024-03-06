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
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/sourcegraph/conc/iter"
)

// ProcessLog processes a single log
// received from the execution client.
func (f *Factory) ProcessLog(
	log *ethtypes.Log,
) (LogContainer, error) {
	// Skip logs that are not registered with the factory.
	// They may be from unregistered contracts (defensive check)
	// or emitted from unregistered events in the registered contracts.
	if !f.IsRegisteredLog(log) {
		//nolint:nilnil // nil is expected here
		return nil, nil
	}

	val, err := f.UnmarshalEthLog(log)
	if err != nil {
		return nil, err
	}
	return &Container{
		value:       &val,
		sig:         log.Topics[0],
		index:       uint64(log.Index),
		blockNumber: log.BlockNumber,
		blockHash:   log.BlockHash,
	}, nil
}

// ProcessLogs processes the logs received
// from the execution client in parallel.
// The order of the logs does not matter
// since the cache will sort them by
// block number and their indices.
func (f *Factory) ProcessLogs(
	logs []ethtypes.Log,
) ([]LogContainer, error) {
	// Process logs in parallel
	containers, err := iter.MapErr(
		logs,
		func(log *ethtypes.Log) (LogContainer, error) {
			return f.ProcessLog(log)
		})

	if err != nil {
		return nil, err
	}

	// Filter out nil values
	noNilContainers := make([]LogContainer, 0, len(containers))
	for _, container := range containers {
		if container != nil {
			noNilContainers = append(noNilContainers, container)
		}
	}

	return noNilContainers, nil
}
