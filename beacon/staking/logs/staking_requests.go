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

	beacontypes "github.com/berachain/beacon-kit/beacon/core/types"
	"github.com/berachain/beacon-kit/beacon/execution/logs"
	stakingabi "github.com/berachain/beacon-kit/contracts/abi"
	enginetypes "github.com/berachain/beacon-kit/engine/types"
	"github.com/berachain/beacon-kit/primitives"
	"github.com/ethereum/go-ethereum/crypto"
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
	DepositContractABI, _ = stakingabi.BeaconDepositContractMetaData.GetAbi()

	// Signature and type of the Deposit event
	// in the deposit contract.
	DepositSig = crypto.Keccak256Hash(
		[]byte(DepositName + "(bytes,bytes,uint64,bytes,uint64)"),
	)
	DepositType = reflect.TypeOf(beacontypes.Deposit{})

	// Signature and type of the Redirect event
	// in the deposit contract.
	RedirectSig = crypto.Keccak256Hash(
		[]byte(RedirectName + "(bytes,bytes,bytes,uint64,uint64)"),
	)
	// RedirectType = reflect.TypeOf(enginetypes.Redirect{}).

	// Signature and type of the Withdraw event
	// in the deposit contract.
	WithdrawalSig = crypto.Keccak256Hash(
		[]byte(WithdrawalName + "(bytes,bytes,bytes,uint64,uint64)"),
	)
	WithdrawalType = reflect.TypeOf(enginetypes.Withdrawal{})
)

// NewStakingRequest returns a log request for the staking service.
func NewStakingRequest(
	depositContractAddress primitives.ExecutionAddress,
) (*logs.LogRequest, error) {
	depositContractAbi, err := stakingabi.BeaconDepositContractMetaData.GetAbi()
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
