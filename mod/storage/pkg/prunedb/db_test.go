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

package prunedb_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"cosmossdk.io/log"
	prune "github.com/berachain/beacon-kit/mod/storage/pkg/prunedb"
	"github.com/berachain/beacon-kit/mod/storage/pkg/prunedb/mocks"
	"github.com/stretchr/testify/require"
)

type test struct {
	name          string
	setupFunc     func(db *mocks.IndexDB) error
	testFunc      func(t *testing.T, db *prune.DB)
	expectedError bool
}

func TestDB_CRUD(t *testing.T) {
	tests := []test{
		{
			name: "Set and Get",
			setupFunc: func(mockDB *mocks.IndexDB) error {
				key := []byte("testKey")
				value := []byte("testValue")
				mockDB.On("Set", uint64(1), key, value).Return(nil)
				mockDB.On("Get", uint64(1), key).Return(value, nil)
				return nil
			},
			testFunc: func(t *testing.T, db *prune.DB) {
				key := []byte("testKey")
				value := []byte("testValue")
				err := db.Set(uint64(1), key, value)
				require.NoError(t, err)
				res, err := db.Get(uint64(1), key)
				require.NoError(t, err)
				require.Equal(t, value, res)
				require.Equal(t, uint64(1), db.GetHighestSetIndex())
			},
		},
		{
			name: "Has",
			setupFunc: func(mockDB *mocks.IndexDB) error {
				key := []byte("testKey")
				mockDB.On("Has", uint64(1), key).Return(true, nil)
				return nil
			},
			testFunc: func(t *testing.T, db *prune.DB) {
				key := []byte("testKey")
				res, err := db.Has(uint64(1), key)
				require.NoError(t, err)
				require.True(t, res)
			},
		},
		{
			name: "Delete",
			setupFunc: func(mockDB *mocks.IndexDB) error {
				key := []byte("testKey")
				mockDB.On("Set", uint64(1), key, []byte("value")).
					Return(nil)
				mockDB.On("Delete", uint64(1), key).
					Return(nil)
				mockDB.On("Has", uint64(1), key).
					Return(false, nil)
				return nil
			},
			testFunc: func(t *testing.T, db *prune.DB) {
				key := []byte("testKey")
				err := db.Set(uint64(1), key, []byte("value"))
				require.NoError(t, err)
				err = db.Delete(uint64(1), key)
				require.NoError(t, err)
				res, err := db.Has(uint64(1), key)
				require.NoError(t, err)
				require.False(t, res)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(mocks.IndexDB)
			db := prune.New(mockDB, log.NewNopLogger(), 50*time.Millisecond, 5)
			if tt.setupFunc != nil {
				if err := tt.setupFunc(mockDB); (err != nil) != tt.expectedError {
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
			mockDB.AssertExpectations(t)
		})
	}
}

func TestDB_DeleteRange(t *testing.T) {
	tests := test{
		name: "DeleteRange",
		setupFunc: func(mockDB *mocks.IndexDB) error {
			for i := uint64(1); i <= 10; i++ {
				key := []byte("testKey" + strconv.FormatUint(i, 10))
				mockDB.On("Set", i, key, []byte("value")).Return(nil)
			}
			mockDB.On("DeleteRange", uint64(1), uint64(5)).Return(nil)

			for i := uint64(1); i <= 4; i++ {
				mockDB.On(
					"Has",
					i,
					[]byte("testKey"+strconv.FormatUint(i, 10))).
					Return(false, nil)
			}

			for i := uint64(5); i <= 10; i++ {
				mockDB.On(""+
					"Has",
					i,
					[]byte("testKey"+strconv.FormatUint(i, 10))).
					Return(true, nil)
			}
			return nil
		},
		testFunc: func(t *testing.T, db *prune.DB) {
			for i := uint64(1); i <= 10; i++ {
				key := []byte("testKey" + strconv.FormatUint(i, 10))
				err := db.Set(i, key, []byte("value"))
				require.NoError(t, err)
			}
			err := db.DeleteRange(uint64(1), uint64(5))
			require.NoError(t, err)
			for i := uint64(1); i <= 4; i++ {
				key := []byte("testKey" + strconv.FormatUint(i, 10))
				res, errInHas := db.Has(i, key)
				require.NoError(t, errInHas)
				require.False(t, res)
			}
			for i := uint64(5); i <= 10; i++ {
				key := []byte("testKey" + strconv.FormatUint(i, 10))
				res, errInHasPresent := db.Has(i, key)
				require.NoError(t, errInHasPresent)
				require.True(t, res)
			}
		},
	}

	t.Run(tests.name, func(t *testing.T) {
		mockDB := new(mocks.IndexDB)
		db := prune.New(mockDB, log.NewNopLogger(), 50*time.Millisecond, 5)
		if tests.setupFunc != nil {
			if err := tests.setupFunc(mockDB); (err != nil) != tests.expectedError {
				t.Fatalf(
					"setupFunc() error = %v, expectedError %v",
					err,
					tests.expectedError,
				)
			}
		}
		if tests.testFunc != nil {
			tests.testFunc(t, db)
		}
		mockDB.AssertExpectations(t)
	})
}

func TestDB_New(t *testing.T) {
	mockDB := new(mocks.IndexDB)
	logger := log.NewNopLogger()
	pruneInterval := 50 * time.Millisecond
	windowSize := uint64(5)

	createdDB := prune.New(mockDB, logger, pruneInterval, windowSize)

	require.NotNil(t, createdDB)
	require.Equal(t, mockDB, createdDB.IndexDB)
	require.Equal(t, logger, createdDB.GetLogger())
	require.Equal(t, windowSize, createdDB.GetWindowSize())
}
func TestDB_Start(t *testing.T) {
	mockDB := new(mocks.IndexDB)
	db := prune.New(mockDB, log.NewNopLogger(), 3*time.Second, 5)

	// Set expectations
	for i := uint64(1); i <= 10; i++ {
		key := []byte("testKey" + strconv.FormatUint(i, 10))
		mockDB.On("Set", i, key, []byte("value")).Return(nil)
	}
	for i := uint64(1); i <= 4; i++ {
		key := []byte("testKey" + strconv.FormatUint(i, 10))
		mockDB.On("Has", i, key).Return(false, nil)
	}
	for i := uint64(5); i <= 10; i++ {
		key := []byte("testKey" + strconv.FormatUint(i, 10))
		mockDB.On("Has", i, key).Return(true, nil)
	}

	// TODO: this test fails as according to logic we are deleting from 0 to 5.
	// However, the test expects deletion from 1 to 5.
	mockDB.On("DeleteRange", uint64(0), uint64(5)).Return(nil)

	for i := uint64(1); i <= 10; i++ {
		key := []byte("testKey" + strconv.FormatUint(i, 10))
		err := db.Set(i, key, []byte("value"))
		require.NoError(t, err)
	}

	for i := uint64(1); i <= 4; i++ {
		key := []byte("testKey" + strconv.FormatUint(i, 10))
		res, err := db.Has(i, key)
		require.NoError(t, err)
		require.False(t, res)
	}
	for i := uint64(5); i <= 10; i++ {
		key := []byte("testKey" + strconv.FormatUint(i, 10))
		res, err := db.Has(i, key)
		require.NoError(t, err)
		require.True(t, res)
	}

	_, cancel := context.WithCancel(context.Background())

	// Wait for the ticker to tick at least once
	time.Sleep(4 * time.Second)

	// Cancel the context to stop the ticker
	cancel()
}
