package initcli

import (
	"encoding/json"
	"os"

	"github.com/berachain/beacon-kit/mod/cli/pkg/v2/flags"
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/spf13/cobra"
)

func Command[
	ConsensusParamsT interface {
		json.Marshaler
		Default() ConsensusParamsT
	},
	GenesisStateT interface {
		json.Marshaler
		Default() GenesisStateT
	},
]() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize the beacon node",
		Long:  "Initialize the beacon node",
		RunE: func(cmd *cobra.Command, args []string) error {
			var cfg cmtcfg.Config

			// TODO: recovery mnemonic

			// TODO: intial height

			// TODO: consensus key

			// TODO: initialize node validator

			// Build genesis state
			genesisFilePath := cfg.GenesisFile()
			overwrite, err := cmd.Flags().GetBool(flags.FlagOverwrite)
			if err != nil {
				overwrite = false
			}
			// Check if the genesis file exists and we're not overwriting it
			if _, err := os.Stat(
				genesisFilePath,
			); !overwrite && !os.IsNotExist(err) {
				return ErrGenesisFileExists
			}

			// TODO: do we need more genesis data than state????
			var cp ConsensusParamsT
			var gs GenesisStateT
			genesis := &Genesis[ConsensusParamsT, GenesisStateT]{
				State:           gs.Default(),
				ConsensusParams: cp.Default(),
			}
			return genesis.Save(genesisFilePath)
		},
	}

	//nolint:lll // it's honestly more clear this way
	cmd.Flags().BoolP(flags.FlagOverwrite, "o", flags.DefaultOverwrite, flags.OverwriteDescription)
	cmd.Flags().BoolP(flags.FlagRecover, "r", flags.DefaultRecover, flags.RecoverDescription)

	return cmd
}
