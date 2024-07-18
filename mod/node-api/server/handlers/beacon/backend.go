package beacon

import (
	"context"

	types "github.com/berachain/beacon-kit/mod/node-api/server/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
)

// Backend is the interface for backend of the beacon API.
type Backend[ValidatorT any] interface {
	GetGenesis(ctx context.Context) (common.Root, error)
	GetStateRoot(
		ctx context.Context,
		stateID string,
	) (common.Bytes32, error)
	GetStateValidators(
		ctx context.Context,
		stateID string,
		id []string,
		status []string,
	) ([]*types.ValidatorData[ValidatorT], error)
	GetStateValidator(
		ctx context.Context,
		stateID string,
		validatorID string,
	) (*types.ValidatorData[ValidatorT], error)
	GetStateValidatorBalances(
		ctx context.Context,
		stateID string,
		id []string,
	) ([]*types.ValidatorBalanceData, error)
	GetBlockRewards(
		ctx context.Context,
		blockID string,
	) (*types.BlockRewardsData, error)
}
