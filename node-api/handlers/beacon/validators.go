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

	"github.com/berachain/beacon-kit/node-api/handlers"
	beacontypes "github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	types "github.com/berachain/beacon-kit/node-api/handlers/types"
	"github.com/berachain/beacon-kit/node-api/handlers/utils"
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
