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
	"fmt"
	"math"
	"strconv"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/types/types"
)

type vector struct {
	Element SSZType
	length  uint64
}

func Vector(element SSZType, length uint64) SSZType {
	return vector{Element: element, length: length}
}

func Bytes(length uint64) SSZType {
	return Vector(U8(), length)
}

func ByteList(length uint64) SSZType {
	return List(U8(), length)
}

func (v vector) ID() types.Type { return types.Vector }

func (v vector) ItemLength() uint64 { return constants.BytesPerChunk }

func (v vector) ItemPosition(p string) (uint64, uint8, uint8, error) {
	i, err := strconv.ParseUint(p, 10, 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("expected index, got name %s", p)
	}
	start := i * v.Element.ItemLength()
	//#nosec:G701 // todo remove float usage.
	return uint64(math.Floor(float64(start) / constants.BytesPerChunk)),
		uint8(start % constants.BytesPerChunk), uint8(start%32 + v.ItemLength()),
		nil
}

func (v vector) HashChunkCount() uint64 {
	totalBytes := v.Length() * v.Element.ItemLength()
	chunks := (totalBytes + constants.BytesPerChunk - 1) / constants.BytesPerChunk
	return chunks
}

// typ.length describes the length for vector types.
func (v vector) Length() uint64 {
	return v.length
}

func (v vector) ElementType(_ string) SSZType {
	return v.Element
}

func (v vector) IsList() bool {
	return false
}
