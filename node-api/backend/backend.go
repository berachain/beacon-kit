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
	"sync/atomic"

	"github.com/berachain/beacon-kit/chain"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/node-core/components/storage"
	"github.com/berachain/beacon-kit/node-core/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	cmtcfg "github.com/cometbft/cometbft/config"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
)

// Backend is the db access layer for the beacon node-api.
// It serves as a wrapper around the storage backend and provides an abstraction
// over building the query context for a given state.
type Backend struct {
	sb   *storage.Backend
	cs   chain.Spec
	node types.ConsensusService

	// genesisValidatorsRoot is cached in the backend.
	genesisValidatorsRoot atomic.Pointer[common.Root]

	// genesisTime is cached here, written to once during initialization!
	genesisTime atomic.Pointer[math.U64]

	// genesisForkVersion is cached here, written to once during initialization!
	genesisForkVersion atomic.Pointer[common.Version]
}

// New creates and returns a new Backend instance.
func New(
	storageBackend *storage.Backend,
	cs chain.Spec,
	cmtCfg *cmtcfg.Config,
) (*Backend, error) {
	b := &Backend{
		sb: storageBackend,
		cs: cs,
	}

	// Load the genesis file from cometbft config.
	appGenesis, err := genutiltypes.AppGenesisFromFile(cmtCfg.GenesisFile())
	if err != nil {
		return nil, err
	}
	gen, err := appGenesis.ToGenesisDoc()
	if err != nil {
		return nil, err
	}

	// Store the genesis time in the backend.
	//#nosec: G115 // Unix time will never be negative.
	genesisTime := math.U64(gen.GenesisTime.Unix())
	b.genesisTime.Store(&genesisTime)

	// Derive the genesis fork version from the genesis time.
	genesisForkVersion := cs.ActiveForkVersionForTimestamp(genesisTime)
	b.genesisForkVersion.Store(&genesisForkVersion)

	return b, nil
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
