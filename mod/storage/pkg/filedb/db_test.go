// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
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

package filedb_test

import (
	"testing"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/log/pkg/noop"
	file "github.com/berachain/beacon-kit/mod/storage/pkg/filedb"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestDB(t *testing.T) {
	tests := []struct {
		name          string
		setupFunc     func(db *file.DB) error
		testFunc      func(t *testing.T, db *file.DB)
		expectedError bool
	}{
		{
			name: "NewDB",
			testFunc: func(t *testing.T, db *file.DB) {
				t.Helper()
				require.NotNil(t, db)
			},
		},
		{
			name: "SetAndGet",
			setupFunc: func(db *file.DB) error {
				return db.Set([]byte("key"), []byte("value"))
			},
			testFunc: func(t *testing.T, db *file.DB) {
				t.Helper()
				retrievedValue, err := db.Get([]byte("key"))
				require.NoError(t, err)
				require.Equal(t, []byte("value"), retrievedValue)
			},
		},
		{
			name: "Has",
			setupFunc: func(db *file.DB) error {
				return db.Set([]byte("key"), []byte("value"))
			},
			testFunc: func(t *testing.T, db *file.DB) {
				t.Helper()
				exists, err := db.Has([]byte("key"))
				require.NoError(t, err)
				require.True(t, exists)
			},
		},
		{
			name: "Delete",
			setupFunc: func(db *file.DB) error {
				return db.Set([]byte("key"), []byte("value"))
			},
			testFunc: func(t *testing.T, db *file.DB) {
				t.Helper()
				exists, err := db.Has([]byte("key"))
				require.NoError(t, err)
				require.True(t, exists)

				err = db.Delete([]byte("key"))
				require.NoError(t, err)

				exists, err = db.Has([]byte("key"))
				require.NoError(t, err)
				require.False(t, exists)
			},
		},
		{
			name: "SetExistingKey",
			setupFunc: func(db *file.DB) error {
				if err := db.Set([]byte("key"), []byte("value1")); err != nil {
					return err
				}
				return db.Set([]byte("key"), []byte("value2"))
			},
			testFunc: func(t *testing.T, db *file.DB) {
				t.Helper()
				retrievedValue, err := db.Get([]byte("key"))
				require.NoError(t, err)
				require.Equal(t, []byte("value2"), retrievedValue)
			},
		},
		{
			name: "GetNonExistingKey",
			testFunc: func(t *testing.T, db *file.DB) {
				t.Helper()
				_, err := db.Get([]byte("non-existing"))
				require.Error(t, err)
			},
			expectedError: true,
		},
		// If the key does not exist, `Has` will return false with error as nil
		{
			name: "HasNonExistingKey",
			testFunc: func(t *testing.T, db *file.DB) {
				t.Helper()
				exists, err := db.Has([]byte("non-existing"))
				require.NoError(t, err)
				require.False(t, exists)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			db := file.NewDB(
				file.WithRootDirectory("/tmp/testdb"),
				file.WithFileExtension("txt"),
				file.WithDirectoryPermissions(0700),
				file.WithLogger(noop.NewLogger[log.Logger]()),
				file.WithAferoFS(fs),
			)

			if tt.setupFunc != nil {
				if err := tt.setupFunc(db); (err != nil) != tt.expectedError {
					t.Fatalf(
						"setupFunc() error = %v, expectedError %v",
						err,
						tt.expectedError,
					)
				}
			}

			if tt.testFunc != nil {
				tt.testFunc(t, db)
			}
		})
	}

	t.Run("NewDBWithInvalidOption", func(t *testing.T) {
		invalidOption := func(_ *file.DB) error {
			return errors.New("invalid option")
		}
		require.Panics(t, func() { file.NewDB(invalidOption) })
	})
}

// Test with `etc` as root directory to cause creation failure
// due to permission denied.
func TestDB_SetExistingKey_CreateError(t *testing.T) {
	test := struct {
		name          string
		setupFunc     func(db *file.DB) error
		testFunc      func(t *testing.T, db *file.DB)
		expectedError bool
	}{
		name: "SetExistingKeyWithCreateError",
		testFunc: func(t *testing.T, db *file.DB) {
			t.Helper()
			err := db.Set([]byte("key"), []byte("value"))
			require.Error(t, err)
			require.ErrorContains(t, err, "failed to create file")
		},
		expectedError: true,
	}

	t.Run(test.name, func(t *testing.T) {
		fs := afero.NewMemMapFs()
		db := file.NewDB(
			file.WithRootDirectory("/etc"),
			file.WithFileExtension("txt"),
			file.WithDirectoryPermissions(0700),
			file.WithLogger(noop.NewLogger[log.Logger]()),
			file.WithAferoFS(fs),
		)

		if test.setupFunc != nil {
			if err := test.setupFunc(db); (err != nil) != test.expectedError {
				require.Error(t, err, "setupFunc() error = %v", err)
			}
		}

		if test.testFunc != nil {
			test.testFunc(t, db)
		}
	})
}

// Test with root directory as a file.
func TestDB_SetHas_NotDirError(t *testing.T) {
	tests := []struct {
		name          string
		setupFunc     func(db *file.DB) error
		testFunc      func(t *testing.T, db *file.DB)
		expectedError bool
	}{
		{
			name: "HasWithError",
			testFunc: func(t *testing.T, db *file.DB) {
				t.Helper()
				value, err := db.Has([]byte("key"))
				require.Error(t, err)
				require.False(t, value)
				require.ErrorContains(t, err, "not a directory")
			},
			expectedError: true,
		},
		{
			name: "SetWithError",
			testFunc: func(t *testing.T, db *file.DB) {
				t.Helper()
				err := db.Set([]byte("key"), []byte("value"))
				require.Error(t, err)
				require.ErrorContains(t, err, "not a directory")
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			db := file.NewDB(
				file.WithRootDirectory("/etc/passwd"),
				file.WithFileExtension("txt"),
				file.WithDirectoryPermissions(0700),
				file.WithLogger(noop.NewLogger[log.Logger]()),
				file.WithAferoFS(fs),
			)

			if tt.setupFunc != nil {
				if err := tt.setupFunc(db); (err != nil) != tt.expectedError {
					require.Error(t, err, "setupFunc() error = %v", err)
				}
			}

			if tt.testFunc != nil {
				tt.testFunc(t, db)
			}
		})
	}
}

// Test with root directory to be created in `etc`
// which will result in permission denied.
func TestDB_Set_MkDirError(t *testing.T) {
	test := struct {
		name          string
		setupFunc     func(db *file.DB) error
		testFunc      func(t *testing.T, db *file.DB)
		expectedError bool
	}{
		name: "SetWithMkdirAllError",
		testFunc: func(t *testing.T, db *file.DB) {
			t.Helper()
			err := db.Set([]byte("key"), []byte("value"))
			require.Error(t, err)
		},
		expectedError: true,
	}

	t.Run(test.name, func(t *testing.T) {
		fs := afero.NewMemMapFs()
		db := file.NewDB(
			file.WithRootDirectory("/etc/test"),
			file.WithFileExtension("txt"),
			file.WithDirectoryPermissions(0700),
			file.WithLogger(noop.NewLogger[log.Logger]()),
			file.WithAferoFS(fs),
		)

		if test.setupFunc != nil {
			if err := test.setupFunc(db); (err != nil) != test.expectedError {
				require.Error(t, err, "setupFunc() error = %v", err)
			}
		}

		if test.testFunc != nil {
			test.testFunc(t, db)
		}
	})
}
