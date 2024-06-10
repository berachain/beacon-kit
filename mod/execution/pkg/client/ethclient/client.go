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

package ethclient

import (
	"context"
	"encoding/json"

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// Eth1Client is a struct that holds the Ethereum 1 client and
// its configuration.
type Eth1Client[
	ExecutionPayloadT interface {
		json.Marshaler
		json.Unmarshaler
		Empty(uint32) ExecutionPayloadT
	},
] struct {
	*ethclient.Client
}

// NewEth1Client creates a new Ethereum 1 client with the provided
// context and options.
func NewEth1Client[
	ExecutionPayloadT interface {
		json.Marshaler
		json.Unmarshaler
		Empty(uint32) ExecutionPayloadT
	},
](client *ethclient.Client) (*Eth1Client[ExecutionPayloadT], error) {
	c := &Eth1Client[ExecutionPayloadT]{
		Client: client,
	}
	return c, nil
}

// NewFromRPCClient creates a new Ethereum 1 client from an RPC client.
func NewFromRPCClient[
	ExecutionPayloadT interface {
		json.Marshaler
		json.Unmarshaler
		Empty(uint32) ExecutionPayloadT
	},
](rpcClient *rpc.Client) (*Eth1Client[ExecutionPayloadT], error) {
	return NewEth1Client[ExecutionPayloadT](ethclient.NewClient(rpcClient))
}

// NewPayloadV3 calls the engine_newPayloadV3 method via JSON-RPC.
func (s *Eth1Client[ExecutionPayloadT]) NewPayloadV3(
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
func (s *Eth1Client[ExecutionPayloadT]) ForkchoiceUpdatedV3(
	ctx context.Context,
	state *engineprimitives.ForkchoiceStateV1,
	attrs engineprimitives.PayloadAttributer,
) (*engineprimitives.ForkchoiceResponseV1, error) {
	return s.forkchoiceUpdateCall(ctx, ForkchoiceUpdatedMethodV3, state, attrs)
}

// forkchoiceUpdateCall is a helper function to call to any version
// of the forkchoiceUpdates method.
func (s *Eth1Client[ExecutionPayloadT]) forkchoiceUpdateCall(
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
func (s *Eth1Client[ExecutionPayloadT]) GetPayloadV3(
	ctx context.Context, payloadID engineprimitives.PayloadID,
) (engineprimitives.BuiltExecutionPayloadEnv[ExecutionPayloadT], error) {
	var t ExecutionPayloadT
	result := &engineprimitives.ExecutionPayloadEnvelope[
		ExecutionPayloadT,
		*engineprimitives.BlobsBundleV1[
			eip4844.KZGCommitment, eip4844.KZGProof, eip4844.Blob,
		],
	]{
		ExecutionPayload: t.Empty(version.Deneb),
	}

	if err := s.Client.Client().CallContext(
		ctx, result, GetPayloadMethodV3, payloadID,
	); err != nil {
		return nil, err
	}

	return result, nil
}

// ExecutionBlockByHash fetches an execution engine block by hash by calling
// eth_blockByHash via JSON-RPC.
func (s *Eth1Client[ExecutionPayloadT]) ExecutionBlockByHash(
	ctx context.Context, hash common.ExecutionHash, withTxs bool,
) (*engineprimitives.Block, error) {
	result := &engineprimitives.Block{}
	err := s.Client.Client().CallContext(
		ctx, result, BlockByHashMethod, hash, withTxs)
	return result, err
}

// ExecutionBlockByNumber fetches an execution engine block by number
// by calling eth_getBlockByNumber via JSON-RPC.
func (s *Eth1Client[ExecutionPayloadT]) ExecutionBlockByNumber(
	ctx context.Context, num rpc.BlockNumber, withTxs bool,
) (*engineprimitives.Block, error) {
	result := &engineprimitives.Block{}
	err := s.Client.Client().CallContext(
		ctx, result, BlockByNumberMethod, num, withTxs)
	return result, err
}

// GetClientVersionV1 calls the engine_getClientVersionV1 method via JSON-RPC.
func (s *Eth1Client[ExecutionPayloadT]) GetClientVersionV1(
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
func (s *Eth1Client[ExecutionPayloadT]) ExchangeCapabilities(
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
