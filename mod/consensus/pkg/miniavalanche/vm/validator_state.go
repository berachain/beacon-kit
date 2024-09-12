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
	"fmt"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow/validators"
)

var (
	_ validators.State           = (*VM)(nil)
	_ validators.SubnetConnector = (*VM)(nil)
)

// Interface needed by ProposerVM to handle network congestion

func (vm *VM) GetMinimumHeight(ctx context.Context) (uint64, error) {
	// For the time being just return height of the latest accepted block
	return vm.GetCurrentHeight(ctx)
}

func (vm *VM) GetCurrentHeight(context.Context) (uint64, error) {
	blk, err := vm.state.GetBlock(vm.state.GetLastAccepted())
	if err != nil {
		return 0, fmt.Errorf("internal error, can't find last accepted block: %w", err)
	}
	return blk.Height(), nil
}

func (vm *VM) GetSubnetID(_ context.Context, chainID ids.ID) (ids.ID, error) {
	// We only have one subnet and one chain
	if chainID != vm.chainCtx.ChainID {
		return ids.Empty, fmt.Errorf("chainID %s unknow, only known chain ID is %s", chainID, vm.chainCtx.ChainID)
	}
	return vm.chainCtx.SubnetID, nil
}

func (vm *VM) GetValidatorSet(
	_ context.Context,
	_ uint64,
	_ ids.ID,
) (map[ids.NodeID]*validators.GetValidatorOutput, error) {
	return vm.validators.GetMap(vm.chainCtx.SubnetID), nil
}

func (vm *VM) ConnectedSubnet(_ context.Context, _ ids.NodeID, subnetID ids.ID) error {
	if subnetID != vm.chainCtx.SubnetID {
		return fmt.Errorf("unknown subnetID %v, only known subnetID %v", subnetID, vm.chainCtx.SubnetID)
	}
	return nil
}
