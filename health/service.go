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

	"github.com/itsdevbear/bolaris/runtime/service"
)

// Service is a health service.
type Service struct {
	service.BaseService

	// svcRegistry is the service registry.
	svcRegistry *service.Registry
}

// Start starts the health service.
func (s *Service) Start(ctx context.Context) {
	go s.reportingLoop(ctx)
}

// retrieveStatuses returns the health status of all services.
func (s *Service) retrieveStatuses() []*serviceStatus {
	var (
		rawStatuses = s.svcRegistry.Statuses()
		statuses    = make([]*serviceStatus, 0, len(rawStatuses))
	)

	for k, v := range rawStatuses {
		s := &serviceStatus{
			Name:    k.String(),
			Healthy: true,
		}

		if v != nil {
			s.Healthy = false
			s.Err = v.Error()
		}
		statuses = append(statuses, s)
	}
	return statuses
}
