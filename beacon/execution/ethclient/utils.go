package eth

import (
	"fmt"
	"os"
	"strings"

	"cosmossdk.io/log"

	"github.com/ethereum/go-ethereum/common"
)

// loadJWTSecret reads the JWT secret from a file and returns it.
// It returns an error if the file cannot be read or if the JWT secret is not valid.
func LoadJWTSecret(filepath string, logger log.Logger) ([]byte, error) {
	// Read the file.
	data, err := os.ReadFile(filepath)
	if err != nil {
		// Return an error if the file cannot be read.
		return nil, err
	}

	// Convert the data to a JWT secret.
	jwtSecret := common.FromHex(strings.TrimSpace(string(data)))

	// Check if the JWT secret is valid.
	if len(jwtSecret) != jwtLength {
		// Return an error if the JWT secret is not valid.
		return nil, fmt.Errorf("failed to load jwt secret from %s", filepath)
	}

	logger.Info("loaded exeuction client jwt secret file", "path", filepath, "crc32")
	return jwtSecret, nil
}
