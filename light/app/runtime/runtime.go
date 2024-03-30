package runtime

import (
	"context"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/async/dispatch"
	"github.com/berachain/beacon-kit/async/notify"
	"github.com/berachain/beacon-kit/beacon/blockchain"
	"github.com/berachain/beacon-kit/beacon/builder"
	localbuilder "github.com/berachain/beacon-kit/beacon/builder/local"
	"github.com/berachain/beacon-kit/beacon/core"
	"github.com/berachain/beacon-kit/beacon/core/blobs"
	"github.com/berachain/beacon-kit/beacon/core/randao"
	"github.com/berachain/beacon-kit/beacon/staking"
	"github.com/berachain/beacon-kit/beacon/sync"
	"github.com/berachain/beacon-kit/cache"
	"github.com/berachain/beacon-kit/config"
	stakingabi "github.com/berachain/beacon-kit/contracts/abi"
	"github.com/berachain/beacon-kit/crypto"
	"github.com/berachain/beacon-kit/engine"
	engineclient "github.com/berachain/beacon-kit/engine/client"
	"github.com/berachain/beacon-kit/health"
	"github.com/berachain/beacon-kit/lib/abi"
	"github.com/berachain/beacon-kit/primitives"
	beaconruntime "github.com/berachain/beacon-kit/runtime"
	_ "github.com/berachain/beacon-kit/runtime/maxprocs"
	"github.com/berachain/beacon-kit/runtime/service"
)

// NewDefaultBeaconKitRuntime creates a new BeaconKitRuntime with the default
// services.
//
//nolint:funlen // This function is long because it sets up the services.
func NewDefaultBeaconKitRuntime(
	appOpts beaconruntime.AppOptions,
	signer crypto.Signer[primitives.BLSSignature],
	logger log.Logger,
	bsb beaconruntime.BeaconStorageBackend,
) (*beaconruntime.BeaconKitRuntime, error) {
	// Set the module as beacon-kit to override the cosmos-sdk naming.
	logger = logger.With("module", "beacon-kit-light")

	// Read the configuration from the application options.
	cfg, err := config.ReadConfigFromAppOpts(appOpts)
	if err != nil {
		return nil, err
	}

	// Build the service dispatcher.
	gcd, err := dispatch.NewGrandCentralDispatch(
		dispatch.WithLogger(logger),
		dispatch.WithDispatchQueue(
			"dispatch.forkchoice",
			dispatch.QueueTypeSerial,
		),
	)
	if err != nil {
		return nil, err
	}

	// Create the base service, we will the create shallow copies for each
	// service.
	baseService := service.NewBaseService(
		cfg, bsb, gcd, logger,
	)

	// Build the client to interact with the Engine API.
	engineClient := engineclient.New(
		engineclient.WithEngineConfig(&cfg.Engine.Config),
		engineclient.WithLogger(logger),
	)

	// TODO: move.
	engineClient.Start(context.Background())

	// Build the Notification Service.
	notificationService := service.New(
		notify.WithBaseService(baseService.ShallowCopy("notify")),
		notify.WithGCD(gcd),
	)

	// Extrac the staking ABI.
	depositABI, err := stakingabi.BeaconDepositContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	// Build the execution engine.
	executionEngine := engine.NewExecutionEngine(engineClient, logger)

	// Build the staking service.
	stakingService := service.New[staking.Service](
		staking.WithBaseService(baseService.ShallowCopy("staking")),
		staking.WithDepositABI(abi.NewWrappedABI(depositABI)),
		staking.WithExecutionEngine(executionEngine),
	)

	// Build the local builder service.
	localBuilder := service.New[localbuilder.Service](
		localbuilder.WithBaseService(baseService.ShallowCopy("local-builder")),
		localbuilder.WithBuilderConfig(&cfg.Builder),
		localbuilder.WithExecutionEngine(executionEngine),
		localbuilder.WithPayloadCache(cache.NewPayloadIDCache()),
		localbuilder.WithValidatorConfig(&cfg.Validator),
	)

	// Build the Blobs Processor.
	blobsProcessor := blobs.NewProcessor()

	// Build the Randao Processor.
	randaoProcessor := randao.NewProcessor(
		randao.WithSigner(signer),
		randao.WithLogger(logger.With("service", "randao")),
		randao.WithConfig(cfg),
	)

	// Build the builder service.
	builderService := service.New[builder.Service](
		builder.WithBaseService(baseService.ShallowCopy("builder")),
		builder.WithBuilderConfig(&cfg.Builder),
		builder.WithLocalBuilder(localBuilder),
		builder.WithRandaoProcessor(randaoProcessor),
		builder.WithSigner(signer),
	)

	// Build the sync service.
	syncService := service.New[sync.Service](
		sync.WithBaseService(baseService.ShallowCopy("sync")),
		sync.WithEngineClient(engineClient),
		sync.WithConfig(sync.DefaultConfig()),
	)

	// Build the blockchain service.
	chainService := service.New[blockchain.Service](
		blockchain.WithBaseService(baseService.ShallowCopy("blockchain")),
		blockchain.WithBlockValidator(core.NewBlockValidator(&cfg.Beacon)),
		blockchain.WithExecutionEngine(executionEngine),
		blockchain.WithLocalBuilder(localBuilder),
		blockchain.WithPayloadValidator(core.NewPayloadValidator(&cfg.Beacon)),
		blockchain.WithStakingService(stakingService),
		blockchain.WithStateProcessor(
			core.NewStateProcessor(
				&cfg.Beacon,
				blobsProcessor,
				randaoProcessor,
			)),
		blockchain.WithSyncService(syncService),
	)

	// Build the service registry.
	svcRegistry := service.NewRegistry(
		service.WithLogger(logger),
		service.WithService(builderService),
		service.WithService(chainService),
		service.WithService(notificationService),
		service.WithService(stakingService),
		service.WithService(syncService),
	)

	// Build the health service.
	healthService := service.New[health.Service](
		health.WithBaseService(baseService.ShallowCopy("health")),
		health.WithServiceRegistry(svcRegistry),
	)

	if err = svcRegistry.RegisterService(healthService); err != nil {
		return nil, err
	}

	// Pass all the services and options into the BeaconKitRuntime.
	return beaconruntime.NewBeaconKitRuntime(
		beaconruntime.WithBeaconStorageBackend(bsb),
		beaconruntime.WithConfig(cfg),
		beaconruntime.WithLogger(logger),
		beaconruntime.WithServiceRegistry(svcRegistry),
	)
}
