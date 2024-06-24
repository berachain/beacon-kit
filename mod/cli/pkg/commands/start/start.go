package start

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"cosmossdk.io/core/transaction"
	serverv2 "cosmossdk.io/server/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewStartCmd[NodeT serverv2.AppI[T], T transaction.Tx](
	appCreator serverv2.AppCreator[NodeT, T],
	server *serverv2.Server[NodeT, T],
	flags []*pflag.FlagSet,
) *cobra.Command {
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

	return startCmd
}
