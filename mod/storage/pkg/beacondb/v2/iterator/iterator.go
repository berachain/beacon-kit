package iterator

import (
	"cosmossdk.io/runtime/v2"
	"cosmossdk.io/store"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb/v2/changeset"
)

type Iterator interface {
	// Domain returns the start (inclusive) and end (exclusive) limits of the iterator.
	// CONTRACT: start, end readonly []byte
	Domain() (start []byte, end []byte)

	// Valid returns whether the current iterator is valid. Once invalid, the Iterator remains
	// invalid forever.
	Valid() bool

	// Next moves the iterator to the next key in the database, as defined by order of iteration.
	// If Valid returns false, this method will panic.
	Next()

	// Key returns the key at the current position. Panics if the iterator is invalid.
	// CONTRACT: key readonly []byte
	Key() (key []byte)

	// Value returns the value at the current position. Panics if the iterator is invalid.
	// CONTRACT: value readonly []byte
	Value() (value []byte)

	// Error returns the last error encountered by the iterator, if any.
	Error() error

	// Close closes the iterator, relasing any allocated resources.
	Close() error
}

type iterator struct {
	changeset *changeset.Changeset
	store     runtime.Store

	start []byte
	end   []byte
}

func New(start, end []byte, store runtime.Store, changeset *changeset.Changeset) store.Iterator {
	return &iterator{
		changeset: changeset,
		store:     store,
		start:     start,
		end:       end,
	}
}

func (i *iterator) Domain() (start []byte, end []byte) {
	return i.start, i.end
}

func (i *iterator) Valid() bool {
	return false
}

func (i *iterator) Next() {
}

func (i *iterator) Key() (key []byte) {
	return nil
}

func (i *iterator) Value() (value []byte) {
	return nil
}

func (i *iterator) Error() error {
	return nil
}

func (i *iterator) Close() error {
	return nil
}
