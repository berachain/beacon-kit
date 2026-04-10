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

func TestRetrievePayload(t *testing.T) {
	t.Parallel()

	chainSpec, err := spec.MainnetChainSpec()
	require.NoError(t, err)

	var (
		slot            = math.Slot(10)
		parentBlockRoot = common.Root{0x01}
		denebTimestamp  = math.U64(1_737_381_600) // before mainnet Deneb1ForkTime
		dummyPayloadID  = engineprimitives.PayloadID{0xab}
	)

	validEnvelope := &mockExecutionPayloadEnvelope[*engineprimitives.BlobsBundleV1]{
		ExecutionPayload: &ctypes.ExecutionPayload{
			Timestamp:   denebTimestamp,
			Withdrawals: engineprimitives.Withdrawals{},
		},
		BlobsBundle: &engineprimitives.BlobsBundleV1{},
	}
	nilWithdrawalsEnvelope := &mockExecutionPayloadEnvelope[*engineprimitives.BlobsBundleV1]{
		ExecutionPayload: &ctypes.ExecutionPayload{
			Timestamp:   denebTimestamp,
			Withdrawals: nil,
		},
		BlobsBundle: &engineprimitives.BlobsBundleV1{},
	}
	nilPayloadEnvelope := &mockExecutionPayloadEnvelope[*engineprimitives.BlobsBundleV1]{
		ExecutionPayload: nil,
	}

	tests := []struct {
		name string

		// If non-nil, seed the PayloadIDCache before calling RetrievePayload.
		cachePayloadID   *engineprimitives.PayloadID
		cacheForkVersion common.Version

		// Stub response from the execution engine's GetPayload.
		eeEnvelope ctypes.BuiltExecutionPayloadEnv

		// If non-nil, cache as the latest verified payload for the same slot.
		verifiedEnvelope ctypes.BuiltExecutionPayloadEnv

		expectedForkVersion common.Version
		wantEnvelope        ctypes.BuiltExecutionPayloadEnv
		wantErr             error
	}{
		{
			name:                "sunny path via PayloadIDCache",
			cachePayloadID:      &dummyPayloadID,
			cacheForkVersion:    version.Deneb(),
			eeEnvelope:          validEnvelope,
			expectedForkVersion: version.Deneb(),
			wantEnvelope:        validEnvelope,
		},
		{
			name:                "nil withdrawals list rejected",
			cachePayloadID:      &dummyPayloadID,
			cacheForkVersion:    version.Deneb(),
			eeEnvelope:          nilWithdrawalsEnvelope,
			expectedForkVersion: version.Deneb(),
			wantErr:             builder.ErrNilWithdrawals,
		},
		{
			name:                "fallback reuses verified payload on fork version match",
			verifiedEnvelope:    validEnvelope,
			expectedForkVersion: version.Deneb(),
			wantEnvelope:        validEnvelope,
		},
		{
			name:                "fallback rejects verified payload on fork version mismatch",
			verifiedEnvelope:    validEnvelope,
			expectedForkVersion: version.Electra(),
			wantErr:             builder.ErrPayloadIDNotFound,
		},
		{
			name:                "fallback skips verified payload with nil execution payload",
			verifiedEnvelope:    nilPayloadEnvelope,
			expectedForkVersion: version.Deneb(),
			wantErr:             builder.ErrPayloadIDNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ee := &stubExecutionEngine{payloadEnvToReturn: tt.eeEnvelope}
			pc := cache.NewPayloadIDCache()
			pb := builder.New(
				&builder.Config{Enabled: true},
				chainSpec,
				noop.NewLogger[any](),
				ee,
				pc,
				&stubAttributesFactory{},
			)

			if tt.cachePayloadID != nil {
				pc.Set(slot, parentBlockRoot, *tt.cachePayloadID, tt.cacheForkVersion)
			}
			if tt.verifiedEnvelope != nil {
				pb.CacheLatestVerifiedPayload(slot, tt.verifiedEnvelope)
			}

			//nolint:govet // shadow err so that parallel tests do not overwrite err.
			envelope, err := pb.RetrievePayload(t.Context(), slot, parentBlockRoot, tt.expectedForkVersion)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantEnvelope, envelope)
		})
	}
}

// HELPERS section

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
	_ context.Context, _ *ctypes.ForkchoiceUpdateRequest, _ bool,
) (*engineprimitives.PayloadID, error) {
	return nil, errStubNotImplemented
}

type stubAttributesFactory struct{}

func (ee *stubAttributesFactory) BuildPayloadAttributes(
	math.U64, engineprimitives.Withdrawals, common.Bytes32, common.Root, *crypto.BLSPubkey,
) (*engineprimitives.PayloadAttributes, error) {
	return nil, errStubNotImplemented
}
