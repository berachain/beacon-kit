package runtime

import (
	"context"

	"cosmossdk.io/log"

	"github.com/berachain/beacon-kit/light/app"
	"github.com/berachain/beacon-kit/light/mod/core"
	"github.com/berachain/beacon-kit/light/mod/runtime/services/blockchain"
	"github.com/berachain/beacon-kit/mod/execution"
	engineclient "github.com/berachain/beacon-kit/mod/execution/client"
	"github.com/berachain/beacon-kit/mod/node-builder/service"
	"github.com/berachain/beacon-kit/mod/node-builder/utils/jwt"
	beaconruntime "github.com/berachain/beacon-kit/mod/runtime"
)

// NewDefaultBeaconKitRuntime creates a new BeaconKitRuntime with the default
// services.
func NewDefaultLightRuntime(
	cfg *app.Config,
	jwtSecret *jwt.Secret,
	bsb beaconruntime.BeaconStorageBackend,
	logger log.Logger,
) (*beaconruntime.BeaconKitRuntime, error) {
	// Set the module as beacon-kit to override the cosmos-sdk naming.
	logger = logger.With("module", "beacon-kit")

	// Create the base service, we will the create shallow copies for each
	// service.
	baseService := service.NewBaseService(
		cfg, bsb, logger,
	)

	// Build the client to interact with the Engine API.
	engineClient := engineclient.New(
		engineclient.WithEngineConfig(&cfg.Engine),
		engineclient.WithJWTSecret(jwtSecret),
		engineclient.WithLogger(logger),
	)

	// TODO: move.
	engineClient.Start(context.Background())

	// // Extrac the staking ABI.
	// depositABI, err := abi.BeaconDepositContractMetaData.GetAbi()
	// if err != nil {
	// 	return nil, err
	// }

	// Build the execution engine.
	executionEngine := execution.NewEngine(engineClient, logger)

	// // Build the staking service.
	// stakingService := service.New[staking.Service](
	// 	staking.WithBaseService(baseService.ShallowCopy("staking")),
	// 	staking.WithDepositABI(depositABI),
	// 	staking.WithExecutionEngine(executionEngine),
	// )

	// // Build the local builder service.
	// localBuilder := service.New[localbuilder.Service](
	// 	localbuilder.WithBaseService(baseService.ShallowCopy("local-builder")),
	// 	localbuilder.WithBuilderConfig(&cfg.Builder),
	// 	localbuilder.WithExecutionEngine(executionEngine),
	// 	localbuilder.WithPayloadCache(cache.NewPayloadIDCache()),
	// )

	// Build the Blob Verifier.
	// TODO: Implement this service by implementing the blob store.
	// blobVerifier, err := da.NewBlobVerifier(kzgTrustedSetup)
	// if err != nil {
	// 	return nil, err
	// }

	// Build the Blobs Processor.
	// blobsProcessor := blobs.NewProcessor(blobVerifier)

	// Build the Randao Processor.
	// TODO: Implement this service by implementing the validator store.
	// randaoProcessor := randao.NewProcessor(
	// 	randao.WithSigner(signer),
	// 	randao.WithLogger(logger.With("service", "randao")),
	// 	randao.WithConfig(cfg),
	// )

	// Build the builder service.
	// builderService := service.New[builder.Service](
	// 	builder.WithBaseService(baseService.ShallowCopy("builder")),
	// 	builder.WithBuilderConfig(&cfg.Builder),
	// 	builder.WithLocalBuilder(localBuilder),
	// 	builder.WithRandaoProcessor(randaoProcessor),
	// 	builder.WithSigner(signer),
	// )

	// Build the blockchain service.
	chainService := service.New[blockchain.Service](
		blockchain.WithBaseService(baseService.ShallowCopy("blockchain")),
		blockchain.WithBlockValidator(core.NewBlockValidator(&cfg.Beacon)),
		blockchain.WithExecutionEngine(executionEngine),
		blockchain.WithPayloadValidator(core.NewPayloadValidator(&cfg.Beacon)),
		blockchain.WithStateProcessor(
			core.NewStateProcessor(
				&cfg.Beacon,
				// blobsProcessor,
				// randaoProcessor,
				logger,
			)),
	)

	// Build the service registry.
	svcRegistry := service.NewRegistry(
		service.WithLogger(logger),
		// service.WithService(builderService),
		service.WithService(chainService),
		// service.WithService(stakingService),
	)

	// Pass all the services and options into the BeaconKitRuntime.
	return beaconruntime.NewBeaconKitRuntime(
		beaconruntime.WithBeaconStorageBackend(bsb),
		beaconruntime.WithConfig(cfg),
		beaconruntime.WithLogger(logger),
		beaconruntime.WithServiceRegistry(svcRegistry),
	)
}
