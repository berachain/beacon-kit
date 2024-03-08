// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
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

package cli_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/berachain/beacon-kit/config/cli"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"
)

func Test_NewGenerateJWTCommand(t *testing.T) {
	t.Run(
		"command should be available and have correct use",
		func(t *testing.T) {
			cmd := cli.NewGenerateJWTCommand()
			require.Equal(t, "generate-jwt-secret", cmd.Use)
		},
	)

	t.Run("should create proper file in current directory", func(t *testing.T) {
		cmd := cli.NewGenerateJWTCommand()
		cmd.SetArgs([]string{"--output-path", cli.DefaultSecretFileName})
		require.NoError(t, cmd.Execute())

		// We check the file has the contents we expect.
		checkAuthFileIntegrity(t, cli.DefaultSecretFileName)

		require.NoError(t, os.RemoveAll(cli.DefaultSecretFileName))
	})

	t.Run("should create proper file in specified folder", func(t *testing.T) {
		customOutput := filepath.Join("data", "jwt.hex")
		cmd := cli.NewGenerateJWTCommand()
		cmd.SetArgs([]string{"--output-path", customOutput})
		require.NoError(t, cmd.Execute())

		// We check the file has the contents we expect.
		checkAuthFileIntegrity(t, customOutput)

		require.NoError(t, os.RemoveAll(filepath.Dir(customOutput)))
	})

	t.Run("creates proper file in nested specified folder", func(t *testing.T) {
		rootDirectory := "data"
		customOutputPath := filepath.Join(
			rootDirectory,
			"nest",
			"nested",
			"jwt.hex",
		)
		cmd := cli.NewGenerateJWTCommand()
		cmd.SetArgs([]string{"--output-path", customOutputPath})
		require.NoError(t, cmd.Execute())

		// We check the file has the contents we expect.
		checkAuthFileIntegrity(t, customOutputPath)

		require.NoError(t, os.RemoveAll(rootDirectory))
	})

	t.Run("should override existing file when flag is set", func(t *testing.T) {
		// Create a temporary file to simulate an existing file
		tempFile, err := os.CreateTemp("", "existing_jwt.hex")
		require.NoError(t, err)
		defer os.Remove(tempFile.Name()) // clean up

		// Write some content to the file to simulate an existing JWT
		_, err = tempFile.WriteString("not_a_jwt_secret")
		require.NoError(t, err)
		tempFile.Close()

		// Execute the command with the --force flag to override the existing
		// file
		cmd := cli.NewGenerateJWTCommand()
		cmd.SetArgs([]string{"--output-path", tempFile.Name()})
		require.NoError(t, cmd.Execute())

		// Check the file has been overridden with the new content
		checkAuthFileIntegrity(t, tempFile.Name())

		require.NoError(t, os.RemoveAll(tempFile.Name()))
	})
}

func checkAuthFileIntegrity(t testing.TB, fPath string) {
	fileInfo, err := os.Stat(fPath)
	require.NoError(t, err)
	require.NotNil(t, fileInfo)

	enc, err := os.ReadFile(fPath) // Updated to use os.ReadFile directly
	require.NoError(t, err)
	decoded, err := hexutil.Decode(string(enc))
	require.NoError(t, err)
	require.Len(t, decoded, 32)
}
