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
//nolint:dupl // it's okay to duplicate code for different types
package math

import (
	"encoding/binary"
	"fmt"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/schema"
)

/* -------------------------------------------------------------------------- */
/*                                Type Definitions                            */
/* -------------------------------------------------------------------------- */

// Ensure types implement schema.SSZObject.
var _ schema.SSZObject[U32] = (*U32)(nil)

// U32 represents a 32-bit unsigned integer that is both SSZ and JSON.
type U32 uint32

/* -------------------------------------------------------------------------- */
/*                                     U32                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the size of the uint32 in bytes.
func (U32) SizeSSZ() int {
	return constants.U32Size
}

// MarshalSSZ marshals the uint32 into SSZ format.
func (u U32) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, constants.U32Size)
	binary.LittleEndian.PutUint32(buf, uint32(u))
	return buf, nil
}

// NewFromSSZ creates a new U32 from SSZ format.
func (U32) NewFromSSZ(buf []byte) (U32, error) {
	if len(buf) != constants.U32Size {
		return 0, fmt.Errorf(
			"invalid buffer length: expected %d, got %d",
			constants.U32Size,
			len(buf),
		)
	}
	return U32(binary.LittleEndian.Uint32(buf)), nil
}

// HashTreeRoot returns the hash tree root of the uint32.
func (u U32) HashTreeRoot() ([32]byte, error) {
	buf := make([]byte, constants.BytesPerChunk)
	binary.LittleEndian.PutUint32(buf[:constants.U32Size], uint32(u))
	return [32]byte(buf), nil
}

// IsFixed returns true if the bool is fixed size.
func (U32) IsFixed() bool {
	return true
}

// Type returns the type of the U32.
func (U32) Type() schema.SSZType {
	return schema.U32()
}

// ChunkCount returns the number of chunks required to store the uint32.
func (U32) ChunkCount() uint64 {
	return 1
}
