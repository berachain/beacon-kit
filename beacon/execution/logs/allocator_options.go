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
	"fmt"
	"reflect"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/itsdevbear/bolaris/runtime/service"
)

// WithABI returns an Option for registering
// the contract ABI with the TypeAllocator.
func WithABI(contractAbi *ethabi.ABI) service.Option[TypeAllocator] {
	return func(a *TypeAllocator) error {
		a.abi = contractAbi
		a.sigToName = make(map[common.Hash]string)
		a.sigToType = make(map[common.Hash]reflect.Type)
		return nil
	}
}

// WithNameAndType returns an Option for registering
// an event name and type under the given even signature
// with the TypeAllocator.
func WithNameAndType(
	sig common.Hash,
	name string,
	t reflect.Type,
) service.Option[TypeAllocator] {
	return func(a *TypeAllocator) error {
		event, ok := a.abi.Events[name]
		if !ok {
			return errors.New("event not found in ABI")
		}
		if event.ID != sig {
			return fmt.Errorf("event %s signature does not match, expected %s, got %s", name, event.ID.Hex(), sig.Hex())
		}
		a.sigToName[sig] = name
		a.sigToType[sig] = t
		return nil
	}
}
