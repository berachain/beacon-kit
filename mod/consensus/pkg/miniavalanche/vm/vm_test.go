package vm

import (
	"context"
	"fmt"
	"testing"

	"github.com/ava-labs/avalanchego/database/memdb"
	"github.com/ava-labs/avalanchego/snow/engine/common"
	"github.com/ava-labs/avalanchego/snow/snowtest"
	"github.com/ava-labs/avalanchego/snow/validators"
	"github.com/ava-labs/avalanchego/utils/logging"
	"github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/consensus/pkg/miniavalanche/middleware"
	"github.com/berachain/beacon-kit/mod/log/pkg/noop"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/async"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"github.com/stretchr/testify/require"
)

var _ types.EventDispatcher = (*genesisDispatcherStub)(nil)

// genesisDispatcherStub allows loading expections around VM-middleware interactions
// upon VM.Initialize
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
		ch, found := d.brokers[async.GenesisDataProcessed]
		if !found {
			return fmt.Errorf("VM has not subscribed for async.GenesisDataProcessed event")
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

func (d *genesisDispatcherStub) Subscribe(eventID async.EventID, ch any) error {
	d.brokers[eventID] = ch
	return nil
}
func (d *genesisDispatcherStub) Unsubscribe(eventID async.EventID, ch any) error { return nil }

func TestVMInitialization(t *testing.T) {
	r := require.New(t)

	// setup VM
	var (
		beaconLogger = noop.NewLogger[any]()
		avaLogger    = logging.NoLog{} // TODO: consolidate logs
		dp           = &genesisDispatcherStub{
			brokers: make(map[async.EventID]any),
		}
		mdw = middleware.NewABCIMiddleware(dp, beaconLogger)
		f   = Factory{
			Config: Config{
				Validators: validators.NewManager(),
			},
			Middleware: mdw,
		}

		ctx      = context.TODO()
		msgChan  = make(chan common.Message, 1)
		chainCtx = snowtest.Context(t, snowtest.PChainID)
		db       = memdb.New()
	)

	vmIntf, err := f.New(avaLogger)
	r.NoError(err)
	r.IsType(vmIntf, &VM{})
	vm := vmIntf.(*VM)

	genesisBytes, err := setupTestGenesis()
	r.NoError(err)

	// Start middleware before initializing VM
	dp.Subscribe(async.GenesisDataReceived, nil) // allows VM to publish genesis data
	r.NoError(mdw.Start(ctx))

	// test initialization
	err = vm.Initialize(ctx, chainCtx, db, genesisBytes, nil, nil, msgChan, nil, nil)
	r.NoError(err)
}

func setupTestGenesis() ([]byte, error) {
	genesisData := &Base64Genesis{
		Validators: []Base64GenesisValidator{
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
	genContent, err := BuildBase64GenesisString(genesisData)
	if err != nil {
		return nil, err
	}

	// unmarshal genesis
	return ParseBase64StringToBytes(genContent)
}
