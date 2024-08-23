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

package cometbft

import (
	"context"
	"errors"
	"fmt"

	pruningtypes "cosmossdk.io/store/pruning/types"
	storetypes "cosmossdk.io/store/types"
)

// File for storing in-package cometbft optional functions,
// for options that need access to non-exported fields of the Service

// SetPruning sets a pruning option on the multistore associated with the app.
func SetPruning(opts pruningtypes.PruningOptions) func(*Service) {
	return func(bapp *Service) { bapp.cms.SetPruning(opts) }
}

// SetMinRetainBlocks returns a Service option function that sets the minimum
// block retention height value when determining which heights to prune during
// ABCI Commit.
func SetMinRetainBlocks(minRetainBlocks uint64) func(*Service) {
	return func(bapp *Service) { bapp.setMinRetainBlocks(minRetainBlocks) }
}

// SetIAVLCacheSize provides a Service option function that sets the size of
// IAVL cache.
func SetIAVLCacheSize(size int) func(*Service) {
	return func(bapp *Service) { bapp.cms.SetIAVLCacheSize(size) }
}

// SetIAVLDisableFastNode enables(false)/disables(true) fast node usage from the
// IAVL store.
func SetIAVLDisableFastNode(disable bool) func(*Service) {
	return func(bapp *Service) { bapp.cms.SetIAVLDisableFastNode(disable) }
}

// SetInterBlockCache provides a Service option function that sets the
// inter-block cache.
func SetInterBlockCache(
	cache storetypes.MultiStorePersistentCache,
) func(*Service) {
	return func(app *Service) { app.setInterBlockCache(cache) }
}

// SetChainID sets the chain ID in cometbft.
func SetChainID(chainID string) func(*Service) {
	return func(app *Service) { app.chainID = chainID }
}

func (app *Service) SetName(name string) {
	app.name = name
}

// SetVersion sets the application's version string.
func (app *Service) SetVersion(v string) {
	app.version = v
}

func (app *Service) SetAppVersion(ctx context.Context, v uint64) error {
	if app.paramStore == nil {
		return errors.
			New("param store must be set to set app version")
	}

	cp, err := app.paramStore.Get(ctx)
	if err != nil {
		return fmt.
			Errorf("failed to get consensus params: %w", err)
	}
	if cp.Version == nil {
		return errors.
			New("version is not set in param store")
	}
	cp.Version.App = v
	return app.paramStore.Set(ctx, cp)
}
