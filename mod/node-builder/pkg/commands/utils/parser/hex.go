// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

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
