package blockchain

import (
	"context"
	"maps"
	"slices"
	"strconv"
	"time"

	"github.com/berachain/beacon-kit/primitives/math"
)

// defaultRetryInterval processes a deposit event.
const defaultRetryInterval = 20 * time.Second

func (s *Service[
	_, ConsensusBlockT, _, _, _, _, _, _, _, _, _, _,
]) depositFetcher(
	ctx context.Context,
	blockNum math.U64,
) {
	if blockNum < s.eth1FollowDistance {
		s.logger.Info(
			"depositFetcher, nothing to fetch",
			"block num", blockNum,
			"eth1FollowDistance", s.eth1FollowDistance,
		)
		return
	}

	s.fetchAndStoreDeposits(ctx, blockNum-s.eth1FollowDistance)
}

func (s *Service[
	_, ConsensusBlockT, _, _, _, _, _, _, _, _, _, _,
]) fetchAndStoreDeposits(
	ctx context.Context,
	blockNum math.U64,
) {
	deposits, err := s.dc.ReadDeposits(ctx, blockNum)
	if err != nil {
		s.logger.Error("Failed to read deposits", "error", err)
		s.metrics.sink.IncrementCounter(
			"beacon_kit.execution.deposit.failed_to_get_block_logs",
			"block_num",
			strconv.FormatUint(blockNum.Unwrap(), 10),
		)
		s.failedBlocksMu.Lock()
		s.failedBlocks[blockNum] = struct{}{}
		s.failedBlocksMu.Unlock()
		return
	}

	if len(deposits) > 0 {
		s.logger.Info(
			"Found deposits on execution layer",
			"block", blockNum, "deposits", len(deposits),
		)
	}

	if err = s.ds.EnqueueDeposits(deposits); err != nil {
		s.logger.Error("Failed to store deposits", "error", err)
		s.failedBlocksMu.Lock()
		s.failedBlocks[blockNum] = struct{}{}
		s.failedBlocksMu.Unlock()
		return
	}

	s.failedBlocksMu.Lock()
	delete(s.failedBlocks, blockNum)
	s.failedBlocksMu.Unlock()
}

func (s *Service[
	_, ConsensusBlockT, _, _, _, _, _, _, _, _, _, _,
]) depositCatchupFetcher(ctx context.Context) {
	ticker := time.NewTicker(defaultRetryInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.failedBlocksMu.RLock()
			failedBlks := slices.Collect(maps.Keys(s.failedBlocks))
			s.failedBlocksMu.RUnlock()
			if len(failedBlks) == 0 {
				continue
			}
			s.logger.Warn(
				"Failed to get deposits from block(s), retrying...",
				"num_blocks",
				failedBlks,
			)

			// Fetch deposits for blocks that failed to be processed.
			for _, blockNum := range failedBlks {
				s.fetchAndStoreDeposits(ctx, blockNum)
			}
		}
	}
}
