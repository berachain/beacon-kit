package genesis

import (
	"encoding/json"
	"fmt"
	"os"
)

// WriteGenesisToJSON writes the given Genesis to a JSON file.
func WriteGenesisToJSON(genesis interface{}, filename string) ([]byte, error) {
	// Convert the Genesis to JSON with indentation
	genesisJSON, err := json.MarshalIndent(genesis, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal genesis: %v", err)
	}

	// Write the JSON data to a file
	err = os.WriteFile(filename, genesisJSON, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to write %s: %v", filename, err)
	}

	return genesisJSON, nil
}
