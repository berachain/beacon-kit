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

package server

import "context"

// publisher is the interface that supports basic event feed operations.
type publisher interface {
	// Start starts the event feed.
	Start(ctx context.Context)
	// Publish publishes the given event to the event feed.
	Publish(ctx context.Context, event any)
	// Subscribe subscribes the given channel to the event feed.
	Subscribe(ch any) error
	// Unsubscribe unsubscribes the given channel from the event feed.
	Unsubscribe(ch any) error
}

// messageRoute is the interface that supports basic message route operations.
type messageRoute interface {
	// RegisterRecipient sets the recipient for the route.
	RegisterRecipient(ch chan any) error
	// SendRequest sends a request to the recipient.
	SendRequest(msg any) error
	// SendResponse sends a response to the recipient.
	SendResponse(msg any) error
	// AwaitResponse awaits a response from the route.
	AwaitResponse(emptyResp any) error
}
