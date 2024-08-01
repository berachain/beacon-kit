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

package builder

import (
	"context"

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// AttributesFactory is the interface for the attributes factory.
type AttributesFactory[BeaconStateT, PayloadAttributesT any] interface {
	BuildPayloadAttributes(
		BeaconStateT,
		math.Slot,
		math.U64,
		common.Root,
	) (PayloadAttributesT, error)
}

// BeaconState defines the interface for accessing various state-related data
// required for block processing.
type BeaconState[
	ExecutionPayloadHeaderT ExecutionPayloadHeader,
	WithdrawalT any,
] interface {
	// GetRandaoMixAtIndex retrieves the RANDAO mix at a specified index.
	GetRandaoMixAtIndex(uint64) (common.Bytes32, error)
	// ExpectedWithdrawals lists the expected withdrawals in the current state.
	ExpectedWithdrawals() ([]WithdrawalT, error)
	// GetLatestExecutionPayloadHeader fetches the most recent execution payload
	// header.
	GetLatestExecutionPayloadHeader() (
		ExecutionPayloadHeaderT, error,
	)
	// ValidatorIndexByPubkey finds the validator index associated with a given
	// BLS public key.
	ValidatorIndexByPubkey(crypto.BLSPubkey) (math.ValidatorIndex, error)
	// GetBlockRootAtIndex retrieves the block root at a specified index.
	GetBlockRootAtIndex(uint64) (common.Root, error)
}

// ExecutionPayload is the interface for the execution payload.
type ExecutionPayload[T constraints.ForkTyped[T]] interface {
	constraints.ForkTyped[T]
	// GetBlockHash returns the block hash.
	GetBlockHash() gethprimitives.ExecutionHash
	// GetFeeRecipient returns the fee recipient.
	GetFeeRecipient() gethprimitives.ExecutionAddress
	// GetParentHash returns the parent hash.
	GetParentHash() gethprimitives.ExecutionHash
}

// ExecutionPayloadHeader is the interface for the execution payload header.
type ExecutionPayloadHeader interface {
	// GetBlockHash returns the block hash.
	GetBlockHash() gethprimitives.ExecutionHash
	// GetParentHash returns the parent hash.
	GetParentHash() gethprimitives.ExecutionHash
}

// PayloadAttributes is the interface for the payload attributes.
type PayloadAttributes[
	SelfT any,
	WithdrawalT any,
] interface {
	engineprimitives.PayloadAttributer
	// New creates a new payload attributes instance.
	New(
		uint32,
		uint64,
		common.Bytes32,
		gethprimitives.ExecutionAddress,
		[]WithdrawalT,
		common.Root,
	) (SelfT, error)
}

// ExecutionEngine is the interface for the execution engine.
type ExecutionEngine[
	ExecutionPayloadT, PayloadAttributesT any, PayloadIDT ~[8]byte,
] interface {
	// GetPayload returns the payload and blobs bundle for the given slot.
	GetPayload(
		ctx context.Context,
		req *engineprimitives.GetPayloadRequest[PayloadIDT],
	) (engineprimitives.BuiltExecutionPayloadEnv[ExecutionPayloadT], error)
	// NotifyForkchoiceUpdate notifies the execution client of a forkchoice
	// update.
	NotifyForkchoiceUpdate(
		ctx context.Context,
		req *engineprimitives.ForkchoiceUpdateRequest[PayloadAttributesT],
	) (*PayloadIDT, *gethprimitives.ExecutionHash, error)
}

// StateProcessor is the interface for the state processor.
type StateProcessor[BeaconStateT any, WithdrawalT any] interface {
	// ProcessState processes the state.
	ExpectedWithdrawals(
		BeaconStateT,
	) ([]WithdrawalT, error)
}
