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

	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

var _ Genesis = (*Geth)(nil)

// Geth is a struct that holds the genesis configuration for Geth.
type Geth struct {
	*core.Genesis
}

func NewGeth() *Geth {
	return &Geth{
		&core.Genesis{
			Config:     defaultGethChainConfig(),
			Nonce:      zeroNonce,
			Timestamp:  zeroHex,
			ExtraData:  common.FromHex(defaultExtraData),
			GasLimit:   defaultGasLimit,
			Difficulty: big.NewInt(zeroHex),
			Alloc:      make(types.GenesisAlloc),
			Coinbase:   common.ExecutionAddress{},
		},
	}
}

func (g *Geth) AddAccount(address string, balance *big.Int) error {
	addr := common.HexToAddress(address)
	if _, ok := g.Alloc[addr]; ok {
		return errAccountAlreadyExists
	}
	g.Alloc[addr] = types.Account{
		Balance: balance,
	}
	return nil
}

func (g *Geth) AddPredeploy(
	address string,
	code []byte,
	balance *big.Int,
	nonce uint64) error {
	addr := common.HexToAddress(address)
	if _, ok := g.Alloc[addr]; ok {
		return errPredeployAlreadyExists
	}
	g.Alloc[addr] = types.Account{
		Code:    code,
		Balance: balance,
		Nonce:   nonce,
	}
	return nil
}

func (g *Geth) WriteJSON(filename string) error {
	return writeGenesisToJSON(g, filename)
}

func defaultGethChainConfig() *params.ChainConfig {
	return &params.ChainConfig{
		ChainID:                       big.NewInt(chainID),
		HomesteadBlock:                big.NewInt(0),
		DAOForkBlock:                  big.NewInt(0),
		DAOForkSupport:                true,
		EIP150Block:                   big.NewInt(0),
		EIP155Block:                   big.NewInt(0),
		EIP158Block:                   big.NewInt(0),
		ByzantiumBlock:                big.NewInt(0),
		ConstantinopleBlock:           big.NewInt(0),
		PetersburgBlock:               big.NewInt(0),
		IstanbulBlock:                 big.NewInt(0),
		MuirGlacierBlock:              big.NewInt(0),
		BerlinBlock:                   big.NewInt(0),
		LondonBlock:                   big.NewInt(0),
		ArrowGlacierBlock:             big.NewInt(0),
		GrayGlacierBlock:              big.NewInt(0),
		MergeNetsplitBlock:            big.NewInt(0),
		ShanghaiTime:                  new(uint64),
		CancunTime:                    new(uint64),
		TerminalTotalDifficulty:       big.NewInt(0),
		TerminalTotalDifficultyPassed: true,
	}
}
