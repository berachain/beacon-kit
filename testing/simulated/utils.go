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
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/berachain/beacon-kit/chain"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/da/kzg"
	"github.com/berachain/beacon-kit/da/kzg/gokzg"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	gethprimitives "github.com/berachain/beacon-kit/geth-primitives"
	"github.com/berachain/beacon-kit/node-core/components/signer"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	mathpkg "github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/require"
)

// testPkey corresponds to address 0x20f33ce90a13a4b5e7697e3544c3083b8f8a51d4 which is prefunded in genesis
const testPkey = "fffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306"
const blobGasPerTx = 131072

// SharedAccessors holds references to common utilities required in tests.
type SharedAccessors struct {
	Ctx        context.Context
	CancelFunc context.CancelFunc
	HomeDir    string
	TestNode   TestNode

	// ElHandle is a dockertest resource handle that should be closed in teardown.
	ElHandle *dockertest.Resource
}

func GetTestKey(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	// Create a test key - copied from go-ethereum.
	testKey, err := crypto.HexToECDSA(testPkey)
	require.NoError(t, err, "failed to create test key for malicious transaction")
	return testKey
}

// GetBlsSigner returns a new BLSSigner using the configuration files in the provided home directory.
func GetBlsSigner(tempHomeDir string) *signer.BLSSigner {
	privValKeyFile := filepath.Join(tempHomeDir, "config", "priv_validator_key.json")
	privValStateFile := filepath.Join(tempHomeDir, "data", "priv_validator_state.json")
	return signer.NewBLSSigner(privValKeyFile, privValStateFile)
}

// CreateBlockWithTransactions creates a new beacon block with the provided transactions.
// This process requires the engine client as we must simulate to obtain the receipts root
func CreateBlockWithTransactions(
	t *require.Assertions,
	simulationClient *execution.SimulationClient,
	origBlock *ctypes.SignedBeaconBlock,
	blsSigner *signer.BLSSigner,
	chainSpec chain.Spec,
	genesisValidatorsRoot common.Root,
	txs []*gethprimitives.Transaction,
	sidecars []*types.BlobTxSidecar,
	// TODO: To form a valid block we need a valid receiptsRootHash and stateRootHash. This can only be obtained by simulating the block
	// Which can be achieved through the eth_simulateV1 API in a future PR. For now, we hardcode this value.
	receiptsRootHash *gethprimitives.ExecutionHash,
	stateRootHash *gethprimitives.ExecutionHash,
) *ctypes.SignedBeaconBlock {

	simulationClient.Simulate(context.TODO(), )

	// Get the current fork version from the slot.
	forkVersion := chainSpec.ActiveForkVersionForSlot(origBlock.GetMessage().Slot)

	// Extract the existing execution payload and related data.
	payload := origBlock.GetMessage().GetBody().ExecutionPayload
	withdrawals := payload.GetWithdrawals()
	withdrawalsHash := gethprimitives.DeriveSha(withdrawals, gethprimitives.NewStackTrie(nil))
	parentRoot := origBlock.GetMessage().GetParentBlockRoot()

	if receiptsRootHash == nil {
		oldReceiptsRootHash := gethprimitives.ExecutionHash(payload.GetReceiptsRoot())
		receiptsRootHash = &oldReceiptsRootHash
	}
	if stateRootHash == nil {
		oldStateRootHash := gethprimitives.ExecutionHash(payload.GetStateRoot())
		stateRootHash = &oldStateRootHash
	}

	totalTxGasUsed := uint64(0)
	for _, tx := range txs {
		totalTxGasUsed += tx.Gas()
	}

	totalBlobGasUsed := uint64(0)
	for _, sidecar := range sidecars {
		totalBlobGasUsed += uint64(len(sidecar.Blobs) * blobGasPerTx)
	}

	// Construct a new execution block header with the provided transactions.
	executionBlock := gethprimitives.NewBlockWithHeader(
		&gethprimitives.Header{
			ParentHash:       gethprimitives.ExecutionHash(payload.GetParentHash()),
			UncleHash:        gethprimitives.EmptyUncleHash,
			Coinbase:         gethprimitives.ExecutionAddress(payload.GetFeeRecipient()),
			Root:             *stateRootHash,
			TxHash:           gethprimitives.DeriveSha(gethprimitives.Transactions(txs), gethprimitives.NewStackTrie(nil)),
			ReceiptHash:      *receiptsRootHash,
			Bloom:            gethprimitives.LogsBloom(payload.GetLogsBloom()),
			Difficulty:       big.NewInt(0),
			Number:           new(big.Int).SetUint64(payload.GetNumber().Unwrap()),
			GasLimit:         payload.GetGasLimit().Unwrap(),
			GasUsed:          totalTxGasUsed,
			Time:             payload.GetTimestamp().Unwrap(),
			BaseFee:          payload.GetBaseFeePerGas().ToBig(),
			Extra:            payload.GetExtraData(),
			MixDigest:        gethprimitives.ExecutionHash(payload.GetPrevRandao()),
			WithdrawalsHash:  &withdrawalsHash,
			ExcessBlobGas:    payload.GetExcessBlobGas().UnwrapPtr(),
			BlobGasUsed:      &totalBlobGasUsed,
			ParentBeaconRoot: (*gethprimitives.ExecutionHash)(&parentRoot),
		},
	).WithBody(gethprimitives.Body{
		Transactions: txs,
		Uncles:       nil,
		Withdrawals:  *(*gethprimitives.Withdrawals)(unsafe.Pointer(&withdrawals)),
	})

	// Convert the execution block into executable data.
	newExecutionData := gethprimitives.BlockToExecutableData(
		executionBlock,
		nil,
		sidecars,
		nil,
	)

	// Convert the executable data into an ExecutionPayload.
	executionPayload, err := executableDataToExecutionPayload(forkVersion, newExecutionData.ExecutionPayload)
	t.NoError(err, "failed to convert executable data to execution payload")

	// Replace the original payload with the new one.
	origBlock.GetMessage().GetBody().ExecutionPayload = executionPayload

	// Update the block's signature over the new payload.
	newBlock, err := ctypes.NewSignedBeaconBlock(
		origBlock.GetMessage(),
		&ctypes.ForkData{
			CurrentVersion:        chainSpec.ActiveForkVersionForSlot(origBlock.GetMessage().Slot),
			GenesisValidatorsRoot: genesisValidatorsRoot,
		},
		chainSpec,
		blsSigner,
	)
	t.NoError(err, "failed to update signature over the new payload")
	return newBlock
}

