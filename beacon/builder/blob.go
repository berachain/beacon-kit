package builder

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/berachain/beacon-kit/beacon/core/types"
	beacontypes "github.com/berachain/beacon-kit/beacon/core/types"
	"github.com/berachain/beacon-kit/crypto/kzg"
	enginetypes "github.com/berachain/beacon-kit/engine/types"
)

// Store the blobs in the blobstore.
func PrepareBlobsHandler(_ sdk.Context,
	height int64, blk beacontypes.BeaconBlock,
	blobs *enginetypes.BlobsBundleV1) ([]byte, error) {

	var blobTx = make([]*types.BlobTxSidecar, 0, len(blobs.Blobs))
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
