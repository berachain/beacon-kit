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
	"math/big"

	"github.com/holiman/uint256"
)

// U256 represents a 256-bit unsigned integer that is both SSZ and JSON.
type U256 = uint256.Int

// NewU256 creates a new U256 from a uint64.
func NewU256(v uint64) *U256 {
	return uint256.NewInt(v)
}

// NewU256FromBigInt creates a new U256 from a big.Int.
func NewU256FromBigInt(b *big.Int) (*U256, error) {
	// Negative integers ought to be rejected by math.NewU256FromBigInt(b)
	// since they cannot be expressed in the U256 type. However this does
	// not seem to happen (see holiman/uint256#115), so guarding here.
	if b.Sign() < 0 {
		return nil, fmt.Errorf(
			"cannot convert negative big.Int %s to uint256",
			b.String(),
		)
	}
	return uint256.MustFromBig(b), nil
}

// U256Hex represents a 256-bit unsigned integer that is marshaled to JSON
// as a hexadecimal string.
type U256Hex uint256.Int

// MarshalJSON implements the json.Marshaler interface.
// It returns the hexadecimal string representation of the U256Hex value.
func (u *U256Hex) MarshalJSON() ([]byte, error) {
	return []byte(`"` + (*uint256.Int)(u).Hex() + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It expects the input to be a hexadecimal string.
func (u *U256Hex) UnmarshalJSON(data []byte) error {
	return (*uint256.Int)(u).UnmarshalJSON(data)
}
