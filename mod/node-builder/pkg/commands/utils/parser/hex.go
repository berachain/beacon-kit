// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
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

package parser

import (
	"encoding/hex"
	"strings"
)

// DecodeFrom0xPrefixedString decodes a 0x prefixed hex string.
// Note: Use of this function would force the input to contain a 0x prefix
// since otherwise it would cause ambiguity in the conversion.
func DecodeFrom0xPrefixedString(data string) ([]byte, error) {
	if !strings.HasPrefix(data, "0x") {
		return nil, ErrInvalid0xPrefixedHexString
	}
	return hex.DecodeString(data[2:])
}

// EncodeTo0xPrefixedString encodes a byte slice to a 0x prefixed hex string.
func EncodeTo0xPrefixedString(data []byte) string {
	return "0x" + hex.EncodeToString(data)
}
