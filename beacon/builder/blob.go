package builder

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/berachain/beacon-kit/beacon/core/types"
	beacontypes "github.com/berachain/beacon-kit/beacon/core/types"
	"github.com/berachain/beacon-kit/crypto/kzg"
	"github.com/berachain/beacon-kit/db"
	"github.com/berachain/beacon-kit/db/file"
	enginetypes "github.com/berachain/beacon-kit/engine/types"
)

// Store the blobs in the blobstore.
func PrepareBlobsHandler(ctx sdk.Context, storage db.DB,
	height int64, blk beacontypes.BeaconBlock,
	blobs *enginetypes.BlobsBundleV1) ([]byte, error) {

	ranger := file.NewRangeDB(storage)
	var blobTx = make([]*types.BlobTxSidecar, 0, len(blobs.Blobs))
	for i, sidecar := range blobs.Blobs {
		//Create Inclusion Proof
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

		if err := ranger.Set(uint64(height), blobs.Commitments[i], sidecar); err != nil {
			return nil, err
		}

		blobTx[i] = blob
	}

	bl := types.BlobSidecars{BlobSidecars: blobTx}

	return bl.MarshalSSZ()
}

// Store the blobs in the blobstore.
func ProcessBlobsHandler(ctx sdk.Context, storage db.DB,
	height int64, blobTx []byte) error {

	bl := types.BlobSidecars{}
	bl.UnmarshalSSZ(blobTx)

	ranger := file.NewRangeDB(storage)
	// Store the blobs under a single height.
	for i, sidecar := range bl.BlobSidecars {
		if err := kzg.VerifyKZGInclusionProof([]byte{}, sidecar, uint64(i)); err != nil {
			return err
		}
		if err := ranger.Set(uint64(height), sidecar.KzgCommitment, sidecar.Blob); err != nil {
			return err
		}
	}

	return nil
}
