package staking

import (
	stakinglogs "github.com/berachain/beacon-kit/beacon/staking/logs"
	coretypes "github.com/ethereum/go-ethereum/core/types"
)

func (s *Service) ProcessBlockEvents(
	logs []*coretypes.Log,
) error {
	var err error
	for _, log := range logs {
		// We only care about logs from the deposit
		if log.Address != s.BeaconCfg().Execution.DepositContractAddress {
			continue
		}

		switch logSig := log.Topics[0]; {
		case logSig == stakinglogs.DepositSig:
			err = s.addDepositToQueue()
		case logSig == stakinglogs.RedirectSig:
			err = s.addRedirectToQueue()
		case logSig == stakinglogs.WithdrawalSig:
			err = s.addWithdrawalToQueue()
		default:
			continue
		}
	}
	return err
}

// addDepositToQueue adds a deposit to the queue.
func (s *Service) addDepositToQueue() error {
	return nil
}

// addRedirectToQueue adds a redirect to the queue.
func (s *Service) addRedirectToQueue() error {
	return nil
}

// addWithdrawalToQueue adds a withdrawal to the queue.
func (s *Service) addWithdrawalToQueue() error {
	return nil
}
