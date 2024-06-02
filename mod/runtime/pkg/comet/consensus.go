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

package comet

import (
	"context"

	math "github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	cmtproto "github.com/cometbft/cometbft/api/cometbft/types/v1"
	cmttypes "github.com/cometbft/cometbft/types"
)

type ChainSpec interface {
	// GetCometBFTConfigForSlot returns the CometBFT configuration for the given
	// slot.
	GetCometBFTConfigForSlot(math.Slot) any
}

// ConsensusParamsStore is a store for consensus parameters.
type ConsensusParamsStore struct {
	cs ChainSpec
}

// NewConsensusParamsStore creates a new ConsensusParamsStore.
func NewConsensusParamsStore(cs ChainSpec) *ConsensusParamsStore {
	return &ConsensusParamsStore{
		cs: cs,
	}
}

// Get retrieves the consensus parameters from the store.
// It returns the consensus parameters and an error, if any.
func (s *ConsensusParamsStore) Get(
	context.Context,
) (cmtproto.ConsensusParams, error) {
	return s.cs.
		GetCometBFTConfigForSlot(0).(*cmttypes.ConsensusParams).
		ToProto(), nil
}

// Has checks if the consensus parameters exist in the store.
// It returns a boolean indicating the presence of the parameters and an error,
// if any.
func (s *ConsensusParamsStore) Has(context.Context) (bool, error) {
	return true, nil
}

// Set stores the given consensus parameters in the store.
// It returns an error, if any.
func (s *ConsensusParamsStore) Set(
	_ context.Context,
	_ cmtproto.ConsensusParams,
) error {
	return nil
}
