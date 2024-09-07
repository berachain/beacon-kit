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
