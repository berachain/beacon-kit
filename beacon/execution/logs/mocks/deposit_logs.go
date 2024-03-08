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

package mocks

import (
	"github.com/berachain/beacon-kit/beacon/staking/logs"
	"github.com/berachain/beacon-kit/contracts/abi"
	"github.com/berachain/beacon-kit/primitives"
	"github.com/cockroachdb/errors"
	coretypes "github.com/ethereum/go-ethereum/core/types"
)

// CreateDepositLogs creates mock deposit logs.
func CreateDepositLogs(
	numDepositLogs int,
	factor int,
	contractAddress primitives.ExecutionAddress,
	blkNum uint64,
) ([]coretypes.Log, error) {
	if numDepositLogs <= 0 || factor <= 0 {
		return nil, errors.New("invalid input")
	}

	depositContractAbi, err := abi.BeaconDepositContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	} else if depositContractAbi == nil {
		return nil, errors.New("abi not found")
	}
	event := depositContractAbi.Events[logs.DepositName]

	// Create deposit logs.
	numLogs := factor*(numDepositLogs-1) + 1

	mockLogs := make([]coretypes.Log, 0, numLogs)
	for i := 0; i < numLogs; i++ {
		var data []byte
		// Create a log from the deposit.
		data, err = event.Inputs.Pack(
			[]byte("pubkey"),
			[]byte("12345678901234567890123456789012"),
			//#nosec:G701 // no overflow
			uint64(i),
			[]byte("signature"),
			//#nosec:G701 // no overflow
			uint64(i),
		)
		if err != nil {
			return nil, err
		}
		log := &coretypes.Log{
			Topics:      []primitives.ExecutionHash{event.ID},
			Data:        data,
			BlockNumber: blkNum,
			BlockHash:   [32]byte{byte(blkNum)},
		}
		if i%factor == 0 {
			log.Address = contractAddress
		}

		mockLogs = append(mockLogs, *log)
	}
	return mockLogs, nil
}
