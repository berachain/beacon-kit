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

package math

import (
	"fmt"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/schema"
)

// Ensure type implements schema.SSZObject.
var _ schema.SSZObject[Bool] = (*Bool)(nil)

type Bool bool

// SizeSSZ returns the size of the bool in bytes.
func (Bool) SizeSSZ() int {
	return constants.BoolSize
}

// MarshalSSZ marshals the bool into SSZ format.
func (b Bool) MarshalSSZ() ([]byte, error) {
	if b {
		return []byte{1}, nil
	}
	return []byte{0}, nil
}

// NewFromSSZ creates a new Bool from SSZ format.
func (Bool) NewFromSSZ(buf []byte) (Bool, error) {
	if len(buf) != constants.BoolSize {
		return false, fmt.Errorf(
			"invalid buffer length: expected %d, got %d",
			constants.BoolSize,
			len(buf),
		)
	}
	return Bool(buf[0] != 0), nil
}

// HashTreeRoot returns the hash tree root of the bool.
func (b Bool) HashTreeRoot() ([32]byte, error) {
	buf := make([]byte, constants.BytesPerChunk)
	if b {
		buf[0] = 1
	}
	return [constants.BytesPerChunk]byte(buf), nil
}

// IsFixed returns true if the bool is fixed size.
func (Bool) IsFixed() bool {
	return true
}

// Type returns the type of the bool.
func (Bool) Type() schema.SSZType {
	return schema.Bool()
}

// ChunkCount returns the number of chunks required to store the bool.
func (Bool) ChunkCount() uint64 {
	return 1
}
