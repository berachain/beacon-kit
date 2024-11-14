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

package flags

import (
	"github.com/berachain/beacon-kit/mod/config"
	"github.com/spf13/cobra"
)

const (
	// Beacon Kit Root Flag.
	beaconKitRoot      = "beacon-kit."
	BeaconKitAcceptTos = beaconKitRoot + "accept-tos"

	// Builder Config.
	builderRoot              = beaconKitRoot + "payload-builder."
	SuggestedFeeRecipient    = builderRoot + "suggested-fee-recipient"
	LocalBuilderEnabled      = builderRoot + "local-builder-enabled"
	LocalBuildPayloadTimeout = builderRoot + "local-build-payload-timeout"

	// Validator Config.
	validatorRoot = beaconKitRoot + "validator."
	Graffiti      = validatorRoot + "graffiti"

	// Engine Config.
	engineRoot              = beaconKitRoot + "engine."
	RPCDialURL              = engineRoot + "rpc-dial-url"
	RPCRetries              = engineRoot + "rpc-retries"
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
	BlockStoreServiceAvailabilityWindow = blockStoreServiceRoot + "availability-window"

	// Node API Config.
	nodeAPIRoot    = beaconKitRoot + "node-api."
	NodeAPIEnabled = nodeAPIRoot + "enabled"
	NodeAPIAddress = nodeAPIRoot + "address"
	NodeAPILogging = nodeAPIRoot + "logging"
)

// AddBeaconKitFlags adds Beacon Kit flags to the given command.
func AddBeaconKitFlags(startCmd *cobra.Command) {
	defaultCfg := config.DefaultConfig()
	startCmd.Flags().String(
		JWTSecretPath,
		defaultCfg.Engine.JWTSecretPath,
		"path to the execution client secret",
	)
	startCmd.Flags().String(
		RPCDialURL, defaultCfg.Engine.RPCDialURL.String(), "RPC dial URL",
	)
	startCmd.Flags().Uint64(
		RPCRetries, defaultCfg.Engine.RPCRetries, "number of RPC retries",
	)
	startCmd.Flags().Duration(
		RPCTimeout, defaultCfg.Engine.RPCTimeout, "RPC timeout duration",
	)
	startCmd.Flags().Duration(
		RPCStartupCheckInterval,
		defaultCfg.Engine.RPCStartupCheckInterval,
		"RPC startup check interval duration",
	)
	startCmd.Flags().Duration(
		RPCJWTRefreshInterval,
		defaultCfg.Engine.RPCJWTRefreshInterval,
		"RPC JWT refresh interval duration",
	)
	startCmd.Flags().String(
		SuggestedFeeRecipient,
		defaultCfg.PayloadBuilder.SuggestedFeeRecipient.Hex(),
		"suggested fee recipient address",
	)
	startCmd.Flags().String(
		KZGTrustedSetupPath,
		defaultCfg.KZG.TrustedSetupPath,
		"path to the KZG trusted setup",
	)
	startCmd.Flags().String(
		KZGImplementation,
		defaultCfg.KZG.Implementation,
		"KZG implementation type",
	)
	startCmd.Flags().String(
		TimeFormat,
		defaultCfg.Logger.TimeFormat,
		"timestamp format for logging",
	)
	startCmd.Flags().String(
		LogLevel,
		defaultCfg.Logger.LogLevel,
		"logging level",
	)
	startCmd.Flags().String(
		Style,
		defaultCfg.Logger.Style,
		"logging output style",
	)
	startCmd.Flags().Bool(
		BlockStoreServiceEnabled,
		defaultCfg.BlockStoreService.Enabled,
		"enable block store service",
	)
	startCmd.Flags().Int(
		BlockStoreServiceAvailabilityWindow,
		defaultCfg.BlockStoreService.AvailabilityWindow,
		"block store service availability window",
	)
	startCmd.Flags().Bool(
		NodeAPIEnabled,
		defaultCfg.NodeAPI.Enabled,
		"enable node API",
	)
	startCmd.Flags().String(
		NodeAPIAddress,
		defaultCfg.NodeAPI.Address,
		"address for the node API",
	)
	startCmd.Flags().Bool(
		NodeAPILogging,
		defaultCfg.NodeAPI.Logging,
		"enable node API logging",
	)
}
