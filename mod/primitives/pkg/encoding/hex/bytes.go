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

package hex

import (
	"encoding/hex"
)

func FromBytes[B ~[]byte](b B) string {
	return prefix + hex.EncodeToString(b)
}

// ToBytes returns the bytes represented by the given hex string.
// An error is returned if the input is not a valid hex string.
func ToBytes(hexStr string) ([]byte, error) {
	strippedInput, err := ValidateBasicHex(hexStr)
	if err != nil {
		return nil, err
	}
	return hex.DecodeString(strippedInput)
}

// MustToBytes returns the bytes represented by the given hex string.
// It panics if the input is not a valid hex string.
func MustToBytes(input string) []byte {
	bz, err := ToBytes(input)
	if err != nil {
		panic(err)
	}
	return bz
}
