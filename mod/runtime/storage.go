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

	"cosmossdk.io/core/appmodule"
	appmodulev2 "cosmossdk.io/core/appmodule/v2"
	"github.com/berachain/beacon-kit/mod/core"
	"github.com/berachain/beacon-kit/mod/core/state"
	"github.com/berachain/beacon-kit/mod/core/state/deneb"
	dastore "github.com/berachain/beacon-kit/mod/da/pkg/store"
	"github.com/berachain/beacon-kit/mod/primitives"
	engineprimitives "github.com/berachain/beacon-kit/mod/primitives-engine"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/consensus"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/storage"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb"
	"github.com/berachain/beacon-kit/mod/storage/pkg/deposit"
	"github.com/berachain/beacon-kit/mod/storage/pkg/filedb"
	bls12381 "github.com/cosmos/cosmos-sdk/crypto/keys/bls12_381"
)

// BeaconStorageBackend is an interface that provides the
// beacon state to the runtime.
type BeaconStorageBackend[ReadOnlyBeaconBlockT, BlobSidecarsT any] interface {
	AvailabilityStore(
		ctx context.Context,
	) core.AvailabilityStore[ReadOnlyBeaconBlockT, BlobSidecarsT]
	BeaconState(ctx context.Context) state.BeaconState
	DepositStore(ctx context.Context) *deposit.KVStore
}

// Keeper maintains the link to data storage and exposes access to the
// underlying `BeaconState` methods for the x/beacon module.
//
// TODO: The keeper will eventually be dissolved.
type Keeper struct {
	cs primitives.ChainSpec
	storage.Backend
}

// TODO: move this.
func DenebPayloadFactory() engineprimitives.ExecutionPayloadHeader {
	return &engineprimitives.ExecutionPayloadHeaderDeneb{}
}

// NewKeeper creates new instances of the Beacon Keeper.
func NewKeeper(
	fdb *filedb.DB,
	env appmodule.Environment,
	cs primitives.ChainSpec,
	ddb *deposit.KVStore,
) *Keeper {
	return &Keeper{
		cs: cs,
		Backend: *storage.NewBackend(cs, dastore.New[consensus.ReadOnlyBeaconBlock](
			cs, filedb.NewRangeDB(fdb),
		), beacondb.New[
			*consensus.Fork,
			*consensus.BeaconBlockHeader,
			engineprimitives.ExecutionPayloadHeader,
			*consensus.Eth1Data,
			*consensus.Validator,
		](env.KVStoreService, DenebPayloadFactory), ddb),
	}
}

// ApplyAndReturnValidatorSetUpdates returns the validator set updates from
// the beacon state.
//
// TODO: this function is horribly inefficient and should be replaced with a
// more efficient implementation, that does not update the entire
// valset every block.
func (k *Keeper) ApplyAndReturnValidatorSetUpdates(
	ctx context.Context,
) ([]appmodulev2.ValidatorUpdate, error) {
	store := k.BeaconStore().WithContext(ctx)
	// Get the public key of the validator
	val, err := store.GetValidatorsByEffectiveBalance()
	if err != nil {
		panic(err)
	}

	validatorUpdates := make([]appmodulev2.ValidatorUpdate, 0)
	for _, validator := range val {
		// TODO: Config
		// Max 100 validators in the active set.
		// TODO: this is kinda hood.
		if validator.EffectiveBalance == 0 {
			var idx math.ValidatorIndex
			idx, err = store.WithContext(ctx).
				ValidatorIndexByPubkey(validator.Pubkey)
			if err != nil {
				return nil, err
			}
			if err = store.WithContext(ctx).
				RemoveValidatorAtIndex(idx); err != nil {
				return nil, err
			}
		}

		// TODO: this works, but there is a bug where if we send a validator to
		// 0 voting power, it can somehow still propose the next block? This
		// feels big bad.
		validatorUpdates = append(validatorUpdates, appmodulev2.ValidatorUpdate{
			PubKey:     validator.Pubkey[:],
			PubKeyType: (&bls12381.PubKey{}).Type(),
			//#nosec:G701 // will not realistically cause a problem.
			Power: int64(validator.EffectiveBalance),
		})
	}

	// Save the store.
	store.Save()
	return validatorUpdates, nil
}

// InitGenesis initializes the genesis state of the module.
func (k *Keeper) InitGenesis(
	ctx context.Context,
	data *deneb.BeaconState,
) ([]appmodulev2.ValidatorUpdate, error) {
	// Load the store.
	store := k.BeaconStore().WithContext(ctx)
	sdb := state.NewBeaconStateFromDB(store, k.cs)
	if err := sdb.WriteGenesisStateDeneb(data); err != nil {
		return nil, err
	}

	// Build ValidatorUpdates for CometBFT.
	validatorUpdates := make([]appmodulev2.ValidatorUpdate, 0)
	blsType := (&bls12381.PubKey{}).Type()
	for _, validator := range data.Validators {
		validatorUpdates = append(validatorUpdates, appmodulev2.ValidatorUpdate{
			PubKey:     validator.Pubkey[:],
			PubKeyType: blsType,
			//#nosec:G701 // will not realistically cause a problem.
			Power: int64(validator.EffectiveBalance),
		})
	}
	return validatorUpdates, nil
}

// ExportGenesis exports the current state of the module as genesis state.
func (k *Keeper) ExportGenesis(_ context.Context) *deneb.BeaconState {
	return &deneb.BeaconState{}
}
