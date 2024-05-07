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

	"github.com/berachain/beacon-kit/mod/beacon/blockchain"
	"github.com/berachain/beacon-kit/mod/beacon/validator"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/consensus"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/abci"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/service"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeaconKitRuntime is a struct that holds the
// service registry.
type BeaconKitRuntime[
	BlobSidecarsT BlobSidecars,
	DepositStoreT DepositStore,
	ReadOnlyBeaconBlockBodyT consensus.ReadOnlyBeaconBlockBody,
	StorageBackendT BeaconStorageBackend[
		BlobSidecarsT,
		DepositStoreT,
		ReadOnlyBeaconBlockBodyT,
	],
] struct {
	logger   log.Logger[any]
	services *service.Registry
	fscp     StorageBackendT
}

// NewBeaconKitRuntime creates a new BeaconKitRuntime
// and applies the provided options.
func NewBeaconKitRuntime[
	BlobSidecarsT BlobSidecars,
	DepositStoreT DepositStore,
	ReadOnlyBeaconBlockBodyT consensus.ReadOnlyBeaconBlockBody,
	StorageBackendT BeaconStorageBackend[
		BlobSidecarsT,
		DepositStoreT,
		ReadOnlyBeaconBlockBodyT,
	],
](
	logger log.Logger[any],
	services *service.Registry,
	fscp StorageBackendT,
) (*BeaconKitRuntime[
	BlobSidecarsT,
	DepositStoreT,
	ReadOnlyBeaconBlockBodyT,
	StorageBackendT,
], error) {
	return &BeaconKitRuntime[
		BlobSidecarsT,
		DepositStoreT,
		ReadOnlyBeaconBlockBodyT,
		StorageBackendT,
	]{
		logger:   logger,
		services: services,
		fscp:     fscp,
	}, nil
}

// StartServices starts the services.
func (r *BeaconKitRuntime[
	BlobSidecarsT,
	DepositStoreT,
	ReadOnlyBeaconBlockBodyT,
	StorageBackendT,
]) StartServices(
	ctx context.Context,
) {
	r.services.StartAll(ctx)
}

// BuildABCIComponents returns the ABCI components for the beacon runtime.
func (r *BeaconKitRuntime[
	BlobSidecarsT,
	DepositStoreT,
	ReadOnlyBeaconBlockBodyT,
	StorageBackendT,
]) BuildABCIComponents() (
	sdk.PrepareProposalHandler, sdk.ProcessProposalHandler,
	sdk.PreBlocker,
) {
	var (
		chainService   *blockchain.Service[BlobSidecarsT]
		builderService *validator.Service[BlobSidecarsT]
	)
	if err := r.services.FetchService(&chainService); err != nil {
		panic(err)
	}

	if err := r.services.FetchService(&builderService); err != nil {
		panic(err)
	}

	if chainService == nil || builderService == nil {
		panic("missing services")
	}

	handler := abci.NewHandler(
		builderService,
		chainService,
	)

	return handler.PrepareProposalHandler,
		handler.ProcessProposalHandler,
		handler.FinalizeBlock
}
