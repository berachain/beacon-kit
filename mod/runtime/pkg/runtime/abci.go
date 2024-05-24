// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package runtime

import (
	"context"
	"encoding/json"

	appmodulev2 "cosmossdk.io/core/appmodule/v2"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/state/deneb"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"github.com/sourcegraph/conc/iter"
)

// TODO: InitGenesis should be calling into the StateProcessor.
func (r BeaconKitRuntime[
	AvailabilityStoreT, BeaconBlockBodyT,
	BeaconStateT, BlobSidecarsT,
	DepositStoreT, StorageBackendT],
) InitGenesis(
	ctx context.Context,
	bz json.RawMessage,
) ([]appmodulev2.ValidatorUpdate, error) {
	data := new(deneb.BeaconState)
	if err := json.Unmarshal(bz, data); err != nil {
		return nil, err
	}

	// Load the store.
	store := r.storageBackend.StateFromContext(ctx)
	if err := store.WriteGenesisStateDeneb(data); err != nil {
		return nil, err
	}

	// Build ValidatorUpdates for CometBFT.
	updates := make([]appmodulev2.ValidatorUpdate, 0)
	for _, validator := range data.Validators {
		updates = append(updates, appmodulev2.ValidatorUpdate{
			PubKey:     validator.Pubkey[:],
			PubKeyType: crypto.CometBLSType,
			Power:      crypto.CometBLSPower,
		},
		)
	}
	return updates, nil
}

// EndBlock returns the validator set updates from the beacon state.
func (r BeaconKitRuntime[
	AvailabilityStoreT, BeaconBlockBodyT,
	BeaconStateT, BlobSidecarsT, DepositStoreT,
	StorageBackendT,
]) EndBlock(
	ctx context.Context,
) ([]appmodulev2.ValidatorUpdate, error) {
	// Process the state transition and produce the required delta from
	// the sync committee.
	updates, err := r.chainService.ProcessStateTransition(
		ctx,
		// TODO: improve the robustness of these types to ensure we
		// don't run into any nil ptr issues.
		r.abciHandler.LatestBeaconBlock,
		r.abciHandler.LatestSidecars,
	)
	if err != nil {
		return nil, err
	}

	// Convert the delta into the appmodule ValidatorUpdate format to
	// pass onto CometBFT.
	return iter.MapErr(
		updates,
		func(
			u **transition.ValidatorUpdate,
		) (appmodulev2.ValidatorUpdate, error) {
			update := *u
			return appmodulev2.ValidatorUpdate{
				PubKey:     update.Pubkey[:],
				PubKeyType: crypto.CometBLSType,
				//#nosec:G701 // this is safe.
				Power: int64(update.EffectiveBalance.Unwrap()),
			}, nil
		},
	)
}
