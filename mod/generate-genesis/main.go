package main

import (
	"github.com/berachain/beacon-kit/mod/generate-genesis/genesis"
	"github.com/cockroachdb/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/cobra"
	"log"
	"math/big"
	"strconv"
)

// createGenesis creates a genesis configuration struct based on the provided genesis format.
func createGenesis(genesisFormat string) (genesis.Genesis, error) {
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
		return nil, errors.New("invalid genesis format: " + genesisFormat)
	}
	return gen, nil
}

func addAccount(gen genesis.Genesis, address string, balance string) error {
	balanceBigInt, success := new(big.Int).SetString(balance, 10)
	if !success {
		return errors.Wrapf(errors.New("failed to convert balance to big.Int"), "balance: %s", balance)

	}
	gen.AddAccount(common.HexToAddress(address), balanceBigInt)
	return nil
}

func addPredeploy(gen genesis.Genesis, predeployAddress string, predeployCode string, predeployBalance string, predeployNonce string) error {
	address := common.HexToAddress(predeployAddress)
	code := common.FromHex(predeployCode)
	balance := new(big.Int)
	balance.SetString(predeployBalance, 10)

	nonce, err := strconv.ParseUint(predeployNonce, 10, 64) // convert string to uint64
	if err != nil {
		return errors.Wrap(err, "failed to convert nonce to uint64")
	}

	gen.AddPredeploy(address, code, balance, nonce)
	return nil
}

func writeGenesis(gen genesis.Genesis, outputFileName string) error {
	err := gen.ToJSON(outputFileName)
	if err != nil {
		return errors.Wrap(err, "failed to write genesis to a file")
	}
	return nil
}

// createGenesisCmd creates a cobra command for generating a genesis.json file.
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
		RunE: func(cmd *cobra.Command, args []string) error {
			gen, err := createGenesis(genesisFormat)
			if err != nil {
				return err

			}

			for index, address := range accountAddresses {
				balance := accountBalances[index]
				err := addAccount(gen, address, balance)
				if err != nil {
					return err
				}
			}

			for counter := range predeployAddresses {
				predeployAddress := predeployAddresses[counter]
				predeployCode := predeployCodes[counter]
				predeployBalance := predeployBalances[counter]
				predeployNonce := predeployNonces[counter]
				err := addPredeploy(gen, predeployAddress, predeployCode, predeployBalance, predeployNonce)
				if err != nil {
					return err

				}
			}

			err = writeGenesis(gen, outputFileName)
			if err != nil {
				return err

			}
			return nil
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
		log.Printf("Error: %v\n", err)
	}
}
