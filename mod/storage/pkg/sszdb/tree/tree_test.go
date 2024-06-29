package tree_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/state/deneb"
	"github.com/berachain/beacon-kit/mod/storage/pkg/sszdb/tree"
	"github.com/stretchr/testify/require"
)

func testBeaconState() (*deneb.BeaconState, error) {
	bz, err := os.ReadFile("../testdata/beacon.ssz")
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

	tree, err := tree.NewTreeFromFastSSZ(state)
	require.NoError(t, err)
	require.NotNil(t, tree)
	require.True(t, bytes.Equal(rootHash[:], tree.CachedHash()))
	require.True(t, bytes.Equal(tree.CachedHash(), tree.Hash()))

	f, err := os.Create("/tmp/beacon.dot")
	require.NoError(t, err)
	defer f.Close()
	tree.DrawTree(f)
}
