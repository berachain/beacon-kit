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
	"encoding/json"
	"fmt"
	"github.com/berachain/beacon-kit/mod/chain-spec/pkg/chain"
	"github.com/berachain/beacon-kit/mod/config/pkg/spec"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"os"
	"strings"
)

const (
	ChainSpecTypeEnvVar = "CHAIN_SPEC"
	DevnetChainSpecType = "devnet"
	BetnetChainSpecType = "betnet"
)

// ProvideChainSpec provides the chain spec based on the environment variable.
func ProvideChainSpec() common.ChainSpec {
	// TODO: This is still pretty hood but we shouldn't deviate to far from upstream
	specType := os.Getenv(ChainSpecTypeEnvVar)
	if specType == "" {
		panic(fmt.Sprintf("environment variable %s not set", ChainSpecTypeEnvVar))
	}

	var chainSpec common.ChainSpec
	if strings.HasPrefix(specType, "file://") {
		specPath := strings.TrimPrefix(specType, "file://")
		if specPath == "" {
			panic(fmt.Sprintf("file path not set: %s", specPath))
		}

		b, err := os.ReadFile(specPath)
		if err != nil {
			panic(fmt.Sprintf("Failed to open chain specification file: %w", err))
		}

		chainSpecInput := new(spec.ChainSpecInput)
		err = json.Unmarshal(b, &chainSpecInput)
		if err != nil {
			panic(fmt.Sprintf("Failed to unmarshal chain specification: %w", err))
		}

		sd := spec.BaseSpec()
		sd.DepositEth1ChainID = chainSpecInput.Eth1ChainID
		chainSpec = chain.NewChainSpec(sd)

		return chainSpec
	}

	switch specType {
	case DevnetChainSpecType:
		chainSpec = spec.DevnetChainSpec()
	case BetnetChainSpecType:
		chainSpec = spec.BetnetChainSpec()
	default:
		chainSpec = spec.TestnetChainSpec()
	}
	return chainSpec

}
