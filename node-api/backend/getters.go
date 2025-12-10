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
	"errors"
	"fmt"
	"runtime"

	"cosmossdk.io/log"
	cometbft "github.com/berachain/beacon-kit/consensus/cometbft/service"
	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
)

// StateAndSlotFromHeight returns the beacon state at a particular slot using query context,
// resolving an input height of -1 to the latest slot.
//
// This returns the beacon state of the version that was committed to disk at the requested slot,
// which has the empty state root in the latest block header. Hence, the most recent state and
// block roots are not updated.
func (b *Backend) StateAndSlotFromHeight(height int64) (ReadOnlyBeaconState, math.Slot, error) {
	if height < -1 {
		return nil, 0, fmt.Errorf("expected height, must be non-negative or -1 to request tip, got %d", height)
	}
	if height == 0 {
		switch err := b.node.IsAppReady(); {
		case err == nil:
			// chain finally ready, time to loading genesis
			if err = b.loadGenesisState(); err != nil {
				return nil, 0, fmt.Errorf("failed loading genesis state: %w", err)
			}
		case errors.Is(err, cometbft.ErrAppNotReady):
			return nil, 0, cometbft.ErrAppNotReady
		default:
			return nil, 0, fmt.Errorf("unable to check whether app is ready: %w", err)
		}

		b.muSt.Lock()
		defer b.muSt.Unlock()

		// Copy the state to ensure clients potential changes won't pollute the state
		// Also we make sure to create the copy in a thread-safe way via the muCms mutex.
		ms := b.cms.CacheMultiStore()
		ctx := sdk.NewContext(ms, true, log.NewNopLogger())
		ephemeralGenesisState := b.genesisState.Protect(ctx)
		return ephemeralGenesisState, 0, nil
	}

	height = max(0, height) // CreateQueryContext uses 0 to pick latest height.
	queryCtx, err := b.node.CreateQueryContext(height, false)
	if err != nil {
		return nil, 0, fmt.Errorf("CreateQueryContext failed: %w", err)
	}
	st := b.sb.StateFromContext(queryCtx)

	var slot math.Slot
	if height > 0 {
		slot = math.Slot(height)
	} else {
		// height must be -1, so pick state slot
		slot, err = st.GetSlot()
		if err != nil {
			return st, slot, fmt.Errorf("GetSlot failed: %w", err)
		}
	}
	return st, slot, nil
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

// GetSignatureBySlot retrieves the block signature for a given slot from the block store.
func (b *Backend) GetSignatureBySlot(slot math.Slot) (crypto.BLSSignature, error) {
	return b.sb.BlockStore().GetSignatureBySlot(slot)
}

func (b *Backend) GetBlobSidecarsAtSlot(slot math.Slot) (datypes.BlobSidecars, error) {
	return b.sb.AvailabilityStore().GetBlobSidecars(slot)
}

func (b *Backend) GetSyncData() (int64 /*latestHeight*/, int64 /*syncToHeight*/) {
	return b.node.GetSyncData()
}

func (b *Backend) GetVersionData() (
	string, // appName
	string, // cometVersion
	string, // os
	string, // arch
) {
	cometVersionInfo := version.NewInfo() // same used in beacond version command

	var (
		appName      = cometVersionInfo.AppName
		cometVersion = cometVersionInfo.Version
		os           = runtime.GOOS
		arch         = runtime.GOARCH
	)

	return appName,
		cometVersion,
		os,
		arch
}
