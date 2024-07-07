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
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/merkleizer"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/proof"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/types"
)

/* -------------------------------------------------------------------------- */
/*                                Type Definitions                            */
/* -------------------------------------------------------------------------- */

// Vector conforms to the SSZEenumerable interface.
var _ types.SSZEnumerable[proof.Field] = (*Container)(nil)

type Container struct {
	fields     []proof.Field
	fieldIndex map[string]uint64
}

// ContainerFromElements creates a new Container from elements.
func ContainerFromElements(fields ...proof.Field) *Container {
	fieldIndex := make(map[string]uint64)
	types := make([]proof.SSZType, len(fields))
	for i, f := range fields {
		fieldIndex[f.GetName()] = uint64(i)
		types[i] = f
	}

	return &Container{
		fields:     fields,
		fieldIndex: fieldIndex,
	}
}

/* -------------------------------------------------------------------------- */
/*                                 BaseSSZType                                */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the size of the container in bytes.
func (c *Container) SizeSSZ() int {
	size := 0
	for _, element := range c.fields {
		size += element.SizeSSZ()
	}
	return size
}

// IsFixed returns true if the container is fixed size.
func (c *Container) IsFixed() bool {
	for _, element := range c.fields {
		if !element.IsFixed() {
			return false
		}
	}
	return true
}

// N returns the N value as defined in the SSZ specification.
func (c *Container) N() uint64 {
	return uint64(len(c.fields))
}

// Type returns the type of the container.
func (*Container) Type() types.Type {
	return types.Container
}

// ChunkCount returns the number of chunks in the container.
func (c *Container) ChunkCount() uint64 {
	return c.N()
}

// Elements returns the elements of the container.
func (c *Container) Elements() []proof.Field {
	return c.fields
}

// GetFieldByName returns the field with the given name.
func (c *Container) GetFieldByName(name string) proof.Field {
	return c.fields[c.fieldIndex[name]]
}

// GetFieldIndex returns the index of the field with the given name.
func (c *Container) GetFieldIndex(name string) uint64 {
	return c.fieldIndex[name]
}

/* -------------------------------------------------------------------------- */
/*                                Merkleization                               */
/* -------------------------------------------------------------------------- */

// HashTreeRoot returns the hash tree root of the container.
func (c *Container) HashTreeRootWith(
	merkleizer VectorMerkleizer[[32]byte, proof.Field],
) ([32]byte, error) {
	return merkleizer.MerkleizeVectorCompositeOrContainer(c.fields)
}

// HashTreeRoot returns the hash tree root of the container.
func (c *Container) HashTreeRoot() ([32]byte, error) {
	return c.HashTreeRootWith(merkleizer.New[[32]byte, proof.Field]())
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
