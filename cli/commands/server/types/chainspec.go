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

package types

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/berachain/beacon-kit/chain"
	"github.com/berachain/beacon-kit/config/spec"
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

const (
	ChainSpecTypeEnvVar           = "CHAIN_SPEC"
	DevnetChainSpecType           = "devnet"
	MainnetChainSpecType          = "mainnet"
	TestnetChainSpecType          = "testnet"
	ConfigurableChainSpecType     = "configurable"
	FlagConfigurableChainSpecPath = "spec"
)

// ChainSpecCreator is a function that allows us to lazily initialize the ChainSpec
type ChainSpecCreator func(AppOptions) (chain.Spec, error)

func CreateChainSpec(appOpts AppOptions) (chain.Spec, error) {
	var (
		chainSpec chain.Spec
		err       error
	)
	switch os.Getenv(ChainSpecTypeEnvVar) {
	case ConfigurableChainSpecType:
		chainSpec, err = handleConfigurableChainSpec(appOpts)
	case DevnetChainSpecType:
		chainSpec, err = spec.DevnetChainSpec()
	case TestnetChainSpecType:
		chainSpec, err = spec.TestnetChainSpec()
	case MainnetChainSpecType:
		chainSpec, err = spec.MainnetChainSpec()
	default:
		chainSpec, err = spec.MainnetChainSpec()
	}
	if err != nil {
		return nil, err
	}
	if chainSpec == nil {
		return nil, errors.New("no chain spec found")
	}
	return chainSpec, nil
}

func handleConfigurableChainSpec(appOpts AppOptions) (chain.Spec, error) {
	specPath := cast.ToString(appOpts.Get(FlagConfigurableChainSpecPath))
	if specPath == "" {
		return nil, fmt.Errorf("expected flag '%s' for chain spec", FlagConfigurableChainSpecPath)
	}
	specData, err := loadSpecData(specPath)
	if err != nil {
		return nil, err
	}
	return chain.NewSpec(specData)
}

// loadSpecData reads the YAML configuration file from the given path using Viper,
// unmarshals it into a SpecData, and then validates that all required fields are set.
func loadSpecData(path string) (*chain.SpecData, error) {
	v := viper.New()
	v.SetConfigFile(path)

	// Tell Viper we're using toml.
	v.SetConfigType("toml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	// List of required keys as defined by your mapstructure tags.
	requiredKeys := []string{
		"max-effective-balance",
		"ejection-balance",
		"effective-balance-increment",
		"hysteresis-quotient",
		"hysteresis-downward-multiplier",
		"hysteresis-upward-multiplier",
		"slots-per-epoch",
		"slots-per-historical-root",
		"min-epochs-to-inactivity-penalty",
		"domain-type-beacon-proposer",
		"domain-type-beacon-attester",
		"domain-type-randao",
		"domain-type-deposit",
		"domain-type-voluntary-exit",
		"domain-type-selection-proof",
		"domain-type-aggregate-and-proof",
		"domain-type-application-mask",
		"deposit-contract-address",
		"max-deposits-per-block",
		"deposit-eth1-chain-id",
		"eth1-follow-distance",
		"target-seconds-per-eth1-block",
		"genesis-time",
		"deneb-one-fork-time",
		"electra-fork-time",
		"epochs-per-historical-vector",
		"epochs-per-slashings-vector",
		"historical-roots-limit",
		"validator-registry-limit",
		"max-withdrawals-per-payload",
		"max-validators-per-withdrawals-sweep",
		"min-epochs-for-blobs-sidecars-request",
		"max-blob-commitments-per-block",
		"max-blobs-per-block",
		"field-elements-per-blob",
		"bytes-per-blob",
		"kzg-commitment-inclusion-proof-depth",
		"validator-set-cap",
		"evm-inflation-address",
		"evm-inflation-per-block",
		"evm-inflation-address-deneb-one",
		"evm-inflation-per-block-deneb-one",
	}

	// Check if all required keys are set in the config.
	for _, key := range requiredKeys {
		if !v.IsSet(key) {
			return nil, fmt.Errorf("missing required configuration key: %s", key)
		}
	}

	var specData chain.SpecData

	// Define a decode hook to convert hex string to ExecutionAddress.
	decodeHookFunc := mapstructure.ComposeDecodeHookFunc(simpleDecodeHook)
	if err := v.Unmarshal(&specData, viper.DecodeHook(decodeHookFunc)); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config into SpecData: %w", err)
	}

	return &specData, nil
}

// simpleDecodeHook is a decode hook that does two things:
//  1. Converts a string into a common.ExecutionAddress (when target type is ExecutionAddress).
//  2. Converts numeric values into a [4]byte value using bytes.FromUint32.
//     Only numeric values are allowed for the domain types.
func simpleDecodeHook(
	f reflect.Type,
	t reflect.Type,
	data interface{},
) (interface{}, error) {
	// Convert string to ExecutionAddress.
	if f.Kind() == reflect.String && t == reflect.TypeOf(common.ExecutionAddress{}) {
		s, ok := data.(string)
		if !ok {
			return nil, fmt.Errorf("expected string for ExecutionAddress but got %T", data)
		}
		// Assume NewExecutionAddressFromHex returns (common.ExecutionAddress, error)
		addr := common.NewExecutionAddressFromHex(s)
		return addr, nil
	}

	// Convert numeric values to a 4-byte domain type (common.DomainType is an alias for [4]byte).
	if t == reflect.TypeOf(bytes.B4{}) {
		var num uint64
		switch v := data.(type) {
		case int:
			num = uint64(v) // #nosec G115: Conversion is not safe but is a trusted config file.
		case int64:
			num = uint64(v) // #nosec G115: Conversion is not safe but is a trusted config file.
		case uint64:
			num = v
		case float64:
			num = uint64(v)
		default:
			return nil, fmt.Errorf("expected numeric value for [4]byte conversion, got %T", data)
		}
		// Use FromUint32 to convert the number to a little-endian [4]byte.
		// #nosec G115: Conversion is not safe but is a trusted config file.
		return bytes.FromUint32(uint32(num)), nil
	}

	return data, nil
}
