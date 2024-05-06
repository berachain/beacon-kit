package staking

import (
	"context"

	"github.com/berachain/beacon-kit/mod/core/state"
)

type BeaconStorageBackend interface {
	BeaconState(context.Context) state.BeaconState
}
