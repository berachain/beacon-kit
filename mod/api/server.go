package rpc

import (
	"context"
	"github.com/berachain/beacon-kit/mod/api/beaconnode"
)

type ChainQuerier interface {
	GetSyncingStatus() beaconnode.GetSyncingStatusRes
	GetStateRandao(params beaconnode.GetStateRandaoParams) beaconnode.GetStateRandaoRes
	GetNodeVersion(ctx context.Context) beaconnode.GetNodeVersionRes
}

type Server struct {
	beaconnode.UnimplementedHandler

	ChainQuerier ChainQuerier
}

func (s Server) GetStateRandao(_ context.Context, params beaconnode.GetStateRandaoParams) (beaconnode.GetStateRandaoRes, error) {
	return s.ChainQuerier.GetStateRandao(params), nil
}

func (s Server) GetSyncingStatus(_ context.Context) (beaconnode.GetSyncingStatusRes, error) {
	return s.ChainQuerier.GetSyncingStatus(), nil
}

func (s Server) GetNodeVersion(ctx context.Context) (beaconnode.GetNodeVersionRes, error) {
	return s.ChainQuerier.GetNodeVersion(ctx), nil
}
