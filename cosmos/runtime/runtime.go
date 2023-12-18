package runtime

import (
	"context"

	"cosmossdk.io/log"

	"github.com/cosmos/cosmos-sdk/baseapp"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/mempool"

	"github.com/itsdevbear/bolaris/beacon/execution"
	eth "github.com/itsdevbear/bolaris/beacon/execution/ethclient"
	proposal "github.com/itsdevbear/bolaris/cosmos/abci/proposal"
	"github.com/itsdevbear/bolaris/cosmos/config"
	"github.com/itsdevbear/bolaris/cosmos/runtime/miner"
	evmkeeper "github.com/itsdevbear/bolaris/cosmos/x/evm/keeper"
)

// EVMKeeper is an interface that defines the methods needed for the EVM setup.
type EVMKeeper interface {
	// Setup initializes the EVM keeper.
	Setup(execution.EngineCaller) error
}

// CosmosApp is an interface that defines the methods needed for the Cosmos setup.
type CosmosApp interface {
	SetPrepareProposal(sdk.PrepareProposalHandler)
	baseapp.ProposalTxVerifier
	SetMempool(mempool.Mempool)
	SetAnteHandler(sdk.AnteHandler)
	SetExtendVoteHandler(sdk.ExtendVoteHandler)
	SetProcessProposal(sdk.ProcessProposalHandler)
	SetVerifyVoteExtensionHandler(sdk.VerifyVoteExtensionHandler)
	ChainID() string
}

// Polaris is a struct that wraps the Polaris struct from the polar package.
type Polaris struct {
	execution.EngineCaller
	// WrappedMiner is a wrapped version of the Miner component.
	WrappedMiner *miner.Miner
	// logger is the underlying logger supplied by the sdk.
	logger log.Logger
}

// New creates a new Polaris runtime from the provided
// dependencies.
func New(
	appOpts servertypes.AppOptions,
	logger log.Logger,
) (*Polaris, error) {
	var err error
	p := &Polaris{
		logger: logger,
	}

	// Read the configuration from the cosmos app options
	cfg, err := config.ReadConfigFromAppOpts(appOpts)
	if err != nil {
		return nil, err
	}

	// p.Service = execution.NewEngineCaller(ethClient)
	jwtSecert, err := eth.LoadJWTSecret(cfg.ExecutionClient.JWTSecretPath, logger)
	if err != nil {
		return nil, err
	}

	opts := []eth.Option{
		eth.WithHTTPEndpointAndJWTSecret(cfg.ExecutionClient.RPCDialURL, jwtSecert),
		eth.WithLogger(logger),
		eth.WithRequiredChainID(cfg.ExecutionClient.RequiredChainID),
	}

	eth1Client, err := eth.NewEth1Client(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	engineCallerOpts := []execution.Option{
		execution.WithBeaconConfig(&cfg.BeaconConfig),
		execution.WithLogger(logger),
	}

	p.EngineCaller = execution.NewEngineCaller(eth1Client, engineCallerOpts...)

	return p, nil
}

// New creates a new Polaris runtime from the provided
// dependencies, panics on error.
func MustNew(appOpts servertypes.AppOptions, logger log.Logger) *Polaris {
	p, err := New(appOpts, logger)
	if err != nil {
		panic(err)
	}
	return p
}

// Build is a function that sets up the Polaris struct.
// It takes a BaseApp and an EVMKeeper as arguments.
// It returns an error if the setup fails.
func (p *Polaris) Build(app CosmosApp, vs baseapp.ValidatorStore, ek *evmkeeper.Keeper) error {
	// todo use `vs` later?
	_ = vs
	mempool := mempool.NewSenderNonceMempool()
	app.SetMempool(mempool)
	p.WrappedMiner = miner.New(p.EngineCaller, ek, p.logger)

	// Create the proposal handler that will be used to fill proposals with
	// transactions and oracle data.
	// proposalHandler := proposal.NewProposalHandler(
	// 	p.logger,
	// 	baseapp.NoOpPrepareProposal(),
	// 	baseapp.NoOpProcessProposal(),
	// 	ve.NewDefaultValidateVoteExtensionsFn(app.ChainID(), vs),
	// 	ve.NewProcessor(p.WrappedMiner, ek, p.logger).ProcessCommitInfo,
	// )

	defaultProposalHandler := baseapp.NewDefaultProposalHandler(mempool, app)
	proposalHandler := proposal.NewHandler(p.WrappedMiner,
		defaultProposalHandler.PrepareProposalHandler(), defaultProposalHandler.ProcessProposalHandler())
	app.SetPrepareProposal(proposalHandler.PrepareProposalHandler)
	app.SetProcessProposal(proposalHandler.ProcessProposalHandler)

	// if err := p.WrappedMiner.SyncEl(context.Background()); err != nil {
	// 	return err
	// }

	// // Create the vote extensions handler that will be used to extend and verify
	// // vote extensions (i.e. oracle data).
	// voteExtensionsHandler := ve.NewVoteExtensionHandler(
	// 	p.logger,
	// 	time.Second,
	// 	p.WrappedMiner,
	// )
	// app.SetExtendVoteHandler(voteExtensionsHandler.ExtendVoteHandler())
	// app.SetVerifyVoteExtensionHandler(voteExtensionsHandler.VerifyVoteExtensionHandler())

	return nil
}

func (p *Polaris) SyncEL(ctx context.Context) error {
	return p.WrappedMiner.SyncEl(ctx)
}
