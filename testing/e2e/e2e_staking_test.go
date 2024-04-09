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

	byteslib "github.com/berachain/beacon-kit/mod/primitives/bytes"
	stakingabi "github.com/berachain/beacon-kit/mod/runtime/services/staking/abi"
	"github.com/berachain/beacon-kit/testing/e2e/suite"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	coretypes "github.com/ethereum/go-ethereum/core/types"
)

const (
	// DepositContractAddress is the address of the deposit contract.
	DepositContractAddress = "0x00000000219ab540356cbb839cbe05303d7705fa"
)

// TestDepositContract tests the deposit contract to attempt staking and
// increasing a validator's consensus power.
func (s *BeaconKitE2ESuite) TestDepositContract() {
	// Get the consensus client.
	client := s.ConsensusClients()["cl-validator-beaconkit-0"]
	s.Require().NotNil(client)

	// Get the public key.
	pubkey, err := client.GetPubKey(s.Ctx())
	s.Require().NoError(err)
	s.Require().Len(pubkey, 48)

	// Get the consensus power.
	_, err = client.GetConsensusPower(s.Ctx())
	s.Require().NoError(err)

	// Bind the deposit contract.
	dc, err := stakingabi.NewBeaconDepositContract(
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

	// Get the chain ID.
	chainID, err := s.JSONRPCBalancer().ChainID(s.Ctx())
	s.Require().NoError(err)

	// Get the block num
	blkNum, err := s.JSONRPCBalancer().BlockNumber(s.Ctx())
	s.Require().NoError(err)

	// Get original evm balance
	balance, err := s.JSONRPCBalancer().BalanceAt(
		s.Ctx(),
		s.GenesisAccount().Address(),
		big.NewInt(int64(blkNum)),
	)
	s.Require().NoError(err)

	// Create a deposit transaction.
	val, _ := big.NewFloat(32e18).Int(nil)
	tx, err := dc.Deposit(&bind.TransactOpts{
		From:   s.GenesisAccount().Address(),
		Value:  val,
		Signer: s.GenesisAccount().SignerFunc(chainID),
	}, pubkey, credentials, 32*suite.OneGwei, signature[:])
	s.Require().NoError(err)

	// Wait for the transaction to be mined.
	var receipt *coretypes.Receipt
	receipt, err = bind.WaitMined(s.Ctx(), s.JSONRPCBalancer(), tx)
	s.Require().NoError(err)
	s.Require().Equal(uint64(1), receipt.Status)
	s.Logger().Info("Deposit transaction mined", "txHash", receipt.TxHash.Hex())

	// Wait for the log to be processed.
	targetBlkNum := blkNum + 5
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

	// // Submit withdrawal
	// tx, err = dc.Withdraw(&bind.TransactOpts{
	// 	From:   s.GenesisAccount().Address(),
	// 	Signer: s.GenesisAccount().SignerFunc(chainID),
	// }, pubkey, credentials, 31*suite.OneGwei)
	// s.Require().NoError(err)

	// receipt, err = bind.WaitMined(s.Ctx(), s.JSONRPCBalancer(), tx)
	// s.Require().NoError(err)
	// s.Require().Equal(uint64(1), receipt.Status)
	// s.Logger().
	// 	Info("Withdraw transaction mined", "txHash", receipt.TxHash.Hex())

	// // Wait for the log to be processed.
	// targetBlkNum += 4
	// err = s.WaitForFinalizedBlockNumber(targetBlkNum)
	// s.Require().NoError(err)

	// // Check to see if new balance is greater than the previous balance
	// postWithdrawBalance, err := s.JSONRPCBalancer().BalanceAt(
	// 	s.Ctx(),
	// 	s.GenesisAccount().Address(),
	// 	big.NewInt(int64(targetBlkNum)),
	// )
	// s.Require().NoError(err)
	// s.Require().Equal(postWithdrawBalance.Cmp(postDepositBalance), 1)

	// // We are withdrawing all the power, so the power should be 0.
	// postWithdrawPower, err := client.GetConsensusPower(s.Ctx())
	// s.Require().NoError(err)
	// s.Require().Equal(postWithdrawPower, suite.OneGwei)
}

// TODO: once RPC ready
// func (s *BeaconKitE2ESuite) TestCreateNewValidator() {
// 	var nut *types.ConsensusClient
// 	for _, cl := range s.ConsensusClients() {
// 		if yes, err := cl.IsActive(s.Ctx()); err != nil && !yes {
// 			nut = cl
// 			break
// 		}
// 	}

// 	// Generate the credentials.
// 	credentials := byteslib.PrependExtendToSize(
// 		s.GenesisAccount().Address().Bytes(),
// 		32,
// 	)
// 	credentials[0] = 0x01

// 	// Bind the deposit contract.
// 	dc, err := stakingabi.NewBeaconDepositContract(
// 		common.HexToAddress(DepositContractAddress),
// 		s.JSONRPCBalancer(),
// 	)
// 	s.Require().NoError(err)

// 	// Get the chain ID.
// 	chainID, err := s.JSONRPCBalancer().ChainID(s.Ctx())
// 	s.Require().NoError(err)

// 	// Get the block num
// 	blkNum, err := s.JSONRPCBalancer().BlockNumber(s.Ctx())
// 	s.Require().NoError(err)

// 	nodePubkey, err := nut.GetPubKey(s.Ctx())
// 	if err != nil {
// 		s.Require().NoError(err)
// 	}

// 	// TODO: get private key from the node
// 	privateKey := [32]byte{}
// 	signer, err := blst.SecretKeyFromBytes(privateKey[:])
// 	s.Require().NoError(err)

// 	msg := beacontypes.DepositMessage{
// 		Pubkey:      primitives.BLSPubkey(nodePubkey),
// 		Credentials: beacontypes.WithdrawalCredentials(credentials),
// 		Amount:      primitives.Gwei(32 * suite.OneGwei),
// 	}

// 	// forkData := forks.NewForkData(

// 	// )

// 	domain, err := forkData.ComputeDomain(primitives.DomainTypeDeposit)
// 	s.Require().NoError(err)

// 	signingRoot, err := primitives.ComputeSigningRoot(&msg, domain)
// 	s.Require().NoError(err)

// 	// Sign the message.
// 	sig := signer.Sign(signingRoot[:]).Marshal()

// 	// Create a deposit transaction.
// 	val, _ := big.NewFloat(32e18).Int(nil)
// 	tx, err := dc.Deposit(&bind.TransactOpts{
// 		From:   s.GenesisAccount().Address(),
// 		Value:  val,
// 		Signer: s.GenesisAccount().SignerFunc(chainID),
// 	}, nodePubkey[:], credentials, 32*suite.OneGwei, sig)
// 	s.Require().NoError(err)

// 	// Wait for the transaction to be mined.
// 	var receipt *coretypes.Receipt
// 	receipt, err = bind.WaitMined(s.Ctx(), s.JSONRPCBalancer(), tx)
// 	s.Require().NoError(err)
// 	s.Require().Equal(uint64(1), receipt.Status)
// 	s.Logger().Info("Deposit transaction mined", "txHash", receipt.TxHash.Hex())

// 	// Wait for the log to be processed.
// 	targetBlkNum := blkNum + 5
// 	err = s.WaitForFinalizedBlockNumber(targetBlkNum)
// 	s.Require().NoError(err)

// 	active, err := nut.IsActive(s.Ctx())
// 	s.Require().NoError(err)
// 	s.Require().True(active)
// }
