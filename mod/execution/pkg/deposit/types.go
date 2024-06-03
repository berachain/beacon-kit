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

package deposit

import (
	"context"
	"math/big"

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// BeaconBlock is an interface for beacon blocks.
type BeaconBlock[
	BeaconBlockBodyT BeaconBlockBody[DepositT, ExecutionPayloadT],
	DepositT interface{ GetIndex() uint64 },
	ExecutionPayloadT interface{ GetNumber() math.U64 },
] interface {
	GetSlot() math.U64
	GetBody() BeaconBlockBodyT
}

type BeaconBlockBody[
	DepositT interface{ GetIndex() uint64 },
	ExecutionPayloadT interface{ GetNumber() math.U64 },
] interface {
	GetDeposits() []DepositT
	GetExecutionPayload() ExecutionPayloadT
}

// BlockEvent is an interface for block events.
type BlockEvent[
	BeaconBlockT BeaconBlock[BeaconBlockBodyT, DepositT, ExecutionPayloadT],
	BeaconBlockBodyT BeaconBlockBody[DepositT, ExecutionPayloadT],
	DepositT interface{ GetIndex() uint64 },
	ExecutionPayloadT interface{ GetNumber() math.U64 },
] interface {
	Context() context.Context
	Block() BeaconBlockT
}

// BlockFeed is an interface for subscribing to block events.
type BlockFeed[
	BeaconBlockT BeaconBlock[BeaconBlockBodyT, DepositT, ExecutionPayloadT],
	BeaconBlockBodyT BeaconBlockBody[DepositT, ExecutionPayloadT],
	BlockEventT BlockEvent[
		BeaconBlockT, BeaconBlockBodyT, DepositT, ExecutionPayloadT],
	DepositT interface{ GetIndex() uint64 },
	ExecutionPayloadT interface{ GetNumber() math.U64 },
	SubscriptionT interface {
		Unsubscribe()
	}] interface {
	Subscribe(chan<- (BlockEventT)) SubscriptionT
}

// Contract is the ABI for the deposit contract.
type Contract[DepositT any] interface {
	// ReadDeposits reads deposits from the deposit contract.
	ReadDeposits(
		ctx context.Context,
		blockNumber math.U64,
	) ([]DepositT, error)
}

// Deposit is an interface for deposits.
type Deposit[DepositT, WithdrawalCredentialsT any] interface {
	// New creates a new deposit.
	New(
		crypto.BLSPubkey,
		WithdrawalCredentialsT,
		math.U64,
		crypto.BLSSignature,
		uint64,
	) DepositT
}

// EthClient is an interface for interacting with the Ethereum 1.0 client.
type EthClient interface {
	BlockByNumber(
		ctx context.Context,
		number *big.Int,
	) (*engineprimitives.Block, error)
}

// Store defines the interface for managing deposit operations.
type Store[DepositT any] interface {
	// PruneFromInclusive prunes the deposit store from the given
	// index for N indexes.
	PruneFromInclusive(index uint64, numPrune uint64) error
	// EnqueueDeposits adds a list of deposits to the deposit store.
	EnqueueDeposits(deposits []DepositT) error
}

type StorageBackend[
	AvailabilityStoreT any,
	BeaconStateT any,
	BlobSidecarsT any,
	DepositStoreT Store[DepositT],
	DepositT any,
] interface {
	// DepositStore returns the deposit store for the given context.
	DepositStore(context.Context) DepositStoreT
}
