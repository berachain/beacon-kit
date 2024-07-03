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
	cmtruntime "github.com/berachain/beacon-kit/mod/runtime/pkg/cometbft"
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
	_ appmodulev2.HasEndBlocker = AppModule[
		transaction.Tx, appmodulev2.ValidatorUpdate,
	]{}

	_ appmodulev2.HasUpdateValidators = AppModule[
		transaction.Tx, appmodulev2.ValidatorUpdate,
	]{}
)

// AppModule implements an application module for the beacon module.
// It is a wrapper around the ABCIMiddleware.
type AppModule[T transaction.Tx, ValidatorUpdateT any] struct {
	abciMiddleware  *components.ABCIMiddleware
	txCodec         transaction.Codec[T]
	msgServer       *cmtruntime.MsgServer
	consensusEngine *cometbft.ConsensusEngine[T, ValidatorUpdateT]
}

// NewAppModule creates a new AppModule object.
func NewAppModule[T transaction.Tx, ValidatorUpdateT any](
	abciMiddleware *components.ABCIMiddleware,
	txCodec transaction.Codec[T],
	msgServer *cmtruntime.MsgServer,
) AppModule[T, ValidatorUpdateT] {
	return AppModule[T, ValidatorUpdateT]{
		abciMiddleware: abciMiddleware,
		txCodec:        txCodec,
		msgServer:      msgServer,
		consensusEngine: cometbft.NewConsensusEngine[T, ValidatorUpdateT](
			txCodec,
			abciMiddleware,
		),
	}
}

// Name is the name of this module.
func (AppModule[_, _]) Name() string {
	return ModuleName
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule[_, _]) ConsensusVersion() uint64 {
	return ConsensusVersion
}

// RegisterInterfaces registers the module's interface types.
func (AppModule[_, _]) RegisterInterfaces(registry.InterfaceRegistrar) {}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (AppModule[_, _]) IsOnePerModuleType() {}

// IsAppModule implements the appmodule.AppModule interface.
func (AppModule[_, _]) IsAppModule() {}

// DefaultGenesis returns default genesis state as raw bytes
// for the beacon module.
func (AppModule[_, _]) DefaultGenesis() json.RawMessage {
	bz, err := json.Marshal(genesis.DefaultGenesisDeneb())
	if err != nil {
		panic(err)
	}
	return bz
}

// RegisterServices registers module services.
func (am AppModule[_, _]) RegisterServices(
	registrar grpc.ServiceRegistrar,
) error {
	// lolololololololololololololololololololololololololololololololololololol
	sdkconsensustypes.RegisterMsgServer(registrar, am.msgServer)
	return nil
}

// ValidateGenesis performs genesis state validation for the beacon module.
func (AppModule[_, _]) ValidateGenesis(_ json.RawMessage) error {
	return nil
}

// ExportGenesis returns the exported genesis state as raw bytes for the
// beacon module.
func (AppModule[_, _]) ExportGenesis(
	_ context.Context,
) (json.RawMessage, error) {
	return json.Marshal(&genesis.Genesis[
		*types.Deposit, *types.ExecutionPayloadHeader,
	]{})
}

// InitGenesis initializes the beacon module's state from a provided genesis
// state.
func (am AppModule[T, ValidatorUpdateT]) InitGenesis(
	ctx context.Context,
	bz json.RawMessage,
) ([]ValidatorUpdateT, error) {
	return am.consensusEngine.InitGenesis(ctx, bz)
}

// EndBlock returns the validator set updates from the beacon state.
func (am AppModule[T, ValidatorUpdateT]) EndBlock(
	ctx context.Context,
) error {
	return am.consensusEngine.EndBlock(ctx)
}

// EndBlock returns the validator set updates from the beacon state.
func (am AppModule[T, ValidatorUpdateT]) UpdateValidators(
	ctx context.Context,
) ([]ValidatorUpdateT, error) {
	return am.consensusEngine.UpdateValidators(ctx)
}

// proto will be sad that tendermint afk if we don't have this here
// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (AppModule[_, _]) AutoCLIOptions() *autocliv1.ModuleOptions {
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
