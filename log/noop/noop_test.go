// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
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

package noop_test

import (
	"testing"

	"github.com/berachain/beacon-kit/log/noop"
)

func TestNewLogger(t *testing.T) {
	logger := noop.NewLogger[any]()
	if logger == nil {
		t.Error("Expected NewLogger to return a non-nil logger")
	}
}

func TestLogger_Info(*testing.T) {
	logger := noop.NewLogger[any]()
	logger.Info("test message", "key1", "value1", "key2", "value2")
	// Since it's a no-op logger, there's nothing to assert. This test just
	// ensures no panic occurs.
}

func TestLogger_Warn(*testing.T) {
	logger := noop.NewLogger[any]()
	logger.Warn("test warning", "key1", "value1", "key2", "value2")
	// Since it's a no-op logger, there's nothing to assert. This test just
	// ensures no panic occurs.
}

func TestLogger_Error(*testing.T) {
	logger := noop.NewLogger[any]()
	logger.Error("test error", "key1", "value1", "key2", "value2")
	// Since it's a no-op logger, there's nothing to assert. This test just
	// ensures no panic occurs.
}

func TestLogger_Debug(*testing.T) {
	logger := noop.NewLogger[any]()
	logger.Debug("test debug", "key1", "value1", "key2", "value2")
	// Since it's a no-op logger, there's nothing to assert. This test just
	// ensures no panic occurs.
}
