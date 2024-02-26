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
	"net/url"
	"time"

	"github.com/itsdevbear/bolaris/config/flags"
	"github.com/itsdevbear/bolaris/io/cli/parser"
)

// Engine conforms to the BeaconKitConfig interface.
var _ BeaconKitConfig[Engine] = &Engine{}

// DefaultEngineConfig is the default configuration for the engine client.
func DefaultEngineConfig() Engine {
	return Engine{
		RPCDialURL: &url.URL{
			Scheme: "http",
			Host:   "localhost:8551",
		},
		RPCRetries:              3,                //nolint:gomnd // default config.
		RPCTimeout:              2 * time.Second,  //nolint:gomnd // default config.
		RPCStartupCheckInterval: 5 * time.Second,  //nolint:gomnd // default config.
		RPCHealthCheckInterval:  5 * time.Second,  //nolint:gomnd // default config.
		RPCJWTRefreshInterval:   30 * time.Second, //nolint:gomnd // default config.
		JWTSecretPath:           "./jwt.hex",
		RequiredChainID:         7, //nolint:gomnd // default config.
	}
}

// Engine is the configuration struct for the execution client.
type Engine struct {
	// RPCDialURL is the HTTP url of the execution client JSON-RPC endpoint.
	RPCDialURL *url.URL
	// RPCRetries is the number of retries before shutting down consensus
	// client.
	RPCRetries uint64
	// RPCTimeout is the RPC timeout for execution client calls.
	RPCTimeout time.Duration
	// RPCStartupCheckInterval is the Interval for the startup check.
	RPCStartupCheckInterval time.Duration
	// HealthCheckInterval is the Interval for the health check.
	RPCHealthCheckInterval time.Duration
	// JWTRefreshInterval is the Interval for the JWT refresh.
	RPCJWTRefreshInterval time.Duration
	// JWTSecretPath is the path to the JWT secret.
	JWTSecretPath string
	// RequiredChainID is the chain id that the consensus client must be
	// connected to.
	RequiredChainID uint64
}

// Parse parses the configuration.
func (c Engine) Parse(parser parser.AppOptionsParser) (*Engine, error) {
	var err error
	if c.RPCDialURL, err = parser.GetURL(flags.RPCDialURL); err != nil {
		return nil, err
	}
	if c.RPCRetries, err = parser.GetUint64(flags.RPCRetries); err != nil {
		return nil, err
	}
	if c.RPCTimeout, err = parser.GetTimeDuration(
		flags.RPCTimeout,
	); err != nil {
		return nil, err
	}
	if c.RPCStartupCheckInterval, err = parser.GetTimeDuration(
		flags.RPCStartupCheckInterval,
	); err != nil {
		return nil, err
	}
	if c.RPCHealthCheckInterval, err = parser.GetTimeDuration(
		flags.RPCHealthCheckInteval,
	); err != nil {
		return nil, err
	}
	if c.RPCJWTRefreshInterval, err = parser.GetTimeDuration(
		flags.RPCJWTRefreshInterval,
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

# Interval for the health check.
rpc-health-check-interval = "{{ .BeaconKit.Engine.RPCHealthCheckInterval }}"

# Interval for the JWT refresh.
rpc-jwt-refresh-interval = "{{ .BeaconKit.Engine.RPCJWTRefreshInterval }}"

# Path to the execution client JWT-secret
jwt-secret-path = "{{.BeaconKit.Engine.JWTSecretPath}}"

# Required chain id for the execution client.
required-chain-id = "{{.BeaconKit.Engine.RequiredChainID}}"
`
}
