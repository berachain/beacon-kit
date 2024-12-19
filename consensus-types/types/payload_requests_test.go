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

package types_test

import (
	"testing"

	"github.com/berachain/beacon-kit/consensus-types/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/stretchr/testify/require"
)

func TestBuildNewPayloadRequest(t *testing.T) {
	executionPayload := types.ExecutionPayload{}
	var versionedHashes []common.ExecutionHash
	parentBeaconBlockRoot := common.Root{}
	optimistic := false

	request := types.BuildNewPayloadRequest(
		&executionPayload,
		versionedHashes,
		&parentBeaconBlockRoot,
		optimistic,
	)

	require.NotNil(t, request)
	require.Equal(t, executionPayload, *request.ExecutionPayload)
	require.Equal(t, versionedHashes, request.VersionedHashes)
	require.Equal(t, &parentBeaconBlockRoot, request.ParentBeaconBlockRoot)
	require.Equal(t, optimistic, request.Optimistic)
}

func TestBuildForkchoiceUpdateRequest(t *testing.T) {
	state := &engineprimitives.ForkchoiceStateV1{}
	payloadAttributes := &engineprimitives.PayloadAttributes{}
	forkVersion := uint32(1)

	request := types.BuildForkchoiceUpdateRequest(
		state,
		payloadAttributes,
		forkVersion,
	)

	require.NotNil(t, request)
	require.Equal(t, state, request.State)
	require.Equal(t, payloadAttributes, request.PayloadAttributes)
	require.Equal(t, forkVersion, request.ForkVersion)
}

func TestBuildGetPayloadRequest(t *testing.T) {
	payloadID := engineprimitives.PayloadID{}
	forkVersion := uint32(1)

	request := types.BuildGetPayloadRequest(payloadID, forkVersion)

	require.NotNil(t, request)
	require.Equal(t, payloadID, request.PayloadID)
	require.Equal(t, forkVersion, request.ForkVersion)
}

func TestHasValidVersionedAndBlockHashesPayloadError(t *testing.T) {
	executionPayload := types.ExecutionPayload{}
	versionedHashes := []common.ExecutionHash{}
	parentBeaconBlockRoot := common.Root{}
	optimistic := false

	request := types.BuildNewPayloadRequest(
		&executionPayload,
		versionedHashes,
		&parentBeaconBlockRoot,
		optimistic,
	)

	err := request.HasValidVersionedAndBlockHashes()
	require.ErrorIs(t, err, engineprimitives.ErrPayloadBlockHashMismatch)
}

func TestHasValidVersionedAndBlockHashesMismatchedHashes(t *testing.T) {
	executionPayload := types.ExecutionPayload{}
	versionedHashes := []common.ExecutionHash{
		{},
	}
	parentBeaconBlockRoot := common.Root{}
	optimistic := false

	request := types.BuildNewPayloadRequest(
		&executionPayload,
		versionedHashes,
		&parentBeaconBlockRoot,
		optimistic,
	)

	err := request.HasValidVersionedAndBlockHashes()
	require.ErrorIs(t, err, engineprimitives.ErrMismatchedNumVersionedHashes)
}
