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

// reportingInterval is the interval at which the health service
// logs the health status of services.
const reportingInterval = 5 * time.Second

// reportingLoop initiates a loop that periodically checks and
// reports the health status of services.
func (s *Service) reportingLoop(ctx context.Context) {
	ticker := time.NewTicker(reportingInterval)
	for {
		select {
		case <-ticker.C:
			s.reportStatuses()
		case <-ctx.Done():
			return
		}
	}
}

// reportStatuses logs the health status of all services.
func (s *Service) reportStatuses() {
	svcStatuses := s.retrieveStatuses()
	for _, svc := range svcStatuses {
		if svc.Healthy {
			s.Logger().
				Info("service is reporting healthy 🌤️ ", "service", svc.Name)
		} else {
			s.Logger().Error("service is reporting unhealthy ⛈️ ",
				"service", svc.Name, "error", svc.Err)
		}
	}
}
