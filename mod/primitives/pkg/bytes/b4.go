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
//
//nolint:dupl // it's okay to duplicate the code for different types
package bytes

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/hex"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/schema"
)

const (
	// B4Size represents a 4-byte size.
	B4Size = 4
)

var _ schema.MinimalSSZObject = (*B4)(nil)

// B4 represents a 4-byte fixed-size byte array.
// For SSZ purposes it is serialized a `Vector[Byte, 4]`.
type B4 [4]byte

// ToBytes4 is a utility function that transforms a byte slice into a fixed
// 4-byte array. If the input exceeds 4 bytes, it gets truncated.
func ToBytes4(input []byte) B4 {
	return B4(ExtendToSize(input, B4Size))
}

/* -------------------------------------------------------------------------- */
/*                                TextMarshaler                               */
/* -------------------------------------------------------------------------- */

// MarshalText implements the encoding.TextMarshaler interface for B4.
func (h B4) MarshalText() ([]byte, error) {
	return []byte(h.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for B4.
func (h *B4) UnmarshalText(text []byte) error {
	return UnmarshalTextHelper(h[:], text)
}

// String returns the hex string representation of B4.
func (h B4) String() string {
	return hex.FromBytes(h[:]).Unwrap()
}

/* -------------------------------------------------------------------------- */
/*                                JSONMarshaler                               */
/* -------------------------------------------------------------------------- */

// UnmarshalJSON implements the json.Unmarshaler interface for B4.
func (h *B4) UnmarshalJSON(input []byte) error {
	return unmarshalJSONHelper(h[:], input)
}

/* -------------------------------------------------------------------------- */
/*                                SSZMarshaler                                */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the size of its SSZ encoding in bytes.
func (h B4) SizeSSZ() int {
	return B4Size
}

// MarshalSSZ implements the SSZ marshaling for B4.
func (h B4) MarshalSSZ() ([]byte, error) {
	return h[:], nil
}

// IsFixed returns true if the length of the B4 is fixed.
func (h B4) IsFixed() bool {
	return true
}

// Type returns the type of the B4.
func (h B4) Type() schema.SSZType {
	return schema.B4()
}

// HashTreeRoot returns the hash tree root of the B4.
func (h B4) HashTreeRoot() ([32]byte, error) {
	var result [32]byte
	copy(result[:], h[:])
	return result, nil
}
