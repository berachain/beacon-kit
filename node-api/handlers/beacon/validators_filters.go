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

package beacon

import (
	"fmt"
	"slices"

	consensustypes "github.com/berachain/beacon-kit/consensus-types/types"
	cometbft "github.com/berachain/beacon-kit/consensus/cometbft/service"
	"github.com/berachain/beacon-kit/errors"
	beacontypes "github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	handlertypes "github.com/berachain/beacon-kit/node-api/handlers/types"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// errStatusFilterMismatch is an error for when a validator status does not
// match the status filter.
var errStatusFilterMismatch = errors.New("validator status does not match status filter")

// FilterValidators is a helper function to provide implementation
// consistency between GetStateValidators and PostStateValidators, since they
// are intended to behave the same way.
func (h *Handler) FilterValidators(slot math.Slot, ids []string, statuses []string) ([]*beacontypes.ValidatorData, error) {
	st, resolvedSlot, err := h.backend.StateAtSlot(slot)
	if err != nil {
		if errors.Is(err, cometbft.ErrAppNotReady) {
			// chain not ready, like when genesis time is set in the future
			return nil, handlertypes.ErrNotFound
		}
		if errors.Is(err, sdkerrors.ErrInvalidHeight) {
			// height requested too high
			return nil, handlertypes.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get state from slot %d: %w", slot, err)
	}

	allVals, err := st.GetValidators()
	if err != nil {
		return nil, fmt.Errorf("failed to get validators: %w", err)
	}

	// Parse all IDs and pubkeys once at the start
	filters := parseValidatorIDs(ids)
	epoch := h.cs.SlotToEpoch(resolvedSlot)

	return filterAndBuildValidatorData(st, allVals, filters, epoch, statuses)
}

type validatorFilters struct {
	indexes []uint64
	pubkeys []crypto.BLSPubkey
}

// parseID attempts to parse a single ID as either a numeric ID or pubkey
func (f *validatorFilters) parseID(id string) {
	// Try parsing as numeric ID first
	if index, err := math.U64FromString(id); err == nil {
		f.indexes = append(f.indexes, index.Unwrap())
		return
	}

	// Try parsing as pubkey
	var pubkey crypto.BLSPubkey
	if err := pubkey.UnmarshalText([]byte(id)); err == nil {
		f.pubkeys = append(f.pubkeys, pubkey)
	}
	// We can skip errors here, since they should not happen.
	// We do validate these ids in ValidateValidatorID.
}

// parseValidatorIDs parses a slice of string IDs into numeric IDs and pubkeys
func parseValidatorIDs(ids []string) *validatorFilters {
	filters := &validatorFilters{
		indexes: make([]uint64, 0, len(ids)),
		pubkeys: make([]crypto.BLSPubkey, 0, len(ids)),
	}

	for _, id := range ids {
		filters.parseID(id)
	}

	return filters
}

// filterAndBuildValidatorData processes all validators and builds their data based on filters
func filterAndBuildValidatorData(
	st *statedb.StateDB,
	validators []*consensustypes.Validator,
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
		case errors.Is(err, errStatusFilterMismatch):
			continue
		default:
			return nil, err
		}
	}

	return validatorData, nil
}

// matchesFilters checks if a validator matches the filters.
func matchesFilters(validator *consensustypes.Validator, index math.U64, filters *validatorFilters) bool {
	// If no filters, accept all validators
	if len(filters.indexes) == 0 && len(filters.pubkeys) == 0 {
		return true
	}

	// Check numeric IDs
	if len(filters.indexes) > 0 && matchesIndex(index, filters.indexes) {
		return true
	}

	// Check pubkeys
	if len(filters.pubkeys) > 0 && matchesPubkey(validator, filters.pubkeys) {
		return true
	}

	return false
}

func matchesPubkey(validator *consensustypes.Validator, parsedPubkeys []crypto.BLSPubkey) bool {
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
	validator *consensustypes.Validator,
	index math.U64,
	epoch math.Epoch,
	statuses []string,
) (*beacontypes.ValidatorData, error) {
	status, err := validator.Status(epoch)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get validator status for validator pubkey %s and index %d", validator.GetPubkey(), index)
	}

	if !matchesStatusFilter(status, statuses) {
		return nil, errStatusFilterMismatch
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
