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

	asynctypes "github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

type BeaconBlockBody[
	DepositT any,
	ExecutionPayloadT ExecutionPayload,
] interface {
	GetDeposits() []DepositT
	GetExecutionPayload() ExecutionPayloadT
}

// BeaconBlock is an interface for beacon blocks.
type BeaconBlock[BeaconBlockBodyT any] interface {
	GetSlot() math.U64
	GetBody() BeaconBlockBodyT
}

// BlockEvent is an interface for block events.
type BlockEvent[
	DepositT any,
	BeaconBlockBodyT BeaconBlockBody[DepositT, ExecutionPayloadT],
	BeaconBlockT BeaconBlock[BeaconBlockBodyT],
	ExecutionPayloadT ExecutionPayload,
] interface {
	ID() asynctypes.EventID
	Is(asynctypes.EventID) bool
	Data() BeaconBlockT
	Context() context.Context
}

// ExecutionPayload is an interface for execution payloads.
type ExecutionPayload interface {
	GetNumber() math.U64
}

// Contract is the ABI for the deposit contract.
type Contract[DepositT any] interface {
	// ReadDeposits reads deposits from the deposit contract.
	ReadDeposits(
		ctx context.Context,
		blockNumber math.U64,
	) ([]DepositT, error)
}

// Deposit is an interface for deposits.
type Deposit[DepositT, WithdrawalCredentialsT any] interface {
	// New creates a new deposit.
	New(
		crypto.BLSPubkey,
		WithdrawalCredentialsT,
		math.U64,
		crypto.BLSSignature,
		uint64,
	) DepositT
	// GetIndex returns the index of the deposit.
	GetIndex() math.U64
}

// Store defines the interface for managing deposit operations.
type Store[DepositT any] interface {
	// Prune prunes the deposit store of [start, end)
	Prune(index uint64, numPrune uint64) error
	// EnqueueDeposits adds a list of deposits to the deposit store.
	EnqueueDeposits(deposits []DepositT) error
}

// TelemetrySink is an interface for sending metrics to a telemetry backend.
type TelemetrySink interface {
	// IncrementCounter increments a counter metric identified by the provided
	// keys.
	IncrementCounter(key string, args ...string)
}
