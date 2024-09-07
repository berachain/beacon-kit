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

// Some methods are required by the interfaces required by Avalanche consensus engine,
// but should never apply to mini-Avalanche case. [errDisabledMethodCalled] signals if
// such methods are ever called
var errDisabledMethodCalled = errors.New("called disabled method")

func (vm *VM) AppGossip(_ context.Context, _ ids.NodeID, _ []byte) error {
	return fmt.Errorf("vm AppGossip: %w", errDisabledMethodCalled)
}

func (vm *VM) AppRequest(_ context.Context, _ ids.NodeID, _ uint32, _ time.Time, _ []byte) error {
	return fmt.Errorf("vm AppRequest: %w", errDisabledMethodCalled)
}

func (vm *VM) AppRequestFailed(_ context.Context, _ ids.NodeID, _ uint32, _ *common.AppError) error {
	return fmt.Errorf("vm AppRequestFailed: %w", errDisabledMethodCalled)
}

func (vm *VM) AppResponse(_ context.Context, _ ids.NodeID, _ uint32, _ []byte) error {
	return fmt.Errorf("vm AppResponse: %w", errDisabledMethodCalled)
}

func (vm *VM) CrossChainAppRequest(_ context.Context, _ ids.ID, _ uint32, _ time.Time, _ []byte) error {
	return fmt.Errorf("vm CrossChainAppRequest: %w", errDisabledMethodCalled)
}

func (vm *VM) CrossChainAppRequestFailed(_ context.Context, _ ids.ID, _ uint32, _ *common.AppError) error {
	return fmt.Errorf("vm CrossChainAppRequestFailed: %w", errDisabledMethodCalled)
}

func (vm *VM) CrossChainAppResponse(_ context.Context, _ ids.ID, _ uint32, _ []byte) error {
	return fmt.Errorf("vm CrossChainAppResponse: %w", errDisabledMethodCalled)
}
