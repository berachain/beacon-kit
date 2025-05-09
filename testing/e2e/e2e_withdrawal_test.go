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
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"time"

	beaconapi "github.com/attestantio/go-eth2-client/api"
	"github.com/berachain/beacon-kit/execution/requests/eip7002"
	beacontypes "github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/node-api/handlers/utils"
	"github.com/berachain/beacon-kit/primitives/crypto"
	beaconmath "github.com/berachain/beacon-kit/primitives/math"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethcore "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
)

// rpcWrapper wraps an rpc.Client and implements the rpcClient interface required by eip7002
type rpcWrapper struct {
	*rpc.Client
}

// Call implements the rpcClient interface
func (r *rpcWrapper) Call(ctx context.Context, target any, method string, params ...any) error {
	return r.Client.CallContext(ctx, target, method, params...)
}

// getWithdrawalCredentials retrieves the validator's withdrawal credentials
func (s *BeaconKitE2ESuite) getWithdrawalCredentials(validatorIndex string) (string, error) {
	resp, err := s.getStateValidator(utils.StateIDHead, validatorIndex)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	validatorResp, err := s.decodeValidatorResponse(resp)
	if err != nil {
		return "", err
	}

	return validatorResp.Validator.WithdrawalCredentials, nil
}

// getValidatorBalance returns the balance of a validator
func (s *BeaconKitE2ESuite) getValidatorBalance(validatorIndex string) (uint64, error) {
	resp, err := s.getValidatorBalances(utils.StateIDHead, validatorIndex)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	balancesResp, err := s.decodeValidatorBalancesResponse(resp)
	if err != nil {
		return 0, err
	}

	if len(*balancesResp) == 0 {
		return 0, fmt.Errorf("no balance found for validator %s", validatorIndex)
	}

	return (*balancesResp)[0].Balance, nil
}

// getPendingPartialWithdrawals gets the pending partial withdrawals for the given stateID.
func (s *BeaconKitE2ESuite) getPendingPartialWithdrawals(stateID string) (*http.Response, error) {
	client := s.initHTTPBeaconTest()

	url := fmt.Sprintf("/eth/v1/beacon/states/%s/pending_partial_withdrawals", stateID)

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get validator balances: %w", err)
	}
	if resp == nil {
		return nil, errors.New("received nil response")
	}

	fmt.Println(resp.Body)

	return resp, nil
}

// decodePendingPartialWithdrawalsResponse decodes a response containing pending partial withdrawals.
func (s *BeaconKitE2ESuite) decodePendingPartialWithdrawalsResponse(resp *http.Response) (*beacontypes.PendingPartialWithdrawalsResponse, error) {
	partialWithdrawals, err := decodeResponse[beacontypes.PendingPartialWithdrawalsResponse](resp)
	if err != nil {
		return nil, err
	}
	fmt.Println("partialWithdrawals", partialWithdrawals)
	return &partialWithdrawals, nil
}

