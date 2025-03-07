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

package cometbft

import (
	"fmt"

	"github.com/berachain/beacon-kit/primitives/encoding/json"
	cmtabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sourcegraph/conc/iter"
)

func (s *Service) initChain(
	req *cmtabci.InitChainRequest,
) (*cmtabci.InitChainResponse, error) {
	if req.ChainId != s.chainID {
		return nil, fmt.Errorf(
			"invalid chain-id on InitChain; expected: %s, got: %s",
			s.chainID,
			req.ChainId,
		)
	}

	// Enforce that request validators is zero. This is because Berachain derives the validators directly from
	// deposits in the genesis file and disregards the validators in genesis file, which is what Comet uses.
	if len(req.Validators) != 0 {
		return nil, fmt.Errorf("expected no validators in initChain request but got %d", len(req.Validators))
	}

	var genesisState map[string]json.RawMessage
	if err := json.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		return nil, fmt.Errorf("failed to unmarshal genesis state: %w", err)
	}
	// Validate the genesis state.
	err := s.ValidateGenesis(genesisState)
	if err != nil {
		return nil, err
	}

	s.logger.Info(
		"InitChain",
		"initialHeight",
		req.InitialHeight,
		"chainID",
		req.ChainId,
	)

	// Set the initial height, which will be used to determine if we are
	// proposing
	// or processing the first block or not.
	s.initialHeight = req.InitialHeight
	if s.initialHeight == 0 {
		s.initialHeight = 1
	}

	// if req.InitialHeight is > 1, then we set the initial version on all
	// stores
	if req.InitialHeight > 1 {
		if err = s.sm.GetCommitMultiStore().
			SetInitialVersion(req.InitialHeight); err != nil {
			return nil, err
		}
	}

	s.finalizeBlockState = s.resetState(s.ctx)

	resValidators, err := s.initChainer(
		s.finalizeBlockState.Context(),
		req.AppStateBytes,
	)
	if err != nil {
		return nil, err
	}

	// NOTE: We don't commit, but FinalizeBlock for block InitialHeight starts
	// from
	// this FinalizeBlockState.
	return &cmtabci.InitChainResponse{
		ConsensusParams: req.ConsensusParams,
		Validators:      resValidators,
		AppHash:         s.sm.GetCommitMultiStore().LastCommitID().Hash,
	}, nil
}

// InitChainer initializes the chain.
func (s *Service) initChainer(
	ctx sdk.Context,
	appStateBytes []byte,
) ([]cmtabci.ValidatorUpdate, error) {
	var genesisState map[string]json.RawMessage
	if err := json.Unmarshal(appStateBytes, &genesisState); err != nil {
		return nil, err
	}

	data := []byte(genesisState["beacon"])
	valUpdates, err := s.Blockchain.ProcessGenesisData(ctx, data)
	if err != nil {
		return nil, err
	}

	return iter.MapErr(
		valUpdates,
		convertValidatorUpdate[cmtabci.ValidatorUpdate],
	)
}
