package db

// IndexDB is a database that allows prefixing by index.
type IndexDB interface {
	Has(index uint64, key []byte) (bool, error)
	Set(index uint64, key []byte, value []byte) error
}
