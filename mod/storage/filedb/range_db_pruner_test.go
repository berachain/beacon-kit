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
	"context"
	"testing"
	"time"

	"cosmossdk.io/log"
	file "github.com/berachain/beacon-kit/mod/storage/filedb"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestRangeDBPruner(t *testing.T) {
	notifyCh := make(chan uint64)
	tests := []struct {
		name        string
		pruneWindow uint64
		setupFunc   func(
			rdb *file.RangeDB) error
		testFunc func(
			t *testing.T, rdb *file.RangeDB)
		expectedError bool
	}{
		{
			name:        "PruneOldIndexes",
			pruneWindow: 5,
			setupFunc: func(
				rdb *file.RangeDB,
			) error {
				for i := uint64(1); i <= 10; i++ {
					if err := rdb.Set(i, []byte("key"), []byte("value")); err != nil {
						return err
					}
					notifyCh <- i
				}
				return nil
			},
			testFunc: func(
				t *testing.T, rdb *file.RangeDB,
			) {
				time.Sleep(
					100 * time.Millisecond,
				) // Wait for pruner to catch up
				for i := uint64(1); i < 5; i++ {
					exists, err := rdb.Has(i, []byte("key"))
					require.NoError(t, err)
					require.False(
						t,
						exists,
						"Index %d should have been pruned",
						i,
					)
				}
				for i := uint64(5); i <= 10; i++ {
					exists, err := rdb.Has(i, []byte("key"))
					require.NoError(t, err)
					require.True(
						t,
						exists,
						"Index %d should not have been pruned",
						i,
					)
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
			pruner := file.NewRangeDBPruner(rdb, tt.pruneWindow, notifyCh)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			pruner.Start(ctx)

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
