package spec

import (
	"os"
	"path/filepath"

	"github.com/berachain/beacon-kit/mod/node-builder/pkg/config/spec"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

const (
	ChainSpecPath = "chainspec.toml"
)

// Commands creates a new command for managing JWT secrets.
func Commands() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "chainspec",
		Short:              "ChainSpec subcommands",
		DisableFlagParsing: false,
		RunE:               client.ValidateCmd,
	}

	cmd.AddCommand(
		NewCreateChainSpecCommand(),
	)

	return cmd
}

func NewCreateChainSpecCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generates a new chainspec config file",
		Long:  ``,
		RunE:  createChainSpec,
	}

	return cmd
}

func createChainSpec(cmd *cobra.Command, _ []string) error {
	homeDir, err := cmd.Flags().GetString(flags.FlagHome)
	if err != nil {
		return err
	}

	configPath := filepath.Join(homeDir, "config")
	chainspecFilePath := filepath.Join(configPath, ChainSpecPath)

	// when chainspec.toml does not exist create and init with default values
	if _, err := os.Stat(chainspecFilePath); os.IsNotExist(err) {
		if err := os.MkdirAll(configPath, os.ModePerm); err != nil {
			return err
		}

		if err = spec.WriteSpecToFile(
			chainspecFilePath,
			spec.LocalnetChainSpec(),
		); err != nil {
			return err
		}
	}

	return nil
}
