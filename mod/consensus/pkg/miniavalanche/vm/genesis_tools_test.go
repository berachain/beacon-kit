package vm

import (
	"embed"
	"encoding/json"
	"log"
	"testing"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed test-eth-genesis.json
	testEthGenesisContent embed.FS

	testValidators      []*Validator
	testEthGenesisBytes []byte
)

func init() {
	// init testEthGenesisBytes
	var err error
	testEthGenesisBytes, err = testEthGenesisContent.ReadFile("test-eth-genesis.json")
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
	testValidators = []*Validator{val0, val1}
}

func TestEthGenesisEncoding(t *testing.T) {
	r := require.New(t)

	// setup genesis
	genesisData := &Base64Genesis{
		Validators: []Base64GenesisValidator{
			{
				NodeID: testValidators[0].NodeID.String(),
				Weight: testValidators[0].Weight,
			},
			{
				NodeID: testValidators[1].NodeID.String(),
				Weight: testValidators[1].Weight,
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
	r.Equal(testValidators, rValidators)
	r.Equal(testEthGenesisBytes, rGenEthData)

	// check eth genesis content
	var genesisState map[string]json.RawMessage
	r.NoError(json.Unmarshal(testEthGenesisBytes, &genesisState))

	beaconMsg, found := genesisState["beacon"]
	r.True(found)

	var data map[string]json.RawMessage
	r.NoError(json.Unmarshal([]byte(beaconMsg), &data))

	_, found = data["fork_version"]
	r.True(found)
	_, found = data["deposits"]
	r.True(found)
	_, found = data["execution_payload_header"]
	r.True(found)
}
