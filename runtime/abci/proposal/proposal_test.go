// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
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

package proposal_test

import (
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/itsdevbear/bolaris/config"
	"github.com/itsdevbear/bolaris/runtime/abci/proposal"
	"github.com/stretchr/testify/assert"
)

func TestRemoveBeaconBlockFromTxs(t *testing.T) {
	tests := []struct {
		name            string
		inputTxs        [][]byte
		payloadPosition uint
		expectedTxs     [][]byte
	}{
		{
			name:            "remove first element",
			inputTxs:        [][]byte{{1}, {2}, {3}},
			payloadPosition: 0,
			expectedTxs:     [][]byte{{2}, {3}},
		},
		{
			name:            "remove last element",
			inputTxs:        [][]byte{{1}, {2}, {3}},
			payloadPosition: 2,
			expectedTxs:     [][]byte{{1}, {2}},
		},
		{
			name:            "remove middle element",
			inputTxs:        [][]byte{{1}, {2}, {3}},
			payloadPosition: 1,
			expectedTxs:     [][]byte{{1}, {3}},
		},
		{
			name:            "remove single element",
			inputTxs:        [][]byte{{1}},
			payloadPosition: 0,
			expectedTxs:     [][]byte{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := proposal.NewHandler(
				&config.Proposal{BeaconKitBlockPosition: tt.payloadPosition}, nil, nil, nil, nil,
			)
			req := &abci.RequestProcessProposal{Txs: tt.inputTxs}
			result := handler.RemoveBeaconBlockFromTxs(req)
			assert.Equal(t, tt.expectedTxs, result.Txs)
		})
	}
}
