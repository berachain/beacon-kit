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

package backend

import (
	"github.com/berachain/beacon-kit/node-api/backend/utils"
	beacontypes "github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/primitives/math"
)

func (b Backend[
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, ValidatorT, _, _, _,
]) ValidatorByID(
	slot math.Slot, id string,
) (*beacontypes.ValidatorData[ValidatorT], error) {
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
	return &beacontypes.ValidatorData[ValidatorT]{
		ValidatorBalanceData: beacontypes.ValidatorBalanceData{
			Index:   index.Unwrap(),
			Balance: balance.Unwrap(),
		},
		Status:    "active_ongoing", // TODO: fix
		Validator: validator,
	}, nil
}

// TODO: filter by status
func (b Backend[
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, ValidatorT, _, _, _,
]) ValidatorsByIDs(
	slot math.Slot, ids []string, _ []string,
) ([]*beacontypes.ValidatorData[ValidatorT], error) {
	validatorsData := make([]*beacontypes.ValidatorData[ValidatorT], 0)
	for _, id := range ids {
		// TODO: we can probably optimize this via a getAllValidators
		// query and then filtering but blocked by the fact that IDs
		// can be indices and the hard type only holds its own pubkey.
		validatorData, err := b.ValidatorByID(slot, id)
		if err != nil {
			return nil, err
		}
		validatorsData = append(validatorsData, validatorData)
	}
	return validatorsData, nil
}

func (b Backend[
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) ValidatorBalancesByIDs(
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
