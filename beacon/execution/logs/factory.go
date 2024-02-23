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

package logs

import (
	"errors"
	"reflect"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

// Factory is a struct that can be used to unmarshal
// Ethereum logs into objects with the appropriate types.
type Factory struct {
	// addressToAbi is a map of contract addresses to their ABIs.
	addressToAbi map[common.Address]*ethabi.ABI
	// sigToName is a map of event signatures to their names.
	sigToName map[common.Hash]string
	// sigToType is a map of event signatures to their corresponding types.
	sigToType map[common.Hash]reflect.Type
}

// NewFactory returns a new Factory.
func NewFactory() *Factory {
	return &Factory{
		addressToAbi: make(map[common.Address]*ethabi.ABI),
		sigToName:    make(map[common.Hash]string),
		sigToType:    make(map[common.Hash]reflect.Type),
	}
}

// RegisterEvent registers information of an event with the factory.
func (f *Factory) RegisterEvent(
	contractAddress common.Address,
	contractAbi *ethabi.ABI,
	eventName string,
	eventType reflect.Type,
) {
	eventID := contractAbi.Events[eventName].ID
	f.addressToAbi[contractAddress] = contractAbi
	f.sigToName[eventID] = eventName
	f.sigToType[eventID] = eventType
}

// UnmarshalEthLog unmarshals an Ethereum log into an object
// with the appropriate type, based on the registered event
// corresponding to the log.
func (f *Factory) UnmarshalEthLog(log *ethtypes.Log) (reflect.Value, error) {
	var (
		contractAbi *ethabi.ABI
		eventType   reflect.Type
		eventName   string
		ok          bool
	)

	// Get the event signature from the log.
	sig := log.Topics[0]

	// Get the ABI, type, and name of the event from the factory.
	// This function only processes logs from contracts and events
	// that have been registered with the factory.
	if contractAbi, ok = f.addressToAbi[log.Address]; !ok {
		return reflect.Value{}, errors.New("abi not found for log address")
	}
	if eventType, ok = f.sigToType[sig]; !ok {
		return reflect.Value{}, errors.New("type not found for log signature")
	}
	if eventName, ok = f.sigToName[sig]; !ok {
		return reflect.Value{}, errors.New("name not found for log signature")
	}

	// Create a new instance of the event type.
	into := reflect.New(eventType).Interface()
	// Unpack the log data into the new instance.
	if err := contractAbi.UnpackIntoInterface(into, eventName, log.Data); err != nil {
		return reflect.Value{}, err
	}
	return reflect.ValueOf(into), nil
}
