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
	"time"

	"github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/errors"
)

// errTimeout is the error returned when a dispatch operation timed out.
//
//nolint:gochecknoglobals // errors
var (
	// errTimeout is the error returned when a dispatch operation timed out.
	errTimeout = func(messageID types.MessageID, timeout time.Duration) error {
		return errors.Newf("message %s reached the max timeout of %s",
			messageID, timeout)
	}

	// errRouteAlreadySet is the error returned when the route is already set.
	errRouteAlreadySet = errors.New("route already set")

	// errRegisteringNilChannel is the error returned when the channel to
	// register is nil.
	errRegisteringNilChannel = func(messageID types.MessageID) error {
		return errors.Newf("cannot register nil channel for route: %s",
			messageID)
	}

	// errReceiverNotReady is the error returned when the receiver channel is
	// full, closed, or not listening.
	errReceiverNotReady = func(messageID types.MessageID) error {
		return errors.Newf(
			"receiver channel is full, closed, or not listening. Route: %s",
			messageID,
		)
	}

	errSendingNilResponse = func(messageID types.MessageID) error {
		return errors.Newf("cannot send nil response for route: %s",
			messageID)
	}

	// errIncompatibleAssignee is the error returned when the assignee is not
	// compatible with the assigner.
	errIncompatibleAssignee = func(
		assigner interface{}, assignee interface{},
	) error {
		return errors.Newf(
			"incompatible assignee, expected: %T, received: %T",
			assigner,
			assignee,
		)
	}
)
