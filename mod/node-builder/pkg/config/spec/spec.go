// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package spec

import (
	"github.com/berachain/beacon-kit/mod/errors"
	viperlib "github.com/berachain/beacon-kit/mod/node-builder/pkg/config/viper"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/chain"
	cmttypes "github.com/cometbft/cometbft/types"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

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
		return nil,
			errors.Newf("invalid application options type: %T", opts)
	}

	type cfgUnmarshaller struct {
		ChainSpec primitives.ChainSpecData `mapstructure:"chain-spec"`
	}
	cfg := cfgUnmarshaller{}
	if err := v.Unmarshal(&cfg,
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
			viperlib.StringToExecutionAddressFunc(),
			viperlib.StringToDomainTypeFunc(),
			viperlib.StringToCometConsensusParamsFunc[*cmttypes.ConsensusParams](),
		)),
	); err != nil {
		return nil, errors.Newf(
			"failed to decode chain-spec configuration: %w",
			err,
		)
	}

	return chain.NewChainSpec(cfg.ChainSpec), nil
}
