// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package engineprimitives

import (
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/constraints"
	fastssz "github.com/ferranbt/fastssz"
)

var (
	_ constraints.SSZRootable = (*Transactions)(nil)
)

// Transactions is a type alias for [][]byte, which is how
// transactions are received in the execution payload.
type Transactions [][]byte

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the SSZ encoded size in bytes for the Transactions.
func (txs Transactions) SizeSSZ() int {
	size := 4 // List offset
	for _, tx := range txs {
		size += 4 + len(tx) // Offset + transaction bytes
	}
	return size
}

// HashTreeRoot returns the hash tree root of the Transactions object.
func (txs Transactions) HashTreeRoot() ([32]byte, error) {
	hh := fastssz.DefaultHasherPool.Get()
	defer fastssz.DefaultHasherPool.Put(hh)
	if err := txs.HashTreeRootWith(hh); err != nil {
		return [32]byte{}, err
	}
	return hh.HashRoot()
}

/* -------------------------------------------------------------------------- */
/*                                   FastSSZ                                  */
/* -------------------------------------------------------------------------- */

// HashTreeRootWith ssz hashes the Transactions object with a hasher.
func (txs Transactions) HashTreeRootWith(hh fastssz.HashWalker) error {
	indx := hh.Index()
	num := uint64(len(txs))
	if num > constants.MaxTxsPerPayload {
		return fastssz.ErrIncorrectListSize
	}
	for _, tx := range txs {
		if len(tx) > int(constants.MaxBytesPerTx) {
			return fastssz.ErrBytesLength
		}
		// Each transaction is hashed as bytes
		root := merkleizeBytesN(tx)
		hh.AppendBytes32(root[:])
	}
	hh.MerkleizeWithMixin(indx, num, constants.MaxTxsPerPayload)
	return nil
}

// GetTree ssz hashes the Transactions object.
func (txs Transactions) GetTree() (*fastssz.Node, error) {
	return fastssz.ProofTree(txs)
}

// merkleizeBytesN returns the hash tree root of a byte slice.
func merkleizeBytesN(data []byte) [32]byte {
	// Merkleize bytes by padding to chunks
	chunkCount := (len(data) + 31) / 32
	if chunkCount == 0 {
		chunkCount = 1
	}

	// Create a temporary hasher
	hh := fastssz.DefaultHasherPool.Get()
	defer fastssz.DefaultHasherPool.Put(hh)

	// Start merkleization
	indx := hh.Index()

	// Append data in 32-byte chunks
	for i := 0; i < len(data); i += 32 {
		var chunk [32]byte
		end := i + 32
		if end > len(data) {
			end = len(data)
		}
		copy(chunk[:], data[i:end])
		hh.Append(chunk[:])
	}

	// Merkleize with length mixin
	hh.MerkleizeWithMixin(indx, uint64(len(data)), uint64(len(data)))
	root, _ := hh.HashRoot()
	return root
}
