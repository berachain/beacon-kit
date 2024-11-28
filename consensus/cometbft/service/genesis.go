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
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/json"
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cometbft/cometbft/node"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
)

const maxExtraDataSize = 32

// BeaconGenesisState represents the structure of the
// beacon module's genesis state.
//
//nolint:lll // JSON tags are long.
type BeaconGenesisState struct {
	ForkVersion            string                       `json:"fork_version"`
	Deposits               []types.Deposit              `json:"deposits"`
	ExecutionPayloadHeader types.ExecutionPayloadHeader `json:"execution_payload_header"`
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

// isZeroBytes returns true if the provided byte slice is all zeros.
func isZeroBytes(b []byte) bool {
	for _, byteValue := range b {
		if byteValue != 0 {
			return false
		}
	}
	return true
}

// isValidForkVersion returns true if the provided fork version is valid.
// Validate fork version format (should be 0x followed by 8 hex characters).
func isValidForkVersion(forkVersion string) bool {
	if !strings.HasPrefix(forkVersion, "0x") || len(forkVersion) != 10 {
		return false
	}
	_, err := hex.DecodeString(forkVersion[2:])
	return err == nil

}

// validateDeposits validates the provided deposits.
func validateDeposits(deposits []types.Deposit) error {
	if len(deposits) == 0 {
		return errors.New("at least one deposit is required")
	}

	seenPubkeys := make(map[string]bool)

	for i, deposit := range deposits {
		depositIndex := deposit.GetIndex()
		//#nosec:G701 // realistically fine in practice.
		// Validate index matches position
		if depositIndex.Unwrap() != uint64(i) {
			return fmt.Errorf(
				"deposit index %d does not match position %d",
				depositIndex,
				i,
			)
		}

		if isZeroBytes(deposit.Pubkey[:]) {
			return fmt.Errorf("deposit %d has zero public key", i)
		}
		// Check for duplicate pubkeys
		pubkeyHex := hex.EncodeToString(deposit.Pubkey[:])
		if seenPubkeys[pubkeyHex] {
			return fmt.Errorf("duplicate pubkey found in deposit %d", i)
		}
		seenPubkeys[pubkeyHex] = true

		if isZeroBytes(deposit.Credentials[:]) {
			return fmt.Errorf(
				"invalid withdrawal credentials length for deposit %d",
				i,
			)
		}

		if deposit.Amount == 0 {
			return fmt.Errorf("deposit %d has zero amount", i)
		}

		if isZeroBytes(deposit.Signature[:]) {
			return fmt.Errorf("invalid signature length for deposit %d", i)
		}
	}

	return nil
}

// validateExecutionHeader validates the provided execution payload header.
func validateExecutionHeader(header types.ExecutionPayloadHeader) error {
	// Validate hash fields are not zero
	zeroHash := common.ExecutionHash{}
	// For genesis block (when block number is 0), ParentHash must be zero
	if !bytes.Equal(header.ParentHash[:], zeroHash[:]) {
		return errors.New("parent hash must be zero for genesis block")
	}

	if bytes.Equal(header.StateRoot[:], zeroHash[:]) {
		return errors.New("state root cannot be zero")
	}
	if bytes.Equal(header.ReceiptsRoot[:], zeroHash[:]) {
		return errors.New("receipts root cannot be zero")
	}
	if bytes.Equal(header.BlockHash[:], zeroHash[:]) {
		return errors.New("block hash cannot be zero")
	}
	if bytes.Equal(header.TransactionsRoot[:], zeroHash[:]) {
		return errors.New("transactions root cannot be zero")
	}

	// Check block number to be 0
	if header.Number != 0 {
		return errors.New("block number must be 0 for genesis block")
	}

	// Validate prevRandao is zero for genesis
	if !bytes.Equal(header.Random[:], zeroHash[:]) {
		return errors.New("prevRandao must be zero for genesis block")
	}

	// Fee recipient can be zero in genesis block
	// No need to validate fee recipient for genesis

	// We don't validate LogsBloom as it can legitimately be
	// all zeros in a genesis block or in blocks with no logs

	// Validate numeric fields
	if header.GasLimit == 0 {
		return errors.New("gas limit cannot be zero")
	}

	// Extra data length check (max 32 bytes)
	if len(header.ExtraData) > maxExtraDataSize {
		return fmt.Errorf(
			"extra data too long: got %d bytes, max 32 bytes",
			len(header.ExtraData),
		)
	}

	// Validate base fee per gas
	if header.BaseFeePerGas == nil {
		return errors.New("base fee per gas cannot be nil")
	}

	// Additional Deneb-specific validations for blob gas
	if header.BlobGasUsed > header.GetGasLimit() {
		return errors.New("blob gas used exceeds gas limit")
	}

	return nil
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

	var beaconGenesis BeaconGenesisState
	if err := json.Unmarshal(beaconGenesisBz, &beaconGenesis); err != nil {
		return fmt.Errorf(
			"failed to unmarshal beacon genesis state: %w",
			err,
		)
	}

	if !isValidForkVersion(beaconGenesis.ForkVersion) {
		return fmt.Errorf("invalid fork version format: %s",
			beaconGenesis.ForkVersion,
		)
	}

	if err := validateDeposits(beaconGenesis.Deposits); err != nil {
		return fmt.Errorf("invalid deposits: %w", err)
	}

	if err := validateExecutionHeader(
		beaconGenesis.ExecutionPayloadHeader,
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
