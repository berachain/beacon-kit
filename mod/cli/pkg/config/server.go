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

package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	sdklog "cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/cli/pkg/utils/context"
	"github.com/berachain/beacon-kit/mod/log"
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func SetupConfigAndContext(
	cmd *cobra.Command,
	customAppConfigTemplate string,
	customAppConfig interface{},
	cmtConfig *cmtcfg.Config,
	logger log.AdvancedLogger[any, sdklog.Logger],
) (*server.Context, error) {
	// initialize the server context
	serverCtx, err := InitializeServerContext(cmd, logger)
	if err != nil {
		return nil, err
	}

	// intercept the comet and app config files
	cometConfig, err := handleConfigs(
		serverCtx.Viper, customAppConfigTemplate, customAppConfig, cmtConfig)
	if err != nil {
		return nil, err
	}
	serverCtx.Config = cometConfig

	// not really sure why we have to bind this again here
	if err = bindCommandLineFlags(
		serverCtx.Viper.GetString(flags.FlagHome),
		cmd,
		serverCtx.Viper,
	); err != nil {
		return nil, fmt.Errorf("error binding command line flags: %w", err)
	}
	return serverCtx, nil
}

// InitializeServerContext returns a new server.Context object with the root
// viper instance. The comet config and app config are merged into the viper
// instance. If the app config is empty, the viper instance is populated with
// the app config values.
func InitializeServerContext(
	cmd *cobra.Command,
	logger log.AdvancedLogger[any, sdklog.Logger],
) (*server.Context, error) {
	// Get the executable name and configure the viper instance so that
	// environmental variables are checked based off that name.
	baseName, err := baseName()
	if err != nil {
		return nil, err
	}
	viper := newPrefixedViper(baseName)

	// bind cobra flags to the viper instance
	if err = bindFlags(cmd, viper); err != nil {
		return nil, fmt.Errorf("error binding flags: %w", err)
	}

	// create a new server.Context with the logger and viper instance
	return context.CreateServerContext(logger, viper), nil
}

// newPrefixedViper creates a new viper instance with the given environment
// prefix, and replaces all (.) and (-) with (_).
func newPrefixedViper(prefix string) *viper.Viper {
	viper := viper.New()
	viper.SetEnvPrefix(prefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()
	return viper
}

// baseName returns the base name of the executable.
// ex: full path /usr/local/bin/myapp -> myapp
func baseName() (string, error) {
	executableName, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to fetch executable name: %w", err)
	}
	return path.Base(executableName), nil
}

// bindFlags binds the command line flags to the viper instance.
func bindFlags(cmd *cobra.Command, viper *viper.Viper) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}
	if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
		return err
	}
	return nil
}

// bindCommandLineFlags binds the command line flags to the viper instance.
func bindCommandLineFlags(
	basename string, cmd *cobra.Command, v *viper.Viper,
) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("bindFlags failed: %v", r)
		}
	}()

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// this should be redundant
		err = v.BindEnv(f.Name, fmt.Sprintf("%s_%s", basename, strings.ToUpper(
			strings.ReplaceAll(f.Name, "-", "_"))))
		if err != nil {
			panic(err)
		}

		err = v.BindPFlag(f.Name, f)
		if err != nil {
			panic(err)
		}

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

// handleConfigs writes a new comet config file and app config file, and
// merges them into the provided viper instance.
func handleConfigs(
	viper *viper.Viper,
	customAppTemplate string,
	customConfig any,
	cometConfig *cmtcfg.Config,
) (*cmtcfg.Config, error) {
	rootDir := viper.GetString(flags.FlagHome)
	configDirPath := filepath.Join(rootDir, "config")
	cmtCfgFile := filepath.Join(configDirPath, "config.toml")

	if err := handleCometConfig(
		viper, cmtCfgFile, cometConfig, rootDir, configDirPath); err != nil {
		return nil, err
	}

	if err := handleAppConfig(
		viper, configDirPath, customAppTemplate, customConfig); err != nil {
		return nil, err
	}

	return cometConfig, nil
}
