package vm

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/ava-labs/avalanchego/ids"

	"github.com/berachain/beacon-kit/mod/consensus/pkg/miniavalanche/block"
	"github.com/berachain/beacon-kit/mod/consensus/pkg/miniavalanche/encoding"
)

var genesisTime = time.Date(2024, 8, 21, 10, 17, 0, 0, time.UTC)

// Genesis is the in-memory representation of genesis
type Genesis struct {
	Validators []*Validator
}

// Base64Genesis is the genesis representation that can be fed to the node via cli flags
type Base64Genesis struct {
	Validators []Base64GenesisValidator
}

type Base64GenesisValidator struct {
	NodeID string
	Weight uint64
	Nonce  uint16
}

// ParseBase64StringToBytes is used while parsing cli flags
func ParseBase64StringToBytes(genesisStr string) ([]byte, error) {
	// Step 1: base64 string to Base64Genesis
	base64Bytes, err := base64.StdEncoding.DecodeString(genesisStr)
	if err != nil {
		return nil, fmt.Errorf("unable to decode base64 genesis content: %w", err)
	}

	base64Gen := &Base64Genesis{}
	if err := encoding.Decode(base64Bytes, base64Gen); err != nil {
		return nil, fmt.Errorf("unable to decode base64 genesis content: %w", err)
	}

	// Step 2: Base64Genesis to InMemoryGenesis
	inMemGen := &Genesis{}
	for i, v := range base64Gen.Validators {
		nodeID, err := ids.NodeIDFromString(v.NodeID)
		if err != nil {
			return nil, fmt.Errorf("unable to turn string %v, pos %d to ids.ID: %w", v.NodeID, i, err)
		}

		val, err := NewValidator(nodeID, v.Weight, v.Nonce)
		if err != nil {
			return nil, fmt.Errorf("failed building validator: %w", err)
		}

		inMemGen.Validators = append(inMemGen.Validators, val)
	}

	// Step 3: InMemoryGenesis to in memory bytes
	bytes, err := encoding.Encode(inMemGen)
	if err != nil {
		return nil, fmt.Errorf("failed encoding genesis data: %w", err)
	}

	return bytes, nil

}

// BuildBase64GenesisString is used in tools to build a genesis (to be fed to a node via cli)
func BuildBase64GenesisString(base64Gen *Base64Genesis) (string, error) {
	bytes, err := encoding.Encode(base64Gen)
	if err != nil {
		return "", fmt.Errorf("failed encoding base64 genesis: %w", err)
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
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

// process genesisBytes and from them build:
// genesis block, the first block in the chain (and only one until we unlock block creation)
// genesis validators, the validators initially responsible for the chain
func parseGenesis(genesisBytes []byte) (*block.StatelessBlock, []*Validator, error) {
	gen, err := parseInMemoryGenesis(genesisBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed parsing genesis: %w", err)
	}

	genBlk, err := buildGenesisBlock(genesisBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed building genesis block: %w", err)
	}

	return genBlk, gen.Validators, nil
}

// build a block from genesis content and keep it as
// first block in the chain.
func buildGenesisBlock(genesisBytes []byte) (*block.StatelessBlock, error) {
	// Genesis block must be parsable as a block, but genesis bytes do no encode a block
	// We create genesis block by using genesis bytes as block content
	// so that genesis block ID depends on genesisBytes
	return block.NewStatelessBlock(ids.Empty, 0, genesisTime, block.BlockContent{GenesisContent: genesisBytes})
}
