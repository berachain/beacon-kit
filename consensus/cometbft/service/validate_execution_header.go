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
	"fmt"

	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/common"
)

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
