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
	"sync"
	"time"

	"github.com/berachain/beacon-kit/mod/async/pkg/types"
)

// Route represents a communication route to a single recipient.
// Invariant: there is exactly no more than one route for each messageID.
type Route[ReqT any, RespT any] struct {
	messageID   types.MessageID
	recipientCh chan types.Message[ReqT]
	responseCh  chan types.Message[RespT]
	timeout     time.Duration
	mu          sync.Mutex
}

// NewRoute creates a new route.
func NewRoute[ReqT any, RespT any](
	messageID types.MessageID,
) *Route[ReqT, RespT] {
	return &Route[ReqT, RespT]{
		messageID:  messageID,
		responseCh: make(chan types.Message[RespT]),
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
	if r.recipientCh != nil {
		return errRouteAlreadySet
	} else if ch == nil {
		return errRegisteringNilChannel(r.messageID)
	}
	typedCh, err := ensureType[chan types.Message[ReqT]](ch)
	if err != nil {
		return err
	}
	r.recipientCh = typedCh
	return nil
}

// SendRequestAsync accepts a future and sends a request to the recipient
// channel. Once the response is available, it will be written to the future.
func (r *Route[ReqT, RespT]) SendRequest(
	req types.BaseMessage, future any,
) error {
	r.sendRequest(req)
	typedFuture, err := ensureType[types.FutureI[RespT]](future)
	if err != nil {
		return err
	}
	go r.populateFuture(typedFuture)
	return nil
}

// SendResponse sends a response to the response channel.
func (r *Route[ReqT, RespT]) SendResponse(resp types.BaseMessage) error {
	typedMsg, err := ensureType[types.Message[RespT]](resp)
	if err != nil {
		return err
	}
	r.responseCh <- typedMsg
	return nil
}

// populateFuture sends a done signal to the future when the response is
// available.
func (r *Route[ReqT, RespT]) populateFuture(future types.FutureI[RespT]) {
	select {
	case resp := <-r.responseCh:
		future.SetResult(resp.Data(), resp.Error())
	case <-time.After(r.timeout):
		errTimeout(r.messageID, r.timeout)
	}
}

// SendRequest sends a request to the recipient.
func (r *Route[ReqT, RespT]) sendRequest(req types.BaseMessage) error {
	typedReq, err := ensureType[types.Message[ReqT]](req)
	if err != nil {
		return err
	}

	select {
	case r.recipientCh <- typedReq:
		return nil
	default:
		// Channel is full or closed
		return errReceiverNotReady(r.messageID)
	}
}
