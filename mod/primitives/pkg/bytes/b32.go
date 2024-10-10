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
)

const (
	// B32Size represents a 32-byte size.
	B32Size = 32
)

// B32 represents a 32-byte fixed-size byte array.
// For SSZ purposes it is serialized a `Vector[Byte, 32]`.
type B32 [32]byte

// ToBytes32 is a utility function that transforms a byte slice into a fixed
// 32-byte array. If the input exceeds 32 bytes, it gets truncated.
func ToBytes32(input []byte) B32 {
	var b B32
	copy(b[:], input)
	return b
}

/* -------------------------------------------------------------------------- */
/*                                TextMarshaler                               */
/* -------------------------------------------------------------------------- */

// MarshalText implements the encoding.TextMarshaler interface for B32.
func (h B32) MarshalText() ([]byte, error) {
	return []byte(h.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for B32.
func (h *B32) UnmarshalText(text []byte) error {
	return UnmarshalTextHelper(h[:], text)
}

// String returns the hex string representation of B32.
func (h B32) String() string {
	return hex.EncodeBytes(h[:])
}

/* -------------------------------------------------------------------------- */
/*                                JSONMarshaler                               */
/* -------------------------------------------------------------------------- */

// UnmarshalJSON implements the json.Unmarshaler interface for B32.
func (h *B32) UnmarshalJSON(input []byte) error {
	return UnmarshalJSONHelper(h[:], input)
}

/* -------------------------------------------------------------------------- */
/*                                SSZMarshaler                                */
/* -------------------------------------------------------------------------- */

// MarshalSSZ implements the SSZ marshaling for B32.
func (h B32) MarshalSSZ() ([]byte, error) {
	return h[:], nil
}

// HashTreeRoot returns the hash tree root of the B32.
func (h B32) HashTreeRoot() B32 {
	return h
}
