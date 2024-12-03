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

package engineprimitives_test

import (
	"testing"

	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/primitives/pkg/common"
	"github.com/berachain/beacon-kit/primitives/pkg/encoding/json"
	"github.com/stretchr/testify/require"
)

func TestPayloadID(t *testing.T) {
	payloadID := engineprimitives.PayloadID{
		0x1,
		0x2,
		0x3,
		0x4,
		0x5,
		0x6,
		0x7,
		0x8,
	}

	// Test marshaling
	marshaledID, err := json.Marshal(payloadID)
	require.NoError(t, err)
	require.Equal(t, "\"0x0102030405060708\"", string(marshaledID))

	// Test unmarshaling
	var unmarshaledID engineprimitives.PayloadID
	err = json.Unmarshal(marshaledID, &unmarshaledID)
	require.NoError(t, err)
	require.Equal(t, payloadID, unmarshaledID)
}

func TestForkchoiceStateV1(t *testing.T) {
	state := &engineprimitives.ForkchoiceStateV1{
		HeadBlockHash:      common.ExecutionHash{0x1},
		SafeBlockHash:      common.ExecutionHash{0x2},
		FinalizedBlockHash: common.ExecutionHash{0x3},
	}
	require.Equal(t, common.ExecutionHash{0x1}, state.HeadBlockHash)
	require.Equal(t, common.ExecutionHash{0x2}, state.SafeBlockHash)
	require.Equal(
		t,
		common.ExecutionHash{0x3},
		state.FinalizedBlockHash,
	)

	// Test marshaling
	marshaledState, err := json.Marshal(state)
	require.NoError(t, err)

	// Test unmarshaling
	var unmarshaledState engineprimitives.ForkchoiceStateV1
	err = json.Unmarshal(marshaledState, &unmarshaledState)
	require.NoError(t, err)
	require.Equal(t, state, &unmarshaledState)
}

func TestPayloadStatusV1(t *testing.T) {
	status := &engineprimitives.PayloadStatusV1{
		Status:          engineprimitives.PayloadStatusValid,
		LatestValidHash: &common.ExecutionHash{0x1},
	}
	require.Equal(t, engineprimitives.PayloadStatusValid, status.Status)
	require.Equal(t, &common.ExecutionHash{0x1}, status.LatestValidHash)

	// Test marshaling
	marshaledStatus, err := json.Marshal(status)
	require.NoError(t, err)

	// Test unmarshaling
	var unmarshaledStatus engineprimitives.PayloadStatusV1
	err = json.Unmarshal(marshaledStatus, &unmarshaledStatus)
	require.NoError(t, err)
	require.Equal(t, status, &unmarshaledStatus)
}
