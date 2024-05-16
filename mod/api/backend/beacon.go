package backend

import (
	"context"

	"github.com/berachain/beacon-kit/mod/primitives"
)

func (h Backend) GetGenesis(ctx context.Context) (primitives.Root, error) {
	return h.getNewStateDB(ctx).GetGenesisValidatorsRoot()
}

func (h Backend) GetStateHashTreeRoot(ctx context.Context, stateId string) {

}

func (h Backend) GetStateFork(ctx context.Context, stateId string) {

}

func (h Backend) GetStateFinalityCheckpoints(stateId string) {

}

func (h Backend) GetStateValidators(stateId string, id []string, status []string) {

}

func (h Backend) GetStateValidator(stateId string, validatorId string) {

}

func (h Backend) GetStateValidatorBalances(stateId string, id []string) {

}

func (h Backend) PostStateValidatorBalances(stateId string) {

}

func (h Backend) GetStateCommittees(stateId string, epoch string, index string, slot string) {

}

func (h Backend) GetStateSyncCommittees(stateId string, epoch string) {

}

func (h Backend) GetStateRandao(stateId string, epoch string) {

}

func (h Backend) GetBlockHeaders(slot string, parentRoot string) {

}

func (h Backend) GetBlockHeader(blockId string) {

}

func (h Backend) GetBlock(blockId string) {

}

func (h Backend) GetBlockRoot(blockId string) {

}

func (h Backend) GetBlockAttestations(blockId string) {

}

func (h Backend) GetBlobSidecars(blockId string, indices []string) {

}

func (h Backend) GetSyncCommiteeAwards(blockId string) {

}

func (h Backend) GetDepositSnapshot() {

}

func (h Backend) GetBlockAwards(blockId string) {

}

func (h Backend) GetAttestationRewards(epoch string, ids []string) {

}

func (h Backend) GetBlindedBlock(blockId string) {

}

func (h Backend) GetPoolAttestations(slot string, committee_index string) {

}

func (h Backend) PostPoolAttestations(thisoneisfucked string) {

}

func (h Backend) PostPoolSyncCommitteeSignature() {

}

func (h Backend) GetPoolVoluntaryExits() {

}

func (h Backend) PostPoolVoluntaryExits() {

}

func (h Backend) GetPoolSignedBLSExecutionChanges() {

}

func (h Backend) PostPoolSignedBLSExecutionChanges() {

}
