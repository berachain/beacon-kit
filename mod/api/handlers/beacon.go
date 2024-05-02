package apihandler

import "net/http"

func (rh RouteHandler) GetGenesis(res http.ResponseWriter, req *http.Request) {
    res.Write([]byte("Genesis"))
}

func (rh RouteHandler) GetStateHashTreeRoot(res http.ResponseWriter, req *http.Request) {
    res.Write([]byte("StateHashTree"))
}

func (rh RouteHandler) GetStateFork(res http.ResponseWriter, req *http.Request) {
    res.Write([]byte("StateFork"))
}

func (rh RouteHandler) GetStateFinalityCheckpoints(res http.ResponseWriter, req *http.Request) {
    res.Write([]byte("StateFinalityCheckpoints"))
}

func (rh RouteHandler) GetStateValidators(res http.ResponseWriter, req *http.Request) { 
    res.Write([]byte("StateValidators"))
}   

func (rh RouteHandler) PostStateValidators(res http.ResponseWriter, req *http.Request) {
    res.Write([]byte("PostStateValidators"))
}

func (rh RouteHandler) GetStateValidator(res http.ResponseWriter, req *http.Request) {  
    res.Write([]byte("StateValidator"))
}

func (rh RouteHandler) GetStateValidatorBalances(res http.ResponseWriter, req *http.Request) {      
    res.Write([]byte("StateValidatorBalances"))
}

func (rh RouteHandler) PostStateValidatorBalances(res http.ResponseWriter, req *http.Request) {
    res.Write([]byte("PostStateValidatorBalances"))
}

func (rh RouteHandler) GetStateCommittees(res http.ResponseWriter, req *http.Request) {     
    res.Write([]byte("StateCommittees"))
}

func (rh RouteHandler) GetStateSyncCommittees(res http.ResponseWriter, req *http.Request) {      
    res.Write([]byte("StateSyncCommittees"))
}

func (rh RouteHandler) GetStateRandao(res http.ResponseWriter, req *http.Request) {      
    res.Write([]byte("StateRandao"))
}

func (rh RouteHandler) GetBlockHeaders(res http.ResponseWriter, req *http.Request) {      
    res.Write([]byte("BlockHeaders"))
}   

func (rh RouteHandler) GetBlockHeader(res http.ResponseWriter, req *http.Request) { 
    res.Write([]byte("BlockHeader"))
}   

func (rh RouteHandler) PostSignedBlindBlockV1(res http.ResponseWriter, req *http.Request) {      
    res.Write([]byte("SignedBlindBlockV1"))
}       

func (rh RouteHandler) PostSignedBlindBlockV2(res http.ResponseWriter, req *http.Request) {     
    res.Write([]byte("SignedBlindBlockV2"))
}   

func (rh RouteHandler) PostSignedBlockV1(res http.ResponseWriter, req *http.Request) {  
    res.Write([]byte("SignedBlockV1"))
}

func (rh RouteHandler) PostSignedBlockV2(res http.ResponseWriter, req *http.Request) {
    res.Write([]byte("SignedBlockV2"))
}

func (rh RouteHandler) GetBlock(res http.ResponseWriter, req *http.Request) {
    res.Write([]byte("Block"))
}

func (rh RouteHandler) GetBlockRoot(res http.ResponseWriter, req *http.Request) {
    res.Write([]byte("BlockRoot"))
}

func (rh RouteHandler) GetStateAttestations(res http.ResponseWriter, req *http.Request) {   
    res.Write([]byte("StateAttestations"))
}

func (rh RouteHandler) GetBlobSidecars(res http.ResponseWriter, req *http.Request) {
    res.Write([]byte("BlobSidecars"))
}

func (rh RouteHandler) GetSyncCommiteeAwards(res http.ResponseWriter, req *http.Request) {
    res.Write([]byte("SyncCommiteeAwards"))
}

func (rh RouteHandler) GetDepositSnapshot(res http.ResponseWriter, req *http.Request) {
    res.Write([]byte("DepositSnapshot"))
}

func (rh RouteHandler) GetBlockAwards(res http.ResponseWriter, req *http.Request) {
    res.Write([]byte("BlockAwards"))
}

func (rh RouteHandler) GetAttestationRewards(res http.ResponseWriter, req *http.Request) {
    res.Write([]byte("AttestationRewards"))
}

func (rh RouteHandler) GetBlindedBlock(res http.ResponseWriter, req *http.Request) {
    res.Write([]byte("BlindedBlock"))
}

func (rh RouteHandler) GetLightClientBootstrapForRoot(res http.ResponseWriter, req *http.Request) {
    res.Write([]byte("LightClientBootstrapForRoot"))
}

func (rh RouteHandler) GetLightClientUpdateForPeriod(res http.ResponseWriter, req *http.Request) {

    res.Write([]byte("LightClientUpdateForPeriod"))
}

func (rh RouteHandler) GetLightClientFinalityUpdate(res http.ResponseWriter, req *http.Request) {
    res.Write([]byte("LightClientFinalityUpdate"))
}

func (rh RouteHandler) GetLightClientOptimisticUpdate(res http.ResponseWriter, req *http.Request) {
    res.Write([]byte("LightClientOptimisticUpdate"))
}

func (rh RouteHandler) GetPoolAttestations(res http.ResponseWriter, req *http.Request) {

    res.Write([]byte("PoolAttestations"))
}

func (rh RouteHandler) PostPoolAttestations(res http.ResponseWriter, req *http.Request) {
    res.Write([]byte("PostPoolAttestations"))
}

func (rh RouteHandler) GetPoolAttesterSlashings(res http.ResponseWriter, req *http.Request) {
    res.Write([]byte("PoolAttesterSlashings"))
}

func (rh RouteHandler) PostPoolAttesterSlashings(res http.ResponseWriter, req *http.Request) {
    res.Write([]byte("PostPoolAttesterSlashings"))
}

func (rh RouteHandler) GetPoolProposerSlashings(res http.ResponseWriter, req *http.Request) {
    res.Write([]byte("PoolProposerSlashings"))
}

func (rh RouteHandler) PostPoolProposerSlashings(res http.ResponseWriter, req *http.Request) {

    res.Write([]byte("PostPoolProposerSlashings"))
}

func (rh RouteHandler) PostPoolSyncCommitteeSignature(res http.ResponseWriter, req *http.Request) {

    res.Write([]byte("PostPoolSyncCommitteeSignature"))
}

func (rh RouteHandler) GetPoolVoluntaryExits(res http.ResponseWriter, req *http.Request) {
    res.Write([]byte("PoolVoluntaryExits"))
}

func (rh RouteHandler) PostPoolVoluntaryExits(res http.ResponseWriter, req *http.Request) {
    res.Write([]byte("PostPoolVoluntaryExits"))
}

func (rh RouteHandler) GetPoolSignedBLSExecutionChanges(res http.ResponseWriter, req *http.Request) {
    res.Write([]byte("PoolSignedBLSExecutionChanges"))
}

func (rh RouteHandler) PostPoolSignedBLSExecutionChanges(res http.ResponseWriter, req *http.Request) {
    res.Write([]byte("PostPoolSignedBLSExecutionChanges"))
}
