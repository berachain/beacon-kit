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

package math

import (
	"fmt"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/schema"
)

/* -------------------------------------------------------------------------- */
/*                                Type Definitions                            */
/* -------------------------------------------------------------------------- */

// Ensure type implements schema.SSZObject.
var _ schema.SSZObject[U8] = (*U8)(nil)

// U8 represents a 8-bit unsigned integer that is both SSZ and JSON.
type U8 uint8

/* -------------------------------------------------------------------------- */
/*                                     U8                                     */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the size of the uint8 in bytes.
func (U8) SizeSSZ() int {
	return constants.U8Size
}

// MarshalSSZ marshals the uint8 into SSZ format.
func (u U8) MarshalSSZ() ([]byte, error) {
	return []byte{byte(u)}, nil
}

// NewFromSSZ creates a new U8 from SSZ format.
func (U8) NewFromSSZ(buf []byte) (U8, error) {
	if len(buf) != constants.U8Size {
		return 0, fmt.Errorf(
			"invalid buffer length: expected %d, got %d",
			constants.U8Size,
			len(buf),
		)
	}

	//#nosec:G701 // the check above protects against overflow.
	return U8(buf[0]), nil
}

// HashTreeRoot returns the hash tree root of the uint8.
func (u U8) HashTreeRoot() ([32]byte, error) {
	buf := make([]byte, constants.BytesPerChunk)
	buf[0] = byte(u)
	return [32]byte(buf), nil
}

// IsFixed returns true if the bool is fixed size.
func (U8) IsFixed() bool {
	return true
}

// Type returns the type of the U8.
func (U8) Type() schema.SSZType {
	return schema.U8()
}

// ChunkCount returns the number of chunks required to store the uint8.
func (U8) ChunkCount() uint64 {
	return 1
}
