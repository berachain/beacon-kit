package api

import (
	"context"

	"github.com/berachain/beacon-kit/api/beaconnode"
)

type MyServer struct {
}

// Eventstream implements beaconnode.Handler.
func (*MyServer) Eventstream(ctx context.Context, params beaconnode.EventstreamParams) (beaconnode.EventstreamRes, error) {
	panic("unimplemented")
}

// GetAggregatedAttestation implements beaconnode.Handler.
func (*MyServer) GetAggregatedAttestation(ctx context.Context, params beaconnode.GetAggregatedAttestationParams) (beaconnode.GetAggregatedAttestationRes, error) {
	panic("unimplemented")
}

// GetAttestationsRewards implements beaconnode.Handler.
func (*MyServer) GetAttestationsRewards(ctx context.Context, req []string, params beaconnode.GetAttestationsRewardsParams) (beaconnode.GetAttestationsRewardsRes, error) {
	panic("unimplemented")
}

// GetAttesterDuties implements beaconnode.Handler.
func (*MyServer) GetAttesterDuties(ctx context.Context, req []string, params beaconnode.GetAttesterDutiesParams) (beaconnode.GetAttesterDutiesRes, error) {
	panic("unimplemented")
}

// GetBlindedBlock implements beaconnode.Handler.
func (*MyServer) GetBlindedBlock(ctx context.Context, params beaconnode.GetBlindedBlockParams) (beaconnode.GetBlindedBlockRes, error) {
	panic("unimplemented")
}

// GetBlobSidecars implements beaconnode.Handler.
func (*MyServer) GetBlobSidecars(ctx context.Context, params beaconnode.GetBlobSidecarsParams) (beaconnode.GetBlobSidecarsRes, error) {
	panic("unimplemented")
}

// GetBlockAttestations implements beaconnode.Handler.
func (*MyServer) GetBlockAttestations(ctx context.Context, params beaconnode.GetBlockAttestationsParams) (beaconnode.GetBlockAttestationsRes, error) {
	panic("unimplemented")
}

// GetBlockHeader implements beaconnode.Handler.
func (*MyServer) GetBlockHeader(ctx context.Context, params beaconnode.GetBlockHeaderParams) (beaconnode.GetBlockHeaderRes, error) {
	panic("unimplemented")
}

// GetBlockHeaders implements beaconnode.Handler.
func (*MyServer) GetBlockHeaders(ctx context.Context, params beaconnode.GetBlockHeadersParams) (beaconnode.GetBlockHeadersRes, error) {
	panic("unimplemented")
}

// GetBlockRewards implements beaconnode.Handler.
func (*MyServer) GetBlockRewards(ctx context.Context, params beaconnode.GetBlockRewardsParams) (beaconnode.GetBlockRewardsRes, error) {
	panic("unimplemented")
}

// GetBlockRoot implements beaconnode.Handler.
func (*MyServer) GetBlockRoot(ctx context.Context, params beaconnode.GetBlockRootParams) (beaconnode.GetBlockRootRes, error) {
	panic("unimplemented")
}

// GetBlockV2 implements beaconnode.Handler.
func (*MyServer) GetBlockV2(ctx context.Context, params beaconnode.GetBlockV2Params) (beaconnode.GetBlockV2Res, error) {
	panic("unimplemented")
}

// GetDebugChainHeadsV2 implements beaconnode.Handler.
func (*MyServer) GetDebugChainHeadsV2(ctx context.Context) (beaconnode.GetDebugChainHeadsV2Res, error) {
	panic("unimplemented")
}

// GetDebugForkChoice implements beaconnode.Handler.
func (*MyServer) GetDebugForkChoice(ctx context.Context) (beaconnode.GetDebugForkChoiceRes, error) {
	panic("unimplemented")
}

// GetDepositContract implements beaconnode.Handler.
func (*MyServer) GetDepositContract(ctx context.Context) (beaconnode.GetDepositContractRes, error) {
	panic("unimplemented")
}

