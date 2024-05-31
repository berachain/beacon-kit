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

package validator

import (
	"context"
	"time"

	engineerrors "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/errors"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
)

// computeStateRoot computes the state root of an outgoing block.
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT, BlobSidecarsT,
]) computeStateRoot(
	ctx context.Context,
	st BeaconStateT,
	blk BeaconBlockT,
) (primitives.Root, error) {
	startTime := time.Now()
	defer s.metrics.measureStateRootComputationTime(startTime)
	if _, err := s.stateProcessor.Transition(
		// TODO: We should think about how having optimistic
		// engine enabled here would affect the proposer when
		// the payload in their block has come from a remote builder.
		&transition.Context{
			Context:                 ctx,
			OptimisticEngine:        true,
			SkipPayloadVerification: true,
			SkipValidateResult:      true,
			SkipValidateRandao:      true,
		},
		st, blk,
	); err != nil {
		return primitives.Root{}, err
	}

	return st.HashTreeRoot()
}

// verifyStateRoot verifies the state root of an incoming block.
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT, BlobSidecarsT,
]) verifyStateRoot(
	ctx context.Context,
	st BeaconStateT,
	blk BeaconBlockT,
) error {
	startTime := time.Now()
	defer s.metrics.measureStateRootVerificationTime(startTime)
	if _, err := s.stateProcessor.Transition(
		// We run with a non-optimistic engine here to ensure
		// that the proposer does not try to push through a bad block.
		&transition.Context{
			Context:                 ctx,
			OptimisticEngine:        false,
			SkipPayloadVerification: false,
			SkipValidateResult:      false,
			SkipValidateRandao:      false,
		},
		st, blk,
	); errors.Is(err, engineerrors.ErrAcceptedPayloadStatus) {
		// It is safe for the validator to ignore this error since
		// the state transition will enforce that the block is part
		// of the canonical chain.
		//
		// TODO: this is only true because we are assuming SSF.
		return nil
	} else if err != nil {
		return err
	}

	return nil
}
