// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
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

package engineprimitives_test

import (
	"encoding/json"
	"testing"

	engineprimitives "github.com/berachain/beacon-kit/mod/primitives-engine"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
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
	require.Equal(t, common.ExecutionHash{0x3}, state.FinalizedBlockHash)

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
