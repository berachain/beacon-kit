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

package config

import (
	cometbft "github.com/berachain/beacon-kit/consensus/cometbft/service"
	"github.com/berachain/beacon-kit/payload/builder"
)

// Consensus clients.
const (
	NumValidators = 5

	ClientValidator0 = "cl-validator-beaconkit-0"
	ClientValidator1 = "cl-validator-beaconkit-1"
	ClientValidator2 = "cl-validator-beaconkit-2"
	ClientValidator3 = "cl-validator-beaconkit-3"
	ClientValidator4 = "cl-validator-beaconkit-4"
)

// DefaultE2ETestConfig provides a default configuration for end-to-end tests,
// pre-populating with a standard set of validators and no additional
// services.
func DefaultE2ETestConfig() *E2ETestConfig {
	return &E2ETestConfig{
		NetworkConfiguration: defaultNetworkConfiguration(),
		NodeSettings:         defaultNodeSettings(),
		AdditionalServices:   defaultAdditionalServices(),
		RPCServiceName:       "el-full-geth-2",
	}
}

// PreconfE2ETestConfig provides a configuration with preconfirmation enabled
// using a dedicated sequencer node.
func PreconfE2ETestConfig() *E2ETestConfig {
	cfg := DefaultE2ETestConfig()
	cfg.NetworkConfiguration.SequencerNode = &Node{
		ElType:  "reth",
		KZGImpl: "crate-crypto/go-kzg-4844",
	}
	cfg.Preconf = PreconfConfig{
		Enabled: true,
	}
	return cfg
}

// PreconfLoadE2ETestConfig provides a configuration for preconf load testing
// with a dedicated sequencer node (matching the devnet YAML topology).
// Uses fewer full/seed nodes than the default to reduce resource contention
// in local Docker environments, which shortens the consensus gap between
// blocks and improves flashblock latency.
func PreconfLoadE2ETestConfig() *E2ETestConfig {
	cfg := DefaultE2ETestConfig()

	// Minimize non-essential nodes to free CPU for consensus. Faster
	// consensus shortens the gap between blocks where no flashblocks
	// are produced, reducing preconf latency.
	cfg.NetworkConfiguration.FullNodes = NodeSet{
		Type: "full",
		Nodes: []Node{{
			ElType:   "reth",
			Replicas: 1,
			KZGImpl:  "crate-crypto/go-kzg-4844",
		}},
	}
	cfg.NetworkConfiguration.SeedNodes = NodeSet{Type: "seed", Nodes: []Node{}}
	cfg.RPCServiceName = "el-full-reth-0"

	cfg.NetworkConfiguration.SequencerNode = &Node{
		ElType:  "reth",
		KZGImpl: "crate-crypto/go-kzg-4844",
	}
	cfg.NetworkConfiguration.PreconfRPCNodes = &NodeSet{
		Type: "preconf-rpc",
		Nodes: []Node{{
			ElType:   "reth",
			Replicas: 1,
			KZGImpl:  "crate-crypto/go-kzg-4844",
		}},
	}
	cfg.Preconf = PreconfConfig{
		Enabled: true,
	}

	// Enable flashblock-monitor to subscribe to the sequencer's WS and
	// output raw flashblock JSON for debugging.
	cfg.AdditionalServices = []AdditionalService{
		{Name: "flashblock-monitor"},
	}

	return cfg
}

func defaultNetworkConfiguration() NetworkConfiguration {
	return NetworkConfiguration{
		Validators: defaultValidators(),
		FullNodes:  defaultFullNodes(),
		SeedNodes:  defaultSeedNodes(),
	}
}

func defaultValidators() NodeSet {
	return NodeSet{
		Type: "validator",
		Nodes: []Node{
			{
				ElType:   "geth",
				Replicas: 3, //nolint:mnd // we want 3 replicas here
				KZGImpl:  "crate-crypto/go-kzg-4844",
			},
			{
				ElType:   "reth",
				Replicas: 2, //nolint:mnd // we want 2 replicas here
				KZGImpl:  "crate-crypto/go-kzg-4844",
			},
		},
	}
}

func defaultFullNodes() NodeSet {
	return NodeSet{
		Type: "full",
		Nodes: []Node{
			{
				ElType:   "reth",
				Replicas: 2, //nolint:mnd // we want 2 replicas here
				KZGImpl:  "crate-crypto/go-kzg-4844",
			},
			{
				ElType:   "geth",
				Replicas: 2, //nolint:mnd // we want 2 replicas here
				KZGImpl:  "crate-crypto/go-kzg-4844",
			},
		},
	}
}

func defaultSeedNodes() NodeSet {
	return NodeSet{
		Type: "seed",
		Nodes: []Node{
			{
				ElType:   "geth",
				Replicas: 1,
				KZGImpl:  "crate-crypto/go-kzg-4844",
			},
		},
	}
}

func defaultNodeSettings() NodeSettings {
	return NodeSettings{
		ExecutionSettings: defaultExecutionSettings(),
		ConsensusSettings: defaultConsensusSettings(),
	}
}

func defaultExecutionSettings() ExecutionSettings {
	return ExecutionSettings{
		Specs: NodeSpecs{
			MinCPU:    0,
			MaxCPU:    0,
			MinMemory: 0,
			MaxMemory: 2048, //nolint:mnd // 2 GB
		},
		Images: map[string]string{
			"geth": "ghcr.io/berachain/bera-geth:latest",
			"reth": "ghcr.io/berachain/bera-reth:af193b66839313cd75b86da0e371baaeb5e814fd",
		},
	}
}

func defaultConsensusSettings() ConsensusSettings {
	var (
		builderCfg = builder.DefaultConfig()
		defaultCfg = cometbft.DefaultConfig()
		consensus  = defaultCfg.Consensus
		p2p        = defaultCfg.P2P
	)

	return ConsensusSettings{
		Specs: NodeSpecs{
			MinCPU:    0,
			MaxCPU:    2000, //nolint:mnd // 2 vCPUs
			MinMemory: 0,
			MaxMemory: 2048, //nolint:mnd // 2 GB
		},
		Images: map[string]string{
			"beaconkit": "beacond:kurtosis-local",
		},
		Config: ConsensusConfig{
			TimeoutPropose:      consensus.TimeoutPropose.String(),
			TimeoutPrevote:      consensus.TimeoutPrevote.String(),
			TimeoutPrecommit:    consensus.TimeoutPrecommit.String(),
			TimeoutCommit:       "0s", // deprecated field, hardcoded to match DefaultConfig()
			MaxNumInboundPeers:  p2p.MaxNumInboundPeers,
			MaxNumOutboundPeers: p2p.MaxNumOutboundPeers,
		},
		AppConfig: AppConfig{
			PayloadTimeout: builderCfg.PayloadTimeout.String(),
		},
	}
}

func defaultAdditionalServices() []AdditionalService {
	return []AdditionalService{}
}
