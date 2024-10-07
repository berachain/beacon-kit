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
//nolint:dupl // it's okay.
package bytes

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/hex"
)

const (
	// B8Size represents an 8-byte size.
	B8Size = 8
)

// B8 represents an 8-byte fixed-size byte array.
// For SSZ purposes it is serialized a `Vector[Byte, 8]`.
type B8 [8]byte

// ToBytes8 is a utility function that transforms a byte slice into a fixed
// 8-byte array. If the input exceeds 8 bytes, it gets truncated.
func ToBytes8(input []byte) B8 {
	return B8(ExtendToSize(input, B8Size))
}

/* -------------------------------------------------------------------------- */
/*                                TextMarshaler                               */
/* -------------------------------------------------------------------------- */

// MarshalText implements the encoding.TextMarshaler interface for B8.
func (h B8) MarshalText() ([]byte, error) {
	return []byte(h.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for B8.
func (h *B8) UnmarshalText(text []byte) error {
	return UnmarshalTextHelper(h[:], text)
}

// String returns the hex string representation of B8.
func (h B8) String() string {
	return hex.EncodeBytes(h[:])
}

/* -------------------------------------------------------------------------- */
/*                                JSONMarshaler                               */
/* -------------------------------------------------------------------------- */

// UnmarshalJSON implements the json.Unmarshaler interface for B8.
func (h *B8) UnmarshalJSON(input []byte) error {
	return UnmarshalJSONHelper(h[:], input)
}

/* -------------------------------------------------------------------------- */
/*                                SSZMarshaler                                */
/* -------------------------------------------------------------------------- */

// MarshalSSZ implements the SSZ marshaling for B8.
func (h B8) MarshalSSZ() ([]byte, error) {
	return h[:], nil
}

// HashTreeRoot returns the hash tree root of the B8.
func (h B8) HashTreeRoot() B32 {
	return ToBytes32(h[:])
}
