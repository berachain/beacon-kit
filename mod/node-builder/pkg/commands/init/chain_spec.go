package initutil

import (
	"os"
	"path/filepath"

	"github.com/berachain/beacon-kit/mod/node-builder/pkg/components"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/config/spec"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/config/viper"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

// createChainSpecCmd creates a chain spec file if one does not
// already exist.
func createChainSpecCmd(cmd *cobra.Command, args []string) error {
	homeDir, err := cmd.Flags().GetString(flags.FlagHome)
	if err != nil {
		return err
	}
	configPath := filepath.Join(homeDir, "config")
	path := filepath.Join(configPath, components.ChainSpecFileName)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Build the chain spec file
		chainspec := spec.LocalnetChainSpecData
		if err = viper.WriteStructToFile(
			configPath,
			components.ChainSpecFileName,
			chainspec,
		); err != nil {
			return err
		}
	}

	return nil
}
