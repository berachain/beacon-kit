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
	"fmt"

	"github.com/berachain/beacon-kit/chain"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/node-core/components/storage"
	"github.com/berachain/beacon-kit/node-core/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
)

// Backend is the db access layer for the beacon node-api.
// It serves as a wrapper around the storage backend and provides an abstraction
// over building the query context for a given state.
type Backend struct {
	sb   *storage.Backend
	cs   chain.Spec
	node types.ConsensusService
	sp   StateProcessor
}

// New creates and returns a new Backend instance.
func New(
	storageBackend *storage.Backend,
	cs chain.Spec,
	sp StateProcessor,
) *Backend {
	return &Backend{
		sb: storageBackend,
		cs: cs,
		sp: sp,
	}
}

// AttachQueryBackend sets the node on the backend for
// querying historical heights.
func (b *Backend) AttachQueryBackend(node types.ConsensusService) {
	b.node = node
}

// GetSlotByBlockRoot retrieves the slot by a block root from the block store.
func (b *Backend) GetSlotByBlockRoot(root common.Root) (math.Slot, error) {
	return b.sb.BlockStore().GetSlotByBlockRoot(root)
}

// GetSlotByStateRoot retrieves the slot by a state root from the block store.
func (b *Backend) GetSlotByStateRoot(root common.Root) (math.Slot, error) {
	return b.sb.BlockStore().GetSlotByStateRoot(root)
}

// GetParentSlotByTimestamp retrieves the parent slot by a given timestamp from
// the block store.
func (b *Backend) GetParentSlotByTimestamp(timestamp math.U64) (math.Slot, error) {
	return b.sb.BlockStore().GetParentSlotByTimestamp(timestamp)
}

// Spec returns the chain spec used by the backend.
func (b *Backend) Spec() (chain.Spec, error) {
	if b.cs == nil {
		return nil, errors.New("chain spec not found")
	}
	return b.cs, nil
}

// stateFromSlot returns the state at the given slot, after also processing the
// next slot to ensure the returned beacon state is up to date.
func (b *Backend) stateFromSlot(slot math.Slot) (*statedb.StateDB, math.Slot, error) {
	st, slot, err := b.stateFromSlotRaw(slot)
	if err != nil {
		return st, slot, fmt.Errorf("stateFromSlotRaw failed: %w", err)
	}

	// Process the slot to update the latest state and block roots.
	targetSlot := slot + 1
	if _, err = b.sp.ProcessSlots(st, targetSlot); err != nil {
		return st, slot, fmt.Errorf("ProcessSlots failed, target slot %d: %w", targetSlot, err)
	}

	// We need to set the slot on the state back since ProcessSlot will update
	// it to slot + 1.
	if err = st.SetSlot(slot); err != nil {
		return st, slot, fmt.Errorf("failed resetting slot to %d: %w", slot, err)
	}
	return st, slot, nil
}

// stateFromSlotRaw returns the state at the given slot using query context,
// resolving an input slot of 0 to the latest slot. It does not process the
// next slot on the beacon state.
func (b *Backend) stateFromSlotRaw(slot math.Slot) (*statedb.StateDB, math.Slot, error) {
	queryCtx, err := b.node.CreateQueryContext(int64(slot), false) // #nosec G115 -- not an issue in practice.
	if err != nil {
		return nil, slot, fmt.Errorf("CreateQueryContext failed: %w", err)
	}
	st := b.sb.StateFromContext(queryCtx)

	// If using height 0 for the query context, make sure to return the latest slot.
	if slot == 0 {
		slot, err = st.GetSlot()
		if err != nil {
			return st, slot, fmt.Errorf("GetSlot failed: %w", err)
		}
	}
	return st, slot, nil
}
