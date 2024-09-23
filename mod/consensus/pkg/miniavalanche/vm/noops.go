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

package vm

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow/engine/common"
)

// VM level gossip should not be needed in mini avalanche, since
// consensus-level gossip is handled to consensus engine (above VM)
// and application level gossip should be handled by pluggable vm

// Some methods are required by Avalanche consensus engine,
// but are not really used in mini-Avalanche.
// [errDisabledMethodCalled] signals if such methods are ever called.
var errDisabledMethodCalled = errors.New("called disabled method")

func (vm *VM) AppGossip(_ context.Context, _ ids.NodeID, _ []byte) error {
	return fmt.Errorf("vm AppGossip: %w", errDisabledMethodCalled)
}

func (vm *VM) AppRequest(
	_ context.Context,
	_ ids.NodeID,
	_ uint32,
	_ time.Time,
	_ []byte,
) error {
	return fmt.Errorf("vm AppRequest: %w", errDisabledMethodCalled)
}

func (vm *VM) AppRequestFailed(
	_ context.Context,
	_ ids.NodeID,
	_ uint32,
	_ *common.AppError,
) error {
	return fmt.Errorf("vm AppRequestFailed: %w", errDisabledMethodCalled)
}

func (vm *VM) AppResponse(
	_ context.Context,
	_ ids.NodeID,
	_ uint32,
	_ []byte,
) error {
	return fmt.Errorf("vm AppResponse: %w", errDisabledMethodCalled)
}

func (vm *VM) CrossChainAppRequest(
	_ context.Context,
	_ ids.ID,
	_ uint32,
	_ time.Time,
	_ []byte,
) error {
	return fmt.Errorf("vm CrossChainAppRequest: %w", errDisabledMethodCalled)
}

func (vm *VM) CrossChainAppRequestFailed(
	_ context.Context,
	_ ids.ID,
	_ uint32,
	_ *common.AppError,
) error {
	return fmt.Errorf("vm CrossChainAppRequestFailed: %w", errDisabledMethodCalled)
}

func (vm *VM) CrossChainAppResponse(
	_ context.Context,
	_ ids.ID,
	_ uint32,
	_ []byte,
) error {
	return fmt.Errorf("vm CrossChainAppResponse: %w", errDisabledMethodCalled)
}
