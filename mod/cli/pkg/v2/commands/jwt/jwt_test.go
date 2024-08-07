// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package jwt_test

import (
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	jwt "github.com/berachain/beacon-kit/mod/cli/pkg/v2/commands/jwt"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func Test_NewGenerateJWTCommand(t *testing.T) {
	t.Run(
		"command should be available and have correct use",
		func(t *testing.T) {
			cmd := jwt.NewGenerateJWTCommand()
			require.Equal(t, "generate", cmd.Use)
		},
	)

	t.Run("should create proper file in current directory", func(t *testing.T) {
		cmd := jwt.NewGenerateJWTCommand()
		cmd.SetArgs([]string{"--output-path", jwt.DefaultSecretFileName})
		require.NoError(t, cmd.Execute())

		// We check the file has the contents we expect.
		checkAuthFileIntegrity(t, jwt.DefaultSecretFileName)

		require.NoError(t, os.RemoveAll(jwt.DefaultSecretFileName))
	})

	t.Run("should create proper file in specified folder", func(t *testing.T) {
		customOutput := filepath.Join("data", "jwt.hex")
		cmd := jwt.NewGenerateJWTCommand()
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
		cmd := jwt.NewGenerateJWTCommand()
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
		cmd := jwt.NewGenerateJWTCommand()
		cmd.SetArgs([]string{"--output-path", tempFile.Name()})
		require.NoError(t, cmd.Execute())

		// Check the file has been overridden with the new content
		checkAuthFileIntegrity(t, tempFile.Name())

		require.NoError(t, os.RemoveAll(tempFile.Name()))
	})
}

func checkAuthFileIntegrity(tb testing.TB, fPath string) {
	tb.Helper()
	fs := afero.NewOsFs()
	fileInfo, err := fs.Stat(fPath)
	require.NoError(tb, err)
	require.NotNil(tb, fileInfo)

	enc, err := afero.ReadFile(fs, fPath)
	require.NoError(tb, err)
	var decoded = make([]byte, hex.DecodedLen(len(enc[2:])))
	_, err = hex.Decode(decoded, enc[2:])
	require.NoError(tb, err)
	require.Len(tb, decoded, 32)
}
