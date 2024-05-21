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

package beacon

import (
	"context"
	"encoding/json"
	"fmt"

	appmodulev2 "cosmossdk.io/core/appmodule/v2"
	"cosmossdk.io/core/registry"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/state/deneb"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/components"
	"github.com/cosmos/cosmos-sdk/types/module"
)

const (
	// ConsensusVersion defines the current x/beacon module consensus version.
	ConsensusVersion = 1
	// ModuleName is the module name constant used in many places.
	ModuleName = "beacon"
)

var (
	_ appmodulev2.AppModule  = AppModule{}
	_ module.HasABCIGenesis  = AppModule{}
	_ module.HasABCIEndBlock = AppModule{}
)

// AppModule implements an application module for the evm module.
type AppModule struct {
	*components.BeaconKitRuntime
}

// NewAppModule creates a new AppModule object.
func NewAppModule(
	runtime *components.BeaconKitRuntime,
) AppModule {
	return AppModule{
		BeaconKitRuntime: runtime,
	}
}

// Name is the name of this module.
func (am AppModule) Name() string {
	return ModuleName
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return ConsensusVersion }

// RegisterInterfaces registers the module's interface types.
func (am AppModule) RegisterInterfaces(registry.InterfaceRegistrar) {}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule) IsOnePerModuleType() {}

// IsAppModule implements the appmodule.AppModule interface.
func (am AppModule) IsAppModule() {}

// DefaultGenesis returns default genesis state as raw bytes
// for the beacon module.
func (AppModule) DefaultGenesis() json.RawMessage {
	defaultGenesis, err := deneb.DefaultBeaconState()
	if err != nil {
		panic(err)
	}

	bz, err := json.Marshal(defaultGenesis)
	if err != nil {
		panic(err)
	}
	return bz
}

// ValidateGenesis performs genesis state validation for the evm module.
func (AppModule) ValidateGenesis(
	bz json.RawMessage,
) error {
	// TODO: this is bad
	data := new(deneb.BeaconState)
	if err := json.Unmarshal(bz, data); err != nil {
		return err
	}

	seenValidators := make(map[[48]byte]struct{})
	for _, validator := range data.Validators {
		if _, ok := seenValidators[validator.Pubkey]; ok {
			return fmt.Errorf(
				"duplicate pubkey in genesis state: %x",
				validator.Pubkey,
			)
		}
		seenValidators[validator.Pubkey] = struct{}{}
	}
	return nil
}

// ExportGenesis returns the exported genesis state as raw bytes for the evm
// module.
func (am AppModule) ExportGenesis(
	_ context.Context,
) (json.RawMessage, error) {
	return json.Marshal(&deneb.BeaconState{})
}
