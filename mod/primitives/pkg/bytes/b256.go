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
	"github.com/prysmaticlabs/gohashtree"
)

const (
	// B256Size represents a 256-byte size.
	B256Size = 256
)

// B256 represents a 256-byte fixed-size byte array.
// For SSZ purposes it is serialized a `Vector[Byte, 256]`.
type B256 [256]byte

// ToBytes256 is a utility function that transforms a byte slice into a fixed
// 256-byte array. If the input exceeds 256 bytes, it gets truncated.
func ToBytes256(input []byte) B256 {
	return B256(ExtendToSize(input, B256Size))
}

/* -------------------------------------------------------------------------- */
/*                                TextMarshaler                               */
/* -------------------------------------------------------------------------- */

// MarshalText implements the encoding.TextMarshaler interface for B256.
func (h B256) MarshalText() ([]byte, error) {
	return []byte(h.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for B256.
func (h *B256) UnmarshalText(text []byte) error {
	return UnmarshalTextHelper(h[:], text)
}

// String returns the hex string representation of B256.
func (h *B256) String() string {
	return hex.FromBytes(h[:]).Unwrap()
}

/* -------------------------------------------------------------------------- */
/*                                JSONMarshaler                               */
/* -------------------------------------------------------------------------- */

// UnmarshalJSON implements the json.Unmarshaler interface for B256.
func (h *B256) UnmarshalJSON(input []byte) error {
	return unmarshalJSONHelper(h[:], input)
}

/* -------------------------------------------------------------------------- */
/*                                SSZMarshaler                                */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the size of its SSZ encoding in bytes.
func (h B256) SizeSSZ() uint32 {
	return B256Size
}

// MarshalSSZ implements the SSZ marshaling for B256.
func (h B256) MarshalSSZ() ([]byte, error) {
	return h[:], nil
}

// HashTreeRoot returns the hash tree root of the B256.
func (h B256) HashTreeRoot() (B32, error) {
	//nolint:mnd // for a tree height of 3 we need 8 working chunks.
	result := make([][32]byte, 8)
	copy(result[0][:], h[:32])
	copy(result[1][:], h[32:64])
	copy(result[2][:], h[64:96])
	copy(result[3][:], h[96:128])
	copy(result[4][:], h[128:160])
	copy(result[5][:], h[160:192])
	copy(result[6][:], h[192:224])
	copy(result[7][:], h[224:256])
	gohashtree.HashChunks(result, result)
	gohashtree.HashChunks(result, result)
	gohashtree.HashChunks(result, result)
	return result[0], nil
}
