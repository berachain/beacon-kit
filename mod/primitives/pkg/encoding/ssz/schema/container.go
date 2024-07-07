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

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/tree/proof"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/types/types"
)

type container struct {
	Fields     []SSZType
	FieldIndex map[string]uint64
}

func Field(name string, typ SSZType) *proof.Field[SSZType] {
	return proof.NewField(name, typ)
}

func Container(fields ...*proof.Field[SSZType]) SSZType {
	fieldIndex := make(map[string]uint64)
	types := make([]SSZType, len(fields))
	for i, f := range fields {
		fieldIndex[f.GetName()] = uint64(i)
		types[i] = f.GetValue()
	}
	return container{Fields: types, FieldIndex: fieldIndex}
}

func (c container) ID() types.Type { return types.Container }

func (c container) ItemLength() uint64 { return chunkSize }

func (c container) Length() uint64 { return uint64(len(c.Fields)) }

func (c container) HashChunkCount() uint64 { return uint64(len(c.Fields)) }

func (c container) child(p string) SSZType {
	return c.Fields[c.FieldIndex[p]]
}

func (c container) position(p string) (uint64, uint8, error) {
	pos, ok := c.FieldIndex[p]
	if !ok {
		return 0, 0, fmt.Errorf("field %s not found", p)
	}
	return pos, 0, nil
}
