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

package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
)

// ChainConfig is the core config which determines the blockchain settings.
//
// ChainConfig is stored in the database on a per block basis. This means
// that any network, identified by its genesis block, can have its own
// set of configuration options.
type ChainConfig struct {
	ChainID *big.Int `json:"chainId"` // chainId identifies the current chain and is used for replay protection

	HomesteadBlock *big.Int `json:"homesteadBlock,omitempty"` // Homestead switch block (nil = no fork, 0 = already homestead)

	DAOForkBlock   *big.Int `json:"daoForkBlock,omitempty"`   // TheDAO hard-fork switch block (nil = no fork)
	DAOForkSupport bool     `json:"daoForkSupport,omitempty"` // Whether the nodes supports or opposes the DAO hard-fork

	// EIP150 implements the Gas price changes (https://github.com/ethereum/EIPs/issues/150)
	EIP150Block *big.Int `json:"eip150Block,omitempty"` // EIP150 HF block (nil = no fork)
	EIP155Block *big.Int `json:"eip155Block,omitempty"` // EIP155 HF block
	EIP158Block *big.Int `json:"eip158Block,omitempty"` // EIP158 HF block

	ByzantiumBlock      *big.Int `json:"byzantiumBlock,omitempty"`      // Byzantium switch block (nil = no fork, 0 = already on byzantium)
	ConstantinopleBlock *big.Int `json:"constantinopleBlock,omitempty"` // Constantinople switch block (nil = no fork, 0 = already activated)
	PetersburgBlock     *big.Int `json:"petersburgBlock,omitempty"`     // Petersburg switch block (nil = same as Constantinople)
	IstanbulBlock       *big.Int `json:"istanbulBlock,omitempty"`       // Istanbul switch block (nil = no fork, 0 = already on istanbul)
	MuirGlacierBlock    *big.Int `json:"muirGlacierBlock,omitempty"`    // Eip-2384 (bomb delay) switch block (nil = no fork, 0 = already activated)
	BerlinBlock         *big.Int `json:"berlinBlock,omitempty"`         // Berlin switch block (nil = no fork, 0 = already on berlin)
	LondonBlock         *big.Int `json:"londonBlock,omitempty"`         // London switch block (nil = no fork, 0 = already on london)
	ArrowGlacierBlock   *big.Int `json:"arrowGlacierBlock,omitempty"`   // Eip-4345 (bomb delay) switch block (nil = no fork, 0 = already activated)
	GrayGlacierBlock    *big.Int `json:"grayGlacierBlock,omitempty"`    // Eip-5133 (bomb delay) switch block (nil = no fork, 0 = already activated)
	MergeNetsplitBlock  *big.Int `json:"mergeNetsplitBlock,omitempty"`  // Virtual fork after The Merge to use as a network splitter

	// Fork scheduling was switched from blocks to timestamps here

	ShanghaiTime  *uint64 `json:"shanghaiTime,omitempty"`  // Shanghai switch time (nil = no fork, 0 = already on shanghai)
	CancunTime    *uint64 `json:"cancunTime,omitempty"`    // Cancun switch time (nil = no fork, 0 = already on cancun)
	PragueTime    *uint64 `json:"pragueTime,omitempty"`    // Prague switch time (nil = no fork, 0 = already on prague)
	OsakaTime     *uint64 `json:"osakaTime,omitempty"`     // Osaka switch time (nil = no fork, 0 = already on osaka)
	BPO1Time      *uint64 `json:"bpo1Time,omitempty"`      // BPO1 switch time (nil = no fork, 0 = already on bpo1)
	BPO2Time      *uint64 `json:"bpo2Time,omitempty"`      // BPO2 switch time (nil = no fork, 0 = already on bpo2)
	BPO3Time      *uint64 `json:"bpo3Time,omitempty"`      // BPO3 switch time (nil = no fork, 0 = already on bpo3)
	BPO4Time      *uint64 `json:"bpo4Time,omitempty"`      // BPO4 switch time (nil = no fork, 0 = already on bpo4)
	BPO5Time      *uint64 `json:"bpo5Time,omitempty"`      // BPO5 switch time (nil = no fork, 0 = already on bpo5)
	AmsterdamTime *uint64 `json:"amsterdamTime,omitempty"` // Amsterdam switch time (nil = no fork, 0 = already on amsterdam)
	VerkleTime    *uint64 `json:"verkleTime,omitempty"`    // Verkle switch time (nil = no fork, 0 = already on verkle)

	// TerminalTotalDifficulty is the amount of total difficulty reached by
	// the network that triggers the consensus upgrade.
	TerminalTotalDifficulty *big.Int `json:"terminalTotalDifficulty,omitempty"`

	DepositContractAddress common.Address `json:"depositContractAddress,omitempty"`

	// EnableVerkleAtGenesis is a flag that specifies whether the network uses
	// the Verkle tree starting from the genesis block. If set to true, the
	// genesis state will be committed using the Verkle tree, eliminating the
	// need for any Verkle transition later.
	//
	// This is a temporary flag only for verkle devnet testing, where verkle is
	// activated at genesis, and the configured activation date has already passed.
	//
	// In production networks (mainnet and public testnets), verkle activation
	// always occurs after the genesis block, making this flag irrelevant in
	// those cases.
	EnableVerkleAtGenesis bool `json:"enableVerkleAtGenesis,omitempty"`

	// Various consensus engines
	Ethash             *params.EthashConfig       `json:"ethash,omitempty"`
	Clique             *params.CliqueConfig       `json:"clique,omitempty"`
	BlobScheduleConfig *params.BlobScheduleConfig `json:"blobSchedule,omitempty"`

	// Berachain config
	Berachain BerachainConfig `json:"berachain,omitempty"`
}

