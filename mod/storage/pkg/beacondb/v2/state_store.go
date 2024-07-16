package beacondb

import (
	"context"

	sdkcollections "cosmossdk.io/collections"
	"cosmossdk.io/runtime/v2"
	"cosmossdk.io/store"
	"github.com/berachain/beacon-kit/mod/errors"
	storectx "github.com/berachain/beacon-kit/mod/storage/pkg/beacondb/v2/context"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb/v2/iterator"
	"github.com/berachain/beacon-kit/mod/storage/pkg/collections"
)

type StateStore struct {
	ctx *storectx.Context
	runtime.Store
}

func NewStore() *StateStore {
	return &StateStore{}
}

func (s *StateStore) AddChange(storeKey []byte, key []byte, value []byte) {
	s.ctx.Changeset.Add(storeKey, key, value, false)
}

func (s *StateStore) QueryState(storeKey, key []byte) ([]byte, error) {
	// first query the change set
	value, found := s.ctx.Changeset.Query(storeKey, key)
	if found {
		return value, nil
	}
	// query the underlying store with the latest version
	version, err := s.GetLatestVersion()
	if err != nil {
		return nil, err
	}
	resp, err := s.Query(storeKey, version, key, false)
	if err != nil {
		// TODO: clean this up
		if errors.Is(err, sdkcollections.ErrNotFound) {
			return nil, collections.ErrNotFound
		}
		return nil, err
	}
	if resp.Value == nil {
		return nil, collections.ErrNotFound
	}
	return resp.Value, nil
}

func (s *StateStore) SetStore(store runtime.Store) {
	s.Store = store
}

func (s *StateStore) Iterator(start, end []byte) (store.Iterator, error) {
	_, readerMap, err := s.StateLatest()
	if err != nil {
		return nil, err
	}
	// get reader with the storeKey from reader map
	reader, err := readerMap.GetReader([]byte(ModuleName))
	if err != nil {
		return nil, err
	}
	// get iterator from reader with the prefixed start and end
	stateIter, err := reader.Iterator(start, end)
	if err != nil {
		return nil, err
	}

	changeSetIter, err := s.ctx.Changeset.Iterator(start, end)
	if err != nil {
		return nil, err
	}

	return iterator.New(start, end, stateIter, changeSetIter), nil
}

// if commit errors should we still reset? maybe just do an
// explicit call instead of defer to prevent that case
// TODO: return store hash
func (s *StateStore) Save() {
	// reset the changeset following the commit
	defer func() {
		s.ctx.Changeset.Flush()
	}()
	size := s.ctx.Changeset.Size()
	if size == 0 {
		return
	}
	s.Store.Commit(s.ctx.Changeset.GetChanges())
}

func (s *StateStore) Context() context.Context {
	return s.ctx
}

func (s *StateStore) WithContext(ctx context.Context) *StateStore {
	storeCtx, ok := ctx.(*storectx.Context)
	if !ok {
		storeCtx = storectx.New(ctx)
	}
	cpy := *s
	s.ctx = storeCtx
	return &cpy
}

// TODO: hack to get around IAVL versioning issues (a call to save must be made
// before any queries to SC)
func (s *StateStore) Init() {
	s.ctx = storectx.New(context.Background())
	s.AddChange([]byte("beacon"), []byte("umst"), []byte("umst"))
	s.Save()
}
