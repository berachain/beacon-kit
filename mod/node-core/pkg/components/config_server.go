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

package components

import (
	"errors"

	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/mod/consensus/pkg/cometbft/service/server/config"
	servertypes "github.com/berachain/beacon-kit/mod/consensus/pkg/cometbft/service/server/types"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// ServerConfigInput is the input for the dependency injection framework.
type ServerConfigInput struct {
	depinject.In
	AppOpts servertypes.AppOptions
}

// ProvideConfig is a function that provides the BeaconConfig to the
// application.
func ProvideServerConfig(in ConfigInput) (*config.Config, error) {
	v, ok := in.AppOpts.(*viper.Viper)
	if !ok {
		return nil, errors.New("invalid application options type")
	}

	cfg := config.Config{}
	if err := v.Unmarshal(&cfg,
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		))); err != nil {
		return nil, err
	}

	return &cfg, nil
}
