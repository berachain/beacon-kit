package commands

import (
	"os"
	"path/filepath"

	"github.com/berachain/beacon-kit/mod/cli/pkg/commands/start"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/types"
	"github.com/spf13/cobra"

	"cosmossdk.io/core/transaction"
	"cosmossdk.io/log"
	serverv2 "cosmossdk.io/server/v2"
)

// DefaultCommandConfig adds a start command to the root command.
func DefaultCommandConfig[NodeT types.Node[T], T transaction.Tx](
	rootCmd *cobra.Command,
	appCreator serverv2.AppCreator[NodeT, T],
	logger log.Logger,
	bkCommands []*cobra.Command,
	components ...serverv2.ServerComponent[NodeT, T],
) (serverv2.CLIConfig, error) {

	// TODO: this is weird, but is how cosmos does it. pls fix later
	server := serverv2.NewServer(logger, components...)
	flags := server.StartFlags()
	startCmd := start.NewStartCmd[NodeT, T](
		appCreator,
		server,
		flags,
	)
	cmds := server.CLICommands()
	cmds.Commands = append(cmds.Commands, bkCommands...)
	cmds.Commands = append(cmds.Commands, startCmd)

	return cmds, nil
}

// AddCommands adds the start command to the root command and sets the
// server context
func AddCommands[NodeT types.Node[T], T transaction.Tx](
	rootCmd *cobra.Command,
	newApp serverv2.AppCreator[NodeT, T],
	logger log.Logger,
	cmdConfig serverv2.CLIConfig,
	components ...serverv2.ServerComponent[NodeT, T],
) error {
	server := serverv2.NewServer(logger, components...)
	originalPersistentPreRunE := rootCmd.PersistentPreRunE
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		home, err := cmd.Flags().GetString(serverv2.FlagHome)
		if err != nil {
			return err
		}

		err = configHandle(server, home, cmd)
		if err != nil {
			return err
		}

		if rootCmd.PersistentPreRun != nil {
			rootCmd.PersistentPreRun(cmd, args)
			return nil
		}

		return originalPersistentPreRunE(cmd, args)
	}

	rootCmd.AddCommand(cmdConfig.Commands...)
	return nil
}

// configHandle writes the default config to the home directory if it does not exist and sets the server context
func configHandle[NodeT types.Node[T], T transaction.Tx](
	s *serverv2.Server[NodeT, T],
	home string,
	cmd *cobra.Command,
) error {
	if _, err := os.Stat(filepath.Join(home, "config")); os.IsNotExist(err) {
		if err = s.WriteConfig(filepath.Join(home, "config")); err != nil {
			return err
		}
	}

	viper, err := serverv2.ReadConfig(filepath.Join(home, "config"))
	if err != nil {
		return err
	}
	viper.Set(serverv2.FlagHome, home)
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	log, err := serverv2.NewLogger(viper, cmd.OutOrStdout())
	if err != nil {
		return err
	}

	return serverv2.SetCmdServerContext(cmd, viper, log)
}
