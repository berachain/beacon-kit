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

package constants

const (
	// BytesPerChunk is the number of bytes per chunk.
	BytesPerChunk = 32

	// BytesPerLengthOffset is the number of bytes per serialized length offset.
	BytesPerLengthOffset = 4

	// BitsPerByte is the number of bits per byte.
	BitsPerByte = 8

	// ByteSize is the size of a single byte.
	ByteSize = 1

	// BoolSize is the size of a boolean in bytes.
	BoolSize = 1

	// U8Size is the size of an 8-bit unsigned integer in bytes.
	U8Size = 1

	// U16Size is the size of a 16-bit unsigned integer in bytes.
	U16Size = 2

	// U32Size is the size of a 32-bit unsigned integer in bytes.
	U32Size = 4

	// U64Size is the size of a 64-bit unsigned integer in bytes.
	U64Size = 8

	// U128Size is the size of a 128-bit unsigned integer in bytes.
	U128Size = 16

	// U256Size is the size of a 256-bit unsigned integer in bytes.
	U256Size = 32
)
