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

package local

import (
	"context"

	"github.com/itsdevbear/bolaris/builder/types"
	"github.com/itsdevbear/bolaris/types/consensus/v1/capella"
	"google.golang.org/grpc"
)

// Client implements the BuilderServiceClient interface to provide a local simulation of
// builder service operations.
var _ types.BuilderServiceClient = &Client{}

// Client wraps the LocalBlockBuilder to adhere to the BuilderServiceClient interface.
type Client struct {
	local *Builder
}

// NewClient creates a new Client with the given BuilderServiceServer.
func NewClient(local *Builder) *Client {
	return &Client{local: local}
}

// RequestBestBlock simulates a request to the best available block from the builder.
// It directly invokes the RequestBestBlock method of the embedded BuilderServiceServer,
// bypassing gRPC call options.
func (c *Client) RequestBestBlock(
	ctx context.Context, in *types.RequestBestBlockRequest, _ ...grpc.CallOption,
) (*types.RequestBestBlockResponse, error) {
	// Directly call the RequestBestBlock method on the embedded BuilderServiceServer.
	// Note: opts are ignored in this local client simulation.
	beaconBlock, shouldOverride, err := c.local.RequestBestBlock(ctx, in.GetSlot())
	if err != nil {
		return nil, err
	}
	return &types.RequestBestBlockResponse{
		Override: shouldOverride,
		Block:    beaconBlock.(*(capella.BeaconKitBlockCapella)),
	}, nil
}
