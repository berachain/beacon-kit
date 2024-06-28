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

package merkleizer

import (
	"encoding/binary"
	"reflect"

	"github.com/prysmaticlabs/gohashtree"
)

// ChunkCountBitListVec returns the number of chunks required to store a bitlist
// or bitvector.
func ChunkCountBitListVec[T any](t []T) uint64 {
	//nolint:mnd // 256 is okay.
	return (uint64(len(t)) + 255) / 256
}

// ChunkCountBasicList returns the number of chunks required to store a list
// or vector of basic types.
func ChunkCountBasicList[RootT ~[32]byte, B SSZObject[RootT]](
	b []B,
	maxCapacity uint64,
) uint64 {
	numItems := uint64(len(b))
	if numItems == 0 {
		return 1
	}
	//#nosec:G103 // its fine.
	size := uint64(b[0].SizeSSZ())
	//nolint:mnd // 32 is okay.
	limit := (maxCapacity*size + 31) / 32
	if limit != 0 {
		return limit
	}

	return numItems
}

// ChunkCountCompositeList returns the number of chunks required to store a
// list or vector of composite types.
func ChunkCountCompositeList[
	RootT ~[32]byte, C SSZObject[RootT],
](
	c []C,
	limit uint64,
) uint64 {
	return max(uint64(len(c)), limit)
}

// ChunkCountContainer returns the number of chunks required to store a
// container.
func ChunkCountContainer[RootT ~[32]byte, C SSZObject[RootT]](
	c C,
) uint64 {
	//#nosec:G701 // This is a safe operation.
	return uint64(reflect.ValueOf(c).NumField())
}

// MixinLength mixes in the length of an element.
func MixinLength[RootT ~[32]byte](element RootT, length uint64) RootT {
	// Mix in the length of the element.
	//
	//nolint:mnd // its okay.
	chunks := make([][32]byte, 2)
	chunks[0] = element
	binary.LittleEndian.PutUint64(chunks[1][:], length)
	if err := gohashtree.Hash(chunks, chunks); err != nil {
		return [32]byte{}
	}
	return chunks[0]
}
