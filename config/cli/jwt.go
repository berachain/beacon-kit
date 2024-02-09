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

package cli

import (
	"errors"
	"path/filepath"

	"github.com/itsdevbear/bolaris/third_party/go-ethereum/common/hexutil"
	"github.com/prysmaticlabs/prysm/v4/crypto/rand"
	"github.com/prysmaticlabs/prysm/v4/io/file"
	"github.com/spf13/cobra"
)

const (
	DefaultSecretFileName = "jwt.hex"
)

// NewGenerateJWTCommand creates a new command for generating a JWT secret.
func NewGenerateJWTCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-jwt-secret",
		Short: "Generates a new JWT authentication secret",
		Long: `This command generates a new JWT authentication secret and writes it to a file.
If no output file path is specified, it uses the default file name 
"jwt.hex" in the current directory.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fileName, err := getFilePath(cmd)
			if err != nil {
				return err
			}
			return generateAuthSecretInFile(cmd, fileName)
		},
	}
	cmd.Flags().StringP("output-path", "o", "", "Optional output file path for the JWT secret")
	return cmd
}

// getFilePath retrieves the file path for the JWT secret from the command flags.
// If no path is specified, it returns the default secret file name.
func getFilePath(cmd *cobra.Command) (string, error) {
	specifiedFilePath, err := cmd.Flags().GetString("output-path")
	if err != nil {
		return "", err
	}
	if specifiedFilePath != "" {
		return specifiedFilePath, nil
	}
	return DefaultSecretFileName, nil // Use default secret file name if no path is specified
}

// generateAuthSecretInFile writes a newly generated JWT secret to a specified file.
func generateAuthSecretInFile(cmd *cobra.Command, fileName string) error {
	var err error
	fileName, err = file.ExpandPath(fileName)
	if err != nil {
		return err
	}
	fileDir := filepath.Dir(fileName)
	exists, err := file.HasDir(fileDir)
	if err != nil {
		return err
	}
	if !exists {
		if err = file.MkdirAll(fileDir); err != nil {
			return err
		}
	}
	secret, err := generateRandomHexString()
	if err != nil {
		return err
	}
	if err = file.WriteFile(fileName, []byte(secret)); err != nil {
		return err
	}
	cmd.Printf("Successfully wrote new JSON-RPC authentication secret to: %s", fileName)
	return nil
}

// generateRandomHexString generates a random 32-byte hex string to be used as a JWT secret.
func generateRandomHexString() (string, error) {
	secret := make([]byte, 32) //nolint:gomnd // 32 bytes.
	randGen := rand.NewGenerator()
	n, err := randGen.Read(secret)
	if err != nil {
		return "", err
	} else if n <= 0 {
		return "", errors.New("rand: unexpected length")
	}
	return hexutil.Encode(secret), nil
}
