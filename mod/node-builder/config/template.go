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

//nolint:lll // it's a template file.
package config

const Template = `
###############################################################################
###                                BeaconKit                                ###
###############################################################################

[beacon-kit.abci]
# Position of the beacon block in the proposal
beacon-block-proposal-position = {{.BeaconKit.ABCI.BeaconBlockPosition}}

# Position of the blob sidecars in the proposal
blob-sidecars-block-proposal-position = {{.BeaconKit.ABCI.BlobSidecarsBlockPosition}}

[beacon-kit.beacon-chain]

########### Gwei Values ###########
# MinDepositAmount is the minimum deposit amount per deposit transaction.
min-deposit-amount = {{.BeaconKit.Beacon.MinDepositAmount}}

# MaxEffectiveBalance is the maximum effective balance allowed for a validator.
max-effective-balance = {{.BeaconKit.Beacon.MaxEffectiveBalance}}

# EffectiveBalanceIncrement is the effective balance increment.
effective-balance-increment = {{.BeaconKit.Beacon.EffectiveBalanceIncrement}}

########### Time Parameters ##########
# SlotsPerEpoch is the number of slots per epoch.
slots-per-epoch = {{.BeaconKit.Beacon.SlotsPerEpoch}}

# SlotsPerHistoricalRoot is the number of slots per historical root.
slots-per-historical-root = {{.BeaconKit.Beacon.SlotsPerHistoricalRoot}}

########### Eth1 Data ###########
# DepositContractAddress is the address of the deposit contract.
deposit-contract-address = "{{.BeaconKit.Beacon.DepositContractAddress}}"

########### Forks ###########
# Electra fork epoch
electra-fork-epoch = {{.BeaconKit.Beacon.ElectraForkEpoch}}

########### State List Lengths ###########
# EpochsPerHistoricalVector is the number of epochs in the historical vector.
epochs-per-historical-vector = {{.BeaconKit.Beacon.EpochsPerHistoricalVector}}

# EpochsPerSlashingsVector is the number of epochs in the slashings vector.
epochs-per-slashings-vector = {{.BeaconKit.Beacon.EpochsPerSlashingsVector}}

########### Max Operations ###########
# MaxDepositsPerBlock specifies the maximum number of deposit operations allowed per block.
max-deposits-per-block = {{.BeaconKit.Beacon.MaxDepositsPerBlock}}

# MaxWithdrawalsPerPayload indicates the maximum number of withdrawal operations allowed in a single payload.
max-withdrawals-per-payload = {{.BeaconKit.Beacon.MaxWithdrawalsPerPayload}}

# MaxBlobsPerBlock specifies the maximum number of blobs allowed per block.
max-blobs-per-block = {{.BeaconKit.Beacon.MaxBlobsPerBlock}}

########### Rewards and Penalties ###########
# ProportionalSlashingMultiplier is the slashing multiplier relative to the base penalty.
proportional-slashing-multiplier = {{.BeaconKit.Beacon.ProportionalSlashingMultiplier}}

########### Deneb Values ###########
# MinEpochsForBlobsSidecarsRequest is the minimum number of epochs the node will keep the blobs for.
min-epochs-for-blobs-sidecars-request = {{.BeaconKit.Beacon.MinEpochsForBlobsSidecarsRequest}}

[beacon-kit.builder]
# Post bellatrix, this address will receive the transaction fees produced by any blocks 
# from this node.
suggested-fee-recipient = "{{.BeaconKit.Builder.SuggestedFeeRecipient}}"

# Graffiti string that will be included in the graffiti field of the beacon block.
graffiti = "{{.BeaconKit.Builder.Graffiti}}"

# LocalBuilderEnabled determines if the local payload builder is enabled.
local-builder-enabled = {{ .BeaconKit.Builder.LocalBuilderEnabled }}

# The timeout for local build payload. This should match, or be slightly less
# than the configured timeout on your execution client. It also must be less than
# timeout_proposal in the CometBFT configuration.
local-build-payload-timeout = "{{ .BeaconKit.Builder.LocalBuildPayloadTimeout }}"

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

# Required chain id for the execution client.
required-chain-id = "{{.BeaconKit.Engine.RequiredChainID}}"

[beacon-kit.kzg]
# Path to the trusted setup path.
trusted-setup-path = "{{.BeaconKit.KZG.TrustedSetupPath}}"

# KZG implementation to use.
# Options are "crate-crypto/go-kzg-4844" or "ethereum/c-kzg-4844".
implementation = "{{.BeaconKit.KZG.Implementation}}"
`
