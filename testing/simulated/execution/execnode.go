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

// ExecNode represents a test instance of an execution client node running inside a Docker container.
// The CmdStrBuilder field is used to generate the command string that initializes and runs the node.
type ExecNode struct {
	homeDir       string
	image         docker.PullImageOptions
	CmdStrBuilder func(genesisFile string) string
}

// NewExecNode returns a new ExecNode instance configured with the given home directory,
// Docker image options, and a command string builder.
func NewExecNode(homeDir string, image docker.PullImageOptions, builder func(genesisFile string) string) *ExecNode {
	return &ExecNode{
		homeDir:       homeDir,
		image:         image,
		CmdStrBuilder: builder,
	}
}

// Start launches the execution client container using dockertest, waits until the client is ready,
// and returns the container resource along with the connection URL for the Auth RPC endpoint.
func (e *ExecNode) Start(t *testing.T, genesisFile string) (*Resource, *url.ConnectionURL, *url.ConnectionURL) {
	t.Helper()

	// Create a new Docker pool.
	pool, err := dockertest.NewPool("")
	require.NoError(t, err, "failed to create Docker pool")
	require.NotNil(t, pool, "Docker pool is nil")

	// Verify that we can connect to the Docker daemon.
	err = pool.Client.Ping()
	require.NoErrorf(t, err, "could not connect to Docker: %s", err)

	// Pull the image if it is not already present.
	err = pool.Client.PullImage(e.image, docker.AuthConfiguration{})
	require.NoError(t, err, "failed to pull image")

	// Resolve the absolute path to the local test files.
	absPath, err := filepath.Abs("../files")
	require.NoError(t, err, "failed to determine absolute path for test files")

	require.NotNil(t, e.CmdStrBuilder, "CmdStrBuilder is nil")
	cmdStr := e.CmdStrBuilder(genesisFile)

	// Run the container with custom commands that initialize and run the client.
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: e.image.Repository,
		Tag:        e.image.Tag,
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
			fmt.Sprintf("%s:/%s", e.homeDir, "testdata"),
			fmt.Sprintf("%s:/%s", absPath, "testing/files"),
		},
	})
	require.NoError(t, err, "failed to run container")

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
	require.NoError(t, err, "Container did not become ready in time")

	return &Resource{Resource: resource}, authRPC, elRPC
}
