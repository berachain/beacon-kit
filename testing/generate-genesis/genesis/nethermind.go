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

package genesis

import (
	"math/big"
	"strconv"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/hex"
)

var _ Genesis = (*Nethermind)(nil)

// Nethermind is a struct that holds the genesis configuration for Nethermind.
type Nethermind struct {
	Name     string             `json:"name"`
	DataDir  string             `json:"dataDir"`
	Engine   Engine             `json:"engine"`
	Params   Params             `json:"params"`
	Genesis  NetherGenesis      `json:"genesis"`
	Accounts map[string]Account `json:"accounts"`
}

type Engine struct {
	Ethash Ethash `json:"Ethash"`
}

type Ethash struct {
	Params ParamsEthash `json:"params"`
}

type ParamsEthash struct {
	MinimumDifficulty      string   `json:"minimumDifficulty"`
	DifficultyBoundDivisor string   `json:"difficultyBoundDivisor"`
	DurationLimit          string   `json:"durationLimit"`
	BlockReward            struct{} `json:"blockReward"`
	HomesteadTransition    string   `json:"homesteadTransition"`
	DaoHardforkTransition  string   `json:"daoHardforkTransition"`
	Eip100bTransition      string   `json:"eip100bTransition"`
	DifficultyBombDelays   struct{} `json:"difficultyBombDelays"`
}

type Params struct {
	GasLimitBoundDivisor          string `json:"gasLimitBoundDivisor"`
	Registrar                     string `json:"registrar"`
	AccountStartNonce             string `json:"accountStartNonce"`
	MaximumExtraDataSize          string `json:"maximumExtraDataSize"`
	MinGasLimit                   string `json:"minGasLimit"`
	NetworkID                     string `json:"networkID"`
	ForkBlock                     string `json:"forkBlock"`
	MaxCodeSize                   string `json:"maxCodeSize"`
	MaxCodeSizeTransition         string `json:"maxCodeSizeTransition"`
	Eip150Transition              string `json:"eip150Transition"`
	Eip160Transition              string `json:"eip160Transition"`
	Eip161abcTransition           string `json:"eip161abcTransition"`
	Eip161dTransition             string `json:"eip161dTransition"`
	Eip155Transition              string `json:"eip155Transition"`
	Eip140Transition              string `json:"eip140Transition"`
	Eip211Transition              string `json:"eip211Transition"`
	Eip214Transition              string `json:"eip214Transition"`
	Eip658Transition              string `json:"eip658Transition"`
	Eip145Transition              string `json:"eip145Transition"`
	Eip1014Transition             string `json:"eip1014Transition"`
	Eip1052Transition             string `json:"eip1052Transition"`
	Eip1283Transition             string `json:"eip1283Transition"`
	Eip1283DisableTransition      string `json:"eip1283DisableTransition"`
	Eip152Transition              string `json:"eip152Transition"`
	Eip1108Transition             string `json:"eip1108Transition"`
	Eip1344Transition             string `json:"eip1344Transition"`
	Eip1884Transition             string `json:"eip1884Transition"`
	Eip2028Transition             string `json:"eip2028Transition"`
	Eip2200Transition             string `json:"eip2200Transition"`
	Eip2565Transition             string `json:"eip2565Transition"`
	Eip2929Transition             string `json:"eip2929Transition"`
	Eip2930Transition             string `json:"eip2930Transition"`
	Eip1559Transition             string `json:"eip1559Transition"`
	Eip3198Transition             string `json:"eip3198Transition"`
	Eip3529Transition             string `json:"eip3529Transition"`
	Eip3541Transition             string `json:"eip3541Transition"`
	Eip4895TransitionTimestamp    string `json:"eip4895TransitionTimestamp"`
	Eip3855TransitionTimestamp    string `json:"eip3855TransitionTimestamp"`
	Eip3651TransitionTimestamp    string `json:"eip3651TransitionTimestamp"`
	Eip3860TransitionTimestamp    string `json:"eip3860TransitionTimestamp"`
	Eip1153TransitionTimestamp    string `json:"eip1153TransitionTimestamp"`
	Eip4788TransitionTimestamp    string `json:"eip4788TransitionTimestamp"`
	Eip4844TransitionTimestamp    string `json:"eip4844TransitionTimestamp"`
	Eip5656TransitionTimestamp    string `json:"eip5656TransitionTimestamp"`
	Eip6780TransitionTimestamp    string `json:"eip6780TransitionTimestamp"`
	TerminalTotalDifficulty       string `json:"terminalTotalDifficulty"`
	TerminalTotalDifficultyPassed bool   `json:"terminalTotalDifficultyPassed"`
}

