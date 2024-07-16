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
type KVStore[BeaconBlockT BeaconBlock] struct {
	store sdkcollections.Map[uint64, BeaconBlockT]
	mu    sync.RWMutex
}

// NewStore creates a new block store.
func NewStore[BeaconBlockT BeaconBlock](
	kvsp store.KVStoreService,
) *KVStore[BeaconBlockT] {
	schemaBuilder := sdkcollections.NewSchemaBuilder(kvsp)
	return &KVStore[BeaconBlockT]{
		store: sdkcollections.NewMap(
			schemaBuilder,
			sdkcollections.NewPrefix([]byte{uint8(0)}),
			KeyBlockPrefix,
			sdkcollections.Uint64Key,
			encoding.SSZValueCodec[BeaconBlockT]{},
		),
	}
}

// Get retrieves the block by a given index from the store.
func (kv *KVStore[BeaconBlockT]) Get(slot uint64) (BeaconBlockT, error) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	var (
		blk BeaconBlockT
		err error
	)
	blk, err = kv.store.Get(context.TODO(), slot)
	return blk, err
}

// Set sets the block by a given index in the store.
func (kv *KVStore[BeaconBlockT]) Set(slot uint64, blk BeaconBlockT) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	return kv.store.Set(context.TODO(), slot, blk)
}
