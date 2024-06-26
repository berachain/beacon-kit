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

package ssz

import (
	"reflect"

	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// padTo function to pad the chunks to the effective limit with zeroed chunks.
func (m *merkleizer[SpecT, RootT, T]) padTo(
	chunks []RootT,
	size math.U64,
) []RootT {
	switch numChunks := math.U64(len(chunks)); {
	case numChunks == size:
		// No padding needed.
		return chunks
	case numChunks > size:
		// Truncate the chunks to the desired size.
		return chunks[:size]
	default:
		// Append zeroed chunks to the end of the list.
		// #nosec:G701 // size - numChunks is always > 0.
		return append(chunks, m.paddingBuffer.Get(int(size-numChunks))...)
	}
}

// pack packs a list of SSZ-marshallable elements into a single byte slice.
func (m *merkleizer[SpecT, RootT, T]) pack(values []T) ([]RootT, error) {
	// Pack each element into separate buffers.
	var packed []byte
	for _, el := range values {
		fieldValue := reflect.ValueOf(el)
		if fieldValue.Kind() == reflect.Ptr {
			fieldValue = fieldValue.Elem()
		}

		if !fieldValue.CanInterface() {
			return nil, errors.Newf(
				"cannot interface with field %v",
				fieldValue,
			)
		}

		// TODO: Do we need a safety check for Basic only here?
		// TODO: use a real interface instead of hood inline.
		el, ok := reflect.ValueOf(el).
			Interface().(interface{ MarshalSSZ() ([]byte, error) })
		if !ok {
			return nil, errors.Newf("unsupported type %T", el)
		}

		// TODO: Do we need a safety check for Basic only here?
		buf, err := el.MarshalSSZ()
		if err != nil {
			return nil, err
		}
		packed = append(packed, buf...)
	}

	root, _, err := m.partitionBytes(packed)
	return root, err
}

// partitionBytes partitions a byte slice into chunks of a given length.
func (m *merkleizer[SpecT, RootT, T]) partitionBytes(input []byte) (
	[]RootT, uint64, error,
) {
	//nolint:mnd // we add 31 in order to round up the division.
	numChunks := max((len(input)+31)/constants.RootLength, 1)
	chunks := m.intermediateBuffer.Get(numChunks)
	for i := range chunks {
		copy(chunks[i][:], input[32*i:])
	}
	// #nosec:G701 // numChunks is always >= 1.
	return chunks, uint64(numChunks), nil
}
