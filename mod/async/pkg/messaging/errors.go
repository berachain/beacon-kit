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

package messaging

import (
	"reflect"
	"time"

	"github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/errors"
)

// ErrTimeout is the error returned when a dispatch operation timed out.
var (
	ErrTimeout = func(messageID types.MessageID, timeout time.Duration) error {
		return errors.Newf("message %s timed out after %s", messageID, timeout)
	}

	ErrRouteAlreadySet = errors.New("route already set")

	ErrRegisteringNilChannel = func(messageID types.MessageID) error {
		return errors.Newf("cannot register nil channel for route: %s", messageID)
	}

	ErrReceiverNotListening = func(messageID types.MessageID) error {
		return errors.Newf("receiver may be registered but not listening for route: %s", messageID)
	}

	errIncompatibleAssigneeType = func(assigner interface{}, assignee interface{}) error {
		assignerType := reflect.TypeOf(assigner)
		assigneeType := reflect.TypeOf(assignee)
		return errors.Newf(
			"incompatible assignee, expected: %T, received: %T",
			assignerType,
			assigneeType,
		)
	}

	errIncompatibleAssignee = func(assigner interface{}, assignee interface{}) error {
		return errors.Newf(
			"incompatible assignee, expected: %T, received: %T",
			assigner,
			assignee,
		)
	}
)
