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
	"encoding/binary"
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/json"

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives/mocks"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/stretchr/testify/require"
)

type MockExecutionPayloadT struct {
	Value string
}

func (m MockExecutionPayloadT) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Value)
}

func (m MockExecutionPayloadT) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &m.Value)
}

func TestExecutionPayloadEnvelope(t *testing.T) {
	blobsBundle := &mocks.BlobsBundle{}

	// Convert the int to a byte slice
	valueBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(valueBytes, uint64(100))

	// Use math.NewU256L
	blockValue, err := math.NewU256L(valueBytes)
	if err != nil {
		t.Fatalf("Failed to convert int to U256L: %v", err)
	}

	envelope := &engineprimitives.ExecutionPayloadEnvelope[
		MockExecutionPayloadT,
		*mocks.BlobsBundle,
	]{
		ExecutionPayload: MockExecutionPayloadT{Value: "test"},
		BlockValue:       blockValue,
		BlobsBundle:      blobsBundle,
		Override:         true,
	}

	payload := envelope.GetExecutionPayload()
	require.Equal(t, envelope.ExecutionPayload, payload)

	value := envelope.GetValue()
	require.Equal(t, envelope.BlockValue, value)

	bundle := envelope.GetBlobsBundle()
	require.Equal(t, envelope.BlobsBundle, bundle)

	override := envelope.ShouldOverrideBuilder()
	require.Equal(t, envelope.Override, override)
}
