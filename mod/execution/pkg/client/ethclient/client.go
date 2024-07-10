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

	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/geth-primitives/pkg/ethclient"
	"github.com/berachain/beacon-kit/mod/geth-primitives/pkg/rpc"
	engineprimitives "github.com/berachain/beacon-kit/mod/interfaces/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
)

// Eth1Client is a struct that holds the Ethereum 1 client and
// its configuration.
type Eth1Client[
	BlobsBundleT any,
	ExecutionPayloadT constraints.EngineType[ExecutionPayloadT],
	ExecutionPayloadEnvelopeT engineprimitives.ExecutionPayloadEnvelope[
		ExecutionPayloadEnvelopeT, BlobsBundleT, ExecutionPayloadT,
	],
] struct {
	*ethclient.Client
}

// NewEth1Client creates a new Ethereum 1 client with the provided
// context and options.
func NewEth1Client[
	BlobsBundleT any,
	ExecutionPayloadT constraints.EngineType[ExecutionPayloadT],
	ExecutionPayloadEnvelopeT engineprimitives.ExecutionPayloadEnvelope[
		ExecutionPayloadEnvelopeT, BlobsBundleT, ExecutionPayloadT,
	],
](client *ethclient.Client) (*Eth1Client[
	BlobsBundleT, ExecutionPayloadT, ExecutionPayloadEnvelopeT,
], error) {
	c := &Eth1Client[
		BlobsBundleT, ExecutionPayloadT, ExecutionPayloadEnvelopeT,
	]{
		Client: client,
	}
	return c, nil
}

// NewFromRPCClient creates a new Ethereum 1 client from an RPC client.
func NewFromRPCClient[
	BlobsBundleT any,
	ExecutionPayloadT constraints.EngineType[ExecutionPayloadT],
	ExecutionPayloadEnvelopeT engineprimitives.ExecutionPayloadEnvelope[
		ExecutionPayloadEnvelopeT, BlobsBundleT, ExecutionPayloadT,
	],
](rpcClient *rpc.Client) (*Eth1Client[
	BlobsBundleT, ExecutionPayloadT, ExecutionPayloadEnvelopeT,
], error) {
	return NewEth1Client[
		BlobsBundleT, ExecutionPayloadT, ExecutionPayloadEnvelopeT,
	](ethclient.NewClient(rpcClient))
}

// ExecutionBlockByHash fetches an execution engine block by hash by calling
// eth_blockByHash via JSON-RPC.
func (s *Eth1Client[_, _, _]) ExecutionBlockByHash(
	ctx context.Context, hash gethprimitives.ExecutionHash, withTxs bool,
) (*gethprimitives.Block, error) {
	result := &gethprimitives.Block{}
	err := s.Client.Client().CallContext(
		ctx, result, BlockByHashMethod, hash, withTxs)
	return result, err
}

// ExecutionBlockByNumber fetches an execution engine block by number
// by calling eth_getBlockByNumber via JSON-RPC.
func (s *Eth1Client[_, _, _]) ExecutionBlockByNumber(
	ctx context.Context, num rpc.BlockNumber, withTxs bool,
) (*gethprimitives.Block, error) {
	result := &gethprimitives.Block{}
	err := s.Client.Client().CallContext(
		ctx, result, BlockByNumberMethod, num, withTxs)
	return result, err
}
