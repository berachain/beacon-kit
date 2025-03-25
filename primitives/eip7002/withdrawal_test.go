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
