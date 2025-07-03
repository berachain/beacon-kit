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

package sszutil

import (
	"fmt"

	"github.com/berachain/beacon-kit/primitives/constraints"
)

// FastSSZUnmarshaler is a minimal interface for types that can unmarshal from fastssz.
type FastSSZUnmarshaler interface {
	UnmarshalSSZ([]byte) error
	ValidateAfterDecodingSSZ() error
}

// Unmarshal is the way we build objects from byte formatted in SSZ encoding.
// This function highlights the common template for SSZ decoding different objects.
func Unmarshal[T FastSSZUnmarshaler](buf []byte, v T) error {
	if err := v.UnmarshalSSZ(buf); err != nil {
		return fmt.Errorf("failed decoding %T: %w", v, err)
	}

	// Note: ValidateAfterDecodingSSZ may change v even if it returns error
	// (depending on the specific implementations)
	return v.ValidateAfterDecodingSSZ()
}

// MarshalItemsEIP7685 marshals a slice of items that satisfy SSZMarshaler according
// to the EIP-7685 standard. It encodes each item individually and appends its bytes
// to the output buffer.
func MarshalItemsEIP7685[T constraints.SSZMarshaler](items []T) ([]byte, error) {
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

// UnmarshalItemsEIP7685 decodes a slice of items from the provided data according
// to the EIP-7685 standard. It assumes that each item is encoded into a fixed number
// of bytes (itemSize) and that newItem returns a new instance of the item.
func UnmarshalItemsEIP7685[T FastSSZUnmarshaler](
	data []byte,
	itemSize int,
	newItem func() T,
) ([]T, error) {
	if len(data)%itemSize != 0 {
		return nil, fmt.Errorf(
			"invalid data length: %d is not a multiple of item size %d",
			len(data), itemSize,
		)
	}
	numItems := len(data) / itemSize
	items := make([]T, 0, numItems)
	for i := 0; i < len(data); i += itemSize {
		chunk := data[i : i+itemSize]
		item := newItem()
		if err := Unmarshal(chunk, item); err != nil {
			return nil, fmt.Errorf(
				"failed to unmarshal item at index %d: %w", i/itemSize, err,
			)
		}
		items = append(items, item)
	}
	return items, nil
}
