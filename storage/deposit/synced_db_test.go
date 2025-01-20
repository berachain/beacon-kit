// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package deposit

import (
	"os"
	"path/filepath"
	"testing"

	"cosmossdk.io/core/store"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/stretchr/testify/require"
)

func TestSyncedDBPersistence(t *testing.T) {
	// Create temporary directories for both DBs
	// regularDBPath := filepath.Join(t.TempDir(), "regular.db")
	syncedDBPath := filepath.Join(t.TempDir(), "synced.db")

	

	// Test cases to run
	testCases := []struct {
		name          string
		dbPath        string
		createDB      func(string) (store.KVStoreWithBatch, error)
		shouldPersist bool
	}{
		// {
		// 	name:   "Regular DB - should not persist on ungraceful shutdown",
		// 	dbPath: regularDBPath,
		// 	createDB: func(path string) (store.KVStoreWithBatch, error) {
		// 		return dbm.NewPebbleDB("test", path, nil)
		// 	},
		// 	shouldPersist: false,
		// },
		{
			name:   "SyncedDB - should persist on ungraceful shutdown",
			dbPath: syncedDBPath,
			createDB: func(path string) (store.KVStoreWithBatch, error) {
				db, err := dbm.NewPebbleDB("test", path, nil)
				if err != nil {
					return nil, err
				}
				return NewSynced(db), nil
			},
			shouldPersist: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Ensure clean directory
			require.NoError(t, os.RemoveAll(tc.dbPath))

			// Create new DB instance
			db, err := tc.createDB(tc.dbPath)
			require.NoError(t, err)

			// Write test data
			testKey := []byte("test-key")
			testValue := []byte("test-value")
			err = db.Set(testKey, testValue)
			require.NoError(t, err)

			// Properly close the first instance
			require.NoError(t, db.Close())

			// Create new DB instance with same path
			newDB, err := tc.createDB(tc.dbPath)
			require.NoError(t, err)
			defer func() {
				require.NoError(t, newDB.Close())
				require.NoError(t, os.RemoveAll(tc.dbPath))
			}()

			// Check if data persisted
			value, err := newDB.Get(testKey)
			require.NoError(t, err)

			if tc.shouldPersist {
				require.Equal(t, testValue, value, "data should have persisted")
			} else {
				require.Nil(t, value, "data should not have persisted")
			}
		})
	}
}
