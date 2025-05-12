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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package e2e_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"

	beaconapi "github.com/attestantio/go-eth2-client/api"
	"github.com/berachain/beacon-kit/execution/requests/eip7002"
	"github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/node-api/handlers/utils"
	"github.com/berachain/beacon-kit/primitives/crypto"
	beaconmath "github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/berachain/beacon-kit/testing/e2e/suite"
	e2etypes "github.com/berachain/beacon-kit/testing/e2e/suite/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethcore "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	// DefaultWithdrawalAmount is the default withdrawal amount in Gwei for testing
	DefaultWithdrawalAmount = 1_000_000_000 // 1 BERA (10^9 Gwei)

	// WithdrawalTxGasLimit is the gas limit for withdrawal transactions
	WithdrawalTxGasLimit = 500_000

	// BlocksToWaitAfterWithdrawal is the number of blocks to wait after a withdrawal
	BlocksToWaitAfterWithdrawal = 3
)

// rpcWrapper wraps an rpc.Client and implements the rpcClient interface required by EIP7002.
// EIP7002 requires a Call method that returns an error.
type rpcWrapper struct {
	*rpc.Client
}

// Call implements the rpcClient interface.
func (r *rpcWrapper) Call(ctx context.Context, target any, method string, params ...any) error {
	return r.Client.CallContext(ctx, target, method, params...)
}

// getPendingPartialWithdrawals calls the beacon node's /eth/v1/beacon/states/{state_id}/pending_partial_withdrawals endpoint
// and returns the list of pending withdrawals data if any
func (s *BeaconKitE2ESuite) getPendingPartialWithdrawals(stateID string) ([]types.PendingPartialWithdrawalData, error) {
	client := s.initHTTPBeaconTest()

	url := fmt.Sprintf("/eth/v1/beacon/states/%s/pending_partial_withdrawals", stateID)

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending partial withdrawals: %w", err)
	}
	if resp == nil {
		return nil, errors.New("received nil response")
	}

	defer resp.Body.Close()

	var response types.PendingPartialWithdrawalsResponse
	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Validate response fields
	s.Require().Equal(response.Version, version.Name(version.Electra()))
	s.Require().False(response.ExecutionOptimistic)
	s.Require().True(response.Finalized)

	withdrawals, ok := response.Data.([]types.PendingPartialWithdrawalData)
	if !ok {
		// Only if direct type assertion fails, fall back to JSON remarshal
		dataBytes, errInMarshal := json.Marshal(response.Data)
		if errInMarshal != nil {
			return nil, fmt.Errorf("failed to marshal data: %w", errInMarshal)
		}

		if errInUnmarshal := json.Unmarshal(dataBytes, &withdrawals); errInUnmarshal != nil {
			return nil, fmt.Errorf("failed to unmarshal withdrawals: %w", errInUnmarshal)
		}
	}

	return withdrawals, nil
}

// findValidatorWithExecutionCredentials finds a validator with execution credentials
func (s *BeaconKitE2ESuite) findValidatorWithExecutionCredentials(client *e2etypes.ConsensusClient) (string, crypto.BLSPubkey, error) {
	// Get the validators to identify one with execution credentials (0x01)
	validatorsResp, err := client.Validators(
		s.Ctx(),
		&beaconapi.ValidatorsOpts{
			State: utils.StateIDHead,
		},
	)
	if err != nil {
		return "", crypto.BLSPubkey{}, fmt.Errorf("failed to get validators: %w", err)
	}
	if len(validatorsResp.Data) == 0 {
		return "", crypto.BLSPubkey{}, errors.New("no validators found")
	}

	// Find a validator with execution withdrawal credentials (starting with 0x01)
	var validatorIndex string
	var blsPubkey crypto.BLSPubkey
	for _, validator := range validatorsResp.Data {
		credentials := validator.Validator.WithdrawalCredentials
		if len(credentials) >= 2 && credentials[0] == 0x01 {
			validatorIndex = fmt.Sprintf("%d", validator.Index)
			// Convert the phase0.BLSPubKey to our crypto.BLSPubkey type
			copy(blsPubkey[:], validator.Validator.PublicKey[:])
			return validatorIndex, blsPubkey, nil
		}
	}

	return "", crypto.BLSPubkey{}, errors.New("no validator with execution withdrawal credentials found")
}

