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

package telemetry

import (
	"context"

	"github.com/cosmos/cosmos-sdk/telemetry"
)

type Config = telemetry.Config

// Service is a telemetry service.
type Service struct {
	m *telemetry.Metrics
}

// NewService creates a new telemetry service.
func NewService(cfg *telemetry.Config) (*Service, error) {
	m, err := telemetry.New(*cfg)
	if err != nil {
		return nil, err
	}
	return &Service{
		m: m,
	}, nil
}

// Name returns the service name.
func (s *Service) Name() string {
	return "telemetry"
}

// Start starts the telemetry service.
func (s *Service) Start(context.Context) error {
	return nil
}

func (s *Service) Stop() error {
	return nil
}
