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

package types_test

import (
	"fmt"
	"testing"

	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/encoding/ssz"
	prysmtypes "github.com/prysmaticlabs/prysm/v5/proto/prysm/v1alpha1"
	"github.com/stretchr/testify/require"
)

func TestPendingPartialWithdrawal_ValidValuesSSZ(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name    string
		pending *types.PendingPartialWithdrawal
	}{
		{
			name: "basic",
			pending: &types.PendingPartialWithdrawal{
				ValidatorIndex:    1,
				Amount:            1000,
				WithdrawableEpoch: 10,
			},
		},
		{
			name: "zero amount",
			pending: &types.PendingPartialWithdrawal{
				ValidatorIndex:    2,
				Amount:            0,
				WithdrawableEpoch: 20,
			},
		},
		{
			name: "max values",
			pending: &types.PendingPartialWithdrawal{
				ValidatorIndex:    1<<64 - 1,
				Amount:            1<<64 - 1,
				WithdrawableEpoch: 1<<64 - 1,
			},
		},
		{
			name: "random-ish values",
			pending: &types.PendingPartialWithdrawal{
				ValidatorIndex:    7,
				Amount:            54321,
				WithdrawableEpoch: 999,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Marshal the original pending partial withdrawal.
			pendingBytes, err := tc.pending.MarshalSSZ()
			require.NoError(t, err)

			// Unmarshal into the prysm type.
			var prysmType prysmtypes.PendingPartialWithdrawal
			err = prysmType.UnmarshalSSZ(pendingBytes)
			require.NoError(t, err)

			// Compare the HashTreeRoots.
			originalHTR := tc.pending.HashTreeRoot()
			prysmHTR, err := prysmType.HashTreeRoot()
			require.NoError(t, err)
			require.Equal(t, originalHTR[:], prysmHTR[:])

			// Marshal the prysm request
			prysmBytes, err := prysmType.MarshalSSZ()
			require.NoError(t, err)

			// Unmarshal back into original type
			var recomputedPending types.PendingPartialWithdrawal
			err = ssz.Unmarshal(prysmBytes, &recomputedPending)
			require.NoError(t, err)
			require.Equal(t, *tc.pending, recomputedPending)
		})
	}
}

//nolint:paralleltest // Invalid SSZ values cannot be run in parallel due to shared zeroalloc.
func TestPendingPartialWithdrawal_InvalidValuesUnmarshalSSZ(t *testing.T) {
	// Build a valid pending partial withdrawal to get a baseline payload.
	validPending := &types.PendingPartialWithdrawal{
		ValidatorIndex:    1,
		Amount:            1000,
		WithdrawableEpoch: 10,
	}
	validBytes, err := validPending.MarshalSSZ()
	require.NoError(t, err)

	// Define several invalid payloads.
	invalidPayloads := [][]byte{
		nil,                       // nil slice
		{},                        // empty slice
		[]byte("this is not ssz"), // arbitrary non-SSZ data
		{0x00, 0x01},              // too short to be valid
		// A truncated valid payload: remove last 5 bytes.
		func() []byte {
			if len(validBytes) > 5 {
				return validBytes[:len(validBytes)-5]
			}
			return validBytes
		}(),
		// A valid payload with extra trailing bytes.
		func() []byte {
			extra := []byte{0xAA, 0xBB, 0xCC, 0xDD}
			return append(validBytes, extra...)
		}(),
	}

	for i, payload := range invalidPayloads {
		i, payload := i, payload // capture loop variables
		t.Run(fmt.Sprintf("invalidPendingPartialWithdrawal_%d", i), func(t *testing.T) {
			// Ensure that Unmarshal does not panic and returns an error.
			require.NotPanics(t, func() {
				var p types.PendingPartialWithdrawal
				err = ssz.Unmarshal(payload, &p)
				require.Error(t, err, "expected error for payload %v", payload)
			})
		})
	}
}
