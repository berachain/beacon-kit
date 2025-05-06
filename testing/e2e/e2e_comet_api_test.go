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
	"sync"

	"github.com/berachain/beacon-kit/testing/e2e/config"
	"github.com/berachain/beacon-kit/testing/e2e/suite/types"
	rpcclient "github.com/cometbft/cometbft/rpc/client"

	sdkcollections "cosmossdk.io/collections"
	"github.com/berachain/beacon-kit/storage/beacondb/keys"
)

// TestABCIInfo compares the ABCI info response among all nodes and cross-checks with the EL.
func (s *BeaconKitE2ESuite) TestABCIInfo() {
	// Wait for execution block 5 to ensure nodes have progressed.
	err := s.WaitForFinalizedBlockNumber(5)
	s.Require().NoError(err)

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
			heightsMap.Store(name, abciInfo.Response.LastBlockHeight)
		}(name, client)
	}

	// Also retrieve height from the EL client.
	elClient := s.JSONRPCBalancer()
	elHeight, err := elClient.BlockNumber(s.Ctx())
	s.Require().NoError(err)
	heightsMap.Store("el", int64(elHeight)) // #nosec G115

	wg.Wait()

	// Check for errors.
	errorsMap.Range(func(key, value interface{}) bool {
		name := key.(string)
		err := value.(error)
		s.Require().NoError(err, "Error getting ABCI info from node %s", name)
		return true
	})

	// Collect heights into a map for comparison.
	heights := make(map[string]int64)
	heightsMap.Range(func(key, value interface{}) bool {
		name := key.(string)
		height := value.(int64)
		heights[name] = height
		s.Logger().Info("Node height", "node", name, "height", height)
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
//
// TODO: verify the proof to ensure full validity.
func (s *BeaconKitE2ESuite) TestABCIQuery() {
	// Wait for execution block 5 to ensure nodes have progressed.
	err := s.WaitForFinalizedBlockNumber(5)
	s.Require().NoError(err)

	// Get all consensus clients.
	clients := s.ConsensusClients()
	s.Require().NotEmpty(clients, "No consensus clients found")

	// Get ABCI query with proof of the fork data from a node.
	abciQuery, err := clients[config.ClientValidator2].ABCIQuery(
		s.Ctx(),
		"store/beacon/key",
		sdkcollections.NewPrefix([]byte{keys.ForkPrefix}),
		rpcclient.ABCIQueryOptions{
			Prove:  true,
			Height: 5,
		},
	)
	s.Require().NoError(err)
	s.Require().NotNil(abciQuery)
	s.Require().Equal(abciQuery.Height, int64(5))
}
