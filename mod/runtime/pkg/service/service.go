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

package service

import (
	"context"
)

// A Service is a component containing all necessary and sufficient
// logic to operate over a particular business domain.
//
// The Service struct is 2-layered facade that contains:
//   - runner: Intercepts events and delegates to the processor. If this
//     service doesn't need to handle events, the runner simply acts as a
//     runner.
//   - processor: Executes the service's domain business logic
//
// Note: In theory, all one really needs to start a "service" is the
// runner, but the below definition makes it clear from a semantics
// perspective that a service is composed of both the runner and a
// processor, while abstracting away the implementation details.
type Service[
	runnerT Runner[processorT],
	processorT any,
] struct {
	runner    runnerT
	processor processorT
}

// NewService creates a new service with the given runner and processor.
func NewService[
	runnerT Runner[processorT],
	processorT any,
](
	runner runnerT,
	processor processorT,
) *Service[runnerT, processorT] {
	return &Service[runnerT, processorT]{
		runner:    runner,
		processor: processor,
	}
}

// Start will attach the processor to the runner and start the runner.
func (b *Service[_, _]) Start(ctx context.Context) error {
	b.runner.AttachProcessor(b.processor)
	return b.runner.Start(ctx)
}

// Name returns the name of the service from the underlying runner.
func (b *Service[_, _]) Name() string {
	return b.runner.Name()
}
