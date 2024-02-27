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
	"github.com/itsdevbear/bolaris/beacon/execution/logs"
	"github.com/itsdevbear/bolaris/beacon/staking"
	"github.com/itsdevbear/bolaris/contracts/abi"
)

var _ logs.Service = (*Service)(nil)

// Service is a mock service for testing.
// It implements the logs.Service interface,
// so that it can send requests to the log factory.
type Service struct{}

// GetLogRequests returns a list of log requests
// to be sent to the log factory.
func (s *Service) GetLogRequests() ([]logs.LogRequest, error) {
	depositContractAddr := ethcommon.HexToAddress("0x1234")
	depositContractAbi, err := abi.StakingMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	allocator := logs.New[logs.TypeAllocator](
		logs.WithABI(depositContractAbi),
		logs.WithNameAndType(
			staking.DepositSig,
			staking.DepositName,
			staking.DepositType,
		),
		logs.WithNameAndType(
			staking.WithdrawalSig,
			staking.WithdrawalName,
			staking.WithdrawalType,
		),
	)

	request := logs.LogRequest{
		ContractAddress: depositContractAddr,
		Allocator:       allocator,
	}

	return []logs.LogRequest{request}, nil
}
