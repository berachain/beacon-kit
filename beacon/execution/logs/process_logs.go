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
	"fmt"
	"reflect"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/sourcegraph/conc/iter"
)

// ProcessLogs processes the logs received from the execution client
// in parallel but returns the values in the same order of the received logs.
func (f *Factory) ProcessLogs(
	logs []ethtypes.Log,
	blkNum uint64,
) ([]*reflect.Value, error) {
	logValues, err := iter.MapErr(
		logs,
		func(log *ethtypes.Log) (*reflect.Value, error) {
			// Skip logs that are not from the block we are processing
			// This should never happen, but defensively check anyway.
			if log.BlockNumber != blkNum {
				return nil, fmt.Errorf(
					"log from different block, expected %d, got %d",
					blkNum, log.BlockNumber,
				)
			}

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
			return &val, nil
		})

	if err != nil {
		return nil, err
	}

	// Filter out nil values
	vals := make([]*reflect.Value, 0, len(logValues))
	for _, val := range logValues {
		if val != nil {
			vals = append(vals, val)
		}
	}

	return vals, nil
}
