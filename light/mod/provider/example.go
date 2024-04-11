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
	"github.com/berachain/beacon-kit/light/mod/state/codec"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/storage/beacondb/collections/encoding"
	"github.com/berachain/beacon-kit/mod/storage/beacondb/keys"
)

// NOTE: THIS FILE IS TEMPORARY AND IS USED TO DEMONSTRATE HOW
// THE LIGHT CLIENT CAN ACCESS STATE DATA FROM THE PROVIDER

//nolint:gochecknoglobals // this will be removed in a later pr
var (
	eth1BlockHashCodec = codec.Codec[
		codec.None, [32]byte,
	]{
		Value: encoding.Bytes32ValueCodec{},
	}

	latestBeaconBlockHeaderCodec = codec.Codec[
		codec.None, *primitives.BeaconBlockHeader,
	]{
		Value: encoding.SSZValueCodec[*primitives.BeaconBlockHeader]{},
	}

	blockRootsCodec = codec.Codec[uint64, [32]byte]{
		Key:   sdkcollections.Uint64Key,
		Value: encoding.Bytes32ValueCodec{},
	}
)

// Querying from sdkcollections.Item.
func (p *Provider) GetEth1BlockHash(
	ctx context.Context,
) primitives.ExecutionHash {
	res, err := p.Query(ctx, []byte(keys.Eth1BlockHashPrefix), 0)
	if err != nil {
		panic(err)
	}

	hash, err := eth1BlockHashCodec.Value.Decode(res)
	if err != nil {
		panic(err)
	}

	return hash
}

// Another example of querying from sdkcollections.Item.
func (p *Provider) GetLatestBlockHeader(
	ctx context.Context,
) *primitives.BeaconBlockHeader {
	res, err := p.Query(ctx, []byte(keys.LatestBeaconBlockHeaderPrefix), 0)
	if err != nil {
		panic(err)
	}

	header, err := latestBeaconBlockHeaderCodec.Value.Decode(res)
	if err != nil {
		panic(err)
	}
	return header
}

// Querying from sdkcollections.Map.
func (p *Provider) GetBlockRootAtIndex(
	ctx context.Context, index uint64,
) primitives.Root {
	bytesKey, err := sdkcollections.EncodeKeyWithPrefix(
		[]byte(keys.BlockRootsPrefix), blockRootsCodec.Key, index,
	)
	if err != nil {
		panic(err)
	}

	res, err := p.Query(ctx, bytesKey, 0)
	if err != nil {
		panic(err)
	}
	if res == nil {
		return primitives.Root{}
	}

	blockRoot, err := blockRootsCodec.Value.Decode(res)
	if err != nil {
		panic(err)
	}

	return blockRoot
}
