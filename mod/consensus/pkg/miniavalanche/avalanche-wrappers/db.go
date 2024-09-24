package avalanchewrappers

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/ava-labs/avalanchego/database"
	dbm "github.com/cosmos/cosmos-db"
)

var (
	_ database.Database = (*DB)(nil)
	_ database.Batch    = (*Batch)(nil)
	_ database.Iterator = (*Iterator)(nil)

	errDBFeatureNotImplementedYet = errors.New("db feature not implemented yet")
)

type DB struct {
	cosmosDB *dbm.PrefixDB
}

func NewDB(cosmosDB *dbm.PrefixDB) *DB {
	return &DB{
		cosmosDB: cosmosDB,
	}
}

func (db *DB) UnderlyingDB() dbm.DB {
	return db.cosmosDB
}

func (db *DB) Has(key []byte) (bool, error) {
	return db.cosmosDB.Has(key)
}

func (db *DB) Get(key []byte) ([]byte, error) {
	res, err := db.cosmosDB.Get(key)
	if res == nil && err == nil {
		return nil, database.ErrNotFound
	}
	return res, err
}

func (db *DB) Put(key []byte, value []byte) error {
	return db.cosmosDB.Set(key, value)
}

func (db *DB) Delete(key []byte) error {
	return db.cosmosDB.Delete(key)
}

func (db *DB) NewBatch() database.Batch {
	return &Batch{
		added:   map[string][]byte{},
		deleted: map[string]struct{}{},
		batch:   db.cosmosDB.NewBatch(),
		db:      db.cosmosDB,
	}
}

func (db *DB) NewIterator() database.Iterator {
	it, err := db.cosmosDB.Iterator(nil, nil)
	if err != nil {
		panic(fmt.Errorf("failed creating dbm iterator: %w", err))
	}
	return &Iterator{
		it: it,
	}
}

func (db *DB) NewIteratorWithStart(start []byte) database.Iterator {
	it, err := db.cosmosDB.Iterator(start, nil)
	if err != nil {
		panic(fmt.Errorf("failed creating dbm iterator with start: %w", err))
	}
	return &Iterator{
		it: it,
	}
}

func (db *DB) NewIteratorWithPrefix(prefix []byte) database.Iterator {
	it, err := db.cosmosDB.Iterator(nil, nil)
	if err != nil {
		panic(fmt.Errorf("failed creating dbm iterator with prefix: %w", err))
	}
	return &Iterator{
		prefix: prefix,
		it:     it,
	}
}

func (db *DB) NewIteratorWithStartAndPrefix(
	start, prefix []byte,
) database.Iterator {
	it, err := db.cosmosDB.Iterator(start, nil)
	if err != nil {
		panic(fmt.Errorf(
			"failed creating dbm iterator with start and prefix: %w",
			err,
		))
	}
	return &Iterator{
		prefix: prefix,
		it:     it,
	}
}

func (db *DB) Compact([]byte, []byte) error {
	return nil // TODO: check if it can be implemented somehow
}

func (db *DB) Close() error {
	return db.cosmosDB.Close()
}

func (db *DB) HealthCheck(context.Context) (interface{}, error) {
	return db.cosmosDB.Stats(), nil
}

type Batch struct {
	added   map[string][]byte   // added or update elements
	deleted map[string]struct{} // deleted elements

	db    dbm.DB
	batch dbm.Batch
}

func (b *Batch) Has(key []byte) (bool, error) {
	if _, found := b.added[string(key)]; found {
		return true, nil
	}
	if _, found := b.deleted[string(key)]; found {
		return false, nil
	}
	return b.db.Has(key)
}

func (b *Batch) Get(key []byte) ([]byte, error) {
	if v, found := b.added[string(key)]; found {
		return v, nil
	}
	if _, found := b.deleted[string(key)]; found {
		return nil, database.ErrNotFound
	}
	return b.db.Get(key)
}

func (b *Batch) Put(key []byte, value []byte) error {
	if err := b.batch.Set(key, value); err != nil {
		return err
	}
	b.added[string(key)] = value
	delete(b.deleted, string(key))
	return nil
}

func (b *Batch) Delete(key []byte) error {
	if err := b.batch.Delete(key); err != nil {
		return err
	}
	delete(b.added, string(key))
	b.deleted[string(key)] = struct{}{}
	return nil
}

func (b *Batch) Size() int {
	return len(b.added) + len(b.deleted)
}

func (b *Batch) Write() error {
	return b.batch.Write()
}

func (b *Batch) Reset() {
	b.added = make(map[string][]byte)
	b.deleted = map[string]struct{}{}
	if err := b.batch.Close(); err != nil {
		panic(fmt.Errorf("failed closing batch: %w", err))
	}
	b.batch = b.db.NewBatch()
}

func (b *Batch) Replay(database.KeyValueWriterDeleter) error {
	return errDBFeatureNotImplementedYet
}

func (b *Batch) Inner() database.Batch {
	return b // TODO: check this mapping is correct
}

type Iterator struct {
	prefix []byte
	it     dbm.Iterator
}

func (i *Iterator) Next() bool {
	switch {
	case len(i.prefix) == 0:
		i.it.Next()
		return i.it.Valid()
	default:
		i.it.Next()
		for i.it.Valid() {
			if bytes.HasPrefix(i.it.Value(), i.prefix) {
				return true
			}
			i.it.Next()
		}
		return false
	}
}

func (i *Iterator) Error() error {
	return i.it.Error()
}

func (i *Iterator) Key() []byte {
	if i.it.Valid() {
		return i.it.Key()
	}

	// there are cases whey Key is called even if Next() is false
	// (see meterDB for instance)
	return nil
}

func (i *Iterator) Value() []byte {
	if i.it.Valid() {
		return i.it.Value()
	}

	// there are cases whey Key is called even if Next() is false
	// (see meterDB for instance)
	return nil
}

func (i *Iterator) Release() {
	i.it.Close()
}
