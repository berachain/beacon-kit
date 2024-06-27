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

import "encoding/binary"

// HashFn is the generic hash function signature.
type HashFn func(input []byte) [32]byte

// HasherFunc defines a structure to hold a hash function and can be used for
// multiple rounds of
// hashing.
type HasherFunc[T ~[32]byte] struct {
	b        [64]byte
	hashFunc HashFn
}

// NewHasherFunc is the constructor for the object
// that fulfills the Hasher interface.
func NewHasherFunc[T ~[32]byte](h HashFn) *HasherFunc[T] {
	return &HasherFunc[T]{
		b:        [64]byte{},
		hashFunc: h,
	}
}

// Hash utilizes the provided hash function for
// the object.
func (h *HasherFunc[T]) Hash(a []byte) T {
	return T(h.hashFunc(a))
}

// Combi appends the two inputs and hashes them.
func (h *HasherFunc[T]) Combi(a, b T) T {
	copy(h.b[:32], a[:])
	copy(h.b[32:], b[:])
	return h.Hash(h.b[:])
}

// MixIn works like Combi, but using an integer as the second input.
//
//nolint:mnd // its okay.
func (h *HasherFunc[T]) MixIn(a T, i uint64) T {
	copy(h.b[:32], a[:])
	copy(h.b[32:], make([]byte, 32))
	binary.LittleEndian.PutUint64(h.b[32:], i)
	return h.Hash(h.b[:])
}