func GetProofAndCommitmentsForBlobs(t *require.Assertions, blobs []*eip4844.Blob, verifier kzg.BlobProofVerifier) ([]eip4844.KZGProof, []eip4844.KZGCommitment) {
	if verifier.GetImplementation() != gokzg.Implementation {
		t.Fail("test expects gokzg implementation")
	}
	gokzgVerifier, ok := verifier.(*gokzg.Verifier)
	if !ok {
		t.Fail("verifier is not of type *gokzg.Verifier")
	}
	commitments := make([]eip4844.KZGCommitment, len(blobs))
	proofs := make([]eip4844.KZGProof, len(blobs))
	for i, blob := range blobs {
		ckzgBlob := TransformBlobToGoKZGBlob(blob)
		commitment, err := gokzgVerifier.BlobToKZGCommitment(ckzgBlob, 1)
		t.NoError(err)
		proof, err := gokzgVerifier.ComputeBlobKZGProof(ckzgBlob, commitment, 1)
		t.NoError(err)
		commitments[i] = eip4844.KZGCommitment(commitment)
		proofs[i] = eip4844.KZGProof(proof)
	}
	return proofs, commitments
}

// CreateBeaconBlockWithBlobs TODO
func CreateBeaconBlockWithBlobs(
	t *require.Assertions,
	spec chain.Spec,
	commitments []eip4844.KZGCommitment,
	beaconBlock *ctypes.BeaconBlock,
	blsSigner *signer.BLSSigner,
	genesisValidatorsRoot common.Root,
) *ctypes.SignedBeaconBlock {
	// Replace the original payload with commitment
	beaconBlock.GetBody().BlobKzgCommitments = commitments

	// Update the signature over the new payload.
	newBlock, err := ctypes.NewSignedBeaconBlock(
		beaconBlock,
		&ctypes.ForkData{
			CurrentVersion:        spec.ActiveForkVersionForSlot(beaconBlock.GetSlot()),
			GenesisValidatorsRoot: genesisValidatorsRoot,
		},
		spec,
		blsSigner,
	)
	t.NoError(err, "failed to update signature over the new payload")
	return newBlock
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
