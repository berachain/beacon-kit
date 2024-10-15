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
	"errors"
	"fmt"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/hex"
)

const B20Size = 20

// B20 represents a 20-byte fixed-size byte array.
type B20 [B20Size]byte

var ErrIncorrectLength = errors.New("incorrect length")

// ToBytes20 converts a byte slice into a fixed 20-byte array.
func ToBytes20(input []byte) (B20, error) {
	if len(input) != B20Size {
		return B20{}, fmt.Errorf("%w: got %d, expected %d", ErrIncorrectLength, len(input), B20Size)
	}
	var b20 B20
	copy(b20[:], input)
	return b20, nil
}

// MarshalText implements encoding.TextMarshaler for B20.
func (h B20) MarshalText() ([]byte, error) {
	return []byte(h.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler for B20.
func (h *B20) UnmarshalText(text []byte) error {
	return UnmarshalTextHelper(h[:], text)
}

// String returns the hex string representation of B20.
func (h B20) String() string {
	return hex.EncodeBytes(h[:])
}

// UnmarshalJSON implements json.Unmarshaler for B20.
func (h *B20) UnmarshalJSON(input []byte) error {
	return UnmarshalJSONHelper(h[:], input)
}

// MarshalSSZ implements SSZ marshaling for B20.
func (h B20) MarshalSSZ() ([]byte, error) {
	return h[:], nil
}

// HashTreeRoot returns the hash tree root of the B20.
func (h B20) HashTreeRoot() (B32, error) {
	return ToBytes32(ExtendToSize(h[:], B32Size))
}
