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

//nolint:testpackage // private functions.
package merkle

import (
	"bytes"
	"encoding/hex"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/crypto/sha256"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/merkle"
	"github.com/berachain/beacon-kit/primitives/merkle/zero"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

const (
	url      = "https://raw.githubusercontent.com/ethereum/EIPs/master/assets/eip-4881/test_cases.yaml"
	filePath = "testdata/eip4881_test_cases.yaml"
)

// EIP 4881 spec test cases from:
// https://github.com/ethereum/EIPs/blob/master/assets/eip-4881/test_cases.yaml.
//
// NOTE: these test cases must be downloaded from Github. If no internet access
// is available, the test cases can be manually downloaded and added to the
// `testdata` folder with filename `eip4881_test_cases.yaml`.
func readTestCases(t *testing.T) []testCase {
	t.Helper()

	var testCases []testCase

	resp, err := http.Get(url)
	if err != nil {
		t.Logf("failed to download test cases (%v), trying local file...", err)

		var enc []byte
		enc, err = os.ReadFile(filePath)
		if err != nil {
			t.Skipf("skipping, local file %s is missing (%v)", filePath, err)
			return nil
		}

		err = yaml.Unmarshal(enc, &testCases)
		require.NoError(t, err)
		require.NotEmpty(t, testCases)
		return testCases
	}
	defer resp.Body.Close()

	err = yaml.NewDecoder(resp.Body).Decode(&testCases)
	require.NoError(t, err)
	require.NotEmpty(t, testCases)
	return testCases
}

type testCase struct {
	DepositData     depositData `yaml:"deposit_data"`
	DepositDataRoot common.Root `yaml:"deposit_data_root"`
	Eth1Data        *eth1Data   `yaml:"eth1_data"`
	BlockHeight     uint64      `yaml:"block_height"`
	Snapshot        snapshot    `yaml:"snapshot"`
}

func (tc *testCase) UnmarshalYAML(value *yaml.Node) error {
	raw := struct {
		DepositData     depositData `yaml:"deposit_data"`
		DepositDataRoot string      `yaml:"deposit_data_root"`
		Eth1Data        *eth1Data   `yaml:"eth1_data"`
		BlockHeight     string      `yaml:"block_height"`
		Snapshot        snapshot    `yaml:"snapshot"`
	}{}
	err := value.Decode(&raw)
	if err != nil {
		return err
	}
	tc.DepositDataRoot, err = hexStringToByteArray(raw.DepositDataRoot)
	if err != nil {
		return err
	}
	tc.DepositData = raw.DepositData
	tc.Eth1Data = raw.Eth1Data
	tc.BlockHeight, err = stringToUint64(raw.BlockHeight)
	if err != nil {
		return err
	}
	tc.Snapshot = raw.Snapshot
	return nil
}

type depositData struct {
	Pubkey                []byte `yaml:"pubkey"`
	WithdrawalCredentials []byte `yaml:"withdrawal_credentials"`
	Amount                uint64 `yaml:"amount"`
	Signature             []byte `yaml:"signature"`
}

func (dd *depositData) UnmarshalYAML(value *yaml.Node) error {
	raw := struct {
		Pubkey                string `yaml:"pubkey"`
		WithdrawalCredentials string `yaml:"withdrawal_credentials"`
		Amount                string `yaml:"amount"`
		Signature             string `yaml:"signature"`
	}{}
	err := value.Decode(&raw)
	if err != nil {
		return err
	}
	dd.Pubkey, err = hexStringToBytes(raw.Pubkey)
	if err != nil {
		return err
	}
	dd.WithdrawalCredentials, err = hexStringToBytes(raw.WithdrawalCredentials)
	if err != nil {
		return err
	}
	dd.Amount, err = strconv.ParseUint(raw.Amount, 10, 64)
	if err != nil {
		return err
	}
	dd.Signature, err = hexStringToBytes(raw.Signature)
	if err != nil {
		return err
	}
	return nil
}

type eth1Data struct {
	DepositRoot  common.Root          `yaml:"deposit_root"`
	DepositCount uint64               `yaml:"deposit_count"`
	BlockHash    common.ExecutionHash `yaml:"block_hash"`
}

func (ed *eth1Data) UnmarshalYAML(value *yaml.Node) error {
	raw := struct {
		DepositRoot  string `yaml:"deposit_root"`
		DepositCount string `yaml:"deposit_count"`
		BlockHash    string `yaml:"block_hash"`
	}{}
	err := value.Decode(&raw)
	if err != nil {
		return err
	}
	ed.DepositRoot, err = hexStringToByteArray(raw.DepositRoot)
	if err != nil {
		return err
	}
	ed.DepositCount, err = stringToUint64(raw.DepositCount)
	if err != nil {
		return err
	}
	var blockHash common.Root
	blockHash, err = hexStringToByteArray(raw.BlockHash)
	if err != nil {
		return err
	}
	ed.BlockHash = common.ExecutionHash(blockHash)
	return nil
}

type snapshot struct {
	DepositTreeSnapshot
}

func (sd *snapshot) UnmarshalYAML(value *yaml.Node) error {
	raw := struct {
		Finalized            []string `yaml:"finalized"`
		DepositRoot          string   `yaml:"deposit_root"`
		DepositCount         string   `yaml:"deposit_count"`
		ExecutionBlockHash   string   `yaml:"execution_block_hash"`
		ExecutionBlockHeight string   `yaml:"execution_block_height"`
	}{}
	err := value.Decode(&raw)
	if err != nil {
		return err
	}
	sd.finalized = make([]common.Root, len(raw.Finalized))
	for i, finalized := range raw.Finalized {
		sd.finalized[i], err = hexStringToByteArray(finalized)
		if err != nil {
			return err
		}
	}
	sd.depositRoot, err = hexStringToByteArray(raw.DepositRoot)
	if err != nil {
		return err
	}
	sd.depositCount, err = stringToUint64(raw.DepositCount)
	if err != nil {
		return err
	}
	var executionHash common.Root
	executionHash, err = hexStringToByteArray(raw.ExecutionBlockHash)
	if err != nil {
		return err
	}
	sd.executionBlock.Hash = common.ExecutionHash(executionHash)
	var depth uint64
	depth, err = stringToUint64(raw.ExecutionBlockHeight)
	if err != nil {
		return err
	}
	sd.executionBlock.Depth = math.U64(depth)
	sd.hasher = merkle.NewHasher[common.Root](sha256.Hash)
	return nil
}

func hexStringToByteArray(s string) (common.Root, error) {
	raw, err := hexStringToBytes(s)
	if err != nil {
		return common.Root{}, err
	}
	if len(raw) != 32 {
		return common.Root{}, errors.New("invalid hex string length")
	}
	b := common.Root{}
	copy(b[:], raw[:32])
	return b, nil
}

func hexStringToBytes(s string) ([]byte, error) {
	return hex.DecodeString(strings.TrimPrefix(s, "0x"))
}

func stringToUint64(s string) (uint64, error) {
	value, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func merkleRootFromBranch(
	hasher merkle.Hasher[common.Root],
	leaf common.Root,
	branch [constants.DepositContractDepth + 1]common.Root,
	index uint64,
) common.Root {
	root := leaf
	for i, l := range branch {
		ithBit := (index >> i) & 0x1
		if ithBit == 1 {
			root = hasher.Combi(l, root)
		} else {
			root = hasher.Combi(root, l)
		}
	}
	return root
}

func checkProof(t *testing.T, tree *DepositTree, index uint64) {
	t.Helper()
	leaf, proof, err := tree.getProof(index)
	require.NoError(t, err)
	calcRoot := merkleRootFromBranch(tree.hasher, leaf, proof, index)
	require.Equal(t, tree.getRoot(), calcRoot)
}

func compareProof(t *testing.T, tree1, tree2 *DepositTree, index uint64) {
	t.Helper()
	require.Equal(t, tree1.getRoot(), tree2.getRoot())
	checkProof(t, tree1, index)
	checkProof(t, tree2, index)
}

func cloneFromSnapshot(
	t *testing.T,
	snapshot DepositTreeSnapshot,
	testCases []testCase,
) *DepositTree {
	t.Helper()
	cp, err := fromSnapshot(snapshot)
	require.NoError(t, err)
	for _, c := range testCases {
		err = cp.pushLeaf(c.DepositDataRoot)
		require.NoError(t, err)
	}
	return cp
}

func TestDepositCases(t *testing.T) {
	tree := NewDepositTree()
	testCases := readTestCases(t)
	var err error
	for _, c := range testCases {
		err = tree.pushLeaf(c.DepositDataRoot)
		require.NoError(t, err)
	}
}

type Test struct {
	DepositDataRoot common.Root
}

func TestRootEquivalence(t *testing.T) {
	var err error
	tree := NewDepositTree()
	testCases := readTestCases(t)

	depositRoots := make([]common.Root, len(testCases[:128]))
	for i, c := range testCases[:128] {
		err = tree.pushLeaf(c.DepositDataRoot)
		require.NoError(t, err)
		depositRoots[i] = c.DepositDataRoot
	}
	originalRoot := tree.HashTreeRoot()

	generatedTree, err := merkle.NewTreeFromLeavesWithDepth(
		depositRoots,
		uint8(constants.DepositContractDepth),
	)
	require.NoError(t, err)

	rootA := generatedTree.HashTreeRoot()
	require.True(t, rootA.Equals(common.NewRootFromBytes(originalRoot[:])))
}

func TestFinalization(t *testing.T) {
	tree := NewDepositTree()
	testCases := readTestCases(t)
	var err error
	for _, c := range testCases[:128] {
		err = tree.pushLeaf(c.DepositDataRoot)
		require.NoError(t, err)
	}
	originalRoot := tree.getRoot()
	require.True(
		t,
		bytes.Equal(testCases[127].Eth1Data.DepositRoot[:], originalRoot[:]),
	)
	err = tree.Finalize(
		testCases[100].Eth1Data.DepositCount-1,
		testCases[100].Eth1Data.BlockHash,
		math.U64(testCases[100].BlockHeight),
	)
	require.NoError(t, err)
	// ensure finalization doesn't change root
	require.Equal(t, tree.getRoot(), originalRoot)
	snapshotData := tree.GetSnapshot()
	require.True(
		t,
		testCases[100].Snapshot.DepositTreeSnapshot.Equals(&snapshotData),
	)
	// create a copy of the tree from a snapshot by replaying
	// the deposits after the finalized deposit
	cp := cloneFromSnapshot(t, snapshotData, testCases[101:128])
	// ensure original and copy have the same root
	require.Equal(t, tree.getRoot(), cp.getRoot())
	//	finalize original again to check double finalization
	err = tree.Finalize(
		testCases[105].Eth1Data.DepositCount-1,
		testCases[105].Eth1Data.BlockHash,
		math.U64(testCases[105].BlockHeight),
	)
	require.NoError(t, err)
	//	root should still be the same
	require.Equal(t, originalRoot, tree.getRoot())
	// create a copy of the tree by taking a snapshot again
	snapshotData = tree.GetSnapshot()
	cp = cloneFromSnapshot(t, snapshotData, testCases[106:128])
	// create a copy of the tree by replaying ALL deposits from nothing
	fullTreeCopy := NewDepositTree()
	for _, c := range testCases[:128] {
		err = fullTreeCopy.pushLeaf(c.DepositDataRoot)
		require.NoError(t, err)
	}
	for i := 106; i < 128; i++ {
		compareProof(t, tree, cp, uint64(i))
		compareProof(t, tree, fullTreeCopy, uint64(i))
	}
}

func TestSnapshotCases(t *testing.T) {
	tree := NewDepositTree()
	testCases := readTestCases(t)
	var err error
	for _, c := range testCases {
		err = tree.pushLeaf(c.DepositDataRoot)
		require.NoError(t, err)
	}
	for _, c := range testCases {
		err = tree.Finalize(
			c.Eth1Data.DepositCount-1,
			c.Eth1Data.BlockHash,
			math.U64(c.BlockHeight),
		)
		require.NoError(t, err)
		s := tree.GetSnapshot()
		require.True(t, c.Snapshot.DepositTreeSnapshot.Equals(&s))
	}
}

func TestInvalidSnapshot(t *testing.T) {
	invalidSnapshot := DepositTreeSnapshot{
		finalized:    nil,
		depositRoot:  zero.Hashes[0],
		depositCount: 0,
		executionBlock: executionBlock{
			Hash:  zero.Hashes[0],
			Depth: 0,
		},
		hasher: merkle.NewHasher[common.Root](sha256.Hash),
	}
	_, err := fromSnapshot(invalidSnapshot)
	require.ErrorContains(t, err, "snapshot root is invalid")
}

func TestEmptyTree(t *testing.T) {
	tree := NewDepositTree()
	require.Equal(
		t,
		"0xd70a234731285c6804c2a4f56711ddb8c82c99740f207854891028af34e27e5e",
		tree.getRoot().Hex(),
	)
}
