// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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

package spec

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/berachain/beacon-kit/chain"
	"github.com/berachain/beacon-kit/cli/commands/server/types"
	"github.com/berachain/beacon-kit/cli/flags"
	viperlib "github.com/berachain/beacon-kit/config/viper"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

const (
	devnet  = "devnet"
	mainnet = "mainnet"
	testnet = "testnet"
	custom  = "custom"
)

// Create creates a chain spec based on the app options config flag for "chain-spec".
// If unset, the default of "mainnet" chain spec is used.
func Create(appOpts types.AppOptions) (chain.Spec, error) {
	var (
		chainSpec chain.Spec
		err       error
	)
	switch cast.ToString(appOpts.Get(flags.ChainSpec)) {
	case custom:
		chainSpec, err = handleCustomChainSpec(appOpts)
	case devnet:
		chainSpec, err = DevnetChainSpec()
	case testnet:
		chainSpec, err = TestnetChainSpec()
	case mainnet:
		fallthrough
	default:
		chainSpec, err = MainnetChainSpec()
	}
	if err != nil {
		return nil, err
	}
	if chainSpec == nil {
		return nil, errors.New("no chain spec found")
	}
	return chainSpec, nil
}

// handleCustomChainSpec loads a custom chain spec from the given app options.
func handleCustomChainSpec(appOpts types.AppOptions) (chain.Spec, error) {
	specPath := cast.ToString(appOpts.Get(flags.ChainSpecFilePath))
	if specPath == "" {
		return nil, fmt.Errorf("expected flag '%s' for chain spec", flags.ChainSpecFilePath)
	}
	specData, err := loadSpecData(specPath)
	if err != nil {
		return nil, err
	}
	return chain.NewSpec(specData)
}

// loadSpecData reads the TOML chain-spec file from the given path using Viper,
// unmarshals it into a SpecData, and validates that all required fields are set.
func loadSpecData(path string) (*chain.SpecData, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("toml")
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	// Ensure all required fields are set.
	specData := chain.SpecData{}
	specType := reflect.TypeOf(specData)
	for i := 0; i < specType.NumField(); i++ {
		tag := specType.Field(i).Tag.Get("mapstructure")
		if tag == "" {
			continue
		}
		if !v.IsSet(strings.Split(tag, ",")[0]) {
			return nil, fmt.Errorf("missing required configuration for key: %s", tag)
		}
	}

	// Define a decode hook to handle addresses and domain types.
	decodeHookFunc := mapstructure.ComposeDecodeHookFunc(
		viperlib.StringToExecutionAddressFunc(),
		viperlib.NumericToDomainTypeFunc(),
	)
	if err := v.Unmarshal(&specData, viper.DecodeHook(decodeHookFunc)); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config into SpecData: %w", err)
	}

	return &specData, nil
}
