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
	"github.com/holiman/uint256"
)

/* -------------------------------------------------------------------------- */
/*                                Type Definitions                            */
/* -------------------------------------------------------------------------- */

// Ensure type implements schema.SSZObject.
var _ schema.SSZObject[*U256] = (*U256)(nil)

// U256 represents a 256-bit unsigned integer that is both SSZ and JSON.
type U256 uint256.Int

/* -------------------------------------------------------------------------- */
/*                                    U256                                    */
/* -------------------------------------------------------------------------- */

func NewU256FromUint64(v uint64) *U256 {
	return (*U256)(uint256.NewInt(0).SetUint64(v))
}

// SizeSSZ returns the size of the U256 in bytes.
func (U256) SizeSSZ() int {
	return constants.U256Size
}

// MarshalSSZ marshals the U256 into SSZ format.
func (u *U256) MarshalSSZ() ([]byte, error) {
	return (*uint256.Int)(u).MarshalSSZ()
}

// NewFromSSZ creates a new U256 from SSZ format.
func (U256) NewFromSSZ(buf []byte) (*U256, error) {
	if len(buf) != constants.U256Size {
		return nil, fmt.Errorf(
			"invalid buffer length: expected %d, got %d",
			constants.U256Size,
			len(buf),
		)
	}
	u := new(uint256.Int)
	return (*U256)(u), u.UnmarshalSSZ(buf)
}

// HashTreeRoot returns the hash tree root of the U256.
func (u *U256) HashTreeRoot() ([32]byte, error) {
	return (*uint256.Int)(u).HashTreeRoot()
}

// IsFixed returns true if the U256 is fixed size.
func (*U256) IsFixed() bool {
	return true
}

// Type returns the type of the U256.
func (*U256) Type() schema.SSZType {
	return schema.U256()
}

// ChunkCount returns the number of chunks required to store the U256.
func (*U256) ChunkCount() uint64 {
	return 1
}

// Unwrap returns the underlying uint256.Int.
func (u *U256) Unwrap() *uint256.Int {
	return (*uint256.Int)(u)
}
