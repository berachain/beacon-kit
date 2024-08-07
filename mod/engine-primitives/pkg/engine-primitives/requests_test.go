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

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives/mocks"
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/stretchr/testify/require"
)

type MockExecutionPayload struct {
}
type MockWithdrawal struct{}

func (m MockExecutionPayload) Empty(uint32) MockExecutionPayload {
	return m
}
func (m MockExecutionPayload) IsNil() bool {
	return false
}
func (m MockExecutionPayload) Version() uint32 {
	return 0
}
func (m MockExecutionPayload) GetPrevRandao() common.Bytes32 {
	return common.Bytes32{}
}
func (m MockExecutionPayload) GetBlockHash() common.ExecutionHash {
	return common.ExecutionHash{}
}
func (m MockExecutionPayload) GetParentHash() common.ExecutionHash {
	return common.ExecutionHash{}
}
func (m MockExecutionPayload) GetNumber() math.U64 {
	return math.U64(0)
}
func (m MockExecutionPayload) GetGasLimit() math.U64 {
	return math.U64(0)
}
func (m MockExecutionPayload) GetGasUsed() math.U64 {
	return math.U64(0)
}
func (m MockExecutionPayload) GetTimestamp() math.U64 {
	return math.U64(0)
}
func (m MockExecutionPayload) GetExtraData() []byte {
	return []byte{}
}
func (m MockExecutionPayload) GetBaseFeePerGas() *math.U256 {
	return &math.U256{}
}

func (
	m MockExecutionPayload,
) GetFeeRecipient() gethprimitives.ExecutionAddress {
	return gethprimitives.ExecutionAddress{}
}
func (m MockExecutionPayload) GetStateRoot() common.Bytes32 {
	return common.Bytes32{}
}
func (m MockExecutionPayload) GetReceiptsRoot() common.Bytes32 {
	return common.Bytes32{}
}
func (m MockExecutionPayload) GetLogsBloom() bytes.B256 {
	return [256]byte{}
}
func (m MockExecutionPayload) GetBlobGasUsed() math.U64 {
	return math.U64(0)
}
func (m MockExecutionPayload) GetExcessBlobGas() math.U64 {
	return math.U64(0)
}
func (m MockExecutionPayload) GetWithdrawals() []MockWithdrawal {
	return []MockWithdrawal{}
}
func (m MockExecutionPayload) GetTransactions() engineprimitives.Transactions {
	return [][]byte{}
}

func (m MockWithdrawal) GetIndex() math.U64 {
	return math.U64(0)
}
func (m MockWithdrawal) GetAmount() math.U64 {
	return math.U64(0)
}
func (m MockWithdrawal) GetAddress() gethprimitives.ExecutionAddress {
	return gethprimitives.ExecutionAddress{}
}
func (m MockWithdrawal) GetValidatorIndex() math.U64 {
	return math.U64(0)
}

func TestBuildNewPayloadRequest(t *testing.T) {
	executionPayload := MockExecutionPayload{}
	var versionedHashes []common.ExecutionHash
	parentBeaconBlockRoot := common.Root{}
	optimistic := false

	request := engineprimitives.BuildNewPayloadRequest(
		executionPayload,
		versionedHashes,
		&parentBeaconBlockRoot,
		optimistic,
	)

	require.NotNil(t, request)
	require.Equal(t, executionPayload, request.ExecutionPayload)
	require.Equal(t, versionedHashes, request.VersionedHashes)
	require.Equal(t, &parentBeaconBlockRoot, request.ParentBeaconBlockRoot)
	require.Equal(t, optimistic, request.Optimistic)
}

func TestBuildForkchoiceUpdateRequest(t *testing.T) {
	state := &engineprimitives.ForkchoiceStateV1{}
	payloadAttributes := &mocks.PayloadAttributer{}
	forkVersion := uint32(1)

	request := engineprimitives.BuildForkchoiceUpdateRequest(
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

	request := engineprimitives.BuildGetPayloadRequest(payloadID, forkVersion)

	require.NotNil(t, request)
	require.Equal(t, payloadID, request.PayloadID)
	require.Equal(t, forkVersion, request.ForkVersion)
}

func TestHasValidVersionedAndBlockHashesPayloadError(t *testing.T) {
	executionPayload := MockExecutionPayload{}
	versionedHashes := []common.ExecutionHash{}
	parentBeaconBlockRoot := common.Root{}
	optimistic := false

	request := engineprimitives.BuildNewPayloadRequest(
		executionPayload,
		versionedHashes,
		&parentBeaconBlockRoot,
		optimistic,
	)

	err := request.HasValidVersionedAndBlockHashes()
	require.ErrorIs(t, err, engineprimitives.ErrPayloadBlockHashMismatch)
}

func TestHasValidVersionedAndBlockHashesMismatchedHashes(t *testing.T) {
	executionPayload := MockExecutionPayload{}
	versionedHashes := []common.ExecutionHash{
		common.ExecutionHash{},
	}
	parentBeaconBlockRoot := common.Root{}
	optimistic := false

	request := engineprimitives.BuildNewPayloadRequest(
		executionPayload,
		versionedHashes,
		&parentBeaconBlockRoot,
		optimistic,
	)

	err := request.HasValidVersionedAndBlockHashes()
	require.ErrorIs(t, err, engineprimitives.ErrMismatchedNumVersionedHashes)
}
