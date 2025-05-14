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

package e2e_test

import (
	"bytes"
	"fmt"
	"sync"

	sdkcollections "cosmossdk.io/collections"
	"github.com/berachain/beacon-kit/storage/beacondb/keys"
	"github.com/berachain/beacon-kit/testing/e2e/config"
	"github.com/berachain/beacon-kit/testing/e2e/suite/types"
	rpcclient "github.com/cometbft/cometbft/rpc/client"
	ics23 "github.com/cosmos/ics23/go"
)

// TestABCIInfo compares the ABCI info response among all nodes and cross-checks with the EL.
func (s *BeaconKitE2ESuite) TestABCIInfo() {
	// Wait for execution block 5 to ensure nodes have progressed.
	s.Require().NoError(s.WaitForFinalizedBlockNumber(5))

	// Get all consensus clients.
	clients := s.ConsensusClients()
	s.Require().NotEmpty(clients, "No consensus clients found")

	// Retrieve heights from all nodes in parallel.
	var (
		wg         sync.WaitGroup
		heightsMap sync.Map
		errorsMap  sync.Map
	)
	for name, client := range clients {
		wg.Add(1)
		go func(name string, client *types.ConsensusClient) {
			defer wg.Done()
			abciInfo, err := client.ABCIInfo(s.Ctx())
			if err != nil {
				errorsMap.Store(name, err)
				return
			}
			heightsMap.Store(name, abciInfo.LastBlockHeight)
		}(name, client)
	}

	// Also retrieve height from the EL client.
	wg.Add(1)
	go func() {
		defer wg.Done()
		elClient := s.JSONRPCBalancer()
		elHeight, err := elClient.BlockNumber(s.Ctx())
		s.Require().NoError(err)
		heightsMap.Store("el", int64(elHeight)) // #nosec G115
	}()

	wg.Wait()

	// Check for errors.
	errorsMap.Range(func(key, value interface{}) bool {
		name := key.(string) //nolint:errcheck // Safe to ignore.
		err := value.(error) //nolint:errcheck // Safe to ignore.
		s.Require().NoError(err, "Error getting ABCI info from node %s", name)
		return true
	})

	// Collect heights into a map for comparison.
	heights := make(map[string]int64)
	heightsMap.Range(func(key, value interface{}) bool {
		name := key.(string)    //nolint:errcheck // Safe to ignore.
		height := value.(int64) //nolint:errcheck // Safe to ignore.
		heights[name] = height
		return true
	})

	// Verify that all heights are within +/- 1 of each other.
	for name1, height1 := range heights {
		for name2, height2 := range heights {
			if name1 == name2 {
				continue
			}

			diff := height1 - height2
			if diff < 0 {
				diff = -diff
			}

			s.Require().LessOrEqual(diff, int64(1),
				"Height difference between nodes %s (%d) and %s (%d) exceeds 1 block",
				name1, height1, name2, height2)
		}
	}
}

// TestABCIQuery checks that the ABCI query response is valid.
func (s *BeaconKitE2ESuite) TestABCIQuery() {
	// Wait for execution block 5 to ensure nodes have progressed.
	s.Require().NoError(s.WaitForFinalizedBlockNumber(5))

	// Get all consensus clients.
	clients := s.ConsensusClients()
	s.Require().NotEmpty(clients, "No consensus clients found")

	// membership

	// Get ABCI query with proof of the fork data from a node.
	key := sdkcollections.NewPrefix([]byte{keys.ForkPrefix})
	abciQuery, err := clients[config.ClientValidator2].ABCIQuery(
		s.Ctx(),
		"store/beacon/key",
		key,
		rpcclient.ABCIQueryOptions{
			Prove:  true,
			Height: 5,
		},
	)
	s.Require().NoError(err)
	s.Require().NotNil(abciQuery)
	s.Require().Equal(abciQuery.Height, int64(5))

	block := int64(6)
	commit, err := clients[config.ClientValidator2].Commit(s.Ctx(), &block)
	if err != nil {
		panic(err)
	}

	proofs := make([]*ics23.CommitmentProof, len(abciQuery.ProofOps.Ops))

	for i := 0; i < len(abciQuery.ProofOps.Ops); i++ {
		proofs[i] = &ics23.CommitmentProof{}
		proofs[i].Unmarshal(abciQuery.ProofOps.Ops[i].Data)
	}

	verifyChainedMembershipProof(
		ics23.CommitmentRoot(commit.SignedHeader.Header.AppHash),
		[]*ics23.ProofSpec{ics23.IavlSpec, ics23.TendermintSpec},
		proofs,
		[][]byte{[]byte("beacon"), key.Bytes()},
		abciQuery.Value,
		0,
	)

	// non-membership

	// Get ABCI query with proof of the fork data from a node.
	key = []byte("oogabooga")
	abciQuery, err = clients[config.ClientValidator2].ABCIQuery(
		s.Ctx(),
		"store/beacon/key",
		key,
		rpcclient.ABCIQueryOptions{
			Prove:  true,
			Height: 5,
		},
	)
	s.Require().NoError(err)
	s.Require().NotNil(abciQuery)
	s.Require().Equal(abciQuery.Height, int64(5))

	block = int64(6)
	commit, err = clients[config.ClientValidator2].Commit(s.Ctx(), &block)
	if err != nil {
		panic(err)
	}

	s.Assert().Empty(abciQuery.Value)

	proofs = make([]*ics23.CommitmentProof, len(abciQuery.ProofOps.Ops))

	for i := 0; i < len(abciQuery.ProofOps.Ops); i++ {
		proofs[i] = &ics23.CommitmentProof{}
		proofs[i].Unmarshal(abciQuery.ProofOps.Ops[i].Data)
	}

	verifyNonMembership(
		proofs,
		[]*ics23.ProofSpec{ics23.IavlSpec, ics23.TendermintSpec},
		commit.SignedHeader.Header.AppHash.Bytes(),
		[][]byte{[]byte("beacon"), key.Bytes()},
	)
}

