package ssz_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/state/deneb"
	. "github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
	"github.com/emicklei/dot"
	"github.com/stretchr/testify/require"
)

func testBeaconState() (*deneb.BeaconState, error) {
	bz, err := os.ReadFile("testdata/beacon.ssz")
	if err != nil {
		return nil, err
	}
	state := &deneb.BeaconState{}
	err = state.UnmarshalSSZ(bz)
	if err != nil {
		return nil, err
	}
	return state, nil
}

func TestTree_Hash(t *testing.T) {
	state, err := testBeaconState()
	require.NoError(t, err)
	rootHash, err := state.HashTreeRoot()
	require.NoError(t, err)
	require.NotNil(t, rootHash)

	tree, err := state.GetRootNode()
	require.NoError(t, err)
	require.NotNil(t, tree)
	require.True(t, bytes.Equal(rootHash[:], tree.CachedHash()))
	require.True(t, bytes.Equal(tree.CachedHash(), tree.Hash()))

	f, err := os.Create("/tmp/beacon.dot")
	require.NoError(t, err)
	defer f.Close()
	DrawTree(tree, f)
}

func DrawTree(n *Node, w io.Writer) {
	n.CachedHash()
	g := dot.NewGraph(dot.Directed)
	drawNode(n, 1, g)
	g.Write(w)
}

func drawNode(n *Node, levelOrder int, g *dot.Graph) dot.Node {
	h := hex.EncodeToString(n.Value)
	dn := g.Node(fmt.Sprintf("n%d", levelOrder)).
		Label(fmt.Sprintf("%d\n%s..%s", levelOrder, h[:3], h[len(h)-3:]))

	if n.Left != nil {
		ln := drawNode(n.Left, 2*levelOrder, g)
		g.Edge(dn, ln).Label("0")
	}
	if n.Right != nil {
		rn := drawNode(n.Right, 2*levelOrder+1, g)
		g.Edge(dn, rn).Label("1")
	}
	return dn
}