// GetDepositSnapshot implements beaconnode.Handler.
func (*MyServer) GetDepositSnapshot(ctx context.Context) (beaconnode.GetDepositSnapshotRes, error) {
	panic("unimplemented")
}

// GetEpochCommittees implements beaconnode.Handler.
func (*MyServer) GetEpochCommittees(ctx context.Context, params beaconnode.GetEpochCommitteesParams) (beaconnode.GetEpochCommitteesRes, error) {
	panic("unimplemented")
}

// GetEpochSyncCommittees implements beaconnode.Handler.
func (*MyServer) GetEpochSyncCommittees(ctx context.Context, params beaconnode.GetEpochSyncCommitteesParams) (beaconnode.GetEpochSyncCommitteesRes, error) {
	panic("unimplemented")
}

// GetForkSchedule implements beaconnode.Handler.
func (*MyServer) GetForkSchedule(ctx context.Context) (beaconnode.GetForkScheduleRes, error) {
	panic("unimplemented")
}

// GetGenesis implements beaconnode.Handler.
func (*MyServer) GetGenesis(ctx context.Context) (beaconnode.GetGenesisRes, error) {
	panic("unimplemented")
}

// GetHealth implements beaconnode.Handler.
func (*MyServer) GetHealth(ctx context.Context, params beaconnode.GetHealthParams) (beaconnode.GetHealthRes, error) {
	panic("unimplemented")
}

// GetLightClientBootstrap implements beaconnode.Handler.
func (*MyServer) GetLightClientBootstrap(ctx context.Context, params beaconnode.GetLightClientBootstrapParams) (beaconnode.GetLightClientBootstrapRes, error) {
	panic("unimplemented")
}

// GetLightClientFinalityUpdate implements beaconnode.Handler.
func (*MyServer) GetLightClientFinalityUpdate(ctx context.Context) (beaconnode.GetLightClientFinalityUpdateRes, error) {
	panic("unimplemented")
}

// GetLightClientOptimisticUpdate implements beaconnode.Handler.
func (*MyServer) GetLightClientOptimisticUpdate(ctx context.Context) (beaconnode.GetLightClientOptimisticUpdateRes, error) {
	panic("unimplemented")
}

// GetLightClientUpdatesByRange implements beaconnode.Handler.
func (*MyServer) GetLightClientUpdatesByRange(ctx context.Context, params beaconnode.GetLightClientUpdatesByRangeParams) (beaconnode.GetLightClientUpdatesByRangeRes, error) {
	panic("unimplemented")
}

// GetLiveness implements beaconnode.Handler.
func (*MyServer) GetLiveness(ctx context.Context, req []string, params beaconnode.GetLivenessParams) (beaconnode.GetLivenessRes, error) {
	panic("unimplemented")
}

// GetNetworkIdentity implements beaconnode.Handler.
func (*MyServer) GetNetworkIdentity(ctx context.Context) (beaconnode.GetNetworkIdentityRes, error) {
	panic("unimplemented")
}

// GetNextWithdrawals implements beaconnode.Handler.
func (*MyServer) GetNextWithdrawals(ctx context.Context, params beaconnode.GetNextWithdrawalsParams) (beaconnode.GetNextWithdrawalsRes, error) {
	panic("unimplemented")
}

// GetNodeVersion implements beaconnode.Handler.
func (*MyServer) GetNodeVersion(ctx context.Context) (beaconnode.GetNodeVersionRes, error) {
	panic("unimplemented")
}

// GetPeer implements beaconnode.Handler.
func (*MyServer) GetPeer(ctx context.Context, params beaconnode.GetPeerParams) (beaconnode.GetPeerRes, error) {
	panic("unimplemented")
}

// GetPeerCount implements beaconnode.Handler.
func (*MyServer) GetPeerCount(ctx context.Context) (beaconnode.GetPeerCountRes, error) {
	panic("unimplemented")
}

// GetPeers implements beaconnode.Handler.
func (*MyServer) GetPeers(ctx context.Context, params beaconnode.GetPeersParams) (beaconnode.GetPeersRes, error) {
	panic("unimplemented")
}

