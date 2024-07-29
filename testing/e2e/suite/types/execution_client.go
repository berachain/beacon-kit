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

	"github.com/berachain/beacon-kit/mod/log"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"
)

// ExecutionClient represents an execution client.
type ExecutionClient struct {
	*WrappedServiceContext
	*JSONRPCConnection
	logger log.AdvancedLogger[any]
}

// NewExecutionClientFromServiceCtx creates a new execution client from a
// service context.
func NewExecutionClientFromServiceCtx(
	serviceCtx *WrappedServiceContext,
	logger log.AdvancedLogger[any],
) *ExecutionClient {
	ec := &ExecutionClient{
		WrappedServiceContext: serviceCtx,
		logger: logger.With(
			"client-name",
			serviceCtx.GetServiceName(),
		),
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
