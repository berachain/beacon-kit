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

package store

import (
	"context"

	sdkcollections "cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"github.com/ethereum/go-ethereum/common"
	"github.com/itsdevbear/bolaris/lib/store/collections"
	"github.com/itsdevbear/bolaris/lib/store/collections/encoding"
	consensusv1 "github.com/itsdevbear/bolaris/types/consensus/v1"
)

// BeaconStore is a wrapper around an sdk.Context
// that provides access to all beacon related data.
type BeaconStore struct {
	ctx context.Context

	// depositQueue is a list of depositQueue that are queued to be processed.
	depositQueue *collections.Queue[*consensusv1.Deposit]

	// fcSafeEth1BlockHash is the safe block hash.
	fcSafeEth1BlockHash sdkcollections.Item[[32]byte]

	// fcFinalizedEth1BlockHash is the finalized block hash.
	fcFinalizedEth1BlockHash sdkcollections.Item[[32]byte]

	// eth1GenesisHash is the Eth1 genesis hash.
	eth1GenesisHash sdkcollections.Item[[32]byte]

	// parentBlockRoot provides access to the previous
	// head block root for block construction as needed
	// by eip-4788.
	parentBlockRoot sdkcollections.Item[[]byte]

	// lastValidHash is the last valid head in the store.
	// TODO: we need to handle this in a better way.
	lastValidHash *common.Hash
}

// NewBeaconStore creates a new instance of BeaconStore.
func NewBeaconStore(
	kvs store.KVStoreService,
) *BeaconStore {
	schemaBuilder := sdkcollections.NewSchemaBuilder(kvs)
	depositQueue := collections.NewQueue[*consensusv1.Deposit](
		schemaBuilder,
		depositQueuePrefix,
		encoding.SSZValueCodec[*consensusv1.Deposit]{},
	)
	fcSafeEth1BlockHash := sdkcollections.NewItem[[32]byte](
		schemaBuilder,
		sdkcollections.NewPrefix(fcSafeEth1BlockHashPrefix),
		fcSafeEth1BlockHashPrefix,
		encoding.Bytes32ValueCodec{},
	)
	fcFinalizedEth1BlockHash := sdkcollections.NewItem[[32]byte](
		schemaBuilder,
		sdkcollections.NewPrefix(fcFinalizedEth1BlockHashPrefix),
		fcFinalizedEth1BlockHashPrefix,
		encoding.Bytes32ValueCodec{},
	)
	eth1GenesisHash := sdkcollections.NewItem[[32]byte](
		schemaBuilder,
		sdkcollections.NewPrefix(eth1GenesisHashPrefix),
		eth1GenesisHashPrefix,
		encoding.Bytes32ValueCodec{},
	)
	parentBlockRoot := sdkcollections.NewItem[[]byte](
		schemaBuilder,
		sdkcollections.NewPrefix(parentBlockRootPrefix),
		parentBlockRootPrefix,
		sdkcollections.BytesValue,
	)
	return &BeaconStore{
		depositQueue:             depositQueue,
		fcSafeEth1BlockHash:      fcSafeEth1BlockHash,
		fcFinalizedEth1BlockHash: fcFinalizedEth1BlockHash,
		eth1GenesisHash:          eth1GenesisHash,
		parentBlockRoot:          parentBlockRoot,
	}
}

// WithContext( returns the BeaconStore with the given context.
func (s *BeaconStore) WithContext(ctx context.Context) *BeaconStore {
	s.ctx = ctx
	return s
}
