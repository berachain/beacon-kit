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
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/berachain/beacon-kit/testing/utils"
	"github.com/stretchr/testify/require"
)

func TestBuildNewPayloadRequestFromFork(t *testing.T) {
	t.Parallel()

	runForAllSupportedVersions(t, func(t *testing.T, v common.Version) {
		block := utils.GenerateValidBeaconBlock(t, v)

		var parentProposerPubKey *crypto.BLSPubkey
		if version.EqualsOrIsAfter(v, version.Electra1()) {
			parentProposerPubKey = &crypto.BLSPubkey{0x01}
		}
		request, err := types.BuildNewPayloadRequestFromFork(block, parentProposerPubKey)
		require.NoError(t, err)
		require.NotNil(t, request)
		require.Equal(t, block.GetBody().GetExecutionPayload(), request.GetExecutionPayload())
		require.Equal(t, block.GetBody().GetBlobKzgCommitments().ToVersionedHashes(), request.GetVersionedHashes())
		require.Equal(t, block.GetParentBlockRoot(), request.GetParentBeaconBlockRoot())

		if version.EqualsOrIsAfter(v, version.Electra()) {
			requests, getErr := block.GetBody().GetExecutionRequests()
			require.NoError(t, getErr)
			list, getErr := types.GetExecutionRequestsList(requests)
			require.NoError(t, getErr)
			executionRequests, getErr := request.GetEncodedExecutionRequests()
			require.NoError(t, getErr)
			require.Equal(t, list, executionRequests)
		}
	})
}

func TestBuildForkchoiceUpdateRequest(t *testing.T) {
	t.Parallel()
	var (
		state       = &engineprimitives.ForkchoiceStateV1{}
		forkVersion = version.Deneb1()
	)
	payloadAttributes, err := engineprimitives.NewPayloadAttributes(
		forkVersion,
		math.U64(time.Now().Truncate(time.Second).Unix()),
		common.Bytes32{0x01},
		common.ExecutionAddress{},
		engineprimitives.Withdrawals{},
		common.Root{},
		nil,
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
	runForAllSupportedVersions(t, func(t *testing.T, v common.Version) {
		block := utils.GenerateValidBeaconBlock(t, v)
		// Remove txs and kzg commitments from body cos not valid
		block.GetBody().SetExecutionPayload(&types.ExecutionPayload{
			Versionable: types.NewVersionable(v),
		})
		block.GetBody().SetBlobKzgCommitments(eip4844.KZGCommitments[common.ExecutionHash]{})

		var parentProposerPubKey *crypto.BLSPubkey
		if version.EqualsOrIsAfter(v, version.Electra1()) {
			parentProposerPubKey = &crypto.BLSPubkey{0x01}
		}
		request, err := types.BuildNewPayloadRequestFromFork(block, parentProposerPubKey)
		require.NoError(t, err)
		require.NotNil(t, request)
		require.Equal(t, block.GetBody().GetExecutionPayload(), request.GetExecutionPayload())
		require.Equal(t, block.GetBody().GetBlobKzgCommitments().ToVersionedHashes(), request.GetVersionedHashes())
		require.Equal(t, block.GetParentBlockRoot(), request.GetParentBeaconBlockRoot())

		if version.EqualsOrIsAfter(v, version.Electra()) {
			requests, getErr := block.GetBody().GetExecutionRequests()
			require.NoError(t, getErr)
			list, getErr := types.GetExecutionRequestsList(requests)
			require.NoError(t, getErr)
			executionRequests, getErr := request.GetEncodedExecutionRequests()
			require.NoError(t, getErr)
			require.Equal(t, list, executionRequests)
		}
		err = request.HasValidVersionedAndBlockHashes()
		require.ErrorIs(t, err, engineprimitives.ErrPayloadBlockHashMismatch)
	})
}

func TestHasValidVersionedAndBlockHashesMismatchedHashes(t *testing.T) {
	t.Parallel()

	runForAllSupportedVersions(t, func(t *testing.T, v common.Version) {
		block := utils.GenerateValidBeaconBlock(t, v)
		// Remove txs and kzg commitments from body cos not valid
		block.GetBody().SetExecutionPayload(&types.ExecutionPayload{
			Versionable: types.NewVersionable(v),
		})
		block.GetBody().SetBlobKzgCommitments(eip4844.KZGCommitments[common.ExecutionHash]{{}})

		var parentProposerPubKey *crypto.BLSPubkey
		if version.EqualsOrIsAfter(v, version.Electra1()) {
			parentProposerPubKey = &crypto.BLSPubkey{0x01}
		}
		request, err := types.BuildNewPayloadRequestFromFork(block, parentProposerPubKey)
		require.NoError(t, err)

		err = request.HasValidVersionedAndBlockHashes()
		require.ErrorIs(t, err, engineprimitives.ErrMismatchedNumVersionedHashes)
	})
}
