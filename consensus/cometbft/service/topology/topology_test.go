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
//

package topology_test

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	cometbft "github.com/berachain/beacon-kit/consensus/cometbft/service"
	"github.com/berachain/beacon-kit/consensus/cometbft/service/topology"
	"github.com/cometbft/cometbft/p2p"
	"github.com/stretchr/testify/require"
)

// Throw-away test to checks all invariants we require to the tree topology
func TestTreeTopologyOver36NodesNetwork(t *testing.T) {
	t.Parallel()

	valCount := 36
	vals, nonVal := generateNodes(t, valCount)
	expectedTopology := map[string][]string{
		vals[0]: []string{vals[1], vals[2], vals[3], vals[4], vals[5], vals[6]},

		vals[1]: []string{vals[0], vals[7], vals[8], vals[9], vals[10], vals[11]},
		vals[2]: []string{vals[0], vals[12], vals[13], vals[14], vals[15], vals[16]},
		vals[3]: []string{vals[0], vals[17], vals[18], vals[19], vals[20], vals[21]},
		vals[4]: []string{vals[0], vals[22], vals[23], vals[24], vals[25], vals[26]},
		vals[5]: []string{vals[0], vals[27], vals[28], vals[29], vals[30], vals[31]},
		vals[6]: []string{vals[0], vals[32], vals[33], vals[34], vals[35]},

		vals[7]:  []string{vals[1]},
		vals[8]:  []string{vals[1]},
		vals[9]:  []string{vals[1]},
		vals[10]: []string{vals[1]},
		vals[11]: []string{vals[1]},
		vals[12]: []string{vals[2]},
		vals[13]: []string{vals[2]},
		vals[14]: []string{vals[2]},
		vals[15]: []string{vals[2]},
		vals[16]: []string{vals[2]},
		vals[17]: []string{vals[3]},
		vals[18]: []string{vals[3]},
		vals[19]: []string{vals[3]},
		vals[20]: []string{vals[3]},
		vals[21]: []string{vals[3]},
		vals[22]: []string{vals[4]},
		vals[23]: []string{vals[4]},
		vals[24]: []string{vals[4]},
		vals[25]: []string{vals[4]},
		vals[26]: []string{vals[4]},
		vals[27]: []string{vals[5]},
		vals[28]: []string{vals[5]},
		vals[29]: []string{vals[5]},
		vals[30]: []string{vals[5]},
		vals[31]: []string{vals[5]},
		vals[32]: []string{vals[6]},
		vals[33]: []string{vals[6]},
		vals[34]: []string{vals[6]},
		vals[35]: []string{vals[6]},
	}

	// check validators connections
	for i, val := range vals {
		p2pCfgIn := cometbft.DefaultConfig()
		p2pCfgIn.P2P.PersistentPeers = topology.Merge(vals)
		valID := extractNodeID(t, val)

		// test
		p2pCfgOut := topology.ShapeTestNetwork(p2pCfgIn.P2P, valID)

		// checks
		expectedPeers := expectedTopology[val]
		gotPeers := topology.SplitAndTrimEmpty(p2pCfgOut.PersistentPeers)
		require.Equal(t, expectedPeers, gotPeers, "mismatch, validator %d", i)

		if i <= 6 {
			// non-leaf can only connect to other validators
			require.Zero(t, p2pCfgOut.MaxNumInboundPeers)
			require.Zero(t, p2pCfgOut.MaxNumOutboundPeers)
		} else {
			// leaves should have default parameters to allow seeds/full-nodes to connect and send traffic
			require.Equal(t, 40, p2pCfgOut.MaxNumInboundPeers)
			require.Equal(t, 10, p2pCfgOut.MaxNumOutboundPeers)
		}
	}

	// check non-validators connections
	{
		p2pCfgIn := cometbft.DefaultConfig()
		p2pCfgIn.P2P.PersistentPeers = topology.Merge(vals)
		valID := extractNodeID(t, nonVal)

		// test
		p2pCfgOut := topology.ShapeTestNetwork(p2pCfgIn.P2P, valID)

		// checks: non-vals should be connected to leaves only
		re := regexp.MustCompile(`devnet-val-(\d+)@`)
		gotPeers := topology.SplitAndTrimEmpty(p2pCfgOut.PersistentPeers)
		for i, p := range gotPeers {
			matches := re.FindStringSubmatch(p)
			require.GreaterOrEqual(t, len(matches), 2, "mismatch, peer %d", i)
			gotIdx, err := strconv.Atoi(matches[1])
			require.NoError(t, err)

			require.Greater(t, gotIdx, 6, "mismatch, peer %d", i)
		}
	}
}

