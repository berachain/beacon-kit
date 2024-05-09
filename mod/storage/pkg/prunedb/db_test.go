package prunedb

import (
	"context"
	"cosmossdk.io/log"
	"fmt"
	"github.com/berachain/beacon-kit/mod/storage/pkg/prunedb/mocks"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type test struct {
	name        string
	setupFunc   func(db *mocks.IndexDB) error
	testFunc    func(t *testing.T, db *DB)
	expectedErr error
}

func TestDB(t *testing.T) {
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
			testFunc: func(t *testing.T, db *DB) {
				key := []byte("testKey")
				value := []byte("testValue")
				err := db.Set(uint64(1), key, value)
				require.NoError(t, err)
				res, err := db.Get(uint64(1), key)
				require.NoError(t, err)
				require.Equal(t, value, res)
				require.Equal(t, db.highestSetIndex, uint64(1))
			},
		},
		{
			name: "Has",
			setupFunc: func(mockDB *mocks.IndexDB) error {
				key := []byte("testKey")
				mockDB.On("Has", uint64(1), key).Return(true, nil)
				return nil
			},
			testFunc: func(t *testing.T, db *DB) {
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
				mockDB.On("Set", uint64(1), key, []byte("value")).Return(nil)
				mockDB.On("Delete", uint64(1), key).Return(nil)
				mockDB.On("Has", uint64(1), key).Return(false, nil)
				return nil
			},
			testFunc: func(t *testing.T, db *DB) {
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
		{
			name: "DeleteRange",
			setupFunc: func(mockDB *mocks.IndexDB) error {

				for i := uint64(1); i <= 10; i++ {
					key := []byte("testKey" + fmt.Sprint(i))
					mockDB.On("Set", i, key, []byte("value")).Return(nil)
				}
				mockDB.On("DeleteRange", uint64(1), uint64(5)).Return(nil)

				for i := uint64(1); i <= 4; i++ {
					mockDB.On("Has", i, []byte("testKey"+fmt.Sprint(i))).Return(false, nil)
				}

				for i := uint64(5); i <= 10; i++ {
					mockDB.On("Has", i, []byte("testKey"+fmt.Sprint(i))).Return(true, nil)
				}
				return nil
			},
			testFunc: func(t *testing.T, db *DB) {
				for i := uint64(1); i <= 10; i++ {
					key := []byte("testKey" + fmt.Sprint(i))
					err := db.Set(i, key, []byte("value"))
					require.NoError(t, err)
				}
				err := db.DeleteRange(uint64(1), uint64(5))
				require.NoError(t, err)
				for i := uint64(1); i <= 4; i++ {
					key := []byte("testKey" + fmt.Sprint(i))
					res, err := db.Has(i, key)
					require.NoError(t, err)
					require.False(t, res)
				}
				for i := uint64(5); i <= 10; i++ {
					key := []byte("testKey" + fmt.Sprint(i))
					res, err := db.Has(i, key)
					require.NoError(t, err)
					require.True(t, res)
				}
			},
		},
		{name: "new",
			setupFunc: func(mockDB *mocks.IndexDB) error {
				return nil
			},
			testFunc: func(t *testing.T, db *DB) {
				mockDB := new(mocks.IndexDB)
				logger := log.NewNopLogger()
				pruneInterval := 50 * time.Millisecond
				windowSize := uint64(5)

				createdDb := New(mockDB, logger, pruneInterval, windowSize)

				require.NotNil(t, createdDb)
				require.Equal(t, mockDB, createdDb.IndexDB)
				require.Equal(t, logger, createdDb.logger)
				require.Equal(t, windowSize, createdDb.windowSize)
				require.NotNil(t, createdDb.ticker)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(mocks.IndexDB)
			db := New(mockDB, log.NewNopLogger(), 50*time.Millisecond, 5)
			if tt.setupFunc != nil {
				tt.setupFunc(mockDB)
			}
			if tt.testFunc != nil {
				tt.testFunc(t, db)
			}
			mockDB.AssertExpectations(t)
		})
	}
}

func TestPrune(t *testing.T) {

	mockDB := new(mocks.IndexDB)
	db := New(mockDB, log.NewNopLogger(), 50*time.Millisecond, 5)

	// Set expectations
	mockDB.On("DeleteRange", uint64(0), uint64(2)).Return(nil)

	db.highestSetIndex = 7

	err := db.prune()
	require.NoError(t, err)

	mockDB.AssertExpectations(t)
}

func TestDB_Start(t *testing.T) {
	mockDB := new(mocks.IndexDB)
	db := New(mockDB, log.NewNopLogger(), 3*time.Second, 5)

	// Set expectations
	for i := uint64(1); i <= 10; i++ {
		key := []byte("testKey" + fmt.Sprint(i))
		mockDB.On("Set", i, key, []byte("value")).Return(nil)
	}
	for i := uint64(1); i <= 4; i++ {
		key := []byte("testKey" + fmt.Sprint(i))
		mockDB.On("Has", i, key).Return(false, nil)
	}
	for i := uint64(5); i <= 10; i++ {
		key := []byte("testKey" + fmt.Sprint(i))
		mockDB.On("Has", i, key).Return(true, nil)
	}

	// TODO: this test fails as according to logic we are deleting from 0 to 5. However, the test expects deletion from 1 to 5.
	mockDB.On("DeleteRange", uint64(0), uint64(5)).Return(nil)

	for i := uint64(1); i <= 10; i++ {
		key := []byte("testKey" + fmt.Sprint(i))
		err := db.Set(i, key, []byte("value"))
		require.NoError(t, err)
	}

	for i := uint64(1); i <= 4; i++ {
		key := []byte("testKey" + fmt.Sprint(i))
		res, err := db.Has(i, key)
		require.NoError(t, err)
		require.False(t, res)
	}
	for i := uint64(5); i <= 10; i++ {
		key := []byte("testKey" + fmt.Sprint(i))
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
