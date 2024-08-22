package db

import (
	"path/filepath"

	dbm "github.com/cosmos/cosmos-db"
)

// OpenDB opens the application database using the appropriate driver.
func OpenDB(rootDir string, backendType dbm.BackendType) (dbm.DB, error) {
	dataDir := filepath.Join(rootDir, "data")
	return dbm.NewDB("application", backendType, dataDir)
}
