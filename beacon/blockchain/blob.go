package blockchain

import (
	"bytes"
	"sync"

	beacontypes "github.com/berachain/beacon-kit/beacon/core/types"
	"github.com/berachain/beacon-kit/db"
	"github.com/berachain/beacon-kit/db/file"
	"github.com/ethereum/go-ethereum/beacon/engine"
)

// Create a pool of bytes.Buffers.
var bufPool = &sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// Store the blobs in the blobstore.
func PrepareBlobsHandler(storage db.DB,
	height int64, blk beacontypes.BeaconBlock,
	blobs *engine.BlobsBundleV1) ([][]byte, error) {

	// store the blobs under a single height.
	ranger := file.NewRangeDB(storage)
	for i, sidecar := range blobs.Blobs {
		if err := ranger.Set(uint64(height), blobs.Commitments[i], sidecar); err != nil {
			return nil, err
		}
	}

	var blobTx = make([][]byte, 0, len(blobs.Blobs))
	for i, sidecar := range blobs.Blobs {
		blobTx[i] = sidecar
	}

	return blobTx, nil
}

// Store the blobs in the blobstore.
func ProcessBlobsHandler(storage db.DB,
	height int64, commitments [][48]byte, blobs [][]byte) error {

	// TODO: before storage handle validation

	ranger := file.NewRangeDB(storage)
	// Store the blobs under a single height.
	for i, sidecar := range blobs {
		if err := ranger.Set(uint64(height), commitments[i][:], sidecar); err != nil {
			return err
		}
	}

	return nil
}
