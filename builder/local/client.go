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

	"github.com/ethereum/go-ethereum/common"
	"github.com/itsdevbear/bolaris/builder/types"
	"google.golang.org/grpc"
)

// This Cleint adheres to the BuilderServiceClient interface.
var _ types.BuilderServiceClient = &Client{}

// Client wraps the LocalBlockBuilder to adhere to the BuilderServiceClient interface.
type Client struct {
	local *Builder
}

// NewClient creates a new Client with the given BuilderServiceServer.
func NewClient(local *Builder) *Client {
	return &Client{local: local}
}

// RequestBestBlock is used to request the best available block from the builder.
// It directly calls the RequestBestBlock method on the local builder, instead of
// making a gRPC call.
//
// NOTE: The gRPC opts are ignored for obvious reasons.
func (c *Client) GetExecutionPayload(
	ctx context.Context, in *types.GetExecutionPayloadRequest, _ ...grpc.CallOption,
) (*types.GetExecutionPayloadResponse, error) {
	// Directly call the RequestBestBlock method on the local builder.
	executionPayload, shouldOverride, err := c.local.GetExecutionPayload(
		ctx, common.BytesToHash(in.GetParentHash()), in.GetSlot(),
	)
	if err != nil {
		return nil, err
	}

	// Return the response.
	return &types.GetExecutionPayloadResponse{
		Override:         shouldOverride,
		PayloadContainer: executionPayload,
	}, nil
}
