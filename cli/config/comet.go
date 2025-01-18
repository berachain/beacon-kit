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

	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/spf13/viper"
)

// handleCometConfig reads the comet config at <cometConfigFile> into the
// provided <viper> instance. If the file does not exist, it will be populated
// with the values from <cometConfig>.
// <cometConfig> will then be updated with the latest values from <viper>.
func handleCometConfig(
	viper *viper.Viper,
	cometConfigFile string,
	cometConfig *cmtcfg.Config,
	rootDir string,
	configDirPath string,
) error {
	_, err := os.Stat(cometConfigFile)
	if os.IsNotExist(err) {
		// file does not exist, we create a new comet config file one
		// with default values.
		cmtcfg.EnsureRoot(rootDir)
		cmtcfg.WriteConfigFile(cometConfigFile, cometConfig)
	} else if err != nil {
		return err
	}

	// read the config.toml file into the viper instance
	viper.SetConfigType("toml")
	viper.SetConfigName("config")
	viper.AddConfigPath(configDirPath)

	if err = viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read in %s: %w", cometConfigFile, err)
	}

	// update the comet config with the latest values from viper
	if err = viper.Unmarshal(cometConfig); err != nil {
		return err
	}

	cometConfig.SetRoot(rootDir)
	return nil
}
