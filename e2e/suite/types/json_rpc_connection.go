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
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/services"
)

// JSONRPCConnection wraps an Ethereum client connection.
// It provides JSON-RPC communication with an Ethereum node.
type JSONRPCConnection struct {
	*ethclient.Client
	isWebSocket bool
}

// NewJSONRPCConnection creates a new JSON-RPC connection.
func NewJSONRPCConnection(
	serviceCtx *services.ServiceContext,
) (*JSONRPCConnection, error) {
	var (
		err  error
		conn = &JSONRPCConnection{
			isWebSocket: true,
		}
	)

	// Start by trying to get the public port for the JSON-RPC WebSocket
	port, ok := serviceCtx.GetPublicPorts()["eth-json-rpc-ws"]
	if !ok {
		// If the WebSocket port isn't available, try the HTTP port
		port, ok = serviceCtx.GetPublicPorts()["eth-json-rpc"]
		if !ok {
			return nil, ErrPublicPortNotFound
		}
		conn.isWebSocket = false
	}

	prefix := "http://"
	if conn.isWebSocket {
		prefix = "ws://"
	}

	if conn.Client, err = ethclient.Dial(
		fmt.Sprintf("%s://0.0.0.0:%d", prefix, port.GetNumber()),
	); err != nil {
		return nil, err
	}

	return conn, nil
}

// IsWebSocket returns true if the connection is a WebSocket.
func (c *JSONRPCConnection) IsWebSocket() bool {
	return c.isWebSocket
}
