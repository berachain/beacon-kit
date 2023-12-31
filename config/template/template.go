// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
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

package template

const (
	ConfigTemplate = `
###############################################################################
###                                 Polaris                                 ###
###############################################################################
# General Polaris settings
[polaris]

[beacon-kit.execution-client]
# HTTP url of the execution client JSON-RPC endpoint.
rpc-dial-url = "{{ .Polaris.ExecutionClient.RPCDialURL }}"

# RPC timeout for execution client requests.
rpc-timeout = "{{ .Polaris.ExecutionClient.RPCTimeout }}"

# Number of retries before shutting down consensus client.
rpc-retries = "{{.Polaris.ExecutionClient.RPCRetries}}"

# Path to the execution client JWT-secret
jwt-secret-path = "{{.Polaris.ExecutionClient.JWTSecretPath}}"

[beacon-kit.beacon-config]
# Altair fork epoch
altair-fork-epoch = {{.Polaris.BeaconConfig.AltairForkEpoch}}

# Bellatrix fork epoch
bellatrix-fork-epoch = {{.Polaris.BeaconConfig.BellatrixForkEpoch}}

# Capella fork epoch
capella-fork-epoch = {{.Polaris.BeaconConfig.CapellaForkEpoch}}

# Deneb fork epoch
deneb-fork-epoch = {{.Polaris.BeaconConfig.DenebForkEpoch}}
`
)
