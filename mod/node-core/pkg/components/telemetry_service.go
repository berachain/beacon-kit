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

	"github.com/berachain/beacon-kit/mod/consensus/pkg/cometbft/service/server/types"
	"github.com/berachain/beacon-kit/mod/observability/pkg/telemetry"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// ProvideTelemetryService is a function that provides a TelemetrySink.
func ProvideTelemetryService(
	opts types.AppOptions,
) (*telemetry.Service, error) {
	v, ok := opts.(*viper.Viper)
	if !ok {
		return nil, errors.New("invalid application options type")
	}

	type cfgUnmarshaller struct {
		Telemetry telemetry.Config `mapstructure:"telemetry"`
	}
	cfg := cfgUnmarshaller{}
	if err := v.Unmarshal(&cfg,
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		))); err != nil {
		return nil, err
	}

	return telemetry.NewService(&cfg.Telemetry)
}
