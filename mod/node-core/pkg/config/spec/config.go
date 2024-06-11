// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
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

package spec

import (
	"os"

	"github.com/berachain/beacon-kit/mod/errors"
	viperlib "github.com/berachain/beacon-kit/mod/node-core/pkg/config/viper"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/chain"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// TODO: This is hood as fuck needs to be improved
// but for now we ball to get CI unblocked.
func FromEnv() primitives.ChainSpec {
	specType := os.Getenv("CHAIN_SPEC")
	chainSpec := TestnetChainSpec()
	if specType == "devnet" {
		chainSpec = DevnetChainSpec()
	}

	return chainSpec
}

// MustReadFromAppOpts reads the configuration options from the given
// application options.
func MustReadFromAppOpts(
	opts servertypes.AppOptions,
) primitives.ChainSpec {
	spec, err := ReadFromAppOpts(opts)
	if err != nil {
		panic(err)
	}
	return spec
}

// ReadFromAppOpts reads the configuration options from the given
// application options.
func ReadFromAppOpts(
	opts servertypes.AppOptions,
) (primitives.ChainSpec, error) {
	v, ok := opts.(*viper.Viper)
	if !ok {
		return nil, errors.Wrapf(ErrInvalidOptionsType, "%v", opts)
	}

	type cfgUnmarshaller struct {
		ChainSpec primitives.ChainSpecData `mapstructure:"chain-spec"`
	}
	cfg := cfgUnmarshaller{}
	if err := v.Unmarshal(&cfg,
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
			viperlib.StringToExecutionAddressFunc(),
			viperlib.StringToDomainTypeFunc(),
		)),
	); err != nil {
		return nil, errors.Newf(
			"failed to decode chain-spec configuration: %w",
			err,
		)
	}

	return chain.NewChainSpec(cfg.ChainSpec), nil
}
