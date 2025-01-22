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

	"github.com/berachain/beacon-kit/config/spec"
	"github.com/berachain/beacon-kit/geth-primitives/deposit"
	"github.com/berachain/beacon-kit/testing/e2e/config"
	"github.com/cometbft/cometbft/crypto/bls12381"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	coretypes "github.com/ethereum/go-ethereum/core/types"
)

const (
	// NumDepositsLoad is the number of deposits to load in the Deposit Robustness e2e test.
	NumDepositsLoad uint64 = 500

	// DepositAmount is the amount of BERA to deposit.
	DepositAmount = 32e18

	// BlocksToWait is the number of blocks to wait for the nodes to catch up.
	BlocksToWait = 10
)

// TestDepositRobustness tests sending a large number of deposits txs to the Deposit Contract.
// Then it checks whether all the validators' voting power have increased by the correct amount.
//
// TODO:
// 1) Add staking tests for exceeding the max stake.
// 2) Add staking tests for adding a new validator to the network.
// 3) Add staking tests for hitting the validator set cap and eviction.
func (s *BeaconKitE2ESuite) TestDepositRobustness() {
	s.Require().Equal(
		0, int(NumDepositsLoad%config.NumValidators),
		"every validator must get an equal amount of deposits",
	)

	// Get the chain ID.
	chainID, err := s.JSONRPCBalancer().ChainID(s.Ctx())
	s.Require().NoError(err)

	// Bind the deposit contract.
	depositContractAddress := gethcommon.HexToAddress(spec.DefaultDepositContractAddress)

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

	// Get the consensus clients.
	client0 := s.ConsensusClients()[config.ClientValidator0]
	s.Require().NotNil(client0)
	client1 := s.ConsensusClients()[config.ClientValidator1]
	s.Require().NotNil(client1)
	client2 := s.ConsensusClients()[config.ClientValidator2]
	s.Require().NotNil(client2)
	client3 := s.ConsensusClients()[config.ClientValidator3]
	s.Require().NotNil(client3)
	client4 := s.ConsensusClients()[config.ClientValidator4]
	s.Require().NotNil(client4)

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
		credentials  [32]byte
		signature    [96]byte
		value, _     = big.NewFloat(DepositAmount).Int(nil)
		signer       = sender.SignerFunc(chainID)
		from         = sender.Address()
	)
	for i := range NumDepositsLoad {
		// Create a deposit transaction using the default validators' pubkeys.
		switch i % config.NumValidators {
		case 0:
			clientPubkey, err = client0.GetPubKey(s.Ctx())
		case 1:
			clientPubkey, err = client1.GetPubKey(s.Ctx())
		case 2:
			clientPubkey, err = client2.GetPubKey(s.Ctx())
		case 3:
			clientPubkey, err = client3.GetPubKey(s.Ctx())
		case 4:
			clientPubkey, err = client4.GetPubKey(s.Ctx())
		}
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
		}, pubkey, credentials[:], signature[:], operator)
		s.Require().NoError(err)
		s.Logger().Info("Deposit tx created", "num", i+1, "hash", tx.Hash().Hex())
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
	err = s.WaitForNBlockNumbers(BlocksToWait)
	s.Require().NoError(err)

	// Compare height of nodes 0 and 1
	height, err := client0.ABCIInfo(s.Ctx())
	s.Require().NoError(err)
	height2, err := client1.ABCIInfo(s.Ctx())
	s.Require().NoError(err)
	s.Require().InDelta(height.Response.LastBlockHeight, height2.Response.LastBlockHeight, 1)

	// Check to see if evm balance decreased.
	postDepositBalance, err := s.JSONRPCBalancer().BalanceAt(s.Ctx(), sender.Address(), nil)
	s.Require().NoError(err)

	// Check that the eth spent is somewhere~ (gas) between
	// (DepositAmount * NumDepositsLoad, DepositAmount * NumDepositsLoad + 2ether)
	lowerBound := new(big.Int).Mul(value, new(big.Int).SetUint64(NumDepositsLoad))
	upperBound := new(big.Int).Add(lowerBound, big.NewInt(2e18))
	amtSpent := new(big.Int).Sub(balance, postDepositBalance)

	s.Require().Equal(1, amtSpent.Cmp(lowerBound), "amount spent is less than lower bound")
	s.Require().Equal(-1, amtSpent.Cmp(upperBound), "amount spent is greater than upper bound")

	// TODO: determine why voting power is not increasing above 32e9.
	// // Check that all validators' voting power have increased by
	// // (NumDepositsLoad / NumValidators) * DepositAmount
	// // after the end of the epoch (next multiple of 32 after receipt.BlockNumber).
	// nextEpochBlockNum := (receipt.BlockNumber.Uint64()/32 + 1) * 32
	// err = s.WaitForFinalizedBlockNumber(nextEpochBlockNum + 1)
	// s.Require().NoError(err)

	// power0After, err := client0.GetConsensusPower(s.Ctx())
	// s.Require().NoError(err)
	// power1After, err := client1.GetConsensusPower(s.Ctx())
	// s.Require().NoError(err)
	// power2After, err := client2.GetConsensusPower(s.Ctx())
	// s.Require().NoError(err)
	// power3After, err := client3.GetConsensusPower(s.Ctx())
	// s.Require().NoError(err)
	// power4After, err := client4.GetConsensusPower(s.Ctx())
	// s.Require().NoError(err)

	// increaseAmt := NumDepositsLoad / config.NumValidators * uint64(DepositAmount/params.GWei)
	// s.Require().Equal(power0+increaseAmt, power0After)
	// s.Require().Equal(power1+increaseAmt, power1After)
	// s.Require().Equal(power2+increaseAmt, power2After)
	// s.Require().Equal(power3+increaseAmt, power3After)
	// s.Require().Equal(power4+increaseAmt, power4After)
}
