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
		conn = &JSONRPCConnection{}
	)

	// If the WebSocket port isn't available, try the HTTP port
	port, ok := serviceCtx.GetPublicPorts()["eth-json-rpc"]
	if !ok {
		return nil, ErrPublicPortNotFound
	}

	if conn.Client, err = ethclient.Dial(
		fmt.Sprintf("http://://0.0.0.0:%d", port.GetNumber()),
	); err != nil {
		return nil, err
	}

	return conn, nil
}

// IsWebSocket returns true if the connection is a WebSocket.
func (c *JSONRPCConnection) IsWebSocket() bool {
	return c.isWebSocket
}
