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

package payloadtime_test

import (
	"testing"
	"time"

	payloadtime "github.com/berachain/beacon-kit/mod/beacon/payload-time"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/stretchr/testify/require"
)

// TestNextTimestampVerifies checks that next payload timestamp
// built via payloadtime.Next always pass payloadtime.Verify.
func TestNextTimestampVerifies(t *testing.T) {
	tests := []struct {
		name        string
		times       func() (time.Time, time.Time)
		expectedErr error
	}{
		{
			name: "Payload timestamp < consensus timestamp",
			times: func() (time.Time, time.Time) {
				consensusTime := time.Now().Truncate(time.Second)
				parentPayloadTimestamp := consensusTime.Add(-10 * time.Second)
				return consensusTime, parentPayloadTimestamp
			},
			expectedErr: nil,
		},
		{
			name: "Payload timestamp == consensus timestamp",
			times: func() (time.Time, time.Time) {
				consensusTime := time.Now().Truncate(time.Second)
				parentPayloadTimestamp := consensusTime
				return consensusTime, parentPayloadTimestamp
			},
			expectedErr: nil,
		},
		{
			name: "Payload timestamp > consensus timestamp",
			times: func() (time.Time, time.Time) {
				consensusTime := time.Now().Truncate(time.Second)
				parentPayloadTimestamp := consensusTime.Add(10 * time.Second)
				return consensusTime, parentPayloadTimestamp
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			consensusTime, parentPayloadTimestamp := tt.times()

			// Optimistic build case
			nextPayload := payloadtime.Next(
				math.U64(consensusTime.Unix()),
				math.U64(parentPayloadTimestamp.Unix()),
				true, // buildOptimistically
			)

			gotErr := payloadtime.Verify(
				math.U64(consensusTime.Unix()),
				math.U64(parentPayloadTimestamp.Unix()),
				nextPayload,
			)
			if tt.expectedErr == nil {
				require.NoError(t, gotErr)
			} else {
				require.ErrorIs(t, tt.expectedErr, gotErr)
			}

			// Just in time build case
			nextPayload = payloadtime.Next(
				math.U64(consensusTime.Unix()),
				math.U64(parentPayloadTimestamp.Unix()),
				false, // buildOptimistically
			)

			gotErr = payloadtime.Verify(
				math.U64(consensusTime.Unix()),
				math.U64(parentPayloadTimestamp.Unix()),
				nextPayload,
			)
			if tt.expectedErr == nil {
				require.NoError(t, gotErr)
			} else {
				require.ErrorIs(t, tt.expectedErr, gotErr)
			}
		})
	}
}
