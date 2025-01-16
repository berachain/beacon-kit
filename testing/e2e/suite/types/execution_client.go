// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package types

import (
	"context"
	"math/big"

	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/log"
	coretypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"
)

// ExecutionClient represents an execution client.
type ExecutionClient struct {
	*WrappedServiceContext
	*JSONRPCConnection
	logger log.Logger
}

// NewExecutionClientFromServiceCtx creates a new execution client from a
// service context.
func NewExecutionClientFromServiceCtx(
	serviceCtx *WrappedServiceContext,
	logger log.Logger,
) *ExecutionClient {
	ec := &ExecutionClient{
		WrappedServiceContext: serviceCtx,
		logger:                logger,
	}

	if err := ec.Connect(); err != nil {
		panic(err)
	}

	return ec
}

func (ec *ExecutionClient) Connect() error {
	jsonRPCConn, err := NewJSONRPCConnection(ec.ServiceContext)
	if err != nil {
		return err
	}

	ec.JSONRPCConnection = jsonRPCConn
	return nil
}

func (ec ExecutionClient) Start(
	ctx context.Context,
	enclaveContext *enclaves.EnclaveContext,
) (*enclaves.StarlarkRunResult, error) {
	res, err := ec.WrappedServiceContext.Start(ctx, enclaveContext)
	if err != nil {
		return nil, err
	}

	return res, ec.Connect()
}

func (ec ExecutionClient) Stop(
	ctx context.Context,
) (*enclaves.StarlarkRunResult, error) {
	return ec.WrappedServiceContext.Stop(ctx)
}

// IsValidator returns true if the execution client is a validator.
// TODO: All nodes are validators rn.
func (ec *ExecutionClient) IsValidator() bool {
	return true
}

// WaitForLatestBlockNumber waits for the head block number to reach the target.
func (ec *ExecutionClient) WaitForLatestBlockNumber(
	ctx context.Context,
	target uint64,
) error {
	if !ec.JSONRPCConnection.isWebSocket {
		return errors.New(
			"cannot wait for block number on non-websocket connection",
		)
	}

	ch := make(chan *coretypes.Header)
	sub, err := ec.SubscribeNewHead(ctx, ch)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			sub.Unsubscribe()
			return ctx.Err()
		case header := <-ch:
			ec.logger.Info(
				"received new head block",
				"number",
				header.Number.Uint64(),
				"",
			)
			if header.Number.Uint64() >= target {
				ec.logger.Info("reached target block number 🎉", "target", target)
				sub.Unsubscribe()
				return nil
			}
		}
	}
}

// WaitForFinalizedBlockNumber waits for the finalized block number to reach the
// target block number.
func (ec *ExecutionClient) WaitForFinalizedBlockNumber(
	ctx context.Context,
	target uint64,
) error {
	finalQuery := big.NewInt(int64(rpc.FinalizedBlockNumber))
retry:
	// Retry until the head block number is at least the target.
	if err := ec.WaitForLatestBlockNumber(ctx, target+1); err != nil {
		return err
	}

	// Just to be safe, check the finalized block number again.
	finalized, err := ec.BlockByNumber(ctx, finalQuery)
	if err != nil {
		return err
	}

	// If the finalized block number is less than the target, retry.
	if finalized.Number().Uint64() < target {
		goto retry
	}
	return nil
}
