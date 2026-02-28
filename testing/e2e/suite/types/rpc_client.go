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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package types

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/services"
)

// RPCClient wraps an eth JSON-RPC connection to an execution layer node.
type RPCClient struct {
	*JSONRPCConnection

	url string
}

func NewRPCClient(
	serviceCtx *services.ServiceContext,
) (*RPCClient, error) {
	var (
		err  error
		conn = &JSONRPCConnection{}
	)

	// Start by trying to get the public port for the JSON-RPC endpoint.
	port, ok := serviceCtx.GetPublicPorts()["eth-json-rpc"]
	if !ok {
		return nil, ErrPublicPortNotFound
	}

	// Then try to connect to the JSON-RPC Endpoint.
	url := fmt.Sprintf("http://0.0.0.0:%d", port.GetNumber())
	if conn.Client, err = ethclient.Dial(url); err != nil {
		return nil, err
	}

	return &RPCClient{
		JSONRPCConnection: conn,
		url:               url,
	}, nil
}

// URL returns the URL of the RPC client.
func (rc *RPCClient) URL() string {
	return rc.url
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
