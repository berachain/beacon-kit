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

package state

import (
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// BeaconStateMarshallable represents an interface for a beacon state
// with generic types.
type BeaconStateMarshallable[
	T any,
	BeaconBlockHeaderT,
	Eth1DataT,
	ExecutionPayloadHeaderT,
	ForkT,
	ValidatorT any,
] interface {
	constraints.SSZMarshallable
	// New returns a new instance of the BeaconStateMarshallable.
	New(
		forkVersion uint32,
		genesisValidatorsRoot common.Bytes32,
		slot math.U64,
		fork ForkT,
		latestBlockHeader BeaconBlockHeaderT,
		blockRoots []common.Bytes32,
		stateRoots []common.Bytes32,
		eth1Data Eth1DataT,
		eth1DepositIndex uint64,
		latestExecutionPayloadHeader ExecutionPayloadHeaderT,
		validators []ValidatorT,
		balances []uint64,
		randaoMixes []common.Bytes32,
		nextWithdrawalIndex uint64,
		nextWithdrawalValidatorIndex math.U64,
		slashings []uint64, totalSlashing math.U64,
	) (T, error)
}

// Validator represents an interface for a validator with generic withdrawal
// credentials. WithdrawalCredentialsT is a type parameter that must implement
// the WithdrawalCredentials interface.
type Validator[WithdrawalCredentialsT WithdrawalCredentials] interface {
	// GetWithdrawalCredentials returns the withdrawal credentials of the
	// validator.
	GetWithdrawalCredentials() WithdrawalCredentialsT
	// IsFullyWithdrawable checks if the validator is fully withdrawable given a
	// certain Gwei amount and epoch.
	IsFullyWithdrawable(amount math.Gwei, epoch math.Epoch) bool
	// IsPartiallyWithdrawable checks if the validator is partially withdrawable
	// given two Gwei amounts.
	IsPartiallyWithdrawable(amount1 math.Gwei, amount2 math.Gwei) bool
}

// Withdrawal represents an interface for a withdrawal.
type Withdrawal[T any] interface {
	New(
		index math.U64,
		validator math.ValidatorIndex,
		address gethprimitives.ExecutionAddress,
		amount math.Gwei,
	) T
}

// WithdrawalCredentials represents an interface for withdrawal credentials.
type WithdrawalCredentials interface {
	// ToExecutionAddress converts the withdrawal credentials to an execution
	// address.
	ToExecutionAddress() (gethprimitives.ExecutionAddress, error)
}
