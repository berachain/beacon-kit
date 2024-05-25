package events

import (
	"context"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
)

// BlockWithState is a struct that contains a block with state.
type BlockWithState struct {
	ctx context.Context
	// Block is the block.
	Block types.BeaconBlock
}
