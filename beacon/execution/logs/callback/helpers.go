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
	"reflect" //#nosec:G702 // reflect is required for ABI parsing.

	"github.com/ethereum/go-ethereum/accounts/abi"
)

// This function matches each Go implementation to the ABI's respective function.
// It searches for the ABI function in the go struct and performs basic validation on
// the implemented function.
func buildIDsToMethods(si Handler, contractImpl reflect.Value) (map[logSig]*method, error) {
	contractEventsABI := si.ABIEvents()
	contractImplType := contractImpl.Type()
	idsToMethods := make(map[logSig]*method)
	for m := 0; m < contractImplType.NumMethod(); m++ {
		implMethod := contractImplType.Method(m)
		eventName := findMatchingEvent(implMethod, contractEventsABI)
		if eventName == "" {
			continue // nothing in the abi matches our go method.
		}

		method := newMethod(si, contractEventsABI[eventName], implMethod)
		idsToMethods[logSig(contractEventsABI[eventName].ID)] = method
	}

	// verify that every abi method has a corresponding precompile implementation
	for _, abiMethod := range contractEventsABI {
		if _, found := idsToMethods[logSig(abiMethod.ID)]; !found {
			return nil, errors.New("error finding method for event")
		}
	}

	return idsToMethods, nil
}

// findMatchingEvent finds the ABI method that matches the given impl method. It returns the
// key in the ABI methods map that matches the impl method.
func findMatchingEvent(
	implMethod reflect.Method, eventABI map[string]abi.Event,
) string {
	// try to match the impl method name to a ABI event name
	for _, abiMethod := range eventABI {
		if implMethod.Name == abiMethod.RawName {
			if tryMatchInputs(implMethod, abiMethod) {
				return abiMethod.Name // we have a match
			}
		}
	}
	return ""
}

// tryMatchInputs returns true iff the argument types match between the Go implementation and the
// ABI method.
func tryMatchInputs(implMethod reflect.Method, abiMethod abi.Event) bool {
	abiMethodNumIn := len(abiMethod.Inputs)
	implMethodNumIn := implMethod.Type.NumIn()

	// First two args of handler implementation are the receiver contract and the Context, so
	// verify that the ABI method has exactly 2 fewer inputs than the implementation method.
	if implMethodNumIn-2 != abiMethodNumIn {
		return false
	}

	// If the function does not take any inputs, no need to check.
	if abiMethodNumIn == 0 {
		return true
	}

	// Validate that the input args types match ABI input arg types, excluding the
	// first two args (receiver contract and Context). // We also need to handle topics
	// separately, since they should all be common.Hash
	var i = 2
	for ; i < implMethodNumIn; i++ {
		if !abiMethod.Inputs[i-2].Indexed {
			break
		}
	}

	for ; i < implMethodNumIn; i++ {
		implMethodParamType := implMethod.Type.In(i)
		abiMethodParamType := abiMethod.Inputs[i-2].Type.GetType()
		if validateArg(
			reflect.New(implMethodParamType).Elem(), reflect.New(abiMethodParamType).Elem(),
		) != nil {
			return false
		}
	}

	return true
}
