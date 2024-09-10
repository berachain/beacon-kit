package vm

import (
	"encoding/base64"
	"fmt"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/berachain/beacon-kit/mod/consensus/pkg/miniavalanche/encoding"
)

// Base64Genesis is the genesis representation that can be fed to the node via cli flags
type Base64Genesis struct {
	Validators []Base64GenesisValidator
	EthGenesis string
}

type Base64GenesisValidator struct {
	NodeID string
	Weight uint64
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
	inMemGen := &Genesis{
		EthGenesis: []byte(base64Gen.EthGenesis),
	}
	for i, v := range base64Gen.Validators {
		nodeID, err := ids.NodeIDFromString(v.NodeID)
		if err != nil {
			return nil, fmt.Errorf("unable to turn string %v, pos %d to ids.ID: %w", v.NodeID, i, err)
		}

		val, err := NewValidator(nodeID, v.Weight)
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
