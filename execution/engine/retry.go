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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package engine

import (
	"context"
	"math/rand/v2"
	"time"

	"github.com/berachain/beacon-kit/errors"
)

const (
	individualRetryTimeout = time.Minute * 1

	// These values are based on the default values specified at
	// https://github.com/grpc/grpc/blob/master/doc/connection-backoff.md.
	baseDelay  = 100 * time.Millisecond
	multiplier = 1.6
	jitter     = 0.2
	maxDelay   = 5 * time.Second
)

// Backoff returns the amount of time to wait before the next retry given the
// number of retries.
// Implementation copied from google.golang.org/grpc@v1.48.0/internal/backoff/backoff.go.
func Backoff(retries int) time.Duration {
	if retries == 0 {
		return baseDelay
	}

	floatBaseDelay := float64(baseDelay)
	floatMaxDelay := float64(maxDelay)

	for floatBaseDelay < floatMaxDelay && retries > 0 {
		floatBaseDelay *= multiplier
		retries--
	}
	if floatBaseDelay > floatMaxDelay {
		floatBaseDelay = floatMaxDelay
	}
	// Randomize floatBaseDelay delays so that if a cluster of requests start at
	// the same time, they won't operate in lockstep.
	//#nosec:G404 // Doesn't have to be secure rng.
	floatBaseDelay *= 1 + jitter*(rand.Float64()*2-1)
	if floatBaseDelay < 0 {
		return 0
	}

	return time.Duration(floatBaseDelay)
}

// retryWithTimeout retries the given function for the duration of the
// supplied retryTimeout until it returns true or an error. In order for the
// function to be retried, it must return false and no error.
func retryWithTimeout(ctx context.Context, retryTimeout time.Duration, fn func(ctx context.Context) (bool, error)) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, retryTimeout)
	defer cancel()

	var retries int
	for {
		select {
		case <-timeoutCtx.Done():
			return errors.New("retryWithTimeout timed out")
		default:
			// Additional timeout for individual calls. We don't want them lasting forever.
			innerCtx, innerCancel := context.WithTimeout(ctx, individualRetryTimeout)
			ok, err := fn(innerCtx)
			innerCancel()
			if ctx.Err() != nil {
				return errors.Wrap(ctx.Err(), "retry canceled")
			}
			if err != nil {
				return err
			}
			if !ok {
				// Wait for backoff and retry loop.
				select {
				case <-ctx.Done():
				case <-time.After(Backoff(retries)):
				}
				retries++
				continue
			}

			// Call has succeeded.
			return nil
		}
	}
}
