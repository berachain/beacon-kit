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

package merkle

import "sync"

// initialBufferSize is the initial size of the internal buffer.
//
// TODO: choose a more appropriate size?
const initialBufferSize = 16

// TODO: remove this once buffer supports multi-threaded multi-use.
var bufferPool = sync.Pool{
	New: func() interface{} {
		return NewBuffer[[32]byte]()
	},
}

// buffer is a re-usable buffer for merkle tree hashing. Prevents
// unnecessary allocations and garbage collection of byte slices.
//
// NOTE: this buffer is ONLY meant to be used in a single thread.
type buffer[RootT ~[32]byte] struct {
	internal []RootT

	// TODO: add a mutex for multi-thread safety.
}

// NewBuffer creates a new buffer with the given capacity.
func NewBuffer[RootT ~[32]byte]() *buffer[RootT] {
	return &buffer[RootT]{
		internal: make([]RootT, initialBufferSize),
	}
}

// Get returns a slice of the internal buffer of roots of the given size.
func (b *buffer[RootT]) Get(size int) []RootT {
	if size > len(b.internal) {
		b.grow(size - len(b.internal))
	}

	return b.internal[:size]
}

// TODO: add a Put method to return the buffer back for multi-threaded multi-use.

// grow resizes the internal buffer by the requested size.
func (b *buffer[RootT]) grow(newSize int) {
	b.internal = append(b.internal, make([]RootT, newSize)...)
}
