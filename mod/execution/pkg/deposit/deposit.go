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

	"github.com/berachain/beacon-kit/mod/execution/pkg/deposit/abi"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// BeaconDepositContract is a struct that holds a pointer to an ABI.
type BeaconDepositContract[
	DepositT any,
	WithdrawalCredentialsT ~[32]byte,
] struct {
	// BeaconDepositContract is a pointer to the codegen ABI binding.
	*abi.BeaconDepositContract

	// newDepositFn is a function that creates a new deposit.
	newDepositFn func(
		pubkey crypto.BLSPubkey,
		credentials WithdrawalCredentialsT,
		amount math.Gwei,
		signature crypto.BLSSignature,
		index uint64,
	) *DepositT
}

// NewBeaconDepositContract creates a new BeaconDepositContract.
func NewBeaconDepositContract[
	DepositT any,
	WithdrawalCredentialsT ~[32]byte,
](
	address common.ExecutionAddress,
	client bind.ContractBackend,
	newDepositFn func(
		pubkey crypto.BLSPubkey,
		credentials WithdrawalCredentialsT,
		amount math.Gwei,
		signature crypto.BLSSignature,
		index uint64,
	) *DepositT,
) (*BeaconDepositContract[
	DepositT,
	WithdrawalCredentialsT,
], error) {
	contract, err := abi.NewBeaconDepositContract(
		address, client,
	)
	if err != nil {
		return nil, err
	}

	return &BeaconDepositContract[
		DepositT,
		WithdrawalCredentialsT,
	]{
		BeaconDepositContract: contract,
	}, nil
}

// GetDeposits gets deposits from the deposit contract.
func (bdc *BeaconDepositContract[
	DepositT,
	WithdrawalCredentialsT,
]) GetDeposits(
	ctx context.Context,
	blockNumber uint64,
) ([]*DepositT, error) {
	logs, err := bdc.FilterDeposit(
		&bind.FilterOpts{
			Context: ctx,
			Start:   blockNumber,
			End:     &blockNumber,
		},
	)
	if err != nil {
		return nil, err
	}

	deposits := make([]*DepositT, 0)
	for logs.Next() {
		deposits = append(deposits, bdc.newDepositFn(
			crypto.BLSPubkey(logs.Event.Pubkey),
			WithdrawalCredentialsT(logs.Event.Credentials),
			math.U64(logs.Event.Amount),
			crypto.BLSSignature(logs.Event.Signature),
			logs.Event.Index,
		))
	}

	return deposits, nil
}
