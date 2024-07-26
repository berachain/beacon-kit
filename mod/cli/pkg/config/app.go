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
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/server/config"
	"github.com/spf13/viper"
)

// handleAppConfig writes the provided <customConfig> to the file at
// <configDirPath>/app.toml, or reads it into the provided <viper> instance
// if it exists.
func handleAppConfig(
	viper *viper.Viper,
	configDirPath string,
	customAppTemplate string,
	appConfig any,
) error {
	// if the app.toml file does not exist, populate it with the values from
	// <appConfig>
	appCfgFilePath := filepath.Join(configDirPath, "app.toml")
	if _, err := os.Stat(appCfgFilePath); os.IsNotExist(err) {
		return writeAppConfig(
			viper,
			appCfgFilePath,
			customAppTemplate,
			appConfig,
		)
	}

	// merge the app.toml file into the viper instance
	viper.SetConfigType("toml")
	viper.SetConfigName("app")
	viper.AddConfigPath(configDirPath)
	if err := viper.MergeInConfig(); err != nil {
		return fmt.Errorf("failed to merge configuration: %w", err)
	}

	return nil
}

// writeAppConfig creates a new configuration file with default
// values at the specified file path <appCfgFilePath>.
func writeAppConfig(
	rootViper *viper.Viper,
	appConfigFilePath string,
	appTemplate string,
	appConfig any,
) error {
	// customAppTemplate == nil ⊕ customConfig == nil
	if (appTemplate != "" && appConfig == nil) ||
		(appTemplate == "" && appConfig != nil) {
		return errors.New("customAppTemplate and customConfig " +
			"should be both nil or not nil")
	}
	//nolint:nestif // not overly complex
	if appTemplate != "" {
		// set the config template
		if err := config.SetConfigTemplate(appTemplate); err != nil {
			return fmt.Errorf("failed to set config template: %w", err)
		}
		// populate appConfig with the values from the viper instance
		err := rootViper.Unmarshal(&appConfig)
		if err != nil {
			return fmt.Errorf("failed to unmarshal app config: %w", err)
		}
		// write the appConfig to the file at appConfigFilePath
		err = config.WriteConfigFile(appConfigFilePath, appConfig)
		if err != nil {
			return fmt.Errorf("failed to write %s: %w", appConfigFilePath, err)
		}
	} else {
		appConf, err := config.ParseConfig(rootViper)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", appConfigFilePath, err)
		}

		err = config.WriteConfigFile(appConfigFilePath, appConf)
		if err != nil {
			return fmt.Errorf("failed to write %s: %w", appConfigFilePath, err)
		}
	}

	return nil
}
