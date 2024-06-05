package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	nodebuilder "github.com/berachain/beacon-kit/mod/node-builder/pkg/node-builder/v2"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func StartCmd[T types.Tx](creator nodebuilder.AppCreator[T], addFlagsFn func(startCmd *cobra.Command)) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the application",
		RunE: func(cmd *cobra.Command, args []string) error {
			// serverCtx := server.GetServerContextFromCmd(cmd)
			// sa := simapp.NewSimApp(serverCtx.Logger, serverCtx.Viper)
			// am := sa.App.AppManager
			// serverCfg := cometbft.Config{CmtConfig: serverCtx.Config, ConsensusAuthority: sa.GetConsensusAuthority()}

			// cometServer := cometbft.NewCometBFTServer[transaction.Tx](
			// 	am,
			// 	sa.GetStore(),
			// 	sa.GetLogger(),
			// 	serverCfg,
			// 	txCodec,
			// )
			ctx := cmd.Context()
			ctx, cancelFn := context.WithCancel(ctx)
			g, _ := errgroup.WithContext(ctx)
			g.Go(func() error {
				sigCh := make(chan os.Signal, 1)
				signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
				sig := <-sigCh
				cancelFn()
				cmd.Printf("caught %s signal\n", sig.String())

				if err := cometServer.Stop(ctx); err != nil {
					cmd.PrintErrln("failed to stop servers:", err)
				}
				return nil
			})

			if err := cometServer.Start(ctx); err != nil {
				return fmt.Errorf("failed to start servers: %w", err)
			}
			return g.Wait()
		},
	}

	addFlagsFn(cmd)
	return cmd
}
