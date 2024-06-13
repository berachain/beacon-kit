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

package runtime

import (
	"context"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/service"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
)

// BeaconKitRuntime is a struct that holds the
// service registry.
type BeaconKitRuntime[
	AvailabilityStoreT AvailabilityStore[BeaconBlockBodyT, BlobSidecarsT],
	BeaconBlockT interface {
		types.RawBeaconBlock[BeaconBlockBodyT]
		NewFromSSZ([]byte, uint32) (BeaconBlockT, error)
		NewWithVersion(
			math.Slot,
			math.ValidatorIndex,
			primitives.Root,
			uint32,
		) (BeaconBlockT, error)
		Empty(uint32) BeaconBlockT
	},
	BeaconBlockBodyT types.RawBeaconBlockBody,
	BeaconStateT core.BeaconState[
		*types.BeaconBlockHeader,
		*types.Eth1Data,
		*types.ExecutionPayloadHeader,
		*types.Fork,
		*types.Validator,
		*engineprimitives.Withdrawal,
	],
	BlobSidecarsT BlobSidecars,
	DepositStoreT DepositStore,
] struct {
	// logger is used for logging within the BeaconKitRuntime.
	logger log.Logger[any]
	// services is a registry of services used by the BeaconKitRuntime.
	services *service.Registry
	// chainSpec defines the chain specifications for the BeaconKitRuntime.
	chainSpec primitives.ChainSpec
}

// NewBeaconKitRuntime creates a new BeaconKitRuntime
// and applies the provided options.
func NewBeaconKitRuntime[
	AvailabilityStoreT AvailabilityStore[BeaconBlockBodyT, BlobSidecarsT],
	BeaconBlockT interface {
		types.RawBeaconBlock[BeaconBlockBodyT]
		NewFromSSZ([]byte, uint32) (BeaconBlockT, error)
		NewWithVersion(
			math.Slot,
			math.ValidatorIndex,
			primitives.Root,
			uint32,
		) (BeaconBlockT, error)
		Empty(uint32) BeaconBlockT
	},
	BeaconBlockBodyT types.RawBeaconBlockBody,
	BeaconStateT core.BeaconState[
		*types.BeaconBlockHeader, *types.Eth1Data,
		*types.ExecutionPayloadHeader, *types.Fork,
		*types.Validator, *engineprimitives.Withdrawal,
	],
	BlobSidecarsT BlobSidecars,
	DepositStoreT DepositStore,
](
	chainSpec primitives.ChainSpec,
	logger log.Logger[any],
	services *service.Registry,
	// storageBackend StorageBackendT,
) (*BeaconKitRuntime[
	AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, DepositStoreT,
], error) {
	return &BeaconKitRuntime[
		AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
		BlobSidecarsT, DepositStoreT,
	]{
		chainSpec: chainSpec,
		logger:    logger,
		services:  services,
	}, nil
}

// StartServices starts the services.
func (r *BeaconKitRuntime[
	AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, DepositStoreT,
]) StartServices(
	ctx context.Context,
) error {
	return r.services.StartAll(ctx)
}
