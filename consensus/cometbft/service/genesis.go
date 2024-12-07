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

package cometbft

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"sort"

	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/encoding/json"
	cmtabci "github.com/cometbft/cometbft/abci/types"
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cometbft/cometbft/node"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
)

func (s *Service[LoggerT]) initChain(
	_ context.Context,
	req *cmtabci.InitChainRequest,
) (*cmtabci.InitChainResponse, error) {
	if req.ChainId != s.chainID {
		return nil, fmt.Errorf(
			"invalid chain-id on InitChain; expected: %s, got: %s",
			s.chainID,
			req.ChainId,
		)
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
		if err = s.sm.CommitMultiStore().
			SetInitialVersion(req.InitialHeight); err != nil {
			return nil, err
		}
	}

	s.finalizeBlockState = s.resetState()

	resValidators, err := s.initChainer(
		s.finalizeBlockState.Context(),
		req.AppStateBytes,
	)
	if err != nil {
		return nil, err
	}

	// check validators
	if len(req.Validators) > 0 {
		if len(req.Validators) != len(resValidators) {
			return nil, fmt.Errorf(
				"len(RequestInitChain.Validators) != len(GenesisValidators) (%d != %d)",
				len(req.Validators),
				len(resValidators),
			)
		}

		sort.Sort(cmtabci.ValidatorUpdates(req.Validators))

		for i := range resValidators {
			if req.Validators[i].Power != resValidators[i].Power {
				return nil, errors.New("mismatched power")
			}
			if !bytes.Equal(
				req.Validators[i].PubKeyBytes, resValidators[i].
					PubKeyBytes) {
				return nil, errors.New("mismatched pubkey bytes")
			}

			if req.Validators[i].PubKeyType !=
				resValidators[i].PubKeyType {
				return nil, errors.New("mismatched pubkey types")
			}
		}
	}

	// NOTE: We don't commit, but FinalizeBlock for block InitialHeight starts
	// from
	// this FinalizeBlockState.
	return &cmtabci.InitChainResponse{
		ConsensusParams: req.ConsensusParams,
		Validators:      resValidators,
		AppHash:         s.sm.CommitMultiStore().LastCommitID().Hash,
	}, nil
}

// DefaultGenesis returns the default genesis state for the application.
func (s *Service[_]) DefaultGenesis() map[string]json.RawMessage {
	// Implement the default genesis state for the application.
	// This should return a map of module names to their respective default
	// genesis states.
	gen := make(map[string]json.RawMessage)
	var err error
	gen["beacon"], err = json.Marshal(types.DefaultGenesisDeneb())
	if err != nil {
		panic(err)
	}
	return gen
}

// ValidateGenesis validates the provided genesis state.
func (s *Service[_]) ValidateGenesis(
	genesisState map[string]json.RawMessage,
) error {
	// Implemented the validation logic for the provided genesis state.
	// This should validate the genesis state for each module in the
	// application.

	// Validate that required modules are present in genesis. Currently,
	// only the beacon module is required.
	beaconGenesisBz, ok := genesisState["beacon"]
	if !ok {
		return errors.New(
			"beacon module genesis state is required but was not found",
		)
	}

	beaconGenesis := &types.Genesis[
		*types.Deposit,
		*types.ExecutionPayloadHeader,
	]{}

	if err := json.Unmarshal(beaconGenesisBz, &beaconGenesis); err != nil {
		return fmt.Errorf(
			"failed to unmarshal beacon genesis state: %w",
			err,
		)
	}

	if !isValidForkVersion(beaconGenesis.GetForkVersion()) {
		return fmt.Errorf("invalid fork version format: %s",
			beaconGenesis.ForkVersion,
		)
	}

	if err := validateDeposits(beaconGenesis.GetDeposits()); err != nil {
		return fmt.Errorf("invalid deposits: %w", err)
	}

	if err := validateExecutionHeader(
		beaconGenesis.GetExecutionPayloadHeader(),
	); err != nil {
		return fmt.Errorf("invalid execution payload header: %w", err)
	}

	return nil
}

// GetGenDocProvider returns a function which returns the genesis doc from the
// genesis file.
func GetGenDocProvider(
	cfg *cmtcfg.Config,
) func() (node.ChecksummedGenesisDoc, error) {
	return func() (node.ChecksummedGenesisDoc, error) {
		appGenesis, err := genutiltypes.AppGenesisFromFile(cfg.GenesisFile())
		if err != nil {
			return node.ChecksummedGenesisDoc{
				Sha256Checksum: []byte{},
			}, err
		}

		gen, err := appGenesis.ToGenesisDoc()
		if err != nil {
			return node.ChecksummedGenesisDoc{
				Sha256Checksum: []byte{},
			}, err
		}
		genbz, err := gen.AppState.MarshalJSON()
		if err != nil {
			return node.ChecksummedGenesisDoc{
				Sha256Checksum: []byte{},
			}, err
		}

		bz, err := json.Marshal(genbz)
		if err != nil {
			return node.ChecksummedGenesisDoc{
				Sha256Checksum: []byte{},
			}, err
		}
		sum := sha256.Sum256(bz)

		return node.ChecksummedGenesisDoc{
			GenesisDoc:     gen,
			Sha256Checksum: sum[:],
		}, nil
	}
}
