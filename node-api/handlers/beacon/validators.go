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
	"errors"
	"fmt"
	"net/http"

	"github.com/berachain/beacon-kit/node-api/handlers"
	beacontypes "github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	types "github.com/berachain/beacon-kit/node-api/handlers/types"
	"github.com/berachain/beacon-kit/node-api/handlers/utils"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
)

func (h *Handler) GetStateValidators(c handlers.Context) (any, error) {
	req, err := utils.BindAndValidate[beacontypes.GetStateValidatorsRequest](
		c, h.Logger(),
	)
	if err != nil {
		return nil, err
	}

	slot, err := utils.SlotFromStateID(req.StateID, h.backend)
	if err != nil {
		return nil, err
	}
	filteredVals, err := h.FilterValidators(slot, req.IDs, req.Statuses)
	if err != nil {
		return nil, fmt.Errorf("failed to filter validators: %w", err)
	}
	return beacontypes.NewResponse(filteredVals), nil
}

func (h *Handler) PostStateValidators(c handlers.Context) (any, error) {
	req, err := utils.BindAndValidate[beacontypes.PostStateValidatorsRequest](
		c, h.Logger(),
	)
	if err != nil {
		return nil, err
	}

	slot, err := utils.SlotFromStateID(req.StateID, h.backend)
	if err != nil {
		return nil, err
	}
	filteredVals, err := h.FilterValidators(slot, req.IDs, req.Statuses)
	if err != nil {
		return nil, fmt.Errorf("failed to filter validators: %w", err)
	}
	return beacontypes.NewResponse(filteredVals), nil
}

func (h *Handler) GetStateValidator(c handlers.Context) (any, error) {
	req, err := utils.BindAndValidate[beacontypes.GetStateValidatorRequest](
		c, h.Logger(),
	)
	if err != nil {
		return nil, err
	}

	// retrieve slot and associated state
	slot, err := utils.SlotFromStateID(req.StateID, h.backend)
	if err != nil {
		if errors.Is(err, utils.ErrNoSlotForStateRoot) {
			return nil, fmt.Errorf("%s: %w", err.Error(), types.ErrNotFound)
		}
		return nil, fmt.Errorf("failed mapping state id %s to slot: %w", req.StateID, err)
	}
	st, resolvedSlot, err := h.backend.StateAtSlot(slot)
	if err != nil {
		return nil, fmt.Errorf("failed to get state from slot %d: %w", slot, err)
	}

	// retrieve validator data
	index, err := validatorIndexByID(st, req.ValidatorID)
	if err != nil {
		if errors.Is(err, utils.ErrNoSlotForStateRoot) {
			return nil, fmt.Errorf("%s: %w", err.Error(), types.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get validator index by id %s: %w", req.ValidatorID, err)
	}

	validator, err := st.ValidatorByIndex(index)
	if err != nil {
		if errors.Is(err, utils.ErrNoSlotForStateRoot) {
			return nil, fmt.Errorf("%s: %w", err.Error(), types.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get validator by index %s: %w", req.ValidatorID, err)
	}

	balance, err := st.GetBalance(index)
	if err != nil {
		return nil, fmt.Errorf("failed to get validator balance for validator pubkey %s and index %d: %w", validator.GetPubkey(), index, err)
	}
	status, err := validator.Status(h.cs.SlotToEpoch(resolvedSlot))
	if err != nil {
		return nil, fmt.Errorf("failed to get validator status for validator pubkey %s and index %d: %w", validator.GetPubkey(), index, err)
	}
	return beacontypes.NewResponse(
		&beacontypes.ValidatorData{
			ValidatorBalanceData: beacontypes.ValidatorBalanceData{
				Index:   index.Unwrap(),
				Balance: balance.Unwrap(),
			},
			Status:    status,
			Validator: beacontypes.ValidatorFromConsensus(validator),
		},
	), nil
}

// ValidatorIndexByID parses a validator index from a string.
// The string can be either a validator index or a validator pubkey.
func validatorIndexByID(st *statedb.StateDB, keyOrIndex string) (math.U64, error) {
	index, err := math.U64FromString(keyOrIndex)
	if err == nil {
		return index, nil
	}
	var key crypto.BLSPubkey
	if err = key.UnmarshalText([]byte(keyOrIndex)); err != nil {
		return math.U64(0), err
	}
	return st.ValidatorIndexByPubkey(key)
}

func (h *Handler) GetStateValidatorBalances(c handlers.Context) (any, error) {
	req, err := utils.BindAndValidate[beacontypes.GetValidatorBalancesRequest](
		c, h.Logger(),
	)
	if err != nil {
		return nil, err
	}
	slot, err := utils.SlotFromStateID(req.StateID, h.backend)

	switch {
	case err == nil:
		// No error, continue
	case errors.Is(err, utils.ErrNoSlotForStateRoot):
		return &handlers.HTTPError{
			Code:    http.StatusNotFound,
			Message: "State not found",
		}, nil
	default:
		return nil, err
	}
	balances, err := h.backend.ValidatorBalancesByIDs(
		slot,
		req.IDs,
	)
	if err != nil {
		return nil, err
	}
	return beacontypes.NewResponse(balances), nil
}

func (h *Handler) PostStateValidatorBalances(c handlers.Context) (any, error) {
	var ids []string
	if err := c.Bind(&ids); err != nil {
		return nil, types.ErrInvalidRequest
	}
	// Get state_id from URL path parameter
	req := beacontypes.PostValidatorBalancesRequest{
		StateIDRequest: types.StateIDRequest{StateID: c.Param("state_id")},
		IDs:            ids,
	}

	if err := c.Validate(&req); err != nil {
		return nil, types.ErrInvalidRequest
	}

	slot, err := utils.SlotFromStateID(req.StateID, h.backend)
	switch {
	case err == nil:
		// No error, continue
	case errors.Is(err, utils.ErrNoSlotForStateRoot):
		return &handlers.HTTPError{
			Code:    http.StatusNotFound,
			Message: "State not found",
		}, nil
	default:
		return nil, err
	}
	balances, err := h.backend.ValidatorBalancesByIDs(
		slot,
		req.IDs,
	)
	if err != nil {
		return nil, err
	}
	return beacontypes.NewResponse(balances), nil
}
