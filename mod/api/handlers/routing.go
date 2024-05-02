package apihandler

import (
	"net/http"
)

// Router is an interface for a router
type Router interface {
	http.Handler
	Use(...func(http.Handler) http.Handler)
	Get(string, http.HandlerFunc)
	Post(string, http.HandlerFunc)
}

// TODO: thoughts - should these be more specific??? aka argument list should
// specify what to take instead of generic writer and request
// need some sort of middleware mapping?? to extract the params from req and pass to handler

// This handler interface specifies routes according to the
// v2.5.0 - Ethereum Proof-of-Stake Consensus Specification v1.4.0 spec.
type BackendRouteHandler interface {
	GetGenesis(res http.ResponseWriter, req *http.Request)
	GetStateHashTreeRoot(res http.ResponseWriter, req *http.Request)
	GetStateFork(res http.ResponseWriter, req *http.Request)
	GetStateFinalityCheckpoints(res http.ResponseWriter, req *http.Request)
	GetStateValidators(res http.ResponseWriter, req *http.Request)
	PostStateValidators(res http.ResponseWriter, req *http.Request)
	GetStateValidator(res http.ResponseWriter, req *http.Request)
	GetStateValidatorBalances(res http.ResponseWriter, req *http.Request)
	PostStateValidatorBalances(res http.ResponseWriter, req *http.Request)
	GetStateCommittees(res http.ResponseWriter, req *http.Request)
	GetStateSyncCommittees(res http.ResponseWriter, req *http.Request)
	GetStateRandao(res http.ResponseWriter, req *http.Request)
	GetBlockHeaders(res http.ResponseWriter, req *http.Request)
	GetBlockHeader(res http.ResponseWriter, req *http.Request)
	PostSignedBlindBlockV1(res http.ResponseWriter, req *http.Request)
	PostSignedBlindBlockV2(res http.ResponseWriter, req *http.Request)
	PostSignedBlockV1(res http.ResponseWriter, req *http.Request)
	PostSignedBlockV2(res http.ResponseWriter, req *http.Request)
	GetBlock(res http.ResponseWriter, req *http.Request)
	GetBlockRoot(res http.ResponseWriter, req *http.Request)
	GetStateAttestations(res http.ResponseWriter, req *http.Request)
	GetBlobSidecars(res http.ResponseWriter, req *http.Request)
	GetSyncCommiteeAwards(res http.ResponseWriter, req *http.Request)
	GetDepositSnapshot(res http.ResponseWriter, req *http.Request)
	GetBlockAwards(res http.ResponseWriter, req *http.Request)
	GetAttestationRewards(res http.ResponseWriter, req *http.Request)
	GetBlindedBlock(res http.ResponseWriter, req *http.Request)
	GetLightClientBootstrapForRoot(res http.ResponseWriter, req *http.Request)
	GetLightClientUpdateForPeriod(res http.ResponseWriter, req *http.Request)
	GetLightClientFinalityUpdate(res http.ResponseWriter, req *http.Request)
	GetLightClientOptimisticUpdate(res http.ResponseWriter, req *http.Request)
	GetPoolAttestations(res http.ResponseWriter, req *http.Request)
	PostPoolAttestations(res http.ResponseWriter, req *http.Request)
	GetPoolAttesterSlashings(res http.ResponseWriter, req *http.Request)
	PostPoolAttesterSlashings(res http.ResponseWriter, req *http.Request)
	GetPoolProposerSlashings(res http.ResponseWriter, req *http.Request)
	PostPoolProposerSlashings(res http.ResponseWriter, req *http.Request)
	PostPoolSyncCommitteeSignature(res http.ResponseWriter, req *http.Request)
	GetPoolVoluntaryExits(res http.ResponseWriter, req *http.Request)
	PostPoolVoluntaryExits(res http.ResponseWriter, req *http.Request)
	GetPoolSignedBLSExecutionChanges(res http.ResponseWriter, req *http.Request)
	PostPoolSignedBLSExecutionChanges(res http.ResponseWriter, req *http.Request)
}

func UseMiddlewares(r Router, middlewares []func(next http.Handler) http.Handler) {
	for _, middleware := range middlewares {
		r.Use(middleware)
	}
}

func AssignRoutes(r Router, handler BackendRouteHandler) {
	assignBeaconRoutes(r, handler)
}

