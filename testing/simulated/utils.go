//go:build simulated

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

package simulated

import (
	"context"
	"fmt"
	"math/big"
	"path/filepath"
	"unsafe"

	"github.com/berachain/beacon-kit/chain"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	gethprimitives "github.com/berachain/beacon-kit/geth-primitives"
	"github.com/berachain/beacon-kit/node-core/components/signer"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	mathpkg "github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/berachain/beacon-kit/testing/simulated/execution"
	gethcommon "github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

// SharedAccessors holds references to common utilities required in tests.
type SharedAccessors struct {
	Ctx        context.Context
	CancelFunc context.CancelFunc
	HomeDir    string
	TestNode   TestNode

	// ElHandle is a dockertest resource handle that should be closed in teardown.
	ElHandle *execution.Resource
}

// GetBlsSigner returns a new BLSSigner using the configuration files in the provided home directory.
func GetBlsSigner(tempHomeDir string) *signer.BLSSigner {
	privValKeyFile := filepath.Join(tempHomeDir, "config", "priv_validator_key.json")
	privValStateFile := filepath.Join(tempHomeDir, "data", "priv_validator_state.json")
	return signer.NewBLSSigner(privValKeyFile, privValStateFile)
}

// CreateInvalidBlock creates a malicious beacon block by injecting an invalid transaction
// into the execution payload. The invalidity stems from the transaction coming from an account
// with fee below base fee.
func CreateInvalidBlock(
	t *require.Assertions,
	signedBeaconBlock *ctypes.SignedBeaconBlock,
	blsSigner *signer.BLSSigner,
	chainSpec chain.Spec,
	genesisValidatorsRoot common.Root,
) *ctypes.SignedBeaconBlock {
	// Get the current fork version from the slot.
	forkVersion := chainSpec.ActiveForkVersionForSlot(signedBeaconBlock.GetMessage().Slot)

	// Create a test key - copied from go-ethereum.
	testKey, err := crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	t.NoError(err, "failed to create test key for malicious transaction")

	// Sign a malicious transaction that is expected to fail.
	maliciousTx, err := gethtypes.SignNewTx(
		testKey,
		gethtypes.NewCancunSigner(big.NewInt(int64(chainSpec.DepositEth1ChainID()))),
		&gethtypes.DynamicFeeTx{
			Nonce:     0,
			To:        &gethcommon.Address{1},
			Value:     big.NewInt(100000000000),
			Gas:       100,
			GasTipCap: big.NewInt(10000000),
			GasFeeCap: big.NewInt(10000000),
			Data:      []byte{},
		},
	)
	t.NoError(err, "failed to sign malicious transaction")

	payload := signedBeaconBlock.GetMessage().GetBody().ExecutionPayload
	parentRoot := signedBeaconBlock.GetMessage().GetParentBlockRoot()

	// Update the ExecutionPayload with the malicious transaction
	maliciousTxBytes, err := maliciousTx.MarshalBinary()
	t.NoError(err, "failed to marshal malicious transaction")
	payload.Transactions = [][]byte{maliciousTxBytes}

	executionBlock, _, err := ctypes.MakeEthBlock(payload, &parentRoot)
	t.NoError(err, "failed to make execution block")

	// Convert the execution block into executable data.
	newExecutionData := gethprimitives.BlockToExecutableData(
		executionBlock,
		nil,
		nil,
		nil,
	)

	// Convert the executable data into an ExecutionPayload.
	executionPayload, err := executableDataToExecutionPayload(forkVersion, newExecutionData.ExecutionPayload)
	t.NoError(err, "failed to convert executable data to execution payload")

	// Replace the original payload with the malicious payload.
	signedBeaconBlock.GetMessage().GetBody().ExecutionPayload = executionPayload

	// Update the signature over the new payload.
	maliciousBlock, err := ctypes.NewSignedBeaconBlock(
		signedBeaconBlock.GetMessage(),
		&ctypes.ForkData{
			CurrentVersion:        chainSpec.ActiveForkVersionForSlot(signedBeaconBlock.GetMessage().Slot),
			GenesisValidatorsRoot: genesisValidatorsRoot,
		},
		chainSpec,
		blsSigner,
	)
	t.NoError(err, "failed to update signature over the new payload")
	return maliciousBlock
}

// executableDataToExecutionPayload converts Ethereum executable data to a beacon execution payload.
// It supports fork versions before Deneb1 and returns an error if the fork version is not supported.
func executableDataToExecutionPayload(
	forkVersion common.Version,
	data *gethprimitives.ExecutableData,
) (*ctypes.ExecutionPayload, error) {
	// Only support fork versions before Deneb1.
	if version.IsBefore(forkVersion, version.Deneb1()) {
		// Convert withdrawals from gethprimitives to engineprimitives.
		withdrawals := make(engineprimitives.Withdrawals, len(data.Withdrawals))
		for i, withdrawal := range data.Withdrawals {
			// #nosec:G103 -- safe conversion assuming the underlying types are compatible.
			withdrawals[i] = (*engineprimitives.Withdrawal)(unsafe.Pointer(withdrawal))
		}

		// Truncate ExtraData if it exceeds the allowed length.
		if len(data.ExtraData) > constants.ExtraDataLength {
			data.ExtraData = data.ExtraData[:constants.ExtraDataLength]
		}

		// Dereference optional fields safely.
		var blobGasUsed, excessBlobGas uint64
		if data.BlobGasUsed != nil {
			blobGasUsed = *data.BlobGasUsed
		}
		if data.ExcessBlobGas != nil {
			excessBlobGas = *data.ExcessBlobGas
		}

		// Convert BaseFeePerGas into a U256 value.
		baseFeePerGas, err := mathpkg.NewU256FromBigInt(data.BaseFeePerGas)
		if err != nil {
			return nil, fmt.Errorf("failed baseFeePerGas conversion: %w", err)
		}

		executionPayload := &ctypes.ExecutionPayload{
			ParentHash:    common.ExecutionHash(data.ParentHash),
			FeeRecipient:  common.ExecutionAddress(data.FeeRecipient),
			StateRoot:     common.Bytes32(data.StateRoot),
			ReceiptsRoot:  common.Bytes32(data.ReceiptsRoot),
			LogsBloom:     [256]byte(data.LogsBloom),
			Random:        common.Bytes32(data.Random),
			Number:        mathpkg.U64(data.Number),
			GasLimit:      mathpkg.U64(data.GasLimit),
			GasUsed:       mathpkg.U64(data.GasUsed),
			Timestamp:     mathpkg.U64(data.Timestamp),
			Withdrawals:   withdrawals,
			ExtraData:     data.ExtraData,
			BaseFeePerGas: baseFeePerGas,
			BlockHash:     common.ExecutionHash(data.BlockHash),
			Transactions:  data.Transactions,
			BlobGasUsed:   mathpkg.U64(blobGasUsed),
			ExcessBlobGas: mathpkg.U64(excessBlobGas),
			EpVersion:     forkVersion,
		}
		return executionPayload, nil
	}
	return nil, ctypes.ErrForkVersionNotSupported
}
