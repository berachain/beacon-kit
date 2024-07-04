package appmanager

import (
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/genesis"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	dastore "github.com/berachain/beacon-kit/mod/da/pkg/store"
	datypes "github.com/berachain/beacon-kit/mod/da/pkg/types"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/middleware"
)

type (
	// AvailabilityStore is a type alias for the availability store.
	AvailabilityStore = dastore.Store[*types.BeaconBlockBody]

	Genesis = genesis.Genesis[*types.Deposit, *types.ExecutionPayloadHeader]

	// ABCIMiddleware is a type alias for the ABCIMiddleware.
	ABCIMiddleware = middleware.ABCIMiddleware[
		*AvailabilityStore,
		*types.BeaconBlock,
		*datypes.BlobSidecars,
		*types.Deposit,
		*types.ExecutionPayload,
		*Genesis,
	]
)
