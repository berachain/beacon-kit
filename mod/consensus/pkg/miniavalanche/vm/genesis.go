package vm

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ava-labs/avalanchego/ids"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/consensus/pkg/miniavalanche/block"
	"github.com/berachain/beacon-kit/mod/consensus/pkg/miniavalanche/encoding"
)

var genesisTime = time.Date(2024, 8, 21, 10, 17, 0, 0, time.UTC)

// Genesis is the in-memory representation of genesis
type Genesis struct {
	Validators []*Validator
	EthGenesis []byte
}

// DefaultEthGenesisBytes implements the default genesis state for the application.
func DefaultEthGenesisBytes() ([]byte, error) {
	var (
		gen = make(map[string]json.RawMessage)
		err error
	)
	gen["beacon"], err = json.Marshal(types.DefaultGenesisDeneb())
	if err != nil {
		return nil, err
	}
	return json.Marshal(gen)
}

// process genesisBytes and from them build:
// genesis block, the first block in the chain (and only one until we unlock block creation)
// genesis validators, the validators initially responsible for the chain
// ethGenesis bytes, to be passed to the middleware
func parseGenesis(genesisBytes []byte) (*block.StatelessBlock, []*Validator, []byte, error) {
	gen, err := parseInMemoryGenesis(genesisBytes)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed parsing genesis: %w", err)
	}

	genBlk, err := buildGenesisBlock(genesisBytes)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed building genesis block: %w", err)
	}

	return genBlk, gen.Validators, gen.EthGenesis, nil
}

// parseInMemoryGenesis is used in VM initialization to retrieve genesis from its bytes
func parseInMemoryGenesis(genesisBytes []byte) (*Genesis, error) {
	inMemGen := &Genesis{}
	if err := encoding.Decode(genesisBytes, inMemGen); err != nil {
		return nil, fmt.Errorf("unable to gob decode genesis content: %w", err)
	}

	// make sure to calculate ID of every validator
	for i, v := range inMemGen.Validators {
		if err := v.initValID(); err != nil {
			return nil, fmt.Errorf("validator pos %d: %w", i, err)
		}
	}

	return inMemGen, nil
}

// build a block from genesis content and keep it as
// first block in the chain.
func buildGenesisBlock(genesisBytes []byte) (*block.StatelessBlock, error) {
	// Genesis block must be parsable as a block, but genesis bytes do no encode a block
	// We create genesis block by using genesis bytes as block content
	// so that genesis block ID depends on genesisBytes
	return block.NewStatelessBlock(ids.Empty, 0, genesisTime, block.BlockContent{GenesisContent: genesisBytes})
}
