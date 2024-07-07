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

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/types/types"
)

const (
	B4Size  = 4
	B8Size  = 8
	B16Size = 16
	B20Size = 20
	B32Size = 32
	B48Size = 48
	B64Size = 64
	B96Size = 96
)

// B4 creates a Vector of 4 bytes (32 bits).
func B4() SSZType {
	return Vector(U8(), B4Size)
}

// B8 creates a Vector of 8 bytes (64 bits).
func B8() SSZType {
	return Vector(U8(), B8Size)
}

// B16 creates a Vector of 16 bytes (128 bits).
func B16() SSZType {
	return Vector(U8(), B16Size)
}

// B20 creates a Vector of 20 bytes (160 bits).
func B20() SSZType {
	return Vector(U8(), B20Size)
}

// B32 creates a Vector of 32 bytes (256 bits).
func B32() SSZType {
	return Vector(U8(), B32Size)
}

// B48 creates a Vector of 48 bytes (384 bits).
func B48() SSZType {
	return Vector(U8(), B48Size)
}

// B64 creates a Vector of 64 bytes (512 bits).
func B64() SSZType {
	return Vector(U8(), B64Size)
}

// B96 creates a Vector of 96 bytes (768 bits).
func B96() SSZType {
	return Vector(U8(), B96Size)
}

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

func (v vector) ItemLength() uint64 { return chunkSize }

func (v vector) HashChunkCount() uint64 {
	totalBytes := v.Length() * v.Element.ItemLength()
	chunks := (totalBytes + chunkSize - 1) / chunkSize
	return chunks
}

func (v vector) child(_ string) SSZType {
	return v.Element
}

func (v vector) Length() uint64 {
	return v.length
}

func (v vector) position(p string) (uint64, uint8, error) {
	i, err := strconv.ParseUint(p, 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("expected index, got name %s", p)
	}
	start := i * v.Element.ItemLength()
	return uint64(math.Floor(float64(start) / chunkSize)),
		uint8(start % chunkSize),
		nil
}

func (v vector) IsList() bool {
	return false
}
