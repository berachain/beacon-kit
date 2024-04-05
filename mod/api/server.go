package rpc

import (
	"context"
	"fmt"
	"github.com/berachain/beacon-kit/mod/api/beaconnode"
	"github.com/berachain/beacon-kit/mod/node-builder/service"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"strconv"
)

type ChainQuerier interface {
	GetSyncingStatus() beaconnode.GetSyncingStatusRes
}

type Server struct {
	ContextGetter func(height int64, prove bool) (sdk.Context, error)
	Service       service.BeaconStorageBackend

	ChainQuerier ChainQuerier
}

func (s Server) Eventstream(ctx context.Context, params beaconnode.EventstreamParams) (beaconnode.EventstreamRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetAggregatedAttestation(ctx context.Context, params beaconnode.GetAggregatedAttestationParams) (beaconnode.GetAggregatedAttestationRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetAttestationsRewards(ctx context.Context, req []string, params beaconnode.GetAttestationsRewardsParams) (beaconnode.GetAttestationsRewardsRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetAttesterDuties(ctx context.Context, req []string, params beaconnode.GetAttesterDutiesParams) (beaconnode.GetAttesterDutiesRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetBlindedBlock(ctx context.Context, params beaconnode.GetBlindedBlockParams) (beaconnode.GetBlindedBlockRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetBlobSidecars(ctx context.Context, params beaconnode.GetBlobSidecarsParams) (beaconnode.GetBlobSidecarsRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetBlockAttestations(ctx context.Context, params beaconnode.GetBlockAttestationsParams) (beaconnode.GetBlockAttestationsRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetBlockHeader(ctx context.Context, params beaconnode.GetBlockHeaderParams) (beaconnode.GetBlockHeaderRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetBlockHeaders(ctx context.Context, params beaconnode.GetBlockHeadersParams) (beaconnode.GetBlockHeadersRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetBlockRewards(ctx context.Context, params beaconnode.GetBlockRewardsParams) (beaconnode.GetBlockRewardsRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetBlockRoot(ctx context.Context, params beaconnode.GetBlockRootParams) (beaconnode.GetBlockRootRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetBlockV2(ctx context.Context, params beaconnode.GetBlockV2Params) (beaconnode.GetBlockV2Res, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetDebugChainHeadsV2(ctx context.Context) (beaconnode.GetDebugChainHeadsV2Res, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetDebugForkChoice(ctx context.Context) (beaconnode.GetDebugForkChoiceRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetDepositContract(ctx context.Context) (beaconnode.GetDepositContractRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetDepositSnapshot(ctx context.Context) (beaconnode.GetDepositSnapshotRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetEpochCommittees(ctx context.Context, params beaconnode.GetEpochCommitteesParams) (beaconnode.GetEpochCommitteesRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetEpochSyncCommittees(ctx context.Context, params beaconnode.GetEpochSyncCommitteesParams) (beaconnode.GetEpochSyncCommitteesRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetForkSchedule(ctx context.Context) (beaconnode.GetForkScheduleRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetGenesis(ctx context.Context) (beaconnode.GetGenesisRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetHealth(ctx context.Context, params beaconnode.GetHealthParams) (beaconnode.GetHealthRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetLightClientBootstrap(ctx context.Context, params beaconnode.GetLightClientBootstrapParams) (beaconnode.GetLightClientBootstrapRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetLightClientFinalityUpdate(ctx context.Context) (beaconnode.GetLightClientFinalityUpdateRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetLightClientOptimisticUpdate(ctx context.Context) (beaconnode.GetLightClientOptimisticUpdateRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetLightClientUpdatesByRange(ctx context.Context, params beaconnode.GetLightClientUpdatesByRangeParams) (beaconnode.GetLightClientUpdatesByRangeRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetLiveness(ctx context.Context, req []string, params beaconnode.GetLivenessParams) (beaconnode.GetLivenessRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetNetworkIdentity(ctx context.Context) (beaconnode.GetNetworkIdentityRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetNextWithdrawals(ctx context.Context, params beaconnode.GetNextWithdrawalsParams) (beaconnode.GetNextWithdrawalsRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetNodeVersion(ctx context.Context) (beaconnode.GetNodeVersionRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetPeer(ctx context.Context, params beaconnode.GetPeerParams) (beaconnode.GetPeerRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetPeerCount(ctx context.Context) (beaconnode.GetPeerCountRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetPeers(ctx context.Context, params beaconnode.GetPeersParams) (beaconnode.GetPeersRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetPoolAttestations(ctx context.Context, params beaconnode.GetPoolAttestationsParams) (beaconnode.GetPoolAttestationsRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetPoolAttesterSlashings(ctx context.Context) (beaconnode.GetPoolAttesterSlashingsRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetPoolBLSToExecutionChanges(ctx context.Context) (beaconnode.GetPoolBLSToExecutionChangesRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetPoolProposerSlashings(ctx context.Context) (beaconnode.GetPoolProposerSlashingsRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetPoolVoluntaryExits(ctx context.Context) (beaconnode.GetPoolVoluntaryExitsRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetProposerDuties(ctx context.Context, params beaconnode.GetProposerDutiesParams) (beaconnode.GetProposerDutiesRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetSpec(ctx context.Context) (beaconnode.GetSpecRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetStateFinalityCheckpoints(ctx context.Context, params beaconnode.GetStateFinalityCheckpointsParams) (beaconnode.GetStateFinalityCheckpointsRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetStateFork(ctx context.Context, params beaconnode.GetStateForkParams) (beaconnode.GetStateForkRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetStateRandao(_ context.Context, params beaconnode.GetStateRandaoParams) (beaconnode.GetStateRandaoRes, error) {
	stateId := params.StateID
	if stateId == "" {
		return nil, fmt.Errorf("state_id is required in URL params")
	}

	stateIdAsInt, err := strconv.ParseUint(stateId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("state_id must be a number")
	}

	ctx, err := s.ContextGetter(int64(stateIdAsInt), false)
	if err != nil {
		return nil, err
	}

	randao, err := s.Service.BeaconState(ctx).GetRandaoMixAtIndex(stateIdAsInt)
	if err != nil {
		return nil, err
	}

	resp := &beaconnode.GetStateRandaoOK{
		ExecutionOptimistic: false,
		Finalized:           true,
		Data: beaconnode.GetStateRandaoOKData{
			Randao: hexutil.Encode(randao[:]),
		},
	}

	return resp, nil
}

func (s Server) GetStateRoot(ctx context.Context, params beaconnode.GetStateRootParams) (beaconnode.GetStateRootRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetStateV2(ctx context.Context, params beaconnode.GetStateV2Params) (beaconnode.GetStateV2Res, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetStateValidator(ctx context.Context, params beaconnode.GetStateValidatorParams) (beaconnode.GetStateValidatorRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetStateValidatorBalances(ctx context.Context, params beaconnode.GetStateValidatorBalancesParams) (beaconnode.GetStateValidatorBalancesRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetStateValidators(ctx context.Context, params beaconnode.GetStateValidatorsParams) (beaconnode.GetStateValidatorsRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetSyncCommitteeDuties(ctx context.Context, req []string, params beaconnode.GetSyncCommitteeDutiesParams) (beaconnode.GetSyncCommitteeDutiesRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetSyncCommitteeRewards(ctx context.Context, req []string, params beaconnode.GetSyncCommitteeRewardsParams) (beaconnode.GetSyncCommitteeRewardsRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetSyncingStatus(ctx context.Context) (beaconnode.GetSyncingStatusRes, error) {
	return s.ChainQuerier.GetSyncingStatus(), nil
}

func (s Server) PostStateValidatorBalances(ctx context.Context, req []string, params beaconnode.PostStateValidatorBalancesParams) (beaconnode.PostStateValidatorBalancesRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) PrepareBeaconCommitteeSubnet(ctx context.Context, req []beaconnode.PrepareBeaconCommitteeSubnetReqItem) (beaconnode.PrepareBeaconCommitteeSubnetRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) PrepareBeaconProposer(ctx context.Context, req []beaconnode.PrepareBeaconProposerReqItem) (beaconnode.PrepareBeaconProposerRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) PrepareSyncCommitteeSubnets(ctx context.Context, req []beaconnode.PrepareSyncCommitteeSubnetsReqItem) (beaconnode.PrepareSyncCommitteeSubnetsRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) ProduceAttestationData(ctx context.Context, params beaconnode.ProduceAttestationDataParams) (beaconnode.ProduceAttestationDataRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) ProduceBlindedBlock(ctx context.Context, params beaconnode.ProduceBlindedBlockParams) (beaconnode.ProduceBlindedBlockRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) ProduceBlockV2(ctx context.Context, params beaconnode.ProduceBlockV2Params) (beaconnode.ProduceBlockV2Res, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) ProduceBlockV3(ctx context.Context, params beaconnode.ProduceBlockV3Params) (beaconnode.ProduceBlockV3Res, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) ProduceSyncCommitteeContribution(ctx context.Context, params beaconnode.ProduceSyncCommitteeContributionParams) (beaconnode.ProduceSyncCommitteeContributionRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) PublishAggregateAndProofs(ctx context.Context, req []beaconnode.PublishAggregateAndProofsReqItem) (beaconnode.PublishAggregateAndProofsRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) PublishBlindedBlock(ctx context.Context, req beaconnode.PublishBlindedBlockReqApplicationOctetStream, params beaconnode.PublishBlindedBlockParams) (beaconnode.PublishBlindedBlockRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) PublishBlindedBlockV2(ctx context.Context, req beaconnode.PublishBlindedBlockV2ReqApplicationOctetStream, params beaconnode.PublishBlindedBlockV2Params) (beaconnode.PublishBlindedBlockV2Res, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) PublishBlock(ctx context.Context, req beaconnode.PublishBlockReqApplicationOctetStream, params beaconnode.PublishBlockParams) (beaconnode.PublishBlockRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) PublishBlockV2(ctx context.Context, req beaconnode.PublishBlockV2ReqApplicationOctetStream, params beaconnode.PublishBlockV2Params) (beaconnode.PublishBlockV2Res, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) PublishContributionAndProofs(ctx context.Context, req []beaconnode.PublishContributionAndProofsReqItem) (beaconnode.PublishContributionAndProofsRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) RegisterValidator(ctx context.Context, req []beaconnode.RegisterValidatorReqItem) (beaconnode.RegisterValidatorRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) SubmitBeaconCommitteeSelections(ctx context.Context, req []beaconnode.SubmitBeaconCommitteeSelectionsReqItem) (beaconnode.SubmitBeaconCommitteeSelectionsRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) SubmitPoolAttestations(ctx context.Context, req []beaconnode.SubmitPoolAttestationsReqItem) (beaconnode.SubmitPoolAttestationsRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) SubmitPoolAttesterSlashings(ctx context.Context, req *beaconnode.SubmitPoolAttesterSlashingsReq) (beaconnode.SubmitPoolAttesterSlashingsRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) SubmitPoolBLSToExecutionChange(ctx context.Context, req []beaconnode.SubmitPoolBLSToExecutionChangeReqItem) (beaconnode.SubmitPoolBLSToExecutionChangeRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) SubmitPoolProposerSlashings(ctx context.Context, req *beaconnode.SubmitPoolProposerSlashingsReq) (beaconnode.SubmitPoolProposerSlashingsRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) SubmitPoolSyncCommitteeSignatures(ctx context.Context, req []beaconnode.SubmitPoolSyncCommitteeSignaturesReqItem) (beaconnode.SubmitPoolSyncCommitteeSignaturesRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) SubmitPoolVoluntaryExit(ctx context.Context, req *beaconnode.SubmitPoolVoluntaryExitReq) (beaconnode.SubmitPoolVoluntaryExitRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) SubmitSyncCommitteeSelections(ctx context.Context, req []beaconnode.SubmitSyncCommitteeSelectionsReqItem) (beaconnode.SubmitSyncCommitteeSelectionsRes, error) {
	//TODO implement me
	panic("implement me")
}
