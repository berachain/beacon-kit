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
	"context"
	"testing"
	"time"

	"cosmossdk.io/log"
	file "github.com/berachain/beacon-kit/mod/storage/pkg/filedb"
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
				t.Helper()
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
