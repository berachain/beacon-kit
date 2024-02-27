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
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	logmocks "github.com/itsdevbear/bolaris/beacon/execution/logs/mocks"
	"github.com/itsdevbear/bolaris/beacon/staking"
	"github.com/itsdevbear/bolaris/contracts/abi"
	"github.com/itsdevbear/bolaris/types/consensus"
)

// CreateDepositLogs creates mock deposit logs.
func CreateDepositLogs(
	numDepositLogs int,
	factor int,
	contractAddress ethcommon.Address,
	blkNum uint64,
) ([]ethtypes.Log, error) {
	stakingAbi, err := abi.StakingMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	// Create deposit logs.
	numLogs := factor*(numDepositLogs-1) + 1

	logs := make([]ethtypes.Log, 0, numLogs)
	for i := 0; i < numLogs; i++ {
		deposit := consensus.NewDeposit(
			[]byte("pubkey"),
			uint64(i),
			[]byte("12345678901234567890"),
		)
		var log *ethtypes.Log
		log, err = logmocks.NewLogFromDeposit(
			stakingAbi.Events[staking.DepositName],
			deposit,
		)
		if err != nil {
			return nil, err
		}

		if i%factor == 0 {
			log.Address = contractAddress
			log.BlockNumber = blkNum
		} else if i%2 == 0 {
			log.BlockNumber = blkNum
		}

		logs = append(logs, *log)
	}
	return logs, nil
}
