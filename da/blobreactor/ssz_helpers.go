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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package blobreactor

// EncodeBlobSidecarsSSZ takes multiple SSZ-encoded BlobSidecar bytes and combines them
// into a single SSZ-encoded BlobSidecars (slice) format.
// The encoding is: 4-byte offset + concatenated sidecars (each 131928 bytes).
//
//nolint:mnd // ok for now
func EncodeBlobSidecarsSSZ(sidecarBzs [][]byte) []byte {
	if len(sidecarBzs) == 0 {
		// Empty list needs 4-byte offset pointing to itself
		return []byte{4, 0, 0, 0}
	}

	// BlobSidecars is encoded as: offset (4 bytes) + data
	// The offset points to where the data starts (after the offset itself)
	offset := uint32(4) // Data starts after the 4-byte offset

	// Calculate total size
	totalSize := 4 // offset
	for _, data := range sidecarBzs {
		totalSize += len(data)
	}

	result := make([]byte, totalSize)

	// Write offset in little-endian
	result[0] = byte(offset)
	result[1] = byte(offset >> 8)
	result[2] = byte(offset >> 16)
	result[3] = byte(offset >> 24)

	// Concatenate all sidecars after the offset
	pos := 4
	for _, data := range sidecarBzs {
		copy(result[pos:], data)
		pos += len(data)
	}

	return result
}
