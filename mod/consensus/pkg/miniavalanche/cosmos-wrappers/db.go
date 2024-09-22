package cosmoswrappers

import (
	"bytes"
	"errors"

	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/avalanchego/database/prefixdb"
	dbm "github.com/cosmos/cosmos-db"
)

var (
	_ dbm.DB       = (*AvaDBWrap)(nil)
	_ dbm.Batch    = (*AvaBatchWrap)(nil)
	_ dbm.Iterator = (*AvaIteratorWrap)(nil)

	errDBFeatureNotImplementedYet = errors.New("db feature not implemented yet")
)

type AvaDBWrap struct {
	db *prefixdb.Database
}

func NewAvaDBWrapper(prefix []byte, db database.Database) *AvaDBWrap {
	return &AvaDBWrap{
		db: prefixdb.New(prefix, db),
	}
}

func (adbw *AvaDBWrap) Get(key []byte) ([]byte, error) {
	res, err := adbw.db.Get(key)
	if errors.Is(err, database.ErrNotFound) {
		return nil, nil
	}
	return res, err
}

func (adbw *AvaDBWrap) Has(key []byte) (bool, error) {
	return adbw.db.Has(key)
}

func (adbw *AvaDBWrap) Set(key []byte, value []byte) error {
	return adbw.db.Put(key, value)
}

func (adbw *AvaDBWrap) SetSync(key []byte, value []byte) error {
	return adbw.db.Put(key, value) // TODO: make sure this is the right mapping
}

func (adbw *AvaDBWrap) Delete(key []byte) error {
	return adbw.db.Delete(key)
}

func (adbw *AvaDBWrap) DeleteSync(key []byte) error {
	return adbw.db.Delete(key) // TODO: make sure this is the right mapping
}

func (adbw *AvaDBWrap) Iterator(start, end []byte) (dbm.Iterator, error) {
	return &AvaIteratorWrap{
		start:   start,
		end:     end,
		hasNext: nil,
		it:      adbw.db.NewIteratorWithStart(start),
	}, nil
}

func (adbw *AvaDBWrap) ReverseIterator([]byte, []byte) (dbm.Iterator, error) {
	return nil, errDBFeatureNotImplementedYet
}

func (adbw *AvaDBWrap) Close() error {
	return adbw.db.Close()
}

func (adbw *AvaDBWrap) NewBatch() dbm.Batch {
	return &AvaBatchWrap{
		batch: adbw.db.NewBatch(),
	}
}

func (adbw *AvaDBWrap) NewBatchWithSize(int) dbm.Batch {
	return &AvaBatchWrap{
		batch: adbw.db.NewBatch(),
	} // TODO: handle size
}

func (adbw *AvaDBWrap) Print() error {
	return nil // TODO print it
}

func (adbw *AvaDBWrap) Stats() map[string]string {
	return map[string]string{} // TODO populate with relevant stats
}

type AvaBatchWrap struct {
	batch database.Batch
}

func (abw *AvaBatchWrap) Set(key, value []byte) error {
	return abw.batch.Put(key, value)
}

func (abw *AvaBatchWrap) Delete(key []byte) error {
	return abw.batch.Delete(key)
}

func (abw *AvaBatchWrap) Write() error {
	return abw.batch.Write()
}

func (abw *AvaBatchWrap) WriteSync() error {
	return abw.batch.Write() // TODO: make sure this is the right mapping
}

func (abw *AvaBatchWrap) Close() error {
	return nil // TODO: make sure this is the right mapping
}

func (abw *AvaBatchWrap) GetByteSize() (int, error) {
	return abw.batch.Size(), nil
}

type AvaIteratorWrap struct {
	start, end []byte
	hasNext    *bool

	it database.Iterator
}

func (aiw *AvaIteratorWrap) Domain() ([]byte, []byte) {
	return aiw.start, aiw.end
}

func (aiw *AvaIteratorWrap) Valid() bool {
	if aiw.hasNext == nil {
		aiw.hasNext = new(bool)
		*aiw.hasNext = aiw.it.Next()
	}
	if !*aiw.hasNext {
		return *aiw.hasNext
	}

	if bytes.Compare(aiw.it.Value(), aiw.end) == 1 {
		*aiw.hasNext = false // we're beyond end
	}
	return *aiw.hasNext
}

func (aiw *AvaIteratorWrap) Next() {
	*aiw.hasNext = aiw.it.Next()
	if bytes.Compare(aiw.it.Value(), aiw.end) == 1 {
		*aiw.hasNext = false // we're beyond end
	}
}

func (aiw *AvaIteratorWrap) Key() []byte {
	return aiw.it.Key()
}

func (aiw *AvaIteratorWrap) Value() []byte {
	return aiw.it.Value()
}

func (aiw *AvaIteratorWrap) Error() error {
	return aiw.it.Error()
}

func (aiw *AvaIteratorWrap) Close() error {
	aiw.it.Release()
	return nil
}
