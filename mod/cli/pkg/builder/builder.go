package builder

import (
	"os"

	"cosmossdk.io/client/v2/autocli"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	cmdlib "github.com/berachain/beacon-kit/mod/cli/pkg/commands"
	"github.com/berachain/beacon-kit/mod/primitives"
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// CLIBuilder is a builder that incrementally constructs a CLI command.
type CLIBuilder[T servertypes.Application] struct {
	depInjectCfg depinject.Config
	name         string
	description  string
	components   []any
	runHandler   func(cmd *cobra.Command,
		customAppConfigTemplate string,
		customAppConfig interface{},
		cmtConfig *cmtcfg.Config,
	) error
	AppCreator   servertypes.AppCreator[T]
	RootCmdSetup func(cmd *cobra.Command,
		mm *module.Manager,
		appCreator servertypes.AppCreator[T],
		chainSpec primitives.ChainSpec,
	)
}

// New returns a new CLIBuilder with the given options.
func New[T servertypes.Application](opts ...Opt[T]) *CLIBuilder[T] {
	cb := &CLIBuilder[T]{}
	for _, opt := range opts {
		opt(cb)
	}
	return cb
}

// Build builds the CLI commands
func (cb *CLIBuilder[T]) Build() (*cmdlib.Root, error) {
	// dependencies for the root command
	var (
		autoCliOpts autocli.AppOptions
		mm          *module.Manager
		clientCtx   client.Context
		chainSpec   primitives.ChainSpec
	)
	// build dependencies for the root command
	if err := depinject.Inject(
		depinject.Configs(
			cb.depInjectCfg,
			depinject.Supply(
				log.NewLogger(os.Stdout),
				viper.GetViper(),
				// empty middleware must be supplied here because it is a direct
				// dependency of the Module
				// emptyABCIMiddleware(),
			),
			depinject.Provide(
				cb.components...,
			),
		),
		&autoCliOpts,
		&mm,
		&clientCtx,
		&chainSpec,
	); err != nil {
		return nil, err
	}

	cmd := &cobra.Command{
		Use:   cb.name,
		Short: cb.description,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// set the default command outputs
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())

			var err error
			clientCtx, err = client.ReadPersistentCommandFlags(
				clientCtx,
				cmd.Flags(),
			)
			if err != nil {
				return err
			}

			customClientTemplate, customClientConfig := InitClientConfig()
			clientCtx, err = config.CreateClientConfig(
				clientCtx,
				customClientTemplate,
				customClientConfig,
			)
			if err != nil {
				return err
			}

			if err = client.SetCmdClientContextHandler(
				clientCtx, cmd,
			); err != nil {
				return err
			}

			return cb.runHandler(
				cmd,
				DefaultAppConfigTemplate(),
				DefaultAppConfig(),
				DefaultCometConfig(),
			)
		},
	}

	cmdlib.DefaultRootCommandSetup(
		cmd,
		mm,
		cb.AppCreator,
		chainSpec,
	)

	if err := autoCliOpts.EnhanceRootCommand(cmd); err != nil {
		return nil, err
	}

	return cmdlib.NewRoot(cmd), nil
}

// InitClientConfig sets up the default client configuration, allowing for
// overrides.
func InitClientConfig() (string, interface{}) {
	clientCfg := config.DefaultConfig()
	clientCfg.KeyringBackend = "test"
	return config.DefaultClientConfigTemplate, clientCfg
}
