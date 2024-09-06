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
	return beacontypes.ValidatorResponse{
		ExecutionOptimistic: false, // stubbed
		Finalized:           false, // stubbed
		Data:                validators,
	}, nil
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
	return beacontypes.ValidatorResponse{
		ExecutionOptimistic: false, // stubbed
		Finalized:           false, // stubbed
		Data:                validators,
	}, nil
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
	return beacontypes.ValidatorResponse{
		ExecutionOptimistic: false, // stubbed
		Finalized:           false, // stubbed
		Data:                balances,
	}, nil
}

//func (h *Handler[_, ContextT, _, _]) PostStateValidatorBalances(
//	c ContextT,
//) (any, error) {
//	h.Logger().Info("PostStateValidatorBalances: Received request")
//
//	// First, let's log the raw request body
//	var rawBody json.RawMessage
//	if err := c.Bind(&rawBody); err != nil {
//		h.Logger().Error("Failed to read raw body", "error", err)
//		return nil, errors.Wrapf(errors.New("err reading raw body"), "err %v", err)
//	}
//	h.Logger().Info("Raw request body", "body", string(rawBody))
//
//	// Try to unmarshal as array of strings
//	var ids []string
//	if err := json.Unmarshal(rawBody, &ids); err == nil {
//		h.Logger().Info("Parsed request as array of strings", "ids", ids)
//		return h.processValidatorBalances(c, "head", ids)
//	}
//
//	// Try to unmarshal as array of integers
//	var indices []uint64
//	if err := json.Unmarshal(rawBody, &indices); err == nil {
//		h.Logger().Info("Parsed request as array of integers", "indices", indices)
//		ids := make([]string, len(indices))
//		for i, index := range indices {
//			ids[i] = fmt.Sprintf("%d", index)
//		}
//		return h.processValidatorBalances(c, "head", ids)
//	}
//
//	// Try to unmarshal as PostValidatorBalancesRequest
//	var req beacontypes.PostValidatorBalancesRequest
//	if err := json.Unmarshal(rawBody, &req); err == nil {
//		h.Logger().Info("Parsed request as PostValidatorBalancesRequest", "request", req)
//		return h.processValidatorBalances(c, req.StateID, req.IDs)
//	}
//
//	h.Logger().Error("Failed to parse request body")
//	return nil, errors.New("invalid request format")
//}

//func (h *Handler[_, ContextT, _, _]) processValidatorBalances(
//	c ContextT,
//	stateID string,
//	ids []string,
//) (any, error) {
//	h.Logger().Info("Processing validator balances", "stateID", stateID, "ids", ids)
//
//	slot, err := utils.SlotFromStateID(stateID, h.backend)
//	if err != nil {
//		return nil, errors.Wrapf(errors.New("err in getting slot"), "slot req err %v %v %v", stateID, slot, err)
//	}
//
//	balances, err := h.backend.ValidatorBalancesBySlot(slot, ids)
//	if err != nil {
//		return nil, errors.Wrapf(errors.New("err in backend"), "err %v", err)
//	}
//
//	return beacontypes.ValidatorResponse{
//		ExecutionOptimistic: false, // stubbed
//		Finalized:           false, // stubbed
//		Data:                balances,
//	}, nil
//}

func (h *Handler[_, ContextT, _, _]) PostStateValidatorBalances(
	c ContextT,
) (any, error) {

	h.Logger().Info("PostStateValidatorBalances: Received request")

	//var ids []string
	//var req beacontypes.PostValidatorBalancesRequest
	//
	//// Try to bind to the new format (array of strings)
	//if err := c.Bind(&ids); err == nil {
	//	req = beacontypes.PostValidatorBalancesRequest{
	//		StateIDRequest: types.StateIDRequest{StateID: "head"},
	//		IDs:            ids,
	//	}
	//} else {
	//	// If binding to array of strings fails, try the original format
	//	if err := c.Bind(&req); err != nil {
	//		return nil, errors.Wrapf(errors.New("err in bind"), "err %v", err)
	//	}
	//	// If IDs field is empty, it might be because the client sent indices as numbers
	//	// We need to convert these to strings
	//	if len(req.IDs) == 0 {
	//		var indices []uint64
	//		if err := c.Bind(&indices); err == nil {
	//			for _, index := range indices {
	//				req.IDs = append(req.IDs, strconv.FormatUint(index, 10))
	//			}
	//		}
	//	}
	//}
	//

	var ids []string
	if err := c.Bind(&ids); err != nil {
		return nil, errors.Wrapf(errors.New("err in func bind"), "err %v", err)
	}

	req := beacontypes.PostValidatorBalancesRequest{
		StateIDRequest: types.StateIDRequest{StateID: "head"},
		IDs:            ids,
	}

	if err := c.Validate(&req); err != nil {
		return nil, errors.Wrapf(errors.New("err in validate"), "err %v", err)
	}

	h.Logger().Info("PostStateValidatorBalances", "req", req)

	slot, err := utils.SlotFromStateID(req.StateID, h.backend)
	if err != nil {
		return nil, errors.Wrapf(errors.New("err in getting slot "), " slot req err %v %v %v", req, slot, err)
	}

	h.Logger().Info("PostStateValidatorBalances", "slot", slot, "req", req)

	balances, err := h.backend.ValidatorBalancesBySlot(slot, req.IDs)
	if err != nil {
		return nil, errors.Wrapf(errors.New("err in backend "), "err %v", err)
	}
	return beacontypes.ValidatorResponse{
		ExecutionOptimistic: false, // stubbed
		Finalized:           false, // stubbed
		Data:                balances,
	}, nil
}
