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

package client

import (
	"net/url"
	"time"

	"github.com/berachain/beacon-kit/mod/node-builder/config/flags"
	"github.com/berachain/beacon-kit/mod/node-builder/utils/cli/parser"
)

const (
	defaultDialURL                 = "http://localhost:8551"
	defaultRPCRetries              = 3
	defaultRPCTimeout              = 2 * time.Second
	defaultRPCStartupCheckInterval = 3 * time.Second
	defaultRPCJWTRefreshInterval   = 30 * time.Second
	//#nosec:G101 // false positive.
	defaultJWTSecretPath   = "./jwt.hex"
	defaultRequiredChainID = 80087
)

// DefaultConfig is the default configuration for the engine client.
func DefaultConfig() Config {
	//#nosec:G703 // ignoring on purpose since it is the default URL.
	dialURL, _ := url.Parse(defaultDialURL)
	return Config{
		RPCDialURL:              dialURL,
		RPCRetries:              defaultRPCRetries,
		RPCTimeout:              defaultRPCTimeout,
		RPCStartupCheckInterval: defaultRPCStartupCheckInterval,
		RPCJWTRefreshInterval:   defaultRPCJWTRefreshInterval,
		JWTSecretPath:           defaultJWTSecretPath,
		RequiredChainID:         defaultRequiredChainID,
	}
}

// Config is the configuration struct for the execution client.
type Config struct {
	// RPCDialURL is the HTTP url of the execution client JSON-RPC endpoint.
	RPCDialURL *url.URL
	// RPCRetries is the number of retries before shutting down consensus
	// client.
	RPCRetries uint64
	// RPCTimeout is the RPC timeout for execution client calls.
	RPCTimeout time.Duration
	// RPCStartupCheckInterval is the Interval for the startup check.
	RPCStartupCheckInterval time.Duration
	// JWTRefreshInterval is the Interval for the JWT refresh.
	RPCJWTRefreshInterval time.Duration
	// JWTSecretPath is the path to the JWT secret.
	JWTSecretPath string
	// RequiredChainID is the chain id that the consensus client must be
	// connected to.
	RequiredChainID uint64
}

// Parse parses the configuration.
func (c Config) Parse(parser parser.AppOptionsParser) (*Config, error) {
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

func (c Config) Template() string {
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
