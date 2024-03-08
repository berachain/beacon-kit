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

package suite

import (
	"context"
	"strings"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/kurtosis"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/services"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/starlark_run_config"
	"github.com/kurtosis-tech/kurtosis/api/golang/engine/lib/kurtosis_context"
)

// SetupSuite executes before the test suite begins execution.
func (s *KurtosisE2ESuite) SetupSuite() {
	s.SetupSuiteWithOptions()
}

// Option is a function that sets a field on the KurtosisE2ESuite.
func (s *KurtosisE2ESuite) SetupSuiteWithOptions(opts ...Option) {
	// Setup some sane defaults.
	s.cfg = kurtosis.DefaultE2ETestConfig()
	s.ctx = context.Background()
	s.logger = log.NewTestLogger(s.T())

	// Apply all the provided options, this allows
	// the test suite to be configured in a flexible manner by
	// the caller (i.e overriding defaults).
	for _, opt := range opts {
		if err := opt(s); err != nil {
			s.Require().NoError(err, "Error applying option")
		}
	}

	var err error
	s.kCtx, err = kurtosis_context.NewKurtosisContextFromLocalEngine()
	s.Require().NoError(err)
	s.logger.Info("destroying any existing enclave...")
	//#nosec:G703 // its okay if this errors out. It will error out
	// if there is no enclave to destroy.
	_ = s.kCtx.DestroyEnclave(s.ctx, "e2e-test-enclave")

	s.logger.Info("creating enclave...")
	s.enclave, err = s.kCtx.CreateEnclave(s.ctx, "e2e-test-enclave")
	s.Require().NoError(err)

	s.logger.Info(
		"spinning up enclave...",
		"num_participants",
		len(s.cfg.Participants),
	)
	result, err := s.enclave.RunStarlarkPackageBlocking(
		s.ctx,
		"../kurtosis",
		starlark_run_config.NewRunStarlarkConfig(
			starlark_run_config.WithSerializedParams(
				string(s.cfg.MustMarshalJSON()),
			),
		),
	)
	s.Require().NoError(err, "Error running Starlark package")
	s.Require().Nil(result.ExecutionError, "Error running Starlark package")
	s.Require().Empty(result.ValidationErrors)
	s.SetupExecutionClients()
}

// SetupExecutionClients sets up the execution clients for the test suite.
func (s *KurtosisE2ESuite) SetupExecutionClients() {
	s.executionClients = make(map[string]*ExecutionClient)
	svrcs, err := s.Enclave().GetServices()
	s.Require().NoError(err, "Error getting services")
	for name, v := range svrcs {
		var serviceCtx *services.ServiceContext
		serviceCtx, err = s.Enclave().GetServiceContext(string(v))
		s.Require().NoError(err, "Error getting service context")
		if strings.HasPrefix(string(name), "el-") {
			if s.executionClients[string(name)], err = NewExecutionClientFromServiceCtx(
				serviceCtx,
				s.logger,
			); err != nil {
				// TODO: Figoure out how to handle clients that purposefully
				// don't expose JSON-RPC.
				s.Require().NoError(err, "Error creating execution client")
			}
		}
	}
}

// TearDownSuite cleans up resources after all tests have been executed.
// this function executes after all tests executed.
func (s *KurtosisE2ESuite) TearDownSuite() {
	s.Logger().Info("destroying enclave...")
	s.Require().NoError(s.kCtx.DestroyEnclave(s.ctx, "e2e-test-enclave"))
}
