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

package block

import (
	"context"
	"sync"

	sdkcollections "cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb/encoding"
)

// var _ pruner.Prunable = (*KVStore[BeaconBlock])(nil)

const KeyBlockPrefix = "block"

type KVStoreProvider struct {
	store.KVStoreWithBatch
}

// OpenKVStore opens a new KV store.
func (p *KVStoreProvider) OpenKVStore(context.Context) store.KVStore {
	return p.KVStoreWithBatch
}

// KVStore is a simple KV store based implementation that stores
// beacon blocks.
type KVStore[BeaconBlockT BeaconBlock[BeaconBlockT]] struct {
	store sdkcollections.Map[uint64, BeaconBlockT]
	mu    sync.RWMutex
	cdc   *encoding.SSZInterfaceCodec[BeaconBlockT]

	earliestSlot uint64
}

// NewStore creates a new block store.
func NewStore[BeaconBlockT BeaconBlock[BeaconBlockT]](
	kvsp store.KVStoreService,
) *KVStore[BeaconBlockT] {
	schemaBuilder := sdkcollections.NewSchemaBuilder(kvsp)
	cdc := &encoding.SSZInterfaceCodec[BeaconBlockT]{}
	return &KVStore[BeaconBlockT]{
		store: sdkcollections.NewMap(
			schemaBuilder,
			sdkcollections.NewPrefix([]byte{uint8(0)}),
			KeyBlockPrefix,
			sdkcollections.Uint64Key,
			cdc,
		),
		mu:           sync.RWMutex{},
		earliestSlot: 0,
		cdc:          cdc,
	}
}

// Get retrieves the block by a given index from the store.
func (kv *KVStore[BeaconBlockT]) Get(slot uint64) (BeaconBlockT, error) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	return kv.store.Get(context.TODO(), slot)
}

// Set sets the block by a given index in the store.
func (kv *KVStore[BeaconBlockT]) Set(slot uint64, blk BeaconBlockT) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	kv.cdc.SetActiveForkVersion(blk.Version())
	return kv.store.Set(context.TODO(), slot, blk)
}

// Prune removes the [start, end) blocks from the store.
func (kv *KVStore[BeaconBlockT]) Prune(start, end uint64) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	for i := max(start, kv.earliestSlot); i < end; i++ {
		// This only errors if the key passed in cannot be encoded.
		if err := kv.store.Remove(context.TODO(), i); err != nil {
			return err
		}
	}
	kv.earliestSlot = end
	return nil
}
