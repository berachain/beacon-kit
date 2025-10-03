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

package blobreactor_test

import (
	"testing"

	"github.com/berachain/beacon-kit/da/blobreactor"
	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/primitives/encoding/ssz"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/stretchr/testify/require"
)

func TestBlobRequest_SSZ(t *testing.T) {
	t.Parallel()
	original := &blobreactor.BlobRequest{Slot: 12345, RequestID: 67890}

	// Marshal
	data, err := original.MarshalSSZ()
	require.NoError(t, err)
	require.Len(t, data, 16, "BlobRequest should be 16 bytes (8 for Slot + 8 for RequestID)")

	// Unmarshal
	decoded := &blobreactor.BlobRequest{}
	err = decoded.UnmarshalSSZ(data)
	require.NoError(t, err)
	require.Equal(t, original.Slot, decoded.Slot)
}

func TestBlobResponse_SSZ(t *testing.T) {
	t.Parallel()
	t.Run("without sidecars", func(t *testing.T) {
		original := &blobreactor.BlobResponse{
			Slot:        42,
			SidecarData: nil,
			HeadSlot:    100,
		}

		// Marshal
		data, err := original.MarshalSSZ()
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(data), 20, "BlobResponse should be at least 20 bytes")

		// Unmarshal
		decoded := &blobreactor.BlobResponse{}
		err = decoded.UnmarshalSSZ(data)
		require.NoError(t, err)
		require.Equal(t, original.Slot, decoded.Slot)
		require.Equal(t, original.HeadSlot, decoded.HeadSlot)
		require.Empty(t, decoded.SidecarData)
	})

	t.Run("with sidecars", func(t *testing.T) {
		// Create a test sidecar and marshal it to SSZ
		sidecars := datypes.BlobSidecars{&datypes.BlobSidecar{Index: 0}}
		sidecarData, err := sidecars.MarshalSSZ()
		require.NoError(t, err)

		original := &blobreactor.BlobResponse{
			Slot:        100,
			SidecarData: sidecarData,
			HeadSlot:    200,
		}

		// Marshal
		data, err := original.MarshalSSZ()
		require.NoError(t, err)
		require.Greater(t, len(data), 20, "BlobResponse with sidecars should be > 20 bytes")

		// Unmarshal
		decoded := &blobreactor.BlobResponse{}
		err = decoded.UnmarshalSSZ(data)
		require.NoError(t, err)
		require.Equal(t, original.Slot, decoded.Slot)
		require.Equal(t, original.HeadSlot, decoded.HeadSlot)
		require.NotEmpty(t, decoded.SidecarData)

		// Verify we can unmarshal the sidecar data back
		var decodedSidecars datypes.BlobSidecars
		err = ssz.Unmarshal(decoded.SidecarData, &decodedSidecars)
		require.NoError(t, err)
		require.Len(t, decodedSidecars, 1)
	})
}

func TestSSZ_InvalidData(t *testing.T) {
	t.Parallel()
	// Test that decoding fails gracefully with invalid data
	t.Run("BlobRequest too short", func(t *testing.T) {
		req := &blobreactor.BlobRequest{}
		err := req.UnmarshalSSZ([]byte{1, 2, 3}) // Only 3 bytes, need 8
		require.Error(t, err)
	})

	t.Run("BlobResponse too short", func(t *testing.T) {
		resp := &blobreactor.BlobResponse{}
		err := resp.UnmarshalSSZ([]byte{1, 2, 3, 4, 5}) // Only 5 bytes, need 20
		require.Error(t, err)
	})
}

func TestMessageWithTypePrefix(t *testing.T) {
	t.Parallel()
	// Test the complete message flow with type prefix
	req := &blobreactor.BlobRequest{Slot: 1000}

	// Marshal SSZ
	sszData, err := req.MarshalSSZ()
	require.NoError(t, err)

	// Add message type prefix (simulating what the reactor does)
	fullMsg := append([]byte{byte(blobreactor.MessageTypeRequest)}, sszData...)

	// Create BlobMessage wrapper
	blobMsg := blobreactor.NewBlobMessage(fullMsg)
	require.NotNil(t, blobMsg)

	// Extract and decode (simulating Receive method)
	msgType := blobreactor.MessageType(blobMsg.Data[0])
	msgData := blobMsg.Data[1:]

	require.Equal(t, blobreactor.MessageTypeRequest, msgType)

	decoded := &blobreactor.BlobRequest{}
	err = decoded.UnmarshalSSZ(msgData)
	require.NoError(t, err)
	require.Equal(t, math.Slot(1000), decoded.Slot)
}
