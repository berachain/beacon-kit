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

package e2e_test

import (
	"math/big"

	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/berachain/beacon-kit/config/spec"
	"github.com/berachain/beacon-kit/geth-primitives/deposit"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/testing/e2e/config"
	"github.com/berachain/beacon-kit/testing/e2e/suite/types"
	"github.com/cometbft/cometbft/crypto/bls12381"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	coretypes "github.com/ethereum/go-ethereum/core/types"
)

const (
	// NumDepositsLoad is the number of deposits to load in the Deposit Robustness e2e test.
	NumDepositsLoad uint64 = 500
)

// Contains pre-test state for validator info to facilitate validation of the post-state.
type ValidatorTestStruct struct {
	Index                 uint64
	Power                 *big.Int
	WithdrawalBalance     *big.Int
	WithdrawalCredentials [32]byte
	Name                  string
	Client                *types.ConsensusClient
}

// TestDepositRobustness tests sending a large number of deposits txs to the Deposit Contract.
// Then it checks whether all the validators' voting power have increased by the correct amount.
//
// TODO:
// 1) Add staking tests for adding a new validator to the network.
// 2) Add staking tests for hitting the validator set cap and eviction.
func (s *BeaconKitE2ESuite) TestDepositRobustness() {
	// TODO: make test use configurable chain spec.
	chainSpec, err := spec.DevnetChainSpec()
	s.Require().NoError(err)

	weiPerGwei := big.NewInt(1e9)

	// This value is determined by the MIN_DEPOSIT_AMOUNT_IN_GWEI variable from the deposit contract.
	contractMinDepositAmountWei := big.NewInt(0).Mul(big.NewInt(10_000), big.NewInt(1e9*1e9))
	depositAmountWei := new(big.Int).Mul(contractMinDepositAmountWei, big.NewInt(100))
	depositAmountGwei := new(big.Int).Div(depositAmountWei, weiPerGwei)

	// Our deposits should be greater than the min deposit amount.
	s.Require().Equal(1, depositAmountWei.Cmp(contractMinDepositAmountWei))

	s.Require().Equal(
		0, int(NumDepositsLoad%config.NumValidators),
		"every validator must get an equal amount of deposits",
	)

	// Get the chain ID.
	chainID, err := s.JSONRPCBalancer().ChainID(s.Ctx())
	s.Require().NoError(err)

	// Bind the deposit contract.
	depositContractAddress := gethcommon.Address(chainSpec.DepositContractAddress())

	dc, err := deposit.NewDepositContract(depositContractAddress, s.JSONRPCBalancer())
	s.Require().NoError(err)

	// Enforce the deposit count at genesis is equal to the number of validators.
	depositCount, err := dc.DepositCount(&bind.CallOpts{
		BlockNumber: big.NewInt(0),
	})
	s.Require().NoError(err)
	s.Require().Equal(uint64(config.NumValidators), depositCount,
		"initial deposit count should match number of validators")

	// Check that the genesis deposits root is not empty. It is important that this value is
	// consistent across all EL nodes to ensure the EL has a consistent view of the CL deposits
	// at genesis. If the EL chain progresses past the genesis block, this is guaranteed.
	genesisRoot, err := dc.GenesisDepositsRoot(&bind.CallOpts{
		BlockNumber: big.NewInt(0),
	})
	s.Require().NoError(err)
	s.Require().False(genesisRoot == [32]byte{})

	apiClient := s.ConsensusClients()[config.ClientValidator0]
	s.Require().NotNil(apiClient)

	// Grab genesis validators to get withdrawal creds.
	s.Require().NoError(apiClient.Connect(s.Ctx()))
	response, err := apiClient.Validators(s.Ctx(), &api.ValidatorsOpts{
		State:   "genesis",
		Indices: []phase0.ValidatorIndex{0, 1, 2, 3, 4},
	})
	s.Require().NoError(err)
	s.Require().NotNil(response, "Validators response should not be nil")
	s.Require().NotNil(response.Data, "Validators data should not be nil")

	vals := response.Data
	s.Require().Len(vals, config.NumValidators)
	s.Require().Len(s.ConsensusClients(), config.NumValidators)

	// Grab pre-state data for each validator.
	validators := make([]*ValidatorTestStruct, config.NumValidators)
	var idx uint64
	for name, client := range s.ConsensusClients() {
		power, cErr := client.GetConsensusPower(s.Ctx())
		s.Require().NoError(cErr)

		s.Require().Contains(vals, phase0.ValidatorIndex(idx))
		val := vals[phase0.ValidatorIndex(idx)]
		s.Require().NotNil(val)
		s.Require().NotNil(val.Validator)
		creds := [32]byte(val.Validator.WithdrawalCredentials)
		withdrawalAddress := gethcommon.Address(creds[12:])
		withdrawalBalance, jErr := s.JSONRPCBalancer().BalanceAt(s.Ctx(), withdrawalAddress, nil)
		s.Require().NoError(jErr)

		// Populate the validators testing struct so we can keep track of the pre-state.
		validators[idx] = &ValidatorTestStruct{
			Index:                 idx,
			Power:                 new(big.Int).SetUint64(power),
			WithdrawalBalance:     withdrawalBalance,
			WithdrawalCredentials: creds,
			Name:                  name,
			Client:                client,
		}
		idx++
	}

	// Sender account
	sender := s.TestAccounts()[0]

	// Get the block num
	blkNum, err := s.JSONRPCBalancer().BlockNumber(s.Ctx())
	s.Require().NoError(err)

	// Get original evm balance
	balance, err := s.JSONRPCBalancer().BalanceAt(
		s.Ctx(), sender.Address(), new(big.Int).SetUint64(blkNum),
	)
	s.Require().NoError(err)

	// Get the nonce.
	nonce, err := s.JSONRPCBalancer().NonceAt(
		s.Ctx(), sender.Address(), new(big.Int).SetUint64(blkNum),
	)
	s.Require().NoError(err)

	var (
		tx           *coretypes.Transaction
		clientPubkey []byte
		pk           *bls12381.PubKey
		signature    [96]byte
		value        = depositAmountWei
		signer       = sender.SignerFunc(chainID)
		from         = sender.Address()
	)
	for i := range NumDepositsLoad {
		// Create a deposit transaction using the default validators' pubkeys.
		curVal := validators[i%config.NumValidators]
		clientPubkey, err = curVal.Client.GetPubKey(s.Ctx())
		s.Require().NoError(err)
		pk, err = bls12381.NewPublicKeyFromBytes(clientPubkey)
		s.Require().NoError(err)
		pubkey := pk.Compress()
		s.Require().Len(pubkey, 48)

		// Only the first deposit for a pubkey has a non-zero operator.
		operator := gethcommon.Address{}
		if i < config.NumValidators {
			operator = from
		}
		tx, err = dc.Deposit(&bind.TransactOpts{
			From:     from,
			Value:    value,
			Signer:   signer,
			Nonce:    new(big.Int).SetUint64(nonce + i),
			GasLimit: 1000000,
			Context:  s.Ctx(),
		}, pubkey, curVal.WithdrawalCredentials[:], signature[:], operator)
		s.Require().NoError(err)
		s.Logger().Info("Deposit tx created", "num", i+1, "hash", tx.Hash().Hex(), "value", value)
	}

	// Wait for the final deposit tx to be mined.
	s.Logger().Info(
		"Waiting for the final deposit tx to be mined",
		"num", NumDepositsLoad, "hash", tx.Hash().Hex(),
	)
	receipt, err := bind.WaitMined(s.Ctx(), s.JSONRPCBalancer(), tx)
	s.Require().NoError(err)
	s.Require().Equal(coretypes.ReceiptStatusSuccessful, receipt.Status)
	s.Logger().Info("Final deposit tx mined successfully", "hash", receipt.TxHash.Hex())

	// Give time for the nodes to catch up.
	err = s.WaitForNBlockNumbers(NumDepositsLoad / chainSpec.MaxDepositsPerBlock())
	s.Require().NoError(err)

	// Compare height of nodes 0 and 1
	height, err := validators[0].Client.ABCIInfo(s.Ctx())
	s.Require().NoError(err)
	height2, err := validators[1].Client.ABCIInfo(s.Ctx())
	s.Require().NoError(err)
	s.Require().InDelta(height.LastBlockHeight, height2.LastBlockHeight, 1)

	// Check to see if evm balance decreased.
	postDepositBalance, err := s.JSONRPCBalancer().BalanceAt(s.Ctx(), sender.Address(), nil)
	s.Require().NoError(err)

	// Check that the eth spent is somewhere~ (gas) between
	// (depositAmountWei * NumDepositsLoad, depositAmountWei * NumDepositsLoad + 2ether)
	lowerBound := new(big.Int).Mul(value, new(big.Int).SetUint64(NumDepositsLoad))
	upperBound := new(big.Int).Add(lowerBound, big.NewInt(2e18))
	amtSpent := new(big.Int).Sub(balance, postDepositBalance)

	s.Require().Equal(1, amtSpent.Cmp(lowerBound), "amount spent is less than lower bound")
	s.Require().Equal(-1, amtSpent.Cmp(upperBound), "amount spent is greater than upper bound")

	// Check that all validators' voting power have increased by
	// (NumDepositsLoad / NumValidators) * depositAmountWei
	// after the end of the epoch (next multiple of SlotsPerEpoch after receipt.BlockNumber).
	blkNum, err = s.JSONRPCBalancer().BlockNumber(s.Ctx())
	s.Require().NoError(err)
	nextEpoch := chainSpec.SlotToEpoch(math.Slot(blkNum)) + 1
	nextEpochBlockNum := nextEpoch.Unwrap() * chainSpec.SlotsPerEpoch()
	err = s.WaitForFinalizedBlockNumber(nextEpochBlockNum + 1)
	s.Require().NoError(err)

	increaseAmt := new(big.Int).Mul(depositAmountGwei, big.NewInt(int64(NumDepositsLoad/config.NumValidators)))

	for _, val := range validators {
		// Consensus Power is in Gwei.
		powerAfterRaw, cErr := val.Client.GetConsensusPower(s.Ctx())
		s.Require().NoError(cErr)
		powerAfter := new(big.Int).SetUint64(powerAfterRaw)
		powerDiff := new(big.Int).Sub(powerAfter, val.Power)

		// withdrawal balance is in Wei, so we'll convert it to Gwei.
		withdrawalAddress := gethcommon.Address(val.WithdrawalCredentials[12:])
		withdrawalBalanceAfter, jErr := s.JSONRPCBalancer().BalanceAt(s.Ctx(), withdrawalAddress, nil)
		s.Require().NoError(jErr)
		withdrawalDiff := new(big.Int).Sub(withdrawalBalanceAfter, val.WithdrawalBalance)
		withdrawalDiff.Div(withdrawalDiff, weiPerGwei)

		// TODO: currently the kurtosis devnet sets the withdrawal address the same for all validators.
		// We simply validate that the balance is NumValidators times larger than we expect it to be.
		withdrawalDiff.Div(withdrawalDiff, new(big.Int).SetUint64(config.NumValidators))

		// Verify input balance is equal to the power + withdrawal balances.
		s.Require().Equal(increaseAmt, new(big.Int).Add(powerDiff, withdrawalDiff))
	}
}
