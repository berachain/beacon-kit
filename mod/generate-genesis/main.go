package main

import (
	"github.com/berachain/beacon-kit/mod/generate-genesis/genesis"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/spf13/cobra"
	"math/big"
	"strconv"
)

// CreateGenesisFileCmd Commands creates a new genesis file
func CreateGenesisFileCmd() *cobra.Command {
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
			g := genesis.NewGenesis()
			var gen interface{}
			switch genesisFormat {
			case "geth":
				gen = g.ToGethGenesis()
			case "nethermind":
				gen = g.ToNethermindGenesis()
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
				switch genesisFormat {
				case "geth":
					g.AddAccount(gen.(*core.Genesis), common.HexToAddress(address), balanceBigInt)
				case "nethermind":
					g.AddAccountNethermind(gen.(*genesis.NethermindGenesis), common.HexToAddress(address), balanceBigInt)
				}
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

				switch genesisFormat {
				case "geth":
					g.AddPredeploy(gen.(*core.Genesis), predeployAddress, predeployCode, nil, balance, nonce)
				case "nethermind":
					g.AddPredeployNethermind(gen.(*genesis.NethermindGenesis), predeployAddress, predeployCode, balance, nonce)
				}

			}
			switch genesisFormat {
			case "geth":
				g.WriteFileToJSON(gen.(*core.Genesis), outputFileName)
			case "nethermind":
				g.WriteNethermindGenesisToJSON(gen.(*genesis.NethermindGenesis), outputFileName)
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
	err := CreateGenesisFileCmd().Execute()
	if err != nil {
		return
	}
}
