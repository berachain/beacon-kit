package evm

import (
	"context"

	v1alpha1 "github.com/berachain/beacon-kit/runtime/modules/beacon/api/v1alpha1"
	"github.com/berachain/beacon-kit/runtime/modules/beacon/keeper"
	"google.golang.org/protobuf/types/known/emptypb"
)

// APIServer is the server API for Query service.
type APIServer struct {
	v1alpha1.UnimplementedQueryServer
	*keeper.Keeper
}

// DepositQueueLength returns the length of the deposit queue.
func (api *APIServer) DepositQueueLength(
	ctx context.Context, _ *emptypb.Empty,
) (*v1alpha1.DepositQueueLengthResponse, error) {
	queueLength, err := api.Keeper.RawBeaconStore(ctx).LengthDeposits()
	if err != nil {
		return nil, err
	}
	return &v1alpha1.DepositQueueLengthResponse{Length: queueLength}, nil
}
