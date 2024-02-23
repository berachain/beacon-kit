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

// Registry is a struct that stores registered information for each contract.
type Registry struct {
	abi       *ethabi.ABI
	sigToName map[common.Hash]string
	sigToType map[common.Hash]reflect.Type
	// We can add a callback here to handle the logs
	// based on their names like the current callback logic.
}

// NewRegistry returns a new Registry.
func NewRegistry(abi *ethabi.ABI) *Registry {
	return &Registry{
		abi:       abi,
		sigToName: make(map[common.Hash]string),
		sigToType: make(map[common.Hash]reflect.Type),
	}
}

// RegisterEvent registers information of an event with the registry.
func (r *Registry) RegisterEvent(
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
func (r *Registry) GetABI() *ethabi.ABI {
	return r.abi
}

func (r *Registry) GetTypeAndName(
	sig common.Hash,
) (reflect.Type, string, error) {
	event, ok := r.sigToType[sig]
	if !ok {
		return nil, "", errors.New("event not found in registry")
	}
	name, ok := r.sigToName[sig]
	if !ok {
		return nil, "", errors.New("event not found in registry")
	}
	return event, name, nil
}
