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

	stakingabi "github.com/berachain/beacon-kit/contracts/abi"
	byteslib "github.com/berachain/beacon-kit/lib/bytes"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	coretypes "github.com/ethereum/go-ethereum/core/types"
)

// TestForgeScriptExecution tests the execution of a forge script
// against the beacon-kit network.
func (s *BeaconKitE2ESuite) TestDepositContract() {
	client := s.ConsensusClients()["cl-validator-beaconkit-0"]
	s.Require().NotNil(client)

	pubkey, err := client.GetPubKey(s.Ctx())
	s.Require().NoError(err)

	_, err = client.GetConsensusPower(s.Ctx())
	s.Require().NoError(err)

	dc, err := stakingabi.NewBeaconDepositContract(
		common.HexToAddress("0x00000000219ab540356cbb839cbe05303d7705fa"),
		s.JSONRPCBalancer(),
	)
	s.Require().NoError(err)

	bz := byteslib.PrependExtendToSize(s.GenesisAccount().Address().Bytes(), 32)
	bz[0] = 0x01

	val, _ := big.NewFloat(32e18).Int(nil)
	tx, err := dc.Deposit(&bind.TransactOpts{
		From:  s.GenesisAccount().Address(),
		Value: val,
	}, pubkey, bz, 32e9, nil)
	s.Require().NoError(err)

	chainID, err := s.JSONRPCBalancer().ChainID(s.Ctx())
	s.Require().NoError(err)
	tx, err = s.GenesisAccount().SignTx(chainID, tx)
	s.Require().NoError(err)

	err = s.JSONRPCBalancer().SendTransaction(s.Ctx(), tx)
	s.Require().NoError(err)

	var receipt *coretypes.Receipt
	receipt, err = bind.WaitMined(s.Ctx(), s.JSONRPCBalancer(), tx)
	s.Require().NoError(err)
	s.Require().Equal(uint64(1), receipt.Status)
}
