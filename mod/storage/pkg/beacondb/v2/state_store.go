package beacondb

import (
	sdkcollections "cosmossdk.io/collections"
	"cosmossdk.io/runtime/v2"
	"cosmossdk.io/store"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb/v2/changeset"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb/v2/iterator"
	"github.com/berachain/beacon-kit/mod/storage/pkg/collections"
)

type StateStore struct {
	runtime.Store
	changeSet *changeset.Changeset
}

func NewStore() *StateStore {
	return &StateStore{
		changeSet: changeset.New(),
	}
}

func (s *StateStore) AddChange(storeKey []byte, key []byte, value []byte) {
	s.changeSet.Add(storeKey, key, value, false)
}

func (s *StateStore) QueryState(storeKey, key []byte) ([]byte, error) {
	// first query the change set
	value, found := s.changeSet.Query(storeKey, key)
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
	// cs := changeset.New()
	// cs.Add([]byte("beacon"), []byte("umst"), []byte("umst"), false)
	// store.Commit(cs.GetChanges())
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

	changeSetIter, err := s.changeSet.Iterator(start, end)
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
		s.changeSet.Flush()
	}()
	if s.changeSet.Size() == 0 {
		return
	}
	s.Store.Commit(s.changeSet.GetChanges())
}

// TODO: hack to get around IAVL versioning issues (a call to save must be made
// before any queries to SC)
func (s *StateStore) Init() {
	s.AddChange([]byte("beacon"), []byte("umst"), []byte("umst"))
	s.Save()
}
