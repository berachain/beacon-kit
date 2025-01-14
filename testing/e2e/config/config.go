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

//nolint:tagliatelle // starlark uses snek case.
package config

import (
	"fmt"
	"github.com/berachain/beacon-kit/primitives/encoding/json"
)

// E2ETestConfig defines the configuration for end-to-end tests, including any
// additional services and validators involved.
type E2ETestConfig struct {
	// NetworkConfiguration specifies the configuration for the network.
	NetworkConfiguration NetworkConfiguration `json:"network_configuration"`
	// NodeSettings specifies the configuration for the nodes in the test.
	NodeSettings NodeSettings `json:"node_settings"`
	// EthJSONRPCEndpoints specifies the RPC endpoints to include in the test.
	EthJSONRPCEndpoints []EthJSONRPCEndpoint `json:"eth_json_rpc_endpoints"`
	// AdditionalServices specifies any extra services that should be included
	// in the test environment.
	AdditionalServices []AdditionalService `json:"additional_services"`
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
	// KZGImpl specifies the KZG implementation to use for the client.
	KZGImpl string `json:"kzg_impl"`
}

// NodeSettings holds the configuration for a single node in the test,
// including client images and types.
type NodeSettings struct {
	// ConsensusSettings holds the configuration for the consensus layer
	// clients.
	ConsensusSettings ConsensusSettings `json:"consensus_settings"`
	// ExecutionSettings holds the configuration for the execution layer
	// clients.
	ExecutionSettings ExecutionSettings `json:"execution_settings"`
}

// ExecutionSettings holds the configuration for the execution layer clients.
type ExecutionSettings struct {
    // Specs holds the node specs for all nodes in the execution layer.
    Specs NodeSpecs `json:"specs"`
    // Images specifies the images available for the execution layer.
    Images map[string]string `json:"images"`
    // ClientConfigs holds specific configurations for different execution clients
    ClientConfigs map[string]interface{} `json:"client_configs,omitempty"`
}

// ConsensusSettings holds the configuration for the consensus layer
// clients.
type ConsensusSettings struct {
	// Specs holds the node specs for all nodes in the consensus layer.
	Specs NodeSpecs `json:"specs"`
	// Images specifies the images available for the consensus layer.
	Images map[string]string `json:"images"`
	// Config specifies the config.toml edits for the consensus layer nodes.
	Config ConsensusConfig `json:"config"`
	// AppConfig specifies the app.toml edits for the consensus layer nodes.
	AppConfig AppConfig `json:"app"`
}

// ConsensusConfig holds the configuration for the consensus layer.
type ConsensusConfig struct {
	// TimeoutPropose specifies the timeout for proposing a block.
	TimeoutPropose string `json:"timeout_propose"`
	// TimeoutPrevote specifies the timeout for prevoting on a block.
	TimeoutPrevote string `json:"timeout_prevote"`
	// TimeoutVote specifies the timeout for precommiting on a block.
	TimeoutPrecommit string `json:"timeout_precommit"`
	// TimeoutCommit specifies the timeout for committing a block.
	TimeoutCommit string `json:"timeout_commit"`
	// MaxNumInboundPeers specifies the maximum number of inbound peers.
	MaxNumInboundPeers int `json:"max_num_inbound_peers"`
	// MaxNumOutboundPeers specifies the maximum number of outbound peers.
	MaxNumOutboundPeers int `json:"max_num_outbound_peers"`
}

// AppConfig holds the configuration for the app layer.
type AppConfig struct {
	// PayloadTimeout specifies the timeout for the payload.
	PayloadTimeout string `json:"payload_timeout"`
	// EnableOptimisticPayloadBuilds enables building the next block's payload
	// optimistically in process-proposal to allow for the execution client to
	// have more time to assemble the block.
	EnableOptimisticPayloadBuilds bool `json:"enable_optimistic_payload_builds"`
}

