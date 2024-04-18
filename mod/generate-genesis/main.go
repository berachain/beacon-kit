package main

import (
	"github.com/berachain/beacon-kit/mod/generate-genesis/genesis"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/cobra"
	"math/big"
	"strconv"
)

// createGenesisCmd creates a new genesis file based on the file format passed to it.
func createGenesisCmd() *cobra.Command {
	var accountAddresses []string
	var accountBalances []string
	var predeployAddresses []string
	var predeployCodes []string
	var predeployBalances []string
	var predeployNonces []string
	var genesisFormat string
	var outputFileName string

	cmd := &cobra.Command{
		Use:   "generate-genesis",
		Short: "Generate a genesis.json file",
		Run: func(cmd *cobra.Command, args []string) {
			var gen genesis.Genesis
			switch genesisFormat {
			case "geth":

				gethGenesis := &genesis.GethGenesis{
					Alloc: make(types.GenesisAlloc),
				}
				gen = gethGenesis
				gethGenesis.CoreGenesis = gethGenesis.ToGethGenesis().CoreGenesis

			case "nethermind":
				nethermindGenesis := &genesis.NethermindGenesis{}
				nethermindGenesis = nethermindGenesis.ToNethermindGenesis()
				gen = nethermindGenesis

			default:
				cmd.PrintErrf("invalid genesis format %v\n", genesisFormat)
				return
			}

			for i, address := range accountAddresses {
				balance := accountBalances[i]
				balanceBigInt, success := new(big.Int).SetString(balance, 10)
				if !success {
					cmd.PrintErrf("failed to convert balance to big.Int %v\n", balance)
				}
				gen.AddAccount(common.HexToAddress(address), balanceBigInt)
			}

			for i := range predeployAddresses {
				predeployAddress := common.HexToAddress(predeployAddresses[i])
				predeployCode := common.FromHex(predeployCodes[i])
				balance := new(big.Int)
				balance.SetString(predeployBalances[i], 10)

				nonce, err := strconv.ParseUint(predeployNonces[i], 10, 64) // convert string to uint64
				if err != nil {
					cmd.PrintErrf("failed to nonce to uint64 %v\n", err)
				}

				gen.AddPredeploy(predeployAddress, predeployCode, balance, nonce)
			}

			err := gen.ToJSON(outputFileName)
			if err != nil {
				cmd.PrintErrf("failed to write file %v\n", err)
				return
			}
		},
	}

	cmd.Flags().StringSliceVarP(&accountAddresses, "account", "a", []string{}, "Account address to add")
	cmd.Flags().StringSliceVarP(&accountBalances, "balance", "b", []string{}, "Balance for the account")
	cmd.Flags().StringSliceVarP(&predeployAddresses, "predeployAddress", "p", []string{}, "Predeploy contract addresses")
	cmd.Flags().StringSliceVarP(&predeployCodes, "code", "c", []string{}, "Codes for the predeploy contract")
	cmd.Flags().StringSliceVarP(&predeployBalances, "predeploybalance", "i", []string{}, "Balances for the predeploy contract")
	cmd.Flags().StringSliceVarP(&predeployNonces, "nonce", "n", []string{}, "Nonces for the predeploy contracts")
	cmd.Flags().StringVarP(&genesisFormat, "format", "f", "geth", "Format of the genesis file (geth or nethermind)")
	cmd.Flags().StringVarP(&outputFileName, "output", "o", "genesis-go-ethereum.json", "Output file name for the genesis file")

	return cmd
}

func main() {
	err := createGenesisCmd().Execute()
	if err != nil {
		return
	}
}
