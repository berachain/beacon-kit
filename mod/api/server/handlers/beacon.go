package handlers

import (
	"net/http"
	types "github.com/berachain/beacon-kit/mod/api/server/types"
	echo "github.com/labstack/echo/v4"
)

func (rh RouteHandlers) GetGenesis(c echo.Context) (err error) {
	rh.Backend.GetGenesis()
	return echo.NewHTTPError(http.StatusNotFound, "Chain genesis info is not yet known")
}

func (rh RouteHandlers) GetStateHashTreeRoot(c echo.Context) (err error) {
	params, err := BindAndValidate[types.StateIdRequest](c)
	if err != nil {
		return err
	}
	rh.Backend.GetStateHashTreeRoot(params.StateId)
	return echo.NewHTTPError(http.StatusNotFound, "State not found")
}

func (rh RouteHandlers) GetStateFork(c echo.Context) (err error) {
	params, err := BindAndValidate[types.StateIdRequest](c)
	if err != nil {
		return err
	}
	rh.Backend.GetStateFork(params.StateId)
	return echo.NewHTTPError(http.StatusNotFound, "State not found")
}

func (rh RouteHandlers) GetStateFinalityCheckpoints(c echo.Context) (err error) {
	params, err := BindAndValidate[types.StateIdRequest](c)
	if err != nil {
		return err
	}
	rh.Backend.GetStateFinalityCheckpoints(params.StateId)
	return echo.NewHTTPError(http.StatusNotFound, "State not found")
}

func (rh RouteHandlers) GetStateValidators(c echo.Context) (err error) {
	params, err := BindAndValidate[types.StateValidatorsGetRequest](c)
	if err != nil {
		return err
	}
	rh.Backend.GetStateValidators(params.StateId, params.Id, params.Status)
	return echo.NewHTTPError(http.StatusNotFound, "State not found")
}

func (rh RouteHandlers) PostStateValidators(c echo.Context) (err error) {
	params, err := BindAndValidate[types.StateValidatorsPostRequest](c)
	if err != nil {
		return err
	}
	rh.Backend.GetStateValidators(params.StateId, params.Ids, params.Statuses)
	return echo.NewHTTPError(http.StatusNotFound, "State not found")
}

func (rh RouteHandlers) GetStateValidator(c echo.Context) (err error) {
	params, err := BindAndValidate[types.StateValidatorRequest](c)
	if err != nil {
		return err
	}
	rh.Backend.GetStateValidator(params.StateId, params.ValidatorId)
	return echo.NewHTTPError(http.StatusNotFound, "State not found")
}

func (rh RouteHandlers) GetStateValidatorBalances(c echo.Context) (err error) {
	return echo.NewHTTPError(http.StatusNotFound, "State not found")
}

func (rh RouteHandlers) PostStateValidatorBalances(c echo.Context) (err error) {
	return echo.NewHTTPError(http.StatusNotFound, "State not found")
}

func (rh RouteHandlers) GetStateCommittees(c echo.Context) (err error) {
	return echo.NewHTTPError(http.StatusNotFound, "State not found")
}

func (rh RouteHandlers) GetStateSyncCommittees(c echo.Context) (err error) {
	return echo.NewHTTPError(http.StatusNotFound, "State not found")
}

func (rh RouteHandlers) GetStateRandao(c echo.Context) (err error) {
	params, err := BindAndValidate[types.RandaoRequest](c)
	if err != nil {
		return err
	}
	rh.Backend.GetStateRandao(params.StateId, params.Epoch)
	return echo.NewHTTPError(http.StatusNotFound, "State not found")
}

func (rh RouteHandlers) GetBlockHeaders(c echo.Context) (err error) {
	params, err := BindAndValidate[types.BeaconHeadersRequest](c)
	if err != nil {
		return err
	}
	rh.Backend.GetBlockHeaders(params.Slot, params.ParentRoot)
	return c.String(http.StatusOK, "BlockHeaders")
}

func (rh RouteHandlers) GetBlockHeader(c echo.Context) (err error) {
	params, err := BindAndValidate[types.BlockIdRequest](c)
	if err != nil {
		return err
	}
	rh.Backend.GetBlockHeader(params.BlockId)
	return echo.NewHTTPError(http.StatusNotFound, "Block not found")
}

func (rh RouteHandlers) GetBlock(c echo.Context) (err error) {
	params, err := BindAndValidate[types.BlockIdRequest](c)
	if err != nil {
		return err
	}
	rh.Backend.GetBlock(params.BlockId)
	return echo.NewHTTPError(http.StatusNotFound, "Block not found")
}

func (rh RouteHandlers) GetBlockRoot(c echo.Context) (err error) {
	params, err := BindAndValidate[types.BlockIdRequest](c)
	if err != nil {
		return err
	}
	rh.Backend.GetBlockRoot(params.BlockId)
	return echo.NewHTTPError(http.StatusNotFound, "Block not found")
}

func (rh RouteHandlers) GetBlockAttestations(c echo.Context) (err error) {
	params, err := BindAndValidate[types.BlockIdRequest](c)
	if err != nil {
		return err
	}
	rh.Backend.GetBlockAttestations(params.BlockId)
	return echo.NewHTTPError(http.StatusNotFound, "Block not found")
}

func (rh RouteHandlers) GetBlobSidecars(c echo.Context) (err error) {
	params, err := BindAndValidate[types.BlobSidecarRequest](c)
	if err != nil {
		return err
	}
	rh.Backend.GetBlobSidecars(params.BlockId, params.Indices)
	return echo.NewHTTPError(http.StatusNotFound, "Block not found")
}

func (rh RouteHandlers) GetSyncCommiteeAwards(c echo.Context) (err error) {
	params, err := BindAndValidate[types.SyncComitteeAwardsRequest](c)
	if err != nil {
		return err
	}
	rh.Backend.GetBlobSidecars(params.BlockId, params.Ids)
	return echo.NewHTTPError(http.StatusNotFound, "Block not found")
}

func (rh RouteHandlers) GetDepositSnapshot(c echo.Context) (err error) {
	rh.Backend.GetDepositSnapshot()
	return echo.NewHTTPError(http.StatusNotFound, "No Finalized Snapshot Available")
}

func (rh RouteHandlers) GetBlockAwards(c echo.Context) (err error) {
	params, err := BindAndValidate[types.BlockIdRequest](c)
	if err != nil {
		return err
	}
	rh.Backend.GetBlockAwards(params.BlockId)
	return echo.NewHTTPError(http.StatusNotFound, "Block not found")
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
