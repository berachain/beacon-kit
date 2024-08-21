// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	corectx "cosmossdk.io/core/context"
	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/consensus/pkg/cometbft/service/server/config"
	types "github.com/berachain/beacon-kit/mod/consensus/pkg/cometbft/service/server/types"
	cmtcmd "github.com/cometbft/cometbft/cmd/cometbft/commands"
	cmtcfg "github.com/cometbft/cometbft/config"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

// ServerContextKey defines the context key used to retrieve a server.Context
// from
// a command's Context.
const ServerContextKey = sdk.ContextKey("server.context")

// Context server context
// Deprecated: Do not use since we use viper to track all config.
type Context struct {
	Viper  *viper.Viper
	Config *cmtcfg.Config
	Logger log.Logger
}

func NewDefaultContext() *Context {
	return NewContext(
		viper.New(),
		cmtcfg.DefaultConfig(),
		log.NewLogger(os.Stdout),
	)
}

func NewContext(
	v *viper.Viper,
	config *cmtcfg.Config,
	logger log.Logger,
) *Context {
	return &Context{v, config, logger}
}

func bindFlags(
	basename string,
	cmd *cobra.Command,
	v *viper.Viper,
) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("bindFlags failed: %v", r)
		}
	}()

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to
		// their equivalent
		// keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
		err = v.BindEnv(
			f.Name,
			fmt.Sprintf(
				"%s_%s",
				basename,
				strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_")),
			),
		)
		if err != nil {
			panic(err)
		}

		err = v.BindPFlag(f.Name, f)
		if err != nil {
			panic(err)
		}

		// Apply the viper config value to the flag when the flag is not set and
		// viper has a value.
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			err = cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
			if err != nil {
				panic(err)
			}
		}
	})

	return err
}

// InterceptConfigsAndCreateContext performs a pre-run function for the root
// daemon
// application command. It will create a Viper literal and a default server
// Context. The server configuration will either be read and parsed
// or created and saved to disk, where the server Context is updated to reflect
// the CometBFT configuration. It takes custom app config template and config
// settings to create a custom CometBFT configuration. If the custom template
// is empty, it uses default-template provided by the server. The Viper literal
// is used to read and parse the application configuration. Command handlers can
// fetch the server Context to get the CometBFT configuration or to get access
// to Viper.
func InterceptConfigsAndCreateContext(
	cmd *cobra.Command,
	customAppConfigTemplate string,
	customAppConfig interface{},
	cmtConfig *cmtcfg.Config,
) (*Context, error) {
	serverCtx := NewDefaultContext()

	// Get the executable name and configure the viper instance so that
	// environmental variables are checked based off that name. The underscore
	// character is used
	// as a separator.
	executableName, err := os.Executable()
	if err != nil {
		return nil, err
	}

	basename := path.Base(executableName)

	// configure the viper instance
	if err := serverCtx.Viper.BindPFlags(cmd.Flags()); err != nil {
		return nil, err
	}
	if err := serverCtx.Viper.BindPFlags(cmd.PersistentFlags()); err != nil {
		return nil, err
	}

	serverCtx.Viper.SetEnvPrefix(basename)
	serverCtx.Viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	serverCtx.Viper.AutomaticEnv()

	// intercept configuration files, using both Viper instances separately
	config, err := interceptConfigs(
		serverCtx.Viper,
		customAppConfigTemplate,
		customAppConfig,
		cmtConfig,
	)
	if err != nil {
		return nil, err
	}

	// return value is a CometBFT configuration object
	serverCtx.Config = config
	if err = bindFlags(basename, cmd, serverCtx.Viper); err != nil {
		return nil, err
	}

	return serverCtx, nil
}

// GetServerContextFromCmd returns a Context from a command or an empty Context
// if it has not been set.
func GetServerContextFromCmd(cmd *cobra.Command) *Context {
	if v := cmd.Context().Value(ServerContextKey); v != nil {
		serverCtxPtr := v.(*Context)
		return serverCtxPtr
	}

	return NewDefaultContext()
}

// SetCmdServerContext sets a command's Context value to the provided argument.
// If the context has not been set, set the given context as the default.
func SetCmdServerContext(cmd *cobra.Command, serverCtx *Context) error {
	var cmdCtx context.Context

	if cmd.Context() == nil {
		cmdCtx = context.Background()
	} else {
		cmdCtx = cmd.Context()
	}

	cmdCtx = context.WithValue(cmdCtx, ServerContextKey, serverCtx)
	cmdCtx = context.WithValue(cmdCtx, corectx.ViperContextKey, serverCtx.Viper)
	cmdCtx = context.WithValue(
		cmdCtx,
		corectx.LoggerContextKey,
		serverCtx.Logger,
	)

	cmd.SetContext(cmdCtx)

	return nil
}

