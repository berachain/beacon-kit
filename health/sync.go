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

package health

import (
	"context"
	"time"
)

// WaitForHealthyOf waits for the specified services to be healthy.
func (s *Service) WaitForHealthyOf(
	ctx context.Context,
	callerID string,
	services ...string,
) error {
	healthyCh := make(chan struct{})
	timeout := time.After(
		3 * time.Second, //nolint:gomnd // TODO: configurable.
	)

	// Start a goroutine to wait for the service(s) to be healthy.
	go func() {
		s.svcRegistry.WaitForHealthy(ctx, services...)
		close(healthyCh)
	}()

	// Wait for either the services to become healthy,
	// athe context to be done, or the timeout.
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-healthyCh:
		// Services became healthy before the timeout.
		return nil
	case <-timeout:
		s.Logger().Warn(
			"detected potential service degradation...",
			"caller", callerID, "services", services)
		return ErrHealthCheckTimeout
	}
}
