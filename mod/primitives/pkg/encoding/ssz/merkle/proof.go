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

import (
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto/sha256"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math/pow"
)

// BuildProofFromLeaves builds a Merkle proof from the given leaves and the
// index of the leaf. The leaves are assumed to be hashed into 32 byte roots.
func BuildProofFromLeaves[RootT ~[32]byte](
	leaves []RootT,
	index uint64,
) ([]RootT, error) {
	tree, depth := newTree(leaves)
	return buildSingleProofFromTree(tree, NewGeneralizedIndex(depth, index))
}

// newTree returns a Merkle tree of the given leaves. Returns an array
// representing the tree nodes by generalized index: [0, 1, 2, 3, 4, 5, 6, 7],
// where each layer is a power of 2. The 0 index is ignored. The 1 index is the
// root. The result will be twice the size as the padded bottom layer for the
// input leaves. Also returns the depth of the tree.
//
// As defined in the Ethereum 2.0 Spec:
// https://github.com/ethereum/consensus-specs/blob/dev/ssz/merkle-proofs.md#generalized-merkle-tree-index
//
//nolint:lll // link.
func newTree[RootT ~[32]byte](leaves []RootT) ([]RootT, uint8) {
	bottomLength := pow.NextPowerOfTwo(uint64(len(leaves)))
	//nolint:mnd // 2 is okay.
	o := make([]RootT, bottomLength*2)
	copy(o[bottomLength:], leaves)

	var hashFn func([]byte) [32]byte
	if bottomLength > 5 { //nolint:mnd // 5 as defined by the library.
		hashFn = sha256.CustomHashFn()
	} else {
		hashFn = sha256.Hash
	}

	for i := bottomLength - 1; i > 0; i-- {
		o[i] = hashFn(append(o[i*2][:], o[i*2+1][:]...))
	}
	return o, log.ILog2Ceil(bottomLength)
}

// buildSingleProofFromTree returns a Merkle proof of the given tree from the
// given leaf index. Tree nodes are assumed to be ordered by generalized index.
//
// Inspired by the Ethereum 2.0 Spec:
// https://github.com/ethereum/consensus-specs/blob/dev/tests/core/pyspec/eth2spec/test/helpers/merkle.py
//
//nolint:lll // link.
func buildSingleProofFromTree[RootT ~[32]byte](
	tree []RootT,
	index GeneralizedIndex,
) ([]RootT, error) {
	treeLen := GeneralizedIndex(len(tree))
	if pow.PrevPowerOfTwo(treeLen) != treeLen {
		return nil, errors.Newf(
			"invalid tree length (%d), must be power of 2", treeLen,
		)
	}
	if index >= treeLen {
		return nil, errors.Newf(
			"generalized index (%d) must be less than tree length (%d)",
			index, treeLen,
		)
	}
	if index < pow.PrevPowerOfTwo(treeLen-1) {
		return nil, errors.Newf(
			"generalized index (%d) must be of a leaf in the tree", index,
		)
	}

	depth := index.Length()
	proof := make([]RootT, depth)
	for i := range depth {
		proof[i] = tree[index.Sibling()]
		index = index.Parent()
	}
	return proof, nil
}
