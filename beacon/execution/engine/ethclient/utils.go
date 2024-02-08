// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package eth

import (
	"fmt"
	"os"
	"strings"

	"cosmossdk.io/log"

	"github.com/itsdevbear/bolaris/third_party/go-ethereum/common"
)

// loadJWTSecret reads the JWT secret from a file and returns it.
// It returns an error if the file cannot be read or if the JWT secret is not valid.
func LoadJWTSecret(filepath string, logger log.Logger) ([]byte, error) {
	// Read the file.
	//#nosec:G304 // false positive.
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
