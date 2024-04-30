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

package bytes

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// ------------------------------ Helpers ------------------------------

// Helper function to unmarshal JSON for various byte types.
func unmarshalJSONHelper(target []byte, input []byte) error {
	bz := hexutil.Bytes{}
	if err := bz.UnmarshalJSON(input); err != nil {
		return err
	}
	if len(bz) != len(target) {
		return fmt.Errorf(
			"incorrect length, expected %d bytes but got %d",
			len(target), len(bz),
		)
	}
	copy(target, bz)
	return nil
}

// Helper function to unmarshal text for various byte types.
func unmarshalTextHelper(target []byte, text []byte) error {
	bz := hexutil.Bytes{}
	if err := bz.UnmarshalText(text); err != nil {
		return err
	}
	if len(bz) != len(target) {
		return fmt.Errorf(
			"incorrect length, expected %d bytes but got %d",
			len(target), len(bz),
		)
	}
	copy(target, bz)
	return nil
}