// GetPoolAttestations implements beaconnode.Handler.
func (*MyServer) GetPoolAttestations(ctx context.Context, params beaconnode.GetPoolAttestationsParams) (beaconnode.GetPoolAttestationsRes, error) {
	panic("unimplemented")
}

// GetPoolAttesterSlashings implements beaconnode.Handler.
func (*MyServer) GetPoolAttesterSlashings(ctx context.Context) (beaconnode.GetPoolAttesterSlashingsRes, error) {
	panic("unimplemented")
}

// GetPoolBLSToExecutionChanges implements beaconnode.Handler.
func (*MyServer) GetPoolBLSToExecutionChanges(ctx context.Context) (beaconnode.GetPoolBLSToExecutionChangesRes, error) {
	panic("unimplemented")
}

// GetPoolProposerSlashings implements beaconnode.Handler.
func (*MyServer) GetPoolProposerSlashings(ctx context.Context) (beaconnode.GetPoolProposerSlashingsRes, error) {
	panic("unimplemented")
}

// GetPoolVoluntaryExits implements beaconnode.Handler.
func (*MyServer) GetPoolVoluntaryExits(ctx context.Context) (beaconnode.GetPoolVoluntaryExitsRes, error) {
	panic("unimplemented")
}

// GetProposerDuties implements beaconnode.Handler.
func (*MyServer) GetProposerDuties(ctx context.Context, params beaconnode.GetProposerDutiesParams) (beaconnode.GetProposerDutiesRes, error) {
	panic("unimplemented")
}

// GetSpec implements beaconnode.Handler.
func (*MyServer) GetSpec(ctx context.Context) (beaconnode.GetSpecRes, error) {
	panic("unimplemented")
}

// GetStateFinalityCheckpoints implements beaconnode.Handler.
func (*MyServer) GetStateFinalityCheckpoints(ctx context.Context, params beaconnode.GetStateFinalityCheckpointsParams) (beaconnode.GetStateFinalityCheckpointsRes, error) {
	panic("unimplemented")
}

// GetStateFork implements beaconnode.Handler.
func (*MyServer) GetStateFork(ctx context.Context, params beaconnode.GetStateForkParams) (beaconnode.GetStateForkRes, error) {
	panic("unimplemented")
}

// GetStateRandao implements beaconnode.Handler.
func (*MyServer) GetStateRandao(ctx context.Context, params beaconnode.GetStateRandaoParams) (beaconnode.GetStateRandaoRes, error) {
	panic("unimplemented")
}

// GetStateRoot implements beaconnode.Handler.
func (*MyServer) GetStateRoot(ctx context.Context, params beaconnode.GetStateRootParams) (beaconnode.GetStateRootRes, error) {
	panic("unimplemented")
}

// GetStateV2 implements beaconnode.Handler.
func (*MyServer) GetStateV2(ctx context.Context, params beaconnode.GetStateV2Params) (beaconnode.GetStateV2Res, error) {
	panic("unimplemented")
}

// GetStateValidator implements beaconnode.Handler.
func (*MyServer) GetStateValidator(ctx context.Context, params beaconnode.GetStateValidatorParams) (beaconnode.GetStateValidatorRes, error) {
	panic("unimplemented")
}

// GetStateValidatorBalances implements beaconnode.Handler.
func (*MyServer) GetStateValidatorBalances(ctx context.Context, params beaconnode.GetStateValidatorBalancesParams) (beaconnode.GetStateValidatorBalancesRes, error) {
	panic("unimplemented")
}

// GetStateValidators implements beaconnode.Handler.
func (*MyServer) GetStateValidators(ctx context.Context, params beaconnode.GetStateValidatorsParams) (beaconnode.GetStateValidatorsRes, error) {
	panic("unimplemented")
}

// GetSyncCommitteeDuties implements beaconnode.Handler.
func (*MyServer) GetSyncCommitteeDuties(ctx context.Context, req []string, params beaconnode.GetSyncCommitteeDutiesParams) (beaconnode.GetSyncCommitteeDutiesRes, error) {
	panic("unimplemented")
}

