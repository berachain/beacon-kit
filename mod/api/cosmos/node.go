package cosmos

import (
	"context"

	tendermintv1beta1 "cosmossdk.io/api/cosmos/base/tendermint/v1beta1"
	"google.golang.org/protobuf/proto"

	"github.com/berachain/beacon-kit/mod/api/beaconnode"
)

func (c ChainQuerier) GetSyncingStatus() beaconnode.GetSyncingStatusRes {
	ctx, err := c.ContextGetter(0, false)
	if err != nil {
		return &beaconnode.GetSyncingStatusInternalServerError{
			Code:        1,
			Message:     err.Error(),
			Stacktraces: nil,
		}
	}

	responseQuery, err := c.callABCIQuery(ctx, tendermintv1beta1.Service_GetSyncing_FullMethodName)
	if err != nil {
		return &beaconnode.GetSyncingStatusInternalServerError{
			Code:        2,
			Message:     err.Error(),
			Stacktraces: nil,
		}
	}

	resp := tendermintv1beta1.GetSyncingResponse{}
	err = proto.Unmarshal(responseQuery.Value, &resp)
	if err != nil {
		return &beaconnode.GetSyncingStatusInternalServerError{
			Code:        3,
			Message:     err.Error(),
			Stacktraces: nil,
		}
	}

	return &beaconnode.GetSyncingStatusOK{
		Data: beaconnode.GetSyncingStatusOKData{
			HeadSlot:     "",
			SyncDistance: "",
			IsSyncing:    resp.Syncing,
			IsOptimistic: false,
			ElOffline:    false,
		},
	}
}

func (c ChainQuerier) GetNodeVersion(_ context.Context) beaconnode.GetNodeVersionRes {
	ctx, err := c.ContextGetter(0, false)
	if err != nil {
		return &beaconnode.GetNodeVersionInternalServerError{
			Code:        1,
			Message:     err.Error(),
			Stacktraces: nil,
		}
	}

	responseQuery, err := c.callABCIQuery(ctx, tendermintv1beta1.Service_GetNodeInfo_FullMethodName)
	if err != nil {
		return &beaconnode.GetNodeVersionInternalServerError{
			Code:        2,
			Message:     err.Error(),
			Stacktraces: nil,
		}
	}

	resp := tendermintv1beta1.GetNodeInfoResponse{}
	err = proto.Unmarshal(responseQuery.Value, &resp)
	if err != nil {
		return &beaconnode.GetNodeVersionInternalServerError{
			Code:        3,
			Message:     err.Error(),
			Stacktraces: nil,
		}
	}

	return &beaconnode.GetNodeVersionOK{
		Data: beaconnode.GetNodeVersionOKData{
			Version: resp.ApplicationVersion.Version,
		},
	}
}
