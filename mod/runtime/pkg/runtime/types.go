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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
)

// AppOptions is an interface that provides the ability to
// retrieve options from the application.
type AppOptions interface {
	Get(string) interface{}
}

// The AvailabilityStore interface is responsible for validating and storing
// sidecars for specific blocks, as well as verifying sidecars that have already
// been stored.
type AvailabilityStore[BeaconBlockBodyT, BlobSidecarsT any] interface {
	// IsDataAvailable ensures that all blobs referenced in the block are
	// securely stored before it returns without an error.
	IsDataAvailable(
		context.Context, math.Slot, BeaconBlockBodyT,
	) bool
	// Persist makes sure that the sidecar remains accessible for data
	// availability checks throughout the beacon node's operation.
	Persist(math.Slot, BlobSidecarsT) error
}

// StorageBackend defines an interface for accessing various storage components
// required by the beacon node.
type StorageBackend[
	AvailabilityStoreT AvailabilityStore[BeaconBlockBodyT, BlobSidecarsT],
	BeaconBlockBodyT,
	BeaconStateT,
	BlobSidecarsT any,
	DepositStoreT DepositStore,
] interface {
	// AvailabilityStore returns the availability store for the given context.
	AvailabilityStore(context.Context) AvailabilityStoreT
	// StateFromContext retrieves the beacon state from the given context.
	StateFromContext(context.Context) BeaconStateT
	// DepositStore returns the deposit store for the given context.
	DepositStore(context.Context) DepositStoreT
}

// BlobSidecars is an interface that represents the sidecars.
type BlobSidecars interface {
	ssz.Marshallable
	IsNil() bool
	Len() int
}

// DepositStore is an interface that provides the
// expected deposits to the runtime.
type DepositStore interface {
	GetDepositsByIndex(
		startIndex uint64,
		numView uint64,
	) ([]*types.Deposit, error)
	EnqueueDeposits(deposits []*types.Deposit) error
	Prune(index uint64, numPrune uint64) error
}

// Service is a struct that can be registered into a ServiceRegistry for
// easy dependency management.
type Service interface {
	// Start spawns any goroutines required by the service.
	Start(ctx context.Context)
	// Status returns error if the service is not considered healthy.
	Status() error
}
