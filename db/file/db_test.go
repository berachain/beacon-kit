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

package file_test

import (
	"testing"

	"cosmossdk.io/log"
	"github.com/itsdevbear/bolaris/db/file"
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
				require.NotNil(t, db)
			},
		},
		{
			name: "SetAndGet",
			setupFunc: func(db *file.DB) error {
				return db.Set([]byte("key"), []byte("value"))
			},
			testFunc: func(t *testing.T, db *file.DB) {
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
				retrievedValue, err := db.Get([]byte("key"))
				require.NoError(t, err)
				require.Equal(t, []byte("value2"), retrievedValue)
			},
		},
		{
			name: "GetNonExistingKey",
			testFunc: func(t *testing.T, db *file.DB) {
				_, err := db.Get([]byte("non-existing"))
				require.Error(t, err)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			db := file.NewDB(
				file.WithRootDirectory("/tmp/testdb"),
				file.WithFileExtension("txt"),
				file.WithDirectoryPermissions(0700),
				file.WithLogger(log.NewNopLogger()),
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
}
