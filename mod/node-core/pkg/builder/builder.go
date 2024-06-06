package builder

import (
	"os"

	"cosmossdk.io/client/v2/autocli"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/da/pkg/kzg/noop"
	dastore "github.com/berachain/beacon-kit/mod/da/pkg/store"
	engineclient "github.com/berachain/beacon-kit/mod/execution/pkg/client"
	node "github.com/berachain/beacon-kit/mod/node-core/pkg"
	cmdlib "github.com/berachain/beacon-kit/mod/node-core/pkg/commands"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/signer"
	"github.com/berachain/beacon-kit/mod/primitives"
	depositdb "github.com/berachain/beacon-kit/mod/storage/pkg/deposit"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/types/module"
	gokzg4844 "github.com/crate-crypto/go-kzg-4844"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Builder[NodeT node.NodeI] struct {
	node NodeT

	name         string
	description  string
	depInjectCfg depinject.Config
	chainSpec    primitives.ChainSpec
}

// NewBuilder returns a new Builder.
func NewBuilder[NodeT node.NodeI](opts ...Opt[NodeT]) *Builder[NodeT] {
	b := &Builder[NodeT]{
		node: node.New[NodeT](),
	}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

// Build builds the application.
func (b *Builder[NodeT]) Build() (NodeT, error) {
	rootCmd, err := b.buildRootCmd()
	if err != nil {
		return b.node, err
	}

	b.node.SetRootCmd(rootCmd)
	return b.node, nil
}

// buildRootCmd builds the root command for the application.
func (b *Builder[NodeT]) buildRootCmd() (*cobra.Command, error) {
	var (
		autoCliOpts autocli.AppOptions
		mm          *module.Manager
		clientCtx   client.Context
	)
	if err := depinject.Inject(
		depinject.Configs(
			b.depInjectCfg,
			depinject.Supply(
				log.NewLogger(os.Stdout),
				viper.GetViper(),
				b.chainSpec,
				&depositdb.KVStore[*types.Deposit]{},
				&engineclient.EngineClient[*types.ExecutionPayload]{},
				&gokzg4844.JSONTrustedSetup{},
				&noop.Verifier{},
				&dastore.Store[types.BeaconBlockBody]{},
				&signer.BLSSigner{},
			),
			depinject.Provide(
				components.ProvideNoopTxConfig,
				components.ProvideClientContext,
				components.ProvideKeyring,
				components.ProvideConfig,
				components.ProvideTelemetrySink,
			),
		),
		&autoCliOpts,
		&mm,
		&clientCtx,
	); err != nil {
		return nil, err
	}

	cmd := &cobra.Command{
		Use:   b.name,
		Short: b.description,
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

			customClientTemplate, customClientConfig := components.InitClientConfig()
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

			return server.InterceptConfigsPreRunHandler(
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
		b.AppCreator,
		b.chainSpec,
	)

	if err := autoCliOpts.EnhanceRootCommand(cmd); err != nil {
		return nil, err
	}

	return cmd, nil
}
