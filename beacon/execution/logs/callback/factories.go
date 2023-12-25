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
	"errors"
	"reflect"
)

// ===========================================================================
// Stateful Container Factory
// ===========================================================================

// CallbackFactory is used to build stateful precompile containers.
type CallbackFactory struct{}

// NewCallbackFactory creates and returns a new `CallbackFactory`.
func NewCallbackFactory() *CallbackFactory {
	return &CallbackFactory{}
}

// Build returns a stateful precompile container for the given base contract implementation. This
// function will return an error if the given contract is not a stateful implementation.
//
// Build implements `AbstractFactory`.
func (sf *CallbackFactory) Build(
	si CallbackHandler,
) (*wrappedCallbackHandler, error) {
	// add precompile methods to stateful container, if any exist
	idsToMethods, err := buildIdsToMethods(si, reflect.ValueOf(si))
	if err != nil {
		return nil, err
	}

	return New(si, idsToMethods)
}

// This function matches each Go implementation of the precompile to the ABI's respective function.
// It searches for the ABI function in the Go precompile contract and performs basic validation on
// the implemented function.
func buildIdsToMethods(si CallbackHandler, contractImpl reflect.Value) (map[methodID]*method, error) {
	contractEventsABI := si.ABIEvents()
	contractImplType := contractImpl.Type()
	idsToMethods := make(map[methodID]*method)
	for m := 0; m < contractImplType.NumMethod(); m++ {
		implMethod := contractImplType.Method(m)
		eventName, err := FindMatchingABIMethod(implMethod, contractEventsABI)
		if err != nil {
			return nil, err
		}

		if eventName == "" {
			continue // nothing in the abi matches our go method.
		}

		method := newMethod(si, contractEventsABI[eventName], implMethod)
		idsToMethods[methodID(contractEventsABI[eventName].ID)] = method
	}

	// verify that every abi method has a corresponding precompile implementation
	for _, abiMethod := range contractEventsABI {
		if _, found := idsToMethods[methodID(abiMethod.ID)]; !found {
			return nil, errors.New("error finding method for event")
		}
	}

	return idsToMethods, nil
}
