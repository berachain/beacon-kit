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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.
//

package bytes

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/hex"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/schema"
	"github.com/prysmaticlabs/gohashtree"
)

const (
	// B48Size represents a 48-byte size.
	B48Size = 48
)

var _ schema.MinimalSSZObject = (*B48)(nil)

// B48 represents a 48-byte fixed-size byte array.
// For SSZ purposes it is serialized a `Vector[Byte, 48]`.
type B48 [48]byte

// ToBytes48 is a utility function that transforms a byte slice into a fixed
// 48-byte array. If the input exceeds 48 bytes, it gets truncated.
func ToBytes48(input []byte) B48 {
	return B48(ExtendToSize(input, B48Size))
}

/* -------------------------------------------------------------------------- */
/*                                TextMarshaler                               */
/* -------------------------------------------------------------------------- */

// MarshalText implements the encoding.TextMarshaler interface for B48.
func (h B48) MarshalText() ([]byte, error) {
	return []byte(h.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for B48.
func (h *B48) UnmarshalText(text []byte) error {
	return UnmarshalTextHelper(h[:], text)
}

// String returns the hex string representation of B48.
func (h B48) String() string {
	return hex.FromBytes(h[:]).Unwrap()
}

/* -------------------------------------------------------------------------- */
/*                                JSONMarshaler                               */
/* -------------------------------------------------------------------------- */

// UnmarshalJSON implements the json.Unmarshaler interface for B48.
func (h *B48) UnmarshalJSON(input []byte) error {
	return unmarshalJSONHelper(h[:], input)
}

/* -------------------------------------------------------------------------- */
/*                                SSZMarshaler                                */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the size of its SSZ encoding in bytes.
func (h B48) SizeSSZ() int {
	return B48Size
}

// MarshalSSZ implements the SSZ marshaling for B48.
func (h B48) MarshalSSZ() ([]byte, error) {
	return h[:], nil
}

// IsFixed returns true if the length of the B48 is fixed.
func (h B48) IsFixed() bool {
	return true
}

// Type returns the type of the B48.
func (h B48) Type() schema.SSZType {
	return schema.B48()
}

func (h B48) HashTreeRoot() ([32]byte, error) {
	//nolint:mnd // for a tree height of 1 we need 2 working chunks.
	result := make([][32]byte, 2)
	copy(result[0][:], h[:32])
	copy(result[1][:], h[32:48])
	gohashtree.HashChunks(result, result)
	return result[0], nil
}
