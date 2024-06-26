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

package bytes

// initialBufferSize is the initial size of the internal buffer.
const initialBufferSize = 64

// Buffer can be used by hashers to get a buffer of 32 byte slices.
type Buffer[RootT ~[32]byte] interface {
	// Get returns a slice of roots of the given size.
	Get(size int) []RootT

	// TODO: add a Put method to return the buffer back for concurrent use.
}

// reusableBuffer is a re-usable buffer. Prevents unnecessary allocations
// and garbage collection of byte slices.
type reusableBuffer[RootT ~[32]byte] struct {
	internal []RootT

	// TODO: add a mutex for multi-thread safety.
}

// NewReusableBuffer creates a new re-usable buffer.
func NewReusableBuffer[RootT ~[32]byte]() Buffer[RootT] {
	return &reusableBuffer[RootT]{
		internal: make([]RootT, initialBufferSize),
	}
}

// Get returns a slice of the internal buffer of roots of the given size.
func (b *reusableBuffer[RootT]) Get(size int) []RootT {
	if size > len(b.internal) {
		b.grow(size - len(b.internal))
	}

	return b.internal[:size]
}

// grow resizes the internal buffer by the requested delta.
func (b *reusableBuffer[RootT]) grow(delta int) {
	b.internal = append(b.internal, make([]RootT, delta)...)
}

// singleuseBuffer is a buffer for a single use case. Allocates new
// memory for each use (call to `Get`).
type singleuseBuffer[RootT ~[32]byte] struct{}

// NewSingleuseBuffer creates a new single-use buffer.
func NewSingleuseBuffer[RootT ~[32]byte]() Buffer[RootT] {
	return &singleuseBuffer[RootT]{}
}

// Get returns a new slice of roots the given size.
func (*singleuseBuffer[RootT]) Get(size int) []RootT {
	return make([]RootT, size)
}
