package injectedconsensus

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

type gethNode struct {
	homeDir string
	image   docker.PullImageOptions
}

// NewGethNode returns interface to avoid direct access to the concrete type in tests.
func NewGethNode(homeDir string, image docker.PullImageOptions) ExecutionClient {
	return &gethNode{homeDir, image}
}

func (g *gethNode) Start(t *testing.T) (*dockertest.Resource, *url.ConnectionURL) {
	t.Helper()
	// Create pool
	pool, err := dockertest.NewPool("")
	require.NoError(t, err)

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
		Repository: "ethereum/client-go",
		Tag:        "latest",
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
				 --verbosity 3
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
		resp.Body.Close()
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
