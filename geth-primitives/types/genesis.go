// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

//nolint:all // This package contains types copied/adapted from go-ethereum (geth).
package types

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/ethereum/go-ethereum/triedb"
	"github.com/ethereum/go-ethereum/triedb/pathdb"
	"github.com/holiman/uint256"
)

// Genesis specifies the header fields, state of a genesis block. It also defines hard
// fork switch-over blocks through the chain configuration.
type Genesis struct {
	Config     *ChainConfig       `json:"config"`
	Nonce      uint64             `json:"nonce"`
	Timestamp  uint64             `json:"timestamp"`
	ExtraData  []byte             `json:"extraData"`
	GasLimit   uint64             `json:"gasLimit"   gencodec:"required"`
	Difficulty *big.Int           `json:"difficulty" gencodec:"required"`
	Mixhash    common.Hash        `json:"mixHash"`
	Coinbase   common.Address     `json:"coinbase"`
	Alloc      types.GenesisAlloc `json:"alloc"      gencodec:"required"`

	// These fields are used for consensus tests. Please don't use them
	// in actual genesis blocks.
	Number        uint64      `json:"number"`
	GasUsed       uint64      `json:"gasUsed"`
	ParentHash    common.Hash `json:"parentHash"`
	BaseFee       *big.Int    `json:"baseFeePerGas"` // EIP-1559
	ExcessBlobGas *uint64     `json:"excessBlobGas"` // EIP-4844
	BlobGasUsed   *uint64     `json:"blobGasUsed"`   // EIP-4844
}

// IsVerkle indicates whether the state is already stored in a verkle
// tree at genesis time.
func (g *Genesis) IsVerkle() bool {
	return g.Config.IsVerkleGenesis()
}

// ToBlock returns the genesis block according to genesis specification.
func (g *Genesis) ToBlock() *Block {
	root, err := hashAlloc(&g.Alloc, g.IsVerkle())
	if err != nil {
		panic(err)
	}
	return g.toBlockWithRoot(root)
}

// UnmarshalJSON unmarshals from JSON.
func (g *Genesis) UnmarshalJSON(input []byte) error {
	type Genesis struct {
		Config        *ChainConfig                               `json:"config"`
		Nonce         *math.HexOrDecimal64                       `json:"nonce"`
		Timestamp     *math.HexOrDecimal64                       `json:"timestamp"`
		ExtraData     *hexutil.Bytes                             `json:"extraData"`
		GasLimit      *math.HexOrDecimal64                       `json:"gasLimit"   gencodec:"required"`
		Difficulty    *math.HexOrDecimal256                      `json:"difficulty" gencodec:"required"`
		Mixhash       *common.Hash                               `json:"mixHash"`
		Coinbase      *common.Address                            `json:"coinbase"`
		Alloc         map[common.UnprefixedAddress]types.Account `json:"alloc"      gencodec:"required"`
		Number        *math.HexOrDecimal64                       `json:"number"`
		GasUsed       *math.HexOrDecimal64                       `json:"gasUsed"`
		ParentHash    *common.Hash                               `json:"parentHash"`
		BaseFee       *math.HexOrDecimal256                      `json:"baseFeePerGas"`
		ExcessBlobGas *math.HexOrDecimal64                       `json:"excessBlobGas"`
		BlobGasUsed   *math.HexOrDecimal64                       `json:"blobGasUsed"`
	}
	var dec Genesis
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.Config != nil {
		g.Config = dec.Config
	}
	if dec.Nonce != nil {
		g.Nonce = uint64(*dec.Nonce)
	}
	if dec.Timestamp != nil {
		g.Timestamp = uint64(*dec.Timestamp)
	}
	if dec.ExtraData != nil {
		g.ExtraData = *dec.ExtraData
	}
	if dec.GasLimit == nil {
		return errors.New("missing required field 'gasLimit' for Genesis")
	}
	g.GasLimit = uint64(*dec.GasLimit)
	if dec.Difficulty == nil {
		return errors.New("missing required field 'difficulty' for Genesis")
	}
	g.Difficulty = (*big.Int)(dec.Difficulty)
	if dec.Mixhash != nil {
		g.Mixhash = *dec.Mixhash
	}
	if dec.Coinbase != nil {
		g.Coinbase = *dec.Coinbase
	}
	if dec.Alloc == nil {
		return errors.New("missing required field 'alloc' for Genesis")
	}
	g.Alloc = make(types.GenesisAlloc, len(dec.Alloc))
	for k, v := range dec.Alloc {
		g.Alloc[common.Address(k)] = v
	}
	if dec.Number != nil {
		g.Number = uint64(*dec.Number)
	}
	if dec.GasUsed != nil {
		g.GasUsed = uint64(*dec.GasUsed)
	}
	if dec.ParentHash != nil {
		g.ParentHash = *dec.ParentHash
	}
	if dec.BaseFee != nil {
		g.BaseFee = (*big.Int)(dec.BaseFee)
	}
	if dec.ExcessBlobGas != nil {
		g.ExcessBlobGas = (*uint64)(dec.ExcessBlobGas)
	}
	if dec.BlobGasUsed != nil {
		g.BlobGasUsed = (*uint64)(dec.BlobGasUsed)
	}
	return nil
}

