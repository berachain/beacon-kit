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
	sdklog "cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/log/pkg/noop"
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// createServerContext initializes a new server.Context with the default comet
// config, and the provided logger and viper instances.
func CreateServerContext[
	LoggerT log.AdvancedLogger[any, LoggerT],
](
	logger LoggerT,
	viper *viper.Viper,
) *server.Context {
	return &server.Context{
		Viper:  viper,
		Config: cmtcfg.DefaultConfig(),
		Logger: logger,
	}
}

// GetServerContextFromCmd returns a Context from a command or an empty Context
// if it has not been set.
func GetServerContextFromCmd[
	LoggerT log.AdvancedLogger[any, LoggerT],
](
	cmd *cobra.Command,
) *server.Context {
	if v := cmd.Context().Value(server.ServerContextKey); v != nil {
		serverCtxPtr, _ := v.(*server.Context)
		return serverCtxPtr
	}

	return CreateServerContext(&noop.Logger[any, LoggerT]{}, viper.New())
}
