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
package kurtosis

import (
	"encoding/json"
)

// E2ETestConfig defines the configuration for end-to-end tests, including any
// additional services and participants involved.
type E2ETestConfig struct {
	// AdditionalServices specifies any extra services that should be included
	// in the test environment.
	AdditionalServices []interface{} `json:"additional_services"`
	// Participants lists the configurations for each participant in the test.
	Participants []Participant `json:"participants"`
}

// Participant holds the configuration for a single participant in the test,
// including client images and types.
type Participant struct {
	// ClClientImage specifies the Docker image to use for the consensus layer
	// client.
	ClClientImage string `json:"cl_client_image"`
	// ClClientType denotes the type of consensus layer client (e.g.,
	// beaconkit).
	ClClientType string `json:"cl_client_type"`
	// ElClientType denotes the type of execution layer client (e.g., reth).
	ElClientType string `json:"el_client_type"`
}

// DefaultE2ETestConfig provides a default configuration for end-to-end tests,
// pre-populating with a standard set of participants and no additional
// services.
func DefaultE2ETestConfig() *E2ETestConfig {
	return &E2ETestConfig{
		AdditionalServices: []interface{}{},
		Participants: []Participant{
			{
				ElClientType:  "geth",
				ClClientImage: "beacond:kurtosis-local",
				ClClientType:  "beaconkit",
			},
			{
				ElClientType:  "reth",
				ClClientImage: "beacond:kurtosis-local",
				ClClientType:  "beaconkit",
			},
			{
				ElClientType:  "reth",
				ClClientImage: "beacond:kurtosis-local",
				ClClientType:  "beaconkit",
			},
			{
				ElClientType:  "reth",
				ClClientImage: "beacond:kurtosis-local",
				ClClientType:  "beaconkit",
			},
		},
	}
}

// AddNodes adds a number of nodes to the E2ETestConfig, using the specified.
func (c *E2ETestConfig) AddNodes(num int, executionClient string) {
	for i := 0; i < num; i++ {
		c.Participants = append(c.Participants, Participant{
			ElClientType:  executionClient,
			ClClientImage: "beacond:kurtosis-local",
			ClClientType:  "beaconkit",
		})
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
