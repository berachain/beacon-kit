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
	"fmt"
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
	NumDepositsLoad        = 500
)

// TestDepositContract tests the deposit contract to attempt staking and
// increasing a validator's consensus power.
func (s *BeaconKitE2ESuite) TestDepositContract() {
	// Get the consensus client.
	client := s.ConsensusClients()[DefaultClient]
	s.Require().NotNil(client)

	// Get the public key.
	pubkey, err := client.GetPubKey(s.Ctx())
	s.Require().NoError(err)
	s.Require().Len(pubkey, 48)

	// Get the consensus power.
	_, err = client.GetConsensusPower(s.Ctx())
	s.Require().NoError(err)

	// Get the block num
	blkNum, err := s.JSONRPCBalancer().BlockNumber(s.Ctx())
	s.Require().NoError(err)

	// Get the chain ID.
	chainID, err := s.JSONRPCBalancer().ChainID(s.Ctx())
	s.Require().NoError(err)

	// Get original evm balance
	balance, err := s.JSONRPCBalancer().BalanceAt(
		s.Ctx(),
		s.GenesisAccount().Address(),
		big.NewInt(int64(blkNum)),
	)
	s.Require().NoError(err)

	nonce, err := s.JSONRPCBalancer().NonceAt(
		s.Ctx(),
		s.GenesisAccount().Address(),
		big.NewInt(int64(blkNum)),
	)
	s.Require().NoError(err)

	// Create a deposit transaction.
	tx, err := s.generateNewDepositTx(
		s.GenesisAccount().Address(),
		s.GenesisAccount().SignerFunc(chainID),
		big.NewInt(int64(nonce)),
	)
	s.Require().NoError(err)

	// Wait for the transaction to be mined.
	var receipt *coretypes.Receipt
	receipt, err = bind.WaitMined(s.Ctx(), s.JSONRPCBalancer(), tx)
	s.Require().NoError(err)
	s.Require().Equal(uint64(1), receipt.Status)
	s.Logger().Info("Deposit transaction mined", "txHash", receipt.TxHash.Hex())

	// Wait for the log to be processed.
	targetBlkNum := blkNum + 10
	err = s.WaitForFinalizedBlockNumber(targetBlkNum)
	s.Require().NoError(err)

	// Check to see if evm balance decreased.
	postDepositBalance, err := s.JSONRPCBalancer().BalanceAt(
		s.Ctx(),
		s.GenesisAccount().Address(),
		big.NewInt(int64(targetBlkNum)),
	)
	s.Require().NoError(err)
	s.Require().Equal(postDepositBalance.Cmp(balance), -1)

	newPower, err := client.GetConsensusPower(s.Ctx())
	s.Require().NoError(err)
	s.Require().Equal(newPower, 32*suite.OneGwei)
}

func (s *BeaconKitE2ESuite) TestDepositRobustness() {
	// Get the consensus client.
	client := s.ConsensusClients()[DefaultClient]
	s.Require().NotNil(client)

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
		s.GenesisAccount().Address(),
		big.NewInt(int64(blkNum)),
	)
	s.Require().NoError(err)

	nonce, err := s.JSONRPCBalancer().NonceAt(
		s.Ctx(),
		s.GenesisAccount().Address(),
		big.NewInt(int64(blkNum)),
	)
	s.Require().NoError(err)

	// TODO: kill a node

	for i := 0; i < NumDepositsLoad; i++ {
		// Create a deposit transaction.
		_, err := s.generateNewDepositTx(
			s.GenesisAccount().Address(),
			s.GenesisAccount().SignerFunc(chainID),
			big.NewInt(int64(nonce+uint64(i))),
		)
		s.Require().NoError(err)
	}

	// wait blocks
	targetBlkNum := blkNum + 10
	err = s.WaitForFinalizedBlockNumber(targetBlkNum)
	s.Require().NoError(err)

	// Check to see if evm balance decreased.
	postDepositBalance, err := s.JSONRPCBalancer().BalanceAt(
		s.Ctx(),
		s.GenesisAccount().Address(),
		big.NewInt(int64(targetBlkNum)),
	)
	s.Require().NoError(err)
	s.Require().Equal(postDepositBalance.Cmp(balance), -1)

	fmt.Println("Balance before deposit: ", balance)
	fmt.Println("Balance after deposit: ", postDepositBalance)

	// Chekc that the balance is somewhere between og - 32e * 500 > 0 < 32e
	// TODO: revive node

	// // Wait for some txs to be processed.

	// // Check to make sure the balance has decreased by the correct amount.
	// newBalance, err := s.JSONRPCBalancer().BalanceAt(
	// 	s.Ctx()
	// 	s.GenesisAccount().Address(),
	// 	big.NewInt(int64(blkNum)),
	// )
	// s.Require().NoError(err)

	// totalAmountDeposited := new(big.Int).Mul(big.NewInt(32*suite.OneGwei), big.NewInt(NumDepositsLoad))

	// expectedBalance := new(big.Int).Sub(balance, totalAmountDeposited)

	// // Wait for the log to be processed.
	// // Check that the total power is less than total amount by X time.
	// // Check that the total power adds up to the total amount by X time.
}

func (s *BeaconKitE2ESuite) generateNewDepositTx(
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

	// Bind the deposit contract.
	dc, err := deposit.NewBeaconDepositContract(
		common.HexToAddress(DepositContractAddress),
		s.JSONRPCBalancer(),
	)
	s.Require().NoError(err)

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
