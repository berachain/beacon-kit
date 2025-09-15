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

	cometbft "github.com/berachain/beacon-kit/consensus/cometbft/service"
	"github.com/berachain/beacon-kit/node-api/handlers"
	beacontypes "github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	handlertypes "github.com/berachain/beacon-kit/node-api/handlers/types"
	"github.com/berachain/beacon-kit/node-api/handlers/utils"
)

func (h *Handler) GetGenesis(handlers.Context) (any, error) {
	st, _, err := h.backend.StateAtSlot(utils.Genesis)
	if err != nil {
		if errors.Is(err, cometbft.ErrAppNotReady) {
			// chain not ready, like when genesis time is set in the future
			return nil, handlertypes.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get state from genesis: %w", err)
	}

	genesisRoot, err := st.GetGenesisValidatorsRoot()
	if err != nil {
		return nil, err
	}

	genesisFork, err := st.GetFork()
	if err != nil {
		return nil, err
	}

	payload, err := st.GetLatestExecutionPayloadHeader()
	if err != nil {
		return nil, err
	}

	return beacontypes.GenesisResponse{
		Data: beacontypes.GenesisData{
			GenesisTime:           payload.GetTimestamp().Base10(),
			GenesisValidatorsRoot: genesisRoot,
			GenesisForkVersion:    genesisFork.CurrentVersion.String(),
		},
	}, nil
}
