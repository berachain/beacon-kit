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

package config

// Consensus clients.
const (
	DefaultClient   = "cl-validator-beaconkit-0"
	AlternateClient = "cl-validator-beaconkit-1"
)

// DefaultE2ETestConfig provides a default configuration for end-to-end tests,
// pre-populating with a standard set of validators and no additional
// services.
func DefaultE2ETestConfig() *E2ETestConfig {
	return &E2ETestConfig{
		NetworkConfiguration: defaultNetworkConfiguration(),
		NodeSettings:         defaultNodeSettings(),
		EthJSONRPCEndpoints:  defaultEthJSONRPCEndpoints(),
		AdditionalServices:   defaultAdditionalServices(),
	}
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
				ElType: "nethermind",
				// TODO: restore once we solve
				//  https://github.com/berachain/beacon-kit/issues/2177
				Replicas: 0, // nethermind cannot keep up with deposits checks
				KZGImpl:  "crate-crypto/go-kzg-4844",
			},
			{
				ElType:   "geth",
				Replicas: 2, //nolint:mnd // we want two replicas here
				KZGImpl:  "crate-crypto/go-kzg-4844",
			},
			{
				ElType:   "reth",
				Replicas: 2, //nolint:mnd // we want two replicas here
				KZGImpl:  "crate-crypto/go-kzg-4844",
			},
			{
				ElType:   "erigon",
				Replicas: 1,
				KZGImpl:  "crate-crypto/go-kzg-4844",
			},
			{
				ElType:   "besu",
				Replicas: 0, // Besu causing flakey tests.
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
				ElType:   "nethermind",
				Replicas: 1,
				KZGImpl:  "crate-crypto/go-kzg-4844",
			},
			{
				ElType:   "reth",
				Replicas: 1,
				KZGImpl:  "crate-crypto/go-kzg-4844",
			},
			{
				ElType:   "geth",
				Replicas: 1,
				KZGImpl:  "crate-crypto/go-kzg-4844",
			},
			{
				ElType:   "erigon",
				Replicas: 1,
				KZGImpl:  "crate-crypto/go-kzg-4844",
			},
			{
				ElType:   "besu",
				Replicas: 1,
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
			"besu":       "hyperledger/besu:24.5.4",
			"erigon":     "erigontech/erigon:v2.60.9",
			"ethereumjs": "ethpandaops/ethereumjs:stable",
			"geth":       "ethereum/client-go:stable",
			"nethermind": "nethermind/nethermind:latest",
			"reth":       "ghcr.io/paradigmxyz/reth:latest",
		},
	}
}

func defaultConsensusSettings() ConsensusSettings {
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
			TimeoutPropose:      "3s",
			TimeoutPrevote:      "1s",
			TimeoutPrecommit:    "1s",
			TimeoutCommit:       "3s",
			MaxNumInboundPeers:  40, //nolint:mnd // 40 inbound peers
			MaxNumOutboundPeers: 10, //nolint:mnd // 10 outbound peers
		},
		AppConfig: AppConfig{
			PayloadTimeout:                "1.5s",
			EnableOptimisticPayloadBuilds: false,
		},
	}
}

func defaultEthJSONRPCEndpoints() []EthJSONRPCEndpoint {
	return []EthJSONRPCEndpoint{
		{
			Type: "blutgang",
			Clients: []string{
				// "el-full-nethermind-0",
				// "el-full-reth-0",
				"el-full-geth-2",
				// "el-full-erigon-3",
				// "el-full-erigon-3",
				// Besu causing flakey tests.
				// "el-full-besu-4",
			},
		},
	}
}

func defaultAdditionalServices() []AdditionalService {
	return []AdditionalService{}
}
