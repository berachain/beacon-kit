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

package schema

import "github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/constants"

const (
	B4Size   = 4
	B8Size   = 8
	B16Size  = 16
	B20Size  = 20
	B32Size  = 32
	B48Size  = 48
	B64Size  = 64
	B96Size  = 96
	B256Size = 256
)

// Basic SSZ types.
// Bool returns an SSZType representing a boolean.
func Bool() SSZType { return basic(constants.BoolSize) }

// U8 returns an SSZType representing an 8-bit unsigned integer.
func U8() SSZType { return basic(constants.U8Size) }

// U16 returns an SSZType representing a 16-bit unsigned integer.
func U16() SSZType { return basic(constants.U16Size) }

// U32 returns an SSZType representing a 32-bit unsigned integer.
func U32() SSZType { return basic(constants.U32Size) }

// U64 returns an SSZType representing a 64-bit unsigned integer.
func U64() SSZType { return basic(constants.U64Size) }

// U128 returns an SSZType representing a 128-bit unsigned integer.
func U128() SSZType { return basic(constants.U128Size) }

// U256 returns an SSZType representing a 256-bit unsigned integer.
func U256() SSZType { return basic(constants.U256Size) }

// B4 creates a DefineByteVector of 4 bytes (32 bits).
func B4() SSZType { return DefineByteVector(B4Size) }

// B8 creates a DefineByteVector of 8 bytes (64 bits).
func B8() SSZType { return DefineByteVector(B8Size) }

// B16 creates a DefineByteVector of 16 bytes (128 bits).
func B16() SSZType { return DefineByteVector(B16Size) }

// B20 creates a DefineByteVector of 20 bytes (160 bits).
func B20() SSZType { return DefineByteVector(B20Size) }

// B32 creates a DefineByteVector of 32 bytes (256 bits).
func B32() SSZType { return DefineByteVector(B32Size) }

// B48 creates a DefineByteVector of 48 bytes (384 bits).
func B48() SSZType { return DefineByteVector(B48Size) }

// B64 creates a DefineByteVector of 64 bytes (512 bits).
func B64() SSZType { return DefineByteVector(B64Size) }

// B96 creates a DefineByteVector of 96 bytes (768 bits).
func B96() SSZType { return DefineByteVector(B96Size) }

// B256 creates a Vector of 256 bytes (2048 bits).
func B256() SSZType { return DefineVector(U8(), B256Size) }
