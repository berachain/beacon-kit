// SPDX-License-Identifier: MIT
//
// Copyright (c) 2025 Berachain Foundation
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
// WdeHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package eip7685

import (
	"fmt"

	"github.com/karalabe/ssz"
)

// MarshalItems marshals a slice of items that satisfy SSZMarshaler.
// It encodes each item individually and appends its bytes to the output buffer.
func MarshalItems[T sszMarshaler](items []T) ([]byte, error) {
	var buf []byte
	for i, item := range items {
		itemBytes, err := item.MarshalSSZ()
		if err != nil {
			return nil, fmt.Errorf("failed to marshal item at index %d: %w", i, err)
		}
		buf = append(buf, itemBytes...)
	}
	return buf, nil
}

// UnmarshalItems decodes a slice of items from the provided data.
// It assumes that each item is encoded into a fixed number of bytes (itemSize)
// and that newItem returns a new instance of the item.
func UnmarshalItems[T sszUnmarshaler](data []byte, itemSize int, newItem func() T) ([]T, error) {
	if len(data)%itemSize != 0 {
		return nil, fmt.Errorf("invalid data length: %d is not a multiple of item size %d", len(data), itemSize)
	}
	numItems := len(data) / itemSize
	items := make([]T, 0, numItems)
	for i := 0; i < len(data); i += itemSize {
		chunk := data[i : i+itemSize]
		item := newItem()
		err := ssz.DecodeFromBytes(chunk, item)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal item at index %d: %w", i/itemSize, err)
		}
		items = append(items, item)
	}
	return items, nil
}
