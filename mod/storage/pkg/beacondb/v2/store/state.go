package store

import (
	"context"

	sdkcollections "cosmossdk.io/collections"
	corestore "cosmossdk.io/core/store"
	"cosmossdk.io/runtime/v2"
	"cosmossdk.io/store"
	storev2 "cosmossdk.io/store/v2"
	"github.com/berachain/beacon-kit/mod/errors"
	storectx "github.com/berachain/beacon-kit/mod/storage/pkg/beacondb/v2/context"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb/v2/iterator"
	"github.com/berachain/beacon-kit/mod/storage/pkg/collections"
)

// TODO: pull this from some cfg or something idk remind me to do it after
// im done getting shafted by server v2
const (
	ModuleName = "beacon"
)

// StateStore is a store for the state of the beacon
type StateStore struct {
	ctx            *storectx.Context
	transientState *BlockStore
	state          runtime.Store
}

// NewStore creates a new state store
func NewStore(blockStore *BlockStore) *StateStore {
	return &StateStore{
		transientState: blockStore,
	}
}

// AddChange adds a change to the changeset
func (s *StateStore) AddChange(storeKey []byte, key []byte, value []byte) {
	s.ctx.Changeset.Add(storeKey, key, value, false)
}

// TODO: CHANGE NAME TO QUERY
// QueryState queries the state store for the given key
// It first queries the changeset, then the block store, then the chain store
func (s *StateStore) QueryState(storeKey, key []byte) ([]byte, error) {
	// QUERY CHANGESET
	value, found := s.ctx.Changeset.Query(storeKey, key)
	if found {
		return value, nil
	}

	// NOT FOUND IN CHANGESET -> QUERY BLOCK STORE
	value, found = s.transientState.Query(storeKey, key)
	if found {
		return value, nil
	}

	// NOT FOUND IN BLOCK STORE -> QUERY COMMIT STORE
	version, err := s.state.GetLatestVersion()
	if err != nil {
		return nil, err
	}
	// if the version is 0, we're in genesis
	if version != 0 {
		var resp storev2.QueryResult
		resp, err = s.state.Query(storeKey, version, key, false)
		if err != nil {
			// TODO: clean this up
			if errors.Is(err, sdkcollections.ErrNotFound) {
				return nil, collections.ErrNotFound
			}
			return nil, err
		}
		if resp.Value != nil {
			return resp.Value, nil
		}
	}
	return nil, collections.ErrNotFound
}

func (s *StateStore) SetStore(store runtime.Store) {
	s.state = store
}

// Iterator returns a combined iterator over the the changeset, the block store,
// and the chain store
func (s *StateStore) Iterator(start, end []byte) (store.Iterator, error) {
	// get chain store iterator
	_, readerMap, err := s.state.StateLatest()
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

	// get chainset iterator
	changeSetIter, err := s.ctx.Changeset.Iterator(start, end)
	if err != nil {
		return nil, err
	}

	// get block store iterator
	blockStoreIter, err := s.transientState.Iterator(start, end)
	if err != nil {
		return nil, err
	}
	return iterator.New(start, end, stateIter, changeSetIter, blockStoreIter), nil
}

// Save commits the changeset to the BlockStore and resets the changeset
func (s *StateStore) Save() {
	// reset the changeset following the commit
	defer func() {
		s.ctx.Changeset.Flush()
	}()
	size := s.ctx.Changeset.Size()
	if size == 0 {
		return
	}
	// commit changes to block store
	s.transientState.Commit(s.ctx.Changeset.GetChanges())
}

// Context returns the context of the StateStore
func (s *StateStore) Context() *storectx.Context {
	return s.ctx
}

// WithContext returns a new StateStore with the given context
// Invariant: the blockStore in the retuned copy must always point to the same
// blockStore as the current state
func (s *StateStore) WithContext(ctx context.Context) *StateStore {
	storeCtx, ok := ctx.(*storectx.Context)
	if !ok {
		storeCtx = storectx.New(ctx)
	}
	// initialize a new StateStore with the current block store and chain store
	storeCopy := NewStore(s.transientState)
	storeCopy.SetStore(s.state)
	// set the context to the given context
	storeCopy.ctx = storeCtx
	return storeCopy
}

// TODO V2: i'll remove this once i confirm whatever sdfjklhsdfjkl
// TODO: hack to get around IAVL versioning issues (a call to save must be made
// before any queries to SC)
func (s *StateStore) Init() {
	s.ctx = storectx.New(context.Background())
	// s.AddChange([]byte("beacon"), []byte("4206969666"), []byte("69"))
	// s.Save()
}

// TODO: remove this when no needed
// StateLatest returns the latest state from the chain store
func (s *StateStore) StateLatest() (uint64, corestore.ReaderMap, error) {
	return s.state.StateLatest()
}
