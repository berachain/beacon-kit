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
	"fmt"
	"math"

	"github.com/spf13/viper"

	pruningtypes "cosmossdk.io/store/pruning/types"

	"github.com/cosmos/cosmos-sdk/telemetry"
)

const (
	defaultMinGasPrices = ""

	// DefaultAPIAddress defines the default address to bind the API server to.
	DefaultAPIAddress = "tcp://localhost:1317"

	// DefaultGRPCAddress defines the default address to bind the gRPC server to.
	DefaultGRPCAddress = "localhost:9090"

	// DefaultGRPCMaxRecvMsgSize defines the default gRPC max message size in
	// bytes the server can receive.
	DefaultGRPCMaxRecvMsgSize = 1024 * 1024 * 10

	// DefaultGRPCMaxSendMsgSize defines the default gRPC max message size in
	// bytes the server can send.
	DefaultGRPCMaxSendMsgSize = math.MaxInt32
)

// BaseConfig defines the server's basic configuration
type BaseConfig struct {
	Pruning           string `mapstructure:"pruning"`
	PruningKeepRecent string `mapstructure:"pruning-keep-recent"`
	PruningInterval   string `mapstructure:"pruning-interval"`

	// HaltHeight contains a non-zero block height at which a node will gracefully
	// halt and shutdown that can be used to assist upgrades and testing.
	//
	// Note: Commitment of state will be attempted on the corresponding block.
	HaltHeight uint64 `mapstructure:"halt-height"`

	// HaltTime contains a non-zero minimum block time (in Unix seconds) at which
	// a node will gracefully halt and shutdown that can be used to assist
	// upgrades and testing.
	//
	// Note: Commitment of state will be attempted on the corresponding block.
	HaltTime uint64 `mapstructure:"halt-time"`

	// MinRetainBlocks defines the minimum block height offset from the current
	// block being committed, such that blocks past this offset may be pruned
	// from CometBFT. It is used as part of the process of determining the
	// ResponseCommit.RetainHeight value during ABCI Commit. A value of 0 indicates
	// that no blocks should be pruned.
	//
	// This configuration value is only responsible for pruning CometBFT blocks.
	// It has no bearing on application state pruning which is determined by the
	// "pruning-*" configurations.
	//
	// Note: CometBFT block pruning is dependent on this parameter in conjunction
	// with the unbonding (safety threshold) period, state pruning and state sync
	// snapshot parameters to determine the correct minimum value of
	// ResponseCommit.RetainHeight.
	MinRetainBlocks uint64 `mapstructure:"min-retain-blocks"`

	// InterBlockCache enables inter-block caching.
	InterBlockCache bool `mapstructure:"inter-block-cache"`

	// IndexEvents defines the set of events in the form {eventType}.{attributeKey},
	// which informs CometBFT what to index. If empty, all events will be indexed.
	IndexEvents []string `mapstructure:"index-events"`

	// IavlCacheSize set the size of the iavl tree cache.
	IAVLCacheSize uint64 `mapstructure:"iavl-cache-size"`

	// IAVLDisableFastNode enables or disables the fast sync node.
	IAVLDisableFastNode bool `mapstructure:"iavl-disable-fastnode"`
}

// Config defines the server's top level configuration
type Config struct {
	BaseConfig `mapstructure:",squash"`

	// Telemetry defines the application telemetry configuration
	Telemetry telemetry.Config `mapstructure:"telemetry"`
}

// DefaultConfig returns server's default configuration.
func DefaultConfig() *Config {
	return &Config{
		BaseConfig: BaseConfig{
			InterBlockCache:     true,
			Pruning:             pruningtypes.PruningOptionDefault,
			PruningKeepRecent:   "0",
			PruningInterval:     "0",
			MinRetainBlocks:     0,
			IndexEvents:         make([]string, 0),
			IAVLCacheSize:       781250,
			IAVLDisableFastNode: false,
		},
		Telemetry: telemetry.Config{
			Enabled:      false,
			GlobalLabels: [][]string{},
		},
	}
}

// GetConfig returns a fully parsed Config object.
func GetConfig(v *viper.Viper) (Config, error) {
	conf := DefaultConfig()
	if err := v.Unmarshal(conf); err != nil {
		return Config{}, fmt.Errorf("error extracting app config: %w", err)
	}
	return *conf, nil
}

// ValidateBasic returns an error if min-gas-prices field is empty in BaseConfig. Otherwise, it returns nil.
func (c Config) ValidateBasic() error {
	// if c.Pruning == pruningtypes.PruningOptionEverything && c.StateSync.SnapshotInterval > 0 {
	// 	return sdkerrors.ErrAppConfig.Wrapf(
	// 		"cannot enable state sync snapshots with '%s' pruning setting", pruningtypes.PruningOptionEverything,
	// 	)
	// }

	return nil
}
