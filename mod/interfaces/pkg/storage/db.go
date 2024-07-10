package db

// DB is the interface for a simple key-value store.
type DB interface {
	Get(key []byte) ([]byte, error)
	Has(key []byte) (bool, error)
	Set(key []byte, value []byte) error
	Delete(key []byte) error

	// TODO: add Batch and full DB stuff.
}
