package deposit

import (
	"context"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/ethereum/go-ethereum/event"
)

type Service[DepositStoreT DepositStore] struct {
	feed *event.FeedOf[types.BlockEvent]
	dc   DepositContract
	sb   StorageBackend[
		any, any, any, DepositStoreT,
	]

	logger log.Logger[any]
}

func NewService[DepositStoreT DepositStore](
	feed *event.FeedOf[types.BlockEvent],
	logger log.Logger[any],
	sb StorageBackend[
		any, any, any, DepositStoreT,
	],
	dc DepositContract,
) *Service[DepositStoreT] {
	return &Service[DepositStoreT]{
		feed:   feed,
		logger: logger,
		sb:     sb,
		dc:     dc,
	}
}

func (s *Service[DepositStoreT]) Start(ctx context.Context) error {
	ch := make(chan types.BlockEvent)
	feed := s.feed.Subscribe(ch)
	go func() {
		for {
			select {
			case <-ctx.Done():
				feed.Unsubscribe()
			case event := <-ch:
				if err := s.handleDepositEvent(event); err != nil {
					s.logger.Error("failed to handle deposit event", "err", err)
				}
			}
		}
	}()

	return nil
}

func (s *Service[DepositStoreT]) Name() string {
	return "deposit-handler"
}

func (s *Service[DepositStoreT]) Status() error {
	return nil
}

func (s *Service[DepositStoreT]) WaitForHealthy(_ context.Context) {}
