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
	"cosmossdk.io/log"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/services"
)

// ExecutionClient represents an execution client.
type ExecutionClient struct {
	*services.ServiceContext
	*JSONRPCConnection
	logger log.Logger
}

// NewExecutionClientFromServiceCtx creates a new execution client from a
// service context.
func NewExecutionClientFromServiceCtx(
	serviceCtx *services.ServiceContext,
	logger log.Logger,
) (*ExecutionClient, error) {
	jsonRPCConn, err := NewJSONRPCConnection(serviceCtx)
	if err != nil {
		return nil, err
	}

	return &ExecutionClient{
		ServiceContext:    serviceCtx,
		JSONRPCConnection: jsonRPCConn,
		logger: logger.With(
			"client-name",
			serviceCtx.GetServiceName(),
		),
	}, nil
}

// IsValidator returns true if the execution client is a validator.
// TODO: All nodes are validators rn.
func (c *ExecutionClient) IsValidator() bool {
	return true
}
