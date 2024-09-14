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

package beacon

import (
	"github.com/berachain/beacon-kit/mod/errors"
	beacontypes "github.com/berachain/beacon-kit/mod/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/mod/node-api/handlers/types"
	"github.com/berachain/beacon-kit/mod/node-api/handlers/utils"
)

func (h *Handler[_, ContextT, _, _]) GetStateValidators(
	c ContextT,
) (any, error) {
	req, err := utils.BindAndValidate[beacontypes.GetStateValidatorsRequest](
		c, h.Logger(),
	)
	if err != nil {
		return nil, err
	}
	// TODO: remove this once status filter is implemented.
	if len(req.Statuses) > 0 {
		return nil, types.ErrNotImplemented
	}
	slot, err := utils.SlotFromStateID(req.StateID, h.backend)
	if err != nil {
		return nil, err
	}
	validators, err := h.backend.ValidatorsByIDs(
		slot,
		req.IDs,
		req.Statuses,
	)
	if err != nil {
		return nil, err
	}
	if len(validators) == 0 {
		return nil, types.ErrNotFound
	}
	return types.Wrap(validators), nil
}

func (h *Handler[_, ContextT, _, _]) PostStateValidators(
	c ContextT,
) (any, error) {
	req, err := utils.BindAndValidate[beacontypes.PostStateValidatorsRequest](
		c, h.Logger(),
	)
	if err != nil {
		return nil, err
	}
	// TODO: remove this once status filter is implemented.
	if len(req.Statuses) > 0 {
		return nil, types.ErrNotImplemented
	}
	slot, err := utils.SlotFromStateID(req.StateID, h.backend)
	if err != nil {
		return nil, err
	}
	validators, err := h.backend.ValidatorsByIDs(
		slot,
		req.IDs,
		req.Statuses,
	)
	if err != nil {
		return nil, err
	}

	// Convert validators to ValidatorData.
	validatorResponse := make(
		[]beacontypes.ValidatorData,
		len(validators),
	)
	for i, validator := range validators {
		validatorResponse[i] = beacontypes.ValidatorData{
			ValidatorBalanceData: beacontypes.ValidatorBalanceData{
				Index:   validator.Index,
				Balance: validator.Balance,
			},
			Status:    validator.Status,
			Validator: validator.Validator,
		}
	}

	return types.Wrap(validatorResponse), nil
}

func (h *Handler[_, ContextT, _, _]) GetStateValidator(
	c ContextT,
) (any, error) {
	req, err := utils.BindAndValidate[beacontypes.GetStateValidatorRequest](
		c, h.Logger(),
	)
	if err != nil {
		return nil, err
	}
	slot, err := utils.SlotFromStateID(req.StateID, h.backend)
	if err != nil {
		return nil, err
	}
	validator, err := h.backend.ValidatorByID(
		slot,
		req.ValidatorID,
	)
	if err != nil {
		return nil, err
	}
	return validator, nil
}

func (h *Handler[_, ContextT, _, _]) GetStateValidatorBalances(
	c ContextT,
) (any, error) {
	req, err := utils.BindAndValidate[beacontypes.GetValidatorBalancesRequest](
		c, h.Logger(),
	)
	if err != nil {
		return nil, err
	}
	slot, err := utils.SlotFromStateID(req.StateID, h.backend)
	if err != nil {
		return nil, err
	}
	balances, err := h.backend.ValidatorBalancesByIDs(
		slot,
		req.IDs,
	)
	if err != nil {
		return nil, err
	}
	return types.Wrap(balances), nil
}

func (h *Handler[_, ContextT, _, _]) PostStateValidatorBalances(
	c ContextT,
) (any, error) {
	var ids []string
	if err := c.Bind(&ids); err != nil {
		return nil, types.ErrInvalidRequest
	}

	// TODO: Find a way to pass the state_id from request.
	// Currently only head is supported.
	req := beacontypes.PostValidatorBalancesRequest{
		StateIDRequest: types.StateIDRequest{StateID: "head"},
		IDs:            ids,
	}

	if err := c.Validate(&req); err != nil {
		return nil, types.ErrInvalidRequest
	}

	slot, err := utils.SlotFromStateID(req.StateID, h.backend)
	if err != nil {
		return nil, errors.Wrapf(err, "err getting slot for req %v ", req)
	}

	h.Logger().Info("PostStateValidatorBalances", "slot", slot, "req", req)

	balances, err := h.backend.ValidatorBalancesByIDs(slot, req.IDs)
	if err != nil {
		return nil, errors.Wrap(err, "err in backend")
	}
	return types.Wrap(balances), nil
}
