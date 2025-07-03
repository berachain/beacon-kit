// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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

	"github.com/berachain/beacon-kit/primitives/encoding/hex"
	"github.com/prysmaticlabs/gohashtree"
	fastssz "github.com/ferranbt/fastssz"
)

const (
	// B48Size represents a 48-byte size.
	B48Size = 48
)

// B48 represents a 48-byte fixed-size byte array.
// For SSZ purposes it is serialized a `Vector[Byte, 48]`.
type B48 [48]byte

// ToBytes48 is a utility function that transforms a byte slice into a fixed
// 48-byte array. It errs if input has not the required size.
func ToBytes48(input []byte) (B48, error) {
	if len(input) != B48Size {
		return B48{}, fmt.Errorf(
			"%w, got %d, expected %d",
			ErrIncorrectLength,
			len(input),
			B48Size,
		)
	}
	return B48(input), nil
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
	return hex.EncodeBytes(h[:])
}

/* -------------------------------------------------------------------------- */
/*                                JSONMarshaler                               */
/* -------------------------------------------------------------------------- */

// UnmarshalJSON implements the json.Unmarshaler interface for B48.
func (h *B48) UnmarshalJSON(input []byte) error {
	return UnmarshalJSONHelper(h[:], input)
}

/* -------------------------------------------------------------------------- */
/*                                SSZMarshaler                                */
/* -------------------------------------------------------------------------- */

// MarshalSSZ implements the SSZ marshaling for B48.
func (h B48) MarshalSSZ() ([]byte, error) {
	return h[:], nil
}

func (h B48) HashTreeRoot() B32 {
	//nolint:mnd // for a tree height of 1 we need 2 working chunks.
	result := make([][32]byte, 2)
	copy(result[0][:], h[:32])
	copy(result[1][:], h[32:48])
	gohashtree.HashChunks(result, result)
	return result[0]
}

/* -------------------------------------------------------------------------- */
/*                              FastSSZ Methods                               */
/* -------------------------------------------------------------------------- */

// MarshalSSZTo ssz marshals the B48 object to a target array.
func (h B48) MarshalSSZTo(buf []byte) ([]byte, error) {
	dst := buf
	dst = append(dst, h[:]...)
	return dst, nil
}

// UnmarshalSSZ ssz unmarshals the B48 object.
func (h *B48) UnmarshalSSZ(buf []byte) error {
	if len(buf) != B48Size {
		return fastssz.ErrSize
	}
	copy(h[:], buf)
	return nil
}

// SizeSSZ returns the ssz encoded size in bytes for the B48 object.
func (h B48) SizeSSZ() int {
	return B48Size
}

// HashTreeRootWith ssz hashes the B48 object with a hasher.
func (h B48) HashTreeRootWith(hh fastssz.HashWalker) error {
	indx := hh.Index()
	hh.PutBytes(h[:])
	hh.Merkleize(indx)
	return nil
}

// GetTree ssz hashes the B48 object.
func (h *B48) GetTree() (*fastssz.Node, error) {
	return fastssz.ProofTree(h)
}
