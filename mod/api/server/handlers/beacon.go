package handlers

import (
	"net/http"
	types "github.com/berachain/beacon-kit/mod/api/server/types"
	echo "github.com/labstack/echo/v4"
)

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
	return c.String(http.StatusOK, "StateValidatorBalances")
}

func (rh RouteHandlers) PostStateValidatorBalances(c echo.Context) (err error) {
	return c.String(http.StatusOK, "PostStateValidatorBalances")
}

func (rh RouteHandlers) GetStateCommittees(c echo.Context) (err error) {
	return c.String(http.StatusOK, "StateCommittees")
}

func (rh RouteHandlers) GetStateSyncCommittees(c echo.Context) (err error) {
	return c.String(http.StatusOK, "StateSyncCommittees")
}

func (rh RouteHandlers) GetStateRandao(c echo.Context) (err error) {
	params, err := BindAndValidate[types.RandaoRequest](c)
	if err != nil {
		return err
	}
	rh.Backend.GetStateRandao(params.StateId, params.Epoch)
	return c.String(http.StatusOK, "StateRandao")
}

func (rh RouteHandlers) GetBlockHeaders(c echo.Context) (err error) {
	params, err := BindAndValidate[types.BeaconHeadersRequest](c)
	if err != nil {
		return err
	}
	rh.Backend.GetStateRandao(params.Slot, params.ParentRoot)
	return c.String(http.StatusOK, "BlockHeaders")
}

func (rh RouteHandlers) GetBlockHeader(c echo.Context) (err error) {
	params, err := BindAndValidate[types.BlockIdRequest](c)
	if err != nil {
		return err
	}
	rh.Backend.GetBlockHeader(params.BlockId)
	return c.String(http.StatusOK, "BlockHeaders")
}

func (rh RouteHandlers) GetBlock(c echo.Context) (err error) {
	params, err := BindAndValidate[types.BlockIdRequest](c)
	if err != nil {
		return err
	}
	rh.Backend.GetBlock(params.BlockId)
	return c.String(http.StatusOK, "Block")
}

func (rh RouteHandlers) GetBlockRoot(c echo.Context) (err error) {
	params, err := BindAndValidate[types.BlockIdRequest](c)
	if err != nil {
		return err
	}
	rh.Backend.GetBlockRoot(params.BlockId)
	return c.String(http.StatusOK, "BlockRoot")
}

func (rh RouteHandlers) GetBlockAttestations(c echo.Context) (err error) {
	params, err := BindAndValidate[types.BlockIdRequest](c)
	if err != nil {
		return err
	}
	rh.Backend.GetBlockAttestations(params.BlockId)
	return c.String(http.StatusOK, "BlockAttestations")
}

func (rh RouteHandlers) GetBlobSidecars(c echo.Context) (err error) {
	params, err := BindAndValidate[types.BlobSidecarRequest](c)
	if err != nil {
		return err
	}
	rh.Backend.GetBlobSidecars(params.BlockId, params.Indices)
	return c.String(http.StatusOK, "GetBlobSidecars")
}

func (rh RouteHandlers) GetSyncCommiteeAwards(c echo.Context) (err error) {
	params, err := BindAndValidate[types.SyncComitteeAwardsRequest](c)
	if err != nil {
		return err
	}
	rh.Backend.GetBlobSidecars(params.BlockId, params.Ids)
	return c.String(http.StatusOK, "GetBlobSidecars")
}

func (rh RouteHandlers) GetDepositSnapshot(c echo.Context) (err error) {
	rh.Backend.GetDepositSnapshot()
	return c.String(http.StatusOK, "GetDepositSnapshot")
}

func (rh RouteHandlers) GetBlockAwards(c echo.Context) (err error) {
	params, err := BindAndValidate[types.BlockIdRequest](c)
	if err != nil {
		return err
	}
	rh.Backend.GetBlockAwards(params.BlockId)
	return c.String(http.StatusOK, "GetBlockAwards")
}

func (rh RouteHandlers) GetAttestationRewards(c echo.Context) (err error) {
	params, err := BindAndValidate[types.GetAttestationRewardsRequest](c)
	if err != nil {
		return err
	}
	rh.Backend.GetAttestationRewards(params.Epoch, params.Ids)
	return c.String(http.StatusOK, "GetAttestationRewards")
}

func (rh RouteHandlers) GetBlindedBlock(c echo.Context) (err error) {
	params, err := BindAndValidate[types.BlockIdRequest](c)
	if err != nil {
		return err
	}
	rh.Backend.GetBlindedBlock(params.BlockId)
	return c.String(http.StatusOK, "GetBlindedBlock")
}

func (rh RouteHandlers) GetPoolAttestations(c echo.Context) (err error) {
	params, err := BindAndValidate[types.GetPoolAttestationRequest](c)
	if err != nil {
		return err
	}
	rh.Backend.GetPoolAttestations(params.Slot, params.CommitteeIndex)
	return c.String(http.StatusOK, "GetPoolAttestations")
}

func (rh RouteHandlers) PostPoolAttestations(c echo.Context) (err error) {
	rh.Backend.PostPoolAttestations("")
	return c.String(http.StatusOK, "PostPoolAttestations")
}

func (rh RouteHandlers) PostPoolSyncCommitteeSignature(c echo.Context) (err error) {
	rh.Backend.PostPoolSyncCommitteeSignature()
	return c.String(http.StatusOK, "PostPoolSyncCommitteeSignature")
}

func (rh RouteHandlers) GetPoolVoluntaryExits(c echo.Context) (err error) {
	rh.Backend.GetPoolVoluntaryExits()
	return c.String(http.StatusOK, "PoolVoluntaryExits")
}

func (rh RouteHandlers) PostPoolVoluntaryExits(c echo.Context) (err error) {
	rh.Backend.PostPoolVoluntaryExits()
	return c.String(http.StatusOK, "PostPoolVoluntaryExits")
}

func (rh RouteHandlers) GetPoolSignedBLSExecutionChanges(c echo.Context) (err error) {
	rh.Backend.GetPoolSignedBLSExecutionChanges()
	return c.String(http.StatusOK, "PoolSignedBLSExecutionChanges")
}

func (rh RouteHandlers) PostPoolSignedBLSExecutionChanges(c echo.Context) (err error) {
	rh.Backend.PostPoolSignedBLSExecutionChanges()
	return c.String(http.StatusOK, "PostPoolSignedBLSExecutionChanges")
}