// NodeSpecs holds the node specs for all nodes in a single layer.
type NodeSpecs struct {
	// MinCPU specifies the minimum number of CPUs to use for all nodes in the
	// layer.
	MinCPU int `json:"min_cpu"`
	// MaxCPU specifies the maximum number of CPUs to use for all nodes in the
	// layer.
	MaxCPU int `json:"max_cpu"`
	// MinMemory specifies the minimum amount of memory to use for all nodes in
	// the layer.
	MinMemory int `json:"min_memory"`
	// MaxMemory specifies the maximum amount of memory to use for all nodes in
	// the layer.
	MaxMemory int `json:"max_memory"`
}

// AdditionalService holds the configuration for an additional service
// to be included in the test.
type AdditionalService struct {
	// Name specifies the name of the additional service.
	Name string `json:"name"`
	// Replicas specifies the number of replicas to use for the service.
	Replicas int `json:"replicas"`
}

// NethermindConfig holds specific configuration for Nethermind client
type NethermindConfig struct {
    // Specific settings for Nethermind
    SyncConfig struct {
        FastSync                  bool `json:"FastSync"`
        DownloadBodiesInFastSync bool `json:"DownloadBodiesInFastSync"`
        DownloadReceiptsInFastSync bool `json:"DownloadReceiptsInFastSync"`
    } `json:"SyncConfig"`
}

// DefaultNethermindConfig returns the default configuration for Nethermind
func DefaultNethermindConfig() *NethermindConfig {
    return &NethermindConfig{
        SyncConfig: struct {
            FastSync                  bool
            DownloadBodiesInFastSync bool
            DownloadReceiptsInFastSync bool
        }{
            FastSync:                  true,
            DownloadBodiesInFastSync: true,
            DownloadReceiptsInFastSync: true,
        },
    }
}

// MustMarshalJSON marshals the E2ETestConfig to JSON, panicking if an error.
func (c *E2ETestConfig) MustMarshalJSON() []byte {
    // Initialize client configs map if needed
    if c.NodeSettings.ExecutionSettings.ClientConfigs == nil {
        c.NodeSettings.ExecutionSettings.ClientConfigs = make(map[string]interface{})
    }

    // Check for Nethermind nodes and add configuration
    for _, nodeSet := range []NodeSet{
        c.NetworkConfiguration.Validators,
        c.NetworkConfiguration.FullNodes,
        c.NetworkConfiguration.SeedNodes,
    } {
        for _, node := range nodeSet.Nodes {
            if node.ElType == "nethermind" && node.Replicas > 0 {
                if _, exists := c.NodeSettings.ExecutionSettings.ClientConfigs["nethermind"]; !exists {
                    c.NodeSettings.ExecutionSettings.ClientConfigs["nethermind"] = DefaultNethermindConfig()
                }
                break
            }
        }
    }

    jsonBytes, err := json.Marshal(c)
    if err != nil {
        panic(err)
    }
    return jsonBytes
}

// ValidateNethermindConfig validates the Nethermind configuration
func (c *E2ETestConfig) ValidateNethermindConfig() error {
    // Check if we have any Nethermind nodes
    nethermindUsed := false
    for _, nodeSet := range []NodeSet{
        c.NetworkConfiguration.Validators,
        c.NetworkConfiguration.FullNodes,
        c.NetworkConfiguration.SeedNodes,
    } {
        for _, node := range nodeSet.Nodes {
            if node.ElType == "nethermind" && node.Replicas > 0 {
                nethermindUsed = true
                break
            }
        }
        if nethermindUsed {
            break
        }
    }

    if !nethermindUsed {
        return nil // No validation needed if Nethermind is not used
    }

    // Check if we have Nethermind configuration
    nethermindConfig, exists := c.NodeSettings.ExecutionSettings.ClientConfigs["nethermind"]
    if !exists {
        return fmt.Errorf("nethermind configuration is missing")
    }

    config, ok := nethermindConfig.(*NethermindConfig)
    if !ok {
        return fmt.Errorf("invalid nethermind configuration type")
    }

    // Validate SyncConfig
    if !config.SyncConfig.FastSync {
        if config.SyncConfig.DownloadBodiesInFastSync || config.SyncConfig.DownloadReceiptsInFastSync {
            return fmt.Errorf("sync options are inconsistent: FastSync is disabled but download options are enabled")
        }
    }

    return nil
}
