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
	"crypto/sha256"
	"encoding/hex"
	"math/big"

	consensustypes "github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/execution/pkg/deposit"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/commands/utils/parser"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/components/signer"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/config/spec"
	byteslib "github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
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
	// TODO: WE ARE NOT MERGING THIS PR UNTIL ALL OF THIS IS GOOD AGAIN.
	s.T().Skip("Placeholder")
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
	dc, err := deposit.NewBeaconDepositContract(
		common.HexToAddress(DepositContractAddress),
		s.JSONRPCBalancer(),
	)
	s.Require().NoError(err)

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

	// The deposit amount
	val, _ := big.NewFloat(32e18).Int(nil)

	// Generate the credentials.
	credentials := byteslib.PrependExtendToSize(
		s.GenesisAccount().Address().Bytes(),
		32,
	)
	credentials[0] = 0x01

	// Generate the signature.
	validatorPrivKey := ""

	validatorPrivKeyBz, err := hex.DecodeString(validatorPrivKey)
	s.Require().NoError(err)
	s.Require().Equal(len(validatorPrivKeyBz), constants.BLSSecretKeyLength)
	blsSigner, err := signer.NewBLSSigner(
		[constants.BLSSecretKeyLength]byte(validatorPrivKeyBz),
	)
	s.Require().NoError(err)

	gweiU64, err := parser.ConvertAmount(val.String())
	s.Require().NoError(err)

	// TODO: fill in the fork data -- get this from RPC when ready
	forkData := &consensustypes.ForkData{
		CurrentVersion:        [4]byte{},
		GenesisValidatorsRoot: [32]byte{},
	}

	_, signature, err := consensustypes.CreateAndSignDepositMessage(
		forkData,
		spec.LocalnetChainSpec().DomainTypeDeposit(),
		blsSigner,
		consensustypes.WithdrawalCredentials(credentials),
		gweiU64,
	)
	s.Require().NoError(err)

	// Compute the deposit root.
	// Convert the amount from big.Int to U256 so that we can ssz it.
	u256, err := math.NewU256LFromBigInt(val)
	s.Require().NoError(err)
	amount, err := u256.MarshalSSZ()
	s.Require().NoError(err)

	pubkeyRoot := sha256.Sum256(append(pubkey, make([]byte, 16)...))
	signatureRoot := sha256.Sum256(append(signature[:64], make([]byte, 16)...))

	part1 := sha256.Sum256(append(pubkeyRoot[:], credentials...))
	part2 := sha256.Sum256(append(amount, signatureRoot[:]...))
	depositRoot := sha256.Sum256(append(part1[:], part2[:]...))

	// Create a deposit transaction.
	tx, err := dc.Deposit(&bind.TransactOpts{
		From:   s.GenesisAccount().Address(),
		Value:  val,
		Signer: s.GenesisAccount().SignerFunc(chainID),
	}, pubkey, credentials, signature[:], depositRoot)
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
// }

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
// 		Pubkey:      crypto.BLSSignaturePubkey(nodePubkey),
// 		Credentials: beacontypes.WithdrawalCredentials(credentials),
// 		Amount:      math.Gwei(32 * suite.OneGwei),
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
