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

package ssz

import (
	"fmt"
	"reflect"

	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/merkleizer"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/tree"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/types"
)

/* -------------------------------------------------------------------------- */
/*                                Type Definitions                            */
/* -------------------------------------------------------------------------- */

// Vector conforms to the SSZEenumerable interface.
var _ types.SSZEnumerable[types.MinimalSSZType] = (*Container)(nil)

type Container struct {
	fieldIndex map[string]int
	elements   []types.MinimalSSZType
}

type ContainerField struct {
	Name  string
	Value types.MinimalSSZType
}

// ContainerFromElements creates a new Container from elements.
func ContainerFromElements(elements ...types.MinimalSSZType) *Container {
	return &Container{
		elements: elements,
	}
}

func ContainerFromFields(fields []ContainerField) *Container {
	elements := make([]types.MinimalSSZType, len(fields))
	fieldIndex := make(map[string]int)
	for i, field := range fields {
		elements[i] = field.Value
		fieldIndex[field.Name] = i
	}
	return &Container{
		elements:   elements,
		fieldIndex: fieldIndex,
	}
}

// NewContainer creates a new Container from any struct, using reflection to get
// all the fields and put them into the elements list.
func NewContainer(v interface{}) (*Container, error) {
	var (
		val        = reflect.ValueOf(v)
		typ        = reflect.TypeOf(v)
		fieldIndex = make(map[string]int)
		elements   []types.MinimalSSZType
		j          int
	)

	// If v is a pointer, get the value it points to
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	// Ensure v is a struct
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input must be a struct or pointer to struct")
	}

	for i := range flattenStructFields(typ) {
		field := typ.Field(i)
		path := field.Tag.Get("ssz-path")
		if path == "" {
			continue
		}

		fieldValue := val.Field(i)
		if field.Type.Kind() == reflect.Ptr {
			fieldValue = fieldValue.Elem()
		}

		if fieldValue.Kind() == reflect.Struct {
			// Recursively add the fields of the struct
			container, err := NewContainer(fieldValue.Interface())
			if err != nil {
				return nil, err
			}
			elements = append(elements, container)
			fieldIndex[path] = j
		} else {
			switch sszType := fieldValue.Interface().(type) {
			case types.MinimalSSZType:
				elements = append(elements, sszType)
				fieldIndex[path] = j
			}
		}

		// else if sszType, ok := fieldValue.Interface().(types.MinimalSSZType); ok {
		// 	elements = append(elements, sszType)
		// 	fieldIndex[path] = j
		// } else {
		// 	return nil, fmt.Errorf("field %s does not implement MinimalSSZType",
		// 		val.Type().Field(i).Name)
		// }
		j++
	}

	return &Container{elements: elements, fieldIndex: fieldIndex}, nil
}

func flattenStructFields(typ reflect.Type) []reflect.StructField {
	var fields []reflect.StructField
	for i := range typ.NumField() {
		field := typ.Field(i)
		if field.Anonymous {
			// flatten embedded struct fields
			embedded := flattenStructFields(field.Type)
			fields = append(fields, embedded...)
		} else {
			fields = append(fields, field)
		}
	}
	return fields
}

/* -------------------------------------------------------------------------- */
/*                                 BaseSSZType                                */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the size of the container in bytes.
func (c *Container) SizeSSZ() int {
	size := 0
	for _, element := range c.elements {
		size += element.SizeSSZ()
	}
	return size
}

// IsFixed returns true if the container is fixed size.
func (c *Container) IsFixed() bool {
	for _, element := range c.elements {
		if !element.IsFixed() {
			return false
		}
	}
	return true
}

// N returns the N value as defined in the SSZ specification.
func (c *Container) N() uint64 {
	return uint64(len(c.elements))
}

// Type returns the type of the container.
func (*Container) Type() types.Type {
	return types.Composite
}

// ChunkCount returns the number of chunks in the container.
func (c *Container) ChunkCount() uint64 {
	return c.N()
}

// Elements returns the elements of the container.
func (c *Container) Elements() []types.MinimalSSZType {
	return c.elements
}

/* -------------------------------------------------------------------------- */
/*                                Merkleization                               */
/* -------------------------------------------------------------------------- */

// HashTreeRoot returns the hash tree root of the container.
func (c *Container) HashTreeRootWith(
	merkleizer VectorMerkleizer[[32]byte, types.MinimalSSZType],
) ([32]byte, error) {
	return merkleizer.MerkleizeVectorCompositeOrContainer(c.elements)
}

// HashTreeRoot returns the hash tree root of the container.
func (c *Container) HashTreeRoot() ([32]byte, error) {
	return c.HashTreeRootWith(merkleizer.New[[32]byte, types.MinimalSSZType]())
}

/* -------------------------------------------------------------------------- */
/*                                Serialization                               */
/* -------------------------------------------------------------------------- */

// MarshalSSZToBytes marshals the VectorBasic into SSZ format.
func (c *Container) MarshalSSZTo(_ []byte) ([]byte, error) {
	return nil, errors.New("not implemented yet")
}

// MarshalSSZ marshals the VectorBasic into SSZ format.
func (c *Container) MarshalSSZ() ([]byte, error) {
	return c.MarshalSSZTo(make([]byte, 0, c.SizeSSZ()))
}

// NewFromSSZ creates a new Container from SSZ format.
func (c *Container) NewFromSSZ(_ []byte) (*Container, error) {
	return nil, errors.New("not implemented yet")
}

type Schema struct {
	cache map[string]uint64
}

func (c *Container) GIndex(gIndex math.U64, path tree.ObjectPath) *tree.Node {
	head, rest := path.Head()
	if index, ok := c.fieldIndex[head]; ok {
		gIndex = gIndex*math.U64(c.N()).NextPowerOfTwo() + math.U64(index)
		field := c.elements[index]
		if gid, ok := field.(tree.GIndexed); ok {
			return gid.GIndex(gIndex, rest)
		} else {
			// TODO calc offset
			return &tree.Node{GIndex: gIndex}
		}
	}
	return nil
}

/*
func (c *Container) Default() *Container {
	for i, element := range c.elements {
		if enum, ok := element.(types.SSZEnumerable[types.MinimalSSZType]); ok {
			c.elements[i] = enum.Default()
		}
	}
}
*/
