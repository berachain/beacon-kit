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
	"bytes"
	"errors"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

// Factory is a struct that can be used to unmarshal
// Ethereum logs into objects with the appropriate types.
type Factory struct {
	// addressToAbi maps contract addresses to Allocators for their events.
	addressToAllocator map[common.Address]*TypeAllocator
}

// NewFactory creates a new Factory with the given options.
func NewFactory(opts ...Option[Factory]) (*Factory, error) {
	f := &Factory{
		addressToAllocator: make(map[common.Address]*TypeAllocator),
	}
	for _, opt := range opts {
		if err := opt(f); err != nil {
			return nil, err
		}
	}
	return f, nil
}

// UnmarshalEthLog unmarshals an Ethereum log into an object
// with the appropriate type, based on the registered event
// corresponding to the log.
func (f *Factory) UnmarshalEthLog(log *ethtypes.Log) (reflect.Value, error) {
	// Get the event signature from the log.
	sig := log.Topics[0]

	// Get the ABI, type, and name of the event from the factory.
	// This function only processes logs from contracts and events
	// that have been registered with the factory.
	allocator, ok := f.addressToAllocator[log.Address]
	if !ok {
		return reflect.Value{}, errors.New(
			"allocator not found for contract address",
		)
	}

	contractAbi := allocator.GetABI()

	eventName, err := allocator.GetName(sig)
	if err != nil {
		return reflect.Value{}, err
	}

	// Allocate an empty object, which we can unmarshal the log data into.
	into, err := allocator.Allocate(sig)
	if err != nil {
		return reflect.Value{}, err
	}

	// Unpack the log data into the new instance.
	intoPtr := into.Interface()
	if err = contractAbi.UnpackIntoInterface(
		intoPtr, eventName, log.Data,
	); err != nil {
		return reflect.Value{}, err
	}

	return reflect.ValueOf(intoPtr), nil
}

// GetRegisteredAddresses returns the addresses of the contracts
// that have been registered with the factory.
func (f *Factory) GetRegisteredAddresses() []common.Address {
	addresses := make([]common.Address, 0, len(f.addressToAllocator))
	for addr := range f.addressToAllocator {
		addresses = append(addresses, addr)
	}
	return addresses
}

// IsRegisteredLog returns true if its corresponding event
// has been registered with the factory.
func (f *Factory) IsRegisteredLog(log *ethtypes.Log) bool {
	allocator, ok := f.addressToAllocator[log.Address]
	if !ok {
		return false
	}

	sig := log.Topics[0]
	contractAbi := allocator.GetABI()
	eventName, err := allocator.GetName(sig)
	if err != nil {
		return false
	}
	return bytes.Equal(sig.Bytes(), contractAbi.Events[eventName].ID.Bytes())
}
