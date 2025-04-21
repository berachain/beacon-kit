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

package blockchain

import (
	"context"
	"encoding/json"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/transition"
)

// ProcessGenesisData processes the genesis state and initializes the beacon state.
func (s *Service) ProcessGenesisData(
	ctx context.Context,
	bytes []byte,
) (transition.ValidatorUpdates, *ctypes.BeaconBlockHeader, common.Root, common.Root, error) {
	// Unmarshal the genesis data.
	genesisData := ctypes.Genesis{}
	if err := json.Unmarshal(bytes, &genesisData); err != nil {
		s.logger.Error("Failed to unmarshal genesis data", "error", err)
		return nil, nil, common.Root{}, common.Root{}, err
	}

	// Call the state processor to initialize the "premined" (genesis) beacon state.
	genesisState := s.storageBackend.StateFromContext(ctx)
	validatorUpdates, err := s.stateProcessor.InitializePreminedBeaconStateFromEth1(
		genesisState,
		genesisData.GetDeposits(),
		genesisData.GetExecutionPayloadHeader(),
		genesisData.GetForkVersion(),
	)
	if err != nil {
		return nil, nil, common.Root{}, common.Root{}, err
	}

	// After deposits are validated, store the genesis deposits in the deposit store.
	if err = s.storageBackend.DepositStore().EnqueueDeposits(
		ctx,
		genesisData.GetDeposits(),
	); err != nil {
		return nil, nil, common.Root{}, common.Root{}, err
	}

	// Get the genesis beacon block header from the state.
	genesisHeader, err := genesisState.GetLatestBlockHeader()
	if err != nil {
		return nil, nil, common.Root{}, common.Root{}, err
	}

	genesisBlockRoot, err := genesisState.GetBlockRootAtIndex(0)
	if err != nil {
		return nil, nil, common.Root{}, common.Root{}, err
	}

	// Get the genesis validators root from the state.
	genesisValidatorsRoot, err := genesisState.GetGenesisValidatorsRoot()
	if err != nil {
		return nil, nil, common.Root{}, common.Root{}, err
	}

	return validatorUpdates, genesisHeader, genesisValidatorsRoot, genesisBlockRoot, nil
}
