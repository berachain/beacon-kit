// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
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
	"crypto/rand"
	"math/big"

	"github.com/berachain/beacon-kit/config/spec"
	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/geth-primitives/deposit"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/testing/e2e/config"
	"github.com/cometbft/cometbft/crypto/bls12381"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	coretypes "github.com/ethereum/go-ethereum/core/types"
)

const (
	// NumDepositsLoad is the number of deposits to load in the Deposit
	// Robustness e2e test.
	NumDepositsLoad = 500

	// DepositAmount is the amount of BERA to deposit.
	DepositAmount = 32e18
)

func (s *BeaconKitE2ESuite) TestDepositRobustness() {
	// Get the consensus client.
	client := s.ConsensusClients()[config.DefaultClient]
	s.Require().NotNil(client)

	client2 := s.ConsensusClients()[config.AlternateClient]
	s.Require().NotNil(client2)

	// Sender account
	sender := s.TestAccounts()[1]

	// Get the block num
	blkNum, err := s.JSONRPCBalancer().BlockNumber(s.Ctx())
	s.Require().NoError(err)

	// Get the chain ID.
	chainID, err := s.JSONRPCBalancer().ChainID(s.Ctx())
	s.Require().NoError(err)

	// Get original evm balance
	balance, err := s.JSONRPCBalancer().BalanceAt(
		s.Ctx(),
		sender.Address(),
		big.NewInt(int64(blkNum)),
	)
	s.Require().NoError(err)

	// TODO: FIX KURTOSIS BUG
	// // Kill node 2
	// _, err = client2.Stop(s.Ctx())
	// s.Require().NoError(err)

	// Bind the deposit contract.
	dc, err := deposit.NewDepositContract(
		gethcommon.HexToAddress(spec.DefaultDepositContractAddress),
		s.JSONRPCBalancer(),
	)
	s.Require().NoError(err)

	// Get the nonce.
	nonce, err := s.JSONRPCBalancer().NonceAt(
		s.Ctx(),
		sender.Address(),
		big.NewInt(int64(blkNum)),
	)
	s.Require().NoError(err)

	var (
		tx      *coretypes.Transaction
		receipt *coretypes.Receipt
	)
	for i := range NumDepositsLoad {
		// Create a deposit transaction. Use the default validators' pubkeys
		// if exists, otherwise pubkey is a random 48 byte slice.
		var pubkey []byte
		var pk *bls12381.PubKey

		switch i {
		case 0:
			pubkey, err = client.GetPubKey(s.Ctx())
			s.Require().NoError(err)
			pk, err = bls12381.NewPublicKeyFromBytes(pubkey)
			s.Require().NoError(err)
			pubkey = pk.Compress()
			s.Require().Len(pubkey, 48)
		case 1:
			pubkey, err = client2.GetPubKey(s.Ctx())
			s.Require().NoError(err)
			pk, err = bls12381.NewPublicKeyFromBytes(pubkey)
			s.Require().NoError(err)
			pubkey = pk.Compress()
			s.Require().Len(pubkey, 48)
		default:
			pubkey = make([]byte, 48)
			_, err = rand.Read(pubkey)
			s.Require().NoError(err)
		}

		tx, err = s.generateNewDepositTx(
			dc,
			sender.Address(),
			sender.SignerFunc(chainID),
			big.NewInt(int64(nonce+uint64(i))),
			pubkey,
		)
		s.Require().NoError(err)
		s.Logger().
			Info("Deposit transaction created", "txHash", tx.Hash().Hex())
		if i == NumDepositsLoad-1 {
			s.Logger().Info(
				"Waiting for deposit transaction to be mined", "txHash",
				tx.Hash().Hex(),
			)
			receipt, err = bind.WaitMined(s.Ctx(), s.JSONRPCBalancer(), tx)
			s.Require().NoError(err)
			s.Require().Equal(coretypes.ReceiptStatusSuccessful, receipt.Status)
			s.Logger().
				Info("Deposit transaction mined", "txHash", receipt.TxHash.Hex())
		}
	}

	// wait blocks
	blkNum, err = s.JSONRPCBalancer().BlockNumber(s.Ctx())
	s.Require().NoError(err)
	targetBlkNum := blkNum + 10
	err = s.WaitForFinalizedBlockNumber(targetBlkNum)
	s.Require().NoError(err)

	// Check to see if evm balance decreased.
	postDepositBalance, err := s.JSONRPCBalancer().BalanceAt(
		s.Ctx(),
		sender.Address(),
		big.NewInt(int64(targetBlkNum)),
	)
	s.Require().NoError(err)

	// Check that the eth spent is somewhere~ (gas) between
	// upper bound: 32ether * 500 + 2ether
	// lower bound: 32ether * 500
	twoEther := big.NewInt(2e18)
	depositAmt, _ := big.NewFloat(DepositAmount).Int(nil)
	totalAmt := new(big.Int).Mul(depositAmt, big.NewInt(NumDepositsLoad))
	upperBound := new(big.Int).Add(totalAmt, twoEther)
	amtSpent := new(big.Int).Sub(balance, postDepositBalance)

	s.Require().Equal(amtSpent.Cmp(totalAmt), 1)
	s.Require().Equal(amtSpent.Cmp(upperBound), -1)

	// TODO: FIX KURTOSIS BUG
	// // Start node 2 again
	// _, err = client2.Start(s.Ctx(), s.Enclave())
	// s.Require().NoError(err)

	// Update client2's reference

	// err = s.SetupConsensusClients()
	// s.Require().NoError(err)
	// client2 = s.ConsensusClients()[AlternateClient]
	// s.Require().NotNil(client2)

	// Give time for the node to catch up
	err = s.WaitForNBlockNumbers(20)
	s.Require().NoError(err)

	// Compare height of nodes 1 and 2
	height, err := client.ABCIInfo(s.Ctx())
	s.Require().NoError(err)
	height2, err := client2.ABCIInfo(s.Ctx())
	s.Require().NoError(err)
	s.Require().
		InDelta(height.Response.LastBlockHeight, height2.Response.LastBlockHeight, 1)
}

func (s *BeaconKitE2ESuite) generateNewDepositTx(
	dc *deposit.DepositContract,
	sender gethcommon.Address,
	signer bind.SignerFn,
	nonce *big.Int,
	pubkey []byte,
) (*coretypes.Transaction, error) {
	// Generate the credentials.
	credentials := types.NewCredentialsFromExecutionAddress(
		common.ExecutionAddress(s.GenesisAccount().Address()),
	)

	// Generate the signature.
	signature := [96]byte{}
	s.Require().Len(signature[:], 96)

	val, _ := big.NewFloat(DepositAmount).Int(nil)
	return dc.Deposit(&bind.TransactOpts{
		From:   sender,
		Value:  val,
		Signer: signer,
		Nonce:  nonce,
	}, pubkey, credentials[:], signature[:], sender)
}
