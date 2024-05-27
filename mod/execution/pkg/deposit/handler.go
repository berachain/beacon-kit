package deposit

import (
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
)

// TODO: eth1FollowDistance should be done actually properly
const eth1FollowDistance = 1

func (s *Service[DepositStoreT]) handleDepositEvent(e types.BlockEvent) error {
	slot := e.Block().GetSlot()
	slot = slot - eth1FollowDistance
	s.logger.Info("💵 processing deposit logs 💵", "slot", slot)
	deposits, err := s.dc.GetDeposits(e.Context(), slot.Unwrap())
	if err != nil {
		return err
	}

	if err := s.sb.DepositStore(e.Context()).EnqueueDeposits(
		deposits,
	); err != nil {
		return err
	}
	return nil
}
