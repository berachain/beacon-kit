// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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

package backend

import (
	"context"

	"github.com/berachain/beacon-kit/consensus-types/deneb"
	dastore "github.com/berachain/beacon-kit/da/store"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/transition"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
	"github.com/berachain/beacon-kit/storage/block"
	depositdb "github.com/berachain/beacon-kit/storage/deposit"
)

// Node is the interface for a node.
type Node[ContextT any] interface {
	// CreateQueryContext creates a query context for a given height and proof
	// flag.
	CreateQueryContext(height int64, prove bool) (ContextT, error)
}

type StateProcessor interface {
	ProcessSlots(*statedb.StateDB, math.Slot) (transition.ValidatorUpdates, error)
}

// StorageBackend is the interface for the storage backend.
type StorageBackend interface {
	AvailabilityStore() *dastore.Store
	BlockStore() *block.KVStore[*deneb.BeaconBlock]
	DepositStore() *depositdb.KVStore
	StateFromContext(context.Context) *statedb.StateDB
}

// Validator represents an interface for a validator with generic withdrawal
// credentials.
type Validator interface {
	// GetWithdrawalCredentials returns the withdrawal credentials of the
	// validator.
	GetWithdrawalCredentials() deneb.WithdrawalCredentials
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
		address common.ExecutionAddress,
		amount math.Gwei,
	) T
}

// WithdrawalCredentials represents an interface for withdrawal credentials.
type WithdrawalCredentials interface {
	// ToExecutionAddress converts the withdrawal credentials to an execution
	// address.
	ToExecutionAddress() (common.ExecutionAddress, error)
}
