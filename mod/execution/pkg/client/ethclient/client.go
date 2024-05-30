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

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// Eth1Client is a struct that holds the Ethereum 1 client and
// its configuration.
type Eth1Client[
	ExecutionPayloadDenebT engineprimitives.ExecutionPayload,
] struct {
	*ethclient.Client
}

// NewEth1Client creates a new Ethereum 1 client with the provided
// context and options.
func NewEth1Client[
	ExecutionPayloadDenebT engineprimitives.ExecutionPayload,
](client *ethclient.Client) (*Eth1Client[ExecutionPayloadDenebT], error) {
	c := &Eth1Client[ExecutionPayloadDenebT]{
		Client: client,
	}
	return c, nil
}

// NewFromRPCClient creates a new Ethereum 1 client from an RPC client.
func NewFromRPCClient[
	ExecutionPayloadDenebT engineprimitives.ExecutionPayload,
](rpcClient *rpc.Client) (*Eth1Client[ExecutionPayloadDenebT], error) {
	return NewEth1Client[ExecutionPayloadDenebT](ethclient.NewClient(rpcClient))
}

// NewPayloadV3 calls the engine_newPayloadV3 method via JSON-RPC.
func (s *Eth1Client[ExecutionPayloadDenebT]) NewPayloadV3(
	ctx context.Context,
	payload any,
	versionedHashes []common.ExecutionHash,
	parentBlockRoot *primitives.Root,
) (*engineprimitives.PayloadStatusV1, error) {
	result := &engineprimitives.PayloadStatusV1{}
	if err := s.Client.Client().CallContext(
		ctx, result, NewPayloadMethodV3, payload, versionedHashes,
		(*common.ExecutionHash)(parentBlockRoot),
	); err != nil {
		return nil, err
	}
	return result, nil
}

// ForkchoiceUpdatedV3 calls the engine_forkchoiceUpdatedV3 method via JSON-RPC.
func (s *Eth1Client[ExecutionPayloadDenebT]) ForkchoiceUpdatedV3(
	ctx context.Context,
	state *engineprimitives.ForkchoiceStateV1,
	attrs engineprimitives.PayloadAttributer,
) (*engineprimitives.ForkchoiceResponseV1, error) {
	return s.forkchoiceUpdateCall(ctx, ForkchoiceUpdatedMethodV3, state, attrs)
}

// forkchoiceUpdateCall is a helper function to call to any version
// of the forkchoiceUpdates method.
func (s *Eth1Client[ExecutionPayloadDenebT]) forkchoiceUpdateCall(
	ctx context.Context,
	method string,
	state *engineprimitives.ForkchoiceStateV1,
	attrs any,
) (*engineprimitives.ForkchoiceResponseV1, error) {
	result := &engineprimitives.ForkchoiceResponseV1{}

	if err := s.Client.Client().CallContext(
		ctx, result, method, state, attrs,
	); err != nil {
		return nil, err
	}

	if (result.PayloadStatus == engineprimitives.PayloadStatusV1{}) {
		return nil, ErrNilResponse
	}

	return result, nil
}

// GetPayloadV3 calls the engine_getPayloadV3 method via JSON-RPC.
func (s *Eth1Client[ExecutionPayloadDenebT]) GetPayloadV3(
	ctx context.Context, payloadID engineprimitives.PayloadID,
) (engineprimitives.BuiltExecutionPayloadEnv, error) {
	result := &engineprimitives.ExecutionPayloadEnvelope[
		ExecutionPayloadDenebT,
		*engineprimitives.BlobsBundleV1[
			eip4844.KZGCommitment, eip4844.KZGProof, eip4844.Blob,
		],
	]{}

	if err := s.Client.Client().CallContext(
		ctx, result, GetPayloadMethodV3, payloadID,
	); err != nil {
		return nil, err
	}
	return result, nil
}

// ExecutionBlockByHash fetches an execution engine block by hash by calling
// eth_blockByHash via JSON-RPC.
func (s *Eth1Client[ExecutionPayloadDenebT]) ExecutionBlockByHash(
	ctx context.Context, hash common.ExecutionHash, withTxs bool,
) (*engineprimitives.Block, error) {
	result := &engineprimitives.Block{}
	err := s.Client.Client().CallContext(
		ctx, result, BlockByHashMethod, hash, withTxs)
	return result, err
}

// ExecutionBlockByNumber fetches an execution engine block by number
// by calling eth_getBlockByNumber via JSON-RPC.
func (s *Eth1Client[ExecutionPayloadDenebT]) ExecutionBlockByNumber(
	ctx context.Context, num rpc.BlockNumber, withTxs bool,
) (*engineprimitives.Block, error) {
	result := &engineprimitives.Block{}
	err := s.Client.Client().CallContext(
		ctx, result, BlockByNumberMethod, num, withTxs)
	return result, err
}

// GetClientVersionV1 calls the engine_getClientVersionV1 method via JSON-RPC.
func (s *Eth1Client[ExecutionPayloadDenebT]) GetClientVersionV1(
	ctx context.Context,
) ([]engineprimitives.ClientVersionV1, error) {
	result := make([]engineprimitives.ClientVersionV1, 0)
	if err := s.Client.Client().CallContext(
		ctx, &result, GetClientVersionV1, nil,
	); err != nil {
		return nil, err
	}
	return result, nil
}

// ExchangeCapabilities calls the engine_exchangeCapabilities method via
// JSON-RPC.
func (s *Eth1Client[ExecutionPayloadDenebT]) ExchangeCapabilities(
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
