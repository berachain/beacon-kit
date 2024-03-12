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

package beacon

import (
	"context"

	sdkcollections "cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	beacontypes "github.com/berachain/beacon-kit/beacon/core/types"
	"github.com/berachain/beacon-kit/lib/store/collections"
	"github.com/berachain/beacon-kit/lib/store/collections/encoding"
)

// Store is a wrapper around an sdk.Context
// that provides access to all beacon related data.
type Store struct {
	ctx context.Context

	// validatorIndex is a sequence that provides the next
	// available index for a new validator.
	validatorIndex sdkcollections.Sequence

	// validatorIndexToPubkey is a map that provides the
	// public key for a given validator index.
	validatorIndexToPubkey *sdkcollections.IndexedMap[
		uint64, []byte, validatorsIndex,
	]

	// depositQueue is a list of depositQueue that are queued to be processed.
	depositQueue *collections.Queue[*beacontypes.Deposit]

	// parentBlockRoot provides access to the previous
	// head block root for block construction as needed
	// by eip-4788.
	parentBlockRoot sdkcollections.Item[[]byte]
}

// Store creates a new instance of Store.
func NewStore(
	kvs store.KVStoreService,
) *Store {
	schemaBuilder := sdkcollections.NewSchemaBuilder(kvs)
	validatorIndex := sdkcollections.NewSequence(
		schemaBuilder,
		sdkcollections.NewPrefix(validatorIndexPrefix),
		validatorIndexPrefix,
	)
	validatorIndexToPubkey := sdkcollections.NewIndexedMap[uint64, []byte](
		schemaBuilder,
		sdkcollections.NewPrefix(validatorIndexToPubkeyPrefix),
		validatorIndexToPubkeyPrefix,
		sdkcollections.Uint64Key,
		sdkcollections.BytesValue,
		newValidatorsIndex(schemaBuilder),
	)
	depositQueue := collections.NewQueue[*beacontypes.Deposit](
		schemaBuilder,
		depositQueuePrefix,
		encoding.SSZValueCodec[*beacontypes.Deposit]{},
	)
	parentBlockRoot := sdkcollections.NewItem[[]byte](
		schemaBuilder,
		sdkcollections.NewPrefix(parentBlockRootPrefix),
		parentBlockRootPrefix,
		sdkcollections.BytesValue,
	)

	return &Store{
		validatorIndex:         validatorIndex,
		validatorIndexToPubkey: validatorIndexToPubkey,
		depositQueue:           depositQueue,
		parentBlockRoot:        parentBlockRoot,
	}
}

// WithContext returns the Store with the given context.
func (s *Store) WithContext(ctx context.Context) *Store {
	s.ctx = ctx
	return s
}
