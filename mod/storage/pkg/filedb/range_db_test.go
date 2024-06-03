// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
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
