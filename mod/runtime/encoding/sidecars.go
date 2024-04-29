package encoding

import datypes "github.com/berachain/beacon-kit/mod/da/types"

func UnmarshalBlobSidecarsFromABCIRequest(
	req ABCIRequest,
	bzIndex uint,
) (*datypes.BlobSidecars, error) {
	if req == nil {
		return nil, ErrNilABCIRequest
	}

	txs := req.GetTxs()

	// Ensure there are transactions in the request and
	// that the request is valid.
	if lenTxs := uint(len(txs)); txs == nil || lenTxs == 0 {
		return nil, ErrNoBeaconBlockInRequest
	} else if bzIndex >= uint(len(txs)) {
		return nil, ErrBzIndexOutOfBounds
	}

	// Extract the beacon block from the ABCI request.
	sidecarBz := txs[bzIndex]
	if sidecarBz == nil {
		return nil, ErrNilBeaconBlockInRequest
	}

	var sidecars datypes.BlobSidecars
	if err := sidecars.UnmarshalSSZ(sidecarBz); err != nil {
		return nil, err
	}

	return &sidecars, nil
}
