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
	"path/filepath"
	"testing"
	"time"
	"unsafe"

	"github.com/berachain/beacon-kit/beacon/blockchain"
	"github.com/berachain/beacon-kit/chain"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	gethprimitives "github.com/berachain/beacon-kit/geth-primitives"
	"github.com/berachain/beacon-kit/node-core/components/signer"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/transition"
	"github.com/berachain/beacon-kit/state-transition/core"
	"github.com/berachain/beacon-kit/testing/simulated/execution"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

// testPkey corresponds to address 0x20f33ce90a13a4b5e7697e3544c3083b8f8a51d4 which is prefunded in genesis
const testPkey = "fffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306"

// SharedAccessors holds references to common utilities required in tests.
type SharedAccessors struct {
	CtxApp         context.Context
	CtxAppCancelFn context.CancelFunc
	CtxComet       context.Context
	HomeDir    string
	TestNode   TestNode

	// ElHandle is a dockertest resource handle that should be closed in teardown.
	ElHandle *execution.Resource
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

func DefaultSimulationInput(t *testing.T, chainSpec chain.Spec, origBlock *ctypes.BeaconBlock, txs []*gethprimitives.Transaction) *execution.SimOpts {
	t.Helper()
	overrideTime := hexutil.Uint64(origBlock.GetTimestamp().Unwrap())
	overrideGasLimit := hexutil.Uint64(30000000)
	overrideFeeRecipient := origBlock.GetBody().GetExecutionPayload().GetFeeRecipient()
	overridePrevRandao := gethcommon.Hash(origBlock.GetBody().GetExecutionPayload().GetPrevRandao())
	overrideBaseFeePerGas := origBlock.GetBody().GetExecutionPayload().GetBaseFeePerGas().ToBig()
	overrideBeaconRoot := gethcommon.HexToHash(origBlock.GetParentBlockRoot().Hex())
	origWithdrawls := origBlock.GetBody().GetExecutionPayload().GetWithdrawals()
	overrideWithdrawals := *(*gethtypes.Withdrawals)(unsafe.Pointer(&origWithdrawls))

	calls, err := execution.TxsToTransactionArgs(chainSpec.DepositEth1ChainID(), txs)
	require.NoError(t, err)
	simulationInput := &execution.SimOpts{
		BlockStateCalls: []*execution.SimBlock{
			{
				Calls: calls,
				BlockOverrides: &execution.BlockOverrides{
					Time:          &overrideTime,
					GasLimit:      &overrideGasLimit,
					FeeRecipient:  (*gethcommon.Address)(&overrideFeeRecipient),
					PrevRandao:    &overridePrevRandao,
					BaseFeePerGas: (*hexutil.Big)(overrideBaseFeePerGas),
					BeaconRoot:    &overrideBeaconRoot,
					Withdrawals:   overrideWithdrawals,
					// TODO: Do we need to override blob base fee?
				},
			},
		},
		Validation:     true,
		TraceTransfers: false,
	}
	return simulationInput
}

// ComputeAndSetInvalidExecutionBlock transforms the current execution payload of latestBlock
// into a new payload (using the invalid transformation) and updates latestBlock with it.
// This will make sure all the fields validated by the CL, i.e. Execution Block Hash, are valid, but does not set
// correct values for fields like the Execution Block StateRoot and ReceiptsRoot as that requires simulation and
// is not validated in the CL.
func ComputeAndSetInvalidExecutionBlock(
	t *testing.T,
	latestBlock *ctypes.BeaconBlock,
	chainSpec chain.Spec,
	txs []*gethprimitives.Transaction,
) *ctypes.BeaconBlock {
	t.Helper()
	forkVersion := chainSpec.ActiveForkVersionForSlot(latestBlock.GetSlot())
	_, sidecars := splitTxs(txs)
	// Use the current execution payload (e.g. for an invalid block, no simulation is done).
	executionPayload := latestBlock.GetBody().GetExecutionPayload()
	// Transform the payload into a Geth block.
	txsBytesArray := make([][]byte, len(txs))
	for i, tx := range txs {
		txBytes, err := tx.MarshalBinary()
		require.NoError(t, err)
		txsBytesArray[i] = txBytes
	}
	executionPayload.Transactions = txsBytesArray
	parentBlockRoot := latestBlock.GetParentBlockRoot()
	execBlock, _, err := ctypes.MakeEthBlock(executionPayload, &parentBlockRoot)
	require.NoError(t, err)
	return setExecutionPayload(t, latestBlock, forkVersion, execBlock, sidecars)
}

// ComputeAndSetValidExecutionBlock simulates a new execution payload based on the provided transactions,
// transforms the simulated block into a Geth-style execution block, and updates the given beacon block
// with the new execution payload. This will correctly set the Execution Block State and Receipts Root using simulation.
// Note: The returned block's state root is not finalized and must be updated via a state transition (see ComputeAndSetStateRoot).
func ComputeAndSetValidExecutionBlock(
	t *testing.T,
	latestBlock *ctypes.BeaconBlock,
	simClient *execution.SimulationClient,
	chainSpec chain.Spec,
	txs []*gethprimitives.Transaction,
) *ctypes.BeaconBlock {
	// Run simulation to get a simulated block.
	baseHeight := int64(latestBlock.GetSlot().Unwrap()) - 1
	simInput := DefaultSimulationInput(t, chainSpec, latestBlock, txs)
	simulatedBlocks, err := simClient.Simulate(context.TODO(), baseHeight, simInput)
	require.NoError(t, err)
	require.Len(t, simulatedBlocks, 1)
	simBlock := simulatedBlocks[0]

	forkVersion := chainSpec.ActiveForkVersionForSlot(latestBlock.GetSlot())
	txsNoSidecar, sidecars := splitTxs(txs)
	origParent := latestBlock.GetParentBlockRoot()

	// Transform the simulated block into a Geth block.
	execBlock := transformSimulatedBlockToGethBlock(simBlock, txsNoSidecar, origParent)
	return setExecutionPayload(t, latestBlock, forkVersion, execBlock, sidecars)
}

// ComputeAndSetStateRoot applies a state transition to the given beacon block.
// It creates a copy of the current state (from the provided storage backend and query context),
// constructs a transition context using the consensus time and proposer address,
// runs the state transition, and then updates the block’s state root based on the new state.
// Returns the updated block or an error.
// TODO: Can we use a mocked execution client for the StateProcessor to avoid doing an unnecessary NewPayload?
func ComputeAndSetStateRoot(
	queryCtx context.Context,
	consensusTime time.Time,
	proposerAddress []byte,
	stateProcessor *core.StateProcessor,
	storageBackend blockchain.StorageBackend,
	block *ctypes.BeaconBlock,
) (*ctypes.BeaconBlock, error) {

	// Copy the current state from the storage backend.
	stateDBCopy := storageBackend.StateFromContext(queryCtx).Copy(queryCtx)

	// Create a transition context with the provided consensus time and proposer address.
	txCtx := transition.NewTransitionCtx(
		queryCtx,
		math.U64(consensusTime.Unix()),
		proposerAddress,
	).WithVerifyPayload(false).
		WithVerifyRandao(false).
		WithVerifyResult(false).
		WithMeterGas(false)

	// Run the state transition.
	_, err := stateProcessor.Transition(txCtx, stateDBCopy, block)
	if err != nil {
		return nil, fmt.Errorf("state transition failed: %w", err)
	}

	// Compute the new state root from the updated state.
	newStateRoot := stateDBCopy.HashTreeRoot()
	block.SetStateRoot(newStateRoot)
	return block, nil
}

// setExecutionPayload converts the given Geth-style block into executable data,
// converts that into an ExecutionPayload using the given fork version, and then
// sets that payload into latestBlock. It returns the updated block.
func setExecutionPayload(
	t *testing.T,
	latestBlock *ctypes.BeaconBlock,
	forkVersion common.Version,
	execBlock *gethtypes.Block,
	sidecars []*gethtypes.BlobTxSidecar, // adjust type as needed
) *ctypes.BeaconBlock {
	// Convert the Geth block into ExecutableData.
	execData := gethprimitives.BlockToExecutableData(execBlock, nil, sidecars, nil)
	// Convert the ExecutableData into our internal ExecutionPayload type.
	execPayload, err := transformExecutableDataToExecutionPayload(forkVersion, execData.ExecutionPayload)
	require.NoError(t, err, "failed to convert executable data")
	// Update the beacon block with the new execution payload.
	latestBlock.GetBody().SetExecutionPayload(execPayload)
	return latestBlock
}
