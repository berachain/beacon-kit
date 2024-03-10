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
// additional services and validators involved.
type E2ETestConfig struct {
	// AdditionalServices specifies any extra services that should be included
	// in the test environment.
	AdditionalServices []interface{} `json:"additional_services"`
	// Validators lists the configurations for each validator in the test.
	Validators []Validator `json:"validators"`
}

// Validator holds the configuration for a single validator in the test,
// including client images and types.
type Validator struct {
	// ClImage specifies the Docker image to use for the consensus layer
	// client.
	ClImage string `json:"cl_image"`
	// ClType denotes the type of consensus layer client (e.g.,
	// beaconkit).
	ClType string `json:"cl_type"`
	// ElType denotes the type of execution layer client (e.g., reth).
	ElType string `json:"el_type"`
}

// DefaultE2ETestConfig provides a default configuration for end-to-end tests,
// pre-populating with a standard set of validators and no additional
// services.
func DefaultE2ETestConfig() *E2ETestConfig {
	return &E2ETestConfig{
		AdditionalServices: []interface{}{},
		Validators: []Validator{
			{
				ElType:  "geth",
				ClImage: "beacond:kurtosis-local",
				ClType:  "beaconkit",
			},
			{
				ElType:  "reth",
				ClImage: "beacond:kurtosis-local",
				ClType:  "beaconkit",
			},
			{
				ElType:  "reth",
				ClImage: "beacond:kurtosis-local",
				ClType:  "beaconkit",
			},
			{
				ElType:  "nethermind",
				ClImage: "beacond:kurtosis-local",
				ClType:  "beaconkit",
			},
		},
	}
}

// AddNodes adds a number of nodes to the E2ETestConfig, using the specified.
func (c *E2ETestConfig) AddNodes(num int, executionClient string) {
	for i := 0; i < num; i++ {
		c.Validators = append(c.Validators, Validator{
			ElType:  executionClient,
			ClImage: "beacond:kurtosis-local",
			ClType:  "beaconkit",
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
