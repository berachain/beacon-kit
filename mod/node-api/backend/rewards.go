package backend

import (
	"context"

	"github.com/berachain/beacon-kit/mod/node-api/server/types"
)

func (h Backend) GetBlockRewards(ctx context.Context, blockId string) (*types.BlockRewardsData, error) {
	return &types.BlockRewardsData{
		ProposerIndex:     123,
		Total:             123,
		Attestations:      123,
		SyncAggregate:     123,
		ProposerSlashings: 123,
		AttesterSlashings: 123,
	}, nil
}
