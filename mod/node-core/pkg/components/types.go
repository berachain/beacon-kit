package components

import (
	"github.com/berachain/beacon-kit/mod/beacon/blockchain"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	dastore "github.com/berachain/beacon-kit/mod/da/pkg/store"
	datypes "github.com/berachain/beacon-kit/mod/da/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/runtime/middleware"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
	depositdb "github.com/berachain/beacon-kit/mod/storage/pkg/deposit"
)

type (
	BeaconState = core.BeaconState[
		*types.BeaconBlockHeader, *types.Eth1Data,
		*types.ExecutionPayloadHeader, *types.Fork,
		*types.Validator, *engineprimitives.Withdrawal,
	]

	Backend = blockchain.StorageBackend[
		*dastore.Store[*types.BeaconBlockBody],
		*types.BeaconBlockBody,
		BeaconState,
		*datypes.BlobSidecars,
		*types.Deposit,
		*depositdb.KVStore[*types.Deposit],
	]

	FinalizeBlockMiddleware = *middleware.FinalizeBlockMiddleware[
		*types.BeaconBlock, BeaconState, *datypes.BlobSidecars,
	]

	ValidatorMiddleware = *middleware.ValidatorMiddleware[
		*dastore.Store[*types.BeaconBlockBody],
		*types.BeaconBlock,
		*types.BeaconBlockBody,
		BeaconState,
		*datypes.BlobSidecars,
		Backend,
	]
)
