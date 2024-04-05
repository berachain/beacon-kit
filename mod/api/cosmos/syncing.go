package cosmos

import (
	tendermintv1beta1 "cosmossdk.io/api/cosmos/base/tendermint/v1beta1"
	"github.com/berachain/beacon-kit/mod/api/beaconnode"
	cometabci "github.com/cometbft/cometbft/abci/types"
	"google.golang.org/protobuf/proto"
)

func (c ChainQuerier) GetSyncingStatus() beaconnode.GetSyncingStatusRes {
	ctx, err := c.ContextGetter(0, false)
	if err != nil {
		return &beaconnode.GetSyncingStatusInternalServerError{
			Code:        1,
			Message:     "Internal Server Error",
			Stacktraces: nil,
		}
	}

	req := cometabci.RequestQuery{
		Path: tendermintv1beta1.Service_GetSyncing_FullMethodName,
	}
	query, err := c.ABCI.Query(ctx, &req)
	if err != nil {
		return &beaconnode.GetSyncingStatusInternalServerError{
			Code:        1,
			Message:     "Internal Server Error",
			Stacktraces: nil,
		}
	}

	resp := tendermintv1beta1.GetSyncingResponse{}
	err = proto.Unmarshal(query.Value, &resp)
	if err != nil {
		return &beaconnode.GetSyncingStatusInternalServerError{
			Code:        1,
			Message:     "Internal Server Error",
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
