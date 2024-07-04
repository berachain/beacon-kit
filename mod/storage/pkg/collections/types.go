package collections

import "cosmossdk.io/store/v2"

type Store interface {
	store.RootStore
	AddChange([]byte, []byte, []byte)
}

type StoreAccessor func() Store
