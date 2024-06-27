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

package middleware

import (
	"context"
	"encoding/json"
	"time"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
)

// BeaconBlock is an interface for accessing the beacon block.
type BeaconBlock[T any] interface {
	constraints.SSZMarshallable
	constraints.Nillable
	GetSlot() math.Slot
	NewFromSSZ([]byte, uint32) (T, error)
}

// BeaconState is an interface for accessing the beacon state.
type BeaconState interface {
	// ValidatorIndexByPubkey returns the validator index for the given pubkey.
	ValidatorIndexByPubkey(
		pubkey crypto.BLSPubkey,
	) (math.ValidatorIndex, error)
	// GetBlockRootAtIndex returns the block root at the given index.
	GetBlockRootAtIndex(
		index uint64,
	) (common.Root, error)
	// ValidatorIndexByCometBFTAddress returns the validator index for the given
	ValidatorIndexByCometBFTAddress(
		cometBFTAddress []byte,
	) (math.ValidatorIndex, error)
}

// BlockchainService defines the interface for interacting with the blockchain
// state and processing blocks.
type BlockchainService[
	BeaconBlockT any,
	BlobSidecarsT constraints.SSZMarshallable,
	DepositT any,
	GenesisT Genesis,
] interface {
	// ProcessGenesisData processes the genesis data and initializes the beacon
	// state.
	ProcessGenesisData(
		context.Context,
		GenesisT,
	) (transition.ValidatorUpdates, error)
	// ProcessBeaconBlock processes the given beacon block and associated
	// blobs sidecars.
	ProcessBeaconBlock(
		context.Context,
		BeaconBlockT,
	) (transition.ValidatorUpdates, error)
	// ReceiveBlock receives a beacon block and
	// associated blobs sidecars for processing.
	ReceiveBlock(
		ctx context.Context,
		blk BeaconBlockT,
	) error
}

// DAService.
type DAService[
	BlobSidecarsT any,
] interface {
	// ProcessSidecars
	ProcessSidecars(
		context.Context,
		BlobSidecarsT,
	) error
	// ReceiveSidecars
	ReceiveSidecars(
		_ context.Context,
		sidecars BlobSidecarsT,
	) error
}

// ExecutionPayloadHeader is the interface for the execution data of a block.
type ExecutionPayloadHeader[T any] interface {
	NewFromJSON([]byte, uint32) (T, error)
}

// Genesis is the interface for the genesis data.
type Genesis interface {
	json.Unmarshaler
}

// TelemetrySink is an interface for sending metrics to a telemetry backend.
type TelemetrySink interface {
	// MeasureSince measures the time since the given time.
	MeasureSince(key string, start time.Time, args ...string)
}

// StorageBackend is an interface for accessing the storage backend.
type StorageBackend[BeaconStateT any] interface {
	StateFromContext(ctx context.Context) BeaconStateT
}
