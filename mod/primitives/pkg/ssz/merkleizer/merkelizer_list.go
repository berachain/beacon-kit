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

package merkleizer

import "github.com/berachain/beacon-kit/mod/primitives/pkg/merkle"

// MerkleizeListBasic implements the SSZ merkleization algorithm for a list of
// basic types.
func (m *merkleizer[SpecT, RootT, T]) MerkleizeListBasic(
	value []T,
	limit ...uint64,
) (RootT, error) {
	packed, err := m.pack(value)
	if err != nil {
		return [32]byte{}, err
	}

	var effectiveLimit uint64
	if len(limit) > 0 {
		effectiveLimit = limit[0]
	} else {
		effectiveLimit = uint64(len(packed))
	}

	root, err := m.Merkleize(
		packed, ChunkCountBasicList[SpecT](value, effectiveLimit),
	)
	if err != nil {
		return [32]byte{}, err
	}
	return merkle.MixinLength(root, uint64(len(value))), nil
}

// MerkleizeListComposite implements the SSZ merkleization algorithm for a list
// of composite types.
func (m *merkleizer[SpecT, RootT, T]) MerkleizeListComposite(
	value []T,
	limit ...uint64,
) (RootT, error) {
	var (
		err  error
		htrs = m.bytesBuffer.Get(len(value))
	)

	for i, el := range value {
		htrs[i], err = el.HashTreeRoot()
		if err != nil {
			return RootT{}, err
		}
	}

	var effectiveLimit uint64
	if len(limit) > 0 {
		effectiveLimit = limit[0]
	} else {
		effectiveLimit = uint64(len(value))
	}

	root, err := m.Merkleize(
		htrs, ChunkCountCompositeList[SpecT](value, effectiveLimit),
	)
	if err != nil {
		return RootT{}, err
	}

	return merkle.MixinLength(root, uint64(len(value))), nil
}
