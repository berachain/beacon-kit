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
	"github.com/berachain/beacon-kit/mod/node-api/handlers/types"
	"github.com/berachain/beacon-kit/mod/node-api/handlers/utils"
)

func (h *Handler[ContextT, _, _]) GetStateValidators(c ContextT) (any, error) {
	params, err := utils.BindAndValidate[StateValidatorsGetRequest](
		c,
	)
	if err != nil {
		return nil, err
	}
	// TODO: remove this once status filter is implemented.
	if len(params.Statuses) > 0 {
		return nil, types.ErrNotImplemented
	}
	slot, err := utils.SlotFromStateID(params.StateID)
	if err != nil {
		return nil, err
	}
	validators, err := h.backend.ValidatorsByIDs(
		slot,
		params.IDs,
		params.Statuses,
	)
	if err != nil {
		return nil, err
	}
	if len(validators) == 0 {
		return nil, types.ErrNotFound
	}
	return ValidatorResponse{
		ExecutionOptimistic: false, // stubbed
		Finalized:           false, // stubbed
		Data:                validators,
	}, nil
}

func (h *Handler[ContextT, _, _]) PostStateValidators(c ContextT) (any, error) {
	params, err := utils.BindAndValidate[StateValidatorsPostRequest](
		c,
	)
	if err != nil {
		return nil, err
	}
	// TODO: remove this once status filter is implemented.
	if len(params.Statuses) > 0 {
		return nil, types.ErrNotImplemented
	}
	slot, err := utils.SlotFromStateID(params.StateID)
	if err != nil {
		return nil, err
	}
	validators, err := h.backend.ValidatorsByIDs(
		slot,
		params.IDs,
		params.Statuses,
	)
	if err != nil {
		return nil, err
	}
	return ValidatorResponse{
		ExecutionOptimistic: false, // stubbed
		Finalized:           false, // stubbed
		Data:                validators,
	}, nil
}

func (h *Handler[ContextT, _, _]) GetStateValidatorBalances(
	c ContextT,
) (any, error) {
	params, err := utils.BindAndValidate[ValidatorBalancesGetRequest](
		c,
	)
	if err != nil {
		return nil, err
	}
	slot, err := utils.SlotFromStateID(params.StateID)
	if err != nil {
		return nil, err
	}
	balances, err := h.backend.ValidatorBalancesByIDs(
		slot,
		params.IDs,
	)
	if err != nil {
		return nil, err
	}
	return ValidatorResponse{
		ExecutionOptimistic: false, // stubbed
		Finalized:           false, // stubbed
		Data:                balances,
	}, nil
}

func (h *Handler[ContextT, _, _]) PostStateValidatorBalances(
	c ContextT,
) (any, error) {
	params, err := utils.BindAndValidate[ValidatorBalancesPostRequest](
		c,
	)
	if err != nil {
		return nil, err
	}
	slot, err := utils.SlotFromStateID(params.StateID)
	if err != nil {
		return nil, err
	}
	balances, err := h.backend.ValidatorBalancesByIDs(
		slot,
		params.IDs,
	)
	if err != nil {
		return nil, err
	}
	return ValidatorResponse{
		ExecutionOptimistic: false, // stubbed
		Finalized:           false, // stubbed
		Data:                balances,
	}, nil
}
