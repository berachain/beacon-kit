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
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package bytes

import (
	"encoding/hex"
)

// B96 represents a 96-byte array.
type B96 [96]byte

// ToBytes96 is a utility function that transforms a byte slice into a fixed
// 96-byte array. If the input exceeds 96 bytes, it gets truncated.
func ToBytes96(input []byte) B96 {
	return B96(ExtendToSize(input, B96Size))
}

// MarshalText implements the encoding.TextMarshaler interface for Bytes96.
func (h B96) MarshalText() ([]byte, error) {
	return []byte(h.String()), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface for Bytes96.
func (h *B96) UnmarshalJSON(input []byte) error {
	return unmarshalJSONHelper(h[:], input)
}

// String returns the hex string representation of Bytes96.
func (h B96) String() string {
	return "0x" + hex.EncodeToString(h[:])
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for Bytes96.
func (h *B96) UnmarshalText(text []byte) error {
	return UnmarshalTextHelper(h[:], text)
}
