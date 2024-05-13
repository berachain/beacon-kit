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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// NewDepositFn is a function that creates a new deposit from
// the given parameters.
type NewDepositFn[
	DepositT any, WithdrawalCredentialsT ~[32]byte,
] func(
	pubkey crypto.BLSPubkey,
	credentials WithdrawalCredentialsT,
	amount math.Gwei,
	signature crypto.BLSSignature,
	index uint64,
) DepositT

// WrappedBeaconDepositContract is a struct that holds a pointer to an ABI.
//
//go:generate go run github.com/ethereum/go-ethereum/cmd/abigen --abi=../../../../contracts/out/BeaconDepositContract.sol/BeaconDepositContract.abi.json --pkg=deposit --type=BeaconDepositContract --out=bdc.go
type WrappedBeaconDepositContract[
	DepositT any,
	WithdrawalCredentialsT ~[32]byte,
] struct {
	// BeaconDepositContract is a pointer to the codegen ABI binding.
	BeaconDepositContract

	// newDepositFn is a function that creates a new deposit.
	newDepositFn NewDepositFn[DepositT, WithdrawalCredentialsT]
}

// NewWrappedBeaconDepositContract creates a new BeaconDepositContract.
func NewWrappedBeaconDepositContract[
	DepositT any,
	WithdrawalCredentialsT ~[32]byte,
](
	address common.ExecutionAddress,
	client bind.ContractBackend,
	newDepositFn NewDepositFn[DepositT, WithdrawalCredentialsT],
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

	if newDepositFn == nil {
		return nil, errors.New("newDepositFn must not be nil")
	}

	return &WrappedBeaconDepositContract[
		DepositT,
		WithdrawalCredentialsT,
	]{
		BeaconDepositContract: *contract,
		newDepositFn:          newDepositFn,
	}, nil
}

// GetDeposits gets deposits from the deposit contract.
func (bdc *WrappedBeaconDepositContract[
	DepositT,
	WithdrawalCredentialsT,
]) GetDeposits(
	ctx context.Context,
	blkNum uint64,
) ([]DepositT, error) {
	logs, err := bdc.FilterDeposit(
		&bind.FilterOpts{
			Context: ctx,
			Start:   blkNum,
			End:     &blkNum,
		},
	)
	if err != nil {
		return nil, err
	}

	deposits := make([]DepositT, 0)
	for logs.Next() {
		deposits = append(deposits, bdc.newDepositFn(
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
