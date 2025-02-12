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

	"cosmossdk.io/collections"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/node-api/backend/utils"
	beacontypes "github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
)

var ErrValidatorNotFound = errors.New("validator not found")

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
		index, valErr := st.ValidatorIndexByPubkey(validator.GetPubkey())
		if valErr != nil {
			return nil, errors.Wrapf(valErr, "failed to get validator index by pubkey")
		}

		// If filtering by IDs, check if this validator matches any ID (pubkey or index)
		if len(ids) > 0 {
			found := false
			for _, id := range ids {
				// Try as pubkey first
				var pubkey crypto.BLSPubkey
				if err := pubkey.UnmarshalText([]byte(id)); err == nil {
					validatorPubkey := validator.GetPubkey()
					if bytes.Equal(validatorPubkey[:], pubkey[:]) {
						found = true
						break
					}
					continue
				}

				// Try as index
				if idxStr := strconv.FormatUint(index.Unwrap(), 10); idxStr == id {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Skip the validator if we are filtering by statuses and this validator is not included.
		status, valErr := validator.Status(b.cs.SlotToEpoch(slot))
		if valErr != nil {
			return nil, errors.Wrapf(valErr, "failed to get validator status")
		}
		if len(statuses) != 0 && !slices.Contains(statuses, status) {
			continue
		}

		balance, valErr := st.GetBalance(index)
		if valErr != nil {
			return nil, errors.Wrapf(valErr, "failed to get validator balance")
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
	st, _, err := b.stateFromSlot(slot)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get state from slot %d", slot)
	}
	index, err := utils.ValidatorIndexByID(st, id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			//nolint:nilnil // The response should be nil without an error.
			return nil, nil
		}
		return nil, err
	}
	validator, err := st.ValidatorByIndex(index)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			//nolint:nilnil // The response should be nil without an error.
			return nil, nil
		}
		return nil, err
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
	st, _, err := b.stateFromSlot(slot)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get state from slot %d", slot)
	}

	// If no IDs provided, return all validator balances
	if len(ids) == 0 {
		rawBalances, errInBalances := st.GetBalances()
		if errInBalances != nil {
			return nil, errInBalances
		}
		// Convert []uint64 to []*ValidatorBalanceData as per the API spec
		balances := make([]*beacontypes.ValidatorBalanceData, len(rawBalances))
		for i, balance := range rawBalances {
			balances[i] = &beacontypes.ValidatorBalanceData{
				Index:   uint64(i), // #nosec:G115 // Safe as i comes from range loop
				Balance: balance,
			}
		}
		return balances, nil
	}

	balances := make([]*beacontypes.ValidatorBalanceData, 0)
	var index math.U64

	for _, id := range ids {
		index, err = utils.ValidatorIndexByID(st, id)
		if err != nil {
			// If public key as id is not found in the state, do not return an error.
			if errors.Is(err, collections.ErrNotFound) {
				continue
			}
			return nil, err
		}
		var balance math.U64
		// TODO: same issue as above, shouldn't error on not found.
		balance, err = st.GetBalance(index)

		if err != nil {
			// if index does not exist and GetBalance returns an error containing "collections: not found"
			// do not return an error.
			if errors.Is(err, collections.ErrNotFound) {
				continue
			}
			return nil, err
		}
		balances = append(balances, &beacontypes.ValidatorBalanceData{
			Index:   index.Unwrap(),
			Balance: balance.Unwrap(),
		})
	}
	return balances, nil
}
