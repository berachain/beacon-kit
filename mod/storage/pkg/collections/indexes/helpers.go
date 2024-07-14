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

package indexes

import (
	"github.com/berachain/beacon-kit/mod/storage/pkg/collections"
)

// iterator defines the minimum set of methods of an index iterator
// required to work with the helpers.
type iterator[K any] interface {
	// PrimaryKey returns the iterator current primary key.
	PrimaryKey() (K, error)
	// Next advances the iterator by one element.
	Next()
	// Valid asserts if the Iterator is valid.
	Valid() bool
	// Close closes the iterator.
	Close() error
}

// ScanValues collects all the values from an Index iterator and the IndexedMap in a lazy way.
// The iterator is closed when this function exits.
func ScanValues[K, V any, I iterator[K], Idx collections.Indexes[K, V]](
	indexedMap *collections.IndexedMapKeeper[K, V, Idx],
	iter I,
	f func(value V) (stop bool),
) error {
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		key, err := iter.PrimaryKey()
		if err != nil {
			return err
		}

		value, err := indexedMap.Get(key)
		if err != nil {
			return err
		}

		stop := f(value)
		if stop {
			return nil
		}
	}

	return nil
}
