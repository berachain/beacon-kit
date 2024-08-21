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

package components

import (
	"cosmossdk.io/core/appmodule/v2"
	"github.com/berachain/beacon-kit/mod/async/pkg/dispatcher"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/consensus/pkg/cometbft"
	consruntimetypes "github.com/berachain/beacon-kit/mod/consensus/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	engineclient "github.com/berachain/beacon-kit/mod/execution/pkg/client"
	"github.com/berachain/beacon-kit/mod/execution/pkg/deposit"
	execution "github.com/berachain/beacon-kit/mod/execution/pkg/engine"
	"github.com/berachain/beacon-kit/mod/log/pkg/phuslu"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/signer"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/services/version"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/async"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/service"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	depositdb "github.com/berachain/beacon-kit/mod/storage/pkg/deposit"
	"github.com/berachain/beacon-kit/mod/storage/pkg/manager"
	"github.com/berachain/beacon-kit/mod/storage/pkg/pruner"
)

/* -------------------------------------------------------------------------- */
/*                                  Services                                  */
/* -------------------------------------------------------------------------- */

type (
	// ConsensusMiddleware is a type alias for the consensus middleware.
	ConsensusMiddleware = cometbft.Middleware[
		*AttestationData,
		*SlashingInfo,
		*SlotData,
	]

	// DBManager is a type alias for the database manager.
	DBManager = manager.DBManager

	// EngineClient is a type alias for the engine client.
	EngineClient = engineclient.EngineClient[
		*ExecutionPayload,
		*PayloadAttributes,
	]

	// EngineClient is a type alias for the engine client.
	ExecutionEngine = execution.Engine[
		*ExecutionPayload,
		*PayloadAttributes,
		PayloadID,
		engineprimitives.Withdrawals,
	]

	// ReportingService is a type alias for the reporting service.
	ReportingService = version.ReportingService
)

/* -------------------------------------------------------------------------- */
/*                                    Types                                   */
/* -------------------------------------------------------------------------- */

type (
	// AttestationData is a type alias for the attestation data.
	AttestationData = types.AttestationData

	// Context is a type alias for the transition context.
	Context = transition.Context

	// Deposit is a type alias for the deposit.
	Deposit = types.Deposit

	// DepositContract is a type alias for the deposit contract.
	DepositContract = deposit.WrappedBeaconDepositContract[
		*Deposit,
		WithdrawalCredentials,
	]

	// DepositStore is a type alias for the deposit store.
	DepositStore = depositdb.KVStore[*Deposit]

	// Eth1Data is a type alias for the eth1 data.
	Eth1Data = types.Eth1Data

	// ExecutionPayload type aliases.
	ExecutionPayload       = types.ExecutionPayload
	ExecutionPayloadHeader = types.ExecutionPayloadHeader

	// Fork is a type alias for the fork.
	Fork = types.Fork

	// ForkData is a type alias for the fork data.
	ForkData = types.ForkData

	// Genesis is a type alias for the Genesis type.
	Genesis = types.Genesis[
		*Deposit,
		*ExecutionPayloadHeader,
	]

	// Logger is a type alias for the logger.
	Logger = phuslu.Logger

	// LoggerConfig is a type alias for the logger config.
	LoggerConfig = phuslu.Config

	// SlotData is a type alias for the incoming slot.
	SlotData = consruntimetypes.SlotData[
		*AttestationData,
		*SlashingInfo,
	]

	// LegacyKey type alias to LegacyKey used for LegacySinger construction.
	LegacyKey = signer.LegacyKey

	// PayloadAttributes is a type alias for the payload attributes.
	PayloadAttributes = engineprimitives.PayloadAttributes[*Withdrawal]

	// PayloadID is a type alias for the payload ID.
	PayloadID = engineprimitives.PayloadID

	// SlashingInfo is a type alias for the slashing info.
	SlashingInfo = types.SlashingInfo

	// Validator is a type alias for the validator.
	Validator = types.Validator

	// Validators is a type alias for the validators.
	Validators = types.Validators

	// ValidatorUpdate is a type alias for the validator update.
	ABCIValidatorUpdate = appmodule.ValidatorUpdate

	// ValidatorUpdate is a type alias for the validator update.
	ValidatorUpdate = transition.ValidatorUpdate

	// ValidatorUpdates is a type alias for the validator updates.
	ValidatorUpdates = transition.ValidatorUpdates

	// Withdrawal is a type alias for the engineprimitives withdrawal.
	Withdrawal = engineprimitives.Withdrawal

	// Withdrawals is a type alias for the engineprimitives withdrawals.
	Withdrawals = engineprimitives.Withdrawals

	// WithdrawalCredentials is a type alias for the withdrawal credentials.
	WithdrawalCredentials = types.WithdrawalCredentials
)

/* -------------------------------------------------------------------------- */
/*                                   Messages                                 */
/* -------------------------------------------------------------------------- */

// Events.
//

type (

	// GenesisDataReceivedEvent is a type alias for the genesis data received
	// event.
	GenesisDataReceivedEvent = async.Event[*Genesis]

	// GenesisDataProcessedEvent is a type alias for the genesis data processed
	// event.
	GenesisDataProcessedEvent = async.Event[transition.ValidatorUpdates]

	// NewSlotEvent is a type alias for the new slot event.
	NewSlotEvent = async.Event[*SlotData]

	// FinalValidatorUpdatesProcessedEvent is a type alias for the final
	// validator updates processed event.
	FinalValidatorUpdatesProcessedEvent = async.Event[transition.ValidatorUpdates]
)

// Messages.
type (
	// GenesisMessage is a type alias for the genesis message.
	GenesisMessage = async.Event[*Genesis]

	// SlotMessage is a type alias for the slot message.
	SlotMessage = async.Event[*SlotData]

	// StatusMessage is a type alias for the status message.
	StatusMessage = async.Event[*service.StatusEvent]
)

/* -------------------------------------------------------------------------- */
/*                                   Dispatcher                               */
/* -------------------------------------------------------------------------- */

type (
	// Dispatcher is a type alias for the dispatcher.
	Dispatcher = dispatcher.Dispatcher
)

/* -------------------------------------------------------------------------- */
/*                                  Pruners                                   */
/* -------------------------------------------------------------------------- */

type (
	// DepositPruner is a type alias for the deposit pruner.
	DepositPruner = pruner.Pruner[*DepositStore]
)
