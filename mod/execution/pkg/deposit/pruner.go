package deposit

import "github.com/berachain/beacon-kit/mod/primitives"

func GetPruneParamsFn[
	WithdrawalCredentialsT any,
	DepositT Deposit[DepositT, WithdrawalCredentialsT],
	BeaconBlockBodyT BeaconBlockBody[DepositT],
	BeaconBlockT BeaconBlock[DepositT, BeaconBlockBodyT],
	BlockEventT BlockEvent[
		DepositT, BeaconBlockBodyT, BeaconBlockT,
	],
](cs primitives.ChainSpec) func(BlockEventT) (uint64, uint64) {
	return func(event BlockEventT) (uint64, uint64) {
		blk := event.Block()
		deposits := blk.GetBody().GetDeposits()
		if len(deposits) == 0 {
			return 0, 0
		}
		index := deposits[len(deposits)-1].GetIndex()

		return max(0, index-cs.MaxDepositsPerBlock()),
			min(index, cs.MaxDepositsPerBlock())
	}
}
