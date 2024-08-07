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
	viperlib "github.com/berachain/beacon-kit/mod/cli/pkg/v2/config/v2/viper"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// AppOptions is from the SDK, we should look to remove its usage.
type AppOptions interface {
	Get(string) interface{}
}

// MustReadConfigFromAppOpts reads the configuration options from the given
// application options.
func MustReadConfigFromAppOpts[ConfigT any](opts AppOptions) *ConfigT {
	cfg, err := ReadConfigFromAppOpts[ConfigT](opts)
	if err != nil {
		panic(err)
	}
	return cfg
}

// ReadConfigFromAppOpts reads the configuration options from the given
// application options.
func ReadConfigFromAppOpts[ConfigT any](opts AppOptions) (*ConfigT, error) {
	v, ok := opts.(*viper.Viper)
	if !ok {
		return nil, errors.Newf("invalid application options type: %T", opts)
	}

	type cfgUnmarshaller struct {
		BeaconKit ConfigT `mapstructure:"beacon-kit"`
	}
	cfg := cfgUnmarshaller{}
	if err := v.Unmarshal(&cfg,
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
			viperlib.StringToExecutionAddressFunc(),
			viperlib.StringToDialURLFunc(),
			viperlib.StringToConnectionURLFunc(),
		))); err != nil {
		return nil, errors.Newf(
			"failed to decode beacon-kit configuration: %w",
			err,
		)
	}

	return &cfg.BeaconKit, nil
}