// interceptConfigs parses and updates a CometBFT configuration file or
// creates a new one and saves it. It also parses and saves the application
// configuration file. The CometBFT configuration file is parsed given a root
// Viper object, whereas the application is parsed with the private
// package-aware
// viperCfg object.
func interceptConfigs(
	rootViper *viper.Viper,
	customAppTemplate string,
	customConfig interface{},
	cmtConfig *cmtcfg.Config,
) (*cmtcfg.Config, error) {
	rootDir := rootViper.GetString(flags.FlagHome)
	configPath := filepath.Join(rootDir, "config")
	cmtCfgFile := filepath.Join(configPath, "config.toml")

	conf := cmtConfig

	switch _, err := os.Stat(cmtCfgFile); {
	case os.IsNotExist(err):
		cmtcfg.EnsureRoot(rootDir)

		if err = conf.ValidateBasic(); err != nil {
			return nil, fmt.Errorf("error in config file: %w", err)
		}

		defaultCometCfg := cmtcfg.DefaultConfig()
		// The SDK is opinionated about those comet values, so we set them here.
		// We verify first that the user has not changed them for not overriding
		// them.
		if conf.Consensus.TimeoutCommit == defaultCometCfg.Consensus.TimeoutCommit {
			conf.Consensus.TimeoutCommit = 5 * time.Second
		}
		if conf.RPC.PprofListenAddress == defaultCometCfg.RPC.PprofListenAddress {
			conf.RPC.PprofListenAddress = "localhost:6060"
		}

		cmtcfg.WriteConfigFile(cmtCfgFile, conf)

	case err != nil:
		return nil, err

	default:
		rootViper.SetConfigType("toml")
		rootViper.SetConfigName("config")
		rootViper.AddConfigPath(configPath)

		if err := rootViper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read in %s: %w", cmtCfgFile, err)
		}
	}

	// Read into the configuration whatever data the viper instance has for it.
	// This may come from the configuration file above but also any of the other
	// sources viper uses.
	if err := rootViper.Unmarshal(conf); err != nil {
		return nil, err
	}

	conf.SetRoot(rootDir)

	appCfgFilePath := filepath.Join(configPath, "app.toml")
	if _, err := os.Stat(appCfgFilePath); os.IsNotExist(err) {
		if (customAppTemplate != "" && customConfig == nil) ||
			(customAppTemplate == "" && customConfig != nil) {
			return nil, errors.New(
				"customAppTemplate and customConfig should be both nil or not nil",
			)
		}

		if customAppTemplate != "" {
			if err := config.SetConfigTemplate(customAppTemplate); err != nil {
				return nil, fmt.Errorf("failed to set config template: %w", err)
			}

			if err = rootViper.Unmarshal(&customConfig); err != nil {
				return nil, fmt.Errorf(
					"failed to parse %s: %w",
					appCfgFilePath,
					err,
				)
			}

			if err := config.WriteConfigFile(appCfgFilePath, customConfig); err != nil {
				return nil, fmt.Errorf(
					"failed to write %s: %w",
					appCfgFilePath,
					err,
				)
			}
		} else {
			appConf, err := config.ParseConfig(rootViper)
			if err != nil {
				return nil, fmt.Errorf("failed to parse %s: %w", appCfgFilePath, err)
			}

			if err := config.WriteConfigFile(appCfgFilePath, appConf); err != nil {
				return nil, fmt.Errorf("failed to write %s: %w", appCfgFilePath, err)
			}
		}
	}

	rootViper.SetConfigType("toml")
	rootViper.SetConfigName("app")
	rootViper.AddConfigPath(configPath)

	if err := rootViper.MergeInConfig(); err != nil {
		return nil, fmt.Errorf("failed to merge configuration: %w", err)
	}

	return conf, nil
}

// AddCommands add server commands.
func AddCommands[T types.Application](
	rootCmd *cobra.Command,
	appCreator types.AppCreator[T],
	opts StartCmdOptions[T],
) {
	cometCmd := &cobra.Command{
		Use:     "comet",
		Aliases: []string{"cometbft", "tendermint"},
		Short:   "CometBFT subcommands",
	}

	cometCmd.AddCommand(
		ShowNodeIDCmd(),
		ShowValidatorCmd(),
		ShowAddressCmd(),
		VersionCmd(),
		cmtcmd.ResetAllCmd,
		cmtcmd.ResetStateCmd,
		BootstrapStateCmd[T](appCreator),
	)

	startCmd := StartCmdWithOptions(appCreator, opts)
	rootCmd.AddCommand(
		startCmd,
		cometCmd,
		version.NewVersionCommand(),
		NewRollbackCmd[T](appCreator),
	)
}

// ListenForQuitSignals listens for SIGINT and SIGTERM. When a signal is
// received,
// the cleanup function is called, indicating the caller can gracefully exit or
// return.
//
// Note, the blocking behavior of this depends on the block argument.
// The caller must ensure the corresponding context derived from the cancelFn is
// used correctly.
func ListenForQuitSignals(
	g *errgroup.Group,
	block bool,
	cancelFn context.CancelFunc,
	logger log.Logger,
) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	f := func() {
		sig := <-sigCh
		cancelFn()

		logger.Info("caught signal", "signal", sig.String())
	}

	if block {
		g.Go(func() error {
			f()
			return nil
		})
	} else {
		go f()
	}
}

// OpenDB opens the application database using the appropriate driver.
func OpenDB(rootDir string, backendType dbm.BackendType) (dbm.DB, error) {
	dataDir := filepath.Join(rootDir, "data")
	return dbm.NewDB("application", backendType, dataDir)
}

func openTraceWriter(traceWriterFile string) (w io.WriteCloser, err error) {
	if traceWriterFile == "" {
		return
	}
	return os.OpenFile(
		traceWriterFile,
		os.O_WRONLY|os.O_APPEND|os.O_CREATE,
		0o666,
	)
}
