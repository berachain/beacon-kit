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
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
)

const (
	// BRIP-0002 - Base Fee Change Denominator
	BerachainBaseFeeChangeDenominator = 48 // 6x increase from the default
)

var (
	// Deposit Contract
	BerachainDepositContractAddress = common.HexToAddress("0x4242424242424242424242424242424242424242")

	// BRIP-0004 - PoL Distributor
	PoLTxGasLimit         = uint64(30_000_000)
	PoLDistributorAddress = common.HexToAddress("0xD2f19a79b026Fb636A7c300bF5947df113940761")
)

// Genesis hashes to enforce below configs on.
var (
	BerachainGenesisHash = common.HexToHash("0xd57819422128da1c44339fc7956662378c17e2213e669b427ac91cd11dfcfb38")
	BepoliaGenesisHash   = common.HexToHash("0x0207661de38f0e54ba91c8286096e72486784c79dc6a9681fc486b38335c042f")
)

func newUint64(val uint64) *uint64 { return &val }

var (
	// BerachainChainConfig is the chain parameters to run a node on the Berachain network.
	BerachainChainConfig = &ChainConfig{
		ChainID:                 big.NewInt(80094),
		HomesteadBlock:          big.NewInt(0),
		DAOForkBlock:            big.NewInt(0),
		DAOForkSupport:          true,
		EIP150Block:             big.NewInt(0),
		EIP155Block:             big.NewInt(0),
		EIP158Block:             big.NewInt(0),
		ByzantiumBlock:          big.NewInt(0),
		ConstantinopleBlock:     big.NewInt(0),
		PetersburgBlock:         big.NewInt(0),
		IstanbulBlock:           big.NewInt(0),
		MuirGlacierBlock:        big.NewInt(0),
		BerlinBlock:             big.NewInt(0),
		LondonBlock:             big.NewInt(0),
		ArrowGlacierBlock:       big.NewInt(0),
		GrayGlacierBlock:        big.NewInt(0),
		TerminalTotalDifficulty: big.NewInt(0),
		MergeNetsplitBlock:      big.NewInt(0),
		ShanghaiTime:            newUint64(0),
		CancunTime:              newUint64(0),
		PragueTime:              newUint64(1749056400),
		DepositContractAddress:  BerachainDepositContractAddress,
		Ethash:                  new(EthashConfig),
		BlobScheduleConfig: &BlobScheduleConfig{
			Cancun: DefaultCancunBlobConfig,
			Prague: DefaultBerachainPragueBlobConfig,
		},
		Berachain: BerachainConfig{
			Prague1: Prague1Config{
				Time:                     newUint64(1756915200), // Sep 03 2025 16:00:00 UTC
				MinimumBaseFeeWei:        big.NewInt(1 * params.GWei),
				BaseFeeChangeDenominator: BerachainBaseFeeChangeDenominator,
				PoLDistributorAddress:    PoLDistributorAddress,
			},
			Prague2: Prague2Config{
				Time:              newUint64(1759248000), // Sep 30 2025 16:00:00 UTC
				MinimumBaseFeeWei: big.NewInt(0),
			},
			Prague3: Prague3Config{
				Time:            newUint64(1762164459), // Nov 03 2025 10:07:39 UTC
				BexVaultAddress: common.HexToAddress("0x4be03f781c497a489e3cb0287833452ca9b9e80b"),
				BlockedAddresses: []common.Address{
					common.HexToAddress("0x9BAD77F1D527CD2D023d33eB3597A456d0c1Ab4a"),
					common.HexToAddress("0xD875De13Dc789B070a9F2a4549fbBb94cCdA4112"),
					common.HexToAddress("0xF8Bec8cB704b8BD427FD209A2058b396C4BC543e"),
					common.HexToAddress("0xF2b63Dbf539f4862a2eA3a04520D4E04ed5b499C"),
					common.HexToAddress("0x506D1f9EFe24f0d47853aDca907EB8d89AE03207"),
					common.HexToAddress("0x045371528a01071d6e5c934d42d641fd3cbe941c"),
					common.HexToAddress("0xF8be2BF5a14f17C897d00b57fb40EcF8b96c543e"),
					common.HexToAddress("0x9BAD91648D4769695591853478E628bCb499AB4A"),
				},
				RescueAddress: common.HexToAddress("0xD276D30592bE512a418f2448e23f9E7F372b32A2"),
			},
			Prague4: Prague4Config{
				Time: newUint64(1762963200), // Nov 12 2025 16:00:00 UTC
			},
		},
	}
	// BepoliaChainConfig contains the chain parameters to run a node on the Bepolia test network.
	BepoliaChainConfig = &ChainConfig{
		ChainID:                 big.NewInt(80069),
		HomesteadBlock:          big.NewInt(0),
		DAOForkBlock:            big.NewInt(0),
		DAOForkSupport:          true,
		EIP150Block:             big.NewInt(0),
		EIP155Block:             big.NewInt(0),
		EIP158Block:             big.NewInt(0),
		ByzantiumBlock:          big.NewInt(0),
		ConstantinopleBlock:     big.NewInt(0),
		PetersburgBlock:         big.NewInt(0),
		IstanbulBlock:           big.NewInt(0),
		MuirGlacierBlock:        big.NewInt(0),
		BerlinBlock:             big.NewInt(0),
		LondonBlock:             big.NewInt(0),
		ArrowGlacierBlock:       big.NewInt(0),
		GrayGlacierBlock:        big.NewInt(0),
		TerminalTotalDifficulty: big.NewInt(0),
		MergeNetsplitBlock:      big.NewInt(0),
		ShanghaiTime:            newUint64(0),
		CancunTime:              newUint64(0),
		PragueTime:              newUint64(1746633600),
		DepositContractAddress:  BerachainDepositContractAddress,
		Ethash:                  new(EthashConfig),
		BlobScheduleConfig: &BlobScheduleConfig{
			Cancun: DefaultCancunBlobConfig,
			Prague: DefaultBerachainPragueBlobConfig,
		},
		Berachain: BerachainConfig{
			Prague1: Prague1Config{
				Time:                     newUint64(1754496000), // Aug 06 2025 16:00:00 UTC
				MinimumBaseFeeWei:        big.NewInt(10 * params.GWei),
				BaseFeeChangeDenominator: BerachainBaseFeeChangeDenominator,
				PoLDistributorAddress:    PoLDistributorAddress,
			},
			Prague2: Prague2Config{
				Time:              newUint64(1758124800), // Sep 17 2025 16:00:00 UTC
				MinimumBaseFeeWei: big.NewInt(0),
			},
		},
	}
)

