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
	"fmt"
	"sync"
	"time"

	"github.com/berachain/beacon-kit/mod/async/pkg/types"
)

// Route represents a communication route to a single recipient.
// Invariant: there is exactly no more than one route for each messageID.
type Route[ReqT any, RespT any] struct {
	messageID  types.MessageID
	recipient  chan ReqT
	responseCh chan RespT
	timeout    time.Duration
	mu         sync.Mutex
}

// NewRoute creates a new route.
func NewRoute[ReqT any, RespT any](
	messageID types.MessageID,
) *Route[ReqT, RespT] {
	return &Route[ReqT, RespT]{
		messageID:  messageID,
		responseCh: make(chan RespT),
		timeout:    defaultRouterTimeout,
		mu:         sync.Mutex{},
	}
}

// MessageID returns the message ID that the route is responsible for.
func (r *Route[ReqT, RespT]) MessageID() types.MessageID {
	return r.messageID
}

// RegisterReceiver sets the recipient for the route.
func (r *Route[ReqT, RespT]) RegisterReceiver(ch any) error {
	if r.recipient != nil {
		return ErrRouteAlreadySet
	} else if ch == nil {
		return ErrRegisteringNilChannel(r.messageID)
	}
	typedCh, err := ensureType[chan ReqT](ch)
	if err != nil {
		return err
	}
	r.recipient = typedCh
	return nil
}

// SendRequest sends a request to the recipient.
func (r *Route[ReqT, RespT]) SendRequest(req types.MessageI) error {
	typedMsg, err := ensureType[ReqT](req)
	if err != nil {
		return err
	}
	select {
	case r.recipient <- typedMsg:
		return nil
	default:
		return ErrReceiverNotListening(r.messageID)
	}
}

// SendResponse sends a response to the response channel.
func (r *Route[ReqT, RespT]) SendResponse(resp types.MessageI) error {
	fmt.Println("SENDING RESPONSE", resp)
	typedMsg, err := ensureType[RespT](resp)
	if err != nil {
		return err
	}
	fmt.Println("GOOD TYPE", typedMsg)
	r.responseCh <- typedMsg
	return nil
}

// AwaitResponse listens for a response and returns it if it is available
// before the timeout. Otherwise, it returns ErrTimeout.
func (r *Route[ReqT, RespT]) AwaitResponse(respPtr any) error {
	fmt.Println("awaiting RESPONSESESESE")
	select {
	case resp := <-r.responseCh:
		fmt.Println("GOT RESPONSESESESE")
		return assign(resp, respPtr)
	case <-time.After(r.timeout):
		fmt.Println("TIMED OUT GREENENENENENENENENNENENENENENENENENEN")
		return ErrTimeout(r.messageID, r.timeout)
	}
}
