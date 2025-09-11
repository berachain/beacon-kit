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
	"encoding/json"
	"errors"
	"fmt"

	corestore "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/store"
	sdkmetrics "cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	"github.com/berachain/beacon-kit/chain"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	cometbft "github.com/berachain/beacon-kit/consensus/cometbft/service"
	"github.com/berachain/beacon-kit/node-core/components/metrics"
	"github.com/berachain/beacon-kit/node-core/components/storage"
	"github.com/berachain/beacon-kit/node-core/types"
	"github.com/berachain/beacon-kit/state-transition/core/state"
	kvstorage "github.com/berachain/beacon-kit/storage"
	"github.com/berachain/beacon-kit/storage/beacondb"
	"github.com/berachain/beacon-kit/storage/db"
	cmtcfg "github.com/cometbft/cometbft/config"
	dbm "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
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

func (b *Backend) checkChainIsReady() error {
	switch err := b.node.IsAppReady(); {
	case err == nil:
		// chain finally ready, time to loading genesis
		return b.loadGenesisState()
	case errors.Is(err, cometbft.ErrAppNotReady):
		return cometbft.ErrAppNotReady
	default:
		return fmt.Errorf("unable to check whether app is ready: %w", err)
	}
}

type backendKVStoreService struct {
	ctx sdk.Context
}

func (kvs *backendKVStoreService) OpenKVStore(context.Context) corestore.KVStore {
	//nolint:contextcheck // fine with tests
	store := sdk.UnwrapSDKContext(kvs.ctx).KVStore(backendStoreKey)
	return kvstorage.NewKVStore(store)
}

//nolint:gochecknoglobals // todo: fix later
var backendStoreKey = storetypes.NewKVStoreKey("backend-genesis")

func (b *Backend) loadGenesisState() error {
	if b.genesisState != nil {
		// genesis state already initialized, we're fine
		return nil
	}

	// 1- Load Genesis bytes
	genesisData, err := parseGenesisBytes(b)
	if err != nil {
		return err
	}

	// 2- Create Genesis Store
	b.genesisState, err = initGenesisState(b.cs)
	if err != nil {
		return fmt.Errorf("backend data loading:%w", err)
	}

	// 3- Reprocess Genesis via state Processor. This is safe
	// since it's done on its own state AND state processor does not
	// make any call to the EVM during genesis processing.
	// Note: we process genesis here as soon as node start, but
	// chain would wait for genesisTime to come if genesisTime
	// is set in the future. We replicate this behaviour with checkChainIsReady.
	if _, err = b.sp.InitializeBeaconStateFromEth1(
		b.genesisState,
		genesisData.GetDeposits(),
		genesisData.GetExecutionPayloadHeader(),
		genesisData.GetForkVersion(),
	); err != nil {
		return fmt.Errorf("failed processing genesis: %w", err)
	}
	return nil
}

func parseGenesisBytes(b *Backend) (ctypes.Genesis, error) {
	appGenesis, err := genutiltypes.AppGenesisFromFile(b.cmtCfg.GenesisFile())
	if err != nil {
		return ctypes.Genesis{}, fmt.Errorf("failed loading app genesis from file: %w", err)
	}
	gen, err := appGenesis.ToGenesisDoc()
	if err != nil {
		return ctypes.Genesis{}, fmt.Errorf("failed parsing app genesis: %w", err)
	}
	var genesisState map[string]json.RawMessage
	if err = json.Unmarshal(gen.AppState, &genesisState); err != nil {
		return ctypes.Genesis{}, fmt.Errorf("failed to unmarshal genesis state: %w", err)
	}
	data := []byte(genesisState["beacon"])

	genesisData := ctypes.Genesis{}
	if err = json.Unmarshal(data, &genesisData); err != nil {
		return ctypes.Genesis{}, fmt.Errorf("failed to unmarshal genesis data: %w", err)
	}
	return genesisData, nil
}

func initGenesisState(cs chain.Spec) (*state.StateDB, error) {
	db, err := db.OpenDB("", dbm.MemDBBackend)
	if err != nil {
		return nil, fmt.Errorf("failed opening mem db: %w", err)
	}
	var (
		nopLog     = log.NewNopLogger()
		nopMetrics = sdkmetrics.NewNoOpMetrics()
	)

	cms := store.NewCommitMultiStore(db, nopLog, nopMetrics)

	cms.MountStoreWithDB(backendStoreKey, storetypes.StoreTypeIAVL, nil)
	if err = cms.LoadLatestVersion(); err != nil {
		return nil, fmt.Errorf("backend data loading: failed to load latest version: %w", err)
	}

	ctx := sdk.NewContext(cms, true, nopLog)
	backendStoreService := &backendKVStoreService{
		ctx: ctx,
	}
	kvStore := beacondb.New(backendStoreService)

	sdkCtx := sdk.NewContext(cms.CacheMultiStore(), true, log.NewNopLogger())
	return state.NewBeaconStateFromDB(
		kvStore.WithContext(sdkCtx),
		cs,
		sdkCtx.Logger(),
		metrics.NewNoOpTelemetrySink(),
	), nil
}
