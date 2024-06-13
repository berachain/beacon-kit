package components

import (
	"github.com/berachain/beacon-kit/mod/async/pkg/event"
	"github.com/berachain/beacon-kit/mod/beacon/blockchain"
	"github.com/berachain/beacon-kit/mod/beacon/validator"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	dablob "github.com/berachain/beacon-kit/mod/da/pkg/blob"
	dastore "github.com/berachain/beacon-kit/mod/da/pkg/store"
	datypes "github.com/berachain/beacon-kit/mod/da/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	engineclient "github.com/berachain/beacon-kit/mod/execution/pkg/client"
	"github.com/berachain/beacon-kit/mod/execution/pkg/deposit"
	execution "github.com/berachain/beacon-kit/mod/execution/pkg/engine"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/signer"
	payloadbuilder "github.com/berachain/beacon-kit/mod/payload/pkg/builder"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/feed"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/runtime"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/runtime/middleware"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb"
	depositdb "github.com/berachain/beacon-kit/mod/storage/pkg/deposit"
	"github.com/berachain/beacon-kit/mod/storage/pkg/manager"
)

type (
	// AvailabilityStore is a type alias for the availability store.
	AvailabilityStore = *dastore.Store[BeaconBlockBody]

	// BeaconBlock type aliases.
	BeaconBlock       = *types.BeaconBlock
	BeaconBlockBody   = *types.BeaconBlockBody
	BeaconBlockHeader = *types.BeaconBlockHeader

	// BeaconState is a type alias for the BeaconState.
	BeaconState = core.BeaconState[
		BeaconBlockHeader, *types.Eth1Data,
		ExecutionPayloadHeader, *types.Fork,
		*types.Validator, Withdrawal,
	]

	// BeaconKitRuntime is a type alias for the BeaconKitRuntime.
	BeaconKitRuntime = runtime.BeaconKitRuntime[
		AvailabilityStore,
		BeaconBlock,
		BeaconBlockBody,
		BeaconState,
		BlobSidecars,
		DepositStore,
		StorageBackend,
	]

	// BlobSidecars is a type alias for the blob sidecars.
	BlobSidecars = *datypes.BlobSidecars

	// BlobProcessor is a type alias for the blob processor.
	BlobProcessor = *dablob.Processor[
		AvailabilityStore,
		BeaconBlockBody,
	]

	// BlockEvent is a type alias for the block event.
	BlockEvent = *feed.Event[BeaconBlock]

	// BlockFeed is a type alias for the block feed.
	BlockFeed = *event.FeedOf[feed.EventID, BlockEvent]

	// ChainService is a type alias for the chain service.
	ChainService = *blockchain.Service[
		AvailabilityStore,
		BeaconBlock,
		BeaconBlockBody,
		BeaconState,
		BlobSidecars,
		Deposit,
		DepositStore,
	]

	// DBManager is a type alias for the database manager.
	DBManager = *manager.DBManager[
		BeaconBlock,
		BlockEvent,
		event.Subscription,
	]

	// Deposit is a type alias for the deposit.
	Deposit = *types.Deposit

	// DepositService is a type alias for the deposit service.
	DepositService = *deposit.Service[
		BeaconBlock,
		BeaconBlockBody,
		BlockEvent,
		Deposit,
		ExecutionPayload,
		event.Subscription,
		types.WithdrawalCredentials,
	]

	// DepositStore is a type alias for the deposit store.
	DepositStore = *depositdb.KVStore[Deposit]

	// EngineClient is a type alias for the engine client.
	EngineClient = *engineclient.EngineClient[ExecutionPayload]

	// EngineClient is a type alias for the engine client.
	ExecutionEngine = *execution.Engine[ExecutionPayload]

	// ExecutionPayload type aliases.
	ExecutionPayload       = *types.ExecutionPayload
	ExecutionPayloadHeader = *types.ExecutionPayloadHeader

	FinalizeBlockMiddleware = *middleware.FinalizeBlockMiddleware[
		BeaconBlock, BeaconState, BlobSidecars,
	]

	// KVStore is a type alias for the KV store.
	KVStore = *beacondb.KVStore[
		*types.Fork, BeaconBlockHeader, ExecutionPayloadHeader,
		*types.Eth1Data, *types.Validator,
	]

	// LegacyKey type alias to LegacyKey used for LegacySinger construction.
	LegacyKey = signer.LegacyKey

	// LocalBuilder is a type alias for the local builder.
	LocalBuilder = *payloadbuilder.PayloadBuilder[
		BeaconState, ExecutionPayload, ExecutionPayloadHeader,
	]

	// StateProcessor is the type alias for the state processor inteface.
	StateProcessor = blockchain.StateProcessor[
		BeaconBlock,
		BeaconState,
		BlobSidecars,
		*transition.Context,
		Deposit,
	]

	// StorageBackend is the type alias for the storage backend interface.
	StorageBackend = blockchain.StorageBackend[
		AvailabilityStore,
		BeaconBlockBody,
		BeaconState,
		BlobSidecars,
		Deposit,
		DepositStore,
	]

	// ValidatorMiddleware is a type alias for the validator middleware.
	ValidatorMiddleware = *middleware.ValidatorMiddleware[
		AvailabilityStore,
		BeaconBlock,
		BeaconBlockBody,
		BeaconState,
		BlobSidecars,
		StorageBackend,
	]

	// ValidatorService is a type alias for the validator service.
	ValidatorService = *validator.Service[
		BeaconBlock,
		BeaconBlockBody,
		BeaconState,
		BlobSidecars,
		DepositStore,
		*types.ForkData,
	]

	// Withdrawal is a type alias for the engineprimitives withdrawal.
	Withdrawal = *engineprimitives.Withdrawal
)
