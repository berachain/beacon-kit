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

	"github.com/sourcegraph/conc/iter"
)

// MarshalComposite marshals a composite object into a byte slice.
// The function is implemented according to the SSZ specification
// at
// https://github.com/ethereum/consensus-specs/blob/dev/ssz/simple-serialize.md#vectors-containers-lists
func MarshalComposite(c SSZTypeGeneric) ([]byte, error) {
	s := NewSerializer()
	elems := s.Elements(c)
	bzs, err := iter.MapErr(
		elems,
		func(elem *SSZTypeGeneric) ([]byte, error) {
			return s.MarshalSSZ(*elem)
		},
	)
	if err != nil {
		return nil, err
	}

	size := len(elems)

	fixedParts := make([][]byte, size)
	variableParts := make([][]byte, size)

	fixedLen := uint32(0)
	prefixSumVariableLen := make([]uint32, size+1)

	for i, v := range elems {
		// Carry over the offset of the previous variable part.
		prefixSumVariableLen[i+1] = prefixSumVariableLen[i]
		// If the type is fixed-size, add the data to the fixed part.
		if IsFixedSize(v) {
			fixedParts[i] = bzs[i]
			fixedLen += uint32(len(bzs[i]))
		} else {
			variableParts[i] = bzs[i]
			// Place a placeholder for the offset of the variable part.
			fixedLen += BytesPerLengthOffset
			prefixSumVariableLen[i+1] += uint32(len(bzs[i]))
		}
	}

	for i, p := range fixedParts {
		if len(p) == 0 {
			// The current position is a variable part,
			// add the offset to the start of the actual data in the heap.
			fixedParts[i] = make([]byte, BytesPerLengthOffset)
			binary.LittleEndian.PutUint32(
				fixedParts[i],
				fixedLen+prefixSumVariableLen[i],
			)
		}
	}

	totalLength := fixedLen + prefixSumVariableLen[size]

	output := make([]byte, totalLength)
	offset := 0
	for _, p := range fixedParts {
		copy(output[offset:], p)
		offset += len(p)
	}

	for _, p := range variableParts {
		copy(output[offset:], p)
		offset += len(p)
	}

	return output, nil
}
