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

//go:build e2e
// +build e2e

package e2e_test

import (
	"context"
	"testing"
	"time"

	"cosmossdk.io/log"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/starlark_run_config"
	"github.com/kurtosis-tech/kurtosis/api/golang/engine/lib/kurtosis_context"
	"github.com/stretchr/testify/suite"
)

// BeaconE2ESuite is a suite of tests simluating a fully function beacon-kit
// network.
type BeaconKitE2ESuite struct {
	suite.Suite

	ctx     context.Context
	kCtx    *kurtosis_context.KurtosisContext
	logger  log.Logger
	enclave *enclaves.EnclaveContext
}

// TestBeaconKitE2ESuite runs the test suite.
func TestBeaconKitE2ESuite(t *testing.T) {
	suite.Run(t, new(BeaconKitE2ESuite))
}

// this function executes before the test suite begins execution
func (s *BeaconKitE2ESuite) SetupSuite() {
	s.ctx = context.Background()
	s.logger = log.NewTestLogger(s.T())
	var err error
	s.kCtx, err = kurtosis_context.NewKurtosisContextFromLocalEngine()
	s.Require().NoError(err)
	s.logger.Info("Destroying any existing enclave...")
	_ = s.kCtx.DestroyEnclave(s.ctx, "e2e-test-enclave")

	s.logger.Info("Creating enclave...")
	s.enclave, err = s.kCtx.CreateEnclave(s.ctx, "e2e-test-enclave")
	s.Require().NoError(err)
	s.Require().NoError(err)
}

// this function executes after all tests executed
func (s *BeaconKitE2ESuite) TearDownSuite() {
	s.logger.Info("Destroying enclave...")
	s.Require().NoError(s.kCtx.DestroyEnclave(s.ctx, "e2e-test-enclave"))
}

// TestBasicStartup tests the basic startup of the beacon-kit network.
func (s *BeaconKitE2ESuite) TestBasicStartup() {
	s.logger.Info("Running Starlark package...")
	_, cancel, err := s.enclave.RunStarlarkPackage(
		s.ctx,
		"../kurtosis",
		starlark_run_config.NewRunStarlarkConfig(),
	)
	defer cancel()
	s.Require().NoError(err, "Error running Starlark package")

	s.logger.Info("Waiting for services to start...")
	services, err := s.enclave.GetServices()
	s.Require().NoError(err, "Error getting services")

	for k, v := range services {
		s.logger.Info("Service started", "service", k, "uuid", v)
	}

	time.Sleep(3 * time.Second)
}
