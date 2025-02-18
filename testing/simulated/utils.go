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
	"testing"
	"unsafe"

	"github.com/berachain/beacon-kit/chain"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	gethprimitives "github.com/berachain/beacon-kit/geth-primitives"
	"github.com/berachain/beacon-kit/node-core/components/signer"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	gethcommon "github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/require"
)

// SharedAccessors holds references to common utilities required in tests
type SharedAccessors struct {
	Ctx        context.Context
	CancelFunc context.CancelFunc
	HomeDir    string
	TestNode   TestNode

	// EL dockertest handles for closing
	ElHandle *dockertest.Resource
}

func GetBlsSigner(tempHomeDir string) *signer.BLSSigner {
	privValKeyFile := filepath.Join(tempHomeDir, "config/priv_validator_key.json")
	privValStateFile := filepath.Join(tempHomeDir, "data/priv_validator_state.json")
	return signer.NewBLSSigner(privValKeyFile, privValStateFile)
}

func CreateInvalidBlock(
	t *testing.T,
	signedBeaconBlock *ctypes.SignedBeaconBlock,
	blsSigner *signer.BLSSigner,
	chainSpec chain.Spec,
	genesisValidatorsRoot common.Root,
) *ctypes.SignedBeaconBlock {
	forkVersion := chainSpec.ActiveForkVersionForSlot(signedBeaconBlock.GetMessage().Slot)
	// Create a transaction from an account that that doesn't have enough balance
	testKey, _ := crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
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
		})
	require.NoError(t, err)

	// Instead of preallocating with length 1 and then appending (which would leave a nil element),
	// initialize the slice with the malicious transaction:
	maliciousTxs := []*gethprimitives.Transaction{maliciousTx}

	payload := signedBeaconBlock.GetMessage().GetBody().ExecutionPayload
	wds := payload.GetWithdrawals()
	withdrawalsHash := gethprimitives.DeriveSha(
		payload.GetWithdrawals(),
		gethprimitives.NewStackTrie(nil),
	)
	parentRoot := signedBeaconBlock.GetMessage().GetParentBlockRoot()
	executionBlock := gethprimitives.NewBlockWithHeader(
		&gethprimitives.Header{
			ParentHash:       gethprimitives.ExecutionHash(payload.GetParentHash()),
			UncleHash:        gethprimitives.EmptyUncleHash,
			Coinbase:         gethprimitives.ExecutionAddress(payload.GetFeeRecipient()),
			Root:             gethprimitives.ExecutionHash(payload.GetStateRoot()),
			TxHash:           gethprimitives.DeriveSha(gethprimitives.Transactions(maliciousTxs), gethprimitives.NewStackTrie(nil)),
			ReceiptHash:      gethprimitives.ExecutionHash(payload.GetReceiptsRoot()),
			Bloom:            gethprimitives.LogsBloom(payload.GetLogsBloom()),
			Difficulty:       big.NewInt(0),
			Number:           new(big.Int).SetUint64(payload.GetNumber().Unwrap()),
			GasLimit:         payload.GetGasLimit().Unwrap(),
			GasUsed:          payload.GetGasUsed().Unwrap(),
			Time:             payload.GetTimestamp().Unwrap(),
			BaseFee:          payload.GetBaseFeePerGas().ToBig(),
			Extra:            payload.GetExtraData(),
			MixDigest:        gethprimitives.ExecutionHash(payload.GetPrevRandao()),
			WithdrawalsHash:  &withdrawalsHash,
			ExcessBlobGas:    payload.GetExcessBlobGas().UnwrapPtr(),
			BlobGasUsed:      payload.GetBlobGasUsed().UnwrapPtr(),
			ParentBeaconRoot: (*gethprimitives.ExecutionHash)(&parentRoot),
		},
	).WithBody(gethprimitives.Body{
		Transactions: maliciousTxs, Uncles: nil, Withdrawals: *(*gethprimitives.Withdrawals)(unsafe.Pointer(&wds)),
	})

	newExecutionData := gethprimitives.BlockToExecutableData(
		executionBlock,
		nil,
		nil,
		nil,
	)

	executionPayload, err := executableDataToExecutionPayload(forkVersion, newExecutionData.ExecutionPayload)
	require.NoError(t, err)

	signedBeaconBlock.GetMessage().GetBody().ExecutionPayload = executionPayload

	// Update the signature over the new payload
	maliciousBlock, err := ctypes.NewSignedBeaconBlock(signedBeaconBlock.GetMessage(), &ctypes.ForkData{
		CurrentVersion:        chainSpec.ActiveForkVersionForSlot(signedBeaconBlock.GetMessage().Slot),
		GenesisValidatorsRoot: genesisValidatorsRoot,
	}, chainSpec, blsSigner)
	require.NoError(t, err)
	return maliciousBlock
}

// executableDataToExecutionPayload converts the eth executable data type to the beacon execution payload.
// Adapted from executableDataToExecutionPayloadHeader.
func executableDataToExecutionPayload(
	forkVersion common.Version,
	data *gethprimitives.ExecutableData,
) (*ctypes.ExecutionPayload, error) {
	var executionPayload *ctypes.ExecutionPayload
	if version.IsBefore(forkVersion, version.Deneb1()) {
		withdrawals := make(
			engineprimitives.Withdrawals,
			len(data.Withdrawals),
		)
		for i, withdrawal := range data.Withdrawals {
			// #nosec:G103 // primitives.Withdrawals are data.Withdrawals with hard types
			withdrawals[i] = (*engineprimitives.Withdrawal)(
				unsafe.Pointer(withdrawal),
			)
		}

		if len(data.ExtraData) > constants.ExtraDataLength {
			data.ExtraData = data.ExtraData[:constants.ExtraDataLength]
		}

		var blobGasUsed uint64
		if data.BlobGasUsed != nil {
			blobGasUsed = *data.BlobGasUsed
		}

		var excessBlobGas uint64
		if data.ExcessBlobGas != nil {
			excessBlobGas = *data.ExcessBlobGas
		}

		baseFeePerGas, err := math.NewU256FromBigInt(data.BaseFeePerGas)
		if err != nil {
			return nil, fmt.Errorf("failed baseFeePerGas conversion: %w", err)
		}

		executionPayload = &ctypes.ExecutionPayload{
			ParentHash:    common.ExecutionHash(data.ParentHash),
			FeeRecipient:  common.ExecutionAddress(data.FeeRecipient),
			StateRoot:     common.Bytes32(data.StateRoot),
			ReceiptsRoot:  common.Bytes32(data.ReceiptsRoot),
			LogsBloom:     [256]byte(data.LogsBloom),
			Random:        common.Bytes32(data.Random),
			Number:        math.U64(data.Number),
			GasLimit:      math.U64(data.GasLimit),
			GasUsed:       math.U64(data.GasUsed),
			Timestamp:     math.U64(data.Timestamp),
			Withdrawals:   withdrawals,
			ExtraData:     data.ExtraData,
			BaseFeePerGas: baseFeePerGas,
			BlockHash:     common.ExecutionHash(data.BlockHash),
			Transactions:  data.Transactions,
			BlobGasUsed:   math.U64(blobGasUsed),
			ExcessBlobGas: math.U64(excessBlobGas),
			EpVersion:     forkVersion,
		}
	} else {
		return nil, ctypes.ErrForkVersionNotSupported
	}
	return executionPayload, nil
}
