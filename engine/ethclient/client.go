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

package eth

import (
	"context"
	"errors"
	"net/url"
	"sync/atomic"
	"time"

	"cosmossdk.io/log"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/itsdevbear/bolaris/third_party/go-ethereum/common"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/execution"
	payloadattribute "github.com/prysmaticlabs/prysm/v4/consensus-types/payload-attribute"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
)

const (
	// jwtLength is the length of the JWT token.
	jwtLength = 32
)

// Eth1Client is a struct that holds the Ethereum 1 client and its configuration.
type Eth1Client struct {
	logger log.Logger
	*ethclient.Client

	connectedETH1        atomic.Bool
	chainID              uint64
	jwtSecret            [32]byte
	startupRetryInterval time.Duration
	jwtRefreshInterval   time.Duration
	healthCheckInterval  time.Duration
	dialURL              *url.URL
}

// NewEth1Client creates a new Ethereum 1 client with the provided context and options.
func NewEth1Client(ctx context.Context, opts ...Option) (*Eth1Client, error) {
	c := &Eth1Client{}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	c.Start(ctx) // TODO: move this so it is on the cmd.Context.
	return c, nil
}

// Start the powchain service's main event loop.
func (s *Eth1Client) Start(ctx context.Context) {
	// Attempt an intial connection.
	s.setupExecutionClientConnection(ctx)

	// We will spin up the execution client connection in a loop until it is connected.
	for !s.ConnectedETH1() {
		// If we enter this loop, the above connection attempt failed.
		s.logger.Info("Waiting for connection to execution client...", "dial-url", s.dialURL.String())
		s.tryConnectionAfter(ctx, s.startupRetryInterval)
	}

	// If we reached this point, the execution client is connected so we can start
	// the health check & jwt refresh loops.
	go s.healthCheckLoop(ctx)
	go s.jwtRefreshLoop(ctx)
}

// RawClient returns the raw Ethereum 1 client.
func (s *Eth1Client) RawClient() *rpc.Client {
	return s.Client.Client()
}

// ConnectedETH1 returns the connection status of the Ethereum 1 client.
func (s *Eth1Client) ConnectedETH1() bool {
	return s.connectedETH1.Load()
}

// NewPayloadV2 calls the engine_newPayloadV2 method via JSON-RPC.
func (s *Eth1Client) NewPayloadV2(
	ctx context.Context, payload *enginev1.ExecutionPayloadCapella,
) (*enginev1.PayloadStatus, error) {
	result := &enginev1.PayloadStatus{}
	if err := s.RawClient().CallContext(ctx, result, "engine_newPayloadV2", payload); err != nil {
		return nil, s.handleRPCError(err)
	}
	return result, nil
}

// NewPayloadV3 calls the engine_newPayloadV3 method via JSON-RPC.
func (s *Eth1Client) NewPayloadV3(
	ctx context.Context, payload *enginev1.ExecutionPayloadDeneb,
	versionedHashes []common.Hash, parentBlockRoot *common.Hash,
) (*enginev1.PayloadStatus, error) {
	result := &enginev1.PayloadStatus{}
	if err := s.RawClient().CallContext(
		ctx, result, "engine_newPayloadV3", payload, versionedHashes, parentBlockRoot,
	); err != nil {
		return nil, s.handleRPCError(err)
	}
	return result, nil
}

// ForkchoiceUpdatedV2 calls the engine_forkchoiceUpdatedV2 method via JSON-RPC.
func (s *Eth1Client) ForkchoiceUpdatedV2(
	ctx context.Context, state *enginev1.ForkchoiceState, attrs payloadattribute.Attributer,
) (*enginev1.PayloadIDBytes, []byte, error) {
	return s.forkchoiceUpdatedCall(ctx, "engine_forkchoiceUpdatedV2", state, attrs)
}

// ForkchoiceUpdatedV3 calls the engine_forkchoiceUpdatedV3 method via JSON-RPC.
func (s *Eth1Client) ForkchoiceUpdatedV3(
	ctx context.Context, state *enginev1.ForkchoiceState, attrs payloadattribute.Attributer,
) (*enginev1.PayloadIDBytes, []byte, error) {
	return s.forkchoiceUpdatedCall(ctx, "engine_forkchoiceUpdatedV3", state, attrs)
}

// forkchoiceUpdatedCall is a helper function to call the forkchoiceUpdated method via JSON-RPC.
func (s *Eth1Client) forkchoiceUpdatedCall(
	ctx context.Context, method string,
	state *enginev1.ForkchoiceState, attrs payloadattribute.Attributer,
) (*enginev1.PayloadIDBytes, []byte, error) {
	result := &ForkchoiceUpdatedResponse{}
	if attrs == nil {
		return nil, nil, errors.New("nil payload attributer")
	}

	if err := s.RawClient().CallContext(
		ctx, result, method, state, attrs,
	); err != nil {
		return nil, nil, s.handleRPCError(err)
	}

	if result.Status == nil {
		return nil, nil, execution.ErrNilResponse
	}

	return result.PayloadID, result.Status.GetLatestValidHash(), nil
}

// ExecutionBlockByHash fetches an execution engine block by hash by calling
// eth_blockByHash via JSON-RPC.
func (s *Eth1Client) ExecutionBlockByHash(ctx context.Context, hash common.Hash, withTxs bool,
) (*enginev1.ExecutionBlock, error) {
	result := &enginev1.ExecutionBlock{}
	err := s.RawClient().CallContext(
		ctx, result, "eth_getBlockByHash", hash, withTxs)
	return result, s.handleRPCError(err)
}

// ExecutionBlockByNumber fetches an execution engine block by number by calling
// eth_getBlockByNumber via JSON-RPC.
func (s *Eth1Client) ExecutionBlockByNumber(ctx context.Context, num rpc.BlockNumber, withTxs bool,
) (*enginev1.ExecutionBlock, error) {
	result := &enginev1.ExecutionBlock{}
	err := s.RawClient().CallContext(
		ctx, result, "eth_getBlockByNumber", num, withTxs)
	return result, s.handleRPCError(err)
}
