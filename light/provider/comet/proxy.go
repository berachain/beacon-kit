package comet

// import (
// 	"errors"
// 	"net/http"

// 	storetypes "cosmossdk.io/store/types"

// 	cometOs "github.com/cometbft/cometbft/libs/os"
// 	lproxy "github.com/cometbft/cometbft/light/proxy"
// 	lrpc "github.com/cometbft/cometbft/light/rpc"

// 	"github.com/berachain/beacon-kit/light/provider/comet/types"
// )

// func StartProxy(c Config) error {
// 	client, err := NewClient(
// 		c.Logger,
// 		c.ChainID,
// 		c.TrustingPeriod,
// 		c.TrustedHeight,
// 		c.TrustedHash,
// 		c.TrustLevel,
// 		c.Sequential,
// 		c.PrimaryAddr,
// 		c.WitnessesAddrs,
// 		c.Directory,
// 		c.ConfirmationFunc,
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	serverCfg := initServerConfig(c.MaxOpenConnections)

// 	opts := []lrpc.Option{
// 		lrpc.KeyPathFn(lrpc.DefaultMerkleKeyPathFn()),
// 		func(c *lrpc.Client) {
// 			c.RegisterOpDecoder(
// 				storetypes.ProofOpIAVLCommitment, storetypes.CommitmentOpDecoder,
// 			)
// 			c.RegisterOpDecoder(
// 				storetypes.ProofOpSimpleMerkleCommitment, storetypes.CommitmentOpDecoder,
// 			)
// 		},
// 	}
// 	proxy, err := lproxy.NewProxy(
// 		client, c.ListeningAddr, c.PrimaryAddr, serverCfg, c.Logger, opts...,
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	// Stop upon receiving SIGTERM or CTRL-C.
// 	cometOs.TrapSignal(c.Logger, func() {
// 		proxy.Listener.Close()
// 	})

// 	c.Logger.Info("Starting proxy...", "laddr", c.ListeningAddr)
// 	go func() {
// 		if err = proxy.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
// 			// Error starting or closing listener:
// 			c.Logger.Error(types.ListenAndServeError, "err", err)
// 		}
// 	}()

// 	return nil
// }
