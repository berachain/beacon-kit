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
	"errors"
	"fmt"
	"runtime"
	"sync"

	"github.com/berachain/beacon-kit/chain"
	cometbft "github.com/berachain/beacon-kit/consensus/cometbft/service"
	"github.com/berachain/beacon-kit/node-core/components/storage"
	"github.com/berachain/beacon-kit/node-core/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/state-transition/core/state"
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cosmos/cosmos-sdk/version"
)

// Backend is the db access layer for the beacon node-api.
// It serves as a wrapper around the storage backend and provides an abstraction
// over building the query context for a given state.
type Backend struct {
	sb     *storage.Backend
	cs     chain.Spec
	cmtCfg *cmtcfg.Config // used to fetch genesis data upon LoadData
	node   types.ConsensusService

	// Genesis related data
	sp           GenesisStateProcessor // only needed to recreate genesis state upon API loading
	muGs         sync.RWMutex          // muGs protects genesisState for concurrent access.
	genesisState *state.StateDB        // caches genesis data to serve API requests
}

// New creates and returns a new Backend instance.
func New(
	storageBackend *storage.Backend,
	sp GenesisStateProcessor,
	cs chain.Spec,
	cmtCfg *cmtcfg.Config,
	consensusService types.ConsensusService,
) *Backend {
	b := &Backend{
		sb:     storageBackend,
		sp:     sp,
		cs:     cs,
		cmtCfg: cmtCfg,
		node:   consensusService,
	}

	// genesis data will be cached in LoadData
	return b
}

// Backend currently calculates and caches some genesis data. These data
// are only needed if the API is active, so their processing happens in `LoadData`
// which should be called only if node-api server is actually started (it would be
// configure to not start).
func (b *Backend) LoadData(_ context.Context) error {
	switch err := b.node.IsAppReady(); {
	case err == nil:
		// chain finally ready, time to loading genesis
		//nolint:contextcheck // loadGenesisState creates its own context for in-memory store
		return b.loadGenesisState()
	case errors.Is(err, cometbft.ErrAppNotReady):
		// start anyhow, we'll init genesis state later on
		return nil
	default:
		return fmt.Errorf("unable to check whether app is ready: %w", err)
	}
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
