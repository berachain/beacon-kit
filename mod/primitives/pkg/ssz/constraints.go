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

type Base[BaseT any] interface {
	// NewFromSSZ creates a new composite type from an SSZ byte slice.
	NewFromSSZ([]byte) (BaseT, error)
	// MarshalSSZ serializes the composite type to an SSZ byte slice.
	MarshalSSZ() ([]byte, error)
	// SizeSSZ returns the size of the composite type when serialized.
	SizeSSZ() int
	// HashTreeRoot returns the hash tree root of the composite type.
	HashTreeRoot() ([32]byte, error)
}

// Basic defines the interface for a basic type.
type Basic[BasicT any] interface {
	Base[BasicT]
	// Then we add an additional restriction to the following:
	~bool | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 /* TODO: 128, 256 */
}

// Composite defines the interface for a composite type.
type Composite[CompositeT any] interface {
	Base[CompositeT]
	IsFixed() bool
}
