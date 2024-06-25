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

package ssz

// Hashable is an interface representing objects that implement HashTreeRoot().
type Hashable[SpecT any, Root ~[32]byte] interface {
	HashTreeRoot() (Root, error)
}

// U64 is an interface for uint64 types that support
// NextPowerOfTwo and ILog2Ceil.
type U64[T ~uint64] interface {
	~uint64
	NextPowerOfTwo() T
	ILog2Ceil() uint8
}

// U128LT represents a 128-bit unsigned integer in
// little-endian byte order.
type U128LT interface {
	~[16]byte
}

// U256LT represents a 256-bit unsigned integer in
// little-endian byte order.
type U256LT interface {
	~[32]byte
}
