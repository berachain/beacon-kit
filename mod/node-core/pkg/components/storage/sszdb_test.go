// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
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
//
// This integration test exercises the sszdb package with beacon state type
// bindings.  It is placed here to avoid adding consensus-type dependencies to
// the sszdb package.

package storage_test

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/storage/pkg/sszdb"
	"github.com/emicklei/dot"
	fastssz "github.com/ferranbt/fastssz"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
)

type beaconStateMarshallable = types.BeaconState[
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
]

func debugDrawDBTree(
	ctx context.Context,
	t *testing.T,
	db *sszdb.Backend,
	filePath string,
) {
	t.Helper()
	f, err := os.Create(filePath)
	require.NoError(t, err)
	require.NoError(t, db.DrawTree(ctx, f))
}

func debugDrawTree(
	t *testing.T,
	treeable interface {
		sszdb.Treeable
		HashTreeRoot() common.Root
	},
	filePath string,
) {
	t.Helper()
	f, err := os.Create(filePath)
	require.NoError(t, err)
	rootNode, err := sszdb.NewTreeFromFastSSZ(treeable)
	require.NoError(t, err)

	rootNode.CachedHash()
	h := treeable.HashTreeRoot()
	require.Truef(t, bytes.Equal(h[:], rootNode.Value),
		"debugDrawTree: expected %x, got %x", h, rootNode.Value)
	g := dot.NewGraph(dot.Directed)
	drawNode(rootNode, 1, g)
	g.Write(f)

	require.NoError(t, f.Close())
}