// https://github.com/cosmos/ibc-go/blob/20326046a09330898fac90540134d8556f4506cc/modules/core/23-commitment/types/merkle.go#L143-L189
func verifyChainedMembershipProof(root []byte, specs []*ics23.ProofSpec, proofs []*ics23.CommitmentProof, keys [][]byte, value []byte, index int) {
	var (
		subroot []byte
		err     error
	)
	// Initialize subroot to value since the proofs list may be empty.
	// This may happen if this call is verifying intermediate proofs after the lowest proof has been executed.
	// In this case, there may be no intermediate proofs to verify and we just check that lowest proof root equals final root
	subroot = value
	for i := index; i < len(proofs); i++ {
		subroot, err = proofs[i].Calculate()
		if err != nil {
			panic(fmt.Sprintf("could not calculate proof root at index %d, merkle tree may be empty. %v", i, err))
		}

		// Since keys are passed in from highest to lowest, we must grab their indices in reverse order
		// from the proofs and specs which are lowest to highest
		key := keys[uint64(len(keys)-1-i)]
		if err != nil {
			panic(fmt.Sprintf("could not retrieve key bytes for key %s: %v", keys[len(keys)-1-i], err))
		}

		ep := proofs[i].GetExist()
		if ep == nil {
			panic(fmt.Sprintf("commitment proof must be existence proof. got: %T at index %d", i, proofs[i]))
		}

		// verify membership of the proof at this index with appropriate key and value
		if err := ep.Verify(specs[i], subroot, key, value); err != nil {
			panic(fmt.Sprintf("failed to verify membership proof at index %d: %v", i, err))
		}
		// Set value to subroot so that we verify next proof in chain commits to this subroot
		value = subroot
	}

	// Check that chained proof root equals passed-in root
	if !bytes.Equal(root, subroot) {
		panic(fmt.Sprintf("proof did not commit to expected root: %X, got: %X. Please ensure proof was submitted with correct proofHeight and to the correct chain.", root, subroot))
	}
}

// https://github.com/cosmos/ibc-go/blob/20326046a09330898fac90540134d8556f4506cc/modules/core/23-commitment/types/merkle.go#L105-L141
func verifyNonMembership(proofs []*ics23.CommitmentProof, specs []*ics23.ProofSpec, root []byte, path [][]byte) {
	// VerifyNonMembership will verify the absence of key in lowest subtree, and then chain inclusion proofs
	// of all subroots up to final root
	subroot, err := proofs[0].Calculate()
	if err != nil {
		panic(fmt.Sprintf("could not calculate root for proof index 0, merkle tree is likely empty. %v", err))
	}

	key := path[uint64(len(path)-1)]

	np := proofs[0].GetNonexist()
	if np == nil {
		panic(fmt.Sprintf("commitment proof must be non-existence proof for verifying non-membership. got: %T", proofs[0]))
	}

	if err := np.Verify(specs[0], subroot, key); err != nil {
		panic(fmt.Sprintf("failed to verify non-membership proof with key %s: %v", string(key), err))
	}

	// Verify chained membership proof starting from index 1 with value = subroot
	verifyChainedMembershipProof(root, specs, proofs, path, subroot, 1)
}
