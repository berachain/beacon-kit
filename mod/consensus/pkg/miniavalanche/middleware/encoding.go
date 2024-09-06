package middleware

import "errors"

var (
	// ErrNilBeaconBlockInRequest is an error for when
	// the beacon block in an abci request is nil.
	ErrNilBeaconBlockInRequest = errors.New("nil beacon block in abci request")

	// ErrNoBeaconBlockInRequest is an error for when
	// there is no beacon block in an abci request.
	ErrNoBeaconBlockInRequest = errors.New("no beacon block in abci request")

	// ErrBzIndexOutOfBounds is an error for when the index
	// is out of bounds.
	ErrBzIndexOutOfBounds = errors.New("bzIndex out of bounds")

	// ErrNilABCIRequest is an error for when the abci request
	// is nil.
	ErrNilABCIRequest = errors.New("nil abci request")
)

func UnmarshalBeaconBlockFromOuterBlock(outerBlk OuterBlock, forkVersion uint32) (BeaconBlockT, error) {
	var blk BeaconBlockT
	if outerBlk == nil {
		return blk, ErrNilABCIRequest
	}

	blkBz, err := outerBlk.GetBeaconBlockBytes()
	if err != nil {
		return nil, err
	}
	if blkBz == nil {
		return blk, ErrNilBeaconBlockInRequest
	}
	return blk.NewFromSSZ(blkBz, forkVersion)
}

func UnmarshalBlobSidecarsFromOuterBlock(outerBlk OuterBlock) (BlobSidecarsT, error) {
	var sidecars BlobSidecarsT
	if outerBlk == nil {
		return sidecars, ErrNilABCIRequest
	}

	sidecarBz, err := outerBlk.GetSidecarsBytes()
	if err != nil {
		return nil, err
	}
	if sidecarBz == nil {
		return sidecars, ErrNilBeaconBlockInRequest
	}

	// TODO: Do some research to figure out how to make this more
	// elegant.
	sidecars = sidecars.Empty()
	return sidecars, sidecars.UnmarshalSSZ(sidecarBz)
}
