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

package blockchain

import (
	"context"
	"encoding/json"

	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/transition"
)

// ProcessGenesisData processes the genesis state and initializes the beacon state.
func (s *Service[
	_, _, _, _, _, _, GenesisT, _, _,
]) ProcessGenesisData(ctx context.Context, bytes []byte) (transition.ValidatorUpdates, error) {
	genesisData := *new(GenesisT)
	if err := json.Unmarshal(bytes, &genesisData); err != nil {
		s.logger.Error("Failed to unmarshal genesis data", "error", err)
		return nil, err
	}

	// Verify that the genesis deposits root are same between CL and EL.
	genesisDepositsCL := genesisData.GetDeposits()
	genesisDepositsRootEL, err := s.depositContract.GetGenesisDepositsRoot(ctx, 0)
	if err != nil {
		s.logger.Error("Failed to get genesis deposits root from EL", "error", err)
		return nil, err
	}
	if !genesisDepositsCL.HashTreeRoot().Equals(genesisDepositsRootEL) {
		return nil, ErrGenesisDepositsRootMismatch
	}

	// Verify that the genesis deposits count are same between CL and EL.
	genesisDepositsCountEL, err := s.depositContract.GetDepositsCount(ctx, 0)
	if err != nil {
		s.logger.Error("Failed to get genesis deposits count from EL", "error", err)
		return nil, err
	}
	if uint64(len(genesisDepositsCL)) != genesisDepositsCountEL {
		return nil, errors.Wrapf(
			ErrGenesisDepositsCountMismatch,
			"CL genesis deposits count %d mismatch EL genesis deposits count %d",
			len(genesisDepositsCL), genesisDepositsCountEL,
		)
	}

	// Initialize the beacon state from the genesis data.
	validatorUpdates, err := s.stateProcessor.InitializePreminedBeaconStateFromEth1(
		s.storageBackend.StateFromContext(ctx),
		genesisDepositsCL,
		genesisData.GetExecutionPayloadHeader(),
		genesisData.GetForkVersion(),
	)
	if err != nil {
		return nil, err
	}

	// After deposits are validated, store the genesis deposits in the deposit store.
	if err = s.storageBackend.DepositStore().EnqueueDeposits(
		genesisData.GetDeposits(),
	); err != nil {
		return nil, err
	}

	return validatorUpdates, nil
}
