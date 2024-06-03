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

package deposit

import (
	"context"
	"errors"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// WrappedBeaconDepositContract is a struct that holds a pointer to an ABI.
//
//go:generate go run github.com/ethereum/go-ethereum/cmd/abigen --abi=../../../../contracts/out/BeaconDepositContract.sol/BeaconDepositContract.abi.json --pkg=deposit --type=BeaconDepositContract --out=contract.abigen.go
type WrappedBeaconDepositContract[
	DepositT Deposit[DepositT, WithdrawalCredentialsT],
	WithdrawalCredentialsT ~[32]byte,
] struct {
	// BeaconDepositContract is a pointer to the codegen ABI binding.
	BeaconDepositContract
}

// NewWrappedBeaconDepositContract creates a new BeaconDepositContract.
func NewWrappedBeaconDepositContract[
	DepositT Deposit[DepositT, WithdrawalCredentialsT],
	WithdrawalCredentialsT ~[32]byte,
](
	address common.ExecutionAddress,
	client bind.ContractBackend,
) (*WrappedBeaconDepositContract[
	DepositT,
	WithdrawalCredentialsT,
], error) {
	contract, err := NewBeaconDepositContract(
		address, client,
	)

	if err != nil {
		return nil, err
	} else if contract == nil {
		return nil, errors.New("contract must not be nil")
	}

	return &WrappedBeaconDepositContract[
		DepositT,
		WithdrawalCredentialsT,
	]{
		BeaconDepositContract: *contract,
	}, nil
}

// ReadDeposits reads deposits from the deposit contract.
func (dc *WrappedBeaconDepositContract[
	DepositT,
	WithdrawalCredentialsT,
]) ReadDeposits(
	ctx context.Context,
	blkNum math.U64,
) ([]DepositT, error) {
	logs, err := dc.FilterDeposit(
		&bind.FilterOpts{
			Context: ctx,
			Start:   uint64(blkNum),
			End:     (*uint64)(&blkNum),
		},
	)
	if err != nil {
		return nil, err
	}

	deposits := make([]DepositT, 0)
	for logs.Next() {
		var d DepositT
		deposits = append(deposits, d.New(
			bytes.ToBytes48(logs.Event.Pubkey),
			WithdrawalCredentialsT(
				bytes.ToBytes32(logs.Event.Credentials)),
			math.U64(logs.Event.Amount),
			bytes.ToBytes96(logs.Event.Signature),
			logs.Event.Index,
		))
	}

	return deposits, nil
}
