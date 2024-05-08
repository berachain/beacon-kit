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

package main

//
//import (
//	"github.com/berachain/beacon-kit/mod/generate-genesis/genesis"
//	"github.com/ethereum/go-ethereum/common"
//	"github.com/spf13/cobra"
//	"math/big"
//	"strconv"
//)
//
//// CreateGenesisFileCmd Commands creates a new genesis file
//func CreateGenesisFileCmd2() *cobra.Command {
//	var accountAddresses []string
//	var accountBalances []string
//	var predeployAddresses []string
//	var predeployCodes []string
//	var predeployBalances []string
//	var predeployNonces []string
//	var genesisFormat string
//	var outputFileName string
//
//	cmd := &cobra.Command{
//		Use:   "generate-genesis",
//		Short: "Generate a genesis.json file",
//		Run: func(cmd *cobra.Command, args []string) {
//			g := genesis.NewGenesis()
//			gen := g.ToGethGenesis()
//
//			for i, address := range accountAddresses {
//				balanceStr := accountBalances[i]
//				balanceBigInt, success := new(big.Int).SetString(balanceStr, 10)
//				if !success {
//					panic("Failed to convert balance to big.Int")
//				}
//				g.AddAccount(gen, common.HexToAddress(address), balanceBigInt)
//			}
//
//			for i := range predeployAddresses {
//				predeployAddressStr := common.HexToAddress(predeployAddresses[i])
//				predeployCodeStr := common.FromHex(predeployCodes[i])
//				balance := new(big.Int)
//				balance.SetString(predeployBalances[i], 10)
//
//				nonce, err := strconv.ParseUint(predeployNonces[i], 10, 64) // convert string to uint64
//				if err != nil {
//					panic("Failed to convert nonce to uint64")
//				}
//
//				g.AddPredeploy(gen, predeployAddressStr, predeployCodeStr, nil, balance, nonce)
//			}
//			g.WriteFileToJSON(gen, outputFileName)
//		},
//	}
//
//	cmd.Flags().StringSliceVarP(&accountAddresses, "account", "a", []string{}, "Account address to add")
//	cmd.Flags().StringSliceVarP(&accountBalances, "balance", "b", []string{}, "Balance for the account")
//	cmd.Flags().StringSliceVarP(&predeployAddresses, "predeploy", "p", []string{}, "Predeploy contract addresses")
//	cmd.Flags().StringSliceVarP(&predeployCodes, "code", "c", []string{}, "Codes for the predeploy contract")
//	cmd.Flags().StringSliceVarP(&predeployBalances, "predeploybalance", "i", []string{}, "Balances for the predeploy contract")
//	cmd.Flags().StringSliceVarP(&predeployNonces, "nonce", "n", []string{}, "Nonces for the predeploy contracts")
//	cmd.Flags().StringVarP(&genesisFormat, "format", "f", "geth", "Format of the genesis file (geth or nethermind)")
//	cmd.Flags().StringVarP(&outputFileName, "output", "o", "genesis.json", "Output file name for the genesis file")
//
//	return cmd
//}
//
////func AddAccountCmd(g *genesis.Genesis) *cobra.Command {
////	var accountAddress string
////	var accountBalance string
////
////	addAccountCmd := &cobra.Command{
////		Use:   "addAccount",
////		Short: "Add an account to the genesis file",
////		Run: func(cmd *cobra.Command, args []string) {
////			//g := genesis.NewGenesis()
////
////			address := common.HexToAddress(accountAddress)
////			balance := new(big.Int)
////			balance.SetString(accountBalance, 10)
////
////			g.AddAccount(address, balance)
////		},
////	}
////
////	addAccountCmd.Flags().StringVarP(&accountAddress, "address", "a", "", "Account address to add")
////	addAccountCmd.Flags().StringVarP(&accountBalance, "balance", "b", "", "Balance for the account")
////
////	return addAccountCmd
////}
//
////func AddPreDeployCmd(g *genesis.Genesis) *cobra.Command {
////	var predeployAddress string
////	var predeployCode string
////	var predeployBalance string
////	var predeployNonce uint64
////
////	addPredeployCmd := &cobra.Command{
////		Use:   "addPredeploy",
////		Short: "Add a predeploy to the genesis file",
////		Run: func(cmd *cobra.Command, args []string) {
////			//g := genesis.NewGenesis()
////
////			address := common.HexToAddress(predeployAddress)
////			code := common.FromHex(predeployCode)
////			balance := new(big.Int)
////			balance.SetString(predeployBalance, 10)
////
////			g.AddPredeploy(address, code, nil, balance, predeployNonce)
////		},
////	}
////
////	addPredeployCmd.Flags().StringVarP(&predeployAddress, "address", "a", "", "Predeploy contract address")
////	addPredeployCmd.Flags().StringVarP(&predeployCode, "code", "c", "", "Code for the predeploy contract")
////	addPredeployCmd.Flags().StringVarP(&predeployBalance, "balance", "b", "", "Balance for the predeploy contract")
////	addPredeployCmd.Flags().Uint64VarP(&predeployNonce, "nonce", "n", 0, "Nonce for the predeploy contract")
////
////	return addPredeployCmd
////}
//
//func main() {
//
//	//g := genesis.NewGenesis()
//	//
//	//rootCmd := &cobra.Command{
//	//	Use:   "generate-genesis",
//	//	Short: "Generate a genesis.json file",
//	//}
//	//
//	////cmd := CreateGenesisFileCmd(g)
//	//rootCmd.AddCommand(AddAccountCmd(g), AddPreDeployCmd(g))
//	//
//	//if err := rootCmd.Execute(); err != nil {
//	//	panic(err)
//	//}
//	//
//
//	err := CreateGenesisFileCmd().Execute()
//	if err != nil {
//		return
//	}
//
//	// Create a new genesis
//	//genesis := genesis.NewGenesis()
//
//	//http.HandleFunc("/addAccount", addAccountHandler)
//	//http.HandleFunc("/addPredeploy", addPredeployHandler)
//	//log.Fatal(http.ListenAndServe(":8080", nil))
//
//	//genesis.AddAccount(common.HexToAddress("0x20f33ce90a13a4b5e7697e3544c3083b8f8a51d4"), balance)                                                                                                                                                                                                                                             // Add an account with balance
//	//genesis.AddPredeploy(common.HexToAddress("0x000F3df6D732807Ef1319fB7B8bB8522d0Beac02"), common.FromHex("0x3373fffffffffffffffffffffffffffffffffffffffe14604d57602036146024575f5ffd5b5f35801560495762001fff810690815414603c575f5ffd5b62001fff01545f5260205ff35b5f5ffd5b62001fff42064281555f359062001fff015500"), nil, big.NewInt(0x0), 0x1) // Add a predeploy
//	//genesisConfig := genesis.ToGethGenesis()                                                                                                                                                                                                                                                                                                   // Convert the genesis to a Geth genesis
//	//genesis.ConvertToJSON(genesisConfig)
//
//}
//
////func addAccountHandler(w http.ResponseWriter, r *http.Request) {
////	// Parse the request parameters
////	address := common.HexToAddress(r.URL.Query().Get("address"))
////	balance, _ := new(big.Int).SetString(r.URL.Query().Get("balance"), 10)
////
////	// Create a new Genesis instance
////	genesis := genesis.NewGenesis()
////
////	// Add the account
////	//genesis.AddAccount1(address, balance)
////
////	// Respond with a success message
////	fmt.Fprintf(w, string(genesis.ConvertToJSON(genesis.ToGethGenesis())))
////}
//
////func addPredeployHandler(w http.ResponseWriter, r *http.Request) {
////	// Parse the request parameters
////	address := common.HexToAddress(r.URL.Query().Get("address"))
////	code := common.FromHex(r.URL.Query().Get("code"))
////	balance, _ := new(big.Int).SetString(r.URL.Query().Get("balance"), 10)
////	nonce, _ := strconv.ParseUint(r.URL.Query().Get("nonce"), 10, 64)
////
////	// Create a new Genesis instance
////	genesis := genesis.NewGenesis()
////
////	// Add the predeploy
////	//genesis.AddPredeploy(address, code, nil, balance, nonce)
////
////	// Respond with a success message
////	fmt.Fprintf(w, string(genesis.ConvertToJSON(genesis.ToGethGenesis())))
////}
