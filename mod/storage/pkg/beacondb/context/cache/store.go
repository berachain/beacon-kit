package cache

import (
	"errors"
	"fmt"
	"time"

	"cosmossdk.io/core/store"
)

var _ store.Writer = (*Store[store.Reader])(nil)

// Store wraps an in-memory cache around an underlying types.KVStore.
type Store[T store.Reader] struct {
	cache     *cache // always ascending sorted
	changeSet *cache
	parent    T

	cacheHits   uint64
	cacheMisses uint64

	iteratorRangeCache *iteratorRangeCache[T]
}

// NewStore creates a new Store object
func NewStore[T store.Reader](parent T) *Store[T] {
	s := &Store[T]{
		cache:       newCache(),
		changeSet:   newCache(),
		parent:      parent,
		cacheHits:   0,
		cacheMisses: 0,
	}
	s.iteratorRangeCache = newIteratorRangeCache(parent, s.cache)
	return s
}

// Get implements types.KVStore.
func (s *Store[T]) Get(key []byte) (value []byte, err error) {
	value, found := s.cache.get(key)
	if found {
		s.cacheHits += 1
		return
	}
	s.cacheMisses += 1
	value, err = s.parent.Get(key)
	if err != nil {
		return nil, err
	}

	// add the value into the cache.
	s.cache.set(key, value)
	return value, nil
}

// Set implements types.KVStore.
func (s *Store[T]) Set(key, value []byte) error {
	if value == nil {
		return errors.New("cannot set a nil value")
	}

	s.cache.set(key, value)
	s.changeSet.set(key, value)
	return nil
}

// Has implements types.KVStore.
func (s *Store[T]) Has(key []byte) (bool, error) {
	tmpValue, found := s.cache.get(key)
	if found {
		return tmpValue != nil, nil
	}
	return s.parent.Has(key)
}

// Delete implements types.KVStore.
func (s *Store[T]) Delete(key []byte) error {
	s.cache.delete(key)
	s.changeSet.delete(key)
	return nil
}

// ----------------------------------------
// Iteration

// Iterator implements types.KVStore.
func (s *Store[T]) Iterator(start, end []byte) (store.Iterator, error) {
	return s.iterator(start, end, true)
}

// ReverseIterator implements types.KVStore.
func (s *Store[T]) ReverseIterator(start, end []byte) (store.Iterator, error) {
	return s.iterator(start, end, false)
}

var (
	ParentDuration time.Duration
	CacheDuration  time.Duration
)

func (s *Store[T]) iterator(start, end []byte, ascending bool) (store.Iterator, error) {
	// If the range has not been synced yet, sync it.
	if !s.iteratorRangeCache.Seen(start, end) {
		if err := s.iteratorRangeCache.SyncForRange(start, end); err != nil {
			return nil, err
		}
	}

	// Return the appropriate iterator.
	if ascending {
		return s.cache.Iterator(start, end)
	}
	return s.cache.ReverseIterator(start, end)
}

func (s *Store[T]) ApplyChangeSets(changes []store.KVPair) error {
	for _, c := range changes {
		if c.Remove {
			err := s.Delete(c.Key)
			if err != nil {
				return err
			}
		} else {
			err := s.Set(c.Key, c.Value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Store[T]) ChangeSets() (cs []store.KVPair, err error) {
	cs = make([]store.KVPair, s.changeSet.size())
	iter, err := s.changeSet.Iterator(nil, nil)
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	i := 0
	for ; iter.Valid(); iter.Next() {
		k, v := iter.Key(), iter.Value()
		cs[i] = store.KVPair{
			Key:    k,
			Value:  v,
			Remove: v == nil, // maybe we can optimistically compute size.
		}
		i++
	}
	fmt.Println("CACHE HITS IN CHANGESET", s.cacheHits)
	fmt.Println("CACHE MISSES IN CHANGESET", s.cacheMisses)
	return cs, nil
}
