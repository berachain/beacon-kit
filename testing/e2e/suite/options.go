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

package suite

import (
	"context"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/testing/e2e/config"
)

// Type Option is a function that sets a field on the KurtosisE2ESuite.
type Option func(*KurtosisE2ESuite) error

// WithConfig sets the E2ETestConfig for the test suite.
func WithConfig(cfg *config.E2ETestConfig) Option {
	return func(s *KurtosisE2ESuite) error {
		s.cfg = cfg
		return nil
	}
}

// WithContext sets the context for the test suite.
func WithContext(ctx context.Context) Option {
	return func(s *KurtosisE2ESuite) error {
		s.ctx = ctx
		return nil
	}
}

// WithLogger sets the logger for the test suite.
func WithLogger(logger log.Logger) Option {
	return func(s *KurtosisE2ESuite) error {
		s.logger = logger
		return nil
	}
}
