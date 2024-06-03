package deposit

import (
	"context"
)

// depositFetcher processes a deposit event.
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BlockEventT, ExecutionPayloadT, SubscriptionT, DepositT,
]) depositFetcher(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case blk := <-s.newBlock:
			// querierBlockNum := eth1Final.Number().Uint64() - uint64(s.eth1FollowDistance)
			querierBlockNum := blk.GetBody().GetExecutionPayload().GetNumber() - s.eth1FollowDistance
			// Use a goroutine to handle deposits asynchronously to improve performance.

			deposits, err := s.dc.ReadDeposits(ctx, querierBlockNum)
			if err != nil {
				s.logger.Error("Failed to read deposits", "error", err)
				continue
			}

			if len(deposits) == 0 {
				s.logger.Info(
					"waiting for deposits from execution layer", "block", querierBlockNum,
				)
			} else {
				s.logger.Info(
					"found deposits on execution layer",
					"block", querierBlockNum, "deposits", len(deposits),
				)
			}

			if err = s.ds.EnqueueDeposits(deposits); err != nil {
				s.logger.Error("Failed to enqueue deposits", "error", err)
				continue
			}
		}
	}
}