type NetherGenesis struct {
	Coinbase   string `json:"coinbase"`
	Difficulty string `json:"difficulty"`
	ExtraData  string `json:"extraData"`
	GasLimit   string `json:"gasLimit"`
	Nonce      string `json:"nonce"`
	Timestamp  string `json:"timestamp"`
}

type Account struct {
	Balance string `json:"balance"`
	Nonce   string `json:"nonce"`
	Code    string `json:"code"`
}

func NewNethermind() *Nethermind {
	return &Nethermind{
		Name:    "Ethereum",
		DataDir: "ethereum",
		Engine: Engine{
			Ethash: Ethash{
				Params: ParamsEthash{
					MinimumDifficulty:      "0x0",
					DifficultyBoundDivisor: "0x0",
					DurationLimit:          "0x0",
				},
			},
		},
		Params: defaultNethermindParams(),
		Genesis: NetherGenesis{
			Coinbase:   "0x0000000000000000000000000000000000000000",
			Difficulty: "0x0",
			ExtraData:  DefaultExtraData,
			GasLimit:   "0x1c9c380",
			Nonce:      "0x0000000000000000",
			Timestamp:  "0x0",
		},
		Accounts: make(map[string]Account),
	}
}

func (n *Nethermind) AddAccount(address string, balance *big.Int) error {
	if _, ok := n.Accounts[address]; ok {
		return errAccountAlreadyExists
	}

	n.Accounts[address] = Account{
		Balance: "0x" + balance.Text(HexBase), // Convert balance to hexadecimal
	}
	return nil
}

func (n *Nethermind) AddPredeploy(
	address string,
	code []byte,
	balance *big.Int,
	nonce uint64) error {
	if _, ok := n.Accounts[address]; ok {
		return errPredeployAlreadyExists
	}

	n.Accounts[address] = Account{
		// convert to hexadecimal
		Balance: "0x" + balance.Text(HexBase),
		Nonce:   "0x" + strconv.FormatUint(nonce, HexBase),
		Code:    "0x" + hex.FromBytes(code).Unwrap(),
	}

	return nil
}

func (n *Nethermind) WriteJSON(filename string) error {
	return writeGenesisToJSON(n, filename)
}

func defaultNethermindParams() Params {
	return Params{
		GasLimitBoundDivisor:          "0x400",
		Registrar:                     "0xe3389675d0338462dC76C6f9A3e432550c36A142",
		AccountStartNonce:             "0x0",
		MaximumExtraDataSize:          "0x20",
		MinGasLimit:                   "0x1c9c380",
		NetworkID:                     "0x138d7",
		ForkBlock:                     "0x0",
		MaxCodeSize:                   "0x6000",
		MaxCodeSizeTransition:         "0x0",
		Eip150Transition:              "0x0",
		Eip160Transition:              "0x0",
		Eip161abcTransition:           "0x0",
		Eip161dTransition:             "0x0",
		Eip155Transition:              "0x0",
		Eip140Transition:              "0x0",
		Eip211Transition:              "0x0",
		Eip214Transition:              "0x0",
		Eip658Transition:              "0x0",
		Eip145Transition:              "0x0",
		Eip1014Transition:             "0x0",
		Eip1052Transition:             "0x0",
		Eip1283Transition:             "0x0",
		Eip1283DisableTransition:      "0x0",
		Eip152Transition:              "0x0",
		Eip1108Transition:             "0x0",
		Eip1344Transition:             "0x0",
		Eip1884Transition:             "0x0",
		Eip2028Transition:             "0x0",
		Eip2200Transition:             "0x0",
		Eip2565Transition:             "0x0",
		Eip2929Transition:             "0x0",
		Eip2930Transition:             "0x0",
		Eip1559Transition:             "0x0",
		Eip3198Transition:             "0x0",
		Eip3529Transition:             "0x0",
		Eip3541Transition:             "0x0",
		Eip4895TransitionTimestamp:    "0x0",
		Eip3855TransitionTimestamp:    "0x0",
		Eip3651TransitionTimestamp:    "0x0",
		Eip3860TransitionTimestamp:    "0x0",
		Eip1153TransitionTimestamp:    "0x0",
		Eip4788TransitionTimestamp:    "0x0",
		Eip4844TransitionTimestamp:    "0x0",
		Eip5656TransitionTimestamp:    "0x0",
		Eip6780TransitionTimestamp:    "0x0",
		TerminalTotalDifficulty:       "0",
		TerminalTotalDifficultyPassed: true,
	}
}
