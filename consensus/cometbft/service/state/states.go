// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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

package state

import (
	"errors"
	"fmt"
)

var ErrNoFinalStateAtHeight = errors.New("no final state at height")

type BlkHash string

type States struct {
	finalStates map[int64]*State
}

func NewStates() *States {
	return &States{
		finalStates: make(map[int64]*State),
	}
}

func (ss *States) MarkAsFinal(h int64, s *State) error {
	// note: for the time being we allow overwriting a final state.
	// We may add checks against this later on (or allow it only in
	// specific cases).
	ss.finalStates[h] = s
	return nil
}

func (ss *States) GetFinalState(h int64) (*State, error) {
	res, found := ss.finalStates[h]
	if !found {
		return nil, fmt.Errorf("%w: height %d", ErrNoFinalStateAtHeight, h)
	}
	return res, nil
}

func (ss *States) DoneWithHeight(h int64) {
	delete(ss.finalStates, h)
}
