package types

import (
	"context"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives"
)

type BackendHandlers interface {
	GetGenesis(ctx context.Context) (primitives.Root, error)
	GetStateRoot(ctx context.Context, stateId string) (primitives.Bytes32, error)
	GetStateValidators(ctx context.Context, stateId string, id []string, status []string) ([]*ValidatorData, error)
	GetStateValidator(ctx context.Context, stateId string, validatorId string) (*ValidatorData, error)
	GetStateValidatorBalances(ctx context.Context, stateId string, id []string) ([]*ValidatorBalanceData, error)
	GetStateCommittees(ctx context.Context, stateId string, epoch string, index string, slot string)
	GetStateSyncCommittees(ctx context.Context, stateId string, epoch string)
	GetBlockHeaders(ctx context.Context, slot string, parentRoot string) (*types.BeaconBlockHeader, error)
	GetBlockHeader(ctx context.Context, blockId string) (*types.BeaconBlockHeader, error)
	GetBlock(ctx context.Context, blockId string)
	GetBlobSidecars(ctx context.Context, blockId string, indices []string)
	GetPoolVoluntaryExits(ctx context.Context)
	PostPoolVoluntaryExits(ctx context.Context)
	GetPoolSignedBLSExecutionChanges(ctx context.Context)
	PostPoolSignedBLSExecutionChanges(ctx context.Context)
	GetBlockProposerDuties(ctx context.Context, epoch string)
	GetConfigSpec(ctx context.Context)
}
