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
	"reflect"
)

// SizeOfBasic returns the size of a basic type.
func SizeOfBasic[RootT ~[32]byte, B Basic[SpecT, RootT], SpecT any](
	b B,
) uint64 {
	// TODO: Boolean maybe this doesnt work.
	return uint64(reflect.TypeOf(b).Size())
}

// ChunkCount returns the number of chunks required to store a value.
func ChunkCountBasic[RootT ~[32]byte, B Basic[SpecT, RootT], SpecT any](
	B,
) uint64 {
	return 1
}

// ChunkCountBitListVec returns the number of chunks required to store a bitlist
// or bitvector.
func ChunkCountBitListVec[T any](t []T) uint64 {
	//nolint:mnd // 256 is okay.
	return (uint64(len(t)) + 255) / 256
}

// ChunkCountBasicList returns the number of chunks required to store a list
// or vector of basic types.
func ChunkCountBasicList[SpecT any, RootT ~[32]byte, B Basic[SpecT, RootT]](
	b []B,
	maxCapacity uint64,
) uint64 {
	numItems := uint64(len(b))
	if numItems == 0 {
		return 1
	}
	size := SizeOfBasic[RootT, B, SpecT](b[0])
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
	SpecT any, RootT ~[32]byte, C Composite[SpecT, RootT],
](
	c []C,
	limit uint64,
) uint64 {
	return max(uint64(len(c)), limit)
}

// ChunkCountContainer returns the number of chunks required to store a
// container.
func ChunkCountContainer[SpecT any, RootT ~[32]byte, C Container[SpecT, RootT]](
	c C,
) uint64 {
	//#nosec:G701 // This is a safe operation.
	return uint64(reflect.ValueOf(c).NumField())
}
