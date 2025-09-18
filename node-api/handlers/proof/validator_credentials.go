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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

//nolint:dupl // False positive detected.
package proof

import (
	"github.com/berachain/beacon-kit/node-api/handlers"
	"github.com/berachain/beacon-kit/node-api/handlers/proof/merkle"
	"github.com/berachain/beacon-kit/node-api/handlers/proof/types"
	"github.com/berachain/beacon-kit/node-api/handlers/utils"
	"github.com/berachain/beacon-kit/primitives/math"
)

// GetValidatorCredentials returns the withdrawal credentials of a validator along with a
// Merkle proof that can be verified against the beacon block root.
func (h *Handler) GetValidatorCredentials(c handlers.Context) (any, error) {
	params, err := utils.BindAndValidate[types.ValidatorCredentialsRequest](c, h.Logger())
	if err != nil {
		return nil, err
	}

	// Validator index is provided as a string path parameter; convert to math.U64.
	validatorIndex, err := math.U64FromString(params.ValidatorIndex)
	if err != nil {
		return nil, err
	}

	slot, beaconState, blockHeader, err := h.resolveTimestampID(params.TimestampID)
	if err != nil {
		return nil, err
	}

	h.Logger().Info(
		"Generating withdrawal credential proof", "slot", slot, "validator_index", validatorIndex,
	)

	// Generate proof for withdrawal credentials in the block.
	bsm, err := beaconState.GetMarshallable()
	if err != nil {
		return nil, err
	}

	credsProof, beaconBlockRoot, err := merkle.ProveWithdrawalCredentialsInBlock(
		validatorIndex, blockHeader, bsm,
	)
	if err != nil {
		return nil, err
	}

	// Fetch the validator to get the withdrawal credentials to include in the response.
	validator, err := beaconState.ValidatorByIndex(validatorIndex)
	if err != nil {
		return nil, err
	}

	return types.ValidatorWithdrawalCredentialsResponse{
		BeaconBlockHeader:              blockHeader,
		BeaconBlockRoot:                beaconBlockRoot,
		ValidatorWithdrawalCredentials: validator.GetWithdrawalCredentials(),
		WithdrawalCredentialsProof:     credsProof,
	}, nil
}
