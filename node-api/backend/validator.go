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

	"cosmossdk.io/collections"
	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/node-api/backend/utils"
	beacontypes "github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
)

// ErrValidatorNotFound is an error for when a validator is not found.
var ErrValidatorNotFound = errors.New("validator not found")

// ErrStatusFilterMismatch is an error for when a validator status does not
// match the status filter.
var ErrStatusFilterMismatch = errors.New("validator status does not match status filter")

type validatorFilters struct {
	numericIDs []uint64
	pubkeys    []crypto.BLSPubkey
}

// parseValidatorIDs parses a slice of string IDs into numeric IDs and pubkeys
func parseValidatorIDs(ids []string) *validatorFilters {
	filters := &validatorFilters{
		numericIDs: make([]uint64, 0, len(ids)),
		pubkeys:    make([]crypto.BLSPubkey, 0, len(ids)),
	}

	for _, id := range ids {
		filters.parseID(id)
	}

	return filters
}

// parseID attempts to parse a single ID as either a numeric ID or pubkey
func (f *validatorFilters) parseID(id string) {
	// Try parsing as numeric ID first
	if index, err := math.U64FromString(id); err == nil {
		f.numericIDs = append(f.numericIDs, index.Unwrap())
		return
	}

	// Try parsing as pubkey
	var pubkey crypto.BLSPubkey
	if err := pubkey.UnmarshalText([]byte(id)); err == nil {
		f.pubkeys = append(f.pubkeys, pubkey)
	}
	// Silently skip invalid IDs
}

// FilteredValidators will grab all of the validators from the state at the
// given slot. It will then filter them by the provided ids and statuses.
func (b *Backend) FilteredValidators(
	slot math.Slot, ids []string, statuses []string,
) ([]*beacontypes.ValidatorData, error) {
	st, resolvedSlot, err := b.StateAtSlot(slot)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get state from slot %d", slot)
	}

	validators, err := st.GetValidators()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get validators")
	}

	// Parse all IDs and pubkeys once at the start
	filters := parseValidatorIDs(ids)
	epoch := b.cs.SlotToEpoch(resolvedSlot)

	return filterAndBuildValidatorData(st, validators, filters, epoch, statuses)
}

// filterAndBuildValidatorData processes all validators and builds their data based on filters
func filterAndBuildValidatorData(
	st *statedb.StateDB,
	validators []*types.Validator,
	filters *validatorFilters,
	epoch math.Epoch,
	statuses []string,
) ([]*beacontypes.ValidatorData, error) {
	validatorData := make([]*beacontypes.ValidatorData, 0, len(validators))

	for _, validator := range validators {
		index, err := st.ValidatorIndexByPubkey(validator.GetPubkey())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get validator index by pubkey %s", validator.GetPubkey())
		}

		if !matchesFilters(validator, index, filters) {
			continue
		}

		data, err := buildValidatorData(st, validator, index, epoch, statuses)
		switch {
		case err == nil:
			validatorData = append(validatorData, data)
		case errors.Is(err, ErrStatusFilterMismatch):
			continue
		default:
			return nil, err
		}
	}

	return validatorData, nil
}

// matchesFilters checks if a validator matches the filters.
func matchesFilters(validator *types.Validator, index math.U64, filters *validatorFilters) bool {
	// If no filters, accept all validators
	if len(filters.numericIDs) == 0 && len(filters.pubkeys) == 0 {
		return true
	}

	// Check numeric IDs
	if len(filters.numericIDs) > 0 && matchesIndex(index, filters.numericIDs) {
		return true
	}

	// Check pubkeys
	if len(filters.pubkeys) > 0 && matchesPubkey(validator, filters.pubkeys) {
		return true
	}

	return false
}

func matchesPubkey(validator *types.Validator, parsedPubkeys []crypto.BLSPubkey) bool {
	validatorPubkey := validator.GetPubkey()
	return slices.Contains(parsedPubkeys, validatorPubkey)
}

func matchesIndex(index math.U64, ids []uint64) bool {
	return slices.Contains(ids, index.Unwrap())
}

func matchesStatusFilter(status string, statuses []string) bool {
	return len(statuses) == 0 || slices.Contains(statuses, status)
}

func buildValidatorData(
	st *statedb.StateDB,
	validator *types.Validator,
	index math.U64,
	epoch math.Epoch,
	statuses []string,
) (*beacontypes.ValidatorData, error) {
	status, err := validator.Status(epoch)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get validator status for validator pubkey %s and index %d", validator.GetPubkey(), index)
	}

	if !matchesStatusFilter(status, statuses) {
		return nil, ErrStatusFilterMismatch
	}

	balance, err := st.GetBalance(index)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get validator balance for validator pubkey %s and index %d", validator.GetPubkey(), index)
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

func (b *Backend) ValidatorByID(slot math.Slot, id string) (*beacontypes.ValidatorData, error) {
	// Get the state at the given slot.
	st, resolvedSlot, err := b.StateAtSlot(slot)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get state from slot %d", slot)
	}
	index, err := utils.ValidatorIndexByID(st, id)
	switch {
	case err == nil:
		// continue processing
	case errors.Is(err, collections.ErrNotFound):
		return nil, ErrValidatorNotFound
	default:
		return nil, errors.Wrapf(err, "failed to get validator index by id %s", id)
	}
	validator, err := st.ValidatorByIndex(index)
	switch {
	case err == nil:
		// continue processing
	case errors.Is(err, collections.ErrNotFound):
		return nil, ErrValidatorNotFound
	default:
		return nil, errors.Wrapf(err, "failed to get validator by index %d", index)
	}
	balance, err := st.GetBalance(index)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get validator balance for validator pubkey %s and index %d", validator.GetPubkey(), index)
	}
	status, err := validator.Status(b.cs.SlotToEpoch(resolvedSlot))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get validator status for validator pubkey %s and index %d", validator.GetPubkey(), index)
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

func (b *Backend) ValidatorBalancesByIDs(slot math.Slot, ids []string) ([]*beacontypes.ValidatorBalanceData, error) {
	// Get the state at the given slot.
	st, _, err := b.StateAtSlot(slot)
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

	var (
		balances = make([]*beacontypes.ValidatorBalanceData, 0, len(ids))
		index    math.U64
	)
	for _, id := range ids {
		index, err = utils.ValidatorIndexByID(st, id)
		switch {
		case err == nil:
			// nothing to do, keep processing
		case errors.Is(err, collections.ErrNotFound):
			// If public key as id is not found in the state
			// we simply skip the index.
			continue
		default:
			return nil, errors.Wrapf(err, "failed to get validator index by id %s", id)
		}

		var balance math.U64
		switch balance, err = st.GetBalance(index); {
		case err == nil:
			balances = append(balances, &beacontypes.ValidatorBalanceData{
				Index:   index.Unwrap(),
				Balance: balance.Unwrap(),
			})
		case errors.Is(err, collections.ErrNotFound):
			// if index does not exist and GetBalance returns
			// "collections: not found" we simply skip the index.
			continue
		default:
			return nil, errors.Wrapf(err, "failed to get validator balance for validator index %d", index)
		}
	}
	return balances, nil
}
