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

package execution

import (
	"context"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/itsdevbear/bolaris/beacon/core/state"
	enginetypes "github.com/itsdevbear/bolaris/engine/types"
	"github.com/itsdevbear/bolaris/primitives"
)

// BeaconStorageBackend is an interface that wraps the basic BeaconState method.
type BeaconStorageBackend interface {
	// BeaconState provides access to the underlying beacon state.
	BeaconState(ctx context.Context) state.BeaconState
}

// LogFactory is an interface that can unmarshal Ethereum logs into objects,
// in the form of reflect.Value, with appropriate types for each type of logs.
type LogFactory interface {
	GetRegisteredAddresses() []primitives.ExecutionAddress
	ProcessLogs(
		logs []ethtypes.Log,
		blkNum uint64,
	) ([]*reflect.Value, error)
}

// FCUConfig is a struct that holds the configuration for a fork choice update.
type FCUConfig struct {
	// HeadEth1Hash is the hash of the head eth1 block we are updating the
	// execution client's head to be.
	HeadEth1Hash common.Hash

	// ProposingSlot is the slot that the execution client should propose a
	// block
	// for if Attributes neither nil nor empty.
	ProposingSlot primitives.Slot

	// Attributes is a list of payload attributes to include in a forkchoice
	// update to the execution client. It is used to signal to the execution
	// client that
	// it should build a payload.
	Attributes enginetypes.PayloadAttributer
}
