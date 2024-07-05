package collections

import (
	"cosmossdk.io/runtime/v2"
)

type Store interface {
	runtime.Store
	AddChange([]byte, []byte, []byte)
}

type StoreAccessor func() Store
