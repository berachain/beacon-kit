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

// TypeRegistry is a struct that stores registered types
// corresponding to events in each contract.
// The information is used to allocate empty objects
// of the registed types, into which the log data
// can be unmarshaled.
type TypeRegistry struct {
	abi       *ethabi.ABI
	sigToName map[common.Hash]string
	sigToType map[common.Hash]reflect.Type
	// We can add a callback here to handle the logs
	// based on their names like the current callback logic.
}

// NewTypeRegistry returns a new TypeRegistry
// for multiple events in a contract.
func NewTypeRegistry(abi *ethabi.ABI) *TypeRegistry {
	return &TypeRegistry{
		abi:       abi,
		sigToName: make(map[common.Hash]string),
		sigToType: make(map[common.Hash]reflect.Type),
	}
}

// RegisterEvent registers type of an event with the registry.
func (r *TypeRegistry) RegisterEvent(
	eventName string,
	eventType reflect.Type,
) error {
	event, ok := r.abi.Events[eventName]
	if !ok {
		return errors.New("event not found in ABI")
	}
	r.sigToName[event.ID] = eventName
	r.sigToType[event.ID] = eventType
	return nil
}

// GetABI returns the ABI of the contract.
func (r *TypeRegistry) GetABI() *ethabi.ABI {
	return r.abi
}

// GetName returns the name of the event
// corresponding to the signature.
func (r *TypeRegistry) GetName(sig common.Hash) (string, error) {
	name, ok := r.sigToName[sig]
	if !ok {
		return "", errors.New("event not found in registry")
	}
	return name, nil
}

// Allocate returns an empty object of the type,
// which is registered with the signature.
func (r *TypeRegistry) Allocate(sig common.Hash) (reflect.Value, error) {
	eventType, ok := r.sigToType[sig]
	if !ok {
		return reflect.Value{}, errors.New("event type not found in registry")
	}
	return reflect.New(eventType), nil
}
