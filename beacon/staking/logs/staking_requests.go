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

	beacontypesv1 "github.com/berachain/beacon-kit/beacon/core/types/v1"
	"github.com/berachain/beacon-kit/beacon/execution/logs"
	"github.com/berachain/beacon-kit/contracts/abi"
	enginetypes "github.com/berachain/beacon-kit/engine/types"
	"github.com/berachain/beacon-kit/primitives"
)

const (
	// Name of the Deposit event
	// in the deposit contract.
	DepositName = "Deposit"

	// Name of the Redirect event
	// in the deposit contract.
	RedirectName = "Redirect"

	// Name the Withdrawal event
	// in the deposit contract.
	WithdrawalName = "Withdrawal"
)

//nolint:gochecknoglobals // Avoid re-allocating these variables.
var (
	// Signature and type of the Deposit event
	// in the deposit contract.
	// keccak256("Deposit(bytes,bytes,uint64)").
	// 0x1f39b85dd1a529b31e0cd61e5609e1feca0e08e2103fe319fbd3dd5a0c7b68df.
	DepositSig = [32]byte{
		0x1f, 0x39, 0xb8, 0x5d, 0xd1, 0xa5, 0x29, 0xb3,
		0x1e, 0x0c, 0xd6, 0x1e, 0x56, 0x09, 0xe1, 0xfe,
		0xca, 0x0e, 0x08, 0xe2, 0x10, 0x3f, 0xe3, 0x19,
		0xfb, 0xd3, 0xdd, 0x5a, 0x0c, 0x7b, 0x68, 0xdf,
	}

	DepositType = reflect.TypeOf(beacontypesv1.Deposit{})

	// Signature and type of the Redirect event
	// in the deposit contract.
	// keccak256("Redirect(bytes,bytes,bytes,uint64)").
	// 0xe161f5842757f257346b360594d094b7fa530f9404e93a80bf18bd8b14f9258f.
	RedirectSig = [32]byte{
		0xe1, 0x61, 0xf5, 0x84, 0x27, 0x57, 0xf2, 0x57,
		0x34, 0x6b, 0x36, 0x05, 0x94, 0xd0, 0x94, 0xb7,
		0xfa, 0x53, 0x0f, 0x94, 0x04, 0xe9, 0x3a, 0x80,
		0xbf, 0x18, 0xbd, 0x8b, 0x14, 0xf9, 0x25, 0x8f,
	}

	// Signature and type of the Withdraw event
	// in the deposit contract.
	// keccak256("Withdrawal(bytes,bytes,bytes,uint64)").
	// 0x33689ed7aade9d90db44595c06c3bb95713b8332a7dec83be61556783c2154b1.
	WithdrawalSig = [32]byte{
		0x33, 0x68, 0x9e, 0xd7, 0xaa, 0xde, 0x9d, 0x90,
		0xdb, 0x44, 0x59, 0x5c, 0x06, 0xc3, 0xbb, 0x95,
		0x71, 0x3b, 0x83, 0x32, 0xa7, 0xde, 0xc8, 0x3b,
		0xe6, 0x15, 0x56, 0x78, 0x3c, 0x21, 0x54, 0xb1,
	}
	WithdrawalType = reflect.TypeOf(enginetypes.Withdrawal{})
)

// NewStakingRequest returns a log request for the staking service.
func NewStakingRequest(
	depositContractAddress primitives.ExecutionAddress,
) (*logs.LogRequest, error) {
	depositContractAbi, err := abi.BeaconDepositContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	allocator := logs.New[logs.TypeAllocator](
		logs.WithABI(depositContractAbi),
		logs.WithNameAndType(DepositSig, DepositName, DepositType),
		logs.WithNameAndType(WithdrawalSig, WithdrawalName, WithdrawalType),
	)

	return &logs.LogRequest{
		ContractAddress: depositContractAddress,
		Allocator:       allocator,
	}, nil
}
