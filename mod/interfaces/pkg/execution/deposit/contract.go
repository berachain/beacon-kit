package deposit

import (
	"context"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// Contract is an interface for reading deposits from a BeaconDepositContract.
type Contract[DepositT any] interface {
	// ReadDeposits reads deposits from the deposit contract.
	ReadDeposits(
		ctx context.Context,
		blkNum math.U64,
	) ([]DepositT, error)
}
