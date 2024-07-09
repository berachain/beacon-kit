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
	"fmt"
	"strconv"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/merkle/proof"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/types/types"
)

// TypeDef represents an interface for Simple Serialize (SSZ) types.
type TypeDef interface {
	// ID returns the type identifier for the SSZ type.
	ID() types.Type
	// ItemLength returns the length of an item in bytes for the SSZ type.
	ItemLength() uint64
	// ItemPosition calculates the position of an item within the SSZ type.
	// It returns the generalized index, start offset, end offset, and any error
	// encountered.
	ItemPosition(p string) (uint64, uint8, uint8, error)
	// ElementType returns the SSZ type of the element at the given path.
	ElementType(p string) TypeDef
	// HashChunkCount returns the number of 32-byte chunks required to represent
	// the SSZ type in a Merkle tree.
	HashChunkCount() uint64
}

/* -------------------------------------------------------------------------- */
/*                                    Basic                                   */
/* -------------------------------------------------------------------------- */

// basic represents a basic SSZ type.
type basic uint64

// ID returns the type ID of the basic type.
func (b basic) ID() types.Type { return types.Basic }

// ItemLength returns the size of the basic type in bytes.
func (b basic) ItemLength() uint64 { return uint64(b) }

// position always returns an error for basic types, as they have no children.
func (b basic) ItemPosition(_ string) (uint64, uint8, uint8, error) {
	return 0, 0, 0, errors.New("basic type has no children")
}

// child returns the basic type itself, as it has no children.
func (b basic) ElementType(_ string) TypeDef { return b }

// Chunks returns the number of 32-byte chunks required to represent the basic
// type.
func (b basic) HashChunkCount() uint64 { return 1 }

/* -------------------------------------------------------------------------- */
/*                                   Vector                                   */
/* -------------------------------------------------------------------------- */

type vector struct {
	elementType TypeDef
	length      uint64
}

func Vector(elementType TypeDef, length uint64) TypeDef {
	return vector{elementType: elementType, length: length}
}

func ByteVector(length uint64) TypeDef {
	return Vector(U8(), length)
}

func (v vector) ID() types.Type { return types.Vector }

func (v vector) ItemLength() uint64 { return constants.BytesPerChunk }

func (v vector) ItemPosition(p string) (uint64, uint8, uint8, error) {
	i, err := strconv.ParseUint(p, 10, 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("expected index, got name %s", p)
	}
	start := i * v.elementType.ItemLength()
	return start / constants.BytesPerChunk,
		//#nosec:G701 // can't overflow.
		uint8(start % constants.BytesPerChunk),
		//#nosec:G701 // can't overflow.
		uint8(start%constants.BytesPerChunk + v.ItemLength()), nil
}

func (v vector) HashChunkCount() uint64 {
	totalBytes := v.Length() * v.elementType.ItemLength()
	chunks := (totalBytes + constants.BytesPerChunk - 1) / constants.BytesPerChunk
	return chunks
}

// typ.length describes the length for vector types.
func (v vector) Length() uint64 {
	return v.length
}

func (v vector) ElementType(_ string) TypeDef {
	return v.elementType
}

/* -------------------------------------------------------------------------- */
/*                                    List                                    */
/* -------------------------------------------------------------------------- */

// List Type.
type list struct {
	elementType TypeDef
	limit       uint64
}

func List(elementType TypeDef, limit uint64) TypeDef {
	return list{elementType: elementType, limit: limit}
}

func ByteList(limit uint64) TypeDef {
	return List(U8(), limit)
}

func (l list) ID() types.Type { return types.List }

func (l list) ItemLength() uint64 { return l.elementType.ItemLength() }

func (l list) HashChunkCount() uint64 {
	totalBytes := l.Length() * l.elementType.ItemLength()
	chunks := (totalBytes + constants.BytesPerChunk - 1) / constants.BytesPerChunk
	return chunks
}

func (l list) ElementType(_ string) TypeDef {
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
	return start / constants.BytesPerChunk,
		//#nosec:G701 // can't overflow.
		uint8(start % constants.BytesPerChunk),
		//#nosec:G701 // can't overflow.
		uint8(start%constants.BytesPerChunk + l.ItemLength()), nil
}

/* -------------------------------------------------------------------------- */
/*                                  Container                                 */
/* -------------------------------------------------------------------------- */

type container struct {
	Fields     []TypeDef
	FieldIndex map[string]uint64
}

func Field(name string, typ TypeDef) *proof.Field[TypeDef] {
	return proof.NewField(name, typ)
}

func Container(fields ...*proof.Field[TypeDef]) TypeDef {
	fieldIndex := make(map[string]uint64)
	types := make([]TypeDef, len(fields))
	for i, f := range fields {
		//#nosec:G701 // todo fix.
		fieldIndex[f.GetName()] = uint64(i)
		types[i] = f.GetValue()
	}
	return container{Fields: types, FieldIndex: fieldIndex}
}

func (c container) ID() types.Type { return types.Container }

func (c container) ItemLength() uint64 { return constants.BytesPerChunk }

func (c container) ItemPosition(p string) (uint64, uint8, uint8, error) {
	pos, ok := c.FieldIndex[p]
	if !ok {
		return 0, 0, 0, fmt.Errorf("field %s not found", p)
	}
	//#nosec:G701 // can't overflow.
	return pos, 0, uint8(c.Fields[pos].ItemLength()), nil
}

func (c container) ElementType(p string) TypeDef {
	return c.Fields[c.FieldIndex[p]]
}

func (c container) Length() uint64 { return uint64(len(c.Fields)) }

func (c container) HashChunkCount() uint64 { return uint64(len(c.Fields)) }
