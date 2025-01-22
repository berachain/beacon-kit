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

package genesis

import (
	"github.com/berachain/beacon-kit/chain"
	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/encoding/json"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/spf13/afero"
)

// Beacon, AppState and Genesis are code duplications that
// collectively reproduce part of genesis file structure

type Beacon struct {
	Deposits types.Deposits `json:"deposits"`
}

type AppState struct {
	Beacon `json:"beacon"`
}

type Genesis struct {
	AppState `json:"app_state"`
}

// GetValidatorRootFromFile returns the genesis validator root from a given beacond genesis file.
func GetValidatorRootFromFile(filepath string, cs chain.Spec) (common.Root, error) {
	// Read the genesis file.
	genesisBz, err := afero.ReadFile(afero.NewOsFs(), filepath)
	if err != nil {
		return common.Root{}, errors.Wrap(err, "failed to genesis json file")
	}

	var genesis Genesis
	// Unmarshal JSON data into the Genesis struct
	err = json.Unmarshal(genesisBz, &genesis)
	if err != nil {
		return common.Root{}, errors.Wrap(err, "failed to unmarshal JSON")
	}

	return ValidatorsRoot(genesis.Deposits, cs), nil
}

func ValidatorsRoot(deposits types.Deposits, cs chain.Spec) common.Root {
	validators := make(types.Validators, len(deposits))
	minEffectiveBalance := math.Gwei(cs.EjectionBalance() + cs.EffectiveBalanceIncrement())

	for i, deposit := range deposits {
		val := types.NewValidatorFromDeposit(
			deposit.Pubkey,
			deposit.Credentials,
			deposit.Amount,
			math.Gwei(cs.EffectiveBalanceIncrement()),
			math.Gwei(cs.MaxEffectiveBalance()),
		)

		// mimic processGenesisActivation
		if val.GetEffectiveBalance() >= minEffectiveBalance {
			val.SetActivationEligibilityEpoch(0)
			val.SetActivationEpoch(0)
		}
		validators[i] = val
	}

	return validators.HashTreeRoot()
}
