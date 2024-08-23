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
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	corectx "cosmossdk.io/core/context"
	sdklog "github.com/berachain/beacon-kit/mod/consensus/pkg/cometbft/service/log"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/log"
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// SetupCommand returns a cobra.Command initialized with a viper
// instance. If the files expected to contain the comet and app
// configs are empty, it will be populated with the values from
// <appConfig> and <cmtConfig>. In either case, the resulting
// values in these files will be merged with viper.
func SetupCommand[
	LoggerT log.AdvancedLogger[LoggerT],
](
	cmd *cobra.Command,
	appTemplate string,
	appConfig any,
	cmtConfig *cmtcfg.Config,
	logger LoggerT,
) error {
	// initialize the server context
	if err := InitializeCmd(cmd, logger); err != nil {
		return err
	}

	if err := handleConfigs(
		client.GetViperFromCmd(cmd),
		appTemplate, appConfig, cmtConfig,
	); err != nil {
		return err
	}

	return nil
}

// InitializeCmd returns a command object with the root viper instance.
// The comet config and app config are merged into the viper instance.
// If the app config is empty, the viper instance is populated with
// the app config values.
func InitializeCmd[
	LoggerT log.AdvancedLogger[LoggerT],
](
	cmd *cobra.Command,
	logger LoggerT,
) error {
	// Get the executable name and configure the viper instance so that
	// environmental variables are checked based off that name.
	baseName, err := baseName()
	if err != nil {
		return err
	}
	viper := newPrefixedViper(baseName)

	// bind cobra flags to the viper instance
	if err = bindFlags(baseName, cmd, viper); err != nil {
		return errors.Wrapf(ErrFlagBind, "error binding flags: %w", err)
	}

	ctx := cmd.Context()
	ctx = context.WithValue(ctx, corectx.ViperContextKey, viper)
	ctx = context.WithValue(
		ctx, corectx.LoggerContextKey, sdklog.WrapSDKLogger(logger),
	)
	cmd.SetContext(ctx)

	return nil
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
func bindFlags(
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
) error {
	rootDir := viper.GetString(flags.FlagHome)
	configDirPath := filepath.Join(rootDir, "config")
	cmtCfgFile := filepath.Join(configDirPath, "config.toml")

	if err := handleCometConfig(
		viper, cmtCfgFile, cometConfig, rootDir, configDirPath,
	); err != nil {
		return err
	}

	return handleAppConfig(
		viper, configDirPath, customAppTemplate, customConfig,
	)
}