// BerachainConfig is the berachain config.
type BerachainConfig struct {
	// Prague1 fork values.
	Prague1 Prague1Config `json:"prague1,omitempty"`

	// Prague2 fork values.
	Prague2 Prague2Config `json:"prague2,omitempty"`

	// Prague3 fork values.
	Prague3 Prague3Config `json:"prague3,omitempty"`

	// Prague4 fork values.
	Prague4 Prague4Config `json:"prague4,omitempty"`
}

// Prague1Config is the config values for the Prague1 fork on Berachain.
type Prague1Config struct {
	// Time is the time of the Prague1 fork.
	Time *uint64 `json:"time,omitempty"` // Prague1 switch time (0 = already on prague1, nil = no fork)
	// BaseFeeChangeDenominator is the base fee change denominator.
	BaseFeeChangeDenominator uint64 `json:"baseFeeChangeDenominator,omitempty"`
	// MinimumBaseFeeWei is the minimum base fee in wei.
	MinimumBaseFeeWei *big.Int `json:"minimumBaseFeeWei,omitempty"`
	// PoLDistributorAddress is the address of the PoL distributor.
	PoLDistributorAddress common.Address `json:"polDistributorAddress,omitempty"`
}

// Prague2Config is the config values for the Prague2 fork on Berachain.
type Prague2Config struct {
	// Time is the time of the Prague2 fork.
	Time *uint64 `json:"time,omitempty"` // Prague2 switch time (0 = already on prague2, nil = no fork)
	// MinimumBaseFeeWei is the minimum base fee in wei.
	MinimumBaseFeeWei *big.Int `json:"minimumBaseFeeWei,omitempty"`
}

// Prague3Config is the config values for the Prague3 fork on Berachain.
type Prague3Config struct {
	// Time is the time of the Prague3 fork.
	Time *uint64 `json:"time,omitempty"` // Prague3 switch time (0 = already on prague3, nil = no fork)
	// BexVaultAddress is the address of the BEX vault.
	BexVaultAddress common.Address `json:"bexVaultAddress,omitempty"`
	// BlockedAddresses is the list of addresses blocked from sending or receiving ERC20 transfers.
	BlockedAddresses []common.Address `json:"blockedAddresses,omitempty"`
	// RescueAddress is the only address that blocked addresses can send to.
	RescueAddress common.Address `json:"rescueAddress,omitempty"`
}

// Prague4Config is the config values for the Prague4 fork on Berachain.
type Prague4Config struct {
	// Time is the time of the Prague4 fork.
	Time *uint64 `json:"time,omitempty"` // Prague4 switch time (0 = already on prague4, nil = no fork)
}

// IsLondon returns whether num is either equal to the London fork block or greater.
func (c *ChainConfig) IsLondon(num *big.Int) bool {
	return isBlockForked(c.LondonBlock, num)
}

// IsShanghai returns whether time is either equal to the Shanghai fork time or greater.
func (c *ChainConfig) IsShanghai(num *big.Int, time uint64) bool {
	return c.IsLondon(num) && isTimestampForked(c.ShanghaiTime, time)
}

// IsCancun returns whether time is either equal to the Cancun fork time or greater.
func (c *ChainConfig) IsCancun(num *big.Int, time uint64) bool {
	return c.IsLondon(num) && isTimestampForked(c.CancunTime, time)
}

// IsPrague returns whether time is either equal to the Prague fork time or greater.
func (c *ChainConfig) IsPrague(num *big.Int, time uint64) bool {
	return c.IsLondon(num) && isTimestampForked(c.PragueTime, time)
}

// IsPrague1 returns whether time is either equal to the Prague1 fork time or greater.
// NOTE: Prague1 is a Berachain fork and must be on Ethereum's Prague fork.
func (c *ChainConfig) IsPrague1(num *big.Int, time uint64) bool {
	return c.IsPrague(num, time) && isTimestampForked(c.Berachain.Prague1.Time, time)
}

// IsVerkleGenesis checks whether the verkle fork is activated at the genesis block.
//
// Verkle mode is considered enabled if the verkle fork time is configured,
// regardless of whether the local time has surpassed the fork activation time.
// This is a temporary workaround for verkle devnet testing, where verkle is
// activated at genesis, and the configured activation date has already passed.
//
// In production networks (mainnet and public testnets), verkle activation
// always occurs after the genesis block, making this function irrelevant in
// those cases.
func (c *ChainConfig) IsVerkleGenesis() bool {
	return c.EnableVerkleAtGenesis
}

// isBlockForked returns whether a fork scheduled at block s is active at the
// given head block. Whilst this method is the same as isTimestampForked, they
// are explicitly separate for clearer reading.
func isBlockForked(s, head *big.Int) bool {
	if s == nil || head == nil {
		return false
	}
	return s.Cmp(head) <= 0
}

// isTimestampForked returns whether a fork scheduled at timestamp s is active
// at the given head timestamp. Whilst this method is the same as isBlockForked,
// they are explicitly separate for clearer reading.
func isTimestampForked(s *uint64, head uint64) bool {
	if s == nil {
		return false
	}
	return *s <= head
}
