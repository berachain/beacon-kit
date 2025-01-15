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

package components

import (
	"fmt"
	"os"

	"github.com/berachain/beacon-kit/chain"
	"github.com/berachain/beacon-kit/config/spec"
)

const (
	ChainSpecTypeEnvVar  = "CHAIN_SPEC"
	DevnetChainSpecType  = "devnet"
	MainnetChainSpecType = "mainnet"
	TestnetChainSpecType = "testnet"
)

// ProvideChainSpec provides the chain spec based on the environment variable.
// Defaults to use Mainnet if no valid chain spec environment variable is set.
func ProvideChainSpec() (chain.Spec, error) {
	var (
		chainSpec chain.Spec
		err       error
	)

	// TODO: replace reading env var with config value.
	switch os.Getenv(ChainSpecTypeEnvVar) {
	case DevnetChainSpecType:
		chainSpec, err = spec.DevnetChainSpec()
	case TestnetChainSpecType:
		chainSpec, err = spec.TestnetChainSpec()
	case MainnetChainSpecType:
		
	default:
		chainSpec, err = spec.MainnetChainSpec()
	}

	if err != nil {
		return nil, err
	}
	if chainSpec == nil {
		panic("chain spec is nil")
	}
	return chainSpec, nil
}
