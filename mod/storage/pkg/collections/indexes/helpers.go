package indexes

import (
	"github.com/berachain/beacon-kit/mod/storage/pkg/collections"
)

// iterator defines the minimum set of methods of an index iterator
// required to work with the helpers.
type iterator[K any] interface {
	// PrimaryKey returns the iterator current primary key.
	PrimaryKey() (K, error)
	// Next advances the iterator by one element.
	Next()
	// Valid asserts if the Iterator is valid.
	Valid() bool
	// Close closes the iterator.
	Close() error
}

// ScanValues collects all the values from an Index iterator and the IndexedMap in a lazy way.
// The iterator is closed when this function exits.
func ScanValues[K, V any, I iterator[K], Idx collections.Indexes[K, V]](
	indexedMap *collections.IndexedMap[K, V, Idx],
	iter I,
	f func(value V) (stop bool),
) error {
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		key, err := iter.PrimaryKey()
		if err != nil {
			return err
		}

		value, err := indexedMap.Get(key)
		if err != nil {
			return err
		}

		stop := f(value)
		if stop {
			return nil
		}
	}

	return nil
}
