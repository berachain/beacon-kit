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

package engine

import (
	"context"
	"io"
	"math/big"
	"sync/atomic"
	"testing"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	engineerrors "github.com/berachain/beacon-kit/engine-primitives/errors"
	"github.com/berachain/beacon-kit/execution/client"
	"github.com/berachain/beacon-kit/execution/client/ethclient"
	rpcclient "github.com/berachain/beacon-kit/execution/client/ethclient/rpc"
	"github.com/berachain/beacon-kit/log/phuslu"
	"github.com/berachain/beacon-kit/node-core/components/metrics"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/stretchr/testify/require"
)

func TestNotifyForkchoiceUpdateNoRetryOnSyncing(t *testing.T) {
	t.Parallel()

	cfg := client.DefaultConfig()
	nopSink := metrics.NewNoOpTelemetrySink()
	logger := phuslu.NewLogger(io.Discard, nil)
	engineClient := client.New(
		&cfg,
		logger,
		nil,
		nopSink,
		big.NewInt(80094),
	)
	stubRPC := &syncingStubRPCClient{}
	engineClient.Client = ethclient.New(stubRPC)

	ee := New(engineClient, logger, nopSink)
	req := ctypes.BuildForkchoiceUpdateRequestNoAttrs(
		&engineprimitives.ForkchoiceStateV1{},
		version.Deneb1(),
	)

	_, err := ee.NotifyForkchoiceUpdate(context.Background(), req, false)
	require.ErrorIs(t, err, engineerrors.ErrSyncingPayloadStatus)
	require.EqualValues(t, 1, atomic.LoadInt32(&stubRPC.calls))
}

var _ rpcclient.Client = (*syncingStubRPCClient)(nil)

type syncingStubRPCClient struct {
	calls int32
}

func (*syncingStubRPCClient) Initialize() error { return nil }

func (*syncingStubRPCClient) Start(context.Context) {}

func (s *syncingStubRPCClient) Call(
	_ context.Context,
	target any,
	_ string,
	_ ...any,
) error {
	atomic.AddInt32(&s.calls, 1)

	resp, ok := target.(*engineprimitives.ForkchoiceResponseV1)
	if !ok {
		return nil
	}
	resp.PayloadStatus = engineprimitives.PayloadStatusV1{
		Status: engineprimitives.PayloadStatusSyncing,
	}
	return nil
}

func (*syncingStubRPCClient) Close() error { return nil }
