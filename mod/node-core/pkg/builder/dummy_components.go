package builder

import (
	consensustypes "github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	dastore "github.com/berachain/beacon-kit/mod/da/pkg/store"
	datypes "github.com/berachain/beacon-kit/mod/da/pkg/types"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/storage"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/runtime"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/runtime/middleware"
	depositdb "github.com/berachain/beacon-kit/mod/storage/pkg/deposit"
)

// dummy components for the builder. these are needed to support the
// direct dependencies needed to supply the Module
var (
	// MIDDLEWARES
	emptyFinalizeBlockMiddlware = &middleware.FinalizeBlockMiddleware[
		*consensustypes.BeaconBlock,
		runtime.BeaconState,
		*datypes.BlobSidecars,
	]{}
	emptyValidatorMiddleware = &middleware.ValidatorMiddleware[
		*dastore.Store[*consensustypes.BeaconBlockBody],
		*consensustypes.BeaconBlock,
		*consensustypes.BeaconBlockBody,
		runtime.BeaconState,
		*datypes.BlobSidecars,
		runtime.Backend,
	]{}
	// STORAGE BACKEND
	emptyStorageBackend = &storage.Backend[
		*dastore.Store[*consensustypes.BeaconBlockBody],
		*consensustypes.BeaconBlock,
		*consensustypes.BeaconBlockBody,
		components.BeaconState,
		*depositdb.KVStore[*consensustypes.Deposit],
	]{}
)
