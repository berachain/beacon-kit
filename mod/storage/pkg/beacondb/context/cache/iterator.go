package cache

import (
	"bytes"
	"encoding/hex"
	"strings"

	"cosmossdk.io/core/store"
)

// iteratorRangeCache is a cache for iterators that have been synced
// for a given range.
type iteratorRangeCache[T store.Reader] struct {
	parent T
	cache  *cache

	syncedRanges map[string]struct{}
}

// newIteratorRangeCache creates a new iterator range cache.
func newIteratorRangeCache[T store.Reader](
	parent T,
	cache *cache,
) *iteratorRangeCache[T] {
	return &iteratorRangeCache[T]{
		parent:       parent,
		cache:        cache,
		syncedRanges: make(map[string]struct{}),
	}
}

// Seen returns true if the iterator range cache has been synced for
// the given range.
func (c *iteratorRangeCache[T]) Seen(start, end []byte) bool {
	if len(c.syncedRanges) == 0 {
		return false
	}
	key := newIteratorRangeCacheKey(start, end)
	if _, ok := c.syncedRanges[key.String()]; ok {
		return true
	}

	for seen := range c.syncedRanges {
		other, err := newIteratorRangeCacheKeyFromString(seen)
		if err != nil {
			continue
		}
		if key.inRange(other) {
			return true
		}
	}
	return false
}

// SyncForRange syncs the values stored in the parent store with
// the values stored in the cache for the given iteration domain.
// The cache 'shadows' the parent, so if the cache contains a
// value for a key, the parent is ignored. Returns a key that
// represents the range of the iteration that was just synced.
//
// side effects: modifies the underlying cache store in place.
func (c *iteratorRangeCache[T]) SyncForRange(
	start, end []byte,
) error {
	parentIter, err := c.parent.Iterator(start, end)
	if err != nil {
		return err
	}
	defer parentIter.Close()

	for parentIter.Valid() {
		if _, ok := c.cache.get(parentIter.Key()); ok {
			parentIter.Next()
			continue
		}
		c.cache.set(parentIter.Key(), parentIter.Value())
		parentIter.Next()
	}
	// mark the range as synced
	c.syncedRanges[newIteratorRangeCacheKey(start, end).String()] = struct{}{}
	return nil
}

// iteratorRangeCacheKey is a key used to track ranges of iterators
// that have been synced already.
type iteratorRangeCacheKey struct {
	start, end []byte
}

// newIteratorRangeCacheKey creates a new iterator range cache key from
// the given start and end bytes.
func newIteratorRangeCacheKey(start, end []byte) *iteratorRangeCacheKey {
	return &iteratorRangeCacheKey{start: start, end: end}
}

// newIteratorRangeCacheKeyFromString creates a new iterator range cache
// key from a string of the form <0xstart>-<0xend>.
func newIteratorRangeCacheKeyFromString(
	key string,
) (*iteratorRangeCacheKey, error) {
	parts := strings.Split(key, "-")
	if len(parts) != 2 {
		return nil, errInvalidIteratorRangeCacheKey
	}
	start, err := hex.DecodeString(parts[0])
	if err != nil {
		return nil, err
	}
	end, err := hex.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}
	return newIteratorRangeCacheKey(start, end), nil
}

// String formats the iterator range cache key as a string
// of the form <0xstart>-<0xend>.
func (k *iteratorRangeCacheKey) String() string {
	return strings.Join(
		[]string{
			hex.EncodeToString(k.start),
			hex.EncodeToString(k.end),
		},
		"-",
	)
}

// inRange returns true if k is within the range of other, inclusive.
// In other words, if k.start is greater than or equal to other.start
// and k.end is less than or equal to other.end.
func (k *iteratorRangeCacheKey) inRange(other *iteratorRangeCacheKey) bool {
	return bytes.Compare(k.start, other.start) >= 0 &&
		bytes.Compare(k.end, other.end) <= 0
}
