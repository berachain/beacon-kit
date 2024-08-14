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

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// Client is an interface for the client.
type Client interface {
	GetLogsAtBlockNumber(
		ctx context.Context,
		number math.U64,
		address common.ExecutionAddress,
	) ([]engineprimitives.Log, error)
}

// WrappedBeaconDepositContract is a struct that holds a pointer to an ABI.
type WrappedBeaconDepositContract[
	DepositT Deposit[DepositT, WithdrawalCredentialsT],
	WithdrawalCredentialsT ~[32]byte,
] struct {
	client  Client
	address common.ExecutionAddress
}

// NewWrappedBeaconDepositContract creates a new BeaconDepositContract.
func NewWrappedBeaconDepositContract[
	DepositT Deposit[DepositT, WithdrawalCredentialsT],
	WithdrawalCredentialsT ~[32]byte,
](
	address common.ExecutionAddress,
	client Client,
) (*WrappedBeaconDepositContract[
	DepositT,
	WithdrawalCredentialsT,
], error) {

	return &WrappedBeaconDepositContract[
		DepositT,
		WithdrawalCredentialsT,
	]{
		client:  client,
		address: address,
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
	logs, err := dc.client.GetLogsAtBlockNumber(
		ctx,
		blkNum,
		dc.address,
	)
	if err != nil {
		return nil, err
	}

	var d DepositT
	deposits := make([]DepositT, 0)
	for _, log := range logs {
		d = d.Empty()
		if err := d.UnmarshalLog(log); err != nil {
			return nil, err
		}
		deposits = append(deposits)
	}

	return deposits, nil
}
