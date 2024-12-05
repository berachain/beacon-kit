// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
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
	"fmt"

	beaconhttp "github.com/attestantio/go-eth2-client/http"
	"github.com/berachain/beacon-kit/errors"
	rpcclient "github.com/cometbft/cometbft/rpc/client"
	httpclient "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"
	"github.com/rs/zerolog"
)

// ConsensusClient represents a consensus client.
type ConsensusClient struct {
	*WrappedServiceContext

	// Comet JSON-RPC client
	rpcclient.Client

	// Beacon node-api client
	BeaconKitNodeClient

	// Cancel function for the context
	cancelFunc context.CancelFunc
}

// NewConsensusClient creates a new consensus client.
func NewConsensusClient(serviceCtx *WrappedServiceContext) *ConsensusClient {
	cc := &ConsensusClient{
		WrappedServiceContext: serviceCtx,
	}

	if err := cc.Connect(context.Background()); err != nil {
		panic(err)
	}

	return cc
}

// Connect connects the consensus client to the consensus client.
func (cc *ConsensusClient) Connect(ctx context.Context) error {
	// Start by trying to get the public port for the comet JSON-RPC.
	cometPort, ok := cc.WrappedServiceContext.GetPublicPorts()["cometbft-rpc"]
	if !ok {
		panic("Couldn't find the public port for the comet JSON-RPC")
	}
	clientURL := fmt.Sprintf("http://0.0.0.0:%d", cometPort.GetNumber())
	client, err := httpclient.New(clientURL)
	if err != nil {
		return err
	}
	cc.Client = client

	// Then try to get the public port for the node API.
	nodePort, ok := cc.WrappedServiceContext.GetPublicPorts()["node-api"]
	if !ok {
		panic("Couldn't find the public port for the node API")
	}
	cancelCtx, cancel := context.WithCancel(ctx)
	cc.BeaconKitNodeClient, err = NewBeaconKitNodeClient(
		cancelCtx,
		beaconhttp.WithAddress(
			fmt.Sprintf("http://0.0.0.0:%d", nodePort.GetNumber()),
		),
		beaconhttp.WithLogLevel(zerolog.DebugLevel),
	)
	if err != nil {
		cancel()
		return err
	}
	cc.cancelFunc = cancel

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

	return res, cc.Connect(ctx)
}

// Stop stops the consensus client.
func (cc ConsensusClient) Stop(
	ctx context.Context,
) (*enclaves.StarlarkRunResult, error) {
	cc.cancelFunc()
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

// TODO: Add helpers for the beacon node-api client (converting from
// go-eth2-client types to beacon-kit consensus types).
