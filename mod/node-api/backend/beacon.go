package backend

import (
	"context"
	"errors"
	"strconv"

	types "github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	serverType "github.com/berachain/beacon-kit/mod/node-api/server/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

func (h Backend) GetGenesis(ctx context.Context) (primitives.Root, error) {
	// needs genesis_time and gensis_fork_version
	return h.getNewStateDB(ctx, "stateID").GetGenesisValidatorsRoot()
}

func (h Backend) GetStateRoot(ctx context.Context, stateId string) (primitives.Bytes32, error) {
	stateDb := h.getNewStateDB(ctx, stateId)
	slot, err := stateDb.GetSlot()
	if err != nil {
		return primitives.Bytes32{}, err
	}
	block, err := stateDb.StateRootAtIndex(slot.Unwrap())
	if err != nil {
		return primitives.Bytes32{}, err
	}
	root, err := block.HashTreeRoot()
	if err != nil {
		return primitives.Bytes32{}, err
	}
	return root, nil
}

func (h Backend) GetStateFork(ctx context.Context, stateId string) (*types.Fork, error) {
	// resolve state_id
	// currently can only get it for the current epoch
	return h.getNewStateDB(ctx, "stateID").GetFork()
}

func (h Backend) GetStateValidators(ctx context.Context, stateId string, id []string, status []string) ([]*serverType.ValidatorData, error) {
	stateDb := h.getNewStateDB(ctx, stateId)
	validators := make([]*serverType.ValidatorData, 0)
	for _, indexOrKey := range id {
		if index, err := getValidatorIndex(stateDb, indexOrKey); err == nil {
			if validator, err := stateDb.ValidatorByIndex(index); err == nil {
				if balance, err := stateDb.GetBalance(index); err == nil {
					validators = append(validators, &serverType.ValidatorData{
						Index:     index,
						Balance:   balance,
						Status:    "active",
						Validator: validator,
					})
				}
			}
		}
	}
	return validators, nil
}

func getValidatorIndex(stateDb StateDB, keyOrIndex string) (math.U64, error) {
	if index, err := strconv.ParseUint(keyOrIndex, 10, 64); err == nil {
		return math.U64(index), nil
	}
	key := crypto.BLSPubkey{}
	key.UnmarshalText([]byte(keyOrIndex))
	index, err := stateDb.ValidatorIndexByPubkey(key)
	if err == nil {
		return math.U64(index), nil
	}
	return math.U64(0), err
}

func (h Backend) GetStateValidator(ctx context.Context, stateId string, validatorId string) (*serverType.ValidatorData, error) {
	stateDb := h.getNewStateDB(ctx, stateId)
	if index, err := getValidatorIndex(stateDb, validatorId); err == nil {
		if validator, err := stateDb.ValidatorByIndex(index); err == nil {
			if balance, err := stateDb.GetBalance(index); err == nil {
				return &serverType.ValidatorData{
					Index:     index,
					Balance:   balance,
					Status:    "active",
					Validator: validator,
				}, nil
			}
		}
	}
	return nil, errors.New("Validator not found")
}

func (h Backend) GetStateValidatorBalances(ctx context.Context, stateId string, id []string) ([]*serverType.ValidatorBalanceData, error) {
	stateDb := h.getNewStateDB(ctx, stateId)
	balances := make([]*serverType.ValidatorBalanceData, 0)
	for _, indexOrKey := range id {
		if index, err := getValidatorIndex(stateDb, indexOrKey); err == nil {
			if balance, err := stateDb.GetBalance(index); err == nil {
				balances = append(balances, &serverType.ValidatorBalanceData{
					Index:   index,
					Balance: balance,
				})
			}
		}
	}
	return balances, nil
}

func (h Backend) GetStateCommittees(ctx context.Context, stateId string, epoch string, index string, slot string) {

}

func (h Backend) GetStateSyncCommittees(ctx context.Context, stateId string, epoch string) {

}

func (h Backend) GetStateRandao(ctx context.Context, stateId string, epoch string) {

}

func (h Backend) GetBlockHeaders(ctx context.Context, slot string, parentRoot string) (*types.BeaconBlockHeader, error) {
	return h.getNewStateDB(context.TODO(), "stateID").GetLatestBlockHeader()
}

func (h Backend) GetBlockHeader(ctx context.Context, blockId string) (*types.BeaconBlockHeader, error) {
	return h.getNewStateDB(context.TODO(), "stateid").GetLatestBlockHeader()
}

func (h Backend) GetBlock(ctx context.Context, blockId string) {
	// return h.getNewStateDB(ctx).G
}

func (h Backend) GetBlockRoot(ctx context.Context, blockId string) (primitives.Bytes32, error) {
	stateDb := h.getNewStateDB(ctx, "stateID")
	slot, err := stateDb.GetSlot()
	if err != nil {
		return primitives.Bytes32{}, err
	}
	block, err := h.getNewStateDB(context.TODO(), "stateID").GetBlockRootAtIndex(slot.Unwrap())
	if err != nil {
		return primitives.Bytes32{}, err
	}
	root, err := block.HashTreeRoot()
	if err != nil {
		return primitives.Bytes32{}, err
	}
	return root, nil
}

func (h Backend) GetBlockAttestations(ctx context.Context, blockId string) {

}

func (h Backend) GetBlobSidecars(ctx context.Context, blockId string, indices []string) {

}

func (h Backend) GetSyncCommiteeAwards(ctx context.Context, blockId string) {

}

func (h Backend) GetDepositSnapshot(ctx context.Context) {

}

func (h Backend) GetBlockAwards(ctx context.Context, blockId string) {

}

func (h Backend) GetAttestationRewards(ctx context.Context, epoch string, ids []string) {

}

func (h Backend) GetBlindedBlock(ctx context.Context, blockId string) {

}

func (h Backend) GetPoolAttestations(ctx context.Context, slot string, committee_index string) {

}

func (h Backend) PostPoolAttestations(ctx context.Context, thisoneisfucked string) {

}

func (h Backend) PostPoolSyncCommitteeSignature(ctx context.Context) {

}

func (h Backend) GetPoolVoluntaryExits(ctx context.Context) {

}

func (h Backend) PostPoolVoluntaryExits(ctx context.Context) {

}

func (h Backend) GetPoolSignedBLSExecutionChanges(ctx context.Context) {

}

func (h Backend) PostPoolSignedBLSExecutionChanges(ctx context.Context) {

}
