// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
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

package types

import (
	"context"
	"fmt"

	"github.com/berachain/beacon-kit/mod/errors"
	rpcclient "github.com/cometbft/cometbft/rpc/client"
	httpclient "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"
)

// ConsensusClient represents a consensus client.
type ConsensusClient struct {
	*WrappedServiceContext
	rpcclient.Client
}

// NewConsensusClient creates a new consensus client.
func NewConsensusClient(serviceCtx *WrappedServiceContext) *ConsensusClient {
	cc := &ConsensusClient{
		WrappedServiceContext: serviceCtx,
	}

	if err := cc.Connect(); err != nil {
		panic(err)
	}

	return cc
}

// Connect connects the consensus client to the consensus client.
func (cc *ConsensusClient) Connect() error {
	// Start by trying to get the public port for the JSON-RPC WebSocket
	port, ok := cc.WrappedServiceContext.GetPublicPorts()["cometbft-rpc"]
	if !ok {
		panic("Couldn't find the public port for the JSON-RPC WebSocket")
	}
	clientURL := fmt.Sprintf("http://0.0.0.0:%d", port.GetNumber())
	client, err := httpclient.New(clientURL)
	if err != nil {
		return err
	}
	cc.Client = client
	return nil
}

// Start starts the consensus client.
func (cc ConsensusClient) Start(
	ctx context.Context,
	enclaveContext *enclaves.EnclaveContext,
) (*enclaves.StarlarkRunResult, error) {
	res, err := cc.WrappedServiceContext.Start(ctx, enclaveContext)
	if err != nil {
		return nil, err
	}

	return res, cc.Connect()
}

// Stop stops the consensus client.
func (cc ConsensusClient) Stop(
	ctx context.Context,
) (*enclaves.StarlarkRunResult, error) {
	return cc.WrappedServiceContext.Stop(ctx)
}

// GetPubKey returns the public key of the validator running on this node.
func (cc ConsensusClient) GetPubKey(ctx context.Context) ([]byte, error) {
	res, err := cc.Client.Status(ctx)
	if err != nil {
		return nil, err
	} else if res.ValidatorInfo.PubKey == nil {
		return nil, errors.New("node public key is nil")
	}

	return res.ValidatorInfo.PubKey.Bytes(), nil
}

// GetConsensusPower returns the consensus power of the node.
func (cc ConsensusClient) GetConsensusPower(
	ctx context.Context,
) (uint64, error) {
	res, err := cc.Client.Status(ctx)
	if err != nil {
		return 0, err
	}

	//#nosec:G701 // VotingPower won't ever be negative.
	return uint64(res.ValidatorInfo.VotingPower), nil
}

// IsActive returns true if the node is an active validator.
func (cc ConsensusClient) IsActive(ctx context.Context) (bool, error) {
	res, err := cc.Client.Status(ctx)
	if err != nil {
		return false, err
	}

	return res.ValidatorInfo.VotingPower > 0, nil
}
