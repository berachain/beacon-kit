package main

import (
	"context"
	"fmt"
	"os"

	clicomponents "github.com/berachain/beacon-kit/mod/cli/pkg/components"
	"github.com/berachain/beacon-kit/mod/config"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/genesis"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/state"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	cometbft "github.com/berachain/beacon-kit/mod/consensus/pkg/comet"
	consruntimetypes "github.com/berachain/beacon-kit/mod/consensus/pkg/types"
	dastore "github.com/berachain/beacon-kit/mod/da/pkg/store"
	datypes "github.com/berachain/beacon-kit/mod/da/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/log/pkg/phuslu"
	"github.com/berachain/beacon-kit/mod/node/pkg/app"
	"github.com/berachain/beacon-kit/mod/node/pkg/app/components"
	"github.com/berachain/beacon-kit/mod/node/pkg/app/components/storage"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	middlewarev2 "github.com/berachain/beacon-kit/mod/runtime/pkg/middleware/v2"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
	statedb "github.com/berachain/beacon-kit/mod/state-transition/pkg/core/state"
	beacondbv2 "github.com/berachain/beacon-kit/mod/storage/pkg/beacondb/v2"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb/v2/store"
	"github.com/berachain/beacon-kit/mod/storage/pkg/block"
	depositdb "github.com/berachain/beacon-kit/mod/storage/pkg/deposit"
)

func run() error {
	ctx := context.Background()
	cfg := config.DefaultConfig()
	cfg.Engine.JWTSecretPath = "./testing/files/jwt.hex"
	appOpts := &components.AppOptions{
		HomeDir: ".tmp/testingd",
	}
	logger := clicomponents.ProvideLogger(clicomponents.LoggerInput{
		Cfg: cfg,
		Out: os.Stdout,
	})
	chainSpec := components.ProvideChainSpec()

	var (
		storageBackend  = &StorageBackend{}
		stateProcessor  = &StateProcessor{}
		consensusClient = &ABCIMiddlewareV2{}
	)
	appBuilder := app.NewBuilder[*StorageBackend, *StateProcessor]()
	appBuilder.WithComponents(components.DefaultComponentsWithStandardTypes()...)
	appBuilder.WithStateProcessor(stateProcessor)
	appBuilder.WithStorageBackend(storageBackend)
	appBuilder.WithConsensusClient(consensusClient)
	app, err := appBuilder.Build(logger, appOpts, cfg)
	if err != nil {
		return err
	}

	fmt.Println("APP BUILT")
	if err := app.Start(ctx); err != nil {
		return err
	}

	consensus := cometbft.NewConsensus(cfg.CometBFT, app.Logger, app, chainSpec)
	return consensus.Start(ctx)
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

type (
	ABCIMiddlewareV2 = middlewarev2.ABCIMiddleware[
		*AttestationData,
		*AvailabilityStore,
		*BeaconBlock,
		*BeaconState,
		*BlobSidecars,
		*Deposit,
		*ExecutionPayload,
		*Genesis,
		*SlashingInfo,
		*SlotData,
		*StorageBackend,
	]

	// AttestationData is a type alias for the attestation data.
	AttestationData = types.AttestationData

	// AvailabilityStore is a type alias for the availability store.
	AvailabilityStore = dastore.Store[*BeaconBlockBody]

	// BeaconBlock type aliases.
	BeaconBlock       = types.BeaconBlock
	BeaconBlockBody   = types.BeaconBlockBody
	BeaconBlockHeader = types.BeaconBlockHeader

	// BeaconState is a type alias for the BeaconState.
	BeaconState = statedb.StateDB[
		*BeaconBlockHeader,
		*BeaconStateMarshallable,
		*Eth1Data,
		*ExecutionPayloadHeader,
		*Fork,
		*StateManager,
		*Validator,
		*Withdrawal,
		WithdrawalCredentials,
	]

	// BeaconStateMarshallable is a type alias for the BeaconStateMarshallable.
	BeaconStateMarshallable = state.BeaconStateMarshallable[
		*BeaconBlockHeader,
		*Eth1Data,
		*ExecutionPayloadHeader,
		*Fork,
		*Validator,
	]

	// BlobSidecars is a type alias for the blob sidecars.
	BlobSidecars = datypes.BlobSidecars

	// BlockStore is a type alias for the block store.
	BlockStore = block.KVStore[*BeaconBlock]

	// Context is a type alias for the transition context.
	Context = transition.Context

	// Deposit is a type alias for the deposit.
	Deposit = types.Deposit

	// DepositStore is a type alias for the deposit store.
	DepositStore = depositdb.KVStore[*Deposit]

	// Eth1Data is a type alias for the eth1 data.
	Eth1Data = types.Eth1Data

	// ExecutionPayload type aliases.
	ExecutionPayload       = types.ExecutionPayload
	ExecutionPayloadHeader = types.ExecutionPayloadHeader

	// Fork is a type alias for the fork.
	Fork = types.Fork

	// ForkData is a type alias for the fork data.
	ForkData = types.ForkData

	// Genesis is a type alias for the genesis.
	Genesis = genesis.Genesis[
		*Deposit,
		*ExecutionPayloadHeader,
	]

	// Logger is a type alias for the logger.
	Logger = phuslu.Logger

	StateStore = store.StateStore

	StateManager = beacondbv2.StateManager[
		*BeaconBlockHeader,
		*Eth1Data,
		*ExecutionPayloadHeader,
		*Fork,
		*Validator,
	]

	// SlashingInfo is a type alias for the slashing info.
	SlashingInfo = types.SlashingInfo

	// SlotData is a type alias for the incoming slot.
	SlotData = consruntimetypes.SlotData[
		*AttestationData,
		*SlashingInfo,
	]

	// StateProcessor is the type alias for the state processor interface.
	StateProcessor = core.StateProcessor[
		*BeaconBlock,
		*BeaconBlockBody,
		*BeaconBlockHeader,
		*BeaconState,
		*Context,
		*Deposit,
		*Eth1Data,
		*ExecutionPayload,
		*ExecutionPayloadHeader,
		*Fork,
		*ForkData,
		*StateManager,
		*Validator,
		*Withdrawal,
		WithdrawalCredentials,
	]

	// StorageBackend is the type alias for the storage backend interface.
	StorageBackend = storage.Backend[
		*AvailabilityStore,
		*BeaconBlock,
		*BeaconBlockBody,
		*BeaconBlockHeader,
		*BeaconState,
		*BeaconStateMarshallable,
		*BlobSidecars,
		*BlockStore,
		*Deposit,
		*DepositStore,
		*Eth1Data,
		*ExecutionPayloadHeader,
		*Fork,
		*StateManager,
		*Validator,
		*Withdrawal,
		WithdrawalCredentials,
	]

	// Validator is a type alias for the validator.
	Validator = types.Validator

	// Withdrawal is a type alias for the engineprimitives withdrawal.
	Withdrawal = engineprimitives.Withdrawal

	// WithdrawalCredentials is a type alias for the withdrawal credentials.
	WithdrawalCredentials = types.WithdrawalCredentials
)