// TestSubmitPartialWithdrawalTransaction tests submitting a partial withdrawal transaction via the withdrawal contract
func (s *BeaconKitE2ESuite) TestSubmitPartialWithdrawalTransaction() {
	// Use timeout context to better control the test
	ctx, cancel := context.WithTimeout(s.Ctx(), suite.DefaultE2ETestTimeout)
	defer cancel()

	// Initialize test client
	client := s.initBeaconTest()

	// Find a validator with execution credentials
	validatorIndex, blsPubkey, err := s.findValidatorWithExecutionCredentials(client)
	s.Require().NoError(err)
	s.T().Logf("Found validator with index %s for withdrawal test", validatorIndex)

	// Set withdrawal amount
	withdrawalAmount := beaconmath.Gwei(DefaultWithdrawalAmount)

	// Create an rpc client using the load balancer URL
	rpcClient, err := rpc.Dial(s.JSONRPCBalancer().URL())
	s.Require().NoError(err)
	defer rpcClient.Close()

	rpcWrapper := &rpcWrapper{Client: rpcClient}

	// Get current block number before withdrawal
	blkNum, err := s.JSONRPCBalancer().BlockNumber(ctx)
	s.Require().NoError(err)
	s.T().Logf("Block number before withdrawal: %d", blkNum)

	// Check for pending partial withdrawals before submitting the transaction
	pendingWithdrawalsBefore, err := s.getPendingPartialWithdrawals(utils.StateIDHead)
	s.Require().NoError(err)
	s.Require().Len(pendingWithdrawalsBefore, 0, "Expected no pending withdrawals initially")

	// Get the withdrawal fee
	fee, err := eip7002.GetWithdrawalFee(ctx, rpcWrapper)
	s.Require().NoError(err)
	s.T().Logf("Withdrawal fee: %s wei", fee.String())

	// Create the withdrawal transaction data
	withdrawalTxData, err := eip7002.CreateWithdrawalRequestData(blsPubkey, withdrawalAmount)
	s.Require().NoError(err)

	// Use a pre-loaded key that has funds
	privateKey, err := ethcrypto.HexToECDSA("fffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306")
	s.Require().NoError(err)

	chainID, err := s.JSONRPCBalancer().ChainID(ctx)
	s.Require().NoError(err)
	signer := gethcore.NewPragueSigner(chainID)

	// Get the sender's nonce
	var nonce hexutil.Uint64
	err = rpcClient.CallContext(ctx, &nonce, "eth_getTransactionCount",
		common.HexToAddress("0x20f33ce90a13a4b5e7697e3544c3083b8f8a51d4"), "latest",
	)
	s.Require().NoError(err)

	// Create and sign the transaction
	tx := gethcore.MustSignNewTx(privateKey, signer, &gethcore.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     uint64(nonce),
		To:        &params.WithdrawalQueueAddress,
		Gas:       WithdrawalTxGasLimit,
		GasFeeCap: big.NewInt(1000000000),
		GasTipCap: big.NewInt(1000000000),
		Value:     fee,
		Data:      withdrawalTxData,
	})

	// Serialize the transaction
	txBytes, err := tx.MarshalBinary()
	s.Require().NoError(err)

	// Send the transaction
	var txHash common.Hash
	err = rpcClient.CallContext(ctx, &txHash, "eth_sendRawTransaction", hexutil.Encode(txBytes))
	s.Require().NoError(err)
	s.T().Logf("Withdrawal transaction submitted: %s", txHash.Hex())

	// Wait for blocks to be mined after submitting the transaction
	err = s.WaitForNBlockNumbers(BlocksToWaitAfterWithdrawal)
	s.Require().NoError(err)

	// Get the transaction receipt
	var receipt map[string]interface{}
	err = rpcClient.CallContext(ctx, &receipt, "eth_getTransactionReceipt", txHash.Hex())
	s.Require().NoError(err)
	s.Require().NotNil(receipt, "Transaction receipt should not be nil")

	// Get block number where the withdrawal transaction was included
	blockNumStr, ok := receipt["blockNumber"].(string)
	s.Require().True(ok, "Block number should be a string")
	blockNum, err := hexutil.DecodeUint64(blockNumStr)
	s.Require().NoError(err)
	s.T().Logf("Withdrawal transaction included in block: %d", blockNum)

	// Check for pending partial withdrawals after submitting the transaction
	pendingWithdrawalsAfter, err := s.getPendingPartialWithdrawals(utils.StateIDHead)
	s.Require().NoError(err)
	s.Require().Len(pendingWithdrawalsAfter, 1, "Expected one pending withdrawal after transaction")

	// Verify the withdrawal details
	s.Require().Equal(validatorIndex, strconv.FormatUint(pendingWithdrawalsAfter[0].ValidatorIndex, 10),
		"Validator index mismatch in pending withdrawal")
	s.Require().Equal(uint64(withdrawalAmount), pendingWithdrawalsAfter[0].Amount,
		"Withdrawal amount mismatch in pending withdrawal")
}
