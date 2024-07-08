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

// List Type.
type list struct {
	elementType SSZType
	limit       uint64
}

func List(elementType SSZType, limit uint64) SSZType {
	return list{elementType: elementType, limit: limit}
}

func (l list) ID() types.Type { return types.List }

func (l list) ItemLength() uint64 { return l.elementType.ItemLength() }

func (l list) HashChunkCount() uint64 {
	totalBytes := l.Length() * l.elementType.ItemLength()
	chunks := (totalBytes + chunkSize - 1) / chunkSize
	return chunks
}

func (l list) ElementType(_ string) SSZType {
	return l.elementType
}

// typ.length describes the limit for list types.
func (l list) Length() uint64 {
	return l.limit
}

// position returns the chunk index and offset for a given list index.
func (l list) ItemPosition(p string) (uint64, uint8, uint8, error) {
	i, err := strconv.ParseUint(p, 10, 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("expected index, got name %s", p)
	}
	start := i * l.elementType.ItemLength()
	//#nosec:G701 // todo remove float usage.
	return uint64(math.Floor(float64(start) / chunkSize)),
		uint8(start % chunkSize), uint8(start%32 + l.ItemLength()),
		nil
}
