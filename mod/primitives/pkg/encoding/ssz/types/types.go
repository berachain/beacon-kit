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

import (
	"go/types"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/schema"
)

/* -------------------------------------------------------------------------- */
/*                              Type Definitions                              */
/* -------------------------------------------------------------------------- */

// SSZTypeDef defines the interface that type definitions must adhere to.
// An SSZTypeDef *REPRESENTS* an underyling type, but it is NOT an instance
// of this type.
type SSZTypeDef interface {
	// ID returns the type identifier for the SSZ type.
	ID() types.Type
	// ItemLength returns the length of an item in bytes for the SSZ type.
	ItemLength() uint64
	// ItemPosition calculates the position of an item within the SSZ type.
	// It returns the generalized index, start offset, end offset, and any error
	// encountered.
	ItemPosition(p string) (uint64, uint8, uint8, error)
	// ElementType returns the SSZ type of the element at the given path.
	ElementType(p string) SSZTypeDef
	// HashChunkCount returns the number of 32-byte chunks required to represent
	// the SSZ type in a Merkle tree.
	HashChunkCount() uint64
}

/* -------------------------------------------------------------------------- */
/*                                 SSZ Objects                                */
/* -------------------------------------------------------------------------- */

// SSZObject defines an interface for SSZ basic types which includes methods for
// determining the size of the SSZ encoding and computing the hash tree root.
type MerkleizableSSZObject[RootT ~[32]byte] interface {
	// SizeSSZ returns the size in bytes of the SSZ-encoded data.
	SizeSSZ() int
	// HashTreeRoot computes and returns the hash tree root of the data as
	// RootT and an error if the computation fails.
	HashTreeRoot() (RootT, error)
	// MarshalSSZ marshals the data into SSZ format.
	MarshalSSZ() ([]byte, error)
}

// MinimalSSZType is the smallest interface of an SSZable type.
type MinimalSSZType interface {
	MerkleizableSSZObject[[32]byte]
	// MarshalSSZ marshals the type into SSZ format.
	IsFixed() bool
	// Type returns the type of the SSZ object.
	Type() schema.SSZType
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
