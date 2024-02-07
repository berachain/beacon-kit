// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2023, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
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

package callback

import (
	"context"
	"reflect" //#nosec:G702 // reflect is required for ABI parsing.

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	evmv1 "github.com/itsdevbear/bolaris/types/evm/v1"
)

// logSig is a fixed length byte array that represents the method ID of a precompile method.
type logSig common.Hash

// method attaches a specific abi.Event to a golang struct to process it.
type method struct {
	// rcvr is the receiver of the method's executable.
	rcvr Handler

	// abiEvent is the abi.Event corresponding to this methods executable.
	abiEvent abi.Event

	// execute is the executable which will execute when `abiEvent` is seen.
	execute reflect.Method
}

// newMethod creates and returns a new `method` with the given abiEvent, and executable.
func newMethod(
	rcvr Handler, abiEvent abi.Event, execute reflect.Method,
) *method {
	return &method{
		rcvr:     rcvr,
		abiEvent: abiEvent,
		execute:  execute,
	}
}

// Call executes the corresponding method attached to a specific log.
func (m *method) Call(ctx context.Context, log *evmv1.Log) error {
	// Unpack the args from the input, if any exist.
	topics := log.GetTopics()
	unpackedArgs, err := m.abiEvent.Inputs.Unpack(log.GetData())
	if err != nil {
		return err
	}

	// Convert topics to common.Hash
	reflectedUnpackedArgs := make([]reflect.Value, 0, (len(topics)-1)+len(unpackedArgs))
	for _, topic := range topics[1:] {
		reflectedUnpackedArgs = append(reflectedUnpackedArgs, reflect.ValueOf(topic))
	}

	// Convert the unpacked args to reflect values.
	for _, unpacked := range unpackedArgs {
		reflectedUnpackedArgs = append(reflectedUnpackedArgs, reflect.ValueOf(unpacked))
	}

	// Call the executable the reflected values.
	results := m.execute.Func.Call(
		append(
			[]reflect.Value{
				reflect.ValueOf(m.rcvr),
				reflect.ValueOf(ctx),
			},
			reflectedUnpackedArgs...,
		),
	)

	// Check if the results contain an error.
	if len(results) > 0 {
		var ok bool
		if err, ok = results[0].Interface().(error); ok {
			return err
		}
	}
	return nil
}
