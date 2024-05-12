package handlers

import (
	"net/http"
	types "github.com/berachain/beacon-kit/mod/api/server/types"
	echo "github.com/labstack/echo/v4"
)

func (rh RouteHandlers) NotImplemented(c echo.Context) (err error) {
	return c.JSON(http.StatusNotImplemented, "Not Implemented")
}

func (rh RouteHandlers) GetGenesis(c echo.Context) (err error) {
	rh.Backend.GetGenesis()
	return c.String(http.StatusOK, "Genesis")
}

func (rh RouteHandlers) GetStateHashTreeRoot(c echo.Context) (err error) {
	params, err := BindAndValidate[types.StateIdRequest](c)
	if err != nil {
		return err
	}
	rh.Backend.GetStateHashTreeRoot(params.StateId)
	return c.String(http.StatusOK, "StateHashTree")
}

func (rh RouteHandlers) GetStateFork(c echo.Context) (err error) {
	params, err := BindAndValidate[types.StateIdRequest](c)
	if err != nil {
		return err
	}
	rh.Backend.GetStateFork(params.StateId)
	return c.String(http.StatusOK, "StateFork")
}

func (rh RouteHandlers) GetStateFinalityCheckpoints(c echo.Context) (err error) {
	params, err := BindAndValidate[types.StateIdRequest](c)
	if err != nil {
		return err
	}
	rh.Backend.GetStateFinalityCheckpoints(params.StateId)
	return c.String(http.StatusOK, "StateFinalityCheckpoints")
}

func (rh RouteHandlers) GetStateValidators(c echo.Context) (err error) {
	params, err := BindAndValidate[types.StateValidatorsGetRequest](c)
	if err != nil {
		return err
	}
	rh.Backend.GetStateValidators(params.StateId, params.Id, params.Status)
	return c.String(http.StatusOK, "StateValidators")
}

func (rh RouteHandlers) PostStateValidators(c echo.Context) (err error) {
	params, err := BindAndValidate[types.StateValidatorsPostRequest](c)
	if err != nil {
		return err
	}
	rh.Backend.GetStateValidators(params.StateId, params.Ids, params.Statuses)
	return c.String(http.StatusOK, "PostStateValidators")
}

func (rh RouteHandlers) GetStateValidator(c echo.Context) (err error) {
	params, err := BindAndValidate[types.StateValidatorRequest](c)
	if err != nil {
		return err
	}
	rh.Backend.GetStateValidator(params.StateId, params.ValidatorId)
	return c.String(http.StatusOK, "StateValidator")
}

func (rh RouteHandlers) GetStateValidatorBalances(c echo.Context) (err error) {
	return c.String(http.StatusNotImplemented, "StateValidatorBalances")
}

func (rh RouteHandlers) PostStateValidatorBalances(c echo.Context) (err error) {
	return c.String(http.StatusNotImplemented, "PostStateValidatorBalances")
}

func (rh RouteHandlers) GetStateCommittees(c echo.Context) (err error) {
	return c.String(http.StatusNotImplemented, "StateCommittees")
}

func (rh RouteHandlers) GetStateSyncCommittees(c echo.Context) (err error) {
	return c.String(http.StatusNotImplemented, "StateSyncCommittees")
}

func (rh RouteHandlers) GetStateRandao(c echo.Context) (err error) {
	return c.String(http.StatusNotImplemented, "StateRandao")
}

func (rh RouteHandlers) GetBlockHeaders(c echo.Context) (err error) {
	return c.String(http.StatusNotImplemented, "BlockHeaders")
}

func (rh RouteHandlers) GetBlockHeader(c echo.Context) (err error) {
	return c.String(http.StatusNotImplemented, "BlockHeader")
}

func (rh RouteHandlers) PostSignedBlindBlockV1(c echo.Context) (err error) {
	return c.String(http.StatusNotImplemented, "SignedBlindBlockV1")
}

