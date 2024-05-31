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
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/runtime/middleware"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/service"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
)

// BeaconKitRuntime is a struct that holds the
// service registry.
type BeaconKitRuntime[
	AvailabilityStoreT AvailabilityStore[types.BeaconBlockBody, BlobSidecarsT],
	BeaconBlockBodyT types.BeaconBlockBody,
	BeaconStateT core.BeaconState[*types.Validator],
	BlobSidecarsT BlobSidecars,
	DepositStoreT DepositStore,
	StorageBackendT StorageBackend[
		AvailabilityStoreT, BeaconBlockBodyT,
		BeaconStateT, BlobSidecarsT, DepositStoreT,
	],
] struct {
	// logger is used for logging within the BeaconKitRuntime.
	logger log.Logger[any]
	// services is a registry of services used by the BeaconKitRuntime.
	services *service.Registry
	// storageBackend is the backend storage interface used by the
	// BeaconKitRuntime.
	storageBackend StorageBackendT
	// chainSpec defines the chain specifications for the BeaconKitRuntime.
	chainSpec primitives.ChainSpec
	// abciFinalizeBlockMiddleware handles ABCI interactions for the
	// BeaconKitRuntime.
	abciFinalizeBlockMiddleware *middleware.FinalizeBlockMiddleware[
		types.BeaconBlock, BeaconStateT, BlobSidecarsT,
	]
	// abciValidatorMiddleware is responsible for forward ABCI requests to the
	// validator service.
	abciValidatorMiddleware *middleware.ValidatorMiddleware[
		BeaconStateT, BlobSidecarsT,
	]
}

// NewBeaconKitRuntime creates a new BeaconKitRuntime
// and applies the provided options.
func NewBeaconKitRuntime[
	AvailabilityStoreT AvailabilityStore[types.BeaconBlockBody, BlobSidecarsT],
	BeaconBlockBodyT types.BeaconBlockBody,
	BeaconStateT core.BeaconState[*types.Validator],
	BlobSidecarsT BlobSidecars,
	DepositStoreT DepositStore,
	StorageBackendT StorageBackend[
		AvailabilityStoreT, BeaconBlockBodyT,
		BeaconStateT, BlobSidecarsT,
		DepositStoreT,
	],
](
	chainSpec primitives.ChainSpec,
	logger log.Logger[any],
	services *service.Registry,
	storageBackend StorageBackendT,
	telemetrySink middleware.TelemetrySink,
) (*BeaconKitRuntime[
	AvailabilityStoreT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, DepositStoreT, StorageBackendT,
], error) {
	var (
		chainService *blockchain.Service[
			AvailabilityStoreT,
			core.BeaconState[*types.Validator],
			BlobSidecarsT,
			DepositStoreT,
		]
		validatorService *validator.Service[
			types.BeaconBlock,
			types.BeaconBlockBody,
			core.BeaconState[*types.Validator],
			BlobSidecarsT,
		]
	)

	if err := services.FetchService(&chainService); err != nil {
		panic(err)
	}

	if err := services.FetchService(&validatorService); err != nil {
		panic(err)
	}

	return &BeaconKitRuntime[
		AvailabilityStoreT, BeaconBlockBodyT, BeaconStateT,
		BlobSidecarsT, DepositStoreT, StorageBackendT,
	]{
		abciFinalizeBlockMiddleware: middleware.
			NewFinalizeBlockMiddleware[
			types.BeaconBlock, BeaconStateT,
		](
			chainSpec,
			chainService,
			telemetrySink,
		),
		abciValidatorMiddleware: middleware.
			NewValidatorMiddleware[
			BeaconStateT, BlobSidecarsT,
		](
			chainSpec,
			validatorService,
			telemetrySink,
		),
		chainSpec:      chainSpec,
		logger:         logger,
		services:       services,
		storageBackend: storageBackend,
	}, nil
}

// StartServices starts the services.
func (r *BeaconKitRuntime[
	AvailabilityStoreT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, DepositStoreT, StorageBackendT,
]) StartServices(
	ctx context.Context,
) error {
	return r.services.StartAll(ctx)
}

// ABCIHandler returns the ABCI handler.
func (r *BeaconKitRuntime[
	AvailabilityStoreT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, DepositStoreT, StorageBackendT,
]) ABCIFinalizeBlockMiddleware() *middleware.FinalizeBlockMiddleware[
	types.BeaconBlock, BeaconStateT, BlobSidecarsT,
] {
	return r.abciFinalizeBlockMiddleware
}

// ABCIValidatorMiddleware returns the ABCI validator middleware.
func (r *BeaconKitRuntime[
	AvailabilityStoreT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, DepositStoreT, StorageBackendT,
]) ABCIValidatorMiddleware() *middleware.ValidatorMiddleware[
	BeaconStateT, BlobSidecarsT,
] {
	return r.abciValidatorMiddleware
}
