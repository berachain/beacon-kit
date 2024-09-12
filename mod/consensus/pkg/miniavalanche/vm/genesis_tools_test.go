package vm

import (
	"log"
	"testing"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/stretchr/testify/require"
)

var (
	testGenesisValidators []*Validator
	testEthGenesisBytes   []byte
)

func init() {
	// init testEthGenesisBytes
	var err error
	testEthGenesisBytes, err = DefaultEthGenesisBytes()
	if err != nil {
		log.Fatal(err)
	}

	// init testValidators
	val0, err := NewValidator(ids.GenerateTestNodeID(), uint64(999))
	if err != nil {
		log.Fatal(err)
	}
	val1, err := NewValidator(ids.GenerateTestNodeID(), uint64(1001))
	if err != nil {
		log.Fatal(err)
	}
	testGenesisValidators = []*Validator{val0, val1}
}

func TestEthGenesisEncoding(t *testing.T) {
	r := require.New(t)

	// setup genesis
	genesisData := &Base64Genesis{
		Validators: []Base64GenesisValidator{
			{
				NodeID: testGenesisValidators[0].NodeID.String(),
				Weight: testGenesisValidators[0].Weight,
			},
			{
				NodeID: testGenesisValidators[1].NodeID.String(),
				Weight: testGenesisValidators[1].Weight,
			},
		},
		EthGenesis: string(testEthGenesisBytes),
	}

	// marshal genesis
	genContent, err := BuildBase64GenesisString(genesisData)
	r.NoError(err)

	// unmarshal genesis
	parsedGenesisData, err := ParseBase64StringToBytes(genContent)
	r.NoError(err)

	_, rValidators, rGenEthData, err := parseGenesis(parsedGenesisData)
	r.NoError(err)
	r.Equal(testGenesisValidators, rValidators)
	r.Equal(testEthGenesisBytes, rGenEthData)
}
