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
	"fmt"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/transition"
	"github.com/berachain/beacon-kit/primitives/version"
)

// ProcessGenesisData processes the genesis state and initializes the beacon state.
func (s *Service) ProcessGenesisData(
	ctx context.Context,
	bytes []byte,
) (transition.ValidatorUpdates, error) {
	genesisData := ctypes.Genesis{}
	if err := json.Unmarshal(bytes, &genesisData); err != nil {
		s.logger.Error("Failed to unmarshal genesis data", "error", err)
		return nil, err
	}

	// Ensure consistency of the genesis timestamp.
	execPayloadHeader := genesisData.GetExecutionPayloadHeader()
	if s.chainSpec.GenesisTime() != execPayloadHeader.GetTimestamp().Unwrap() {
		return nil, fmt.Errorf(
			"mismatch between chain spec genesis time (%d) and execution payload header time (%d)",
			s.chainSpec.GenesisTime(),
			execPayloadHeader.GetTimestamp().Unwrap(),
		)
	}

	// Ensure consistency of the genesis fork version.
	genesisVersion := genesisData.GetForkVersion()
	if !version.Equals(genesisVersion, s.chainSpec.GenesisForkVersion()) {
		return nil, fmt.Errorf(
			"fork mismatch between CL genesis file version (%s) and chain spec genesis version (%s)",
			genesisVersion, s.chainSpec.GenesisForkVersion(),
		)
	}

	// Initialize the beacon state from the genesis deposits.
	validatorUpdates, err := s.stateProcessor.InitializeBeaconStateFromEth1(
		s.storageBackend.StateFromContext(ctx),
		genesisData.GetDeposits(),
		execPayloadHeader,
		genesisVersion,
	)
	if err != nil {
		return nil, err
	}

	// After deposits are validated, store the genesis deposits in the deposit store.
	if err = s.storageBackend.DepositStore().EnqueueDeposits(
		ctx,
		genesisData.GetDeposits(),
	); err != nil {
		return nil, err
	}

	return validatorUpdates, nil
}
