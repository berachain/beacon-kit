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
	"fmt"
	"net/http"
	"time"

	beraclient "github.com/berachain/beacon-kit/gethlib/ethclient"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/services"
)

// ExecutionClient represents an execution client.
type ExecutionClient struct {
	*beraclient.Client

	serviceCtx *services.ServiceContext
}

// NewExecutionClientFromServiceCtx creates a new execution client from a
// service context.
func NewExecutionClient(serviceCtx *services.ServiceContext) *ExecutionClient {
	return &ExecutionClient{
		serviceCtx: serviceCtx,
	}
}

func (ec *ExecutionClient) Start(ctx context.Context) error {
	port, ok := ec.serviceCtx.GetPublicPorts()["eth-json-rpc"]
	if !ok {
		return ErrPublicPortNotFound
	}

	ethClient, err := DialWithPooling(fmt.Sprintf("http://://0.0.0.0:%d", port.GetNumber()))
	if err != nil {
		return err
	}

	ec.Client = beraclient.Wrap(ethClient)
	return nil
}

// DialWithPooling creates an ethclient with an HTTP transport configured
// for high connection reuse, avoiding ephemeral port exhaustion under
// heavy concurrent load.
func DialWithPooling(url string) (*ethclient.Client, error) {
	const (
		poolMaxIdleConns    = 256
		poolIdleTimeoutSecs = 90
	)
	httpClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        poolMaxIdleConns,
			MaxIdleConnsPerHost: poolMaxIdleConns,
			IdleConnTimeout:     poolIdleTimeoutSecs * time.Second,
		},
	}
	rpcClient, err := rpc.DialOptions(
		context.Background(), url, rpc.WithHTTPClient(httpClient),
	)
	if err != nil {
		return nil, err
	}
	return ethclient.NewClient(rpcClient), nil
}
