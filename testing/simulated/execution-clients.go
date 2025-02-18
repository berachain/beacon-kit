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

package simulated

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

type GethNode struct {
	homeDir string
	image   docker.PullImageOptions
}

func NewGethNode(homeDir string, image docker.PullImageOptions) *GethNode {
	return &GethNode{homeDir, image}
}

func (g *GethNode) Start(t *testing.T) (*dockertest.Resource, *url.ConnectionURL) {
	t.Helper()
	// Create pool
	pool, err := dockertest.NewPool("")
	require.NoError(t, err)
	require.NotNil(t, pool)

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	require.NoErrorf(t, err, "Could not connect to Docker: %s", err)

	// Pull the Geth image (if not present). This can speed up future runs.
	err = pool.Client.PullImage(g.image, docker.AuthConfiguration{})
	require.NoError(t, err)

	absPath, err := filepath.Abs("../files")
	require.NoError(t, err)

	// Run container with a custom Cmd that does BOTH `init` and `run`
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: g.image.Repository,
		Tag:        g.image.Tag,
		// We'll chain these commands with bash -c
		// Override the default entrypoint to be /bin/sh instead of geth:
		Entrypoint: []string{"/bin/sh"},
		Cmd: []string{
			"-c",
			`
			geth init --datadir /tmp/gethdata /testdata/eth-genesis.json && 
			geth --http --http.addr 0.0.0.0 --http.api eth,net,web3 \
				 --authrpc.addr 0.0.0.0 \
				 --authrpc.jwtsecret /testing/files/jwt.hex \
				 --authrpc.vhosts '*' \
				 --datadir /tmp/gethdata \
				 --ipcpath /tmp/gethdata/geth.ipc \
				 --syncmode full \
				 --verbosity 4
			`,
		},
		ExposedPorts: []string{"8545/tcp", "8551/tcp", "30303/tcp"},
		Mounts: []string{
			// bind-mount local testdata => container /testdata
			fmt.Sprintf("%s:/%s", g.homeDir, "testdata"),
			fmt.Sprintf("%s:/%s", absPath, "testing/files"),
		},
	})
	require.NoError(t, err)

	elRPC, err := url.NewFromRaw("http://" + resource.GetHostPort("8545/tcp"))
	require.NoError(t, err)
	authRPC, err := url.NewFromRaw("http://" + resource.GetHostPort("8551/tcp"))
	require.NoError(t, err)

	t.Log(authRPC.String())

	// Wait until the container is ready (i.e., Geth is listening on the RPC port)
	err = pool.Retry(func() error {
		//nolint:noctx // it's just a test
		resp, httpErr := http.Get(elRPC.String())
		if httpErr != nil {
			return httpErr
		}
		readerErr := resp.Body.Close()
		if readerErr != nil {
			return readerErr
		}
		return nil
	})
	require.NoError(t, err)
	return resource, authRPC
}

func ValidGethImage() docker.PullImageOptions {
	return docker.PullImageOptions{
		Repository: "ethereum/client-go",
		Tag:        "latest",
	}
}
