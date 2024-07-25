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
	beacontypes "github.com/berachain/beacon-kit/mod/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/mod/node-api/handlers/types"
	"github.com/berachain/beacon-kit/mod/node-api/handlers/utils"
)

func (h *Handler[_, ContextT, _, _]) GetGenesis(_ ContextT) (any, error) {
	genesisRoot, err := h.backend.GenesisValidatorsRoot(utils.Genesis)
	if err != nil {
		return nil, err
	}
	if len(genesisRoot) == 0 {
		return nil, types.ErrNotFound
	}
	return types.Wrap(beacontypes.GenesisData{
		GenesisTime:           "1590832934", // stub
		GenesisValidatorsRoot: genesisRoot,
		GenesisForkVersion:    "0x00000000", // stub
	}), nil
}
