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

import (
	"github.com/itsdevbear/bolaris/config/flags"
	engineclient "github.com/itsdevbear/bolaris/engine/client"
	"github.com/itsdevbear/bolaris/io/cli/parser"
)

// Engine conforms to the BeaconKitConfig interface.
var _ BeaconKitConfig[Engine] = &Engine{}

// DefaultEngineConfig is the default configuration for the engine client.
func DefaultEngineConfig() Engine {
	return Engine{
		Config: engineclient.DefaultConfig(),
	}
}

// Engine is the configuration struct for the execution client.
type Engine struct {
	engineclient.Config
}

// Parse parses the configuration.
func (c Engine) Parse(parser parser.AppOptionsParser) (*Engine, error) {
	var err error
	if c.Config.RPCDialURL, err = parser.GetURL(flags.RPCDialURL); err != nil {
		return nil, err
	}
	if c.Config.RPCRetries, err = parser.GetUint64(flags.RPCRetries); err != nil {
		return nil, err
	}
	if c.Config.RPCTimeout, err = parser.GetTimeDuration(
		flags.RPCTimeout,
	); err != nil {
		return nil, err
	}
	if c.Config.RPCStartupCheckInterval, err = parser.GetTimeDuration(
		flags.RPCStartupCheckInterval,
	); err != nil {
		return nil, err
	}

	if c.Config.RPCJWTRefreshInterval, err = parser.GetTimeDuration(
		flags.RPCJWTRefreshInterval,
	); err != nil {
		return nil, err
	}
	if c.Config.JWTSecretPath, err = parser.GetString(
		flags.JWTSecretPath,
	); err != nil {
		return nil, err
	}
	if c.Config.RequiredChainID, err = parser.GetUint64(
		flags.RequiredChainID,
	); err != nil {
		return nil, err
	}
	return &c, nil
}

func (c Engine) Template() string {
	return `
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
`
}