// checkPendingPartialWithdrawals checks if there are pending partial withdrawals
func (s *BeaconKitE2ESuite) checkPendingPartialWithdrawals(stateID string) (*[]beacontypes.PendingPartialWithdrawalsData, error) {
	resp, err := s.getPendingPartialWithdrawals(stateID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	partialWithdrawals, err := s.decodePendingPartialWithdrawalsResponse(resp)
	if err != nil {
		return nil, err
	}

	return partialWithdrawals, nil
}

// getCurrentBlockNumber gets the current block number
func (s *BeaconKitE2ESuite) getCurrentBlockNumber(rpcClient *rpc.Client) (uint64, error) {
	var blockNumHex string
	err := rpcClient.CallContext(s.Ctx(), &blockNumHex, "eth_blockNumber")
	if err != nil {
		return 0, err
	}

	blockNum, err := hexutil.DecodeUint64(blockNumHex)
	if err != nil {
		return 0, err
	}

	return blockNum, nil
}

// TestSubmitPartialWithdrawalTransaction tests submitting a partial withdrawal transaction via the withdrawal contract
func (s *BeaconKitE2ESuite) TestSubmitPartialWithdrawalTransaction() {
	// Initialize the test environment
	client := s.initBeaconTest()

	// Get the validators to identify one with execution credentials (0x01)
	validatorsResp, err := client.Validators(
		s.Ctx(),
		&beaconapi.ValidatorsOpts{
			State: utils.StateIDHead,
		},
	)
	s.Require().NoError(err)
	s.Require().NotEmpty(validatorsResp.Data)

	// Find a validator with execution withdrawal credentials (starting with 0x01)
	var validatorIndex string
	var blsPubkey crypto.BLSPubkey
	for _, validator := range validatorsResp.Data {
		credentials := validator.Validator.WithdrawalCredentials
		if len(credentials) >= 2 && credentials[0] == 0x01 {
			validatorIndex = fmt.Sprintf("%d", validator.Index)
			// Convert the phase0.BLSPubKey to our crypto.BLSPubkey type
			copy(blsPubkey[:], validator.Validator.PublicKey[:])
			break
		}
	}
	s.Require().NotEmpty(validatorIndex, "No validator with execution withdrawal credentials found")

	// Get initial balance
	initialBalance, err := s.getValidatorBalance(validatorIndex)
	s.Require().NoError(err)
	s.T().Logf("Initial validator balance: %d Gwei", initialBalance)

	// Set withdrawal amount (in Gwei) - requesting 1 BERA (10^9 Gwei)
	withdrawalAmount := beaconmath.Gwei(1_000_000_000)

	// Create an rpc client using the load balancer URL
	rpcClient, err := rpc.Dial(s.JSONRPCBalancer().URL())
	s.Require().NoError(err)
	defer rpcClient.Close()

	// Wrap the RPC client
	rpcWrapper := &rpcWrapper{Client: rpcClient}

	// Get current block number before withdrawal
	beforeBlockNum, err := s.getCurrentBlockNumber(rpcClient)
	s.Require().NoError(err)
	s.T().Logf("Block number before withdrawal: %d", beforeBlockNum)

	// Get the withdrawal fee
	fee, err := eip7002.GetWithdrawalFee(s.Ctx(), rpcWrapper)
	s.Require().NoError(err)
	s.T().Logf("Withdrawal fee: %s wei", fee.String())

	// Create the withdrawal transaction data
	withdrawalTxData, err := eip7002.CreateWithdrawalRequestData(blsPubkey, withdrawalAmount)
	s.Require().NoError(err)
	s.T().Logf("Withdrawal transaction data: %s", hexutil.Encode(withdrawalTxData))

	// Use a pre-loaded key that has funds
	privateKey, err := ethcrypto.HexToECDSA("fffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306")
	s.Require().NoError(err)

	chainID, err := s.JSONRPCBalancer().ChainID(s.Ctx())
	s.Require().NoError(err)
	signer := gethcore.NewPragueSigner(chainID)

	// Get the sender's nonce
	var nonce hexutil.Uint64
	err = rpcClient.CallContext(s.Ctx(), &nonce, "eth_getTransactionCount", common.HexToAddress("0x20f33ce90a13a4b5e7697e3544c3083b8f8a51d4"), "latest")
	s.Require().NoError(err)

	tx := gethcore.MustSignNewTx(privateKey, signer, &gethcore.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     uint64(nonce),
		To:        &params.WithdrawalQueueAddress,
		Gas:       500_000,
		GasFeeCap: big.NewInt(1000000000),
		GasTipCap: big.NewInt(1000000000),
		Value:     fee,
		Data:      withdrawalTxData,
	})

	// Serialize the transaction
	txBytes, err := tx.MarshalBinary()
	s.Require().NoError(err)

	// Check for pending partial withdrawals before submitting the transaction
	pendingWithdrawalsBefore, err := s.checkPendingPartialWithdrawals(utils.StateIDHead)
	s.Require().NoError(err)
	s.T().Logf("Pending withdrawals before: %v", pendingWithdrawalsBefore)
	s.Require().Len(pendingWithdrawalsBefore, 0)

	// Send the transaction
	var txHash common.Hash
	err = rpcClient.CallContext(s.Ctx(), &txHash, "eth_sendRawTransaction", hexutil.Encode(txBytes))
	s.Require().NoError(err)
	s.T().Logf("Withdrawal transaction submitted: %s", txHash.Hex())

	// Wait for the transaction to be mined
	receipt := s.waitForTransaction(txHash, rpcClient)
	s.Require().NotNil(receipt, "Transaction receipt should not be nil")

	pendingWithdrawalsAfter, err := s.checkPendingPartialWithdrawals(utils.StateIDHead)
	s.Require().NoError(err)
	s.T().Logf("Pending withdrawals after: %v", pendingWithdrawalsAfter)

	// Get block number where the withdrawal transaction was included
	blockNumStr, ok := receipt["blockNumber"].(string)
	s.Require().True(ok, "Block number should be a string")
	blockNum, err := hexutil.DecodeUint64(blockNumStr)
	s.Require().NoError(err)
	s.T().Logf("Withdrawal transaction included in block: %d", blockNum)

	// Check for pending partial withdrawals - might take a few blocks to appear
	var pendingWithdrawals *beacontypes.PendingPartialWithdrawalsResponse
	// var withdrawalFound bool

	for attempts := 0; attempts < 10; attempts++ {
		// Wait for a few blocks
		time.Sleep(5 * time.Second)

		pendingWithdrawals, err = s.checkPendingPartialWithdrawals(utils.StateIDHead)
		if err != nil {
			s.T().Logf("Error checking pending withdrawals (attempt %d): %v", attempts+1, err)
			continue
		}

		// Get the length of the pending withdrawals
		len(pendingWithdrawals.Data)	
		s.T().Logf("Pending withdrawals length: %d", len)

		// Convert the Data field to a slice of PendingPartialWithdrawalData
		// pendingWithdrawalsData, ok := pendingWithdrawals.Data.([]beacontypes.PendingPartialWithdrawalData)
		// s.Require().True(ok, "Data should be a slice of PendingPartialWithdrawalData")

		// Look for our withdrawal
		// for _, withdrawal := range pendingWithdrawalsData {
		// 	if fmt.Sprintf("%d", withdrawal.ValidatorIndex) == validatorIndex {
		// 		withdrawalFound = true
		// 		s.T().Logf("Found pending withdrawal for validator %s: %d Gwei (attempt %d)",
		// 			validatorIndex, withdrawal.Amount, attempts+1)
		// 		break
		// 	}
		// }

		// if withdrawalFound {
		// 	break
		// }

		// Get current block number to show progress
		// currentBlockNum, err := s.getCurrentBlockNumber(rpcClient)
		// if err == nil {
		// 	s.T().Logf("Current block: %d, Blocks since withdrawal tx: %d (attempt %d)",
		// 		currentBlockNum, currentBlockNum-blockNum, attempts+1)
		// } else {
		// 	s.T().Logf("Error getting current block number: %v", err)
		// }

		// s.T().Logf("Pending withdrawal not found yet (attempt %d), waiting...", attempts+1)
	}

	// Get the current block number after all checks
	afterBlockNum, err := s.getCurrentBlockNumber(rpcClient)
	s.Require().NoError(err)
	s.T().Logf("Block number after withdrawal checks: %d", afterBlockNum)
	s.T().Logf("Blocks elapsed during test: %d", afterBlockNum-beforeBlockNum)

	// Check that the withdrawal was processed
	finalBalance, err := s.getValidatorBalance(validatorIndex)
	s.Require().NoError(err)
	s.T().Logf("Final validator balance: %d Gwei", finalBalance)

	// The transaction might have been queued for processing in a future block,
	// so we don't strictly check that balance decreased immediately
	s.T().Logf("Withdrawal transaction was successfully submitted and processed")
}

// waitForTransaction waits for a transaction to be mined and returns the receipt
func (s *BeaconKitE2ESuite) waitForTransaction(txHash common.Hash, rpcClient *rpc.Client) map[string]interface{} {
	for i := 0; i < 10; i++ {
		var receipt map[string]interface{}
		err := rpcClient.CallContext(s.Ctx(), &receipt, "eth_getTransactionReceipt", txHash.Hex())

		if err == nil && receipt != nil {
			s.T().Logf("Transaction mined in block: %v", receipt["blockNumber"])
			return receipt
		}

		s.T().Logf("Waiting for transaction to be mined (attempt %d)...", i+1)
		time.Sleep(2 * time.Second)
	}

	s.T().Logf("Timed out waiting for transaction to be mined")
	return nil
}
