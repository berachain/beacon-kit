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
	"slices"

	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto/sha256"
)

// Inspired by the Ethereum 2.0 spec:
// https://github.com/ethereum/consensus-specs/blob/dev/ssz/merkle-proofs.md#merkle-multiproofs
//
//nolint:lll // link.

// CalculateRoot calculates the Merkle root from the leaf and proof.
func CalculateRoot[RootT ~[32]byte](
	index GeneralizedIndex,
	leaf RootT,
	proof []RootT,
) (RootT, error) {
	if len(proof) != index.Length() {
		return RootT{},
			errors.Wrapf(ErrUnexpectedProofLength, "expected proof length %d, received %d", index.Length(),
				len(proof))
	}
	for i, h := range proof {
		if index.IndexBit(i) {
			leaf = sha256.Hash(append(h[:], leaf[:]...))
		} else {
			leaf = sha256.Hash(append(leaf[:], h[:]...))
		}
	}
	return leaf, nil
}

// VerifyProof verifies the Merkle proof for the given
// leaf, proof, and root.
func VerifyProof[RootT ~[32]byte](
	index GeneralizedIndex,
	leaf RootT,
	proof []RootT,
	root RootT,
) (bool, error) {
	calculated, err := CalculateRoot(index, leaf, proof)
	return calculated == root, err
}

// CalculateMultiRoot calculates the Merkle root for multiple leaves with
// their corresponding proofs and indices.
func CalculateMultiRoot[RootT ~[32]byte](
	indices GeneralizedIndices,
	leaves []RootT,
	proof []RootT,
) (RootT, error) {
	if len(leaves) != len(indices) {
		return RootT{}, errors.Wrapf(
			ErrMistmatchLeavesIndicesLength,
			"mismatched leaves and indices length: %d != %d",
			len(leaves), len(indices),
		)
	}

	helperIndices := indices.GetHelperIndices()
	if len(proof) != len(helperIndices) {
		return RootT{}, errors.New(
			"mismatched proof and helper indices length",
		)
	}

	objects := make(map[GeneralizedIndex]RootT)
	for i, index := range indices {
		objects[index] = leaves[i]
	}
	for i, index := range helperIndices {
		objects[index] = proof[i]
	}

	// Extract keys into slice to traverse in descending order.
	keys := make(GeneralizedIndices, 0, len(objects))
	for k := range objects {
		keys = append(keys, k)
	}
	slices.SortFunc(keys, GeneralizedIndexReverseComparator)

	return hashRoot(objects, keys), nil
}

// hashRoot hashes the objects in the given keys to the root.
func hashRoot[RootT ~[32]byte](
	objects map[GeneralizedIndex]RootT,
	keys GeneralizedIndices,
) RootT {
	var hashFn func([]byte) [32]byte
	if len(keys) > 5 { //nolint:mnd // 5 as defined by the library.
		hashFn = sha256.CustomHashFn()
	} else {
		hashFn = sha256.Hash
	}

	var (
		pos   int
		left  RootT
		right RootT
	)
	for pos < len(keys) {
		k := keys[pos]
		if _, ok := objects[k]; !ok {
			pos++
			continue
		}
		if _, ok := objects[k.Sibling()]; !ok {
			pos++
			continue
		}
		parent := k.Parent()
		if _, ok := objects[parent]; ok {
			pos++
			continue
		}

		if k%2 == 0 {
			left = objects[k]
			right = objects[k.Sibling()]
		} else {
			left = objects[k.Sibling()]
			right = objects[k]
		}
		objects[parent] = hashFn(append(left[:], right[:]...))
		keys = append(keys, parent)
		pos++
	}
	return objects[1]
}

// VerifyMultiproof verifies the Merkle multiproof by comparing the
// calculated root with the provided root.
func VerifyMultiproof[RootT ~[32]byte](
	indices GeneralizedIndices,
	leaves []RootT,
	proof []RootT,
	root RootT,
) bool {
	calculatedRoot, err := CalculateMultiRoot(indices, leaves, proof)
	if err != nil {
		return false
	}
	return calculatedRoot == root
}
