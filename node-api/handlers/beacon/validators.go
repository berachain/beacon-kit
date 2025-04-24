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
	"net/http"

	cerrors "github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/node-api/backend"
	"github.com/berachain/beacon-kit/node-api/handlers"
	beacontypes "github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	types "github.com/berachain/beacon-kit/node-api/handlers/types"
	"github.com/berachain/beacon-kit/node-api/handlers/utils"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
)

var ErrNoSlotForStateRoot = errors.New("slot not found at state root")

// getStateValidators is a helper function to provide implementation
// consistency between GetStateValidators and PostStateValidators, since they
// are intended to behave the same way.
func (h *Handler) getStateValidators(stateID string, ids []string, statuses []string) (any, error) {
	if stateID == utils.StateIDGenesis {
		genesisState := h.backend.GenesisState()
		genesisValidators, err := genesisState.GetValidators()
		if err != nil {
			return nil, err
		}
		validators, err := h.backend.FilteredValidatorsAtGenesis(
			genesisValidators,
			genesisState,
			ids,
			statuses,
		)
		if err != nil {
			return nil, err
		}
		return beacontypes.NewResponse(validators), nil
	}
	slot, err := utils.SlotFromStateID(stateID, h.backend)
	if err != nil {
		return nil, err
	}
	validators, err := h.backend.FilteredValidators(
		slot,
		ids,
		statuses,
	)
	if err != nil {
		return nil, err
	}

	return beacontypes.NewResponse(validators), nil
}

func (h *Handler) GetStateValidators(c handlers.Context) (any, error) {
	req, err := utils.BindAndValidate[beacontypes.GetStateValidatorsRequest](
		c, h.Logger(),
	)
	if err != nil {
		return nil, err
	}
	return h.getStateValidators(req.StateID, req.IDs, req.Statuses)
}

func (h *Handler) PostStateValidators(c handlers.Context) (any, error) {
	req, err := utils.BindAndValidate[beacontypes.PostStateValidatorsRequest](
		c, h.Logger(),
	)
	if err != nil {
		return nil, err
	}
	return h.getStateValidators(req.StateID, req.IDs, req.Statuses)
}

func (h *Handler) GetStateValidator(c handlers.Context) (any, error) {
	req, err := utils.BindAndValidate[beacontypes.GetStateValidatorRequest](
		c, h.Logger(),
	)
	if err != nil {
		return nil, err
	}
	if req.StateID == utils.StateIDGenesis {
		st := h.backend.GenesisState()
		validators, err := st.GetValidators()
		if err != nil {
			return nil, err
		}
		return beacontypes.NewResponse(validators), nil
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
	validator, err := h.backend.ValidatorByID(
		slot,
		req.ValidatorID,
	)
	switch {
	case errors.Is(err, backend.ErrValidatorNotFound):
		return &handlers.HTTPError{
			Code:    http.StatusNotFound,
			Message: "Validator not found",
		}, nil
	case err != nil:
		return nil, err
	default:
		return beacontypes.NewResponse(validator), nil
	}
}

func (h *Handler) GetStateValidatorBalances(c handlers.Context) (any, error) {
	req, err := utils.BindAndValidate[beacontypes.GetValidatorBalancesRequest](
		c, h.Logger(),
	)
	if err != nil {
		return nil, err
	}

	var st *statedb.StateDB
	if req.StateID == utils.StateIDGenesis {
		st = h.backend.GenesisState()
	} else {
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
		st, _, err = h.backend.StateAtSlot(slot)
		if err != nil {
			return nil, cerrors.Wrapf(err, "failed to get state from slot %d", slot)
		}
	}

	balances, err := h.backend.ValidatorBalancesByIDs(
		st,
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

	var st *statedb.StateDB

	if req.StateID == utils.StateIDGenesis {
		st = h.backend.GenesisState()
	} else {
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
		st, _, err = h.backend.StateAtSlot(slot)
		if err != nil {
			return nil, cerrors.Wrapf(err, "failed to get state from slot %d", slot)
		}
	}

	balances, err := h.backend.ValidatorBalancesByIDs(
		st,
		req.IDs,
	)
	if err != nil {
		return nil, err
	}
	return beacontypes.NewResponse(balances), nil
}
