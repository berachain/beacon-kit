package genesis

import (
	"encoding/json"
	"os"
)

// writeGenesisToJSON writes the given Genesis to a JSON file.
func writeGenesisToJSON(genesis Genesis, filename string) error {
	// Convert the Genesis to JSON with indentation
	genesisJSON, err := json.MarshalIndent(genesis, "", "  ")
	if err != nil {
		return err
	}

	// Write the JSON data to a file
	if err = os.WriteFile(filename, genesisJSON, 0644); err != nil {
		return err
	}

	return nil
}
