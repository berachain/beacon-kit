// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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
