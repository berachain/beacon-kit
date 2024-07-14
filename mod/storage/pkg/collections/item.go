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

package collections

import (
	"cosmossdk.io/collections/codec"
)

// An ItemKeeper is an intermediary that holds a certain configuration,
// and uses that configuration to interact with the underlying store.
type ItemKeeper[V any] struct {
	storeKey      []byte
	key           []byte
	valueCodec    codec.ValueCodec[V]
	storeAccessor StoreAccessor
}

func NewItemKeeper[V any](
	storeKey []byte,
	key []byte,
	valueCodec codec.ValueCodec[V],
	storeAccessor StoreAccessor,
) ItemKeeper[V] {
	return ItemKeeper[V]{
		storeKey:      storeKey,
		key:           key,
		valueCodec:    valueCodec,
		storeAccessor: storeAccessor,
	}
}

// Get retrieves the value from the store, and returns the decoded value.
func (i *ItemKeeper[V]) Get() (V, error) {
	var result V
	res, err := i.storeAccessor().QueryState(i.storeKey, i.key)
	if err != nil {
		return result, err
	}

	return i.valueCodec.Decode(res)
}

// Set sets the value in the store with the encoded key and value.
func (i *ItemKeeper[V]) Set(value V) error {
	encodedValue, err := i.valueCodec.Encode(value)
	if err != nil {
		return err
	}
	i.storeAccessor().AddChange(i.storeKey, i.key, encodedValue)
	return nil
}
