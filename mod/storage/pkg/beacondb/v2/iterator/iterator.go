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

package iterator

import (
	"cosmossdk.io/store"
)

// Iterator is a wrapper around the store and changeset iterators,
// and provides an iterator which iterates over all changes in the
// changeset and store, skipping duplicates.
type iterator struct {
	changeset  store.Iterator
	chainStore store.Iterator
	blockStore store.Iterator

	start []byte
	end   []byte

	seen map[string]struct{}
}

func New(
	start, end []byte,
	chainStore store.Iterator,
	changeset store.Iterator,
	blockStore store.Iterator,
) store.Iterator {
	return &iterator{
		changeset:  changeset,
		chainStore: chainStore,
		start:      start,
		end:        end,
		blockStore: blockStore,
		seen:       make(map[string]struct{}),
	}
}

func (i *iterator) Domain() (start []byte, end []byte) {
	return i.start, i.end
}

func (i *iterator) Valid() bool {
	return i.changeset.Valid() || i.blockStore.Valid() || i.chainStore.Valid()
}

func (i *iterator) Next() {
	// get next value from changeset if valid
	if i.changeset.Valid() {
		i.changeset.Next()
		if i.changeset.Valid() {
			i.seen[string(i.changeset.Key())] = struct{}{}
		}
		return
	}

	// otherwise, iterate over the block store until a valid key is found
	for i.blockStore.Valid() {
		i.blockStore.Next()
		if !i.blockStore.Valid() {
			break
		}
		if _, ok := i.seen[string(i.blockStore.Key())]; !ok {
			i.seen[string(i.blockStore.Key())] = struct{}{}
			return
		}
	}

	// finally, iterate over the chain store until a valid key is found
	for i.chainStore.Valid() {
		i.chainStore.Next()
		if !i.chainStore.Valid() {
			break
		}
		if _, ok := i.seen[string(i.chainStore.Key())]; !ok {
			return
		}
	}
}

func (i *iterator) Key() (key []byte) {
	if i.changeset.Valid() {
		return i.changeset.Key()
	} else if i.blockStore.Valid() {
		return i.blockStore.Key()
	}
	return i.chainStore.Key()
}

func (i *iterator) Value() (value []byte) {
	if i.changeset.Valid() {
		return i.changeset.Value()
	} else if i.blockStore.Valid() {
		return i.blockStore.Value()
	}
	return i.chainStore.Value()
}

func (i *iterator) Error() error {
	if i.changeset.Valid() {
		return i.changeset.Error()
	} else if i.blockStore.Valid() {
		return i.blockStore.Error()
	}
	return i.chainStore.Error()
}

func (i *iterator) Close() error {
	if err := i.changeset.Close(); err != nil {
		return err
	}
	if err := i.blockStore.Close(); err != nil {
		return err
	}
	return i.chainStore.Close()
}
