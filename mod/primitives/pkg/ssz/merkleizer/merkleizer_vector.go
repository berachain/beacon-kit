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

// MerkleizeVectorBasic implements the SSZ merkleization algorithm
// for a vector of basic types.
func (m *merkleizer[RootT, T]) MerkleizeVectorBasic(
	value []T,
) (RootT, error) {
	// merkleize(pack(value))
	// if value is a basic object or a vector of basic objects.
	packed, _, err := pack[RootT](value)
	if err != nil {
		return [32]byte{}, err
	}
	return m.Merkleize(packed)
}

// MerkleizeVectorComposite implements the SSZ merkleization algorithm for a
// vector
// of composite types.
func (m *merkleizer[RootT, T]) MerkleizeVectorComposite(
	value []T,
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
	return m.Merkleize(htrs)
}
