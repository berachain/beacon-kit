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
			fs := afero.NewMemMapFs()
			db := file.NewDB(
				file.WithRootDirectory("/tmp/testdb"),
				file.WithFileExtension("txt"),
				file.WithDirectoryPermissions(0700),
				file.WithLogger(log.NewNopLogger()),
				file.WithAferoFS(fs),
			)
			rdb := file.NewRangeDB(db)

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

			rdb := file.NewRangeDB(tt.db)

			err := rdb.DeleteRange(1, 4)
			require.Error(t, err)
			require.Equal(t,
				"rangedb: delete range not supported for this db",
				err.Error())
		})
	}
}
