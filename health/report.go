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
	"reflect"
	"time"
)

func (s *Service) reportingLoop(ctx context.Context) {
	ticker := time.NewTicker(healthCheckInterval)
	for {
		select {
		case <-ticker.C:
			s.reportStatuses()
		case <-ctx.Done():
			return
		}
	}
}

func (s *Service) reportStatuses() {
	// Define a local struct to hold service type and error information.
	type report struct {
		svc reflect.Type // Service type
		err error        // Error, if any
	}

	// Initialize slices to hold healthy and unhealthy service types.
	var (
		healthy   = make([]reflect.Type, 0) // List of healthy services
		unhealthy = make(
			[]report,
			0,
		) // List of unhealthy services with errors
	)

	// Iterate over the statuses returned by the service registry.
	for svcType, err := range s.svcRegistry.Statuses() {
		if err == nil {
			// If there's no error, consider the service healthy.
			healthy = append(healthy, svcType)
		} else {
			// If there's an error, consider the service unhealthy.
			unhealthy = append(unhealthy, report{svcType, err})
		}
	}

	// Log the healthy services.
	for _, svcType := range healthy {
		s.Logger().Info("reporting healthy 🍉", "service", svcType)
	}

	// Log the unhealthy services and their errors.
	for _, r := range unhealthy {
		s.Logger().Error("reporting unhealthy 🍉",
			"service", r.svc, "err", r.err)
	}
}
