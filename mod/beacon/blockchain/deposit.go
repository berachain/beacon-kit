package blockchain

import (
	"context"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// ProcessLogsInETH1Block gets logs in the Eth1 block
// received from the execution client and processes them to
// convert them into appropriate objects that can be consumed
// by other services.
func (s *Service[BeaconStateT, BlobSidecarsT, DepositStoreT]) ProcessLogsInETH1Block(
	ctx context.Context,
	blockNumber math.U64,
) error {
	deposits, err := s.bdc.GetDeposits(ctx, blockNumber.Unwrap())
	if err != nil {
		return err
	}

	return s.bsb.DepositStore(ctx).EnqueueDeposits(deposits)
}

// PruneDepositEvents prunes deposit events.
func (s *Service[BeaconStateT, BlobSidecarsT, DepositStoreT]) PruneDepositEvents(idx uint64) error {
	return s.bsb.DepositStore(context.Background()).PruneToIndex(idx)
}
