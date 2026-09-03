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

package blobreactor_test

import (
	"testing"

	"github.com/berachain/beacon-kit/da/blobreactor"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/stretchr/testify/require"
)

func TestMessages_ByRootRequestRoundTrip(t *testing.T) {
	t.Parallel()
	msg := &blobreactor.SidecarsByRootRequest{
		RequestID: 0xdeadbeef12345678,
		Slot:      42,
		BlockRoot: common.Root{0x01, 0x02},
	}
	bz, err := msg.MarshalSSZ()
	require.NoError(t, err)

	var decoded blobreactor.SidecarsByRootRequest
	require.NoError(t, decoded.UnmarshalSSZ(bz))
	require.Equal(t, *msg, decoded)
}

func TestMessages_ByRangeRequestRoundTrip(t *testing.T) {
	t.Parallel()
	msg := &blobreactor.SidecarsByRangeRequest{
		RequestID: 7,
		StartSlot: 100,
		Count:     32,
	}
	bz, err := msg.MarshalSSZ()
	require.NoError(t, err)

	var decoded blobreactor.SidecarsByRangeRequest
	require.NoError(t, decoded.UnmarshalSSZ(bz))
	require.Equal(t, *msg, decoded)
}

func TestMessages_ByRangeRequestCountValidation(t *testing.T) {
	t.Parallel()
	for _, count := range []uint64{0, blobreactor.MaxRequestedSlots + 1} {
		msg := &blobreactor.SidecarsByRangeRequest{RequestID: 1, StartSlot: 1, Count: count}
		bz, err := msg.MarshalSSZ()
		require.NoError(t, err)

		var decoded blobreactor.SidecarsByRangeRequest
		require.Error(t, decoded.UnmarshalSSZ(bz), "count %d must be rejected", count)
	}
}

func TestMessages_PushRoundTrip(t *testing.T) {
	t.Parallel()
	msg := &blobreactor.SidecarsPush{
		BlockRoot:   common.Root{0xaa},
		SidecarData: []byte{1, 2, 3, 4},
	}
	bz, err := msg.MarshalSSZ()
	require.NoError(t, err)

	var decoded blobreactor.SidecarsPush
	require.NoError(t, decoded.UnmarshalSSZ(bz))
	require.Equal(t, msg.BlockRoot, decoded.BlockRoot)
	require.Equal(t, msg.SidecarData, decoded.SidecarData)
}

func TestMessages_PushRejectsEmptyData(t *testing.T) {
	t.Parallel()
	msg := &blobreactor.SidecarsPush{BlockRoot: common.Root{0xaa}}
	bz, err := msg.MarshalSSZ()
	require.NoError(t, err)

	var decoded blobreactor.SidecarsPush
	require.Error(t, decoded.UnmarshalSSZ(bz))
}

func TestMessages_HaveRoundTrip(t *testing.T) {
	t.Parallel()
	msg := &blobreactor.SidecarsHave{BlockRoot: common.Root{0x11, 0x22}}
	bz, err := msg.MarshalSSZ()
	require.NoError(t, err)

	var decoded blobreactor.SidecarsHave
	require.NoError(t, decoded.UnmarshalSSZ(bz))
	require.Equal(t, *msg, decoded)
}

func TestMessages_ResponseRoundTrip(t *testing.T) {
	t.Parallel()
	msg := &blobreactor.SidecarsResponse{
		RequestID:     99,
		SidecarChunks: [][]byte{{1, 2, 3}, {4, 5}},
	}
	bz, err := msg.MarshalSSZ()
	require.NoError(t, err)

	var decoded blobreactor.SidecarsResponse
	require.NoError(t, decoded.UnmarshalSSZ(bz))
	require.Equal(t, msg.RequestID, decoded.RequestID)
	require.Equal(t, msg.SidecarChunks, decoded.SidecarChunks)
}

func TestMessages_ResponseEmptyChunks(t *testing.T) {
	t.Parallel()
	msg := &blobreactor.SidecarsResponse{RequestID: 1}
	bz, err := msg.MarshalSSZ()
	require.NoError(t, err)

	var decoded blobreactor.SidecarsResponse
	require.NoError(t, decoded.UnmarshalSSZ(bz))
	require.Empty(t, decoded.SidecarChunks)
}

// Garbage bytes must not decode.
func TestMessages_GarbageRejected(t *testing.T) {
	t.Parallel()
	garbage := []byte{0xff, 0xfe, 0xfd}
	require.Error(t, new(blobreactor.SidecarsByRootRequest).UnmarshalSSZ(garbage))
	require.Error(t, new(blobreactor.SidecarsByRangeRequest).UnmarshalSSZ(garbage))
	require.Error(t, new(blobreactor.SidecarsHave).UnmarshalSSZ(garbage))
	require.Error(t, new(blobreactor.SidecarsPush).UnmarshalSSZ(garbage))
	require.Error(t, new(blobreactor.SidecarsResponse).UnmarshalSSZ(garbage))
}
