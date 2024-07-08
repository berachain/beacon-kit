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

package schema

import (
	"errors"
	"strings"
)

// ObjectPath represents a path to an object in a Merkle tree.
type ObjectPath[RootT ~[32]byte] string

// Split returns the path split by "/".
func (p ObjectPath[_]) Split() []string {
	return strings.Split(string(p), "/")
}

// GetGeneralizedIndex converts a path to a generalized index representing its
// position in the Merkle tree.
func (p ObjectPath[RootT]) GetGeneralizedIndex(
	typ SSZType,
) (GeneralizedIndex[RootT], uint8, error) {
	gIndex := GeneralizedIndex[RootT](1)
	offset := uint8(0)
	for _, part := range p.Split() {
		if typ.ID().IsBasic() {
			return 0, 0, errors.New("cannot descend further from basic type")
		}
		if part == "__len__" {
			if !typ.ID().IsList() {
				return 0, 0, errors.New("__len__ is only valid for List types")
			}
			gIndex = gIndex.RightChild()
		} else {
			pos, start, _, err := typ.ItemPosition(part)
			if err != nil {
				return 0, 0, err
			}
			gIndex = GeneralizedIndex[RootT](
				uint64(gIndex)*getBaseIndex(typ)*nextPowerOfTwo(typ.HashChunkCount()) + pos,
			)
			typ = typ.child(part)
			offset = start
		}
	}

	return gIndex, offset, nil
}

// getBaseIndex returns the base index for a given SSZ type.
// For list types, it returns 2, for all other types it returns 1.
func getBaseIndex(typ SSZType) uint64 {
	if typ.ID().IsList() {
		return 2
	}
	return 1
}
