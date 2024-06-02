package deposit

import "context"

// PrunerFn is a function that prunes the deposit store up to the specified index.
func PrunerFn[
	BeaconBlockT interface {
		GetBody() BeaconBlockBodyT
	},
	BeaconBlockBodyT interface {
		GetDeposits() []DepositT
	},
	DepositT interface {
		GetIndex() uint64
	},
](ctx context.Context, blk BeaconBlockT) uint64 {
	deposits := blk.GetBody().GetDeposits()
	if len(deposits) == 0 {
		return 0
	}

	return deposits[len(deposits)-1].GetIndex()
}
