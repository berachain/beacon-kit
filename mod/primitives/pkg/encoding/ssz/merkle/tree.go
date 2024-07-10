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

import "github.com/berachain/beacon-kit/mod/primitives/pkg/math/pow"

// New returns a Merkle tree of the given leaves.
// As defined in the Ethereum 2.0 Spec:
// https://github.com/ethereum/consensus-specs/blob/dev/ssz/merkle-proofs.md#generalized-merkle-tree-index
//
//nolint:lll // link.
func NewTree[LeafT ~[32]byte](
	leaves []LeafT,
	hashFn func([]byte) LeafT,
) []LeafT {
	/*
	   Return an array representing the tree nodes by generalized index:
	   [0, 1, 2, 3, 4, 5, 6, 7], where each layer is a power of 2. The 0 index is ignored. The 1 index is the root.
	   The result will be twice the size as the padded bottom layer for the input leaves.
	*/
	bottomLength := pow.NextPowerOfTwo(uint64(len(leaves)))
	//nolint:mnd // 2 is okay.
	o := make([]LeafT, bottomLength*2)
	copy(o[bottomLength:], leaves)
	for i := bottomLength - 1; i > 0; i-- {
		o[i] = hashFn(append(o[i*2][:], o[i*2+1][:]...))
	}
	return o
}
