// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package deposit

import (
	"context"
	"errors"
	"fmt"

	gethprimitives "github.com/berachain/beacon-kit/geth-primitives"
	"github.com/berachain/beacon-kit/geth-primitives/bind"
	"github.com/berachain/beacon-kit/geth-primitives/deposit"
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
)

// WrappedDepositContract is a struct that holds a pointer to an ABI.
type WrappedDepositContract[
	DepositT Deposit[DepositT, WithdrawalCredentialsT],
	WithdrawalCredentialsT ~[32]byte,
] struct {
	// DepositContractFilterer is a pointer to the codegen ABI binding.
	deposit.DepositContractFilterer
}

// NewWrappedDepositContract creates a new DepositContract.
func NewWrappedDepositContract[
	DepositT Deposit[DepositT, WithdrawalCredentialsT],
	WithdrawalCredentialsT ~[32]byte,
](
	address common.ExecutionAddress,
	client bind.ContractFilterer,
) (*WrappedDepositContract[
	DepositT,
	WithdrawalCredentialsT,
], error) {
	contract, err := deposit.NewDepositContractFilterer(
		gethprimitives.ExecutionAddress(address), client,
	)

	if err != nil {
		return nil, err
	} else if contract == nil {
		return nil, errors.New("contract must not be nil")
	}

	return &WrappedDepositContract[
		DepositT,
		WithdrawalCredentialsT,
	]{
		DepositContractFilterer: *contract,
	}, nil
}

// ReadDeposits reads deposits from the deposit contract.
func (dc *WrappedDepositContract[
	DepositT,
	WithdrawalCredentialsT,
]) ReadDeposits(
	ctx context.Context,
	blkNum math.U64,
) ([]DepositT, error) {
	logs, err := dc.FilterDeposit(
		&bind.FilterOpts{
			Context: ctx,
			Start:   blkNum.Unwrap(),
			End:     (*uint64)(&blkNum),
		},
	)
	if err != nil {
		return nil, err
	}

	deposits := make([]DepositT, 0)
	for logs.Next() {
		var (
			cred   bytes.B32
			pubKey bytes.B48
			d      DepositT
			sign   bytes.B96
		)
		pubKey, err = bytes.ToBytes48(logs.Event.Pubkey)
		if err != nil {
			return nil, fmt.Errorf("failed reading pub key: %w", err)
		}
		cred, err = bytes.ToBytes32(logs.Event.Credentials)
		if err != nil {
			return nil, fmt.Errorf("failed reading credentials: %w", err)
		}
		sign, err = bytes.ToBytes96(logs.Event.Signature)
		if err != nil {
			return nil, fmt.Errorf("failed reading signature: %w", err)
		}
		deposits = append(deposits, d.New(
			pubKey,
			WithdrawalCredentialsT(cred),
			math.U64(logs.Event.Amount),
			sign,
			logs.Event.Index,
		))
	}

	return deposits, nil
}