func drawNode(n *sszdb.Node, levelOrder int, g *dot.Graph) dot.Node {
	h := hex.EncodeToString(n.Value)
	dn := g.Node(fmt.Sprintf("n%d", levelOrder)).
		Label(fmt.Sprintf("%d\n%s..%s", levelOrder, h[:3], h[len(h)-3:]))
	if n.Value == nil {
		dn = dn.Attr("color", "red")
	}

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

func testBeaconState(t *testing.T) (*beaconStateMarshallable, error) {
	t.Helper()
	bz, err := os.ReadFile("testdata/beacon_state.ssz")
	if err != nil {
		return nil, err
	}
	beacon := &beaconStateMarshallable{}
	err = beacon.UnmarshalSSZ(bz)
	if err != nil {
		return nil, err
	}
	roundtrip, err := beacon.MarshalSSZ()
	require.NoError(t, err)
	require.True(t, bytes.Equal(bz, roundtrip))
	return beacon, nil
}

func TestTree_Basics(t *testing.T) {
	beaconState, err := testBeaconState(t)
	require.NoError(t, err)

	debugDrawTree(t, beaconState, "testdata/beacon_state_start.dot")

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

func newBeaconState() *beaconStateMarshallable {
	return &beaconStateMarshallable{
		GenesisValidatorsRoot: [32]byte{7, 7, 7, 7},
		Slot:                  777,
		Fork: &types.Fork{
			PreviousVersion: [4]byte{1, 2, 3, 4},
			CurrentVersion:  [4]byte{5, 6, 7, 8},
			Epoch:           123,
		},
		LatestBlockHeader: &types.BeaconBlockHeader{
			Slot:            777,
			ProposerIndex:   123,
			ParentBlockRoot: [32]byte{1, 2, 3, 4},
			StateRoot:       [32]byte{5, 6, 7, 8},
			BodyRoot:        [32]byte{9, 10, 11, 12},
		},
		BlockRoots: []common.Root{
			{1, 2, 3, 4},
			{5, 6, 7, 8},
			{9, 10, 11, 12},
			{13, 14, 15, 16},
		},
		StateRoots: []common.Root{},
		LatestExecutionPayloadHeader: &types.ExecutionPayloadHeader{
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
		},
		Eth1Data: &types.Eth1Data{
			DepositRoot:  [32]byte{1, 2, 3, 4},
			DepositCount: 123,
			BlockHash:    [32]byte{5, 6, 7, 8},
		},
		Validators: []*types.Validator{
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
		},
		Balances: []uint64{1000, 2000},
		RandaoMixes: []common.Bytes32{
			{1, 2, 3, 4},
			{5, 6, 7, 8},
		},
		Slashings: []uint64{
			100, 200,
		},
		TotalSlashing: 300,
	}
}

func Test_SchemaDB(t *testing.T) {
	beacon := newBeaconState()
	dir := t.TempDir() + "/sszdb.db"
	t.Logf("db path: %s", dir)
	db, err := sszdb.NewBackend(sszdb.BackendConfig{Path: dir})
	require.NoError(t, err)

	ctx := context.TODO()
	err = db.SaveMonolith(beacon)
	require.NoError(t, err)
	err = db.Commit(ctx)
	require.NoError(t, err)

	debugDrawTree(t, beacon, "testdata/beacon_state_test_start.dot")

	beaconDB, err := sszdb.NewSchemaDB(db, beacon)
	require.NoError(t, err)
	assertRootHash := func() {
		hash := beacon.HashTreeRoot()
		hashSSZ, err := beaconDB.Get(1, 0)
		require.NoError(t, err)
		require.True(t, bytes.Equal(hash[:], hashSSZ))
	}

	genesisValidatorsRoot, err := beaconDB.GetPath(
		ctx,
		"genesis_validators_root",
	)
	require.NoError(t, err)
	require.True(
		t,
		bytes.Equal(beacon.GenesisValidatorsRoot[:], genesisValidatorsRoot),
	)

	slotBz, err := beaconDB.GetPath(ctx, "slot")
	require.NoError(t, err)
	slot := math.U64(fastssz.UnmarshallUint64(slotBz))
	require.Equal(t, beacon.Slot, slot)

	bz, err := beaconDB.GetPath(ctx, "fork")
	require.NoError(t, err)
	beaconBz, err := beacon.Fork.MarshalSSZ()
	require.NoError(t, err)
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

	blockRoots, err := sszdb.DecodeListOfStaticElements(
		ctx,
		beaconDB,
		"block_roots",
		32,
		func(b []byte) (common.Root, error) {
			return common.Root(b), nil
		},
	)
	require.NoError(t, err)
	require.Equal(t, beacon.BlockRoots, blockRoots)

	// validators
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

	debugDrawDBTree(ctx, t, beaconDB.Backend, "testdata/empty.dot")

	err = beaconDB.SetListElementRaw(ctx, "balances", 2, []byte{1})
	require.NoError(t, err)
	beacon.Balances = append(beacon.Balances, 1)
	debugDrawDBTree(ctx, t, beaconDB.Backend, "testdata/one.dot")
	require.NoError(t, beaconDB.Commit(ctx))
	assertRootHash()

	err = beaconDB.SetListElementRaw(ctx, "balances", 3, []byte{2})
	beacon.Balances = append(beacon.Balances, 2)
	require.NoError(t, err)
	debugDrawDBTree(ctx, t, beaconDB.Backend, "testdata/two.dot")
	require.NoError(t, beaconDB.Commit(ctx))
	assertRootHash()

	validator2 := &types.Validator{
		Pubkey:                     [48]byte{10, 11, 12, 13},
		WithdrawalCredentials:      [32]byte{13, 14, 15, 16},
		EffectiveBalance:           789,
		Slashed:                    false,
		ActivationEligibilityEpoch: 789,
		ActivationEpoch:            789,
		ExitEpoch:                  789,
		WithdrawableEpoch:          789,
	}
	err = beaconDB.SetListElementObject(ctx, "validators", 2, validator2)
	require.NoError(t, err)
	beacon.Validators = append(beacon.Validators, validator2)

	err = beaconDB.SetListElementRaw(ctx, "balances", 4, []byte{3})
	require.NoError(t, err)
	beacon.Balances = append(beacon.Balances, 3)

	require.NoError(t, beaconDB.Commit(ctx))
	assertRootHash()

	// execution payload header
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

	// Test Hashes and single node in list retrieval
	hash := beacon.HashTreeRoot()
	hashSSZ, err := beaconDB.Get(1, 0)
	require.NoError(t, err)
	hashSSZ2, err := beaconDB.Hash(ctx)
	require.NoError(t, err)
	require.True(t, bytes.Equal(hash[:], hashSSZ2))
	require.True(t, bytes.Equal(hash[:], hashSSZ))

	beacon.BlockRoots = append(beacon.BlockRoots, common.Root{7, 7, 7, 7})
	hash = beacon.HashTreeRoot()
	require.False(t, bytes.Equal(hash[:], hashSSZ))
	for i, root := range beacon.BlockRoots {
		err = beaconDB.SetListElementRaw(
			ctx,
			"block_roots",
			uint64(i),
			root[:],
		)
		require.NoError(t, err)
	}
	require.NoError(t, db.Commit(ctx))
	hashSSZ, err = beaconDB.Get(1, 0)
	require.NoError(t, err)
	require.True(t, bytes.Equal(hash[:], hashSSZ))

	// now try an append
	beacon.BlockRoots = append(beacon.BlockRoots, common.Root{8, 8, 8, 8})
	hash = beacon.HashTreeRoot()

	require.False(t, bytes.Equal(hash[:], hashSSZ))
	err = beaconDB.SetListElementRaw(
		ctx,
		"block_roots",
		uint64(len(beacon.BlockRoots)-1),
		[]byte{8, 8, 8, 8},
	)
	require.NoError(t, err)
	require.NoError(t, db.Commit(ctx))
	hashSSZ, err = beaconDB.Get(1, 0)
	require.NoError(t, err)
	require.Truef(
		t,
		bytes.Equal(hash[:], hashSSZ),
		"expected %x, got %x",
		hash[:],
		hashSSZ,
	)
}

func Test_Empty_Save(t *testing.T) {
	dir := t.TempDir() + "/sszdb.db"
	db, err := sszdb.NewBackend(sszdb.BackendConfig{Path: dir})
	require.NoError(t, err)
	emptyState := (&beaconStateMarshallable{}).Empty()
	schemaDB, err := sszdb.NewSchemaDB(db, emptyState)
	require.NoError(t, err)
	ctx := context.TODO()

	stateHash := emptyState.HashTreeRoot()
	dbHash, err := schemaDB.Hash(ctx)
	require.NoError(t, err)
	require.True(t, bytes.Equal(stateHash[:], dbHash))
}

func Test_PartialHashing(t *testing.T) {
	dir := t.TempDir() + "/sszdb.db"
	db, err := sszdb.NewBackend(sszdb.BackendConfig{Path: dir})
	require.NoError(t, err)
	emptyState := (&beaconStateMarshallable{}).Empty()
	schemaDB, err := sszdb.NewSchemaDB(db, emptyState)
	require.NoError(t, err)
	ctx := context.TODO()

	stateHash := emptyState.HashTreeRoot()
	dbHash, err := schemaDB.Hash(ctx)
	require.NoError(t, err)
	require.True(t, bytes.Equal(stateHash[:], dbHash))

	assertHashEqual := func() {
		sh := emptyState.HashTreeRoot()
		dh, dbErr := schemaDB.Hash(ctx)
		require.NoError(t, dbErr)
		require.True(t, bytes.Equal(sh[:], dh))
	}

	nextState := newBeaconState()
	err = schemaDB.SetRaw(
		ctx,
		"genesis_validators_root",
		nextState.GenesisValidatorsRoot[:],
	)
	emptyState.GenesisValidatorsRoot = nextState.GenesisValidatorsRoot
	assertHashEqual()
	require.NoError(t, err)

	err = schemaDB.SetRaw(
		ctx,
		"slot",
		fastssz.MarshalUint64(nil, uint64(nextState.Slot)),
	)
	emptyState.Slot = nextState.Slot
	assertHashEqual()
	require.NoError(t, err)

	err = schemaDB.SetObject(ctx, "fork", nextState.Fork)
	emptyState.Fork = nextState.Fork
	assertHashEqual()
	require.NoError(t, err)

	err = schemaDB.SetObject(
		ctx,
		"latest_block_header",
		nextState.LatestBlockHeader,
	)
	emptyState.LatestBlockHeader = nextState.LatestBlockHeader
	assertHashEqual()
	require.NoError(t, err)

	for i, root := range nextState.BlockRoots {
		err = schemaDB.SetListElementRaw(
			ctx,
			"block_roots",
			uint64(i),
			root[:],
		)
		require.NoError(t, err)
	}
	emptyState.BlockRoots = nextState.BlockRoots
	assertHashEqual()

	for i, root := range nextState.StateRoots {
		err = schemaDB.SetListElementRaw(
			ctx,
			"state_roots",
			uint64(i),
			root[:],
		)
		require.NoError(t, err)
	}
	emptyState.StateRoots = nextState.StateRoots
	assertHashEqual()

	err = schemaDB.SetObject(ctx, "eth1_data", nextState.Eth1Data)
	require.NoError(t, err)
	emptyState.Eth1Data = nextState.Eth1Data
	assertHashEqual()

	err = schemaDB.SetObject(
		ctx,
		"latest_execution_payload_header",
		nextState.LatestExecutionPayloadHeader,
	)
	require.NoError(t, err)
	emptyState.LatestExecutionPayloadHeader =
		nextState.LatestExecutionPayloadHeader
	assertHashEqual()

	for i, validator := range nextState.Validators {
		err = schemaDB.SetListElementObject(
			ctx,
			"validators",
			uint64(i),
			validator,
		)
		require.NoError(t, err)
	}
	emptyState.Validators = nextState.Validators
	assertHashEqual()

	for i, balance := range nextState.Balances {
		err = schemaDB.SetListElementRaw(
			ctx,
			"balances",
			uint64(i),
			fastssz.MarshalUint64(nil, balance),
		)
		require.NoError(t, err)
	}
	emptyState.Balances = nextState.Balances
	assertHashEqual()

	for i, mix := range nextState.RandaoMixes {
		err = schemaDB.SetListElementRaw(
			ctx,
			"randao_mixes",
			uint64(i),
			mix[:],
		)
		require.NoError(t, err)
	}
	emptyState.RandaoMixes = nextState.RandaoMixes
	assertHashEqual()

	err = schemaDB.SetRaw(
		ctx,
		"next_withdrawal_index",
		fastssz.MarshalUint64(nil, nextState.NextWithdrawalIndex),
	)
	require.NoError(t, err)

	err = schemaDB.SetRaw(
		ctx,
		"next_withdrawal_validator_index",
		fastssz.MarshalUint64(
			nil,
			nextState.NextWithdrawalValidatorIndex.Unwrap(),
		),
	)
	require.NoError(t, err)

	for i, slashing := range nextState.Slashings {
		err = schemaDB.SetListElementRaw(
			ctx,
			"slashings",
			uint64(i),
			[]byte{byte(slashing)},
		)
		require.NoError(t, err)
	}

	err = schemaDB.SetRaw(
		ctx,
		"total_slashing",
		fastssz.MarshalUint64(nil, nextState.TotalSlashing.Unwrap()),
	)
	require.NoError(t, err)

	dbHash, err = schemaDB.Hash(ctx)
	require.NoError(t, err)
	stateHash = nextState.HashTreeRoot()
	require.True(t, bytes.Equal(stateHash[:], dbHash))
}
