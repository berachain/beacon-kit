package types

import "github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"

type BlockBundle struct {
	Block    *types.BeaconBlock
	Sidecars *BlobSidecars
}

func (bb *BlockBundle) GetBeaconBlock() *types.BeaconBlock {
	return bb.Block
}

func (bb *BlockBundle) GetSidecars() *BlobSidecars {
	return bb.Sidecars
}

func (bb *BlockBundle) SetBeaconBlock(block *types.BeaconBlock) {
	bb.Block = block
}

func (bb *BlockBundle) SetSidecars(sidecars *BlobSidecars) {
	bb.Sidecars = sidecars
}
