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

package config

import (
	"github.com/itsdevbear/bolaris/config/flags"
	"github.com/itsdevbear/bolaris/config/parser"
)

// Execution conforms to the BeaconKitConfig interface.
var _ BeaconKitConfig[Execution] = &Execution{}

// DefaultExecutionConfig returns the default configuration for the execution client.
func DefaultExecutionConfig() Execution {
	return Execution{
		RPCDialURL:      "http://localhost:8551",
		RPCTimeout:      5, //nolint:gomnd // default config.
		RPCRetries:      3, //nolint:gomnd // default config.
		JWTSecretPath:   "./jwt.hex",
		RequiredChainID: 7, //nolint:gomnd // default config.
	}
}

// Execution is the configuration struct for the execution client.
type Execution struct {
	// RPCDialURL is the HTTP url of the execution client JSON-RPC endpoint.
	RPCDialURL string
	// RPCTimeout is the RPC timeout for execution client requests.
	RPCTimeout uint64
	// RPCRetries is the number of retries before shutting down consensus client.
	RPCRetries uint64
	// JWTSecretPath is the path to the JWT secret.
	JWTSecretPath string
	// RequiredChainID is the chain id that the consensus client must be connected to.
	RequiredChainID uint64
}

// Parse parses the configuration.
func (c Execution) Parse(parser parser.AppOptionsParser) (*Execution, error) {
	var err error
	if c.RPCDialURL, err = parser.GetString(flags.RPCDialURL); err != nil {
		return nil, err
	}
	if c.RPCRetries, err = parser.GetUint64(flags.RPCRetries); err != nil {
		return nil, err
	}
	if c.RPCTimeout, err = parser.GetUint64(
		flags.RPCTimeout,
	); err != nil {
		return nil, err
	}
	if c.JWTSecretPath, err = parser.GetString(
		flags.JWTSecretPath,
	); err != nil {
		return nil, err
	}
	if c.RequiredChainID, err = parser.GetUint64(
		flags.RequiredChainID,
	); err != nil {
		return nil, err
	}
	return &c, nil
}

func (c Execution) Template() string {
	return `
[beacon-kit.execution-client]
# HTTP url of the execution client JSON-RPC endpoint.
rpc-dial-url = "{{ .BeaconKit.Execution.RPCDialURL }}"

# RPC timeout for execution client requests.
rpc-timeout = "{{ .BeaconKit.Execution.RPCTimeout }}"

# Number of retries before shutting down consensus client.
rpc-retries = "{{.BeaconKit.Execution.RPCRetries}}"

# Path to the execution client JWT-secret
jwt-secret-path = "{{.BeaconKit.Execution.JWTSecretPath}}"

# Required chain id for the execution client.
required-chain-id = "{{.BeaconKit.Execution.RequiredChainID}}"
`
}
