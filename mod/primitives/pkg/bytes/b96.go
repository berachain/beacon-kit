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
	"fmt"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/hex"
	"github.com/prysmaticlabs/gohashtree"
)

const (
	// B96Size represents a 96-byte size.
	B96Size = 96
)

// B96 represents a 96-byte fixed-size byte array.
// For SSZ purposes it is serialized a `Vector[Byte, 96]`.
type B96 [96]byte

// ToBytes96 is a utility function that transforms a byte slice into a fixed
// 96-byte array. It errs if input has not the required size.
func ToBytes96(input []byte) (B96, error) {
	if len(input) != B96Size {
		return B96{}, fmt.Errorf(
			"%w, got %d, expected %d",
			ErrIncorrectLength,
			len(input),
			B96Size,
		)
	}
	return B96(input), nil
}

/* -------------------------------------------------------------------------- */
/*                                TextMarshaler                               */
/* -------------------------------------------------------------------------- */

// MarshalText implements the encoding.TextMarshaler interface for B96.
func (h B96) MarshalText() ([]byte, error) {
	return []byte(h.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for B96.
func (h *B96) UnmarshalText(text []byte) error {
	return UnmarshalTextHelper(h[:], text)
}

// String returns the hex string representation of B96.
func (h *B96) String() string {
	return hex.EncodeBytes(h[:])
}

/* -------------------------------------------------------------------------- */
/*                                JSONMarshaler                               */
/* -------------------------------------------------------------------------- */

// UnmarshalJSON implements the json.Unmarshaler interface for B96.
func (h *B96) UnmarshalJSON(input []byte) error {
	return UnmarshalJSONHelper(h[:], input)
}

/* -------------------------------------------------------------------------- */
/*                                SSZMarshaler                                */
/* -------------------------------------------------------------------------- */

// MarshalSSZ implements the SSZ marshaling for B96.
func (h B96) MarshalSSZ() ([]byte, error) {
	return h[:], nil
}

// HashTreeRoot returns the hash tree root of the B96.
func (h B96) HashTreeRoot() B32 {
	//nolint:mnd // for a tree height of 2 we need 4 working chunks.
	result := make([][32]byte, 4)
	copy(result[0][:], h[:32])
	copy(result[1][:], h[32:64])
	copy(result[2][:], h[64:96])
	gohashtree.HashChunks(result, result)
	gohashtree.HashChunks(result, result)
	return result[0]
}
