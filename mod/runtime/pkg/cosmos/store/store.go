package store

import storetypes "cosmossdk.io/store/types"

type CommitMultiStore struct {
	storetypes.CommitMultiStore
}

func NewCommitMultiStore(ms storetypes.CommitMultiStore) *CommitMultiStore {
	return &CommitMultiStore{ms}
}

func (cms *CommitMultiStore) CacheMultiStore() storetypes.CacheMultiStore {
	return cms.CommitMultiStore.CacheMultiStore()
}
