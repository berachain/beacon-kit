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

package flags

import (
	"fmt"
	"time"

	"github.com/berachain/beacon-kit/config"
	"github.com/spf13/cobra"
)

const (
	// Beacon Kit Root Flag.
	beaconKitRoot   = "beacon-kit."
	ShutdownTimeout = beaconKitRoot + "shutdown-timeout"

	// Builder Config.
	builderRoot           = beaconKitRoot + "payload-builder."
	SuggestedFeeRecipient = builderRoot + "suggested-fee-recipient"
	BuilderEnabled        = builderRoot + "enabled"
	BuildPayloadTimeout   = builderRoot + "payload-timeout"

	// Validator Config.
	validatorRoot = beaconKitRoot + "validator."
	Graffiti      = validatorRoot + "graffiti"

	// Engine Config.
	engineRoot              = beaconKitRoot + "engine."
	RPCDialURL              = engineRoot + "rpc-dial-url"
	RPCRetryInterval        = engineRoot + "rpc-retry-interval"
	RPCMaxRetryInterval     = engineRoot + "rpc-max-retry-interval"
	RPCTimeout              = engineRoot + "rpc-timeout"
	RPCStartupCheckInterval = engineRoot + "rpc-startup-check-interval"
	RPCHealthCheckInteval   = engineRoot + "rpc-health-check-interval"
	RPCJWTRefreshInterval   = engineRoot + "rpc-jwt-refresh-interval"
	JWTSecretPath           = engineRoot + "jwt-secret-path"

	// KZG Config.
	kzgRoot             = beaconKitRoot + "kzg."
	KZGTrustedSetupPath = kzgRoot + "trusted-setup-path"
	KZGImplementation   = kzgRoot + "implementation"

	// Logger Config.
	loggerRoot = beaconKitRoot + "logger."
	TimeFormat = loggerRoot + "time-format"
	LogLevel   = loggerRoot + "log-level"
	Style      = loggerRoot + "style"

	// Block Store Service Config.
	blockStoreServiceRoot               = beaconKitRoot + "block-store-service."
	BlockStoreServiceEnabled            = blockStoreServiceRoot + "enabled"
	BlockStoreServiceAvailabilityWindow = blockStoreServiceRoot +
		"availability-window"

	// Node API Config.
	nodeAPIRoot    = beaconKitRoot + "node-api."
	NodeAPIEnabled = nodeAPIRoot + "enabled"
	NodeAPIAddress = nodeAPIRoot + "address"
	NodeAPILogging = nodeAPIRoot + "logging"

	// BLS Config.
	PrivValidatorKeyFile   = "priv_validator_key_file"
	PrivValidatorStateFile = "priv_validator_state_file"

	// CometBFT Config.
	cometbftRoot     = beaconKitRoot + "cometbft."
	SBTUpgradeHeight = cometbftRoot + "sbt-upgrade-height"
	SBTUpgradeTime   = cometbftRoot + "sbt-upgrade-time"
)

// Custom type to handle time parsing
type TimeFlag struct {
	Time time.Time
}

// Implement the pflag.Value interface
func (t *TimeFlag) Set(value string) error {
	parsedTime, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return fmt.Errorf("invalid time format: %w", err)
	}
	t.Time = parsedTime
	return nil
}

func (t *TimeFlag) String() string { return t.Time.Format(time.RFC3339Nano) }
func (t *TimeFlag) Type() string   { return "time" }

// AddBeaconKitFlags implements servertypes.ModuleInitFlags interface.
func AddBeaconKitFlags(startCmd *cobra.Command) {
	defaultCfg := config.DefaultConfig()
	startCmd.Flags().Duration(
		ShutdownTimeout,
		defaultCfg.ShutdownTimeout,
		"maximum time to wait for the node to gracefully shutdown before forcing an exit",
	)
	startCmd.Flags().String(
		JWTSecretPath,
		defaultCfg.Engine.JWTSecretPath,
		"path to the execution client secret",
	)
	startCmd.Flags().String(
		RPCDialURL, defaultCfg.Engine.RPCDialURL.String(), "rpc dial url",
	)
	startCmd.Flags().Duration(
		RPCRetryInterval, defaultCfg.Engine.RPCRetryInterval, "initial rpc retry interval",
	)
	startCmd.Flags().Duration(
		RPCMaxRetryInterval, defaultCfg.Engine.RPCMaxRetryInterval, "max rpc retry interval",
	)
	startCmd.Flags().Duration(
		RPCTimeout, defaultCfg.Engine.RPCTimeout, "rpc timeout",
	)
	startCmd.Flags().Duration(
		RPCStartupCheckInterval,
		defaultCfg.Engine.RPCStartupCheckInterval,
		"rpc startup check interval",
	)
	startCmd.Flags().Duration(
		RPCJWTRefreshInterval,
		defaultCfg.Engine.RPCJWTRefreshInterval,
		"rpc jwt refresh interval",
	)
	startCmd.Flags().Bool(
		BuilderEnabled,
		defaultCfg.PayloadBuilder.Enabled,
		"payload builder enabled",
	)
	startCmd.Flags().Duration(
		BuildPayloadTimeout,
		defaultCfg.PayloadBuilder.PayloadTimeout,
		"payload builder timeout",
	)
	startCmd.Flags().String(
		SuggestedFeeRecipient,
		defaultCfg.PayloadBuilder.SuggestedFeeRecipient.Hex(),
		"suggested fee recipient",
	)
	startCmd.Flags().String(
		KZGTrustedSetupPath,
		defaultCfg.KZG.TrustedSetupPath,
		"kzg trusted setup path",
	)
	startCmd.Flags().String(
		KZGImplementation,
		defaultCfg.KZG.Implementation,
		"kzg implementation",
	)
	startCmd.Flags().String(
		TimeFormat,
		defaultCfg.Logger.TimeFormat,
		"time format",
	)
	startCmd.Flags().String(
		LogLevel,
		defaultCfg.Logger.LogLevel,
		"log level",
	)
	startCmd.Flags().String(
		Style,
		defaultCfg.Logger.Style,
		"style",
	)
	startCmd.Flags().Bool(
		BlockStoreServiceEnabled,
		defaultCfg.BlockStoreService.Enabled,
		"block service enabled",
	)
	startCmd.Flags().Int(
		BlockStoreServiceAvailabilityWindow,
		defaultCfg.BlockStoreService.AvailabilityWindow,
		"block service availability window",
	)
	startCmd.Flags().Bool(
		NodeAPIEnabled,
		defaultCfg.NodeAPI.Enabled,
		"node api enabled",
	)
	startCmd.Flags().String(
		NodeAPIAddress,
		defaultCfg.NodeAPI.Address,
		"node api address",
	)
	startCmd.Flags().Bool(
		NodeAPILogging,
		defaultCfg.NodeAPI.Logging,
		"node api logging",
	)
	startCmd.Flags().Int64(
		SBTUpgradeHeight, 0, "height at which SBT is enabled",
	)
	startCmd.Flags().Var(
		&TimeFlag{}, SBTUpgradeTime, "time at which SBT is enabled (block's timestamp)",
	)
}
