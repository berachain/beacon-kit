package blockchain

import (
	"bytes"
	"sync"

	"github.com/berachain/beacon-kit/beacon/core/types"
	beacontypes "github.com/berachain/beacon-kit/beacon/core/types"
	"github.com/berachain/beacon-kit/crypto/kzg"
	"github.com/berachain/beacon-kit/db"
	"github.com/berachain/beacon-kit/db/file"
	enginetypes "github.com/berachain/beacon-kit/engine/types"
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
	blobs *enginetypes.BlobsBundleV1) ([][]byte, error) {

	ranger := file.NewRangeDB(storage)
	var blobTx = make([]types.BlobTxSidecar, 0, len(blobs.Blobs))
	for i, sidecar := range blobs.Blobs {
		//Create Inclusion Proof
		ic, err := kzg.MerkleProofKZGCommitment(blk, i)
		if err != nil {
			return nil, err
		}
		blob := types.BlobTxSidecar{
			Blob:           sidecar,
			KzgCommitment:  blobs.Commitments[i],
			KzgProof:       blobs.Proofs[i],
			InclusionProof: ic,
		}

		if err := ranger.Set(uint64(height), blobs.Commitments[i], sidecar); err != nil {
			return nil, err
		}

		blobTx[i] = blob
	}

	//TODO: ssz encode the blobTx

	return blobTx, nil
}

// Store the blobs in the blobstore.
func ProcessBlobsHandler(storage db.DB,
	height int64, commitments [][48]byte, blobs [][]byte) error {

	// TODO: verify blob inclusion. Since we include it in the block itself not a sub network its not required

	ranger := file.NewRangeDB(storage)
	// Store the blobs under a single height.
	for i, sidecar := range blobs {
		if err := ranger.Set(uint64(height), commitments[i][:], sidecar); err != nil {
			return err
		}
	}

	return nil
}
