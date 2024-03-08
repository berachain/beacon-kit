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
	"errors"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	beacontypes "github.com/itsdevbear/bolaris/beacon/core/types"
	"github.com/itsdevbear/bolaris/beacon/staking/logs"
	"github.com/itsdevbear/bolaris/contracts/abi"
	"github.com/itsdevbear/bolaris/primitives"
)

// CreateDepositLogs creates mock deposit logs.
func CreateDepositLogs(
	numDepositLogs int,
	factor int,
	contractAddress primitives.ExecutionAddress,
	blkNum uint64,
) ([]ethtypes.Log, error) {
	if numDepositLogs <= 0 || factor <= 0 {
		return nil, errors.New("invalid input")
	}

	stakingAbi, err := abi.StakingMetaData.GetAbi()
	if err != nil {
		return nil, err
	} else if stakingAbi == nil {
		return nil, errors.New("abi not found")
	}

	// Create deposit logs.
	numLogs := factor*(numDepositLogs-1) + 1

	mockLogs := make([]ethtypes.Log, 0, numLogs)
	for i := 0; i < numLogs; i++ {
		deposit := beacontypes.NewDeposit(
			[]byte("pubkey"),
			//#nosec:G701 // no overflow
			uint64(i),
			[]byte("12345678901234567890"),
		)
		var log *ethtypes.Log
		events := stakingAbi.Events
		if events == nil {
			return nil, errors.New("events not found")
		}
		log, err = NewLogFromDeposit(
			events[logs.DepositName],
			deposit,
		)
		if err != nil {
			return nil, err
		}

		log.BlockNumber = blkNum
		log.BlockHash = [32]byte{byte(blkNum)}
		if i%factor == 0 {
			log.Address = contractAddress
		}

		mockLogs = append(mockLogs, *log)
	}
	return mockLogs, nil
}
