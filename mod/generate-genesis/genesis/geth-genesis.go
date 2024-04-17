package genesis

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"log"
	"math/big"
	"os"
)

type Genesis struct {
	alloc types.GenesisAlloc
}

// NewGenesis creates a new Genesis instance with an initialized alloc field.
func NewGenesis() *Genesis {
	return &Genesis{
		alloc: make(types.GenesisAlloc),
	}
}

// TO-DO : Remove core package and use the internal package "primtives" instead
func (g *Genesis) AddAccount(genesis *core.Genesis, address common.Address, balance *big.Int) {
	g.alloc[address] = types.Account{
		Balance: balance,
	}
	genesis.Alloc = g.alloc
}

func (g *Genesis) AddPredeploy(genesis *core.Genesis, address common.Address, code []byte, storage map[common.Hash]common.Hash, balance *big.Int, nonce uint64) {
	g.alloc[address] = types.Account{
		Code:    code,
		Storage: storage,
		Balance: balance,
		Nonce:   nonce,
	}
	genesis.Alloc = g.alloc
}

func (g *Genesis) ToGethGenesis() *core.Genesis {
	chainConfig := &params.ChainConfig{
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

	return &core.Genesis{
		Config:     chainConfig,
		Nonce:      0x0000000000000000,
		Timestamp:  0x0,
		ExtraData:  common.FromHex("0x0000000000000000000000000000000000000000000000000000000000000000658bdf435d810c91414ec09147daa6db624063790000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
		GasLimit:   0x1c9c380,
		Difficulty: big.NewInt(0x0),
		Alloc:      g.alloc,
		Coinbase:   common.Address{},
	}
}
func (g *Genesis) WriteFileToJSON(genesis *core.Genesis, fileName string) ([]byte, error) {
	// Convert the genesis to JSON with indentation
	genesisJSON, err := json.MarshalIndent(genesis, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal genesis: %v", err)
		return nil, err
	}
	// Write the JSON data to a file
	err = os.WriteFile(fileName, genesisJSON, 0644)
	if err != nil {
		log.Fatalf("Failed to write genesis.json: %v", err)
		return nil, err
	}
	return genesisJSON, nil
}
