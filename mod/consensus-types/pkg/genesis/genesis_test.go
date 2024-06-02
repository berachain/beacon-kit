// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

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
