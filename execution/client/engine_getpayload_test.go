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

//nolint:testpackage // constructs EngineClient with unexported fields to test the GetPayload guard.
package client

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	engineerrors "github.com/berachain/beacon-kit/engine-primitives/errors"
	"github.com/berachain/beacon-kit/execution/client/ethclient"
	"github.com/berachain/beacon-kit/execution/client/ethclient/rpc"
	"github.com/berachain/beacon-kit/log/noop"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/stretchr/testify/require"
)

// TestGetPayloadRejectsNilExecutionPayload ensures that a malicious or buggy
// execution client returning a JSON "null" for "executionPayload" is rejected
// with a clean error at the trust boundary.
func TestGetPayloadRejectsNilExecutionPayload(t *testing.T) {
	t.Parallel()

	logger := noop.NewLogger[any]()
	ec := &EngineClient{
		Client:  ethclient.New(nullPayloadRPCClient{}),
		cfg:     &Config{RPCTimeout: MinRPCTimeout},
		logger:  logger,
		metrics: newClientMetrics(noopTelemetrySink{}, logger),
	}

	envelope, err := ec.GetPayload(
		t.Context(),
		engineprimitives.PayloadID{},
		version.Deneb1(),
	)
	require.ErrorIs(t, err, engineerrors.ErrNilExecutionPayload)
	// the envelope decoded fine, but its inner payload is nil.
	require.NotNil(t, envelope)
	require.Nil(t, envelope.GetExecutionPayload())
}

// TestGetPayloadRejectsNilBlockValue ensures that a malicious or buggy
// execution client returning a JSON "null" for "blockValue" is rejected with a
// clean error at the trust boundary.
func TestGetPayloadRejectsNilBlockValue(t *testing.T) {
	t.Parallel()

	logger := noop.NewLogger[any]()
	ec := &EngineClient{
		Client:  ethclient.New(nullBlockValueRPCClient{}),
		cfg:     &Config{RPCTimeout: MinRPCTimeout},
		logger:  logger,
		metrics: newClientMetrics(noopTelemetrySink{}, logger),
	}

	envelope, err := ec.GetPayload(
		t.Context(),
		engineprimitives.PayloadID{},
		version.Deneb1(),
	)
	require.ErrorIs(t, err, engineerrors.ErrNilBlockValue)
	require.NotNil(t, envelope)
	require.Nil(t, envelope.GetBlockValue())
}

type nullPayloadRPCClient struct{}

var _ rpc.Client = nullPayloadRPCClient{}

func (nullPayloadRPCClient) Initialize() error     { return nil }
func (nullPayloadRPCClient) Start(context.Context) {}
func (nullPayloadRPCClient) Close() error          { return nil }
func (nullPayloadRPCClient) Call(_ context.Context, target any, _ string, _ ...any) error {
	const resp = `{"executionPayload":null,"blockValue":"0x0",` +
		`"blobsBundle":{"commitments":[],"proofs":[],"blobs":[]},` +
		`"shouldOverrideBuilder":false}`
	return json.Unmarshal([]byte(resp), target)
}

type nullBlockValueRPCClient struct{}

var _ rpc.Client = nullBlockValueRPCClient{}

func (nullBlockValueRPCClient) Initialize() error     { return nil }
func (nullBlockValueRPCClient) Start(context.Context) {}
func (nullBlockValueRPCClient) Close() error          { return nil }
func (nullBlockValueRPCClient) Call(_ context.Context, target any, _ string, _ ...any) error {
	// executionPayload is omitted so the envelope keeps the non-nil empty
	// payload from NewEmptyExecutionPayloadEnvelope; only blockValue is null.
	const resp = `{"blockValue":null,` +
		`"blobsBundle":{"commitments":[],"proofs":[],"blobs":[]},` +
		`"shouldOverrideBuilder":false}`
	return json.Unmarshal([]byte(resp), target)
}

type noopTelemetrySink struct{}

func (noopTelemetrySink) IncrementCounter(string, ...string)        {}
func (noopTelemetrySink) MeasureSince(string, time.Time, ...string) {}
