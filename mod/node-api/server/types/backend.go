package types

import (
	"context"

	"github.com/berachain/beacon-kit/mod/primitives"
)

type BackendHandlers interface {
	GetGenesis(ctx context.Context) (primitives.Root, error)
	GetStateRoot(ctx context.Context, stateId string) (primitives.Bytes32, error)
	GetStateValidators(ctx context.Context, stateId string, id []string, status []string) ([]*ValidatorData, error)
	GetStateValidator(ctx context.Context, stateId string, validatorId string) (*ValidatorData, error)
	GetStateValidatorBalances(ctx context.Context, stateId string, id []string) ([]*ValidatorBalanceData, error)
	GetBlockRewards(ctx context.Context, blockId string) (*BlockRewardsData, error)
}