// GetSyncCommitteeRewards implements beaconnode.Handler.
func (*MyServer) GetSyncCommitteeRewards(ctx context.Context, req []string, params beaconnode.GetSyncCommitteeRewardsParams) (beaconnode.GetSyncCommitteeRewardsRes, error) {
	panic("unimplemented")
}

// GetSyncingStatus implements beaconnode.Handler.
func (*MyServer) GetSyncingStatus(ctx context.Context) (beaconnode.GetSyncingStatusRes, error) {
	panic("unimplemented")
}

// PostStateValidatorBalances implements beaconnode.Handler.
func (*MyServer) PostStateValidatorBalances(ctx context.Context, req []string, params beaconnode.PostStateValidatorBalancesParams) (beaconnode.PostStateValidatorBalancesRes, error) {
	panic("unimplemented")
}

// PrepareBeaconCommitteeSubnet implements beaconnode.Handler.
func (*MyServer) PrepareBeaconCommitteeSubnet(ctx context.Context, req []beaconnode.PrepareBeaconCommitteeSubnetReqItem) (beaconnode.PrepareBeaconCommitteeSubnetRes, error) {
	panic("unimplemented")
}

// PrepareBeaconProposer implements beaconnode.Handler.
func (*MyServer) PrepareBeaconProposer(ctx context.Context, req []beaconnode.PrepareBeaconProposerReqItem) (beaconnode.PrepareBeaconProposerRes, error) {
	panic("unimplemented")
}

// PrepareSyncCommitteeSubnets implements beaconnode.Handler.
func (*MyServer) PrepareSyncCommitteeSubnets(ctx context.Context, req []beaconnode.PrepareSyncCommitteeSubnetsReqItem) (beaconnode.PrepareSyncCommitteeSubnetsRes, error) {
	panic("unimplemented")
}

// ProduceAttestationData implements beaconnode.Handler.
func (*MyServer) ProduceAttestationData(ctx context.Context, params beaconnode.ProduceAttestationDataParams) (beaconnode.ProduceAttestationDataRes, error) {
	panic("unimplemented")
}

// ProduceBlindedBlock implements beaconnode.Handler.
func (*MyServer) ProduceBlindedBlock(ctx context.Context, params beaconnode.ProduceBlindedBlockParams) (beaconnode.ProduceBlindedBlockRes, error) {
	panic("unimplemented")
}

// ProduceBlockV2 implements beaconnode.Handler.
func (*MyServer) ProduceBlockV2(ctx context.Context, params beaconnode.ProduceBlockV2Params) (beaconnode.ProduceBlockV2Res, error) {
	panic("unimplemented")
}

// ProduceBlockV3 implements beaconnode.Handler.
func (*MyServer) ProduceBlockV3(ctx context.Context, params beaconnode.ProduceBlockV3Params) (beaconnode.ProduceBlockV3Res, error) {
	panic("unimplemented")
}

// ProduceSyncCommitteeContribution implements beaconnode.Handler.
func (*MyServer) ProduceSyncCommitteeContribution(ctx context.Context, params beaconnode.ProduceSyncCommitteeContributionParams) (beaconnode.ProduceSyncCommitteeContributionRes, error) {
	panic("unimplemented")
}

// PublishAggregateAndProofs implements beaconnode.Handler.
func (*MyServer) PublishAggregateAndProofs(ctx context.Context, req []beaconnode.PublishAggregateAndProofsReqItem) (beaconnode.PublishAggregateAndProofsRes, error) {
	panic("unimplemented")
}

// PublishBlindedBlock implements beaconnode.Handler.
func (*MyServer) PublishBlindedBlock(ctx context.Context, req beaconnode.PublishBlindedBlockReqApplicationOctetStream, params beaconnode.PublishBlindedBlockParams) (beaconnode.PublishBlindedBlockRes, error) {
	panic("unimplemented")
}

// PublishBlindedBlockV2 implements beaconnode.Handler.
func (*MyServer) PublishBlindedBlockV2(ctx context.Context, req beaconnode.PublishBlindedBlockV2ReqApplicationOctetStream, params beaconnode.PublishBlindedBlockV2Params) (beaconnode.PublishBlindedBlockV2Res, error) {
	panic("unimplemented")
}

