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

package template

//nolint:lll // template.
const TomlTemplate = `
###############################################################################
###                                BeaconKit                                ###
###############################################################################

[beacon-kit.engine]
# HTTP url of the execution client JSON-RPC endpoint.
rpc-dial-url = "{{ .BeaconKit.Engine.RPCDialURL }}"

# Number of retries before shutting down consensus client.
rpc-retries = "{{.BeaconKit.Engine.RPCRetries}}"

# RPC timeout for execution client requests.
rpc-timeout = "{{ .BeaconKit.Engine.RPCTimeout }}"

# Interval for the startup check.
rpc-startup-check-interval = "{{ .BeaconKit.Engine.RPCStartupCheckInterval }}"

# Interval for the JWT refresh.
rpc-jwt-refresh-interval = "{{ .BeaconKit.Engine.RPCJWTRefreshInterval }}"

# Path to the execution client JWT-secret
jwt-secret-path = "{{.BeaconKit.Engine.JWTSecretPath}}"

[beacon-kit.logger]
# TimeFormat is a string that defines the format of the time in the logger.
time-format = "{{.BeaconKit.Logger.TimeFormat}}"

# LogLevel is the level of logging. Logger will log messages with verbosity up 
# to LogLevel.
log-level = "{{.BeaconKit.Logger.LogLevel}}"

# Style is the style of the logger.
style = "{{.BeaconKit.Logger.Style}}"

[beacon-kit.kzg]
# Path to the trusted setup path.
trusted-setup-path = "{{.BeaconKit.KZG.TrustedSetupPath}}"

# KZG implementation to use.
# Options are "crate-crypto/go-kzg-4844" or "ethereum/c-kzg-4844".
implementation = "{{.BeaconKit.KZG.Implementation}}"

[beacon-kit.payload-builder]
# Enabled determines if the local payload builder is enabled.
enabled = {{ .BeaconKit.PayloadBuilder.Enabled }}

# Post bellatrix, this address will receive the transaction fees produced by any blocks 
# from this node.
suggested-fee-recipient = "{{.BeaconKit.PayloadBuilder.SuggestedFeeRecipient}}"

# The timeout for local build payload. This should match, or be slightly less
# than the configured timeout on your execution client. It also must be less than
# timeout_proposal in the CometBFT configuration.
payload-timeout = "{{ .BeaconKit.PayloadBuilder.PayloadTimeout }}"

[beacon-kit.validator]
# Graffiti string that will be included in the graffiti field of the beacon block.
graffiti = "{{.BeaconKit.Validator.Graffiti}}"

# EnableOptimisticPayloadBuilds enables building the next block's payload optimistically in
# process-proposal to allow for the execution client to have more time to assemble the block.
enable-optimistic-payload-builds = "{{.BeaconKit.Validator.EnableOptimisticPayloadBuilds}}"

[beacon-kit.block-service]
# Enabled determines if the block service is enabled.
enabled = "{{ .BeaconKit.BlockService.Enabled }}"

# PrunerEnabled determines if the block pruner is enabled.
pruner-enabled = "{{ .BeaconKit.BlockService.PrunerEnabled }}"

# AvailabilityWindow is the number of slots to keep in the store.
availability-window = "{{ .BeaconKit.BlockService.AvailabilityWindow }}"
`
