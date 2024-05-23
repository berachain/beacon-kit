package encoding

import (
	"reflect"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
)

// ExtractBlobsAndBlockFromRequest extracts the blobs and block from an ABCI request.
func ExtractBlobsAndBlockFromRequest[BlobsSidecarsT ssz.Marshallable](
	req ABCIRequest,
	beaconBlkIndex uint,
	blobSidecarsIndex uint,
	forkVersion uint32,
) (types.BeaconBlock, BlobsSidecarsT, error) {
	var blobs BlobsSidecarsT
	if req == nil {
		return nil, blobs, ErrNilABCIRequest
	}

	blk, err := UnmarshalBeaconBlockFromABCIRequest(req, beaconBlkIndex, forkVersion)
	if err != nil {
		return nil, blobs, err
	}

	blobs, err = UnmarshalBlobSidecarsFromABCIRequest[BlobsSidecarsT](req, blobSidecarsIndex)
	if err != nil {
		return nil, blobs, err
	}

	return blk, blobs, nil
}

// UnmarshalBeaconBlockFromABCIRequest extracts a beacon block from an ABCI request.
func UnmarshalBeaconBlockFromABCIRequest(
	req ABCIRequest,
	bzIndex uint,
	forkVersion uint32,
) (types.BeaconBlock, error) {
	if req == nil {
		return nil, ErrNilABCIRequest
	}

	txs := req.GetTxs()
	lenTxs := uint(len(txs))

	// Ensure there are transactions in the request and that the request is valid.
	if txs == nil || lenTxs == 0 {
		return nil, ErrNoBeaconBlockInRequest
	}
	if bzIndex >= lenTxs {
		return nil, ErrBzIndexOutOfBounds
	}

	// Extract the beacon block from the ABCI request.
	blkBz := txs[bzIndex]
	if blkBz == nil {
		return nil, ErrNilBeaconBlockInRequest
	}
	return types.BeaconBlockFromSSZ(blkBz, forkVersion)
}

// UnmarshalBlobSidecarsFromABCIRequest extracts blob sidecars from an ABCI request.
func UnmarshalBlobSidecarsFromABCIRequest[
	T interface{ UnmarshalSSZ([]byte) error },
](
	req ABCIRequest,
	bzIndex uint,
) (T, error) {
	var sidecars T

	sidecars, ok := reflect.New(reflect.TypeOf(sidecars).Elem()).Interface().(T)
	if !ok {
		return sidecars, ErrInvalidType
	}

	if req == nil {
		return sidecars, ErrNilABCIRequest
	}

	txs := req.GetTxs()
	lenTxs := uint(len(txs))

	// Ensure there are transactions in the request and that the request is valid.
	if txs == nil || lenTxs == 0 {
		return sidecars, ErrNoBeaconBlockInRequest
	}
	if bzIndex >= lenTxs {
		return sidecars, ErrBzIndexOutOfBounds
	}

	// Extract the blob sidecars from the ABCI request.
	sidecarBz := txs[bzIndex]
	if sidecarBz == nil {
		return sidecars, ErrNilBeaconBlockInRequest
	}

	err := sidecars.UnmarshalSSZ(sidecarBz)
	return sidecars, err
}
