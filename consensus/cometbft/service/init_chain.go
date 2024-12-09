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
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/encoding/json"
	cmtabci "github.com/cometbft/cometbft/abci/types"
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cometbft/cometbft/node"
	sdk "github.com/cosmos/cosmos-sdk/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/sourcegraph/conc/iter"
)

//nolint:gocognit // its fine.
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

// InitChainer initializes the chain.
func (s *Service[LoggerT]) initChainer(
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

//
// VALIDATION FUNCTIONS
//

// validateDeposits performs validation of the provided deposits.
// It ensures:
// - At least one deposit is present
// - No duplicate public keys
// Returns an error with details if any validation fails.
func validateDeposits(deposits []*types.Deposit) error {
	if len(deposits) == 0 {
		return errors.New("at least one deposit is required")
	}

	seenPubkeys := make(map[string]struct{})

	// In genesis, we have 1:1 mapping between deposits and validators. Hence,
	// we check for duplicate public key.
	for i, deposit := range deposits {
		if deposit == nil {
			return fmt.Errorf("deposit %d is nil", i)
		}

		// Check for duplicate pubkeys
		pubkeyHex := hex.EncodeToString(deposit.Pubkey[:])
		if _, seen := seenPubkeys[pubkeyHex]; seen {
			return fmt.Errorf("duplicate pubkey found in deposit %d", i)
		}
		seenPubkeys[pubkeyHex] = struct{}{}
	}

	return nil
}

// maxExtraDataSize defines the maximum allowed size in bytes for the ExtraData
// field in the execution payload header.
const maxExtraDataSize = 32

// validateExecutionHeader validates the provided execution payload header
// for the genesis block.
func validateExecutionHeader(header *types.ExecutionPayloadHeader) error {
	if header == nil {
		return errors.New("execution payload header cannot be nil")
	}

	// Check block number to be 0
	if header.Number != 0 {
		return errors.New("block number must be 0 for genesis block")
	}

	if header.GasLimit == 0 {
		return errors.New("gas limit cannot be zero")
	}

	if header.GasUsed != 0 {
		return errors.New("gas used must be zero for genesis block")
	}

	if header.BaseFeePerGas == nil {
		return errors.New("base fee per gas cannot be nil")
	}

	// Additional Deneb-specific validations for blob gas
	if header.BlobGasUsed != 0 {
		return errors.New(
			"blob gas used must be zero for genesis block",
		)
	}
	if header.ExcessBlobGas != 0 {
		return errors.New(
			"excess blob gas must be zero for genesis block",
		)
	}

	if header.BlobGasUsed > header.GasLimit {
		return fmt.Errorf("blob gas used (%d) exceeds gas limit (%d)",
			header.BlobGasUsed, header.GasLimit,
		)
	}

	// Validate hash fields are not zero
	zeroHash := common.ExecutionHash{}
	emptyTrieRoot := common.Bytes32(
		common.NewExecutionHashFromHex(
			"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
		))

	// For genesis block (when block number is 0), ParentHash must be zero
	if !bytes.Equal(header.ParentHash[:], zeroHash[:]) {
		return errors.New("parent hash must be zero for genesis block")
	}

	if header.ReceiptsRoot != emptyTrieRoot {
		return errors.New(
			"receipts root must be empty trie root for genesis block",
		)
	}

	if bytes.Equal(header.BlockHash[:], zeroHash[:]) {
		return errors.New("block hash cannot be zero")
	}

	// Validate prevRandao is zero for genesis
	var zeroBytes32 common.Bytes32
	if !bytes.Equal(header.Random[:], zeroBytes32[:]) {
		return errors.New("prevRandao must be zero for genesis block")
	}

	// Fee recipient can be zero in genesis block
	// No need to validate fee recipient for genesis

	// We don't validate LogsBloom as it can legitimately be
	// all zeros in a genesis block or in blocks with no logs

	// Extra data length check (max 32 bytes)
	if len(header.ExtraData) > maxExtraDataSize {
		return fmt.Errorf(
			"extra data too long: got %d bytes, max 32 bytes",
			len(header.ExtraData),
		)
	}

	return nil
}

const expectedHexLength = 8

// isValidForkVersion returns true if the provided fork version is valid.
// A valid fork version must:
// - Start with "0x"
// - Be followed by exactly 8 hexadecimal characters.
func isValidForkVersion(forkVersion common.Version) bool {
	forkVersionStr := forkVersion.String()
	if !strings.HasPrefix(forkVersionStr, "0x") {
		return false
	}

	// Remove "0x" prefix and verify remaining characters
	hexPart := strings.TrimPrefix(forkVersionStr, "0x")

	// Should have exactly 8 characters after 0x prefix
	if len(hexPart) != expectedHexLength {
		return false
	}

	// Verify it's a valid hex number
	_, err := hex.DecodeString(hexPart)
	return err == nil
}
