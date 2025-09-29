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

package topology

import (
	"math"
	"math/rand"
	"strings"

	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cometbft/cometbft/p2p"
)

const (
	peersListSeparator = ","
	peerListCutSet     = " "
)

// Assume p2p.PersistentPeers has the **full list of validators**, **ordered by pub key**
func ShapeTestNetwork(p2p *cmtcfg.P2PConfig, thisNodeID p2p.ID) *cmtcfg.P2PConfig {
	return buildTreeTopology(p2p, thisNodeID)
}

func buildTreeTopology(p2p *cmtcfg.P2PConfig, thisNodeID p2p.ID) *cmtcfg.P2PConfig {
	// extract validators
	vals := SplitAndTrimEmpty(p2p.PersistentPeers)
	if len(vals) == 0 {
		return p2p
	}

	// build topology, i.e. map of validator -> list of all node it must connect (and only those)
	topology := buildTree(vals)

	// find out in which position is the validator.
	// If this node is not a validator, connect it to a random leaf
	var (
		persistentPeers      []string
		unconditionalPeerIDs []string

		maxNumInboundPeers, maxNumOutboundPeers = 40, 10 // default values used for non-validators
	)
	if valIdx := retrieveValidatorIdx(thisNodeID, vals); valIdx >= 0 {
		persistentPeers = topology[valIdx]
		for _, p := range persistentPeers {
			parts := strings.SplitN(p, "@", 2) //nolint:mnd // format is ID@IP:PORT
			if len(parts) != 2 {               //nolint:mnd // format is ID@IP:PORT
				panic("don't know how to part this peer to retrieve peerID: " + p)
			}
			unconditionalPeerIDs = append(unconditionalPeerIDs, parts[0])
		}

		// in order to enforce topology strictly, I limit input and outbound node count to
		// just what the topology require. Skip any other validator, seeds, full-nodes.
		// Leaves should have some space to allow seeds/full nodes to connect
		if !isLeaf(topology, valIdx) {
			maxNumInboundPeers = 0
			maxNumOutboundPeers = 0
		}
	} else {
		// node is not a validator. In this topology, connect it to 2 leaves at random
		// leaves nodes are those indexes not connected to root directly
		low, up := len(topology[0])+1, len(vals)-1
		leaf1 := rand.Intn(up-low+1) + low //#nosec: G404 // first index
		leaf2 := rand.Intn(up-low+1) + low //#nosec: G404 // second index, must be different from first one
		for leaf2 == leaf1 {
			leaf2 = rand.Intn(up-low+1) + low //#nosec: G404 // second index, must be different from first one
		}

		persistentPeers = append(persistentPeers, vals[leaf1])
		persistentPeers = append(persistentPeers, vals[leaf2])
	}

	p2p.PersistentPeers = Merge(persistentPeers)
	p2p.UnconditionalPeerIDs = Merge(unconditionalPeerIDs)
	p2p.MaxNumOutboundPeers = maxNumOutboundPeers
	p2p.MaxNumInboundPeers = maxNumInboundPeers
	return p2p
}

func buildTree(vals []string) [][]string {
	// Shape like a tree, with 2 layers
	// layer 0: 1 node, root
	// layer 1: fanOut nodes, so that root has fanOut links
	// layer 2: up to fanOut-1 nodes per each layer 1 node, so fanOut*(fanOut-1), so that layer 1 nodes have fanOut links
	// So total number of nodes is: 1 + fanOut + fanOut*(fanOut-1) = fanOut^2+1
	// We rename N as fanOut to appease then linter
	fanOut := int(math.Ceil(math.Sqrt(float64(len(vals) - 1))))
	topology := make([][]string, len(vals)) // rows are validators, columns are list of peers per validator

	// setup root, just connect it to the first N validators (but itself)
	topology[0] = append(topology[0], vals[1:fanOut+1]...)

	// layer 1 has N nodes, each has N-1 peers, up to len(vals)
	startIdx := fanOut + 1
	for i := range fanOut {
		endIdx := min(startIdx+(fanOut-1), len(vals))
		topology[i+1] = append(topology[i+1], vals[0])
		topology[i+1] = append(topology[i+1], vals[startIdx:endIdx]...)
		startIdx = endIdx
	}

	// layer 2: explicitly enforce peering to the right layer1 node
	for idx := fanOut + 1; idx < len(vals); idx++ {
		parentIdx := (idx-(fanOut+1))/(fanOut-1) + 1 // by Bar-Bera
		topology[idx] = []string{vals[parentIdx]}
	}
	return topology
}

func retrieveValidatorIdx(thisNodeID p2p.ID, vals []string) int {
	idx := -1 // negative int if node is not a validator
	for i, v := range vals {
		if strings.Contains(v, string(thisNodeID)) {
			idx = i
			break
		}
	}
	return idx
}

func isLeaf(topology [][]string, valIdx int) bool {
	if valIdx < 0 {
		return false // not a validator, so not a leaf
	}

	return valIdx > len(topology[0])
}

// adapted from CometBFT, just hard-coding sep and cutset
func SplitAndTrimEmpty(s string) []string {
	if s == "" {
		return []string{}
	}

	spl := strings.Split(s, peersListSeparator)
	nonEmptyStrings := make([]string, 0, len(spl))
	for i := range spl {
		element := strings.Trim(spl[i], peerListCutSet)
		if element != "" {
			nonEmptyStrings = append(nonEmptyStrings, element)
		}
	}
	return nonEmptyStrings
}

// exported to be used by unit tests
func Merge(items []string) string {
	if len(items) == 0 {
		return ""
	}
	var res string
	for _, i := range items {
		res += i + peersListSeparator
	}
	return res[:len(res)-1] // drop final sep
}
