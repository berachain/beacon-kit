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
	"fmt"
	"reflect"

	"github.com/berachain/beacon-kit/primitives"
	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
)

// WithABI returns an Option for registering
// the contract ABI with the TypeAllocator.
func WithABI(contractAbi *ethabi.ABI) Option[TypeAllocator] {
	return func(a *TypeAllocator) error {
		a.abi = contractAbi
		a.sigToName = make(map[primitives.ExecutionHash]string)
		a.sigToType = make(map[primitives.ExecutionHash]reflect.Type)
		return nil
	}
}

// WithNameAndType returns an Option for registering
// an event name and type under the given even signature
// with the TypeAllocator.
// NOTE: WithABI must be called before this function.
func WithNameAndType(
	sig primitives.ExecutionHash,
	name string,
	t reflect.Type,
) Option[TypeAllocator] {
	return func(a *TypeAllocator) error {
		event, ok := a.abi.Events[name]
		if !ok {
			return fmt.Errorf("event %s not found in ABI", name)
		}
		if event.ID != sig {
			return fmt.Errorf(
				"event %s signature does not match, expected %s, got %s",
				name, event.ID.Hex(), sig.Hex(),
			)
		}
		a.sigToName[sig] = name
		a.sigToType[sig] = t
		return nil
	}
}
