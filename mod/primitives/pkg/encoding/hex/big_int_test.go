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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package hex_test

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/hex"
	"github.com/stretchr/testify/require"
)

//nolint:lll // long strings
func TestBigIntRoundTrip(t *testing.T) {
	// assume FromBigInt only called on non-negative big.Int
	tests := []struct {
		name     string
		input    *big.Int
		expected []byte
	}{
		{
			name:     "zero value",
			input:    big.NewInt(0),
			expected: []byte("0x0"),
		},
		{
			name:     "positive value",
			input:    big.NewInt(12345),
			expected: []byte("0x3039"),
		},
		{
			name:     "large positive value",
			input:    new(big.Int).SetBytes(bytes.Repeat([]byte{0xff}, 32)),
			expected: []byte("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hex.MarshalBigIntText(tt.input)
			require.Equal(t, tt.expected, result)

			_, err := hex.ValidateBasicHex(result)
			require.NoError(t, err)

			var dec *big.Int

			if tt.input.Sign() >= 0 {
				dec, err = hex.UnmarshalBigIntText(result)
			} else {
				dec, err = hex.UnmarshalBigIntText(result)
				dec = dec.Neg(dec)
			}

			require.NoError(t, err)
			require.Zero(t, dec.Cmp(tt.input))
		})
	}
}
