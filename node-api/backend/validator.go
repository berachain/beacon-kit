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
	"bytes"
	"slices"
	"strconv"

	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/node-api/backend/utils"
	beacontypes "github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
)

// FilteredValidators will grab all of the validators from the state at the
// given slot. It will then filter them by the provided ids and statuses.
func (b Backend) FilteredValidators(
	slot math.Slot, ids []string, statuses []string,
) ([]*beacontypes.ValidatorData, error) {
	st, _, err := b.stateFromSlot(slot)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get state from slot %d", slot)
	}

	validators, err := st.GetValidators()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get validators")
	}

	validatorData := make([]*beacontypes.ValidatorData, 0, len(validators))
	for _, validator := range validators {
		data, errInProcess := b.processValidator(st, validator, slot, ids, statuses)
		if errInProcess != nil {
			return nil, errors.Wrapf(errInProcess, "failed to process validator")
		}
		if data != nil {
			validatorData = append(validatorData, data)
		}
	}
	return validatorData, nil
}

func (b Backend) processValidator(
	st *statedb.StateDB,
	validator *types.Validator,
	slot math.Slot,
	ids []string,
	statuses []string,
) (*beacontypes.ValidatorData, error) {
	index, err := st.ValidatorIndexByPubkey(validator.GetPubkey())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get validator index by pubkey")
	}

	if !b.matchesIDFilter(validator, index, ids) {
		//nolint:nilnil // no data to return
		return nil, nil
	}

	status, err := validator.Status(b.cs.SlotToEpoch(slot))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get validator status")
	}

	if !b.matchesStatusFilter(status, statuses) {
		//nolint:nilnil // no data to return
		return nil, nil
	}

	return b.buildValidatorData(st, validator, index, status)
}

func (b Backend) matchesIDFilter(validator *types.Validator, index math.U64, ids []string) bool {
	if len(ids) == 0 {
		return true
	}

	for _, id := range ids {
		if b.matchesPubkey(validator, id) || b.matchesIndex(index, id) {
			return true
		}
	}
	return false
}

func (b Backend) matchesPubkey(validator *types.Validator, id string) bool {
	var pubkey crypto.BLSPubkey
	if err := pubkey.UnmarshalText([]byte(id)); err != nil {
		return false
	}
	validatorPubkey := validator.GetPubkey()
	return bytes.Equal(validatorPubkey[:], pubkey[:])
}

func (b Backend) matchesIndex(index math.U64, id string) bool {
	return strconv.FormatUint(index.Unwrap(), 10) == id
}

func (b Backend) matchesStatusFilter(status string, statuses []string) bool {
	return len(statuses) == 0 || slices.Contains(statuses, status)
}

func (b Backend) buildValidatorData(
	st *statedb.StateDB,
	validator *types.Validator,
	index math.U64,
	status string,
) (*beacontypes.ValidatorData, error) {
	balance, err := st.GetBalance(index)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get validator balance")
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

func (b Backend) ValidatorByID(
	slot math.Slot, id string,
) (*beacontypes.ValidatorData, error) {
	st, _, err := b.stateFromSlot(slot)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get state from slot %d", slot)
	}
	index, err := utils.ValidatorIndexByID(st, id)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get validator index by id %s", id)
	}
	validator, err := st.ValidatorByIndex(index)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get validator by index %d", index)
	}
	balance, err := st.GetBalance(index)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get validator balance")
	}
	status, err := validator.Status(b.cs.SlotToEpoch(slot))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get validator status")
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
		return nil, errors.Wrapf(err, "failed to get state from slot %d", slot)
	}
	balances := make([]*beacontypes.ValidatorBalanceData, 0)
	for _, id := range ids {
		index, err = utils.ValidatorIndexByID(st, id)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get validator index by id %s", id)
		}
		var balance math.U64
		// TODO: same issue as above, shouldn't error on not found.
		balance, err = st.GetBalance(index)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get validator balance")
		}
		balances = append(balances, &beacontypes.ValidatorBalanceData{
			Index:   index.Unwrap(),
			Balance: balance.Unwrap(),
		})
	}
	return balances, nil
}
