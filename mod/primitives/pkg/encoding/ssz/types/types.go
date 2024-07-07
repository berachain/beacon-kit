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

package types

type Type uint8

const (
	Basic Type = iota
	Elements
	Container
)

// IsBasic returns true if the type is a basic type.
func (t Type) IsBasic() bool {
	return t == Basic
}

// IsComposite returns true if the type is a composite type.
func (t Type) IsComposite() bool {
	return t == Elements || t == Container
}

// MinimalSSZType is the smallest interface of an SSZable type.
type MinimalSSZType interface {
	// MarshalSSZ marshals the type into SSZ format.
	IsFixed() bool
	// Type returns the type of the SSZ object.
	Type() Type
	// SizeSSZ returns the size of the type in bytes.
	SizeSSZ() int
	// ItemLength returns the number of bytes required to represent the type.
	ItemLength() uint64

	// TODO: Do we want these off the minimal?
	//
	// HashTreeRoot returns the hash tree root of the composite type.
	HashTreeRoot() ([32]byte, error)
	// MarshalSSZ marshals the type into SSZ format.
	MarshalSSZ() ([]byte, error)
}

// SSZType is the interface for all SSZ types.
type SSZType[T any] interface {
	MinimalSSZType
	// ChunkCount returns the number of chunks required to store the type.
	ChunkCount() uint64
	NewFromSSZ([]byte) (T, error)
}

// SSZEnumerable is the interface for all SSZ enumerable types must implement.
type SSZEnumerable[
	ElementT any,
] interface {
	MinimalSSZType
	// N returns the N value as defined in the SSZ specification.
	N() uint64
	// Elements returns the elements of the enumerable type.
	Elements() []ElementT
}
