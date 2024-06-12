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
	"encoding/hex"
	"encoding/json"
	"os"

	"github.com/berachain/beacon-kit/mod/errors"
	viperlib "github.com/berachain/beacon-kit/mod/node-core/pkg/config/viper"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/chain"
	cmttypes "github.com/cometbft/cometbft/types"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

const SpecTypeEnv = "CHAIN_SPEC"

// TODO: This is hood as fuck needs to be improved
// but for now we ball to get CI unblocked.
func FromEnvToAppOpts() primitives.ChainSpec {
	specType := os.Getenv(SpecTypeEnv)
	chainSpec := TestnetChainSpec()
	if specType == "devnet" {
		chainSpec = DevnetChainSpec()
	}

	return chain.NewChainSpec(encodeSpecData(chainSpec.SpecData()))
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
		return nil, errors.Wrap(err, ErrFailedToUnmarshalConfig)
	}
	chainSpec, err := decodeSpecData(cfg.ChainSpec)
	if err != nil {
		return nil, errors.Wrap(err, ErrFailedToDecodeSpec)
	}

	return chain.NewChainSpec(chainSpec), nil
}

// encodeSpecData encodes the chain spec data to be stored in a
// decoder-friendly format.
func encodeSpecData(data primitives.ChainSpecData) primitives.ChainSpecData {
	bz, err := json.Marshal(data.CometBFTValues)
	if err != nil {
		panic(errors.Wrap(err, ErrFailedToEncodeSpec))
	}

	data.CometBFTValues = hex.EncodeToString(bz)
	return data
}

// decodeSpecData decodes the chain spec data from the encoded format
// to the original format.
func decodeSpecData(data primitives.ChainSpecData) (
	primitives.ChainSpecData, error,
) {
	if data.CometBFTValues == nil {
		return data, nil
	}

	var params cmttypes.ConsensusParams
	values, err := hex.DecodeString(cast.ToString(data.CometBFTValues))
	if err != nil {
		return primitives.ChainSpecData{}, err
	}

	if err = json.Unmarshal(values, &params); err != nil {
		return primitives.ChainSpecData{}, err
	}

	data.CometBFTValues = &params
	return data, nil
}
