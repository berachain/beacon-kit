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
	"io"
	"math/big"
	"net/http"

	beaconapi "github.com/attestantio/go-eth2-client/api"
	"github.com/berachain/beacon-kit/execution/requests/eip7002"
	"github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/node-api/handlers/utils"
	"github.com/berachain/beacon-kit/primitives/crypto"
	beaconmath "github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
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

	return resp, nil
}

// checkPendingPartialWithdrawals checks if there are pending partial withdrawals
func (s *BeaconKitE2ESuite) checkPendingPartialWithdrawals(stateID string) ([]types.PendingPartialWithdrawalData, error) {
	resp, err := s.getPendingPartialWithdrawals(stateID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// For debugging
	s.T().Logf("Raw getPendingPartialWithdrawals response: %s", string(body))

	var pendingPartialWithdrawals types.PendingPartialWithdrawalsResponse

	err = json.Unmarshal(body, &pendingPartialWithdrawals)
	if err != nil {
		return nil, err
	}

	s.Require().Equal(pendingPartialWithdrawals.Version, version.Name(version.Electra()))
	s.Require().False(pendingPartialWithdrawals.ExecutionOptimistic)
	s.Require().True(pendingPartialWithdrawals.Finalized)

	// Parse the actual withdrawals data
	var withdrawals []types.PendingPartialWithdrawalData

	dataBytes, err := json.Marshal(pendingPartialWithdrawals.Data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(dataBytes, &withdrawals)
	if err != nil {
		return nil, err
	}

	// Even if there are no pending withdrawals, the format should be valid
	s.T().Logf("Number of pending partial withdrawals: %d", len(withdrawals))

	return withdrawals, nil
}

// TestSubmitPartialWithdrawalTransaction tests submitting a partial withdrawal transaction via the withdrawal contract
func (s *BeaconKitE2ESuite) TestSubmitPartialWithdrawalTransaction() {
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

	// Set withdrawal amount (in Gwei) - requesting 1 BERA (10^9 Gwei)
	withdrawalAmount := beaconmath.Gwei(1_000_000_000)

	// Create an rpc client using the load balancer URL
	rpcClient, err := rpc.Dial(s.JSONRPCBalancer().URL())
	s.Require().NoError(err)
	defer rpcClient.Close()

	// Wrap the RPC client
	rpcWrapper := &rpcWrapper{Client: rpcClient}

	// Get current block number before withdrawal
	blkNum, err := s.JSONRPCBalancer().BlockNumber(s.Ctx())
	s.Require().NoError(err)
	s.T().Logf("Block number before withdrawal: %d", blkNum)

	// Get the withdrawal fee
	fee, err := eip7002.GetWithdrawalFee(s.Ctx(), rpcWrapper)
	s.Require().NoError(err)
	s.T().Logf("Withdrawal fee: %s wei", fee.String())

	// Create the withdrawal transaction data
	withdrawalTxData, err := eip7002.CreateWithdrawalRequestData(blsPubkey, withdrawalAmount)
	s.Require().NoError(err)

	// Use a pre-loaded key that has funds
	privateKey, err := ethcrypto.HexToECDSA("fffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306")
	s.Require().NoError(err)

	chainID, err := s.JSONRPCBalancer().ChainID(s.Ctx())
	s.Require().NoError(err)
	signer := gethcore.NewPragueSigner(chainID)

	// Get the sender's nonce
	var nonce hexutil.Uint64
	err = rpcClient.CallContext(s.Ctx(), &nonce, "eth_getTransactionCount",
		common.HexToAddress("0x20f33ce90a13a4b5e7697e3544c3083b8f8a51d4"), "latest",
	)
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

	// Check for pending partial withdrawals before submitting the transaction.
	// This is to ensure that the withdrawal is not already in the queue.
	pendingWithdrawalsBefore, err := s.checkPendingPartialWithdrawals(utils.StateIDHead)
	s.Require().NoError(err)
	s.Require().Len(pendingWithdrawalsBefore, 0)

	// Send the transaction
	var txHash common.Hash
	err = rpcClient.CallContext(s.Ctx(), &txHash, "eth_sendRawTransaction", hexutil.Encode(txBytes))
	s.Require().NoError(err)
	s.T().Logf("Withdrawal transaction submitted: %s", txHash.Hex())

	// wait for 3 blocks to be mined after submitting the transaction
	err = s.WaitForNBlockNumbers(3)
	s.Require().NoError(err)

	// Now get the transaction receipt
	var receipt map[string]interface{}
	err = rpcClient.CallContext(s.Ctx(), &receipt, "eth_getTransactionReceipt", txHash.Hex())
	s.Require().NoError(err)
	s.Require().NotNil(receipt, "Transaction receipt should not be nil")

	// Get block number where the withdrawal transaction was included
	blockNumStr, ok := receipt["blockNumber"].(string)
	s.Require().True(ok, "Block number should be a string")
	blockNum, err := hexutil.DecodeUint64(blockNumStr)
	s.Require().NoError(err)
	s.T().Logf("Withdrawal transaction included in block: %d", blockNum)

	pendingWithdrawalsAfter, err := s.checkPendingPartialWithdrawals(utils.StateIDHead)
	s.Require().NoError(err)
	s.Require().Len(pendingWithdrawalsAfter, 1)
	s.Require().Equal(validatorIndex, fmt.Sprintf("%d", pendingWithdrawalsAfter[0].ValidatorIndex))
	s.Require().Equal(uint64(withdrawalAmount), pendingWithdrawalsAfter[0].Amount)
}