func assignBeaconRoutes(r Router, handler BackendRouteHandler) {
	r.Get("/eth/v1/beacon/genesis", handler.GetGenesis)
	r.Get("/eth/v1/beacon/states/{state_id}/root", handler.GetStateHashTreeRoot)
	r.Get("/eth/v1/beacon/states/{state_id}/fork", handler.GetStateFork)
	r.Get("/eth/v1/beacon/states/{state_id}/finality_checkpoints", handler.GetStateFinalityCheckpoints)
	r.Get("/eth/v1/beacon/states/{state_id}/validators", handler.GetStateValidators)
	r.Post("/eth/v1/beacon/states/{state_id}/validators", handler.PostStateValidators)
	r.Get("/eth/v1/beacon/states/{state_id}/validators/{validator_id}", handler.GetStateValidator)
	r.Get("/eth/v1/beacon/states/{state_id}/validators/validator_balances", handler.GetStateValidatorBalances)
	r.Post("/eth/v1/beacon/states/{state_id}/validators/validator_balances", handler.PostStateValidatorBalances)
	r.Get("/eth/v1/beacon/states/{state_id}/committees", handler.GetStateCommittees)
	r.Get("/eth/v1/beacon/states/{state_id}/sync_committees", handler.GetStateSyncCommittees)
	r.Get("/eth/v1/beacon/states/{state_id}/randao", handler.GetStateRandao)
	r.Get("/eth/v1/beacon/headers", handler.GetBlockHeaders)
	r.Get("/eth/v1/beacon/headers/{block_id}", handler.GetBlockHeader)
	r.Post("/eth/v1/beacon/blocks/blinded_blocks", handler.PostSignedBlindBlockV1)
	r.Post("/eth/v2/beacon/blocks/blinded_blocks", handler.PostSignedBlindBlockV2)
	r.Post("/eth/v1/beacon/blocks", handler.PostSignedBlockV1)
	r.Post("/eth/v2/beacon/blocks", handler.PostSignedBlockV2)
	r.Get("/eth/v1/beacon/blocks/{block_id}", handler.GetBlock)
	r.Get("/eth/v1/beacon/blocks/{block_id}/root", handler.GetBlockRoot)
	r.Get("/eth/v1/beacon/blocks/{block_id}/attestations", handler.GetStateAttestations)
	r.Get("/eth/v1/beacon/blob_sidecars/{block_id}", handler.GetBlobSidecars)
	r.Get("/eth/v1/beacon/rewards/sync_committee/{block_id}", handler.GetSyncCommiteeAwards)
	r.Get("/eth/v1/beacon/deposit_snapshot", handler.GetDepositSnapshot)
	r.Get("/eth/v1/beacon/rewards/block/{block_id}", handler.GetBlockAwards)
	r.Post("/eth/v1/beacon/rewards/attestation/{block_id}", handler.GetAttestationRewards)
	r.Get("/eth/v1/beacon/blinded_blocks/{block_id}", handler.GetBlindedBlock)
	r.Get("/eth/v1/beacon/light_client/bootstrap/{block_root}", handler.GetLightClientBootstrapForRoot)
	r.Get("/eth/v1/beacon/light_client/updates", handler.GetLightClientUpdateForPeriod)
	r.Get("/eth/v1/beacon/light_client/finality_update", handler.GetLightClientFinalityUpdate)
	r.Get("/eth/v1/beacon/light_client/optimistic_update", handler.GetLightClientOptimisticUpdate)
	r.Get("/eth/v1/beacon/pool/attestations", handler.GetPoolAttestations)
	r.Post("/eth/v1/beacon/pool/attestations", handler.PostPoolAttestations)
	r.Get("/eth/v1/beacon/pool/attester_slashings", handler.GetPoolAttesterSlashings)
	r.Post("/eth/v1/beacon/pool/attester_slashings", handler.PostPoolAttesterSlashings)
	r.Get("/eth/v1/beacon/pool/proposer_slashings", handler.GetPoolProposerSlashings)
	r.Post("/eth/v1/beacon/pool/proposer_slashings", handler.PostPoolProposerSlashings)
	r.Post("/eth/v1/beacon/pool/sync_committes", handler.PostPoolSyncCommitteeSignature)
	r.Get("/eth/v1/beacon/pool/voluntary_exits", handler.GetPoolVoluntaryExits)
	r.Post("/eth/v1/beacon/pool/voluntary_exits", handler.PostPoolVoluntaryExits)
	r.Get("/eth/v1/beacon/pool/bls_to_execution_changes", handler.GetPoolSignedBLSExecutionChanges)
	r.Post("/eth/v1/beacon/pool/bls_to_execution_changes", handler.PostPoolSignedBLSExecutionChanges)
}
