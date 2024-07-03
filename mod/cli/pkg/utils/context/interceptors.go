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

package context

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	sdklog "cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/log"
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// InterceptConfigsAndCreateContext intercepts Comet and App Config files and
// creates a new server.Context object. It returns an error if the configuration
// files cannot be read or parsed.
func InterceptConfigsAndCreateContext(
	cmd *cobra.Command,
	customAppConfigTemplate string,
	customAppConfig interface{},
	cmtConfig *cmtcfg.Config,
	logger log.AdvancedLogger[any, sdklog.Logger],
) (*server.Context, error) {
	serverCtx := newDefaultContextWithLogger(logger)

	// Get the executable name and configure the viper instance so that
	// environmental variables are checked based off that name.
	// The underscore character is used as a separator.
	executableName, err := os.Executable()
	if err != nil {
		return nil, errors.Newf("failed to fetch executable name: %w", err)
	}

	basename := path.Base(executableName)

	// configure the viper instance
	if err = serverCtx.Viper.BindPFlags(cmd.Flags()); err != nil {
		return nil, err
	}
	if err = serverCtx.Viper.BindPFlags(cmd.PersistentFlags()); err != nil {
		return nil, err
	}

	serverCtx.Viper.SetEnvPrefix(basename)
	serverCtx.Viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	serverCtx.Viper.AutomaticEnv()

	// intercept configuration files, using both Viper instances separately
	config, err := interceptConfigs(
		serverCtx.Viper, customAppConfigTemplate, customAppConfig, cmtConfig)
	if err != nil {
		return nil, err
	}

	// return value is a CometBFT configuration object
	serverCtx.Config = config
	if err = bindFlags(basename, cmd, serverCtx.Viper); err != nil {
		return nil, errors.Newf("error binding flags for basename '%s': %w",
			basename, err)
	}

	return serverCtx, nil
}

// newDefaultContextWithLogger returns a new server.Context with the default
// configuration.
func newDefaultContextWithLogger(
	logger log.AdvancedLogger[any, sdklog.Logger],
) *server.Context {
	return &server.Context{
		Viper:  viper.New(),
		Config: cmtcfg.DefaultConfig(),
		Logger: logger,
	}
}

// TODO: Call this function from ProvideConfig to solve all our AppOpts problems
// This will allow us ingest ChainSpec into app.toml, and set logger config
// at build-time.

// WriteAppConfig creates a new configuration file with default values if it
// does not exist, and write it to the specified file path. If the config file
// exists it skips.
//

func WriteAppConfig(
	rootViper *viper.Viper,
	configPath string,
	customAppTemplate string,
	customConfig interface{},
) error {
	appCfgFilePath := filepath.Join(configPath, "app.toml")
	// check if the application configuration file exists
	if _, err := os.Stat(appCfgFilePath); os.IsNotExist(err) {
		return populateConfigFileWithDefaults(
			rootViper,
			appCfgFilePath,
			customAppTemplate,
			customConfig,
			err,
		)
	}
	return nil
}

// populateConfigFileWithDefaults creates a new configuration file with default
// values at the specified file path <appCfgFilePath>.
func populateConfigFileWithDefaults(
	rootViper *viper.Viper,
	appCfgFilePath string,
	customAppTemplate string,
	customConfig interface{},
	statError error,
) error {
	if (customAppTemplate != "" && customConfig == nil) ||
		(customAppTemplate == "" && customConfig != nil) {
		return errors.New("customAppTemplate and customConfig " +
			"should be both nil or not nil")
	}
	//nolint:nestif // not overly complex
	if customAppTemplate != "" {
		// set the configuration template
		if config.SetConfigTemplate(customAppTemplate) != nil {
			return errors.Newf("failed to set config template: %w", statError)
		}

		if rootViper.Unmarshal(&customConfig) != nil {
			return errors.Newf("failed to parse %s: %w",
				appCfgFilePath, statError)
		}

		if config.WriteConfigFile(appCfgFilePath, customConfig) != nil {
			return errors.Newf("failed to write %s: %w", appCfgFilePath,
				statError)
		}
	} else {
		appConf, err := config.ParseConfig(rootViper)
		if err != nil {
			return errors.Newf("failed to parse %s: %w", appCfgFilePath, err)
		}

		if config.WriteConfigFile(appCfgFilePath, appConf) != nil {
			return errors.Newf("failed to write %s: %w", appCfgFilePath, err)
		}
	}
	return nil
}

