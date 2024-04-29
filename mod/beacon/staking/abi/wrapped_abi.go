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

package abi

import (
	"errors"

	engineprimitives "github.com/berachain/beacon-kit/mod/primitives-engine"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

// WrappedABI extends the abi.ABI struct to provide
// additional functionality for working with wrapped ABIs.
type WrappedABI struct {
	*abi.ABI
}

// NewWrappedABI creates a new WrappedABI instance from a given abi.ABI
// reference.
func NewWrappedABI(a *abi.ABI) *WrappedABI {
	return &WrappedABI{a}
}

// UnpackLogs attempts to unpack log entries for a specific event.
func (wabi WrappedABI) UnpackLogs(
	out interface{},
	event string,
	log engineprimitives.Log,
) error {
	// Anonymous events are not supported.
	if len(log.Topics) == 0 {
		return errors.New("abi: cannot unpack anonymous event")
	}
	// Check for event signature mismatch.
	if log.Topics[0] != wabi.Events[event].ID {
		return errors.New("abi: event signature mismatch")
	}
	// Unpack the data if present.
	if len(log.Data) > 0 {
		if err := wabi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return err
		}
	}
	// Collect indexed arguments for the event.
	var indexed abi.Arguments
	for _, arg := range wabi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	// Parse the topics to extract indexed event parameters.
	return abi.ParseTopics(out, indexed, log.Topics[1:])
}
