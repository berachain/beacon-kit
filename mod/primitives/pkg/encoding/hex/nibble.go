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

package hex

import (
	"math/big"
)

// decodeNibble decodes a single hexadecimal nibble (half-byte) into uint64.
func decodeNibble(in byte) uint64 {
	// uint64 conversion here is safe
	switch {
	case in >= '0' && in <= '9':
		//#nosec G701 // The resulting value will be in the range 0-9.
		return uint64(in - hexBaseOffset)
	case in >= 'A' && in <= 'F':
		//#nosec G701 // The resulting value will be in the range 10-15.
		return uint64(in - hexAlphaOffsetUpper)
	case in >= 'a' && in <= 'f':
		//#nosec G701 // The resulting value will be in the range 10-15.
		return uint64(in - hexAlphaOffsetLower)
	default:
		return badNibble
	}
}

// getBigWordNibbles returns the number of nibbles required for big.Word.
//
//nolint:mnd // this is fine xD
func getBigWordNibbles() (int, error) {
	// This is a weird way to compute the number of nibbles required for
	// big.Word. The usual way would be to use constant arithmetic but go vet
	// can't handle that
	b, _ := new(big.Int).SetString("FFFFFFFFFF", 16)
	switch len(b.Bits()) {
	case 1:
		return 16, nil
	case 2:
		return 8, nil
	default:
		return 0, ErrInvalidBigWordSize
	}
}
