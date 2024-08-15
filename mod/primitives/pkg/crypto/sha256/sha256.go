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

package sha256

import (
	"hash"
	"sync"

	"github.com/minio/sha256-simd"
)

// sha256Pool is a pool of sha256 hash functions.
//
//nolint:gochecknoglobals // needed for pool.
var sha256Pool = sync.Pool{New: func() interface{} {
	return sha256.New()
}}

// Hash defines a function that returns the sha256 checksum of the data passed
// in. Adheres to the crypto.HashFn signature.
// https://github.com/ethereum/consensus-specs/blob/v0.9.3/specs/core/0_beacon-chain.md#hash
//
//nolint:lll // url.
func Hash(data []byte) [32]byte {
	h, ok := sha256Pool.Get().(hash.Hash)
	if !ok {
		h = sha256.New()
	}
	defer sha256Pool.Put(h)
	h.Reset()

	var b [32]byte
	//#nosec:G104 bet
	h.Write(data)
	h.Sum(b[:0])
	return b
}

// CustomHashFn provides a hash function utilizing
// an internal hasher. It is not thread-safe as the same
// hasher instance is reused.
//
// Note: This method is more efficient only if the callback
// is invoked more than 5 times.
func CustomHashFn() func([]byte) [32]byte {
	hasher, ok := sha256Pool.Get().(hash.Hash)
	if !ok {
		hasher = sha256.New()
	} else {
		hasher.Reset()
	}
	var h [32]byte

	return func(data []byte) [32]byte {
		//#nosec:G104 // bet
		hasher.Write(data)
		hasher.Sum(h[:0])
		hasher.Reset()

		return h
	}
}
