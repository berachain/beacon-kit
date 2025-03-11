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

package client

import (
	"time"

	"github.com/berachain/beacon-kit/primitives/net/url"
)

const (
	MinRPCTimeout = 2 * time.Second

	defaultDialURL                 = "http://localhost:8551"
	defaultRPCRetries              = 0
	defaultRPCRetryInterval        = 100 * time.Millisecond
	defaultRPCMaxRetryInterval     = 10 * time.Second
	defaultRPCStartupCheckInterval = 3 * time.Second
	defaultRPCJWTRefreshInterval   = 30 * time.Second
	//#nosec:G101 // false positive.
	defaultJWTSecretPath = "./jwt.hex"
)

// DefaultConfig is the default configuration for the engine client.
func DefaultConfig() Config {
	//#nosec:G703 // ignoring on purpose since it is the default URL.
	dialURL, _ := url.NewFromRaw(defaultDialURL)
	return Config{
		RPCDialURL:              dialURL,
		RPCRetries:              defaultRPCRetries,
		RPCRetryInterval:        defaultRPCRetryInterval,
		RPCMaxRetryInterval:     defaultRPCMaxRetryInterval,
		RPCTimeout:              MinRPCTimeout,
		RPCStartupCheckInterval: defaultRPCStartupCheckInterval,
		RPCJWTRefreshInterval:   defaultRPCJWTRefreshInterval,
		JWTSecretPath:           defaultJWTSecretPath,
	}
}

// Config is the configuration struct for the execution client.
type Config struct {
	// RPCDialURL is the HTTP url of the execution client JSON-RPC endpoint.
	RPCDialURL *url.ConnectionURL `mapstructure:"rpc-dial-url"`
	// RPCRetries is the number of retries before shutting down consensus
	// client. A value of 0 will retry infinitely.
	RPCRetries uint64 `mapstructure:"rpc-retries"`
	// RPCRetryInterval is the initial RPC backoff for repeated execution client calls.
	RPCRetryInterval time.Duration `mapstructure:"rpc-retry-interval"`
	// MaxRPCRetryInterval is the maximum RPC backoff for repeated execution client calls.
	RPCMaxRetryInterval time.Duration `mapstructure:"rpc-max-retry-interval"`
	// RPCTimeout is the RPC timeout for individual execution client calls.
	RPCTimeout time.Duration `mapstructure:"rpc-timeout"`
	// RPCStartupCheckInterval is the Interval for the startup check.
	RPCStartupCheckInterval time.Duration `mapstructure:"rpc-startup-check-interval"`
	// JWTRefreshInterval is the Interval for the JWT refresh.
	RPCJWTRefreshInterval time.Duration `mapstructure:"rpc-jwt-refresh-interval"`
	// JWTSecretPath is the path to the JWT secret.
	JWTSecretPath string `mapstructure:"jwt-secret-path"`
}
