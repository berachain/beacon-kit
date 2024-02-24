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

package staking

import (
	"reflect"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/itsdevbear/bolaris/beacon/execution/logs"
	"github.com/itsdevbear/bolaris/contracts/abi"
	enginev1 "github.com/itsdevbear/bolaris/engine/types/v1"
	consensusv1 "github.com/itsdevbear/bolaris/types/consensus/v1"
)

//nolint:gochecknoglobals // Avoid re-allocating these variables.
var (
	// Name, signature, and type of the Deposit event
	// in the deposit contract.
	depositName = "Deposit"
	depositSig  = ethcrypto.Keccak256Hash(
		[]byte("Deposit(bytes,bytes,uint64)"),
	)
	depositType = reflect.TypeOf(consensusv1.Deposit{})

	// Name, signature, and type of the Withdrawal event
	// in the deposit contract.
	withdrawalName = "Withdrawal"
	withdrawalSig  = ethcrypto.Keccak256Hash(
		[]byte("Withdrawal(bytes,bytes,uint64)"),
	)
	withdrawalType = reflect.TypeOf(enginev1.Withdrawal{})
)

// GetLogRequests returns a list of log requests from the staking service
// to be sent to the log factory in the execution service.
func (s *Service) GetLogRequests() ([]logs.LogRequest, error) {
	depositContractAddr := s.BeaconCfg().Execution.DepositContractAddress
	depositContractAbi, err := abi.StakingMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	allocator := logs.New[logs.TypeAllocator](
		logs.WithABI(depositContractAbi),
		logs.WithNameAndType(depositSig, depositName, depositType),
		logs.WithNameAndType(withdrawalSig, withdrawalName, withdrawalType),
	)

	request := logs.LogRequest{
		ContractAddress: depositContractAddr,
		Allocator:       allocator,
	}

	return []logs.LogRequest{request}, nil
}
