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
	"reflect"

	beacontypesv1 "github.com/itsdevbear/bolaris/beacon/core/types/v1"
	"github.com/itsdevbear/bolaris/beacon/execution/logs"
	"github.com/itsdevbear/bolaris/contracts/abi"
	enginetypes "github.com/itsdevbear/bolaris/engine/types"
	"github.com/itsdevbear/bolaris/primitives"
)

const (
	// Name of the Deposit event
	// in the deposit contract.
	DepositName = "Deposit"

	// Name the Withdrawal event
	// in the deposit contract.
	WithdrawalName = "Withdrawal"
)

//nolint:gochecknoglobals // Avoid re-allocating these variables.
var (
	// Signature and type of the Deposit event
	// in the deposit contract.
	// keccak256("Deposit(bytes,bytes,uint64)").
	DepositSig = [32]byte{
		0x16,
		0x32,
		0x44,
		0xa8,
		0x52,
		0xf0,
		0x99,
		0x31,
		0x5d,
		0x72,
		0xdc,
		0xfb,
		0xb5,
		0xb1,
		0x03,
		0x1c,
		0xa0,
		0x36,
		0x55,
		0x43,
		0xf2,
		0xac,
		0x18,
		0x49,
		0xbd,
		0xb6,
		0x9b,
		0x01,
		0xd8,
		0x64,
		0x8b,
		0x18,
	}

	DepositType = reflect.TypeOf(beacontypesv1.Deposit{})

	// Signature and type of the Withdrawal event
	// in the deposit contract.
	// keccak256("Withdrawal(bytes,bytes,uint64)").
	WithdrawalSig = [32]byte{
		0x3c,
		0xd2,
		0x41,
		0x0b,
		0x5f,
		0x33,
		0xd3,
		0x96,
		0x69,
		0x54,
		0x5e,
		0x9f,
		0x38,
		0xba,
		0x4d,
		0x4c,
		0x63,
		0x18,
		0xf2,
		0xb8,
		0xf1,
		0xa3,
		0x3f,
		0x00,
		0x1b,
		0xf6,
		0xc0,
		0x3b,
		0x2a,
		0xb1,
		0x80,
		0xb4,
	}
	WithdrawalType = reflect.TypeOf(enginetypes.Withdrawal{})
)

// NewStakingRequest returns a log request for the staking service.
func NewStakingRequest(
	depositContractAddress primitives.ExecutionAddress,
) (*logs.LogRequest, error) {
	stakingAbi, err := abi.StakingMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	allocator := logs.New[logs.TypeAllocator](
		logs.WithABI(stakingAbi),
		logs.WithNameAndType(DepositSig, DepositName, DepositType),
		logs.WithNameAndType(WithdrawalSig, WithdrawalName, WithdrawalType),
	)

	return &logs.LogRequest{
		ContractAddress: depositContractAddress,
		Allocator:       allocator,
	}, nil
}
