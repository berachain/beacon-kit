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

package vm_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/ava-labs/avalanchego/snow/engine/common"
	"github.com/ava-labs/avalanchego/snow/snowtest"
	"github.com/ava-labs/avalanchego/snow/validators"
	"github.com/ava-labs/avalanchego/utils/logging"
	"github.com/berachain/beacon-kit/mod/async/pkg/types"
	avalanchewrappers "github.com/berachain/beacon-kit/mod/consensus/pkg/miniavalanche/avalanche-wrappers"
	"github.com/berachain/beacon-kit/mod/consensus/pkg/miniavalanche/middleware"
	"github.com/berachain/beacon-kit/mod/consensus/pkg/miniavalanche/vm"
	"github.com/berachain/beacon-kit/mod/log/pkg/noop"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/async"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	cosmosdb "github.com/cosmos/cosmos-db"
	"github.com/stretchr/testify/require"
)

var _ types.EventDispatcher = (*genesisDispatcherStub)(nil)

type genesisDispatcherStub struct {
	brokers map[async.EventID]any
}

func (d *genesisDispatcherStub) Publish(event async.BaseEvent) error {
	_, found := d.brokers[event.ID()]
	if !found {
		return fmt.Errorf("can't publish event %v", event.ID())
	}

	if event.ID() == async.GenesisDataReceived {
		// as soon as VM published the genesis, send back
		// the async.GenesisDataProcessed event
		var ch any
		ch, found = d.brokers[async.GenesisDataProcessed]
		if !found {
			return errors.New("async.GenesisDataProcessed not subscribed")
		}
		genCh, ok := ch.(chan async.Event[transition.ValidatorUpdates])
		if !ok {
			return fmt.Errorf("unexpected channel type %T", ch)
		}

		go func() {
			noUpdatesEvt := async.NewEvent[transition.ValidatorUpdates](
				context.TODO(),
				async.GenesisDataProcessed,
				transition.ValidatorUpdates{},
			)
			genCh <- noUpdatesEvt
		}()
	}

	return nil
}

func (d *genesisDispatcherStub) Subscribe(
	eventID async.EventID,
	ch any,
) error {
	d.brokers[eventID] = ch
	return nil
}
func (d *genesisDispatcherStub) Unsubscribe(_ async.EventID, _ any) error {
	return nil
}

func TestVMInitialization(t *testing.T) {
	r := require.New(t)

	// setup VM
	var (
		beaconLogger = noop.NewLogger[any]()
		avaLogger    = logging.NoLog{} // TODO: consolidate logs
		dp           = &genesisDispatcherStub{
			brokers: make(map[async.EventID]any),
		}
		mdw      = middleware.NewABCIMiddleware(dp, beaconLogger)
		cosmosDB = cosmosdb.NewPrefixDB(cosmosdb.NewMemDB(), vm.BerachainDBPrefix)
		db       = avalanchewrappers.NewDB(cosmosDB)
		f        = vm.Factory{
			Config: vm.Config{
				Validators: validators.NewManager(),
			},
			BaseDB:     db,
			Middleware: mdw,
			StoreKey:   *components.ProvideKVStoreKey(),
		}

		ctx      = context.TODO()
		msgChan  = make(chan common.Message, 1)
		chainCtx = snowtest.Context(t, snowtest.PChainID)
	)

	vmIntf, err := f.New(avaLogger)
	r.NoError(err)
	r.IsType(&vm.VM{}, vmIntf)
	vm, _ := vmIntf.(*vm.VM)

	genesisBytes, err := setupTestGenesis()
	r.NoError(err)

	// Start middleware before initializing VM
	r.NoError(dp.Subscribe(async.GenesisDataReceived, nil))
	r.NoError(mdw.Start(ctx))

	// test initialization
	err = vm.Initialize(
		ctx,
		chainCtx,
		db,
		genesisBytes,
		nil,
		nil,
		msgChan,
		nil,
		nil,
	)
	r.NoError(err)
}

func setupTestGenesis() ([]byte, error) {
	genesisData := &vm.Base64Genesis{
		Validators: []vm.Base64GenesisValidator{
			{
				NodeID: testGenesisValidators[0].NodeID.String(),
				Weight: testGenesisValidators[0].Weight,
			},
			{
				NodeID: testGenesisValidators[1].NodeID.String(),
				Weight: testGenesisValidators[1].Weight,
			},
		},
		EthGenesis: string(testEthGenesisBytes),
	}

	// marshal genesis
	genContent, err := vm.BuildBase64GenesisString(genesisData)
	if err != nil {
		return nil, err
	}

	// unmarshal genesis
	return vm.ParseBase64StringToBytes(genContent)
}
