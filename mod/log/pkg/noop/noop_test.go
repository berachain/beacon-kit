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

package noop_test

import (
	"testing"

	"github.com/berachain/beacon-kit/mod/log/pkg/noop"
)

func TestNewLogger(t *testing.T) {
	logger := noop.NewLogger()
	if logger == nil {
		t.Error("Expected NewLogger to return a non-nil logger")
	}
}

func TestLogger_Info(*testing.T) {
	logger := noop.NewLogger()
	logger.Info("test message", "key1", "value1", "key2", "value2")
	// Since it's a no-op logger, there's nothing to assert. This test just
	// ensures no panic occurs.
}

func TestLogger_Warn(*testing.T) {
	logger := noop.NewLogger()
	logger.Warn("test warning", "key1", "value1", "key2", "value2")
	// Since it's a no-op logger, there's nothing to assert. This test just
	// ensures no panic occurs.
}

func TestLogger_Error(*testing.T) {
	logger := noop.NewLogger()
	logger.Error("test error", "key1", "value1", "key2", "value2")
	// Since it's a no-op logger, there's nothing to assert. This test just
	// ensures no panic occurs.
}

func TestLogger_Debug(*testing.T) {
	logger := noop.NewLogger()
	logger.Debug("test debug", "key1", "value1", "key2", "value2")
	// Since it's a no-op logger, there's nothing to assert. This test just
	// ensures no panic occurs.
}
