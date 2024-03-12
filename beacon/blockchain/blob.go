package blockchain

import (
	"bytes"
	"sync"

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

// func sidecarFileKey(sidecar *eth.BlobSidecar) string {
// 	return path.Join(bs.baseDir, fmt.Sprintf(
// 		"%d_%x_%d_%x.blob",
// 		sidecar.Slot,
// 		sidecar.BlockRoot,
// 		sidecar.Index,
// 		sidecar.KzgCommitment,
// 	))
// }

// Store the blobs in the blobstore.
func PrepareBlobsHandler(storage db.DB,
	height int64, commitments [][48]byte,
	blobs *engine.BlobsBundleV1) ([][]byte, error) {

	// TODO: before storage handle validation

	// store the blobs under a single height.
	ranger := file.NewRangeDB[int64](storage)
	for i, sidecar := range blobs.Blobs {
		if err := ranger.Set(height, blobs.Commitments[i], sidecar); err != nil {
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

	ranger := file.NewRangeDB[int64](storage)
	// Store the blobs under a single height.
	for i, sidecar := range blobs {
		if err := ranger.Set(height, commitments[i][:], sidecar); err != nil {
			return err
		}
	}

	return nil
}
