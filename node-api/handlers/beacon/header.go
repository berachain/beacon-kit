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

	"github.com/berachain/beacon-kit/node-api/handlers"
	beacontypes "github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/node-api/handlers/utils"
	"github.com/berachain/beacon-kit/primitives/math"
)

func (h *Handler) GetBlockHeaders(c handlers.Context) (any, error) {
	req, err := utils.BindAndValidate[beacontypes.GetBlockHeadersRequest](c, h.Logger())
	if err != nil {
		return nil, err
	}
	slot, err := math.U64FromString(req.Slot)
	if err != nil {
		return nil, fmt.Errorf("failed formatting request slot to type: %w", err)
	}
	return makeBlockHeaderResponse(h.backend, slot)
}

func (h *Handler) GetBlockHeaderByID(c handlers.Context) (any, error) {
	req, err := utils.BindAndValidate[beacontypes.GetBlockHeaderRequest](c, h.Logger())
	if err != nil {
		return nil, err
	}
	slot, err := utils.SlotFromBlockID(req.BlockID, h.backend)
	if err != nil {
		return nil, fmt.Errorf("failed retrieving slot from block ID %s: %w", req.BlockID, err)
	}

	return makeBlockHeaderResponse(h.backend, slot)
}

func makeBlockHeaderResponse(backend Backend, slot math.Slot) (any, error) {
	header, err := backend.BlockHeaderAtSlot(slot)
	if err != nil {
		return nil, fmt.Errorf("failed retrieving header at slot %d: %w", slot, err)
	}

	return beacontypes.NewResponse(
		&beacontypes.BlockHeaderResponse{
			Root:      header.GetBodyRoot(),
			Canonical: true,
			Header: &beacontypes.SignedBeaconBlockHeader{
				Message:   beacontypes.BeaconBlockHeaderFromConsensus(header),
				Signature: "", // TODO: implement
			},
		}), nil
}
