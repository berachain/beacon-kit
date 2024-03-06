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

package enginetypes

import (
	"encoding/json"
	"reflect"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// PayloadID is a custom type representing Payload IDs in a fixed-size byte
// array.
type PayloadID [8]byte

// MarshalJSON converts a PayloadID to a JSON-encoded hexadecimal string.
func (b PayloadID) MarshalJSON() ([]byte, error) {
	return json.Marshal(hexutil.Bytes(b[:]))
}

// UnmarshalJSON improves the conversion of a JSON-encoded hexadecimal string
// to a PayloadID by directly utilizing UnmarshalFixedJSON.
func (b *PayloadID) UnmarshalJSON(enc []byte) error {
	if err := hexutil.UnmarshalFixedJSON(
		reflect.TypeOf(PayloadID{}), enc, b[:],
	); err != nil {
		return err
	}
	return nil
}
