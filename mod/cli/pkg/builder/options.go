package builder

import (
	"cosmossdk.io/depinject"
	cmdlib "github.com/berachain/beacon-kit/mod/cli/pkg/commands"
	cmtcfg "github.com/cometbft/cometbft/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cobra"
)

// Opt is a type that defines a function that modifies CLIBuilder.
type Opt[T servertypes.Application] func(*CLIBuilder[T])

// WithName sets the name for the CLIBuilder
func WithName[T servertypes.Application](name string) Opt[T] {
	return func(cb *CLIBuilder[T]) {
		cb.name = name
	}
}

// WithDescription sets the description for the CLIBuilder
func (cb *CLIBuilder[T]) WithDescription(description string) Opt[T] {
	return func(cb *CLIBuilder[T]) {
		cb.description = description
	}
}

func (cb *CLIBuilder[T]) WithDepInjectConfig(cfg depinject.Config) Opt[T] {
	return func(cb *CLIBuilder[T]) {
		cb.depInjectCfg = cfg
	}
}

// WithComponents sets the components for the CLIBuilder
func (cb *CLIBuilder[T]) WithComponents(components []any) Opt[T] {
	return func(cb *CLIBuilder[T]) {
		cb.components = components
	}
}

func (cb *CLIBuilder[T]) WithRunHandler(
	runHandler func(cmd *cobra.Command,
		customAppConfigTemplate string,
		customAppConfig interface{},
		cmtConfig *cmtcfg.Config,
	) error,
) Opt[T] {
	return func(cb *CLIBuilder[T]) {
		cb.runHandler = runHandler
	}
}

// WithDefaultRootCommandSetup sets the root command setup func to the default
func (cb *CLIBuilder[T]) WithDefaultRootCommandSetup() Opt[T] {
	return func(cb *CLIBuilder[T]) {
		cb.RootCmdSetup = cmdlib.DefaultRootCommandSetup
	}
}

// WithAppCreator sets the cosmos app creator for the CLIBuilder
func (cb *CLIBuilder[T]) WithAppCreator(
	appCreator servertypes.AppCreator[T],
) Opt[T] {
	return func(cb *CLIBuilder[T]) {
		cb.AppCreator = appCreator
	}
}
