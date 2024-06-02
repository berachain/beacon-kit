// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

//nolint:tagliatelle // starlark uses snek case.
package config

import (
	"encoding/json"
)

// E2ETestConfig defines the configuration for end-to-end tests, including any
// additional services and validators involved.
type E2ETestConfig struct {
	// AdditionalServices specifies any extra services that should be included
	// in the test environment.
	AdditionalServices []AdditionalService `json:"additional_services"`
	// NetworkConfiguration specifies the configuration for the network.
	NetworkConfiguration NetworkConfiguration `json:"network_configuration"`
	// NodeSettings specifies the configuration for the nodes in the test.
	NodeSettings NodeSettings `json:"node_settings"`
	// EthJSONRPCEndpoints specifies the RPC endpoints to include in the test.
	EthJSONRPCEndpoints []EthJSONRPCEndpoint `json:"eth_json_rpc_endpoints"`
}

type NetworkConfiguration struct {
	// Validators lists the configurations for each validator in the test.
	Validators NodeSet `json:"validators"`
	// FullNodes specifies the number of full nodes to include in the test.
	FullNodes NodeSet `json:"full_nodes"`
	// SeedNodes specifies the number of seed nodes to include in the test.
	SeedNodes NodeSet `json:"seed_nodes"`
}

type EthJSONRPCEndpoint struct {
	Type    string   `json:"type"`
	Clients []string `json:"clients"`
}

// NodeSet holds nodes that have a distinct role in the network.
type NodeSet struct {
	// Type is the type of node set.
	Type string `json:"type"`
	// Nodes is a list of nodes in the set.
	Nodes []Node `json:"nodes"`
}

// Node holds the configuration for a single node in the test,
// including client images and types.
type Node struct {
	// ElType denotes the type of execution layer client (e.g., reth).
	ElType string `json:"el_type"`
	// Replicas specifies the number of replicas to use for the client.
	Replicas int `json:"replicas"`
}

// NodeSettings holds the configuration for a single node in the test,
// including client images and types.
type NodeSettings struct {
	// ExecutionSettings holds the configuration for the execution layer
	// clients.
	ExecutionSettings ExecutionSettings `json:"execution_settings"`
	// ConsensusSettings holds the configuration for the consensus layer
	// clients.
	ConsensusSettings ConsensusSettings `json:"consensus_settings"`
}

// ExecutionSettings holds the configuration for the execution layer clients.
type ExecutionSettings struct {
	// Images specifies the images to use for the execution layer clients.
	Images map[string]string `json:"images"`
}

// ConsensusSettings holds the configuration for the consensus layer clients.
type ConsensusSettings struct {
	// Images specifies the images to use for the consensus layer clients.
	Images map[string]string `json:"images"`
}

// AdditionalService holds the configuration for an additional service
// to be included in the test.
type AdditionalService struct {
	// Name specifies the name of the additional service.
	Name string `json:"name"`
	// Replicas specifies the number of replicas to use for the service.
	Replicas int `json:"replicas"`
}

// DefaultE2ETestConfig provides a default configuration for end-to-end tests,
// pre-populating with a standard set of validators and no additional
// services.
func DefaultE2ETestConfig() *E2ETestConfig {
	return &E2ETestConfig{
		AdditionalServices: []AdditionalService{
			{
				Name:     "tx-fuzz",
				Replicas: 1,
			},
		},
		NetworkConfiguration: NetworkConfiguration{
			Validators: NodeSet{
				Type: "validator",
				Nodes: []Node{
					{
						ElType:   "nethermind",
						Replicas: 0,
					},
					{
						ElType:   "geth",
						Replicas: 1,
					},
					{
						ElType:   "reth",
						Replicas: 2, //nolint:mnd // 2 replicas
					},
					{
						ElType:   "erigon",
						Replicas: 1,
					},
					{
						ElType:   "besu",
						Replicas: 0,
					},
				},
			},
			FullNodes: NodeSet{
				Type: "full",
				Nodes: []Node{
					{
						ElType:   "nethermind",
						Replicas: 0,
					},
					{
						ElType:   "reth",
						Replicas: 2, //nolint:mnd // 2 replicas
					},
					{
						ElType:   "geth",
						Replicas: 1,
					},
					{
						ElType:   "erigon",
						Replicas: 1,
					},
					{
						ElType:   "besu",
						Replicas: 0,
					},
				},
			},
			SeedNodes: NodeSet{
				Type: "seed",
				Nodes: []Node{
					{
						ElType:   "reth",
						Replicas: 1,
					},
				},
			},
		},
		NodeSettings: NodeSettings{
			ExecutionSettings: ExecutionSettings{
				Images: map[string]string{
					"besu":       "hyperledger/besu:latest",
					"erigon":     "thorax/erigon:latest",
					"ethereumjs": "ethpandaops/ethereumjs:stable",
					"geth":       "ethereum/client-go:latest",
					"nethermind": "nethermind/nethermind:latest",
					"reth":       "ghcr.io/paradigmxyz/reth:latest",
				},
			},
			ConsensusSettings: ConsensusSettings{
				Images: map[string]string{
					"beaconkit": "beacond:kurtosis-local",
				},
			},
		},
		EthJSONRPCEndpoints: []EthJSONRPCEndpoint{
			{
				Type: "blutgang",
				Clients: []string{
					// "el-full-nethermind-0",
					"el-full-reth-0",
					"el-full-reth-1",
					"el-full-geth-2",
					// "el-full-erigon-3",
					// Besu causing flakey tests.
					// "el-full-besu-4",
				},
			},
		},
	}
}

// MustMarshalJSON marshals the E2ETestConfig to JSON, panicking if an error.
func (c *E2ETestConfig) MustMarshalJSON() []byte {
	jsonBytes, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}

	return jsonBytes
}
