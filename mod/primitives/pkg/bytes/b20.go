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
)

const (
	// B20Size represents a 20-byte size.
	B20Size = 20
)

var _ schema.MinimalSSZObject = (*B20)(nil)

// B20 represents a 20-byte fixed-size byte array.
// For SSZ purposes it is serialized a `Vector[Byte, 20]`.
type B20 [20]byte

// ToBytes20 is a utility function that transforms a byte slice into a fixed
// 20-byte array. If the input exceeds 20 bytes, it gets truncated.
func ToBytes20(input []byte) B20 {
	return B20(ExtendToSize(input, B20Size))
}

/* -------------------------------------------------------------------------- */
/*                                TextMarshaler                               */
/* -------------------------------------------------------------------------- */

// MarshalText implements the encoding.TextMarshaler interface for B20.
func (h B20) MarshalText() ([]byte, error) {
	return []byte(h.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for B20.
func (h *B20) UnmarshalText(text []byte) error {
	return UnmarshalTextHelper(h[:], text)
}

// String returns the hex string representation of B20.
func (h *B20) String() string {
	return hex.FromBytes(h[:]).Unwrap()
}

/* -------------------------------------------------------------------------- */
/*                                JSONMarshaler                               */
/* -------------------------------------------------------------------------- */

// UnmarshalJSON implements the json.Unmarshaler interface for B20.
func (h *B20) UnmarshalJSON(input []byte) error {
	return unmarshalJSONHelper(h[:], input)
}

/* -------------------------------------------------------------------------- */
/*                                SSZMarshaler                                */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the size of its SSZ encoding in bytes.
func (h B20) SizeSSZ() int {
	return B20Size
}

// MarshalSSZ implements the SSZ marshaling for B20.
func (h B20) MarshalSSZ() ([]byte, error) {
	return h[:], nil
}

// IsFixed returns true if the length of the B20 is fixed.
func (h B20) IsFixed() bool {
	return true
}

// Type returns the type of the B20.
func (h B20) Type() schema.SSZType {
	return schema.B20()
}

// HashTreeRoot returns the hash tree root of the B20.
func (h B20) HashTreeRoot() ([32]byte, error) {
	var result [32]byte
	copy(result[:], h[:])
	return result, nil
}
