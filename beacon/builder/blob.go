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
	blobs *enginetypes.BlobsBundleV1,
) ([]byte, error) {
	var blobTx = make([]*types.BlobSidecar, len(blobs.Blobs))
	for i := range len(blobs.Blobs) {
		// Create Inclusion Proof
		inclusionProof, err := kzg.MerkleProofKZGCommitment(blk, i)
		if err != nil {
			return nil, err
		}

		blob := &types.BlobSidecar{
			Blob:           blobs.Blobs[i],
			KzgCommitment:  blobs.Commitments[i],
			KzgProof:       blobs.Proofs[i],
			InclusionProof: inclusionProof,
		}

		blobTx[i] = blob
	}

	bl := types.BlobSidecars{Sidecars: blobTx}

	return bl.MarshalSSZ()
}
