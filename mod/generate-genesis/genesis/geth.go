package genesis

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
)

type GethGenesis struct {
	Alloc       types.GenesisAlloc
	CoreGenesis *core.Genesis
}

func (g *GethGenesis) AddAccount(address common.Address, balance *big.Int) {
	g.Alloc[address] = types.Account{
		Balance: balance,
	}
	g.CoreGenesis.Alloc = g.Alloc
}

func (g *GethGenesis) AddPredeploy(address common.Address, code []byte, balance *big.Int, nonce uint64) {
	g.Alloc[address] = types.Account{
		Code:    code,
		Balance: balance,
		Nonce:   nonce,
	}
	g.CoreGenesis.Alloc = g.Alloc
}

func (g *GethGenesis) ToJSON(filename string) error {
	_, err := WriteGenesisToJSON(g.CoreGenesis, filename)
	return err
}

func (g *GethGenesis) initializeGethGenesis() *GethGenesis {
	gg := &GethGenesis{
		Alloc: make(types.GenesisAlloc),
		CoreGenesis: &core.Genesis{
			Config:     &params.ChainConfig{},
			Nonce:      0x0000000000000000,
			Timestamp:  0x0,
			ExtraData:  []byte{},
			GasLimit:   0x1c9c380,
			Difficulty: big.NewInt(0x0),
			Alloc:      make(types.GenesisAlloc),
			Coinbase:   common.Address{},
		},
	}
	return gg
}

func (g *GethGenesis) populateChainConfig(gg *GethGenesis) {
	gg.CoreGenesis.Config = &params.ChainConfig{
		ChainID:                       big.NewInt(80087), // 80087 is the chain ID for Berachain
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

func (g *GethGenesis) ToGethGenesis() *GethGenesis {
	gg := g.initializeGethGenesis()
	g.populateChainConfig(gg)
	return gg
}
