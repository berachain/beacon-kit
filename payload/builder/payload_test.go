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

package builder_test

import (
	"context"
	"testing"

	"github.com/berachain/beacon-kit/config/spec"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/log/noop"
	"github.com/berachain/beacon-kit/payload/builder"
	"github.com/berachain/beacon-kit/payload/cache"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/stretchr/testify/require"
)

type mockExecutionPayloadEnvelope[BlobsBundleT engineprimitives.BlobsBundle] struct {
	ExecutionPayload  *ctypes.ExecutionPayload         `json:"executionPayload"`
	BlockValue        *math.U256                       `json:"blockValue"`
	BlobsBundle       BlobsBundleT                     `json:"blobsBundle"`
	ExecutionRequests []ctypes.EncodedExecutionRequest `json:"executionRequests"`
	Override          bool                             `json:"shouldOverrideBuilder"`
}

func (m mockExecutionPayloadEnvelope[BlobsBundleT]) GetExecutionPayload() *ctypes.ExecutionPayload {
	return m.ExecutionPayload
}

func (m mockExecutionPayloadEnvelope[BlobsBundleT]) GetBlockValue() *math.U256 {
	return m.BlockValue
}

func (m mockExecutionPayloadEnvelope[BlobsBundleT]) GetBlobsBundle() engineprimitives.BlobsBundle {
	return m.BlobsBundle
}

func (m mockExecutionPayloadEnvelope[BlobsBundleT]) GetEncodedExecutionRequests() []ctypes.EncodedExecutionRequest {
	return m.ExecutionRequests
}

func (m mockExecutionPayloadEnvelope[BlobsBundleT]) ShouldOverrideBuilder() bool {
	return m.Override
}

// TODO cluster these tests into a single test table
func TestRetrievePayloadSunnyPath(t *testing.T) {
	t.Parallel()

	chainSpec, err := spec.MainnetChainSpec()
	require.NoError(t, err)

	// Create payload builder
	var (
		logger = noop.NewLogger[any]()
		cfg    = &builder.Config{Enabled: true}
		ee     = &stubExecutionEngine{}
		cache  = cache.NewPayloadIDCache()
		af     = &stubAttributesFactory{}
	)
	pb := builder.New(
		cfg,
		chainSpec,
		logger,
		ee,
		cache,
		af,
	)

	// create inputs and set expectations
	var (
		ctx             = context.TODO()
		slot            = math.Slot(2025)
		parentBlockRoot = common.Root{0xff, 0xaa}
		dummyPayloadID  = engineprimitives.PayloadID{0xab}

		expectedPayload = &mockExecutionPayloadEnvelope[*engineprimitives.BlobsBundleV1]{
			ExecutionPayload: &ctypes.ExecutionPayload{
				Withdrawals: engineprimitives.Withdrawals{},
			},
			BlobsBundle: &engineprimitives.BlobsBundleV1{},
		}
	)

	// set expectations
	cache.Set(slot, parentBlockRoot, dummyPayloadID, version.Deneb())
	ee.payloadEnvToReturn = expectedPayload

	// test and checks
	payload, err := pb.RetrievePayload(ctx, slot, parentBlockRoot)
	require.NoError(t, err)
	require.Equal(t, expectedPayload, payload)
}

func TestRetrievePayloadNilWithdrawalsListRejected(t *testing.T) {
	t.Parallel()

	chainSpec, err := spec.MainnetChainSpec()
	require.NoError(t, err)

	// Create payload builder
	var (
		logger = noop.NewLogger[any]()
		cfg    = &builder.Config{Enabled: true}
		ee     = &stubExecutionEngine{}
		cache  = cache.NewPayloadIDCache()
		af     = &stubAttributesFactory{}
	)
	pb := builder.New(
		cfg,
		chainSpec,
		logger,
		ee,
		cache,
		af,
	)

	// create inputs
	var (
		ctx             = context.TODO()
		slot            = math.Slot(2025)
		parentBlockRoot = common.Root{0xff, 0xaa}
		dummyPayloadID  = engineprimitives.PayloadID{0xab}

		faultyPayload = &mockExecutionPayloadEnvelope[*engineprimitives.BlobsBundleV1]{
			ExecutionPayload: &ctypes.ExecutionPayload{
				Withdrawals: nil, // empty withdrawals are fine, nil list should be rejected
			},
			BlobsBundle: &engineprimitives.BlobsBundleV1{},
		}
	)

	// set expectations
	cache.Set(slot, parentBlockRoot, dummyPayloadID, version.Deneb())
	ee.payloadEnvToReturn = faultyPayload

	// test and checks
	_, err = pb.RetrievePayload(ctx, slot, parentBlockRoot)
	require.ErrorIs(t, builder.ErrNilWithdrawals, err)
}

// HELPERS section

var errStubNotImplemented = errors.New("stub not implemented")

type stubExecutionEngine struct {
	payloadEnvToReturn ctypes.BuiltExecutionPayloadEnv
	errToReturn        error
}

func (ee *stubExecutionEngine) GetPayload(
	context.Context, *ctypes.GetPayloadRequest,
) (ctypes.BuiltExecutionPayloadEnv, error) {
	return ee.payloadEnvToReturn, ee.errToReturn
}

func (ee *stubExecutionEngine) NotifyForkchoiceUpdate(
	context.Context, *ctypes.ForkchoiceUpdateRequest,
) (*engineprimitives.PayloadID, error) {
	return nil, errStubNotImplemented
}

type stubAttributesFactory struct{}

func (ee *stubAttributesFactory) BuildPayloadAttributes(
	math.U64, engineprimitives.Withdrawals, common.Bytes32, common.Root, *crypto.BLSPubkey,
) (*engineprimitives.PayloadAttributes, error) {
	return nil, errStubNotImplemented
}
