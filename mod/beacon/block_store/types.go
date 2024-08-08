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

package blockstore

import (
	asynctypes "github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// BeaconBlock is a generic interface for a beacon block.
type BeaconBlock interface {
	constraints.SSZMarshaler
	// GetSlot returns the slot of the block.
	GetSlot() math.U64
}

// BlockStore is a generic interface for a block store.
type BlockStore[BeaconBlockT BeaconBlock] interface {
	// Set sets a block at a given index.
	Set(index math.Slot, blk BeaconBlockT) error
}

// Event is an interface for block events.
type Event[BeaconBlockT BeaconBlock] interface {
	// Type returns the type of the event.
	Type() asynctypes.EventID
	// Is returns true if the event is of the given type.
	Is(asynctypes.EventID) bool
	// Data returns the data of the event.
	Data() BeaconBlockT
}

// EventFeed is a generic interface for sending events.
type EventFeed[EventT any] interface {
	// Subscribe returns a channel that will receive events.
	Subscribe() (chan EventT, error)
}
