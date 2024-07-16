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
	"errors"
	"strings"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/schema"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math/pow"
)

// ObjectPath represents a path to an object in a Merkle tree.
type ObjectPath[GeneralizedIndexT ~uint64, RootT ~[32]byte] string

// Split returns the path split by "/".
func (p ObjectPath[_, _]) Split() []string {
	return strings.Split(string(p), "/")
}

// GetGeneralizedIndex converts a path to a generalized index representing its
// position in the Merkle tree.
//
//	"""
//
// Converts a path (eg. `[7, "foo", 3]` for `x[7].foo[3]`, `[12, "bar",
// "__len__"]` for `len(x[12].bar)`) into the generalized index representing its
// position in the Merkle tree.
// """
// root = GeneralizedIndex(1)
// for p in path:
//
//	assert not issubclass(typ, BasicValue)  # If we descend to a basic type, the
//
// path cannot continue further
//
//	if p == '__len__':
//		typ = uint64
//		assert issubclass(typ, (List, ByteList))
//		root = GeneralizedIndex(root * 2 + 1)
//	else:
//		pos, _, _ = get_item_position(typ, p)
//		base_index = (GeneralizedIndex(2) if issubclass(typ, (List, ByteList)) else
//
// GeneralizedIndex(1)) 		root = GeneralizedIndex(root * base_index *
// get_power_of_two_ceil(chunk_count(typ)) + pos)
//
//	typ = get_elem_type(typ, p)
//
// return root.
func (p ObjectPath[GeneralizedIndexT, RootT]) GetGeneralizedIndex(
	typ schema.SSZType,
) (schema.SSZType, GeneralizedIndexT, uint8, error) {
	gIndex := uint64(1)
	offset := uint8(0)
	for _, part := range p.Split() {
		// If we descend to a basic type, the path cannot continue further
		if typ.ID().IsBasic() {
			return nil, 0, 0, errors.New(
				"cannot descend further from basic type",
			)
		}

		if part == "__len__" {
			if !typ.ID().IsList() {
				return nil, 0, 0, errors.New(
					"__len__ is only valid for List types",
				)
			}
			typ = schema.U64()
			//nolint:mnd // from spec.
			gIndex = gIndex*2 + 1
		} else {
			pos, start, _, err := typ.ItemPosition(part)
			if err != nil {
				return nil, 0, 0, err
			}

			gIndex = gIndex*getBaseIndex(typ)*
				pow.NextPowerOfTwo(typ.HashChunkCount()) + pos
			typ = typ.ElementType(part)
			offset = start
		}
	}

	return typ, GeneralizedIndexT(gIndex), offset, nil
}

// getBaseIndex returns the base index for a given SSZ type.
// For list types, it returns 2, for all other types it returns 1.
func getBaseIndex(typ schema.SSZType) uint64 {
	if typ.ID().IsList() {
		//nolint:mnd // 2 is allowed.
		return 2
	}
	return 1
}