var (
	// DefaultCancunBlobConfig is the default blob configuration for the Cancun fork.
	DefaultCancunBlobConfig = &BlobConfig{
		Target:         3,
		Max:            6,
		UpdateFraction: 3338477,
	}
	// DefaultPragueBlobConfig is the default blob configuration for the Prague fork.
	DefaultPragueBlobConfig = &BlobConfig{
		Target:         6,
		Max:            9,
		UpdateFraction: 5007716,
	}
	// DefaultBerachainPragueBlobConfig is the default blob configuration for the Prague fork
	// on Berachain networks.
	DefaultBerachainPragueBlobConfig = &BlobConfig{
		Target:         3,
		Max:            6,
		UpdateFraction: 3338477,
	}
	// DefaultOsakaBlobConfig is the default blob configuration for the Osaka fork.
	DefaultOsakaBlobConfig = &BlobConfig{
		Target:         6,
		Max:            9,
		UpdateFraction: 5007716,
	}
	// DefaultBPO1BlobConfig is the default blob configuration for the BPO1 fork.
	DefaultBPO1BlobConfig = &BlobConfig{
		Target:         10,
		Max:            15,
		UpdateFraction: 8346193,
	}
	// DefaultBPO2BlobConfig is the default blob configuration for the BPO2 fork.
	DefaultBPO2BlobConfig = &BlobConfig{
		Target:         14,
		Max:            21,
		UpdateFraction: 11684671,
	}
	// DefaultBPO3BlobConfig is the default blob configuration for the BPO3 fork.
	DefaultBPO3BlobConfig = &BlobConfig{
		Target:         21,
		Max:            32,
		UpdateFraction: 20609697,
	}
	// DefaultBPO4BlobConfig is the default blob configuration for the BPO4 fork.
	DefaultBPO4BlobConfig = &BlobConfig{
		Target:         14,
		Max:            21,
		UpdateFraction: 13739630,
	}
	// DefaultBlobSchedule is the latest configured blob schedule for Ethereum mainnet.
	DefaultBlobSchedule = &BlobScheduleConfig{
		Cancun: DefaultCancunBlobConfig,
		Prague: DefaultPragueBlobConfig,
		Osaka:  DefaultOsakaBlobConfig,
	}
)

