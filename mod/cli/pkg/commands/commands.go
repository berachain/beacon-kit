package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"

	"cosmossdk.io/core/transaction"
	"cosmossdk.io/log"
	serverv2 "cosmossdk.io/server/v2"
)

// DefaultCommandConfig adds a start command to the root command.
func DefaultCommandConfig(
	rootCmd *cobra.Command,
	appCreator serverv2.AppCreator[transaction.Tx],
	logger log.Logger,
	components ...serverv2.ServerComponent[transaction.Tx],
) (serverv2.CLIConfig, error) {

	server := serverv2.NewServer(logger, components...)
	flags := server.StartFlags()

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Run the application",
		RunE: func(cmd *cobra.Command, args []string) error {
			v := serverv2.GetViperFromCmd(cmd)
			l := serverv2.GetLoggerFromCmd(cmd)

			for _, startFlags := range flags {
				if err := v.BindPFlags(startFlags); err != nil {
					return err
				}
			}

			if err := v.BindPFlags(cmd.Flags()); err != nil {
				return err
			}

			app := appCreator(l, v)

			if err := server.Init(app, v, l); err != nil {
				return err
			}

			srvConfig := serverv2.Config{StartBlock: true}
			ctx := cmd.Context()
			ctx = context.WithValue(ctx, serverv2.ServerContextKey, srvConfig)
			ctx, cancelFn := context.WithCancel(ctx)
			go func() {
				sigCh := make(chan os.Signal, 1)
				signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
				sig := <-sigCh
				cancelFn()
				cmd.Printf("caught %s signal\n", sig.String())

				if err := server.Stop(ctx); err != nil {
					cmd.PrintErrln("failed to stop servers:", err)
				}
			}()

			if err := server.Start(ctx); err != nil {
				return fmt.Errorf("failed to start servers: %w", err)
			}

			return nil
		},
	}

	cmds := server.CLICommands()
	cmds.Commands = append(cmds.Commands, startCmd)

	return cmds, nil
}

// AddCommands adds the start command to the root command and sets the
// server context
func AddCommands(
	rootCmd *cobra.Command,
	newApp serverv2.AppCreator[transaction.Tx],
	logger log.Logger,
	components ...serverv2.ServerComponent[transaction.Tx],
) error {
	cmds, err := DefaultCommandConfig(rootCmd, newApp, logger, components...)
	if err != nil {
		return err
	}

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

	rootCmd.AddCommand(cmds.Commands...)
	return nil
}

// configHandle writes the default config to the home directory if it does not exist and sets the server context
func configHandle(s *serverv2.Server, home string, cmd *cobra.Command) error {
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
