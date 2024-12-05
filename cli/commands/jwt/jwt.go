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

package jwt

import (
	"os"
	"path/filepath"

	"github.com/berachain/beacon-kit/node-core/components"
	"github.com/berachain/beacon-kit/primitives/net/jwt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

const (
	DefaultSecretFileName = "jwt.hex"
	FlagOutputPath        = "output-path"
	FlagInputPath         = "input-path"
	ConfigFolder          = "config"
)

// Commands creates a new command for managing JWT secrets.
func Commands() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "jwt",
		Short:                      "JWT subcommands",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2, //nolint:mnd // from sdk.
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		NewGenerateJWTCommand(),
		NewValidateJWTCommand(),
	)

	return cmd
}

// NewGenerateJWTCommand creates a new command for generating a JWT secret.
func NewGenerateJWTCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generates a new JWT authentication secret",
		Long: `This command generates a new JWT authentication secret and
writes it to a file. If no output file path is specified, it uses the default
file name "jwt.hex" in the current directory.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			// Get the file path from the command flags.
			outputPath, err := getFilePath(cmd, FlagOutputPath)
			if err != nil {
				return err
			}

			return generateAuthSecretInFile(cmd, outputPath)
		},
	}
	cmd.Flags().StringP(
		FlagOutputPath, "o", "", "Optional output file path for the JWT secret")
	return cmd
}

// NewValidateJWTCommand creates a new command for validating a JWT secret.
func NewValidateJWTCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validates a JWT secret conforms to Engine-RPC requirements",
		Long: `This command validates a JWT secret by checking if the JWT secret
is formatted properly. If no output file path is specified, it uses the default
file name "jwt.hex" in the current directory.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			// Get the file path from the command flags.
			inputPath, err := getFilePath(cmd, FlagInputPath)
			if err != nil {
				return err
			}

			return validateJWTSecret(cmd, inputPath)
		},
	}

	cmd.Flags().StringP(
		FlagInputPath, "i", "", "Optional input file path for the JWT secret",
	)
	return cmd
}

// getFilePath retrieves the file path for the JWT secret from the command flag.
// If no path is specified, it returns the default secret file name.
func getFilePath(cmd *cobra.Command, path string) (string, error) {
	specifiedFilePath, err := cmd.Flags().GetString(path)
	if err != nil {
		return "", err
	}
	if specifiedFilePath != "" {
		return specifiedFilePath, nil
	}

	// If no path is specified, try to get the cosmos client context and use
	// the configured home directory to write the secret to the default file
	// name.
	clientCtx, ok := cmd.Context().
		Value(client.ClientContextKey).(*client.Context)
	if !ok {
		return "", ErrNoClientCtx
	}
	specifiedFilePath = filepath.Join(
		clientCtx.HomeDir, ConfigFolder, DefaultSecretFileName,
	)

	// Use default secret file name if no path is specified
	return specifiedFilePath, nil
}

// generateAuthSecretInFile writes a newly generated JWT secret
// to a specified file.
func generateAuthSecretInFile(cmd *cobra.Command, fileName string) error {
	var err error
	fs := afero.NewOsFs()
	fileDir := filepath.Dir(fileName)
	exists, err := afero.DirExists(fs, fileDir)
	if err != nil {
		return err
	}

	if !exists {
		if err = fs.MkdirAll(fileDir, os.ModePerm); err != nil {
			return err
		}
	}

	secret, err := jwt.NewRandom()
	if err != nil {
		return err
	}

	if err = afero.WriteFile(
		fs, fileName, []byte(secret.Hex()), os.ModePerm,
	); err != nil {
		return err
	}

	cmd.Printf(
		"Successfully wrote new JSON-RPC authentication secret to: %s",
		fileName,
	)
	return nil
}

func validateJWTSecret(cmd *cobra.Command, filePath string) error {
	_, err := components.LoadJWTFromFile(filePath)
	if err != nil {
		return err
	}

	cmd.Printf("Successfully validated JWT secret")
	return nil
}