// NetworkNames are user friendly names to use in the chain spec banner.
var NetworkNames = map[string]string{
	BerachainChainConfig.ChainID.String(): "berachain mainnet",
	BepoliaChainConfig.ChainID.String():   "bepolia testnet",
}

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
	Ethash             *EthashConfig       `json:"ethash,omitempty"`
	Clique             *CliqueConfig       `json:"clique,omitempty"`
	BlobScheduleConfig *BlobScheduleConfig `json:"blobSchedule,omitempty"`

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

// String implements the stringer interface.
func (o *BerachainConfig) String() string {
	banner := "berachain"
	if o.Prague1.Time != nil {
		banner += fmt.Sprintf("(%s)", o.Prague1)
	}
	if o.Prague2.Time != nil {
		banner += fmt.Sprintf("(%s)", o.Prague2)
	}
	if o.Prague3.Time != nil {
		banner += fmt.Sprintf("(%s)", o.Prague3)
	}
	if o.Prague4.Time != nil {
		banner += fmt.Sprintf("(%s)", o.Prague4)
	}
	return banner
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

// String implements the stringer interface.
func (c Prague1Config) String() string {
	banner := "prague1"
	if c.Time != nil {
		banner += fmt.Sprintf(
			"(time: %v, baseFeeChangeDenominator: %v, minimumBaseFeeWei: %v, polDistributorAddress: %v)",
			*c.Time, c.BaseFeeChangeDenominator, c.MinimumBaseFeeWei, c.PoLDistributorAddress,
		)
	}
	return banner
}

// Prague2Config is the config values for the Prague2 fork on Berachain.
type Prague2Config struct {
	// Time is the time of the Prague2 fork.
	Time *uint64 `json:"time,omitempty"` // Prague2 switch time (0 = already on prague2, nil = no fork)
	// MinimumBaseFeeWei is the minimum base fee in wei.
	MinimumBaseFeeWei *big.Int `json:"minimumBaseFeeWei,omitempty"`
}

// String implements the stringer interface.
func (c Prague2Config) String() string {
	banner := "prague2"
	if c.Time != nil {
		banner += fmt.Sprintf("(time: %v, minimumBaseFeeWei: %v)", *c.Time, c.MinimumBaseFeeWei)
	}
	return banner
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

// String implements the stringer interface.
func (c Prague3Config) String() string {
	banner := "prague3"
	if c.Time != nil {
		blocked := make([]string, 0, len(c.BlockedAddresses))
		for _, addr := range c.BlockedAddresses {
			blocked = append(blocked, addr.String())
		}
		banner += fmt.Sprintf("(time: %v, bexVaultAddress: %v, blockedAddresses: [%s], rescueAddress: %v)", *c.Time, c.BexVaultAddress, strings.Join(blocked, ", "), c.RescueAddress)
	}
	return banner
}

// Prague4Config is the config values for the Prague4 fork on Berachain.
type Prague4Config struct {
	// Time is the time of the Prague4 fork.
	Time *uint64 `json:"time,omitempty"` // Prague4 switch time (0 = already on prague4, nil = no fork)
}

// String implements the stringer interface.
func (c Prague4Config) String() string {
	banner := "prague4"
	if c.Time != nil {
		banner += fmt.Sprintf("(time: %v)", *c.Time)
	}
	return banner
}

// EthashConfig is the consensus engine configs for proof-of-work based sealing.
type EthashConfig struct{}

// String implements the stringer interface, returning the consensus engine details.
func (c EthashConfig) String() string {
	return "ethash"
}

// CliqueConfig is the consensus engine configs for proof-of-authority based sealing.
type CliqueConfig struct {
	Period uint64 `json:"period"` // Number of seconds between blocks to enforce
	Epoch  uint64 `json:"epoch"`  // Epoch length to reset votes and checkpoint
}

// String implements the stringer interface, returning the consensus engine details.
func (c CliqueConfig) String() string {
	return fmt.Sprintf("clique(period: %d, epoch: %d)", c.Period, c.Epoch)
}

// String implements the fmt.Stringer interface, returning a string representation
// of ChainConfig.
func (c *ChainConfig) String() string {
	result := fmt.Sprintf("ChainConfig{ChainID: %v", c.ChainID)

	// Add block-based forks
	if c.HomesteadBlock != nil {
		result += fmt.Sprintf(", HomesteadBlock: %v", c.HomesteadBlock)
	}
	if c.DAOForkBlock != nil {
		result += fmt.Sprintf(", DAOForkBlock: %v", c.DAOForkBlock)
	}
	if c.EIP150Block != nil {
		result += fmt.Sprintf(", EIP150Block: %v", c.EIP150Block)
	}
	if c.EIP155Block != nil {
		result += fmt.Sprintf(", EIP155Block: %v", c.EIP155Block)
	}
	if c.EIP158Block != nil {
		result += fmt.Sprintf(", EIP158Block: %v", c.EIP158Block)
	}
	if c.ByzantiumBlock != nil {
		result += fmt.Sprintf(", ByzantiumBlock: %v", c.ByzantiumBlock)
	}
	if c.ConstantinopleBlock != nil {
		result += fmt.Sprintf(", ConstantinopleBlock: %v", c.ConstantinopleBlock)
	}
	if c.PetersburgBlock != nil {
		result += fmt.Sprintf(", PetersburgBlock: %v", c.PetersburgBlock)
	}
	if c.IstanbulBlock != nil {
		result += fmt.Sprintf(", IstanbulBlock: %v", c.IstanbulBlock)
	}
	if c.MuirGlacierBlock != nil {
		result += fmt.Sprintf(", MuirGlacierBlock: %v", c.MuirGlacierBlock)
	}
	if c.BerlinBlock != nil {
		result += fmt.Sprintf(", BerlinBlock: %v", c.BerlinBlock)
	}
	if c.LondonBlock != nil {
		result += fmt.Sprintf(", LondonBlock: %v", c.LondonBlock)
	}
	if c.ArrowGlacierBlock != nil {
		result += fmt.Sprintf(", ArrowGlacierBlock: %v", c.ArrowGlacierBlock)
	}
	if c.GrayGlacierBlock != nil {
		result += fmt.Sprintf(", GrayGlacierBlock: %v", c.GrayGlacierBlock)
	}
	if c.MergeNetsplitBlock != nil {
		result += fmt.Sprintf(", MergeNetsplitBlock: %v", c.MergeNetsplitBlock)
	}

	// Add timestamp-based forks
	if c.ShanghaiTime != nil {
		result += fmt.Sprintf(", ShanghaiTime: %v", *c.ShanghaiTime)
	}
	if c.CancunTime != nil {
		result += fmt.Sprintf(", CancunTime: %v", *c.CancunTime)
	}
	if c.PragueTime != nil {
		result += fmt.Sprintf(", PragueTime: %v", *c.PragueTime)
	}
	if c.OsakaTime != nil {
		result += fmt.Sprintf(", OsakaTime: %v", *c.OsakaTime)
	}
	if c.BPO1Time != nil {
		result += fmt.Sprintf(", BPO1Time: %v", *c.BPO1Time)
	}
	if c.BPO2Time != nil {
		result += fmt.Sprintf(", BPO2Time: %v", *c.BPO2Time)
	}
	if c.BPO3Time != nil {
		result += fmt.Sprintf(", BPO3Time: %v", *c.BPO3Time)
	}
	if c.BPO4Time != nil {
		result += fmt.Sprintf(", BPO4Time: %v", *c.BPO4Time)
	}
	if c.BPO5Time != nil {
		result += fmt.Sprintf(", BPO5Time: %v", *c.BPO5Time)
	}
	if c.AmsterdamTime != nil {
		result += fmt.Sprintf(", AmsterdamTime: %v", *c.AmsterdamTime)
	}
	if c.VerkleTime != nil {
		result += fmt.Sprintf(", VerkleTime: %v", *c.VerkleTime)
	}
	result += "}"
	return result
}

// Description returns a human-readable description of ChainConfig.
func (c *ChainConfig) Description() string {
	var banner string

	// Create some basic network config output
	network := NetworkNames[c.ChainID.String()]
	if network == "" {
		network = "unknown"
	}
	banner += fmt.Sprintf("Chain ID:  %v (%s)\n", c.ChainID, network)
	switch {
	case c.Ethash != nil:
		banner += "Consensus: Beacon (proof-of-stake), merged from Ethash (proof-of-work)\n"
	case c.Clique != nil:
		banner += "Consensus: Beacon (proof-of-stake), merged from Clique (proof-of-authority)\n"
	default:
		banner += "Consensus: unknown\n"
	}
	banner += "\n"

	// Create a list of forks with a short description of them. Forks that only
	// makes sense for mainnet should be optional at printing to avoid bloating
	// the output for testnets and private networks.
	banner += "Pre-Merge hard forks (block based):\n"
	banner += fmt.Sprintf(" - Homestead:                   #%-8v\n", c.HomesteadBlock)
	if c.DAOForkBlock != nil {
		banner += fmt.Sprintf(" - DAO Fork:                    #%-8v\n", c.DAOForkBlock)
	}
	banner += fmt.Sprintf(" - Tangerine Whistle (EIP 150): #%-8v\n", c.EIP150Block)
	banner += fmt.Sprintf(" - Spurious Dragon/1 (EIP 155): #%-8v\n", c.EIP155Block)
	banner += fmt.Sprintf(" - Spurious Dragon/2 (EIP 158): #%-8v\n", c.EIP158Block)
	banner += fmt.Sprintf(" - Byzantium:                   #%-8v\n", c.ByzantiumBlock)
	banner += fmt.Sprintf(" - Constantinople:              #%-8v\n", c.ConstantinopleBlock)
	banner += fmt.Sprintf(" - Petersburg:                  #%-8v\n", c.PetersburgBlock)
	banner += fmt.Sprintf(" - Istanbul:                    #%-8v\n", c.IstanbulBlock)
	if c.MuirGlacierBlock != nil {
		banner += fmt.Sprintf(" - Muir Glacier:                #%-8v\n", c.MuirGlacierBlock)
	}
	banner += fmt.Sprintf(" - Berlin:                      #%-8v\n", c.BerlinBlock)
	banner += fmt.Sprintf(" - London:                      #%-8v\n", c.LondonBlock)
	if c.ArrowGlacierBlock != nil {
		banner += fmt.Sprintf(" - Arrow Glacier:               #%-8v\n", c.ArrowGlacierBlock)
	}
	if c.GrayGlacierBlock != nil {
		banner += fmt.Sprintf(" - Gray Glacier:                #%-8v\n", c.GrayGlacierBlock)
	}
	banner += "\n"

	// Add a special section for the merge as it's non-obvious
	banner += "Merge configured:\n"
	banner += fmt.Sprintf(" - Total terminal difficulty:  %v\n", c.TerminalTotalDifficulty)
	if c.MergeNetsplitBlock != nil {
		banner += fmt.Sprintf(" - Merge netsplit block:       #%-8v\n", c.MergeNetsplitBlock)
	}
	banner += "\n"

	// Create a list of forks post-merge
	banner += "Post-Merge hard forks (timestamp based):\n"
	if c.ShanghaiTime != nil {
		banner += fmt.Sprintf(" - Shanghai:                    @%-10v\n", *c.ShanghaiTime)
	}
	if c.CancunTime != nil {
		banner += fmt.Sprintf(" - Cancun:                      @%-10v blob: (%s)\n", *c.CancunTime, c.BlobScheduleConfig.Cancun)
	}
	if c.PragueTime != nil {
		banner += fmt.Sprintf(" - Prague:                      @%-10v blob: (%s)\n", *c.PragueTime, c.BlobScheduleConfig.Prague)
	}
	if c.Berachain.Prague1.Time != nil {
		banner += fmt.Sprintf(" - Prague1:                     %-10v (https://github.com/berachain/BRIPs/blob/main/meta/BRIP-0004.md)\n", c.Berachain.Prague1)
	}
	if c.Berachain.Prague2.Time != nil {
		banner += fmt.Sprintf(" - Prague2:                     %-10v\n", c.Berachain.Prague2)
	}
	if c.Berachain.Prague3.Time != nil {
		banner += fmt.Sprintf(" - Prague3:                     %-10v\n", c.Berachain.Prague3)
	}
	if c.Berachain.Prague4.Time != nil {
		banner += fmt.Sprintf(" - Prague4:                     %-10v\n", c.Berachain.Prague4)
	}
	if c.OsakaTime != nil {
		banner += fmt.Sprintf(" - Osaka:                       @%-10v blob: (%s)\n", *c.OsakaTime, c.BlobScheduleConfig.Osaka)
	}
	if c.BPO1Time != nil {
		banner += fmt.Sprintf(" - BPO1:                        @%-10v blob: (%s)\n", *c.BPO1Time, c.BlobScheduleConfig.BPO1)
	}
	if c.BPO2Time != nil {
		banner += fmt.Sprintf(" - BPO2:                        @%-10v blob: (%s)\n", *c.BPO2Time, c.BlobScheduleConfig.BPO2)
	}
	if c.BPO3Time != nil {
		banner += fmt.Sprintf(" - BPO3:                        @%-10v blob: (%s)\n", *c.BPO3Time, c.BlobScheduleConfig.BPO3)
	}
	if c.BPO4Time != nil {
		banner += fmt.Sprintf(" - BPO4:                        @%-10v blob: (%s)\n", *c.BPO4Time, c.BlobScheduleConfig.BPO4)
	}
	if c.BPO5Time != nil {
		banner += fmt.Sprintf(" - BPO5:                        @%-10v blob: (%s)\n", *c.BPO5Time, c.BlobScheduleConfig.BPO5)
	}
	if c.AmsterdamTime != nil {
		banner += fmt.Sprintf(" - Amsterdam:									 @%-10v blob: (%s)\n", *c.AmsterdamTime, c.BlobScheduleConfig.Amsterdam)
	}
	if c.VerkleTime != nil {
		banner += fmt.Sprintf(" - Verkle:                      @%-10v blob: (%s)\n", *c.VerkleTime, c.BlobScheduleConfig.Verkle)
	}
	banner += fmt.Sprintf("\nAll fork specifications can be found at https://ethereum.github.io/execution-specs/src/ethereum/forks/\n")
	return banner
}

// BlobConfig specifies the target and max blobs per block for the associated fork.
type BlobConfig struct {
	Target         int    `json:"target"`
	Max            int    `json:"max"`
	UpdateFraction uint64 `json:"baseFeeUpdateFraction"`
}

// String implement fmt.Stringer, returning string format blob config.
func (bc *BlobConfig) String() string {
	if bc == nil {
		return "nil"
	}
	return fmt.Sprintf("target: %d, max: %d, fraction: %d", bc.Target, bc.Max, bc.UpdateFraction)
}

// BlobScheduleConfig determines target and max number of blobs allow per fork.
type BlobScheduleConfig struct {
	Cancun    *BlobConfig `json:"cancun,omitempty"`
	Prague    *BlobConfig `json:"prague,omitempty"`
	Osaka     *BlobConfig `json:"osaka,omitempty"`
	Verkle    *BlobConfig `json:"verkle,omitempty"`
	BPO1      *BlobConfig `json:"bpo1,omitempty"`
	BPO2      *BlobConfig `json:"bpo2,omitempty"`
	BPO3      *BlobConfig `json:"bpo3,omitempty"`
	BPO4      *BlobConfig `json:"bpo4,omitempty"`
	BPO5      *BlobConfig `json:"bpo5,omitempty"`
	Amsterdam *BlobConfig `json:"amsterdam,omitempty"`
}

// IsHomestead returns whether num is either equal to the homestead block or greater.
func (c *ChainConfig) IsHomestead(num *big.Int) bool {
	return isBlockForked(c.HomesteadBlock, num)
}

// IsDAOFork returns whether num is either equal to the DAO fork block or greater.
func (c *ChainConfig) IsDAOFork(num *big.Int) bool {
	return isBlockForked(c.DAOForkBlock, num)
}

// IsEIP150 returns whether num is either equal to the EIP150 fork block or greater.
func (c *ChainConfig) IsEIP150(num *big.Int) bool {
	return isBlockForked(c.EIP150Block, num)
}

// IsEIP155 returns whether num is either equal to the EIP155 fork block or greater.
func (c *ChainConfig) IsEIP155(num *big.Int) bool {
	return isBlockForked(c.EIP155Block, num)
}

// IsEIP158 returns whether num is either equal to the EIP158 fork block or greater.
func (c *ChainConfig) IsEIP158(num *big.Int) bool {
	return isBlockForked(c.EIP158Block, num)
}

// IsByzantium returns whether num is either equal to the Byzantium fork block or greater.
func (c *ChainConfig) IsByzantium(num *big.Int) bool {
	return isBlockForked(c.ByzantiumBlock, num)
}

// IsConstantinople returns whether num is either equal to the Constantinople fork block or greater.
func (c *ChainConfig) IsConstantinople(num *big.Int) bool {
	return isBlockForked(c.ConstantinopleBlock, num)
}

// IsMuirGlacier returns whether num is either equal to the Muir Glacier (EIP-2384) fork block or greater.
func (c *ChainConfig) IsMuirGlacier(num *big.Int) bool {
	return isBlockForked(c.MuirGlacierBlock, num)
}

// IsPetersburg returns whether num is either
// - equal to or greater than the PetersburgBlock fork block,
// - OR is nil, and Constantinople is active
func (c *ChainConfig) IsPetersburg(num *big.Int) bool {
	return isBlockForked(c.PetersburgBlock, num) || c.PetersburgBlock == nil && isBlockForked(c.ConstantinopleBlock, num)
}

// IsIstanbul returns whether num is either equal to the Istanbul fork block or greater.
func (c *ChainConfig) IsIstanbul(num *big.Int) bool {
	return isBlockForked(c.IstanbulBlock, num)
}

// IsBerlin returns whether num is either equal to the Berlin fork block or greater.
func (c *ChainConfig) IsBerlin(num *big.Int) bool {
	return isBlockForked(c.BerlinBlock, num)
}

// IsLondon returns whether num is either equal to the London fork block or greater.
func (c *ChainConfig) IsLondon(num *big.Int) bool {
	return isBlockForked(c.LondonBlock, num)
}

// IsArrowGlacier returns whether num is either equal to the Arrow Glacier (EIP-4345) fork block or greater.
func (c *ChainConfig) IsArrowGlacier(num *big.Int) bool {
	return isBlockForked(c.ArrowGlacierBlock, num)
}

// IsGrayGlacier returns whether num is either equal to the Gray Glacier (EIP-5133) fork block or greater.
func (c *ChainConfig) IsGrayGlacier(num *big.Int) bool {
	return isBlockForked(c.GrayGlacierBlock, num)
}

// IsTerminalPoWBlock returns whether the given block is the last block of PoW stage.
func (c *ChainConfig) IsTerminalPoWBlock(parentTotalDiff *big.Int, totalDiff *big.Int) bool {
	if c.TerminalTotalDifficulty == nil {
		return false
	}
	return parentTotalDiff.Cmp(c.TerminalTotalDifficulty) < 0 && totalDiff.Cmp(c.TerminalTotalDifficulty) >= 0
}

// IsPostMerge reports whether the given block number is assumed to be post-merge.
// Here we check the MergeNetsplitBlock to allow configuring networks with a PoW or
// PoA chain for unit testing purposes.
func (c *ChainConfig) IsPostMerge(blockNum uint64, timestamp uint64) bool {
	mergedAtGenesis := c.TerminalTotalDifficulty != nil && c.TerminalTotalDifficulty.Sign() == 0
	return mergedAtGenesis ||
		c.MergeNetsplitBlock != nil && blockNum >= c.MergeNetsplitBlock.Uint64() ||
		c.ShanghaiTime != nil && timestamp >= *c.ShanghaiTime
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

// IsPrague2 returns whether time is either equal to the Prague2 fork time or greater.
// NOTE: Prague2 is a Berachain fork and must be on Ethereum's Prague fork.
func (c *ChainConfig) IsPrague2(num *big.Int, time uint64) bool {
	return c.IsPrague(num, time) && isTimestampForked(c.Berachain.Prague2.Time, time)
}

// IsPrague3 returns whether time is either equal to the Prague3 fork time or greater.
// NOTE: Prague3 is a Berachain fork and must be on Ethereum's Prague fork.
func (c *ChainConfig) IsPrague3(num *big.Int, time uint64) bool {
	return c.IsPrague(num, time) && isTimestampForked(c.Berachain.Prague3.Time, time)
}

// IsPrague4 returns whether time is either equal to the Prague4 fork time or greater.
// NOTE: Prague4 is a Berachain fork and must be on Ethereum's Prague fork.
func (c *ChainConfig) IsPrague4(num *big.Int, time uint64) bool {
	return c.IsPrague(num, time) && isTimestampForked(c.Berachain.Prague4.Time, time)
}

// IsOsaka returns whether time is either equal to the Osaka fork time or greater.
func (c *ChainConfig) IsOsaka(num *big.Int, time uint64) bool {
	return c.IsLondon(num) && isTimestampForked(c.OsakaTime, time)
}

// IsBPO1 returns whether time is either equal to the BPO1 fork time or greater.
func (c *ChainConfig) IsBPO1(num *big.Int, time uint64) bool {
	return c.IsLondon(num) && isTimestampForked(c.BPO1Time, time)
}

// IsBPO2 returns whether time is either equal to the BPO2 fork time or greater.
func (c *ChainConfig) IsBPO2(num *big.Int, time uint64) bool {
	return c.IsLondon(num) && isTimestampForked(c.BPO2Time, time)
}

// IsBPO3 returns whether time is either equal to the BPO3 fork time or greater.
func (c *ChainConfig) IsBPO3(num *big.Int, time uint64) bool {
	return c.IsLondon(num) && isTimestampForked(c.BPO3Time, time)
}

// IsBPO4 returns whether time is either equal to the BPO4 fork time or greater.
func (c *ChainConfig) IsBPO4(num *big.Int, time uint64) bool {
	return c.IsLondon(num) && isTimestampForked(c.BPO4Time, time)
}

// IsBPO5 returns whether time is either equal to the BPO5 fork time or greater.
func (c *ChainConfig) IsBPO5(num *big.Int, time uint64) bool {
	return c.IsLondon(num) && isTimestampForked(c.BPO5Time, time)
}

// IsAmsterdam returns whether time is either equal to the Amsterdam fork time or greater.
func (c *ChainConfig) IsAmsterdam(num *big.Int, time uint64) bool {
	return c.IsLondon(num) && isTimestampForked(c.AmsterdamTime, time)
}

// IsVerkle returns whether time is either equal to the Verkle fork time or greater.
func (c *ChainConfig) IsVerkle(num *big.Int, time uint64) bool {
	return c.IsLondon(num) && isTimestampForked(c.VerkleTime, time)
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

// IsEIP4762 returns whether eip 4762 has been activated at given block.
func (c *ChainConfig) IsEIP4762(num *big.Int, time uint64) bool {
	return c.IsVerkle(num, time)
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
