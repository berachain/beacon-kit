// SPDX-License-Identifier: MIT
//
// Copyright (c) 2025 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package eip7002_test

import (
	"encoding/hex"
	"math"
	"testing"

	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/eip7002"
	beaconmath "github.com/berachain/beacon-kit/primitives/math"
	"github.com/stretchr/testify/require"
)

func TestCreateWithdrawalRequestData(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		pubKey         string
		withdrawAmount uint64 // Assuming math.U64 is an alias for uint64
		expected       string
	}{
		{
			name:           "Normal case",
			pubKey:         "acaf2e8ec309513be835104abc43c8ab27e0665701482d3ce11c592e6ec22910804e8378b0be0f6eb92f452d086599fd",
			withdrawAmount: 10,
			expected:       "0xacaf2e8ec309513be835104abc43c8ab27e0665701482d3ce11c592e6ec22910804e8378b0be0f6eb92f452d086599fd000000000000000a",
		},
		{
			name:           "All zeros pubkey with zero withdrawal",
			pubKey:         "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			withdrawAmount: 0,
			expected:       "0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			name:           "All f's pubkey with max withdrawal",
			pubKey:         "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			withdrawAmount: math.MaxUint64,
			expected:       "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			blsPubKeyBytes, err := hex.DecodeString(tt.pubKey)
			require.NoError(t, err)
			blsPubKey := crypto.BLSPubkey(blsPubKeyBytes)
			// Call the function under test.
			result, err := eip7002.CreateWithdrawalRequestData(blsPubKey, beaconmath.U64(tt.withdrawAmount))
			require.NoError(t, err)
			// Compare the resulting bytes with the expected output.
			require.Equal(t, tt.expected, result.String())
		})
	}
}
