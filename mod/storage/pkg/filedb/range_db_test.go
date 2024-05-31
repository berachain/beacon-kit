// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package filedb_test

import (
	"testing"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/errors"
	file "github.com/berachain/beacon-kit/mod/storage/pkg/filedb"
	"github.com/berachain/beacon-kit/mod/storage/pkg/interfaces/mocks"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// =========================== BASIC OPERATIONS ============================

func TestRangeDB(t *testing.T) {
	tests := []struct {
		name          string
		setupFunc     func(rdb *file.RangeDB) error
		testFunc      func(t *testing.T, rdb *file.RangeDB)
		expectedError bool
	}{
		{
			name: "Get",
			setupFunc: func(rdb *file.RangeDB) error {
				return rdb.Set(1, []byte("testKey"), []byte("testValue"))
			},
			testFunc: func(t *testing.T, rdb *file.RangeDB) {
				t.Helper()
				gotValue, err := rdb.Get(1, []byte("testKey"))
				require.NoError(t, err)
				require.Equal(t, []byte("testValue"), gotValue)
			},
		},
		{
			name: "Has",
			setupFunc: func(rdb *file.RangeDB) error {
				return rdb.Set(1, []byte("testKey"), []byte("testValue"))
			},
			testFunc: func(t *testing.T, rdb *file.RangeDB) {
				t.Helper()
				exists, err := rdb.Has(1, []byte("testKey"))
				require.NoError(t, err)
				require.True(t, exists)
			},
		},
		{
			name: "Set",
			setupFunc: func(_ *file.RangeDB) error {
				return nil // No setup required
			},
			testFunc: func(t *testing.T, rdb *file.RangeDB) {
				t.Helper()
				err := rdb.Set(1, []byte("testKey"), []byte("testValue"))
				require.NoError(t, err)

				exists, err := rdb.Has(1, []byte("testKey"))
				require.NoError(t, err)
				require.True(t, exists)
			},
		},
		{
			name: "Delete",
			setupFunc: func(rdb *file.RangeDB) error {
				return rdb.Set(1, []byte("testKey"), []byte("testValue"))
			},
			testFunc: func(t *testing.T, rdb *file.RangeDB) {
				t.Helper()
				err := rdb.Delete(1, []byte("testKey"))
				require.NoError(t, err)

				exists, err := rdb.Has(1, []byte("testKey"))
				require.NoError(t, err)
				require.False(t, exists)
			},
		},
		{
			name: "DeleteRange",
			setupFunc: func(rdb *file.RangeDB) error {
				for index := uint64(1); index <= 5; index++ {
					if err := rdb.Set(
						index, []byte("testKey"), []byte("testValue"),
					); err != nil {
						return err
					}
				}
				return nil
			},
			testFunc: func(t *testing.T, rdb *file.RangeDB) {
				t.Helper()
				err := rdb.DeleteRange(1, 4)
				require.NoError(t, err)

				for index := uint64(1); index <= 3; index++ {
					var exists bool
					exists, err = rdb.Has(index, []byte("testKey"))
					require.NoError(t, err)
					require.False(t, exists)
				}

				for index := uint64(4); index <= 5; index++ {
					var exists bool
					exists, err = rdb.Has(index, []byte("testKey"))
					require.NoError(t, err)
					require.True(t, exists)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rdb := file.NewRangeDB(newTestFDB(), 1)

			if tt.setupFunc != nil {
				if err := tt.setupFunc(rdb); (err != nil) != tt.expectedError {
					t.Fatalf(
						"setupFunc() error = %v, expectedError %v",
						err,
						tt.expectedError,
					)
				}
			}

			if tt.testFunc != nil {
				tt.testFunc(t, rdb)
			}
		})
	}
}

func TestExtractIndex(t *testing.T) {
	tests := []struct {
		name        string
		prefixedKey []byte
		expectedIdx uint64
		expectedErr error
	}{
		{
			name:        "ValidKey",
			prefixedKey: []byte("12345/testKey"),
			expectedIdx: 12345,
			expectedErr: nil,
		},
		{
			name:        "InvalidKeyFormat",
			prefixedKey: []byte("testKey"),
			expectedIdx: 0,
			expectedErr: errors.New("invalid key format"),
		},
		{
			name:        "InvalidIndex",
			prefixedKey: []byte("abc/testKey"),
			expectedIdx: 0,
			expectedErr: errors.New(
				"invalid index: strconv.ParseUint: parsing \"abc\": invalid syntax",
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()
			idx, err := file.ExtractIndex(tt.prefixedKey)
			require.Equal(t, tt.expectedIdx, idx)
			if tt.expectedErr != nil {
				require.ErrorContains(t, err, tt.expectedErr.Error())
			}
		})
	}
}

// =========================== PRUNING =====================================

func TestRangeDB_DeleteRange_NotSupported(t *testing.T) {
	tests := []struct {
		name string
		db   *mocks.DB
	}{
		{
			name: "DeleteRangeNotSupported",
			db:   new(mocks.DB),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()
			tt.db.On("DeleteRange", mock.Anything, mock.Anything).
				Return(errors.New("rangedb: delete range not supported for this db"))

			rdb := file.NewRangeDB(tt.db, 1)

			err := rdb.DeleteRange(1, 4)
			require.Error(t, err)
			require.Equal(t,
				"rangedb: delete range not supported for this db",
				err.Error())
		})
	}
}

func TestRangeDB_Prune(t *testing.T) {
	tests := []struct {
		name          string
		dataWindow    uint64
		setupFunc     func(rdb *file.RangeDB) error
		pruneIndex    uint64
		expectedError bool
		testFunc      func(t *testing.T, rdb *file.RangeDB)
	}{
		{
			name:       "PruneNoOp",
			dataWindow: 5,
			setupFunc: func(rdb *file.RangeDB) error {
				if err := populateTestDB(rdb, 1, 5); err != nil {
					return err
				}
				return nil
			},
			pruneIndex:    3, // index less than prune window, should be no-op
			expectedError: false,
			testFunc: func(t *testing.T, rdb *file.RangeDB) {
				t.Helper()
				requireExist(t, rdb, 1, 5)
			},
		},
		{
			name:       "PruneWithDeleteRange",
			dataWindow: 2,
			setupFunc: func(rdb *file.RangeDB) error {
				if err := populateTestDB(rdb, 1, 5); err != nil {
					return err
				}
				return nil
			},
			pruneIndex:    6, // index 6 with window 2 means delete from 0 to 3
			expectedError: false,
			testFunc: func(t *testing.T, rdb *file.RangeDB) {
				t.Helper()
				requireNotExist(t, rdb, 0, 3)
				requireExist(t, rdb, 4, 5)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rdb := file.NewRangeDB(newTestFDB(), tt.dataWindow)

			if tt.setupFunc != nil {
				if err := tt.setupFunc(rdb); (err != nil) != tt.expectedError {
					t.Fatalf(
						"setupFunc() error = %v, expectedError %v",
						err,
						tt.expectedError,
					)
				}
			}

			err := rdb.Prune(tt.pruneIndex)
			if (err != nil) != tt.expectedError {
				t.Fatalf(
					"Prune() error = %v, expectedError %v",
					err,
					tt.expectedError,
				)
			}

			if tt.testFunc != nil {
				tt.testFunc(t, rdb)
			}
		})
	}
}

// =========================== INVARIANTS ================================
// invariant: all indexes up to the firstNonNilIndex should be nil.
//
//nolint:gocognit //23
func TestRangeDB_Invariants(t *testing.T) {
	// we ignore errors for most of the tests below because we want to ensure
	// that the invariants hold in exceptional circumstances.
	tests := []struct {
		name      string
		setupFunc func(rdb *file.RangeDB) error
		testFunc  func(t *testing.T, rdb *file.RangeDB)
	}{
		{
			name: "Populate from empty",
			setupFunc: func(rdb *file.RangeDB) error {
				if err := populateTestDB(rdb, 1, 5); err != nil {
					return err
				}
				return nil
			},
			testFunc: func(t *testing.T, rdb *file.RangeDB) {
				t.Helper()
				requireNotExist(t, rdb, 0, min(rdb.FirstNonNilIndex()-1, 0))
			},
		},
		{
			name: "Delete from populated",
			setupFunc: func(rdb *file.RangeDB) error {
				if err := populateTestDB(rdb, 1, 5); err != nil {
					return err
				}
				return nil
			},
			testFunc: func(t *testing.T, rdb *file.RangeDB) {
				t.Helper()
				_ = rdb.Delete(2, []byte("key"))
				requireNotExist(t, rdb, 0, min(rdb.FirstNonNilIndex()-1, 0))
			},
		},
		{
			name: "Prune from populated",
			setupFunc: func(rdb *file.RangeDB) error {
				if err := populateTestDB(rdb, 1, 10); err != nil {
					return err
				}
				return nil
			},
			testFunc: func(t *testing.T, rdb *file.RangeDB) {
				t.Helper()
				_ = rdb.Prune(3)
				requireNotExist(t, rdb, 0, min(rdb.FirstNonNilIndex()-1, 0))
			},
		},
		{
			name: "DeleteRange from populated",
			setupFunc: func(rdb *file.RangeDB) error {
				if err := populateTestDB(rdb, 1, 10); err != nil {
					return err
				}
				return nil
			},
			testFunc: func(t *testing.T, rdb *file.RangeDB) {
				t.Helper()
				_ = rdb.DeleteRange(1, 5) // ignore error
				requireNotExist(t, rdb, 0, min(rdb.FirstNonNilIndex()-1, 0))
			},
		},
		{
			name: "Populate, Prune, Set round trip",
			setupFunc: func(rdb *file.RangeDB) error {
				if err := populateTestDB(rdb, 1, 30); err != nil {
					return err
				}
				return nil
			},
			testFunc: func(t *testing.T, rdb *file.RangeDB) {
				t.Helper()
				if err := rdb.Prune(25); err != nil {
					t.Fatalf("Prune() error = %v", err)
				}
				_ = populateTestDB(rdb, 5, 10)
				requireNotExist(t, rdb, 0, min(rdb.FirstNonNilIndex()-1, 0))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rdb := file.NewRangeDB(newTestFDB(), 2)

			if tt.setupFunc != nil {
				if err := tt.setupFunc(rdb); err != nil {
					// enforce invariant integrity on error
					requireNotExist(t, rdb, 0, min(rdb.FirstNonNilIndex()-1, 0))
				}
			}
			if tt.testFunc != nil {
				tt.testFunc(t, rdb)
			}
		})
	}
}

// MockDB that errors after a certain number of calls.
type ErrMockDB struct {
	*file.DB
	callsUntilError int
}

func NewErrMockDB(calls int) *ErrMockDB {
	return &ErrMockDB{
		DB:              newTestFDB(),
		callsUntilError: calls,
	}
}

func (m *ErrMockDB) RemoveAll(_ string) error {
	if m.callsUntilError == 0 {
		return errors.New("mocked error")
	}
	m.callsUntilError--
	return nil
}

// explicit enforcement of the RangeDB invariant on Pruning error.
func TestRangeDB_Invariant_Err(t *testing.T) {
	rdb := file.NewRangeDB(NewErrMockDB(2), 2)
	// populate the DB
	_ = populateTestDB(rdb, 1, 5)
	// first call to prune will move the firstNonNilIndex to 3
	_ = rdb.Prune(3)
	// second call to prune will err
	_ = rdb.Prune(10)
	// enforce invariant
	requireNotExist(t, rdb, 0, min(rdb.FirstNonNilIndex()-1, 0))
}

// =============================== HELPERS ==================================

// newTestFDB returns a new file DB instance with an in-memory filesystem.
func newTestFDB() *file.DB {
	fs := afero.NewMemMapFs()
	return file.NewDB(
		file.WithRootDirectory("/tmp/testdb"),
		file.WithFileExtension("txt"),
		file.WithDirectoryPermissions(0700),
		file.WithLogger(log.NewNopLogger()),
		file.WithAferoFS(fs),
	)
}

// requireNotExist requires the indexes from `from` to `to` to be empty.
//
//nolint:unparam //lol
func requireNotExist(t *testing.T, rdb *file.RangeDB, from uint64, to uint64) {
	t.Helper()
	for i := from; i <= to; i++ {
		exists, err := rdb.Has(i, []byte("key"))
		require.NoError(t, err)
		require.False(t, exists, "Index %d should have been pruned", i)
	}
}

// requireExist requires the indexes from `from` to `to` not be empty.
func requireExist(t *testing.T, rdb *file.RangeDB, from uint64, to uint64) {
	t.Helper()
	for i := from; i <= to; i++ {
		exists, err := rdb.Has(i, []byte("key"))
		require.NoError(t, err)
		require.True(t, exists, "Index %d should not have been pruned", i)
	}
}

// populateTestDB populates the test DB with indexes from `from` to `to`.
func populateTestDB(rdb *file.RangeDB, from, to uint64) error {
	for i := from; i <= to; i++ {
		if err := rdb.Set(i, []byte("key"), []byte("value")); err != nil {
			return err
		}
	}
	return nil
}
