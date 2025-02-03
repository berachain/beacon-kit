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

package backend

import (
	"slices"

	"github.com/berachain/beacon-kit/node-api/backend/utils"
	beacontypes "github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/primitives/math"
)

func (b Backend) FilteredValidators(
	slot math.Slot, ids []string, statuses []string,
) ([]*beacontypes.ValidatorData, error) {
	// TODO: to adhere to the spec, this shouldn't error if the error
	// is not found, but i can't think of a way to do that without coupling
	// db impl to the api impl.
	st, _, err := b.stateFromSlot(slot)
	if err != nil {
		return nil, err
	}

	// Convert requested ids (can be validator index or pubkey) to validator index only.
	validatorIndicies := make([]uint64, 0, len(ids))
	for _, id := range ids {
		validatorIndex, vErr := utils.ValidatorIndexByID(st, id)
		if vErr != nil {
			return nil, vErr
		}
		validatorIndicies = append(validatorIndicies, validatorIndex.Unwrap())
	}

	validators, err := st.GetValidators()
	if err != nil {
		return nil, err
	}

	// Filter on validator indexes and statuses.
	validatorData := make([]*beacontypes.ValidatorData, 0, len(validators))
	for _, validator := range validators {
		// Skip the validator if we are filtering by indicies and this validator is not included.
		index, valErr := st.ValidatorIndexByPubkey(validator.GetPubkey())
		if valErr != nil {
			return nil, err
		}
		if len(validatorIndicies) != 0 && !slices.Contains(validatorIndicies, index.Unwrap()) {
			continue
		}

		// Skip the validator if we are filtering by statuses and this validator is not included.
		status, valErr := utils.GetValidatorStatus(b.cs.SlotToEpoch(slot), validator)
		if valErr != nil {
			return nil, valErr
		}
		if len(statuses) != 0 && !slices.Contains(statuses, status) {
			continue
		}

		balance, valErr := st.GetBalance(index)
		if valErr != nil {
			return nil, valErr
		}
		validatorData = append(validatorData, &beacontypes.ValidatorData{
			ValidatorBalanceData: beacontypes.ValidatorBalanceData{
				Index:   index.Unwrap(),
				Balance: balance.Unwrap(),
			},
			Status:    status,
			Validator: beacontypes.ValidatorFromConsensus(validator),
		})
	}
	return validatorData, nil
}

func (b Backend) ValidatorByID(
	slot math.Slot, id string,
) (*beacontypes.ValidatorData, error) {
	// TODO: to adhere to the spec, this shouldn't error if the error
	// is not found, but i can't think of a way to do that without coupling
	// db impl to the api impl.
	st, _, err := b.stateFromSlot(slot)
	if err != nil {
		return nil, err
	}
	index, err := utils.ValidatorIndexByID(st, id)
	if err != nil {
		return nil, err
	}
	validator, err := st.ValidatorByIndex(index)
	if err != nil {
		return nil, err
	}
	balance, err := st.GetBalance(index)
	if err != nil {
		return nil, err
	}
	status, err := utils.GetValidatorStatus(b.cs.SlotToEpoch(slot), validator)
	if err != nil {
		return nil, err
	}
	return &beacontypes.ValidatorData{
		ValidatorBalanceData: beacontypes.ValidatorBalanceData{
			Index:   index.Unwrap(),
			Balance: balance.Unwrap(),
		},
		Status:    status,
		Validator: beacontypes.ValidatorFromConsensus(validator),
	}, nil
}

func (b Backend) ValidatorBalancesByIDs(
	slot math.Slot, ids []string,
) ([]*beacontypes.ValidatorBalanceData, error) {
	var index math.U64
	st, _, err := b.stateFromSlot(slot)
	if err != nil {
		return nil, err
	}
	balances := make([]*beacontypes.ValidatorBalanceData, 0)
	for _, id := range ids {
		index, err = utils.ValidatorIndexByID(st, id)
		if err != nil {
			return nil, err
		}
		var balance math.U64
		// TODO: same issue as above, shouldn't error on not found.
		balance, err = st.GetBalance(index)
		if err != nil {
			return nil, err
		}
		balances = append(balances, &beacontypes.ValidatorBalanceData{
			Index:   index.Unwrap(),
			Balance: balance.Unwrap(),
		})
	}
	return balances, nil
}
