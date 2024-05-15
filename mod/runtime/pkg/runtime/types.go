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

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
	ssz "github.com/ferranbt/fastssz"
)

// AppOptions is an interface that provides the ability to
// retrieve options from the application.
type AppOptions interface {
	Get(string) interface{}
}

// BeaconStorageBackend is an interface that provides the
// beacon state to the runtime.
type BeaconStorageBackend[
	BeaconBlockT,
	BeaconStateT,
	BlobSidecarsT any,
	DepositStoreT DepositStore,
] interface {
	AvailabilityStore(
		ctx context.Context,
	) core.AvailabilityStore[BeaconBlockT, BlobSidecarsT]
	StateFromContext(ctx context.Context) BeaconStateT
	DepositStore(ctx context.Context) DepositStoreT
}

// BlobSidecars is an interface that represents the sidecars.
type BlobSidecars interface {
	ssz.Marshaler
	ssz.Unmarshaler
	Len() int
}

type Config interface{}

// DepositStore is an interface that provides the
// expected deposits to the runtime.
type DepositStore interface {
	ExpectedDeposits(
		numView uint64,
	) ([]*types.Deposit, error)
	EnqueueDeposits(deposits []*types.Deposit) error
	DequeueDeposits(
		numDequeue uint64,
	) ([]*types.Deposit, error)
	PruneToIndex(
		index uint64,
	) error
}

// Service is a struct that can be registered into a ServiceRegistry for
// easy dependency management.
type Service interface {
	// Start spawns any goroutines required by the service.
	Start(ctx context.Context)
	// Status returns error if the service is not considered healthy.
	Status() error
}
