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

package config

//nolint:lll // template.
const Template = `
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

# Interval for checking client sync status
sync-check-interval = "{{ .BeaconKit.Engine.SyncCheckInterval }}"

# Path to the execution client JWT-secret
jwt-secret-path = "{{.BeaconKit.Engine.JWTSecretPath}}"

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
`
