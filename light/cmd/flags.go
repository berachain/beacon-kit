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

package cmd

import "time"

// Flag Names.
const (
	listenAddr         = "listening-address"
	engineURL          = "engine"
	primaryAddr        = "primary-addr"
	witnessAddrsJoined = "witness-addr"
	dir                = "dir"
	maxOpenConnections = "max-open-connections"

	sequential     = "sequential-verification"
	trustingPeriod = "trust-period"
	trustedHeight  = "height"
	trustedHash    = "hash"
	trustLevel     = "trust-level"

	jwtSecretPath = "jwt-secret"
	logLevel      = "log-level"
)

// Default Flag Values.
const (
	defaultListeningAddress = "tcp://localhost:26658"
	defaultEngineURL        = "http://localhost:8552"
	defaultPrimaryAddress   = "tcp://localhost:26657"
	defaultWitnessAddresses = "http://localhost:26657"
	defaultDir              = ".tmp/.beacon-light"
	defaultMaxOpenConn      = 900
	defaultTrustPeriod      = 168 * time.Hour
	defaultTrustedHeight    = 1
	defaultLogLevel         = "info"
	defaultTrustLevel       = "1/3"
	defaultSequential       = false
	defaultJWTSecretPath    = "./beacond/jwt.hex"
)

// Flag Descriptions.
const (
	listenAddrDesc         = "serve the proxy on the given address"
	engineURLDesc          = "connect to the execution client at this address"
	primaryAddrDesc        = "connect to a beacond node at this address"
	witnessAddrsJoinedDesc = `beacond nodes to cross-check the primary node,
	comma-separated`
	dirDesc                = "specify the directory"
	maxOpenConnectionsDesc = `maximum number of simultaneous connections
	(including WebSocket)`
	trustingPeriodDesc = `trusting period that headers can be verified within.
	Should be significantly less than the unbonding period`
	trustedHeightDesc = "Trusted header's height"
	trustedHashDesc   = "Trusted header's hash"
	logLevelDesc      = "Log level, info or debug (Default: info) "
	trustLevelDesc    = "trust level. Must be between 1/3 and 3/3"
	sequentialDesc    = `sequential verification.
	Verify all headers sequentially as opposed to using skipping verification`
	jwtSecretDesc = "Path to the JWT secret file"
)

// Log Level Flags.
const (
	logLevelInfo  = "info"
	logLevelDebug = "debug"
)
