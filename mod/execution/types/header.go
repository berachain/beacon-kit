// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package types

import (
	"unsafe"

	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/uint256"
	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/core/types"
)

// IsValidVersionAndBlockHashes checks if the version and block hashes are
// valid.
// As per the Ethereum 2.0 specification:
// https://github.com/ethereum/consensus-specs/blob/v1.4.0-beta.2/specs/deneb/beacon-chain.md#is_valid_block_hash
// https://github.com/ethereum/consensus-specs/blob/v1.4.0-beta.2/specs/deneb/beacon-chain.md#is_valid_versioned_hashes
//
//nolint:lll
func IsValidVersionAndBlockHashes(
	payload *ExecutableDataDeneb,
	parentBeaconBlockRoot *primitives.Root,
	versionedHashes []primitives.ExecutionHash,
) error {
	data := engine.ExecutableData{
		ParentHash:    payload.ParentHash,
		FeeRecipient:  payload.FeeRecipient,
		StateRoot:     payload.StateRoot,
		ReceiptsRoot:  payload.ReceiptsRoot,
		LogsBloom:     payload.LogsBloom,
		Random:        payload.Random,
		Number:        payload.Number,
		GasLimit:      payload.GasLimit,
		GasUsed:       payload.GasUsed,
		Timestamp:     payload.Timestamp,
		ExtraData:     payload.ExtraData,
		BaseFeePerGas: uint256.LittleFromBigEndian(payload.BaseFeePerGas).Big(),
		BlockHash:     payload.BlockHash,
		Transactions:  payload.Transactions,
		Withdrawals:   *(*[]*types.Withdrawal)(unsafe.Pointer(&payload.Withdrawals)),
		BlobGasUsed:   &payload.BlobGasUsed,
		ExcessBlobGas: &payload.ExcessBlobGas,
	}
	_, err := engine.ExecutableDataToBlock(
		data,
		versionedHashes,
		(*primitives.ExecutionHash)(parentBeaconBlockRoot),
	)
	return err
}
