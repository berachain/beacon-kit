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

package buffer

// initialBufferSize is the initial size of the internal buffer.
const initialBufferSize = 64

// ReusableBuffer is a re-usable buffer for merkle tree hashing. Prevents
// unnecessary allocations and garbage collection of byte slices.
//
// NOTE: this buffer is currently only safe for use in a single thread.
type ReusableBuffer[RootT ~[32]byte] struct {
	internal []RootT
}

// NewReusableBuffer creates a new re-usable buffer for merkle tree hashing.
func NewReusableBuffer[RootT ~[32]byte]() *ReusableBuffer[RootT] {
	return &ReusableBuffer[RootT]{
		internal: make([]RootT, initialBufferSize),
	}
}

// Get returns a slice of the internal buffer of roots of the given size.
func (b *ReusableBuffer[RootT]) Get(size int) []RootT {
	if delta := size - len(b.internal); delta > 0 {
		b.grow(delta)
	}

	return b.internal[:size]
}

// grow resizes the internal buffer by the requested delta.
func (b *ReusableBuffer[RootT]) grow(delta int) {
	b.internal = append(b.internal, make([]RootT, delta)...)
}

// singleuseBuffer is a buffer for a single use case. Allocates new
// memory for each use (call to `Get`).
//
// NOTE: this buffer is only used for testing.
type SingleUseBuffer[RootT ~[32]byte] struct{}

// NewSingleuseBuffer creates a new single-use buffer.
func NewSingleuseBuffer[RootT ~[32]byte]() *SingleUseBuffer[RootT] {
	return &SingleUseBuffer[RootT]{}
}

// Get returns a new slice of roots the given size.
func (b *SingleUseBuffer[RootT]) Get(size int) []RootT {
	return make([]RootT, size)
}
