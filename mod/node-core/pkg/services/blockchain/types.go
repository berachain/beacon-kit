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

package blockchain

import (
	"context"

	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// BeaconBlock represents a beacon block interface.
type BeaconBlock[
	BeaconBlockBodyT BeaconBlockBody[ExecutionPayloadT],
	ExecutionPayloadT any,
] interface {
	constraints.SSZMarshallableRootable
	constraints.Nillable
	// GetSlot returns the slot of the beacon block.
	GetSlot() math.Slot
	// GetParentBlockRoot returns the parent block root of the beacon block.
	GetParentBlockRoot() common.Root
	// GetStateRoot returns the state root of the beacon block.
	GetStateRoot() common.Root
	// GetBody returns the body of the beacon block.
	GetBody() BeaconBlockBodyT
}

// BeaconBlockBody represents the interface for the beacon block body.
type BeaconBlockBody[ExecutionPayloadT any] interface {
	constraints.SSZMarshallableRootable
	constraints.Nillable
	// GetExecutionPayload returns the execution payload of the beacon block
	// body.
	GetExecutionPayload() ExecutionPayloadT
}

// EventFeed is a generic interface for sending events.
type EventFeed[EventT any] interface {
	// Publish sends an event and returns an error if any occurred.
	Publish(ctx context.Context, event EventT) error
	// Subscribe returns a channel that will receive events.
	Subscribe() (chan EventT, error)
}

// ExecutionPayload is the interface for the execution payload.
type ExecutionPayload interface {
	ExecutionPayloadHeader
}

// ExecutionPayloadHeader is the interface for the execution payload header.
type ExecutionPayloadHeader interface {
	// GetTimestamp returns the timestamp.
	GetTimestamp() math.U64
	// GetBlockHash returns the block hash.
	GetBlockHash() gethprimitives.ExecutionHash
	// GetParentHash returns the parent hash.
	GetParentHash() gethprimitives.ExecutionHash
}

// Genesis is the interface for the genesis.
type Genesis[DepositT any, ExecutionPayloadHeaderT any] interface {
	// GetForkVersion returns the fork version.
	GetForkVersion() common.Version
	// GetDeposits returns the deposits.
	GetDeposits() []DepositT
	// GetExecutionPayloadHeader returns the execution payload header.
	GetExecutionPayloadHeader() ExecutionPayloadHeaderT
}
