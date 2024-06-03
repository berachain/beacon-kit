package components

import (
	"fmt"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	dastore "github.com/berachain/beacon-kit/mod/da/pkg/store"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/components/storage"
	primitives "github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
	depositdb "github.com/berachain/beacon-kit/mod/storage/pkg/deposit"
)

// ProvideStorageBackend
func ProvideStorageBackend(
	avs *dastore.Store[types.BeaconBlockBody],
	depositStore *depositdb.KVStore[*types.Deposit],
	chainSpec primitives.ChainSpec,
) (*storage.Backend[
	*dastore.Store[types.BeaconBlockBody],
	types.BeaconBlockBody,
	core.BeaconState[*types.BeaconBlockHeader,
		*types.ExecutionPayloadHeader, *types.Fork, *types.Validator, *engineprimitives.Withdrawal],
	*depositdb.KVStore[*types.Deposit],
], error) {
	fmt.Println("NEW STORAGE BACKEND")
	storageBackend := storage.NewBackend[
		*dastore.Store[types.BeaconBlockBody],
		types.BeaconBlockBody,
		core.BeaconState[
			*types.BeaconBlockHeader, *types.ExecutionPayloadHeader, *types.Fork,
			*types.Validator, *engineprimitives.Withdrawal,
		],
	](
		chainSpec,
		avs,
		nil,
		depositStore,
	)
	return storageBackend, nil
}
