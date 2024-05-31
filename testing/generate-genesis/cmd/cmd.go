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

package cmd

import (
	"math/big"
	"strconv"

	"github.com/berachain/beacon-kit/testing/generate-genesis/genesis"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

// CreateEthGenesisCmd creates a cobra command for generating a genesis.json file.
func CreateEthGenesisCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-genesis <format>",
		Short: "Generate a genesis.json file",
		Args:  cobra.ExactArgs(1),
		RunE:  createEthGenesisCmdFunc,
	}

	cmd.Flags().StringSliceP(
		accountAddressesFlag, "a",
		[]string{},
		"Account addresses to add",
	)
	cmd.Flags().StringSliceP(
		accountBalancesFlag, "b",
		[]string{},
		"Account balances to add",
	)
	cmd.Flags().StringSliceP(
		predeployAddressesFlag, "p",
		[]string{},
		"Predeploy contract addresses to add",
	)
	cmd.Flags().StringSliceP(
		predeployCodesFlag, "c",
		[]string{},
		"Predeploy contract codes to add",
	)
	cmd.Flags().StringSliceP(
		predeployBalancesFlag, "i",
		[]string{},
		"Predeploy contract balances to add",
	)
	cmd.Flags().StringSliceP(
		predeployNoncesFlag, "n",
		[]string{},
		"Predeploy contract nonces to add",
	)
	cmd.Flags().StringP(
		outputFileFlag, "o",
		"eth-genesis.json",
		"Output file name for the genesis file",
	)

	return cmd
}

func createEthGenesisCmdFunc(cmd *cobra.Command, args []string) error {
	var gen genesis.Genesis
	switch args[0] {
	case "geth":
		gen = genesis.NewGeth()
	case "nethermind":
		gen = genesis.NewNethermind()
	default:
		return errInvalidEthGenesisFormat
	}

	accountAddresses, accountBalances, predeployAddresses, predeployCodes, predeployBalances, predeployNonces, outputFile, err := sanitizeFlags(cmd)
	if err != nil {
		return err
	}

	for index, address := range accountAddresses {
		balance, ok := new(big.Int).SetString(accountBalances[index], 10)
		if !ok {
			return errInvalidAccountBalance
		}
		if err := gen.AddAccount(address, balance); err != nil {
			return err
		}
	}

	for i := range predeployAddresses {
		code := common.FromHex(predeployCodes[i])
		balance := new(big.Int)
		balance.SetString(predeployBalances[i], 10)

		nonce, err := strconv.ParseUint(predeployNonces[i], 10, 64) // convert string to uint64
		if err != nil {
			return err
		}

		if err = gen.AddPredeploy(predeployAddresses[i], code, balance, nonce); err != nil {
			return err
		}
	}

	return gen.WriteJSON(outputFile)
}
