package sszdb_test

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/storage/pkg/sszdb"
	"github.com/emicklei/dot"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
)

func drawTree(n *sszdb.Node, w io.Writer) {
	n.CachedHash()
	g := dot.NewGraph(dot.Directed)
	drawNode(n, 1, g)
	g.Write(w)
}

func drawNode(n *sszdb.Node, levelOrder int, g *dot.Graph) dot.Node {
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

func testBeaconState(t *testing.T) (*types.BeaconState[
	*types.BeaconBlockHeader,
	*types.Eth1Data,
	*types.ExecutionPayloadHeader,
	*types.Fork,
	*types.Validator,
	types.BeaconBlockHeader,
	types.Eth1Data,
	types.ExecutionPayloadHeader,
	types.Fork,
	types.Validator,
], error) {
	t.Helper()
	bz, err := os.ReadFile("testdata/beacon_state.ssz")
	if err != nil {
		return nil, err
	}
	beacon := &types.BeaconState[
		*types.BeaconBlockHeader,
		*types.Eth1Data,
		*types.ExecutionPayloadHeader,
		*types.Fork,
		*types.Validator,
		types.BeaconBlockHeader,
		types.Eth1Data,
		types.ExecutionPayloadHeader,
		types.Fork,
		types.Validator,
	]{}
	err = beacon.UnmarshalSSZ(bz)
	if err != nil {
		return nil, err
	}
	roundtrip, err := beacon.MarshalSSZ()
	require.True(t, bytes.Equal(bz, roundtrip))
	return beacon, nil
}

func TestTree_Basics(t *testing.T) {
	beaconState, err := testBeaconState(t)
	require.NoError(t, err)

	f, err := os.Create("testdata/beacon_state.dot")
	require.NoError(t, err)
	defer f.Close()
	rootNode, err := sszdb.NewTreeFromFastSSZ(beaconState)
	require.NoError(t, err)
	drawTree(rootNode, f)

	dir := t.TempDir() + "/sszdb.db"
	db, err := sszdb.NewBackend(sszdb.BackendConfig{Path: dir})
	require.NoError(t, err)

	ctx := context.TODO()
	err = db.SaveMonolith(beaconState)
	require.NoError(t, err)
	err = db.Commit(ctx)
	require.NoError(t, err)

	schemaDB, err := sszdb.NewSchemaDB(db, beaconState)
	require.NoError(t, err)
	require.NotNil(t, schemaDB)

	payloadHeaderBz, err := schemaDB.GetPath(
		ctx,
		"latest_execution_payload_header",
	)
	require.NoError(t, err)
	bz, err := beaconState.LatestExecutionPayloadHeader.MarshalSSZ()
	require.NoError(t, err)
	require.True(t, bytes.Equal(payloadHeaderBz, bz))
}

func Test_SchemaDB(t *testing.T) {
	beacon := &types.BeaconState[
		*types.BeaconBlockHeader,
		*types.Eth1Data,
		*types.ExecutionPayloadHeader,
		*types.Fork,
		*types.Validator,
		types.BeaconBlockHeader,
		types.Eth1Data,
		types.ExecutionPayloadHeader,
		types.Fork,
		types.Validator,
	]{
		Fork:              &types.Fork{},
		LatestBlockHeader: &types.BeaconBlockHeader{},
		BlockRoots:        []common.Root{},
		StateRoots:        []common.Root{},
		LatestExecutionPayloadHeader: &types.ExecutionPayloadHeader{
			BaseFeePerGas: uint256.MustFromDecimal("123456"),
		},
		Eth1Data:    &types.Eth1Data{},
		Validators:  []*types.Validator{},
		Balances:    []uint64{},
		RandaoMixes: []common.Bytes32{},
		Slashings:   []uint64{},
	}

	beacon.GenesisValidatorsRoot = [32]byte{7, 7, 7, 7}
	beacon.Slot = 777
	beacon.Fork = &types.Fork{
		PreviousVersion: [4]byte{1, 2, 3, 4},
		CurrentVersion:  [4]byte{5, 6, 7, 8},
		Epoch:           123,
	}
	beacon.LatestBlockHeader = &types.BeaconBlockHeader{
		Slot:            777,
		ProposerIndex:   123,
		ParentBlockRoot: [32]byte{1, 2, 3, 4},
		StateRoot:       [32]byte{5, 6, 7, 8},
		BodyRoot:        [32]byte{9, 10, 11, 12},
	}
	beacon.BlockRoots = []common.Root{
		{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}, {13, 14, 15, 16},
	}
	beacon.Validators = []*types.Validator{
		{
			Pubkey:                     [48]byte{1, 2, 3, 4},
			WithdrawalCredentials:      [32]byte{5, 6, 7, 8},
			EffectiveBalance:           123,
			Slashed:                    true,
			ActivationEligibilityEpoch: 123,
			ActivationEpoch:            123,
			ExitEpoch:                  123,
			WithdrawableEpoch:          123,
		},
		{
			Pubkey:                     [48]byte{9, 10, 11, 12},
			WithdrawalCredentials:      [32]byte{13, 14, 15, 16},
			EffectiveBalance:           456,
			Slashed:                    false,
			ActivationEligibilityEpoch: 456,
			ActivationEpoch:            456,
			ExitEpoch:                  456,
			WithdrawableEpoch:          456,
		},
	}
	beacon.LatestExecutionPayloadHeader = &types.ExecutionPayloadHeader{
		StateRoot:     [32]byte{1, 2, 3, 4},
		ReceiptsRoot:  [32]byte{5, 6, 7, 8},
		Random:        [32]byte{13, 14, 15, 16},
		LogsBloom:     [256]byte{17, 18, 19, 20},
		Number:        123,
		GasLimit:      456,
		GasUsed:       789,
		Timestamp:     101112,
		ExtraData:     []byte{29, 30, 31, 32, 35},
		BaseFeePerGas: uint256.MustFromDecimal("123456"),
	}

	dir := t.TempDir() + "/sszdb.db"
	t.Logf("db path: %s", dir)
	db, err := sszdb.NewBackend(sszdb.BackendConfig{Path: dir})
	require.NoError(t, err)

	ctx := context.TODO()
	err = db.SaveMonolith(beacon)
	require.NoError(t, err)
	err = db.Commit(ctx)
	require.NoError(t, err)

	f, err := os.Create("testdata/beacon_state.dot")
	require.NoError(t, err)
	defer f.Close()
	rootNode, err := sszdb.NewTreeFromFastSSZ(beacon)
	require.NoError(t, err)
	drawTree(rootNode, f)

	beaconDB, err := sszdb.NewSchemaDB(db, beacon)
	require.NoError(t, err)

	genesisValidatorsRoot, err := beaconDB.GetPath(
		ctx,
		"genesis_validators_root",
	)
	require.NoError(t, err)
	require.True(
		t,
		bytes.Equal(beacon.GenesisValidatorsRoot[:], genesisValidatorsRoot),
	)

	slot, err := beaconDB.GetSlot(ctx)
	require.NoError(t, err)
	require.Equal(t, beacon.Slot, slot)

	bz, err := beaconDB.GetPath(ctx, "fork")
	require.NoError(t, err)
	beaconBz, err := beacon.Fork.MarshalSSZ()
	fork := &types.Fork{}
	require.NoError(t, fork.UnmarshalSSZ(bz))
	require.True(t, bytes.Equal(bz, beaconBz))
	require.Equal(t, beacon.Fork, fork)

	blockHeaderBz, err := beaconDB.GetPath(ctx, "latest_block_header")
	require.NoError(t, err)
	blockHeader := &types.BeaconBlockHeader{}
	err = blockHeader.UnmarshalSSZ(blockHeaderBz)
	require.NoError(t, err)
	require.Equal(t, beacon.LatestBlockHeader, blockHeader)

	blockRoots, err := beaconDB.GetBlockRoots(ctx)
	require.NoError(t, err)
	require.Equal(t, beacon.BlockRoots, blockRoots)

	validator0Bz, err := beaconDB.GetPath(ctx, "validators/0")
	require.NoError(t, err)
	validator0 := &types.Validator{}
	err = validator0.UnmarshalSSZ(validator0Bz)
	require.NoError(t, err)
	require.Equal(t, beacon.Validators[0], validator0)

	validator1Bz, err := beaconDB.GetPath(ctx, "validators/1")
	require.NoError(t, err)
	validator1 := &types.Validator{}
	err = validator1.UnmarshalSSZ(validator1Bz)
	require.NoError(t, err)
	require.Equal(t, beacon.Validators[1], validator1)

	executionPayloadHeader := &types.ExecutionPayloadHeader{}
	executionPayloadHeaderBz, err := beaconDB.GetPath(
		ctx,
		"latest_execution_payload_header",
	)
	require.NoError(t, err)
	err = executionPayloadHeader.UnmarshalSSZ(executionPayloadHeaderBz)
	require.NoError(t, err)
	require.Equal(
		t,
		beacon.LatestExecutionPayloadHeader,
		executionPayloadHeader,
	)
}
