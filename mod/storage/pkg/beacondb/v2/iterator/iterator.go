package iterator

import (
	"fmt"

	"cosmossdk.io/store"
)

// Iterator is a wrapper around the store and changeset iterators,
// and provides an iterator which iterates over all changes in the
// changeset and store, skipping duplicates.
type iterator struct {
	changeset store.Iterator
	store     store.Iterator

	start []byte
	end   []byte

	seen map[string]struct{}
}

func New(start, end []byte, store store.Iterator, changeset store.Iterator) store.Iterator {
	return &iterator{
		changeset: changeset,
		store:     store,
		start:     start,
		end:       end,
		seen:      make(map[string]struct{}),
	}
}

func (i *iterator) Domain() (start []byte, end []byte) {
	return i.start, i.end
}

func (i *iterator) Valid() bool {
	return i.changeset.Valid() || i.store.Valid()
}

func (i *iterator) Next() {
	if i.changeset.Valid() {
		i.changeset.Next()
		i.seen[string(i.changeset.Key())] = struct{}{}
		return
	}
	i.store.Next()
	for i.store.Valid() {
		if _, ok := i.seen[string(i.store.Key())]; !ok {
			break
		}
		i.store.Next()
	}
}

func (i *iterator) Key() (key []byte) {
	if i.changeset.Valid() {
		fmt.Println("changeset key: ", i.changeset.Key())
		return i.changeset.Key()
	}
	fmt.Println("store key: ", i.store.Key())
	return i.store.Key()
}

func (i *iterator) Value() (value []byte) {
	if i.changeset.Valid() {
		return i.changeset.Value()
	}
	return i.store.Value()
}

func (i *iterator) Error() error {
	if i.changeset.Valid() {
		return i.changeset.Error()
	}
	return i.store.Error()
}

func (i *iterator) Close() error {
	if err := i.changeset.Close(); err != nil {
		return err
	}
	return i.store.Close()
}
