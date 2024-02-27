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

	"github.com/ethereum/go-ethereum/common"
	"github.com/itsdevbear/bolaris/beacon/execution/logs"
	enginev1 "github.com/itsdevbear/bolaris/engine/types/v1"
	consensusv1 "github.com/itsdevbear/bolaris/types/consensus/v1"
)

const (
	// Name of the Deposit event
	// in the deposit contract.
	depositName = "Deposit"

	// Name the Withdrawal event
	// in the deposit contract.
	withdrawalName = "Withdrawal"
)

//nolint:gochecknoglobals,lll // Avoid re-allocating these variables.
var (
	// Signature and type of the Deposit event
	// in the deposit contract.
	// keccak256("Deposit(bytes,bytes,uint64)").
	DepositSig  = common.HexToHash("163244a852f099315d72dcfbb5b1031ca0365543f2ac1849bdb69b01d8648b18")
	depositType = reflect.TypeOf(consensusv1.Deposit{})

	// Signature and type of the Withdrawal event
	// in the deposit contract.
	// keccak256("Withdrawal(bytes,bytes,uint64)").
	WithdrawalSig  = common.HexToHash("3cd2410b5f33d39669545e9f38ba4d4c6318f2b8f1a33f001bf6c03b2ab180b4")
	withdrawalType = reflect.TypeOf(enginev1.Withdrawal{})
)

// GetLogRequests returns a list of log requests from the staking service
// to be sent to the log factory in the execution service.
func (s *Service) GetLogRequests() ([]logs.LogRequest, error) {
	depositContractAddr := s.BeaconCfg().Execution.DepositContractAddress

	allocator := logs.New[logs.TypeAllocator](
		logs.WithABI(s.abi),
		logs.WithNameAndType(DepositSig, depositName, depositType),
		logs.WithNameAndType(WithdrawalSig, withdrawalName, withdrawalType),
	)

	request := logs.LogRequest{
		ContractAddress: depositContractAddr,
		Allocator:       allocator,
	}

	return []logs.LogRequest{request}, nil
}
