package provider

import (
	"context"

	sdkcollections "cosmossdk.io/collections"
	beacontypes "github.com/berachain/beacon-kit/beacon/core/types"
	"github.com/berachain/beacon-kit/lib/store/collections/encoding"
	"github.com/berachain/beacon-kit/light/store"
	"github.com/berachain/beacon-kit/primitives"
	"github.com/ethereum/go-ethereum/common"
)

const (
	finalized_key                 = "fc_finalized"
	blockRootsPrefix              = "block_roots"
	latestBeaconBlockHeaderPrefix = "latest_beacon_block_header"
)

var (
	latestBeaconBlockHeaderCodec = store.Codec[store.None, *beacontypes.BeaconBlockHeader]{
		Value: encoding.SSZValueCodec[*beacontypes.BeaconBlockHeader]{},
	}

	blockRootsCodec = store.Codec[uint64, [32]byte]{
		Key:   sdkcollections.Uint64Key,
		Value: encoding.Bytes32ValueCodec{},
	}
)

func (p *Provider) GetTrustedEth1Hash() common.Hash {
	res, err := p.Query(context.Background(), []byte(finalized_key), p.latestBlockHeight)
	if err != nil {
		panic(err)
	}
	return common.BytesToHash(res)
}

func (p *Provider) GetBlockRootAtIndex(index uint64) primitives.Root {
	bytesKey, err := sdkcollections.EncodeKeyWithPrefix(
		[]byte(blockRootsPrefix), blockRootsCodec.Key, index,
	)
	if err != nil {
		panic(err)
	}

	res, err := p.Query(context.Background(), bytesKey, 0)
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

func (p *Provider) GetLatestBlockHeader() *beacontypes.BeaconBlockHeader {
	res, err := p.Query(context.Background(), []byte(latestBeaconBlockHeaderPrefix), p.latestBlockHeight)
	if err != nil {
		panic(err)
	}

	header, err := latestBeaconBlockHeaderCodec.Value.Decode(res)
	if err != nil {
		panic(err)
	}
	return header
}
