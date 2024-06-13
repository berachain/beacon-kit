package components

import (
	"github.com/berachain/beacon-kit/mod/beacon/blockchain"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	dastore "github.com/berachain/beacon-kit/mod/da/pkg/store"
	datypes "github.com/berachain/beacon-kit/mod/da/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/runtime"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
	depositdb "github.com/berachain/beacon-kit/mod/storage/pkg/deposit"
)

// BeaconState is a type alias for the BeaconState.
type BeaconState = core.BeaconState[
	*types.BeaconBlockHeader, *types.Eth1Data,
	*types.ExecutionPayloadHeaderDeneb, *types.Fork,
	*types.Validator, *engineprimitives.Withdrawal,
]

// BeaconKitRuntime is a type alias for the BeaconKitRuntime.
type BeaconKitRuntime = runtime.BeaconKitRuntime[
	*dastore.Store[*types.BeaconBlockBody[*types.ExecutionPayload]],
	*types.BeaconBlock[*types.ExecutionPayload],
	*types.BeaconBlockBody[*types.ExecutionPayload],
	BeaconState,
	*datypes.BlobSidecars,
	*depositdb.KVStore[*types.Deposit],
	*types.ExecutionPayloadHeaderDeneb,
	blockchain.StorageBackend[
		*dastore.Store[*types.BeaconBlockBody[*types.ExecutionPayload]],
		*types.BeaconBlockBody[*types.ExecutionPayload],
		BeaconState,
		*datypes.BlobSidecars,
		*types.Deposit,
		*depositdb.KVStore[*types.Deposit],
	],
]
