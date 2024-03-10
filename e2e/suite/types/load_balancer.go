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
	"strings"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/services"
)

// LoadBalancer represents a group of eth JSON-RPC endpoints
// behind an NGINX load balancer.
type LoadBalancer struct {
	*JSONRPCConnection
}

func NewLoadBalancer(
	serviceCtx *services.ServiceContext,
) (*LoadBalancer, error) {
	conn := &JSONRPCConnection{
		isWebSocket: false,
	}

	// Start by trying to get the public port for the JSON-RPC WebSocket
	jsonRPC, ok := serviceCtx.GetPublicPorts()["http"]
	if !ok {
		return nil, ErrPublicPortNotFound
	}
	// Split the string to get the port
	str := strings.Split(jsonRPC.String(), ":")
	if len(str) == 0 {
		return nil, ErrPublicPortNotFound
	}
	str = strings.Split(str[1], "/")
	if len(str) == 0 {
		return nil, ErrPublicPortNotFound
	}
	port := str[0]
	ethClient, err := ethclient.Dial(
		"http://0.0.0.0" + ":" + port,
	)
	if err != nil {
		return nil, err
	}
	conn.Client = ethClient
	return &LoadBalancer{JSONRPCConnection: conn}, nil
}