func (rh RouteHandlers) PostSignedBlindBlockV2(c echo.Context) (err error) {
	return c.String(http.StatusNotImplemented, "SignedBlindBlockV2")
}

func (rh RouteHandlers) PostSignedBlockV1(c echo.Context) (err error) {
	return c.String(http.StatusNotImplemented, "SignedBlockV1")
}

func (rh RouteHandlers) PostSignedBlockV2(c echo.Context) (err error) {
	return c.String(http.StatusNotImplemented, "SignedBlockV2")
}

func (rh RouteHandlers) GetBlock(c echo.Context) (err error) {
	return c.String(http.StatusNotImplemented, "Block")
}

func (rh RouteHandlers) GetBlockRoot(c echo.Context) (err error) {
	return c.String(http.StatusNotImplemented, "BlockRoot")
}

func (rh RouteHandlers) GetStateAttestations(c echo.Context) (err error) {
	return c.String(http.StatusNotImplemented, "StateAttestations")
}

func (rh RouteHandlers) GetBlobSidecars(c echo.Context) (err error) {
	return c.String(http.StatusNotImplemented, "BlobSidecars")
}

func (rh RouteHandlers) GetSyncCommiteeAwards(c echo.Context) (err error) {
	return c.String(http.StatusNotImplemented, "SyncCommiteeAwards")
}

func (rh RouteHandlers) GetDepositSnapshot(c echo.Context) (err error) {
	return c.String(http.StatusNotImplemented, "DepositSnapshot")
}

func (rh RouteHandlers) GetBlockAwards(c echo.Context) (err error) {
	return c.String(http.StatusNotImplemented, "BlockAwards")
}

func (rh RouteHandlers) GetAttestationRewards(c echo.Context) (err error) {
	return c.String(http.StatusNotImplemented, "AttestationRewards")
}

func (rh RouteHandlers) GetBlindedBlock(c echo.Context) (err error) {
	return c.String(http.StatusNotImplemented, "BlindedBlock")
}

func (rh RouteHandlers) GetPoolAttestations(c echo.Context) (err error) {
	return c.String(http.StatusNotImplemented, "PoolAttestations")
}

func (rh RouteHandlers) PostPoolAttestations(c echo.Context) (err error) {
	return c.String(http.StatusNotImplemented, "PostPoolAttestations")
}

func (rh RouteHandlers) GetPoolAttesterSlashings(c echo.Context) (err error) {
	return c.String(http.StatusNotImplemented, "PoolAttesterSlashings")
}

func (rh RouteHandlers) PostPoolAttesterSlashings(c echo.Context) (err error) {
	return c.String(http.StatusNotImplemented, "PostPoolAttesterSlashings")
}

func (rh RouteHandlers) GetPoolProposerSlashings(c echo.Context) (err error) {
	return c.String(http.StatusNotImplemented, "PoolProposerSlashings")
}

func (rh RouteHandlers) PostPoolProposerSlashings(c echo.Context) (err error) {
	return c.String(http.StatusNotImplemented, "PostPoolProposerSlashings")
}

func (rh RouteHandlers) PostPoolSyncCommitteeSignature(c echo.Context) (err error) {
	return c.String(http.StatusNotImplemented, "PostPoolSyncCommitteeSignature")
}

func (rh RouteHandlers) GetPoolVoluntaryExits(c echo.Context) (err error) {
	return c.String(http.StatusNotImplemented, "PoolVoluntaryExits")
}

func (rh RouteHandlers) PostPoolVoluntaryExits(c echo.Context) (err error) {
	return c.String(http.StatusNotImplemented, "PostPoolVoluntaryExits")
}

func (rh RouteHandlers) GetPoolSignedBLSExecutionChanges(c echo.Context) (err error) {
	return c.String(http.StatusNotImplemented, "PoolSignedBLSExecutionChanges")
}

func (rh RouteHandlers) PostPoolSignedBLSExecutionChanges(c echo.Context) (err error) {
	return c.String(http.StatusNotImplemented, "PostPoolSignedBLSExecutionChanges")
}