func TestTreeTopologyOver5NodesNetwork(t *testing.T) {
	t.Parallel()

	valCount := 5
	vals, nonVal := generateNodes(t, valCount)
	expectedTopology := map[string][]string{
		vals[0]: []string{vals[1], vals[2]},

		vals[1]: []string{vals[0], vals[3]},
		vals[2]: []string{vals[0], vals[4]},

		vals[3]: []string{vals[1]},
		vals[4]: []string{vals[2]},
	}

	for i, val := range vals {
		p2pCfgIn := cometbft.DefaultConfig()
		p2pCfgIn.P2P.PersistentPeers = topology.Merge(vals)
		valID := extractNodeID(t, val)

		// test
		p2pCfgOut := topology.ShapeTestNetwork(p2pCfgIn.P2P, valID)

		// checks
		expectedPeers := expectedTopology[val]
		gotPeers := topology.SplitAndTrimEmpty(p2pCfgOut.PersistentPeers)
		require.Equal(t, expectedPeers, gotPeers, "mismatch, validator %d", i)

		if i <= 2 {
			require.Zero(t, p2pCfgOut.MaxNumInboundPeers)
			require.Zero(t, p2pCfgOut.MaxNumOutboundPeers)
		} else {
			// leaves should have default parameters to allow seeds/full-nodes to connect and send traffic
			require.Equal(t, 40, p2pCfgOut.MaxNumInboundPeers)
			require.Equal(t, 10, p2pCfgOut.MaxNumOutboundPeers)
		}
	}

	// check non-validators connections
	{
		p2pCfgIn := cometbft.DefaultConfig()
		p2pCfgIn.P2P.PersistentPeers = topology.Merge(vals)
		valID := extractNodeID(t, nonVal)

		// test
		p2pCfgOut := topology.ShapeTestNetwork(p2pCfgIn.P2P, valID)

		// checks: non-vals should be connected to leaves only
		re := regexp.MustCompile(`devnet-val-(\d+)@`)
		gotPeers := topology.SplitAndTrimEmpty(p2pCfgOut.PersistentPeers)
		for i, p := range gotPeers {
			matches := re.FindStringSubmatch(p)
			require.GreaterOrEqual(t, len(matches), 2, "mismatch, peer %d", i)
			gotIdx, err := strconv.Atoi(matches[1])
			require.NoError(t, err)

			require.Greater(t, gotIdx, 2, "mismatch, peer %d", i)
		}
	}
}

func generateNodes(t *testing.T, valCount int) ([]string, string) {
	t.Helper()

	length := p2p.IDByteLength  // number of hex characters
	byteLen := (length + 1) / 2 // convert to bytes (round up)

	// non validator
	nonValIdx := valCount + 1
	b := make([]byte, byteLen)
	_, err := rand.Read(b)
	require.NoError(t, err, "index %d", nonValIdx)
	nodeID := hex.EncodeToString(b)[:length]

	// append NODEID@IP:PORT at the result
	nonVal := nodeID + fmt.Sprintf("devnet-val-%d@1.2.3.%d:26656", nonValIdx, nonValIdx+10)

	res := make([]string, 0)
	if valCount == 0 {
		return res, nonVal
	}

	for i := range valCount {
		// generate random nodeID
		b = make([]byte, byteLen)
		_, err = rand.Read(b)
		require.NoError(t, err, "index %d", i)
		nodeID = hex.EncodeToString(b)[:length]

		// append NODEID@IP:PORT at the result
		res = append(res, nodeID+fmt.Sprintf("devnet-val-%d@1.2.3.%d:26656", i, i+10))
	}

	return res, nonVal
}

func extractNodeID(t *testing.T, peer string) p2p.ID {
	t.Helper()

	parts := strings.SplitN(peer, "@", 2)
	require.Len(t, parts, 2)
	return p2p.ID(parts[0])
}
