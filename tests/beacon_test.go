package e2e_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	e2e "github.com/cometbft/cometbft/test/e2e/pkg"
	"github.com/stretchr/testify/require"
)

type BeaconHTTPClient struct {
	*http.Client
	baseURL string
}

func getNodeIP(node e2e.Node) string {
	ip := node.ExternalIP.String()
	if node.ExternalIP.To4() == nil {
		// IPv6 addresses must be wrapped in [] to avoid conflict with : port separator
		ip = fmt.Sprintf("[%v]", ip)
	}
	return ip
}

// Tests that all nodes have peered with some other nodes, regardless of discovery method.
func TestNodeAPIGenesis(t *testing.T) {
	testNode(t, func(t *testing.T, node e2e.Node) {
		t.Log("Testing node-api genesis for node", node.Name)
		t.Helper()
		// Seed nodes shouldn't necessarily mesh with the entire network.
		if node.Mode == e2e.ModeSeed || node.Mode == e2e.ModeFull {
			return
		}
		fmt.Println("Testing node-api genesis for node", node.Name)
		ip := getNodeIP(node)

		nIP := fmt.Sprintf("http://%v:350%d", ip, extractValidatorIndices(node.Name))

		fmt.Println("Connecting to beacon node at", nIP)
		bclient := &BeaconHTTPClient{
			Client: &http.Client{
				Timeout: time.Second * 10,
			},
			baseURL: nIP,
		}

		resp, err := bclient.Get(bclient.baseURL + "/eth/v1/beacon/genesis")
		if err != nil {
			t.Log(err)
			t.Fatalf("failed to get validator info")
		}
		if resp == nil {
			t.Fatalf("received nil response")
		}

	})
}

func TestStateValidator(t *testing.T) {
	testNode(t, func(t *testing.T, node e2e.Node) {
		if node.Mode == e2e.ModeSeed || node.Mode == e2e.ModeFull {
			return
		}
		ip := getNodeIP(node)

		nIP := fmt.Sprintf("http://%v:350%d", ip, extractValidatorIndices(node.Name))
		bclient := &BeaconHTTPClient{
			Client: &http.Client{
				Timeout: time.Second * 10,
			},
			baseURL: nIP,
		}
		resp, err := bclient.Get(fmt.Sprintf(bclient.baseURL+"/eth/v1/beacon/states/%s/validators/%s", "10", node.PrivvalKey.PubKey().Address().String()))
		if err != nil {
			require.NoError(t, err)
		}
		if resp == nil {
			require.NotNil(t, resp, "received nil response")
		}

	})

}
