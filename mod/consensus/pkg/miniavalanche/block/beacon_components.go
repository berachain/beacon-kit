package block

import (
	"errors"

	"github.com/berachain/beacon-kit/mod/consensus/pkg/miniavalanche"
)

var (
	// errNilBeaconBlock is an error for when
	// the beacon block in an abci request is nil.
	errNilBeaconBlock = errors.New("nil beacon block")

	// errNilBlock is an error for when the abci request
	// is nil.
	errNilBlock = errors.New("nil abci request")
)

func (b *StatelessBlock) GetBeaconBlock(forkVersion uint32) (miniavalanche.BeaconBlockT, error) {
	var blk miniavalanche.BeaconBlockT
	if b == nil {
		return blk, errNilBlock
	}
	if b.BlkContent.BeaconBlockByte == nil {
		return blk, errNilBeaconBlock
	}
	return blk.NewFromSSZ(b.BlkContent.BeaconBlockByte, forkVersion)
}

func (b *StatelessBlock) GetBlobSidecars() (miniavalanche.BlobSidecarsT, error) {
	var sidecars miniavalanche.BlobSidecarsT
	if b == nil {
		return sidecars, errNilBlock
	}
	if b.BlkContent.BlobsBytes == nil {
		return sidecars, errNilBeaconBlock
	}

	// TODO: Do some research to figure out how to make this more
	// elegant.
	sidecars = sidecars.Empty()
	return sidecars, sidecars.UnmarshalSSZ(b.BlkContent.BlobsBytes)
}
