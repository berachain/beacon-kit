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

package beacon

import (
	"context"
	"encoding/json"
	"fmt"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	consensusv1 "cosmossdk.io/api/cosmos/consensus/v1"
	appmodulev2 "cosmossdk.io/core/appmodule/v2"
	"cosmossdk.io/core/registry"
	"cosmossdk.io/core/transaction"
	sdkconsensustypes "cosmossdk.io/x/consensus/types"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/genesis"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/consensus/pkg/cometbft"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/comet"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"google.golang.org/grpc"
)

const (
	// ConsensusVersion defines the current x/beacon module consensus version.
	ConsensusVersion = 1
	// ModuleName is the module name constant used in many places.
	ModuleName = "beacon"
)

var (
	_ appmodulev2.AppModule = AppModule[
		transaction.Tx, appmodulev2.ValidatorUpdate,
	]{}
	_ module.HasABCIGenesis = AppModule[
		transaction.Tx, appmodulev2.ValidatorUpdate,
	]{}
	_ module.HasABCIEndBlock = AppModule[
		transaction.Tx, appmodulev2.ValidatorUpdate,
	]{}
)

// AppModule implements an application module for the beacon module.
// It is a wrapper around the ABCIMiddleware.
type AppModule[T transaction.Tx, ValidatorUpdateT any] struct {
	ABCIMiddleware *components.ABCIMiddleware
	TxCodec        transaction.Codec[T]
	msgServer      *comet.MsgServer
}

// NewAppModule creates a new AppModule object.
func NewAppModule[T transaction.Tx, ValidatorUpdateT any](
	abciMiddleware *components.ABCIMiddleware,
	txCodec transaction.Codec[T],
	msgServer *comet.MsgServer,
) AppModule[T, ValidatorUpdateT] {
	return AppModule[T, ValidatorUpdateT]{
		ABCIMiddleware: abciMiddleware,
		TxCodec:        txCodec,
		msgServer:      msgServer,
	}
}

// Name is the name of this module.
func (am AppModule[T, ValidatorUpdateT]) Name() string {
	return ModuleName
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule[T, ValidatorUpdateT]) ConsensusVersion() uint64 {
	return ConsensusVersion
}

// RegisterInterfaces registers the module's interface types.
func (am AppModule[T, ValidatorUpdateT]) RegisterInterfaces(registry.InterfaceRegistrar) {}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule[T, ValidatorUpdateT]) IsOnePerModuleType() {}

// IsAppModule implements the appmodule.AppModule interface.
func (am AppModule[T, ValidatorUpdateT]) IsAppModule() {}

// DefaultGenesis returns default genesis state as raw bytes
// for the beacon module.
func (AppModule[T, ValidatorUpdateT]) DefaultGenesis() json.RawMessage {
	bz, err := json.Marshal(
		genesis.DefaultGenesisDeneb(),
	)
	if err != nil {
		panic(err)
	}
	return bz
}

// RegisterServices registers module services.
func (am AppModule[T, ValidatorUpdateT]) RegisterServices(registrar grpc.ServiceRegistrar) error {
	// lolololololololololololololololololololololololololololololololololololol
	sdkconsensustypes.RegisterMsgServer(registrar, am.msgServer)
	return nil
}

// ValidateGenesis performs genesis state validation for the beacon module.
func (AppModule[T, ValidatorUpdateT]) ValidateGenesis(
	_ json.RawMessage,
) error {
	return nil
}

// ExportGenesis returns the exported genesis state as raw bytes for the
// beacon module.
func (am AppModule[T, ValidatorUpdateT]) ExportGenesis(
	_ context.Context,
) (json.RawMessage, error) {
	return json.Marshal(
		&genesis.Genesis[
			*types.Deposit, *types.ExecutionPayloadHeader,
		]{},
	)
}

// InitGenesis initializes the beacon module's state from a provided genesis
// state.
func (am AppModule[T, ValidatorUpdateT]) InitGenesis(
	ctx context.Context,
	bz json.RawMessage,
) ([]ValidatorUpdateT, error) {
	return cometbft.NewConsensusEngine[T, ValidatorUpdateT](
		am.TxCodec,
		am.ABCIMiddleware,
	).InitGenesis(ctx, bz)
}

// EndBlock returns the validator set updates from the beacon state.
func (am AppModule[T, ValidatorUpdateT]) EndBlock(
	ctx context.Context,
) ([]ValidatorUpdateT, error) {
	return cometbft.NewConsensusEngine[T, ValidatorUpdateT](
		am.TxCodec,
		am.ABCIMiddleware,
	).EndBlock(ctx)
}

// proto will be sad that tendermint afk if we don't have this here
// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule[T, ValidatorUpdateT]) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: consensusv1.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Use:       "update-params-proposal [params]",
					Short:     "Submit a proposal to update consensus module params. Note: the entire params must be provided.",
					Example:   fmt.Sprintf(`%s tx consensus update-params-proposal '{ params }'`, version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "block"},
						{ProtoField: "evidence"},
						{ProtoField: "validator"},
						{ProtoField: "abci"},
					},
					GovProposal: true,
				},
			},
		},
	}
}