// toBlockWithRoot constructs the genesis block with the given genesis state root.
func (g *Genesis) toBlockWithRoot(root common.Hash) *Block {
	head := &Header{
		Number:     new(big.Int).SetUint64(g.Number),
		Nonce:      types.EncodeNonce(g.Nonce),
		Time:       g.Timestamp,
		ParentHash: g.ParentHash,
		Extra:      g.ExtraData,
		GasLimit:   g.GasLimit,
		GasUsed:    g.GasUsed,
		BaseFee:    g.BaseFee,
		Difficulty: g.Difficulty,
		MixDigest:  g.Mixhash,
		Coinbase:   g.Coinbase,
		Root:       root,
	}
	if g.GasLimit == 0 {
		head.GasLimit = params.GenesisGasLimit
	}
	if g.Difficulty == nil {
		if g.Config != nil && g.Config.Ethash == nil {
			head.Difficulty = big.NewInt(0)
		} else if g.Mixhash == (common.Hash{}) {
			head.Difficulty = params.GenesisDifficulty
		}
	}
	if g.Config != nil && g.Config.IsLondon(common.Big0) {
		if g.BaseFee != nil {
			head.BaseFee = g.BaseFee
		} else {
			head.BaseFee = new(big.Int).SetUint64(params.InitialBaseFee)
		}
	}
	var (
		withdrawals []*types.Withdrawal
	)
	if conf := g.Config; conf != nil {
		num := big.NewInt(int64(g.Number))
		if conf.IsShanghai(num, g.Timestamp) {
			head.WithdrawalsHash = &types.EmptyWithdrawalsHash
			withdrawals = make([]*types.Withdrawal, 0)
		}
		if conf.IsCancun(num, g.Timestamp) {
			// EIP-4788: The parentBeaconBlockRoot of the genesis block is always
			// the zero hash. This is because the genesis block does not have a parent
			// by definition.
			head.ParentBeaconRoot = new(common.Hash)
			// EIP-4844 fields
			head.ExcessBlobGas = g.ExcessBlobGas
			head.BlobGasUsed = g.BlobGasUsed
			if head.ExcessBlobGas == nil {
				head.ExcessBlobGas = new(uint64)
			}
			if head.BlobGasUsed == nil {
				head.BlobGasUsed = new(uint64)
			}
		} else {
			if g.ExcessBlobGas != nil {
				log.Warn("Invalid genesis, unexpected ExcessBlobGas set before Cancun, allowing it for testing purposes")
				head.ExcessBlobGas = g.ExcessBlobGas
			}
		}
		if conf.IsPrague(num, g.Timestamp) {
			head.RequestsHash = &types.EmptyRequestsHash
		}
		if conf.IsPrague1(num, g.Timestamp) {
			// BRIP-0004: The parentProposerPubkey of the genesis block is always
			// the zero pubkey. This is because the genesis block does not have a parent
			// by definition.
			head.ParentProposerPubkey = new(ExecutionPubkey)
		}
	}
	return NewBlock(head, &Body{Withdrawals: withdrawals}, nil, trie.NewStackTrie(nil))
}

// hashAlloc computes the state root according to the genesis specification.
func hashAlloc(ga *types.GenesisAlloc, isVerkle bool) (common.Hash, error) {
	// If a genesis-time verkle trie is requested, create a trie config
	// with the verkle trie enabled so that the tree can be initialized
	// as such.
	var config *triedb.Config
	if isVerkle {
		config = &triedb.Config{
			PathDB:   pathdb.Defaults,
			IsVerkle: true,
		}
	}
	// Create an ephemeral in-memory database for computing hash,
	// all the derived states will be discarded to not pollute disk.
	emptyRoot := types.EmptyRootHash
	if isVerkle {
		emptyRoot = types.EmptyVerkleHash
	}
	db := rawdb.NewMemoryDatabase()
	statedb, err := state.New(emptyRoot, state.NewDatabase(triedb.NewDatabase(db, config), nil))
	if err != nil {
		return common.Hash{}, err
	}
	for addr, account := range *ga {
		if account.Balance != nil {
			statedb.AddBalance(addr, uint256.MustFromBig(account.Balance), tracing.BalanceIncreaseGenesisBalance)
		}
		statedb.SetCode(addr, account.Code, tracing.CodeChangeGenesis)
		statedb.SetNonce(addr, account.Nonce, tracing.NonceChangeGenesis)
		for key, value := range account.Storage {
			statedb.SetState(addr, key, value)
		}
	}
	return statedb.Commit(0, false, false)
}
