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
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/log/pkg/phuslu"
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetViperFromCmd(cmd *cobra.Command) *viper.Viper {
	value := cmd.Context().Value(ViperContextKey)
	v, ok := value.(*viper.Viper)
	if !ok {
		return viper.New()
	}
	return v
}

func GetLoggerFromCmd[
	LoggerT log.AdvancedLogger[LoggerT],
](cmd *cobra.Command) LoggerT {
	v := cmd.Context().Value(LoggerContextKey)
	logger, ok := v.(LoggerT)
	if !ok {
		//nolint:errcheck // should be safe
		return any(phuslu.NewLogger(cmd.OutOrStdout(), nil)).(LoggerT)
	}
	return logger
}

func GetConfigFromCmd(cmd *cobra.Command) *cmtcfg.Config {
	v := cmd.Context().Value(ViperContextKey)
	viper, ok := v.(*viper.Viper)
	if !ok {
		return cmtcfg.DefaultConfig()
	}
	return GetConfigFromViper(viper)
}

func GetConfigFromViper(v *viper.Viper) *cmtcfg.Config {
	conf := cmtcfg.DefaultConfig()
	err := v.Unmarshal(conf)
	rootDir := v.GetString(flags.FlagHome)
	if err != nil {
		return cmtcfg.DefaultConfig().SetRoot(rootDir)
	}
	return conf.SetRoot(rootDir)
}
