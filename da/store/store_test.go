//go:build race
// +build race

package store

import (
	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/config/spec"
	"github.com/berachain/beacon-kit/consensus-types/types"

	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/storage/filedb"
	"os"
	"testing"
)

func TestStore_PersistRace(t *testing.T) {
	tmpFilePath := "/tmp/store_test"

	// Make sure we start fresh
	err := os.RemoveAll(tmpFilePath)
	if err != nil {
		t.Fatal(err)
	}

	// Remove db when we're done
	defer os.RemoveAll(tmpFilePath)

	logger := log.NewNopLogger()
	chainSpec, err := spec.DevnetChainSpec()
	if err != nil {
		t.Fatal(err)
	}
	s := New[*types.BeaconBlockBody](
		filedb.NewRangeDB(
			filedb.NewDB(filedb.WithRootDirectory(tmpFilePath),
				filedb.WithFileExtension("ssz"),
				filedb.WithDirectoryPermissions(0700),
				filedb.WithLogger(logger),
			),
		),
		logger.With("service", "da-store"),
		chainSpec,
	)

	// This many blobs is not currently possible, but it doesn't hurt eh
	sc := make([]*datypes.BlobSidecar, 20)
	for i := range sc {
		sc[i] = &datypes.BlobSidecar{
			Index:             uint64(i),
			BeaconBlockHeader: &types.BeaconBlockHeader{},
		}
	}
	sidecars := datypes.BlobSidecars{
		Sidecars: sc,
	}

	// Multiple writes to DB
	err = s.Persist(0, &sidecars)
	err = s.Persist(1, &sidecars)
	err = s.Prune(0, 1)
	err = s.Persist(0, &sidecars)
}
