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

package ethclient_test

import (
	"context"
	"testing"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/execution/client/ethclient"
	"github.com/berachain/beacon-kit/execution/client/ethclient/rpc"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/stretchr/testify/require"
)

// TestGetPayloadV3NeverReturnsEmptyPayload shows that execution payload
// returned by ethClient is not nil.
func TestGetPayloadV3NeverReturnsEmptyPayload(t *testing.T) {
	c := ethclient.New(&stubRPCClient{t: t})

	var (
		ctx         = context.Background()
		payloadID   engineprimitives.PayloadID
		forkVersion = version.Deneb1()
	)

	pe, err := c.GetPayloadV3(ctx, payloadID, forkVersion)
	require.NoError(t, err)

	// check that execution payload is not nil
	require.False(t, pe.GetExecutionPayload().IsNil())
}

// TestNewPayloadWithValidVersion tests that NewPayload correctly handles Deneb version.
func TestNewPayloadWithValidVersion(t *testing.T) {
	c := ethclient.New(&stubRPCClient{t: t})
	ctx := context.Background()

	payload := &ctypes.ExecutionPayload{EpVersion: version.Deneb1()}
	versionedHashes := []common.ExecutionHash{}
	var parentBlockRoot *common.Root

	_, err := c.NewPayload(ctx, payload, versionedHashes, parentBlockRoot)
	require.NoError(t, err)
}

// TestNewPayloadWithInvalidVersion tests that NewPayload returns ErrInvalidVersion for Capella.
func TestNewPayloadWithInvalidVersion(t *testing.T) {
	c := ethclient.New(&stubRPCClient{t: t})
	ctx := context.Background()

	payload := &ctypes.ExecutionPayload{EpVersion: version.Capella()}
	versionedHashes := []common.ExecutionHash{}
	var parentBlockRoot *common.Root

	_, err := c.NewPayload(ctx, payload, versionedHashes, parentBlockRoot)
	require.ErrorIs(t, err, ethclient.ErrInvalidVersion)
}

// TestForkchoiceUpdatedWithValidVersion tests that ForkchoiceUpdated correctly handles Deneb version.
func TestForkchoiceUpdatedWithValidVersion(t *testing.T) {
	c := ethclient.New(&stubRPCClient{t: t})
	ctx := context.Background()

	state := &engineprimitives.ForkchoiceStateV1{}
	attrs := struct{}{}
	forkVersion := version.Deneb1()

	_, err := c.ForkchoiceUpdated(ctx, state, attrs, forkVersion)
	require.NoError(t, err)
}

// TestForkchoiceUpdatedWithInvalidVersion tests that ForkchoiceUpdated returns ErrInvalidVersion for Capella.
func TestForkchoiceUpdatedWithInvalidVersion(t *testing.T) {
	c := ethclient.New(&stubRPCClient{t: t})
	ctx := context.Background()

	state := &engineprimitives.ForkchoiceStateV1{}
	attrs := struct{}{}
	forkVersion := version.Capella()

	_, err := c.ForkchoiceUpdated(ctx, state, attrs, forkVersion)
	require.ErrorIs(t, err, ethclient.ErrInvalidVersion)
}

// TestGetPayloadWithValidVersion tests that GetPayload correctly handles >= Deneb version.
func TestGetPayloadWithValidVersion(t *testing.T) {
	c := ethclient.New(&stubRPCClient{t: t})
	ctx := context.Background()

	var payloadID engineprimitives.PayloadID
	forkVersion := version.Deneb1()

	_, err := c.GetPayload(ctx, payloadID, forkVersion)
	require.NoError(t, err)
}

// TestGetPayloadWithInvalidVersion tests that GetPayload returns ErrInvalidVersion for Capella.
func TestGetPayloadWithInvalidVersion(t *testing.T) {
	c := ethclient.New(&stubRPCClient{t: t})
	ctx := context.Background()

	var payloadID engineprimitives.PayloadID
	forkVersion := version.Capella()

	_, err := c.GetPayload(ctx, payloadID, forkVersion)
	require.ErrorIs(t, err, ethclient.ErrInvalidVersion)
}

var _ rpc.Client = (*stubRPCClient)(nil)

type stubRPCClient struct {
	t *testing.T
}

func (tc *stubRPCClient) Start(context.Context) {}
func (tc *stubRPCClient) Call(_ context.Context, target any, _ string, _ ...any) error {
	tc.t.Helper()
	require.NotNil(tc.t, target)

	// If calling ForkchoiceUpdated, set the PayloadStatus to not empty to
	// avoid returning ErrNilResponse.
	if fcu, ok := target.(*engineprimitives.ForkchoiceResponseV1); ok {
		fcu.PayloadStatus = engineprimitives.PayloadStatusV1{
			Status: "not empty",
		}
	}

	return nil
}
func (tc *stubRPCClient) Close() error { return nil }
