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
func WithDescription[T servertypes.Application](description string) Opt[T] {
	return func(cb *CLIBuilder[T]) {
		cb.description = description
	}
}

// WithDepInjectConfig sets the depinject config for the CLIBuilder
func WithDepInjectConfig[T servertypes.Application](
	cfg depinject.Config) Opt[T] {
	return func(cb *CLIBuilder[T]) {
		cb.depInjectCfg = cfg
	}
}

// WithComponents sets the components for the CLIBuilder
func WithComponents[T servertypes.Application](components []any) Opt[T] {
	return func(cb *CLIBuilder[T]) {
		cb.components = components
	}
}

// WithModuleDeps populates the slice of direct module dependencies to be
// supplied to depinject
func WithModuleDeps[T servertypes.Application](deps []any) Opt[T] {
	return func(cb *CLIBuilder[T]) {
		cb.moduleDeps = deps
	}
}

// WithRunHandler sets the run handler for the CLIBuilder
func WithRunHandler[T servertypes.Application](
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
func WithDefaultRootCommandSetup[T servertypes.Application]() Opt[T] {
	return func(cb *CLIBuilder[T]) {
		cb.rootCmdSetup = cmdlib.DefaultRootCommandSetup
	}
}

// WithAppCreator sets the cosmos app creator for the CLIBuilder
func WithAppCreator[T servertypes.Application](
	appCreator servertypes.AppCreator[T],
) Opt[T] {
	return func(cb *CLIBuilder[T]) {
		cb.appCreator = appCreator
	}
}
