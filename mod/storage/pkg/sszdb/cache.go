package sszdb

// A layer over SchemaDB mirroring the cache layer in the Cosmos SDK KVStore.
type CacheDB struct {
	*SchemaDB
	cache map[uint64][]byte
}

func NewCacheDB(schemaDB *SchemaDB) *CacheDB {
	return &CacheDB{
		SchemaDB: schemaDB,
		cache:    make(map[uint64][]byte),
	}
}
