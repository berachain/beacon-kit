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

	"cosmossdk.io/collections"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/node-api/backend"
	"github.com/berachain/beacon-kit/node-api/handlers"
	beacontypes "github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/node-api/handlers/types"
	"github.com/berachain/beacon-kit/node-api/handlers/utils"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
)

func (h *Handler) GetStateValidatorBalances(c handlers.Context) (any, error) {
	req, err := utils.BindAndValidate[beacontypes.GetValidatorBalancesRequest](
		c, h.Logger(),
	)
	if err != nil {
		return nil, err
	}

	height, err := utils.StateIDToHeight(req.StateID, h.backend)
	if err != nil {
		return nil, err
	}
	balances, err := h.getValidatorBalance(height, req.IDs)
	if err != nil {
		return nil, err
	}
	return beacontypes.NewResponse(balances), nil
}

func (h *Handler) PostStateValidatorBalances(c handlers.Context) (any, error) {
	var ids []string
	if err := c.Bind(&ids); err != nil {
		return nil, fmt.Errorf("%s: %w", err.Error(), types.ErrInvalidRequest)
	}
	// Get state_id from URL path parameter
	req := beacontypes.PostValidatorBalancesRequest{
		StateIDRequest: types.StateIDRequest{StateID: c.Param("state_id")},
		IDs:            ids,
	}

	if err := c.Validate(&req); err != nil {
		return nil, types.ErrInvalidRequest
	}

	height, err := utils.StateIDToHeight(req.StateID, h.backend)
	if err != nil {
		return nil, err
	}
	balances, err := h.getValidatorBalance(height, req.IDs)
	if err != nil {
		return nil, err
	}
	return beacontypes.NewResponse(balances), nil
}

func (h *Handler) getValidatorBalance(height int64, validatorIDs []string) ([]*beacontypes.ValidatorBalanceData, error) {
	st, _, err := h.backend.StateAndSlotFromHeight(height)
	if err != nil {
		return nil, fmt.Errorf("failed to get state from height %d: %w", height, err)
	}

	// If no IDs provided, return all validator balances
	if len(validatorIDs) == 0 {
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
		balances = make([]*beacontypes.ValidatorBalanceData, 0, len(validatorIDs))
		index    math.U64
	)
	for _, id := range validatorIDs {
		index, err = validatorIndexByID(st, id)
		switch {
		case err == nil:
			// nothing to do, keep processing
		case errors.Is(err, collections.ErrNotFound):
			// If public key as id is not found in the state
			// we simply skip the index.
			continue
		default:
			return nil, fmt.Errorf("failed to get validator index by id %s: %w", id, err)
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
			return nil, fmt.Errorf("failed to get validator balance for validator index %d: %w", index, err)
		}
	}
	return balances, nil
}

// ValidatorIndexByID parses a validator index from a string.
// The string can be either a validator index or a validator pubkey.
func validatorIndexByID(st backend.ReadOnlyBeaconState, keyOrIndex string) (math.U64, error) {
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