// PublishBlock implements beaconnode.Handler.
func (*MyServer) PublishBlock(ctx context.Context, req beaconnode.PublishBlockReqApplicationOctetStream, params beaconnode.PublishBlockParams) (beaconnode.PublishBlockRes, error) {
	panic("unimplemented")
}

// PublishBlockV2 implements beaconnode.Handler.
func (*MyServer) PublishBlockV2(ctx context.Context, req beaconnode.PublishBlockV2ReqApplicationOctetStream, params beaconnode.PublishBlockV2Params) (beaconnode.PublishBlockV2Res, error) {
	panic("unimplemented")
}

// PublishContributionAndProofs implements beaconnode.Handler.
func (*MyServer) PublishContributionAndProofs(ctx context.Context, req []beaconnode.PublishContributionAndProofsReqItem) (beaconnode.PublishContributionAndProofsRes, error) {
	panic("unimplemented")
}

// RegisterValidator implements beaconnode.Handler.
func (*MyServer) RegisterValidator(ctx context.Context, req []beaconnode.RegisterValidatorReqItem) (beaconnode.RegisterValidatorRes, error) {
	panic("unimplemented")
}

// SubmitBeaconCommitteeSelections implements beaconnode.Handler.
func (*MyServer) SubmitBeaconCommitteeSelections(ctx context.Context, req []beaconnode.SubmitBeaconCommitteeSelectionsReqItem) (beaconnode.SubmitBeaconCommitteeSelectionsRes, error) {
	panic("unimplemented")
}

// SubmitPoolAttestations implements beaconnode.Handler.
func (*MyServer) SubmitPoolAttestations(ctx context.Context, req []beaconnode.SubmitPoolAttestationsReqItem) (beaconnode.SubmitPoolAttestationsRes, error) {
	panic("unimplemented")
}

// SubmitPoolAttesterSlashings implements beaconnode.Handler.
func (*MyServer) SubmitPoolAttesterSlashings(ctx context.Context, req *beaconnode.SubmitPoolAttesterSlashingsReq) (beaconnode.SubmitPoolAttesterSlashingsRes, error) {
	panic("unimplemented")
}

// SubmitPoolBLSToExecutionChange implements beaconnode.Handler.
func (*MyServer) SubmitPoolBLSToExecutionChange(ctx context.Context, req []beaconnode.SubmitPoolBLSToExecutionChangeReqItem) (beaconnode.SubmitPoolBLSToExecutionChangeRes, error) {
	panic("unimplemented")
}

// SubmitPoolProposerSlashings implements beaconnode.Handler.
func (*MyServer) SubmitPoolProposerSlashings(ctx context.Context, req *beaconnode.SubmitPoolProposerSlashingsReq) (beaconnode.SubmitPoolProposerSlashingsRes, error) {
	panic("unimplemented")
}

// SubmitPoolSyncCommitteeSignatures implements beaconnode.Handler.
func (*MyServer) SubmitPoolSyncCommitteeSignatures(ctx context.Context, req []beaconnode.SubmitPoolSyncCommitteeSignaturesReqItem) (beaconnode.SubmitPoolSyncCommitteeSignaturesRes, error) {
	panic("unimplemented")
}

// SubmitPoolVoluntaryExit implements beaconnode.Handler.
func (*MyServer) SubmitPoolVoluntaryExit(ctx context.Context, req *beaconnode.SubmitPoolVoluntaryExitReq) (beaconnode.SubmitPoolVoluntaryExitRes, error) {
	panic("unimplemented")
}

// SubmitSyncCommitteeSelections implements beaconnode.Handler.
func (*MyServer) SubmitSyncCommitteeSelections(ctx context.Context, req []beaconnode.SubmitSyncCommitteeSelectionsReqItem) (beaconnode.SubmitSyncCommitteeSelectionsRes, error) {
	panic("unimplemented")
}

var _ beaconnode.Handler = (*MyServer)(nil)

func main() {

}
