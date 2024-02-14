// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	storetypes "cosmossdk.io/store/types"
	"github.com/cometbft/cometbft/crypto/merkle"
	"github.com/cometbft/cometbft/libs/log"
	cometOs "github.com/cometbft/cometbft/libs/os"
	lproxy "github.com/cometbft/cometbft/light/proxy"
	lrpc "github.com/cometbft/cometbft/light/rpc"
	cmtypes "github.com/cometbft/cometbft/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/itsdevbear/bolaris/light"
	"github.com/itsdevbear/bolaris/light/provider"
	"github.com/spf13/cobra"
)

func runProxy(cmd *cobra.Command, args []string) error {
	// Initialize logger.
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	var option log.Option
	if logLevel == "info" {
		option, _ = log.AllowLevel("info")
	} else {
		option, _ = log.AllowLevel("debug")
	}
	logger = log.NewFilter(logger, option)

	chainID := args[0]
	logger.Info("Creating client...", "chainID", chainID)

	tl, err := cmd.Flags().GetString(trustLevel)
	if err != nil {
		return err
	}
	sequential, err := cmd.Flags().GetBool(sequential)
	if err != nil {
		return err
	}
	trustedHeight, err := cmd.Flags().GetInt64(trustedHeight)
	if err != nil {
		return err
	}
	trustedHash, err := cmd.Flags().GetBytesHex(trustedHash)
	if err != nil {
		return err
	}
	trustingPeriod, err := cmd.Flags().GetDuration(trustingPeriod)
	if err != nil {
		return err
	}
	maxOpenConnections, err := cmd.Flags().GetInt(maxOpenConnections)
	if err != nil {
		return err
	}
	witnesses, err := cmd.Flags().GetString(witnessAddrsJoined)
	if err != nil {
		return err
	}

	// witness addresses are the other addresses other than the primary node
	// used to cross check and verify the primary node's headers and etc.
	var witnessesAddrs []string
	if witnessAddrsJoined != "" {
		witnessesAddrs = strings.Split(witnesses, ",")
	}

	if len(witnessesAddrs) == 0 {
		witnessesAddrs = []string{"http://localhost:26657"}
	}

	primaryAddr, err := cmd.Flags().GetString(primaryAddr)
	if err != nil {
		return err
	}

	if primaryAddr == "" {
		primaryAddr = "tcp://localhost:26657"
	}

	client, err := light.NewClient(
		logger,
		chainID,
		trustingPeriod,
		trustedHeight,
		trustedHash,
		tl,
		sequential,
		primaryAddr,
		witnessesAddrs,
		dir,
		NewConfirmationFunc(cmd),
	)
	if err != nil {
		return err
	}
	config := light.InitServerConfig(maxOpenConnections)
	listenAddr_ := "tcp://localhost:26658"
	// Create the proxy server.
	// this is a tendermint light client proxy server
	p, err := lproxy.NewProxy(
		client,
		listenAddr_,
		primaryAddr,
		config,
		logger,
		func(c *lrpc.Client) {
			c.RegisterOpDecoder(
				storetypes.ProofOpIAVLCommitment, storetypes.CommitmentOpDecoder,
			)
			c.RegisterOpDecoder(
				storetypes.ProofOpSimpleMerkleCommitment, storetypes.CommitmentOpDecoder,
			)
		},
		lrpc.KeyPathFn(MerkleKeyPathFn()),
	)
	if err != nil {
		return err
	}

	// Stop upon receiving SIGTERM or CTRL-C.
	cometOs.TrapSignal(logger, func() {
		p.Listener.Close()
	})

	c := lrpc.NewClient(p.Client, client, func(c *lrpc.Client) {
		c.RegisterOpDecoder(
			storetypes.ProofOpIAVLCommitment, storetypes.CommitmentOpDecoder,
		)
		c.RegisterOpDecoder(
			storetypes.ProofOpSimpleMerkleCommitment, storetypes.CommitmentOpDecoder,
		)
	},
		lrpc.KeyPathFn(MerkleKeyPathFn()),
	)

	// querier := NewQuerier(p.Client)
	cosmosProvider := provider.CosmosProvider{
		RPCClient: c,
	}

	c.Start()
	go ListenForNewBlocks(cmd.Context(), &cosmosProvider)

	// Unmarshal raw bytes to proto.Message
	logger.Info("Starting proxy...", "laddr", listenAddr)
	// Start the proxy server.
	go func() {
		if err = p.ListenAndServe(); errors.Is(err, http.ErrServerClosed) {
			// Error starting or closing listener:
			logger.Error("proxy ListenAndServe", "err", err)
		} else {
			logger.Error("proxy server closed", "error", err)
		}
	}()
	<-cmd.Context().Done()
	return err
}

func ListenForNewBlocks(ctx context.Context, prov *provider.CosmosProvider) {
	sub, err := prov.RPCClient.Subscribe(ctx, "subscriber", cmtypes.EventQueryNewBlockHeader.String())
	if err != nil {
		panic(err)
	}
	for {
		select {
		case <-ctx.Done():
			return
		case <-sub:
			resA, err := prov.RPCClient.ABCIInfo(ctx)
			if err != nil {
				fmt.Println(err)
				continue
			}
			resp, _, err := prov.RunGRPCQuery(
				context.Background(),
				"/store/evm/key",
				[]byte("fc_finalized"),
				resA.Response.LastBlockHeight-1,
				false,
			)
			if err != nil {
				fmt.Println("Error querying the store:", err)
				continue
			}
			fmt.Println("\033[32mFinalized Block on the Execution Client:", common.BytesToHash(resp.Value), "\033[0m")
		default:
			continue
		}
	}
}

const expectedMatchLength = 2

// DefaultMerkleKeyPathFn creates a function used to generate merkle key paths
// from a path string and a key. This is the default used by the cosmos SDK.
// This merkle key paths are required when verifying /abci_query calls.
func MerkleKeyPathFn() lrpc.KeyPathFunc {
	// regexp for extracting store name from /abci_query path
	storeNameRegexp := regexp.MustCompile(`\/store\/(.+)\/key`)

	return func(path string, key []byte) (merkle.KeyPath, error) {
		matches := storeNameRegexp.FindStringSubmatch(path)
		if len(matches) != expectedMatchLength {
			return nil, fmt.Errorf("can't find store name in %s using %s", path, storeNameRegexp)
		}
		storeName := matches[1]

		kp := merkle.KeyPath{}
		kp = kp.AppendKey([]byte(storeName), merkle.KeyEncodingURL)
		kp = kp.AppendKey(key, merkle.KeyEncodingURL)
		return kp, nil
	}
}

func NewConfirmationFunc(cmd *cobra.Command) func(string) bool {
	return func(action string) bool {
		cmd.Println(action)
		scanner := bufio.NewScanner(os.Stdin)
		for {
			scanner.Scan()
			response := scanner.Text()
			switch response {
			case "y", "Y":
				return true
			case "n", "N":
				return false
			default:
				cmd.Println("please input 'Y' or 'n' and press ENTER")
			}
		}
	}
}

type Querier struct {
	// The client to use for querying.
	client *lrpc.Client
}

func NewQuerier(client *lrpc.Client) *Querier {
	return &Querier{client: client}
}
