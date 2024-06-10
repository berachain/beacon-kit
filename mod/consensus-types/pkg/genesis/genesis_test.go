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

package genesis_test

import (
	"testing"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/genesis"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	"github.com/stretchr/testify/require"
)

func TestDefaultGenesisDeneb(t *testing.T) {
	g := genesis.DefaultGenesisDeneb()
	if g.ForkVersion != version.FromUint32[primitives.Version](version.Deneb) {
		t.Errorf(
			"Expected fork version %v, but got %v",
			version.FromUint32[primitives.Version](
				version.Deneb,
			),
			g.ForkVersion,
		)
	}

	if len(g.Deposits) != 0 {
		t.Errorf("Expected no deposits, but got %v", len(g.Deposits))
	}
	// add assertions for ExecutionPayloadHeader
	payloadHeader := g.ExecutionPayloadHeader
	if payloadHeader == nil {
		t.Errorf("Expected ExecutionPayloadHeader to be non-nil")
	}

	require.Equal(t, common.ZeroHash, payloadHeader.ParentHash,
		"Unexpected ParentHash")
	require.Equal(t, common.ZeroAddress, payloadHeader.FeeRecipient,
		"Unexpected FeeRecipient")
	require.Equal(t, math.U64(30000000), payloadHeader.GasLimit,
		"Unexpected GasLimit")
	require.Equal(t, math.U64(0), payloadHeader.GasUsed,
		"Unexpected GasUsed")
	require.Equal(t, math.U64(0), payloadHeader.Timestamp,
		"Unexpected Timestamp")
}

func TestDefaultGenesisExecutionPayloadHeaderDeneb(t *testing.T) {
	header, err := genesis.DefaultGenesisExecutionPayloadHeaderDeneb()
	require.NoError(t, err)
	require.NotNil(t, header)
}
