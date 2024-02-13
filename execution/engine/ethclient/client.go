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

package ethclient

import (
	"context"
	"errors"
	"net/url"
	"sync/atomic"
	"time"

	"cosmossdk.io/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/itsdevbear/bolaris/io/jwt"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
)

// Eth1Client is a struct that holds the Ethereum 1 client and its configuration.
type Eth1Client struct {
	GethRPCClient
	*ethclient.Client

	logger               log.Logger
	isConnected          atomic.Bool
	chainID              uint64
	jwtSecret            *jwt.Secret
	startupRetryInterval time.Duration
	jwtRefreshInterval   time.Duration
	healthCheckInterval  time.Duration
	dialURL              *url.URL
}

// NewEth1Client creates a new Ethereum 1 client with the provided context and options.
func NewEth1Client(opts ...Option) (*Eth1Client, error) {
	c := &Eth1Client{}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// Start the powchain service's main event loop.
func (s *Eth1Client) Start(ctx context.Context) {
	// Attempt an intial connection.
	s.setupExecutionClientConnection(ctx)

	// We will spin up the execution client connection in a loop until it is connected.
	for !s.isConnected.Load() {
		// If we enter this loop, the above connection attempt failed.
		s.logger.Info("Waiting for connection to execution client...", "dial-url", s.dialURL.String())
		s.tryConnectionAfter(ctx, s.startupRetryInterval)
	}

	// If we reached this point, the execution client is connected so we can start
	// the health check & jwt refresh loops.
	go s.healthCheckLoop(ctx)
	go s.jwtRefreshLoop(ctx)
}

// IsConnected returns the connection status of the Ethereum 1 client.
func (s *Eth1Client) IsConnected() bool {
	return s.isConnected.Load()
}

// NewPayloadV3 calls the engine_newPayloadV3 method via JSON-RPC.
func (s *Eth1Client) NewPayloadV3(
	ctx context.Context, payload *enginev1.ExecutionPayloadDeneb,
	versionedHashes []common.Hash, parentBlockRoot *common.Hash,
) (*enginev1.PayloadStatus, error) {
	result := &enginev1.PayloadStatus{}
	if err := s.GethRPCClient.CallContext(
		ctx, result, NewPayloadMethodV3, payload, versionedHashes, parentBlockRoot,
	); err != nil {
		return nil, s.handleRPCError(err)
	}
	return result, nil
}

// NewPayloadV2 calls the engine_newPayloadV2 method via JSON-RPC.
func (s *Eth1Client) NewPayloadV2(
	ctx context.Context, payload *enginev1.ExecutionPayloadCapella,
) (*enginev1.PayloadStatus, error) {
	result := &enginev1.PayloadStatus{}
	if err := s.GethRPCClient.CallContext(
		ctx, result, NewPayloadMethodV2, payload,
	); err != nil {
		return nil, s.handleRPCError(err)
	}
	return result, nil
}

// ForkchoiceUpdatedV3 calls the engine_forkchoiceUpdatedV3 method via JSON-RPC.
func (s *Eth1Client) ForkchoiceUpdatedV3(
	ctx context.Context, state *enginev1.ForkchoiceState, attrs *enginev1.PayloadAttributesV3,
) (*ForkchoiceUpdatedResponse, error) {
	return s.forkchoiceUpdateCall(ctx, ForkchoiceUpdatedMethodV3, state, attrs)
}

// ForkchoiceUpdatedV2 calls the engine_forkchoiceUpdatedV2 method via JSON-RPC.
func (s *Eth1Client) ForkchoiceUpdatedV2(
	ctx context.Context, state *enginev1.ForkchoiceState, attrs *enginev1.PayloadAttributesV2,
) (*ForkchoiceUpdatedResponse, error) {
	return s.forkchoiceUpdateCall(ctx, ForkchoiceUpdatedMethodV2, state, attrs)
}

// forkchoiceUpdateCall is a helper function to call to any version of the forkchoiceUpdated
// method.
func (s *Eth1Client) forkchoiceUpdateCall(
	ctx context.Context, method string, state *enginev1.ForkchoiceState, attrs any,
) (*ForkchoiceUpdatedResponse, error) {
	result := &ForkchoiceUpdatedResponse{}

	if err := s.GethRPCClient.CallContext(
		ctx, result, method, state, attrs,
	); err != nil {
		return nil, s.handleRPCError(err)
	}

	if result.Status == nil {
		return nil, ErrNilResponse
	} else if result.ValidationError != "" {
		s.logger.Error(
			"Got validation error in forkChoiceUpdated",
			"err", errors.New(result.ValidationError),
		)
	}

	return result, nil
}

// GetPayloadV3 calls the engine_getPayloadV3 method via JSON-RPC.
func (s *Eth1Client) GetPayloadV3(
	ctx context.Context, payloadID enginev1.PayloadIDBytes,
) (*enginev1.ExecutionPayloadDenebWithValueAndBlobsBundle, error) {
	result := &enginev1.ExecutionPayloadDenebWithValueAndBlobsBundle{}
	if err := s.GethRPCClient.CallContext(
		ctx, result, GetPayloadMethodV3, payloadID,
	); err != nil {
		return nil, s.handleRPCError(err)
	}
	return result, nil
}

// GetPayloadV2 calls the engine_getPayloadV2 method via JSON-RPC.
func (s *Eth1Client) GetPayloadV2(
	ctx context.Context, payloadID enginev1.PayloadIDBytes,
) (*enginev1.ExecutionPayloadCapellaWithValue, error) {
	result := &enginev1.ExecutionPayloadCapellaWithValue{}
	if err := s.GethRPCClient.CallContext(
		ctx, result, GetPayloadMethodV2, payloadID,
	); err != nil {
		return nil, s.handleRPCError(err)
	}
	return result, nil
}

// ExecutionBlockByHash fetches an execution engine block by hash by calling
// eth_blockByHash via JSON-RPC.
func (s *Eth1Client) ExecutionBlockByHash(ctx context.Context, hash common.Hash, withTxs bool,
) (*enginev1.ExecutionBlock, error) {
	result := &enginev1.ExecutionBlock{}
	err := s.GethRPCClient.CallContext(
		ctx, result, BlockByHashMethod, hash, withTxs)
	return result, s.handleRPCError(err)
}

// ExecutionBlockByNumber fetches an execution engine block by number by calling
// eth_getBlockByNumber via JSON-RPC.
func (s *Eth1Client) ExecutionBlockByNumber(ctx context.Context, num rpc.BlockNumber, withTxs bool,
) (*enginev1.ExecutionBlock, error) {
	result := &enginev1.ExecutionBlock{}
	err := s.GethRPCClient.CallContext(
		ctx, result, BlockByNumberMethod, num, withTxs)
	return result, s.handleRPCError(err)
}
