// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package e2e_test

import (
	"math/big"

	"github.com/berachain/beacon-kit/mod/execution/pkg/deposit"
	byteslib "github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/testing/e2e/suite"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	coretypes "github.com/ethereum/go-ethereum/core/types"
)

const (
	// DepositContractAddress is the address of the deposit contract.
	DepositContractAddress = "0x00000000219ab540356cbb839cbe05303d7705fa"
	DefaultClient          = "cl-validator-beaconkit-0"
	AlternateClient        = "cl-validator-beaconkit-1"
	NumDepositsLoad        = 500
)

func (s *BeaconKitE2ESuite) TestDepositRobustness() {
	// Get the consensus client.
	client := s.ConsensusClients()[DefaultClient]
	s.Require().NotNil(client)

	client2 := s.ConsensusClients()[AlternateClient]
	s.Require().NotNil(client2)

	// Sender account
	sender := s.TestAccounts()[1]

	// Get the public key.
	pubkey, err := client.GetPubKey(s.Ctx())
	s.Require().NoError(err)
	s.Require().Len(pubkey, 48)

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

	nonce, err := s.JSONRPCBalancer().NonceAt(
		s.Ctx(),
		sender.Address(),
		big.NewInt(int64(blkNum)),
	)
	s.Require().NoError(err)

	// Kill node 2
	_, err = client2.Stop(s.Ctx())
	s.Require().NoError(err)

	// Bind the deposit contract.
	dc, err := deposit.NewBeaconDepositContract(
		common.HexToAddress(DepositContractAddress),
		s.JSONRPCBalancer(),
	)
	s.Require().NoError(err)

	for i := range NumDepositsLoad {
		var receipt *coretypes.Receipt
		var tx *coretypes.Transaction
		// Create a deposit transaction.
		tx, err = s.generateNewDepositTx(
			dc,
			sender.Address(),
			sender.SignerFunc(chainID),
			big.NewInt(int64(nonce+uint64(i))),
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
			s.Require().Equal(uint64(1), receipt.Status)
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
	// upper bound: 32ether * 500 + 1ether
	// lower bound: 32ether * 500
	oneEther := big.NewInt(1e18)
	totalAmt := new(big.Int).Mul(oneEther, big.NewInt(NumDepositsLoad*32))
	upperBound := new(big.Int).Add(totalAmt, oneEther)
	amtSpent := new(big.Int).Sub(balance, postDepositBalance)

	s.Require().Equal(amtSpent.Cmp(totalAmt), 1)
	s.Require().Equal(amtSpent.Cmp(upperBound), -1)

	// Start node 2 again
	_, err = client2.Start(s.Ctx(), s.Enclave())
	s.Require().NoError(err)

	// Update client2's reference
	err = s.SetupConsensusClients()
	s.Require().NoError(err)
	client2 = s.ConsensusClients()[AlternateClient]
	s.Require().NotNil(client2)

	// Give time for the node to catch up
	err = s.WaitForNBlockNumbers(5)
	s.Require().NoError(err)

	// Compare height of node 1 and 2
	height, err := client.ABCIInfo(s.Ctx())
	s.Require().NoError(err)
	height2, err := client2.ABCIInfo(s.Ctx())
	s.Require().NoError(err)
	s.Require().
		Equal(height.Response.LastBlockHeight, height2.Response.LastBlockHeight)
}

func (s *BeaconKitE2ESuite) generateNewDepositTx(
	dc *deposit.BeaconDepositContract,
	sender common.Address,
	signer bind.SignerFn,
	nonce *big.Int,
) (*coretypes.Transaction, error) {
	// Get the consensus client.
	client := s.ConsensusClients()[DefaultClient]
	s.Require().NotNil(client)

	pubkey, err := client.GetPubKey(s.Ctx())
	s.Require().NoError(err)
	s.Require().Len(pubkey, 48)

	// Generate the credentials.
	credentials := byteslib.PrependExtendToSize(
		s.GenesisAccount().Address().Bytes(),
		32,
	)
	credentials[0] = 0x01

	// Generate the signature.
	signature := [96]byte{}
	s.Require().Len(signature[:], 96)

	val, _ := big.NewFloat(32e18).Int(nil)
	return dc.Deposit(&bind.TransactOpts{
		From:   sender,
		Value:  val,
		Signer: signer,
		Nonce:  nonce,
	}, pubkey, credentials, 32*suite.OneGwei, signature[:])
}
