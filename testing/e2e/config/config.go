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
	AdditionalServices []any `json:"additional_services"`
	// Validators lists the configurations for each validator in the test.
	Validators []Node `json:"validators"`
	// FullNodes specifies the number of full nodes to include in the test.
	FullNodes []Node `json:"full_nodes"`
	// RPCEndpoints specifies the RPC endpoints to include in the test.
	RPCEndpoints []RPCEndpoint `json:"rpc_endpoints"`
}

type RPCEndpoint struct {
	Type     string   `json:"type"`
	Services []string `json:"services"`
}

// Validator holds the configuration for a single validator in the test,
// including client images and types.
type Node struct {
	// ClImage specifies the Docker image to use for the consensus layer
	// client.
	ClImage string `json:"cl_image"`
	// ClType denotes the type of consensus layer client (e.g.,
	// beaconkit).
	ClType string `json:"cl_type"`
	// ElType denotes the type of execution layer client (e.g., reth).
	ElType string `json:"el_type"`
	// Replicas specifies the number of replicas to use for the client.
	Replicas int `json:"replicas"`
}

// DefaultE2ETestConfig provides a default configuration for end-to-end tests,
// pre-populating with a standard set of validators and no additional
// services.
func DefaultE2ETestConfig() *E2ETestConfig {
	return &E2ETestConfig{
		AdditionalServices: []any{
			"tx-fuzz",
		},
		Validators: []Node{
			{
				ElType:   "nethermind",
				ClImage:  "beacond:kurtosis-local",
				ClType:   "beaconkit",
				Replicas: 1,
			},
			{
				ElType:   "geth",
				ClImage:  "beacond:kurtosis-local",
				ClType:   "beaconkit",
				Replicas: 1,
			},
			{
				ElType:   "reth",
				ClImage:  "beacond:kurtosis-local",
				ClType:   "beaconkit",
				Replicas: 1,
			},
			// {
			// 	ElType:   "erigon",
			// 	ClImage:  "beacond:kurtosis-local",
			// 	ClType:   "beaconkit",
			// 	Replicas: 1,
			// },
			// {
			// 	ElType:   "besu",
			// 	ClImage:  "beacond:kurtosis-local",
			// 	ClType:   "beaconkit",
			// 	Replicas: 1,
			// },
		},
		FullNodes: []Node{
			{
				ElType:   "nethermind",
				ClImage:  "beacond:kurtosis-local",
				ClType:   "beaconkit",
				Replicas: 1,
			},
			{
				ElType:   "reth",
				ClImage:  "beacond:kurtosis-local",
				ClType:   "beaconkit",
				Replicas: 1,
			},
			{
				ElType:   "geth",
				ClImage:  "beacond:kurtosis-local",
				ClType:   "beaconkit",
				Replicas: 1,
			},
			// {
			// 	ElType:   "erigon",
			// 	ClImage:  "beacond:kurtosis-local",
			// 	ClType:   "beaconkit",
			// 	Replicas: 1,
			// },
			// {
			// 	ElType:   "besu",
			// 	ClImage:  "beacond:kurtosis-local",
			// 	ClType:   "beaconkit",
			// 	Replicas: 1,
			// },
		},
		RPCEndpoints: []RPCEndpoint{
			{
				Type: "nginx",
				Services: []string{
					"el-full-nethermind-0:8545",
					"el-full-reth-1:8545",
					"el-full-geth-2:8545",
					// "el-full-erigon-3:8545",
					// Besu causing flakey tests.
					// "el-full-besu-4:8545",
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
