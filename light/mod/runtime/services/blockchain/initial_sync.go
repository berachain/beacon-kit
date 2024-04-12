package blockchain

import (
	"context"

	"github.com/berachain/beacon-kit/mod/config/version"
	"github.com/berachain/beacon-kit/mod/execution"
	"github.com/berachain/beacon-kit/mod/primitives/engine"
	"github.com/ethereum/go-ethereum/common"
)

// InitialSync performs the initial sync of the blockchain.
// This is done by fetching the latest exec payload from the store
// and sending a forkchoice update to the execution engine.
func (s *Service) initialSync(ctx context.Context) error {
	s.Logger().Info("starting lightchain service initial sync")

	// Get the latest execution payload from the store.
	latestPayload, err := s.BeaconState(ctx).GetLatestExecutionPayload()
	if err != nil {
		return err
	}

	latestHeader, err := s.BeaconState(ctx).GetLatestBlockHeader()
	if err != nil {
		return err
	}

	// Notify the execution client of a new payload.
	_, err = s.ee.VerifyAndNotifyNewPayload(
		ctx,
		execution.BuildNewPayloadRequest(
			latestPayload,
			[]common.Hash{},
			&latestHeader.ParentRoot,
		),
	)
	if err != nil {
		return err
	}

	eth1BlockHash := latestPayload.GetBlockHash()

	// Send a forkchoice update to the execution engine.
	_, _, err = s.ee.NotifyForkchoiceUpdate(
		ctx,
		&execution.ForkchoiceUpdateRequest{
			State: &engine.ForkchoiceState{
				HeadBlockHash:      eth1BlockHash,
				SafeBlockHash:      eth1BlockHash,
				FinalizedBlockHash: eth1BlockHash,
			},
			ForkVersion: version.Deneb,
		},
	)
	return err
}
