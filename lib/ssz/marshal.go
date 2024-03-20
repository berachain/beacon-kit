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

package ssz

import (
	"encoding/binary"

	"github.com/berachain/beacon-kit/lib/ssz/common"
)

// MarshalComposite marshals a composite object into a byte slice.
func MarshalComposite(c common.Composite) ([]byte, error) {
	elems := c.Elements()
	size := len(elems)

	fixedParts := make([][]byte, size)
	variableParts := make([][]byte, size)

	fixedLen := uint32(0)
	prefixSumVariableLen := make([]uint32, size+1)

	for i, v := range elems {
		bz, err := v.Marshal()
		if err != nil {
			return nil, err
		}
		if !common.IsVariableSize(v.Type()) {
			fixedParts[i] = bz
			fixedLen += uint32(len(bz))
		} else {
			variableParts[i] = bz
			fixedLen += common.BytesPerLengthOffset
			prefixSumVariableLen[i+1] = prefixSumVariableLen[i] + uint32(len(bz))
		}
	}

	for i, p := range fixedParts {
		if len(p) == 0 {
			binary.LittleEndian.PutUint32(fixedParts[i], fixedLen+prefixSumVariableLen[i])
		}
	}

	bz := make([]byte, 0)
	for _, p := range fixedParts {
		bz = append(bz, p...)
	}
	for _, p := range variableParts {
		bz = append(bz, p...)
	}
	return bz, nil
}
