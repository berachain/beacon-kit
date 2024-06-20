package attributes

import "github.com/berachain/beacon-kit/mod/primitives/pkg/common"

// BeaconState is an interface for accessing the beacon state.
type BeaconState[WithdrawalT any] interface {
	// ExpectedWithrawals returns the expected withdrawals.
	ExpectedWithdrawals() ([]WithdrawalT, error)
	// GetRandaoMixAtIndex returns the randao mix at the given index.
	GetRandaoMixAtIndex(index uint64) (common.Root, error)
}
