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

package app

import (
	"context"
	"fmt"
	"os"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/light/mod/core/state"
	"github.com/berachain/beacon-kit/light/mod/provider"
	"github.com/berachain/beacon-kit/light/mod/runtime"
	"github.com/berachain/beacon-kit/light/mod/storage/beacondb"
	"github.com/berachain/beacon-kit/mod/core"
	beaconstate "github.com/berachain/beacon-kit/mod/core/state"
	"github.com/berachain/beacon-kit/mod/da"
	"github.com/berachain/beacon-kit/mod/node-builder/utils/jwt"
	filedb "github.com/berachain/beacon-kit/mod/storage/filedb"
)

// RunLightNode starts the light node with the given configuration.
// In the future, this will build the light node runtime and start it.
func RunLightNode(ctx context.Context, config *Config) {
	// builds the runtime and starts the node
	client := provider.New(*config.Provider)
	if err := client.Start(); err != nil {
		panic(err)
	}
	// if err := client.Start(); err != nil {
	// 	panic(err)
	// }
	// time.Sleep(3 * time.Second)

	kvStore := beacondb.New(client).WithContext(ctx)

	// for {
	// 	payload, err := kvStore.GetLatestExecutionPayload()
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	fmt.Println("block hash", payload.GetBlockHash())
	// 	time.Sleep(5 * time.Second)
	// }

	fdb := filedb.NewDB(
		filedb.WithRootDirectory(config.Comet.Directory),
		filedb.WithDirectoryPermissions(os.ModePerm),
	)
	availabilityStore := da.NewStore(nil, fdb)

	backend := NewLightStorageBackend(kvStore, availabilityStore)

	secret, err := jwt.LoadFromFile(config.Beacon.Engine.JWTSecretPath)
	if err != nil {
		panic(err)
	}

	// create the runtime
	rt, err := runtime.NewDefaultBeaconLightRuntime(
		config.Beacon,
		secret,
		backend,
		client,
		log.NewLogger(os.Stdout),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("starting runtime")

	go rt.StartServices(ctx)

	<-ctx.Done()
	// subscribe to light block
	// ch, err := client.SubscribeToLightBlock(ctx)
	// if err != nil {
	// 	panic(err)
	// }

	// for {
	// 	select {
	// 	case block := <-ch:
	// 		// process the light block
	// 		log.Println(block)
	// 		log.Println("trustedEth1Hash", client.GetEth1BlockHash(ctx))
	// 		log.Println("latestBlockHeader", client.GetLatestBlockHeader(ctx))
	// 		log.Println("blockRootAtIndex", client.GetBlockRootAtIndex(ctx, 0))
	// 	case <-ctx.Done():
	// 		// context cancelled, exit the loop
	// 		return
	// 	}
	// }
}

type LightStorageBackend struct {
	// BeaconDB is the backend for the beacon chain.
	beaconDB *beacondb.KVStore

	// AvailabilityDB is the backend for the availability store.
	availabilityDB *da.Store
}

func NewLightStorageBackend(
	beaconDB *beacondb.KVStore,
	availabilityDB *da.Store,
) *LightStorageBackend {
	return &LightStorageBackend{
		beaconDB:       beaconDB,
		availabilityDB: availabilityDB,
	}
}

func (l *LightStorageBackend) AvailabilityStore(ctx context.Context) core.AvailabilityStore {
	return l.availabilityDB
}

func (l *LightStorageBackend) BeaconState(ctx context.Context) beaconstate.BeaconState {
	return state.NewBeaconStateFromDB(l.beaconDB)
}
