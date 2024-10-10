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

import "math/big"

// FromBigInt encodes bigint as a hex string with 0x prefix.
// Precondition: bigint is non-negative.
func FromBigInt(bigint *big.Int) string {
	switch sign := bigint.Sign(); {
	case sign == 0:
		return Prefix + "0"
	case sign > 0:
		return Prefix + bigint.Text(hexBase)
	default:
		// this return should never reach if precondition is met
		return Prefix + bigint.Text(hexBase)[1:]
	}
}

// ToBigInt decodes a hex string with 0x prefix.
func ToBigInt(hexStr string) (*big.Int, error) {
	raw, err := formatAndValidateNumber(hexStr)
	if err != nil {
		return nil, err
	}
	if len(raw) > nibblesPer256Bits {
		return nil, ErrBig256Range
	}
	bigWordNibbles, err := getBigWordNibbles()
	if err != nil {
		return nil, err
	}
	words := make([]big.Word, len(raw)/bigWordNibbles+1)
	end := len(raw)
	for i := range words {
		start := end - bigWordNibbles
		if start < 0 {
			start = 0
		}
		for ri := start; ri < end; ri++ {
			nib := decodeNibble(raw[ri])
			if nib == badNibble {
				return nil, ErrInvalidString
			}
			words[i] *= 16
			words[i] += big.Word(nib)
		}
		end = start
	}
	dec := new(big.Int).SetBits(words)
	return dec, nil
}

// MustToBigInt decodes a hex string with 0x prefix.
// It panics for invalid input.
func MustToBigInt(hexStr string) *big.Int {
	bi, err := ToBigInt(hexStr)
	if err != nil {
		panic(err)
	}
	return bi
}