// interceptConfigs parses and updates a CometBFT configuration file or
// creates a new one and saves it. It also parses and saves the application
// configuration file. The CometBFT configuration file is parsed given a root
// Viper object, whereas the application is parsed with the private
// package-aware viperCfg object.
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

	if err := writeCometConfig(
		rootViper,
		cmtCfgFile,
		conf,
		rootDir,
		configPath); err != nil {
		return nil, err
	}

	// Read into the configuration whatever data the viper instance has for it.
	// This may come from the configuration file above but also any of the other
	// sources viper uses.
	if err := rootViper.Unmarshal(conf); err != nil {
		return nil, err
	}

	conf.SetRoot(rootDir)

	if err := WriteAppConfig(
		rootViper,
		configPath,
		customAppTemplate,
		customConfig); err != nil {
		return nil, err
	}

	rootViper.SetConfigType("toml")
	rootViper.SetConfigName("app")
	rootViper.AddConfigPath(configPath)

	if err := rootViper.MergeInConfig(); err != nil {
		return nil, errors.Newf("failed to merge configuration: %w", err)
	}

	return conf, nil
}

// writeCometConfig creates a new comet config file one with default values.
// If the file exists, it reads and merges it into the provided Viper
// instance.
func writeCometConfig(
	rootViper *viper.Viper,
	cmtCfgFile string,
	conf *cmtcfg.Config,
	rootDir string,
	configPath string,
) error {
	// check if the configuration file exists
	switch _, err := os.Stat(cmtCfgFile); {
	case os.IsNotExist(err):
		// create new config file with default values
		cmtcfg.EnsureRoot(rootDir)

		if err = conf.ValidateBasic(); err != nil {
			return errors.Newf("error in config file: %w", err)
		}

		defaultCometCfg := cmtcfg.DefaultConfig()
		// The SDK is opinionated about those comet values, so we set them here.
		// We verify first that the user has not changed them for not overriding
		// them.
		if conf.Consensus.TimeoutCommit ==
			defaultCometCfg.Consensus.TimeoutCommit {
			//nolint:mnd // 5 second timeout
			conf.Consensus.TimeoutCommit = 5 * time.Second
		}
		if conf.RPC.PprofListenAddress ==
			defaultCometCfg.RPC.PprofListenAddress {
			conf.RPC.PprofListenAddress = "localhost:6060"
		}
		// Write the configuration file to the config directory.
		cmtcfg.WriteConfigFile(cmtCfgFile, conf)

	case err != nil:
		return err

	default:
		// read in the configuration file
		rootViper.SetConfigType("toml")
		rootViper.SetConfigName("config")
		rootViper.AddConfigPath(configPath)

		if err = rootViper.ReadInConfig(); err != nil {
			return errors.Newf("failed to read in %s: %w", cmtCfgFile, err)
		}
	}
	return nil
}

// bindFlags binds the command line flags to the viper instance.
func bindFlags(
	basename string, cmd *cobra.Command, v *viper.Viper,
) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Newf("bindFlags failed: %v", r)
		}
	}()

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to
		// their equivalent keys with underscores, e.g. --favorite-color to
		// STING_FAVORITE_COLOR
		err = v.BindEnv(f.Name, fmt.Sprintf("%s_%s", basename, strings.ToUpper(
			strings.ReplaceAll(f.Name, "-", "_"))))
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
