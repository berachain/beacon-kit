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

package merkleizer

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
)

// pack packs a list of SSZ-marshallable elements into a single byte slice.
func pack[
	RootT ~[32]byte,
	T interface {
		MarshalSSZ() ([]byte, error)
	},
](
	values []T,
) ([]RootT, uint64, error) {
	// pack(values): Given ordered objects of the same basic type:
	// Serialize values into bytes.
	// If not aligned to a multiple of BYTES_PER_CHUNK bytes,
	// right-pad with zeroes to the next multiple.
	// Partition the bytes into BYTES_PER_CHUNK-byte chunks.
	// Return the chunks.
	var packed []byte
	for _, el := range values {
		buf, err := el.MarshalSSZ()
		if err != nil {
			return nil, 0, err
		}
		packed = append(packed, buf...)
	}

	chunks, numChunks := chunkifyBytes[RootT](packed)
	return chunks, numChunks, nil
}

// chunkifyBytes partitions a byte slice into chunks of a given length.
func chunkifyBytes[RootT ~[32]byte](input []byte) (
	[]RootT, uint64,
) {
	//nolint:mnd // we add 31 in order to round up the division.
	numChunks := max((len(input)+31)/constants.RootLength, 1)
	// TODO: figure out how to safely chunk these bytes.
	chunks := make([]RootT, numChunks)
	for i := range chunks {
		copy(chunks[i][:], input[32*i:])
	}
	//#nosec:G701 // numChunks is always >= 1.
	return chunks, uint64(numChunks)
}

// packBits packs a list of SSZ-marshallable bitlists into a single byte slice.
//
//nolint:unused // todo eventually implement this function.
func packBits[
	RootT ~[32]byte,
	T interface {
		MarshalSSZ() ([]byte, error)
	},
]([]T) ([]RootT, error) {
	// pack_bits(bits): Given the bits of bitlist or bitvector, get
	// bitfield_bytes by packing them in bytes and aligning to the start.
	// The length-delimiting bit for bitlists is excluded. Then return pack
	// (bitfield_bytes).
	panic("not yet implemented")
}
