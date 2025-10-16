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

package suite

import (
	"context"
	"fmt"

	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/starlark_run_config"
)

// StopService stops a running service in the Kurtosis enclave.
func (s *KurtosisE2ESuite) StopService(ctx context.Context, serviceName string) error {
	s.logger.Info("Stopping service", "service", serviceName)

	script := fmt.Sprintf(`
def run(plan):
	plan.stop_service("%s")
`, serviceName)

	result, err := s.enclave.RunStarlarkScriptBlocking(ctx, script, starlark_run_config.NewRunStarlarkConfig())
	if err != nil {
		return fmt.Errorf("failed to stop service %s: %w", serviceName, err)
	}

	if result.ExecutionError != nil {
		return fmt.Errorf("error stopping service %s: %s", serviceName, result.ExecutionError.String())
	}

	if len(result.ValidationErrors) > 0 {
		return fmt.Errorf("validation error stopping service %s: %s", serviceName, result.ValidationErrors[0].String())
	}

	s.logger.Info("Service stopped successfully", "service", serviceName)
	return nil
}

// StartService starts a stopped service in the Kurtosis enclave.
func (s *KurtosisE2ESuite) StartService(ctx context.Context, serviceName string) error {
	s.logger.Info("Starting service", "service", serviceName)

	script := fmt.Sprintf(`
def run(plan):
	plan.start_service("%s")
`, serviceName)

	result, err := s.enclave.RunStarlarkScriptBlocking(ctx, script, starlark_run_config.NewRunStarlarkConfig())
	if err != nil {
		return fmt.Errorf("failed to start service %s: %w", serviceName, err)
	}

	if result.ExecutionError != nil {
		return fmt.Errorf("error starting service %s: %s", serviceName, result.ExecutionError.String())
	}

	if len(result.ValidationErrors) > 0 {
		return fmt.Errorf("validation error starting service %s: %s", serviceName, result.ValidationErrors[0].String())
	}

	s.logger.Info("Service started successfully", "service", serviceName)
	return nil
}

// RestartService stops and then starts a service.
func (s *KurtosisE2ESuite) RestartService(ctx context.Context, serviceName string) error {
	if err := s.StopService(ctx, serviceName); err != nil {
		return err
	}

	return s.StartService(ctx, serviceName)
}
