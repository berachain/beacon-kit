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
		if node.Mode == e2e.ModeSeed {
			return
		}
		fmt.Println("Testing node-api genesis for node", node.Name)
		ip := getNodeIP(node)

		nIP := fmt.Sprintf("http://%v:3500", ip)
		client, err := node.Client()
		require.NoError(t, err)

		fmt.Println("Connecting to beacon node at", nIP)
		bclient := &BeaconHTTPClient{
			Client: &http.Client{
				Timeout: time.Second * 10,
			},
			baseURL: nIP,
		}

		resp, err := bclient.Get(bclient.baseURL + "/eth/v1/beacon/genesis") //fmt.Sprintf("/eth/v1/beacon/states/%s/validators/%s", "10", node.PrivvalKey.PubKey().Address().String()))
		if err != nil {
			t.Log(err)
			t.Fatalf("failed to get validator info")
		}
		if resp == nil {
			t.Fatalf("received nil response")
		}

		netInfo, err := client.NetInfo(ctx)
		require.NoError(t, err)

		seen := map[string]bool{}
		for _, n := range node.Testnet.Nodes {
			seen[n.Name] = n.Name == node.Name // we've clearly seen ourself
		}
		for _, peerInfo := range netInfo.Peers {
			peer := node.Testnet.LookupNode(peerInfo.NodeInfo.Moniker)
			require.NotNil(t, peer, "unknown node %v", peerInfo.NodeInfo.Moniker)
			require.Equal(t, peer.InternalIP.String(), peerInfo.RemoteIP,
				"unexpected IP address for peer %v", peer.Name)
			seen[peerInfo.NodeInfo.Moniker] = true
		}

		for name := range seen {
			require.True(t, seen[name], "node %v not peered with %v", node.Name, name)
		}
	})
}
