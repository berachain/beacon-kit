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

package cmd

import (
	"path/filepath"
	"time"

	"github.com/itsdevbear/bolaris/io/file"
	"github.com/spf13/cobra"
)

// beacond light cosmoshub-3 --primary-addr http://0.0.0.0:8545
// --witness-addr http://144.76.61.201:26657/ --trusted-height 5940895
// --trusted-hash 8663FBD3FB9DCE3D8E461EA521C38256F6EAF85D4FA492BAE26D5863F53CA15

const (
	shortDescription = "Run a light client proxy server, verifying CometBFT rpc"
	longDescription  = `Run a light client proxy server, verifying CometBFT rpc.

	All calls that can be tracked back to a block header by a proof
	will be verified before passing them back to the caller. Other than
	that, it will present the same interface as a full CometBFT node.
	
	Furthermore to the chainID, a fresh instance of a light client will
	need a primary RPC address, a trusted hash and height and witness RPC addresses
	(if not using sequential verification). To restart the node, thereafter
	only the chainID is required.
	`
	example = `light cosmoshub-4 -primary-addr http://52.57.29.196:26657
	-witness-addr http://public-seed-node.cosmoshub.certus.one:26657
	--height 962118 --hash 28B97BE9F6DE51AC69F70E0B7BFD7E5C9CD1A595B7DC31AFF27C50D4948020CD`
)

const (
	// Flag Names.
	listenAddr         = "listening-address"
	primaryAddr        = "primary-addr"
	witnessAddrsJoined = "witness-addr"
	dir                = "dir"
	maxOpenConnections = "max-open-connections"

	sequential     = "sequential-verification"
	trustingPeriod = "trust-period"
	trustedHeight  = "trusted-height"
	trustedHash    = "trusted-hash"
	trustLevel     = "trust-level"

	logLevel = "log-level"

	// Default Flag Values.
	defaultListeningAddress = "tcp://localhost:8888"
	defaultPrimaryAddress   = ""
	defaultWitnessAddresses = "http://localhost:26657"
	defaultDir              = ".beacon-light"
	defaultMaxOpenConn      = 900
	defaultTrustPeriod      = 168 * time.Hour
	defaultTrustedHeight    = 1
	defaultLogLevel         = "info"
	defaultTrustLevel       = "1/3"
	defaultSequential       = false

	// Flag Descriptions.
	listenAddrDesc         = "serve the proxy on the given address"
	primaryAddrDesc        = "connect to a Tendermint node at this address"
	witnessAddrsJoinedDesc = "tendermint nodes to cross-check the primary node, comma-separated"
	dirDesc                = "specify the directory"
	maxOpenConnectionsDesc = "maximum number of simultaneous connections (including WebSocket)."
	trustingPeriodDesc     = `trusting period that headers can be verified within.
	Should be significantly less than the unbonding period`
	trustedHeightDesc = "Trusted header's height"
	trustedHashDesc   = "Trusted header's hash"
	logLevelDesc      = "Log level, info or debug (Default: info) "
	trustLevelDesc    = "trust level. Must be between 1/3 and 3/3"
	sequentialDesc    = `sequential verification.
	Verify all headers sequentially as opposed to using skipping verification`
)

func LightClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "light [chainID]",
		Short:   shortDescription,
		Long:    longDescription,
		RunE:    runProxy,
		Args:    cobra.ExactArgs(1),
		Example: example,
	}

	homeDir, err := file.HomeDir()
	if err != nil {
		panic(err)
	}

	cmd.Flags().String(listenAddr, defaultListeningAddress, listenAddrDesc)
	cmd.Flags().String(primaryAddr, defaultPrimaryAddress, primaryAddrDesc)
	cmd.Flags().String(witnessAddrsJoined, defaultWitnessAddresses, witnessAddrsJoinedDesc)
	cmd.Flags().String(dir, filepath.Join(homeDir, defaultDir), dirDesc)
	cmd.Flags().Int(maxOpenConnections, defaultMaxOpenConn, maxOpenConnectionsDesc)
	cmd.Flags().Duration(trustingPeriod, defaultTrustPeriod, trustingPeriodDesc)
	cmd.Flags().Int64(trustedHeight, defaultTrustedHeight, trustedHeightDesc)
	cmd.Flags().BytesHex(trustedHash, []byte{}, trustedHashDesc)
	cmd.Flags().String(logLevel, defaultLogLevel, logLevelDesc)
	cmd.Flags().String(trustLevel, defaultTrustLevel, trustLevelDesc)
	cmd.Flags().Bool(sequential, defaultSequential, sequentialDesc)

	return cmd
}
