package store

import (
	"cosmossdk.io/store"

	"github.com/ethereum/go-ethereum/common"
)

type Genesis struct {
	store store.KVStore
}

func NewGenesis(store store.KVStore) *Genesis {
	return &Genesis{
		store: store,
	}
}

func (f *Genesis) Store(eth1GenesisHash string) error {
	f.store.Set([]byte("eth1_genesis_hash"), []byte(eth1GenesisHash))
	return nil
}

func (f *Genesis) Retrieve() common.Hash {
	bz := f.store.Get([]byte("eth1_genesis_hash"))
	return common.HexToHash(string(bz))
}
