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

package ethclient

import (
	"context"

	ethengine "github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	enginev1 "github.com/itsdevbear/bolaris/engine/types/v1"
)

// Eth1Client is a struct that holds the Ethereum 1 client and
// its configuration.
type Eth1Client struct {
	*ethclient.Client
}

// NewEth1Client creates a new Ethereum 1 client with the provided
// context and options.
func NewEth1Client(client *ethclient.Client) (*Eth1Client, error) {
	c := &Eth1Client{
		Client: client,
	}
	return c, nil
}

// NewPayloadV3 calls the engine_newPayloadV3 method via JSON-RPC.
func (s *Eth1Client) NewPayloadV3(
	ctx context.Context, payload *enginev1.ExecutionPayloadDeneb,
	versionedHashes []common.Hash, parentBlockRoot *common.Hash,
) (*enginev1.PayloadStatus, error) {
	result := &enginev1.PayloadStatus{}
	if err := s.Client.Client().CallContext(
		ctx, result, NewPayloadMethodV3, payload, versionedHashes, parentBlockRoot,
	); err != nil {
		return nil, err
	}
	return result, nil
}

// ForkchoiceUpdatedV3 calls the engine_forkchoiceUpdatedV3 method via JSON-RPC.
func (s *Eth1Client) ForkchoiceUpdatedV3(
	ctx context.Context,
	state *enginev1.ForkchoiceState,
	attrs *enginev1.PayloadAttributesV3,
) (*ForkchoiceUpdatedResponse, error) {
	return s.forkchoiceUpdateCall(ctx, ForkchoiceUpdatedMethodV3, state, attrs)
}

// forkchoiceUpdateCall is a helper function to call to any version
// of the forkchoiceUpdates method.
func (s *Eth1Client) forkchoiceUpdateCall(
	ctx context.Context,
	method string,
	state *enginev1.ForkchoiceState,
	attrs any,
) (*ForkchoiceUpdatedResponse, error) {
	result := &ForkchoiceUpdatedResponse{}

	if err := s.Client.Client().CallContext(
		ctx, result, method, state, attrs,
	); err != nil {
		return nil, err
	}

	if result.Status == nil {
		return nil, ErrNilResponse
	}

	return result, nil
}

// GetPayloadV3 calls the engine_getPayloadV3 method via JSON-RPC.
func (s *Eth1Client) GetPayloadV3(
	ctx context.Context, payloadID enginev1.PayloadIDBytes,
) (*enginev1.ExecutionPayloadContainer, error) {
	result := &enginev1.ExecutionPayloadDenebWithValueAndBlobsBundle{}
	if err := s.Client.Client().CallContext(
		ctx, result, GetPayloadMethodV3, payloadID,
	); err != nil {
		return nil, err
	}
	return &enginev1.ExecutionPayloadContainer{
		Payload: &enginev1.ExecutionPayloadContainer_Deneb{
			Deneb: result.GetPayload(),
		},
		PayloadValue:          result.GetValue(),
		BlobsBundle:           result.GetBlobsBundle(),
		ShouldOverrideBuilder: result.GetShouldOverrideBuilder(),
	}, nil
}

// ExecutionBlockByHash fetches an execution engine block by hash by calling
// eth_blockByHash via JSON-RPC.
func (s *Eth1Client) ExecutionBlockByHash(
	ctx context.Context, hash common.Hash, withTxs bool,
) (*enginev1.ExecutionBlock, error) {
	result := &enginev1.ExecutionBlock{}
	err := s.Client.Client().CallContext(
		ctx, result, BlockByHashMethod, hash, withTxs)
	return result, err
}

// ExecutionBlockByNumber fetches an execution engine block by number
// by calling eth_getBlockByNumber via JSON-RPC.
func (s *Eth1Client) ExecutionBlockByNumber(
	ctx context.Context, num rpc.BlockNumber, withTxs bool,
) (*enginev1.ExecutionBlock, error) {
	result := &enginev1.ExecutionBlock{}
	err := s.Client.Client().CallContext(
		ctx, result, BlockByNumberMethod, num, withTxs)
	return result, err
}

// GetClientVersionV1 calls the engine_getClientVersionV1 method via JSON-RPC.
func (s *Eth1Client) GetClientVersionV1(
	ctx context.Context,
) ([]ethengine.ClientVersionV1, error) {
	result := make([]ethengine.ClientVersionV1, 0)
	if err := s.Client.Client().CallContext(
		ctx, &result, GetClientVersionV1, nil,
	); err != nil {
		return nil, err
	}
	return result, nil
}

// ExchangeCapabilities calls the engine_exchangeCapabilities method via
// JSON-RPC.
func (s *Eth1Client) ExchangeCapabilities(
	ctx context.Context,
	capabilities []string,
) ([]string, error) {
	result := make([]string, 0)
	if err := s.Client.Client().CallContext(
		ctx, &result, ExchangeCapabilities, &capabilities,
	); err != nil {
		return nil, err
	}
	return result, nil
}
