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

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
)

// BeaconState defines the interface for accessing various state-related data
// required for block processing.
type BeaconState interface {
	// GetRandaoMixAtIndex retrieves the RANDAO mix at a specified index.
	GetRandaoMixAtIndex(uint64) (common.Bytes32, error)
	// ExpectedWithdrawals lists the expected withdrawals in the current state.
	ExpectedWithdrawals() (engineprimitives.Withdrawals, error)
	// GetLatestExecutionPayloadHeader fetches the most recent execution payload
	// header.
	GetLatestExecutionPayloadHeader() (*ctypes.ExecutionPayloadHeader, error)
	// ValidatorIndexByPubkey finds the validator index associated with a given
	// BLS public key.
	ValidatorIndexByPubkey(crypto.BLSPubkey) (math.ValidatorIndex, error)
	// GetBlockRootAtIndex retrieves the block root at a specified index.
	GetBlockRootAtIndex(uint64) (common.Root, error)
}

type PayloadCache[PayloadIDT, RootT, SlotT any] interface {
	Get(slot SlotT, stateRoot RootT) (PayloadIDT, bool)
	Has(slot SlotT, stateRoot RootT) bool
	Set(slot SlotT, stateRoot RootT, pid PayloadIDT)
	UnsafePrunePrior(slot SlotT)
}

// ExecutionPayload is the interface for the execution payload.
type ExecutionPayload[T any] interface {
	constraints.ForkTyped[T]
	// GetBlockHash returns the block hash.
	GetBlockHash() common.ExecutionHash
	// GetFeeRecipient returns the fee recipient.
	GetFeeRecipient() common.ExecutionAddress
	// GetParentHash returns the parent hash.
	GetParentHash() common.ExecutionHash
}

// AttributesFactory is the interface for the attributes factory.
type AttributesFactory[
	BeaconStateT any,
] interface {
	BuildPayloadAttributes(
		st BeaconStateT,
		slot math.U64,
		timestamp uint64,
		prevHeadRoot [32]byte,
	) (*engineprimitives.PayloadAttributes, error)
}

// PayloadAttributes is the interface for the payload attributes.
type PayloadAttributes[
	SelfT any,
] interface {
	engineprimitives.PayloadAttributer
	// New creates a new payload attributes instance.
	New(
		uint32,
		uint64,
		common.Bytes32,
		common.ExecutionAddress,
		engineprimitives.Withdrawals,
		common.Root,
	) (SelfT, error)
}

// ExecutionEngine is the interface for the execution engine.
type ExecutionEngine[
	PayloadIDT ~[8]byte,
] interface {
	// GetPayload returns the payload and blobs bundle for the given slot.
	GetPayload(
		ctx context.Context,
		req *ctypes.GetPayloadRequest[PayloadIDT],
	) (ctypes.BuiltExecutionPayloadEnv, error)
	// NotifyForkchoiceUpdate notifies the execution client of a forkchoice
	// update.
	NotifyForkchoiceUpdate(
		ctx context.Context,
		req *ctypes.ForkchoiceUpdateRequest,
	) (*PayloadIDT, *common.ExecutionHash, error)
}
