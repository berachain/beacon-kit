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

package types_test

import (
	"testing"
	"time"

	"github.com/berachain/beacon-kit/consensus-types/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/stretchr/testify/require"
)

func TestBuildNewPayloadRequest(t *testing.T) {
	t.Parallel()
	var (
		executionPayload = &types.ExecutionPayload{
			Versionable: types.NewVersionable(version.Deneb1()),
		}
		versionedHashes       []common.ExecutionHash
		parentBeaconBlockRoot = common.Root{}
	)

	request := types.BuildNewPayloadRequest(
		executionPayload,
		versionedHashes,
		&parentBeaconBlockRoot,
	)

	require.NotNil(t, request)
	require.Equal(t, executionPayload, request.ExecutionPayload)
	require.Equal(t, versionedHashes, request.VersionedHashes)
	require.Equal(t, &parentBeaconBlockRoot, request.ParentBeaconBlockRoot)
}

func TestBuildForkchoiceUpdateRequest(t *testing.T) {
	t.Parallel()
	var (
		state       = &engineprimitives.ForkchoiceStateV1{}
		forkVersion = version.Deneb1()
	)
	payloadAttributes, err := engineprimitives.NewPayloadAttributes(
		forkVersion,
		uint64(time.Now().Truncate(time.Second).Unix()),
		common.Bytes32{0x01},
		common.ExecutionAddress{},
		engineprimitives.Withdrawals{},
		common.Root{},
	)
	require.NoError(t, err)

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
	t.Parallel()
	payloadID := engineprimitives.PayloadID{}
	forkVersion := version.Altair()

	request := types.BuildGetPayloadRequest(payloadID, forkVersion)

	require.NotNil(t, request)
	require.Equal(t, payloadID, request.PayloadID)
	require.Equal(t, forkVersion, request.ForkVersion)
}

func TestHasValidVersionedAndBlockHashesPayloadError(t *testing.T) {
	t.Parallel()
	var (
		executionPayload = &types.ExecutionPayload{
			Versionable: types.NewVersionable(version.Deneb1()),
		}
		versionedHashes       = []common.ExecutionHash{}
		parentBeaconBlockRoot = common.Root{}
	)

	request := types.BuildNewPayloadRequest(
		executionPayload,
		versionedHashes,
		&parentBeaconBlockRoot,
	)

	err := request.HasValidVersionedAndBlockHashes()
	require.ErrorIs(t, err, engineprimitives.ErrPayloadBlockHashMismatch)
}

func TestHasValidVersionedAndBlockHashesMismatchedHashes(t *testing.T) {
	t.Parallel()
	var (
		executionPayload = &types.ExecutionPayload{
			Versionable: types.NewVersionable(version.Deneb1()),
		}
		versionedHashes       = []common.ExecutionHash{{}}
		parentBeaconBlockRoot = common.Root{}
	)

	request := types.BuildNewPayloadRequest(
		executionPayload,
		versionedHashes,
		&parentBeaconBlockRoot,
	)

	err := request.HasValidVersionedAndBlockHashes()
	require.ErrorIs(t, err, engineprimitives.ErrMismatchedNumVersionedHashes)
}
