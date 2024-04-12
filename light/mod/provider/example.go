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

package provider

import (
	"context"

	sdkcollections "cosmossdk.io/collections"
	lightkeys "github.com/berachain/beacon-kit/light/mod/storage/beacondb/keys"
	"github.com/berachain/beacon-kit/light/mod/storage/codec"
	"github.com/berachain/beacon-kit/mod/execution/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/storage/beacondb/collections/encoding"
	"github.com/berachain/beacon-kit/mod/storage/beacondb/keys"
)

// NOTE: THIS FILE IS TEMPORARY AND IS USED TO DEMONSTRATE HOW
// THE LIGHT CLIENT CAN ACCESS STATE DATA FROM THE PROVIDER

//nolint:gochecknoglobals // this will be removed in a later pr
var (
	latestExecutionPayloadCodec = codec.NewItem[types.ExecutionPayload](
		keys.LatestExecutionPayloadPrefix,
		encoding.SSZInterfaceCodec[types.ExecutionPayload]{
			Factory: func() types.ExecutionPayload {
				return &types.ExecutableDataDeneb{}
			},
		},
	)

	latestBeaconBlockHeaderCodec = codec.NewItem[*primitives.BeaconBlockHeader](
		keys.LatestBeaconBlockHeaderPrefix,
		encoding.SSZValueCodec[*primitives.BeaconBlockHeader]{},
	)

	blockRootsCodec = codec.NewMap[uint64, [32]byte](
		keys.BlockRootsPrefix,
		sdkcollections.Uint64Key,
		encoding.Bytes32ValueCodec{},
	)
)

// Querying from sdkcollections.Item.
func (p *Provider) GetEth1BlockHash(
	ctx context.Context,
) primitives.ExecutionHash {
	res, err := p.Query(
		ctx,
		lightkeys.BeaconStoreKey,
		[]byte(keys.LatestExecutionPayloadPrefix),
		0,
	)
	if err != nil {
		panic(err)
	}

	payload, err := latestExecutionPayloadCodec.Decode(res)
	if err != nil {
		panic(err)
	}

	return payload.GetBlockHash()
}

// Another example of querying from sdkcollections.Item.
func (p *Provider) GetLatestBlockHeader(
	ctx context.Context,
) *primitives.BeaconBlockHeader {
	res, err := p.Query(
		ctx,
		lightkeys.BeaconStoreKey,
		[]byte(keys.LatestBeaconBlockHeaderPrefix),
		0,
	)
	if err != nil {
		panic(err)
	}

	header, err := latestBeaconBlockHeaderCodec.Decode(res)
	if err != nil {
		panic(err)
	}
	return header
}

// Querying from sdkcollections.Map.
func (p *Provider) GetBlockRootAtIndex(
	ctx context.Context, index uint64,
) primitives.Root {
	key, err := blockRootsCodec.Key(index)
	if err != nil {
		panic(err)
	}

	res, err := p.Query(ctx, lightkeys.BeaconStoreKey, key, 0)
	if err != nil {
		panic(err)
	}
	if res == nil {
		return primitives.Root{}
	}

	blockRoot, err := blockRootsCodec.Decode(res)
	if err != nil {
		panic(err)
	}

	return blockRoot
}
