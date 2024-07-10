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
var _ schema.SSZObject[U16] = (*U16)(nil)

// U16 represents a 16-bit unsigned integer that is both SSZ and JSON.
type U16 uint16

/* -------------------------------------------------------------------------- */
/*                                     U16                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the size of the uint16 in bytes.
func (U16) SizeSSZ() int {
	return constants.U16Size
}

// MarshalSSZ marshals the uint16 into SSZ format.
func (u U16) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, constants.U16Size)
	binary.LittleEndian.PutUint16(buf, uint16(u))
	return buf, nil
}

// NewFromSSZ creates a new U16 from SSZ format.
func (U16) NewFromSSZ(buf []byte) (U16, error) {
	if len(buf) != constants.U16Size {
		return 0, fmt.Errorf(
			"invalid buffer length: expected %d, got %d",
			constants.U16Size,
			len(buf),
		)
	}
	return U16(binary.LittleEndian.Uint16(buf)), nil
}

// HashTreeRoot returns the hash tree root of the uint16.
func (u U16) HashTreeRoot() ([32]byte, error) {
	buf := make([]byte, constants.BytesPerChunk)
	binary.LittleEndian.PutUint16(buf[:constants.U16Size], uint16(u))
	return [32]byte(buf), nil
}

// IsFixed returns true if the bool is fixed size.
func (U16) IsFixed() bool {
	return true
}

// Type returns the type of the U16.
func (U16) Type() schema.SSZType {
	return schema.U16()
}

// ChunkCount returns the number of chunks required to store the uint16.
func (U16) ChunkCount() uint64 {
	return 1
}
