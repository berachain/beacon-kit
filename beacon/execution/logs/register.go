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
)

// RegisterABI registers a contract ABI with the factory.
func (f *Factory) RegisterABI(contractAddress common.Address, abi *ethabi.ABI) {
	f.addressToRegistry[contractAddress] = NewRegistry(abi)
}

// RegisterEvent registers information of an event with the factory.
func (f *Factory) RegisterEvent(
	contractAddress common.Address,
	eventName string,
	eventType reflect.Type,
) error {
	registry, ok := f.addressToRegistry[contractAddress]
	if !ok {
		return errors.New("registry not found for contract address")
	}
	return registry.RegisterEvent(eventName, eventType)
}

// RegisterWithCallback allows a service to register a contract ABI
// and event information with the factory.
func (f *Factory) RegisterWithCallback(registrant LogRegistrant) error {
	contractAddresses := registrant.GetAddresses()
	for _, contractAddress := range contractAddresses {
		contractAbi := registrant.GetABI(contractAddress)
		f.RegisterABI(contractAddress, contractAbi)
		for eventName, event := range contractAbi.Events {
			eventID := event.ID
			eventType := registrant.GetType(eventID)
			if eventType == nil {
				continue
			}
			if err := f.RegisterEvent(
				contractAddress, eventName, eventType,
			); err != nil {
				return err
			}
		}
	}
	return nil
}
