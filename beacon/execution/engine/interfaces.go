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

package engine

import (
	"context"

	gethcoretypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/itsdevbear/bolaris/third_party/go-ethereum/common"
	enginev1 "github.com/itsdevbear/bolaris/third_party/prysm/proto/engine/v1"
	"github.com/itsdevbear/bolaris/types/consensus/interfaces"
	"github.com/itsdevbear/bolaris/types/consensus/primitives"
	payloadattribute "github.com/prysmaticlabs/prysm/v4/consensus-types/payload-attribute"
)

// Caller defines a client that can interact with an Ethereum
// execution node's engine engineClient via JSON-RPC.
type Caller interface {
	// Generic Methods
	//
	// ConnectedETH1 returns true if the client is connected to the execution node.
	ConnectedETH1() bool

	// Engine API Related Methods
	//
	// NewPayload creates a new payload for the Ethereum execution node.
	NewPayload(ctx context.Context, payload interfaces.ExecutionData,
		versionedHashes []common.Hash, parentBlockRoot *common.Hash) ([]byte, error)
	// ForkchoiceUpdated updates the fork choice of the Ethereum execution node.
	ForkchoiceUpdated(
		ctx context.Context, state *enginev1.ForkchoiceState, attrs payloadattribute.Attributer,
	) (*enginev1.PayloadIDBytes, []byte, error)
	// GetPayload retrieves the payload from the Ethereum execution node.
	GetPayload(ctx context.Context, payloadID primitives.PayloadID,
		slot primitives.Slot) (interfaces.ExecutionData, *enginev1.BlobsBundle, bool, error)
	// ExecutionBlockByHash retrieves the execution block by its hash.
	ExecutionBlockByHash(ctx context.Context, hash common.Hash,
		withTxs bool) (*enginev1.ExecutionBlock, error)

	// Eth Namespace Methods
	//
	// BlockByHash retrieves the block by its hash.
	HeaderByHash(ctx context.Context, hash common.Hash) (*gethcoretypes.Header, error)
}
