//go:build simulated

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

package execution

import (
	"fmt"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/berachain/beacon-kit/primitives/net/url"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/stretchr/testify/require"
)

// GethNode represents a test instance of a Geth node configured
// to run inside a Docker container.
type GethNode struct {
	homeDir string
	image   docker.PullImageOptions
}

// NewGethNode returns a new GethNode instance configured with the given
// home directory and Docker image options.
func NewGethNode(homeDir string, image docker.PullImageOptions) *GethNode {
	return &GethNode{
		homeDir: homeDir,
		image:   image,
	}
}

// Start launches the Geth node container using dockertest, waits until the node is ready,
// and returns the container resource and the connection URL for the Auth RPC endpoint.
func (g *GethNode) Start(t *testing.T, genesisFile string) (*Resource, *url.ConnectionURL) {
	t.Helper()

	// Create a new Docker pool.
	pool, err := dockertest.NewPool("")
	require.NoError(t, err, "failed to create Docker pool")
	require.NotNil(t, pool, "Docker pool is nil")

	// Verify that we can connect to the Docker daemon.
	err = pool.Client.Ping()
	require.NoErrorf(t, err, "could not connect to Docker: %s", err)

	// Pull the Geth image if it is not already present.
	err = pool.Client.PullImage(g.image, docker.AuthConfiguration{})
	require.NoError(t, err, "failed to pull Geth image")

	// Resolve the absolute path to the local test files.
	absPath, err := filepath.Abs("../files")
	require.NoError(t, err, "failed to determine absolute path for test files")

	// Use the passed genesisFile variable in the command.
	cmdStr := fmt.Sprintf(`
		geth init --datadir /tmp/gethdata /testdata/%s && 
		geth --http --http.addr 0.0.0.0 --http.api eth,net,web3,debug \
			 --authrpc.addr 0.0.0.0 \
			 --authrpc.jwtsecret /testing/files/jwt.hex \
			 --authrpc.vhosts '*' \
			 --datadir /tmp/gethdata \
			 --ipcpath /tmp/gethdata/geth.ipc \
			 --syncmode full \
			 --verbosity 4 \
			 --nodiscover
	`, genesisFile)

	// Run the container with custom commands that initialize and run Geth.
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: g.image.Repository,
		Tag:        g.image.Tag,
		// Override the default entrypoint to use /bin/sh so we can chain commands.
		Entrypoint: []string{"/bin/sh"},
		Cmd: []string{
			"-c",
			cmdStr,
		},
		// Expose required ports for EL RPC, Auth RPC, and P2P communication.
		ExposedPorts: []string{"8545/tcp", "8551/tcp", "30303/tcp"},
		// Bind mount the local test data and JWT files to the container.
		Mounts: []string{
			fmt.Sprintf("%s:/%s", g.homeDir, "testdata"),
			fmt.Sprintf("%s:/%s", absPath, "testing/files"),
		},
	})
	require.NoError(t, err, "failed to run Geth container")

	// Build the connection URLs for EL RPC and Auth RPC.
	elRPC, err := url.NewFromRaw("http://" + resource.GetHostPort("8545/tcp"))
	require.NoError(t, err, "failed to create EL RPC URL")
	authRPC, err := url.NewFromRaw("http://" + resource.GetHostPort("8551/tcp"))
	require.NoError(t, err, "failed to create Auth RPC URL")

	t.Logf("Auth RPC URL: %s", authRPC.String())

	// Wait until the EL RPC endpoint is available by retrying HTTP GET requests.
	err = pool.Retry(func() error {
		resp, httpErr := http.Get(elRPC.String())
		if httpErr != nil {
			return httpErr
		}
		defer resp.Body.Close()
		return nil
	})
	require.NoError(t, err, "Geth container did not become ready in time")

	return &Resource{Resource: resource}, authRPC
}

// ValidGethImage returns the default Docker image options for the Geth node.
func ValidGethImage() docker.PullImageOptions {
	return docker.PullImageOptions{
		Repository: "ethereum/client-go",
		Tag:        "latest",
	}
}

// ValidGethImageWithSimulate returns the default Docker image options for the Geth node with Simulate API
// Build references commit https://github.com/ethereum/go-ethereum/tree/2407255bb3032dc17205ba0d648270357c98b713
// TODO: Remove once https://github.com/ethereum/go-ethereum/pull/31304/files is merged.
func ValidGethImageWithSimulate() docker.PullImageOptions {
	return docker.PullImageOptions{
		Repository: "ghcr.io/berachain/geth-simulate",
		Tag:        "latest",
	}
}
