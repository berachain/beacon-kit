package builder

import (
	"github.com/berachain/beacon-kit/beacon/core/types"
	beacontypes "github.com/berachain/beacon-kit/beacon/core/types"
	"github.com/berachain/beacon-kit/crypto/kzg"
	enginetypes "github.com/berachain/beacon-kit/engine/types"
)

// PrepareBlobsHandler is responsible for attaching an inclusion proof to the
// blob sidecar
func PrepareBlobsHandler(
	height int64, blk beacontypes.BeaconBlock,
	blobs *enginetypes.BlobsBundleV1) ([]byte, error) {

	var blobTx = make([]*types.BlobTxSidecar, len(blobs.Blobs))
	for i, sidecar := range blobs.Blobs {
		// Create Inclusion Proof
		ic, err := kzg.MerkleProofKZGCommitment(blk, i)
		if err != nil {
			return nil, err
		}
		blob := &types.BlobTxSidecar{
			Blob:           sidecar,
			KzgCommitment:  blobs.Commitments[i],
			KzgProof:       blobs.Proofs[i],
			InclusionProof: ic,
		}

		blobTx[i] = blob
	}

	bl := types.BlobSidecars{BlobSidecars: blobTx}

	return bl.MarshalSSZ()
}
