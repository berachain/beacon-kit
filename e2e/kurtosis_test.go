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
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/starlark_run_config"
	"github.com/kurtosis-tech/kurtosis/api/golang/engine/lib/kurtosis_context"
)

func TestE2EKurtosisMultiNode(t *testing.T) {
	ctx := context.Background()
	logger := log.NewTestLogger(t)
	kCtx, err := kurtosis_context.NewKurtosisContextFromLocalEngine()
	if err != nil {
		t.Fatalf("Error instantiating Kurtosis context: %v", err)
	}

	logger.Info("Destroying any existing enclave...")
	_ = kCtx.DestroyEnclave(ctx, "e2e-test-enclave")

	logger.Info("Creating enclave...")
	enclave, err := kCtx.CreateEnclave(ctx, "e2e-test-enclave")
	if err != nil {
		t.Fatalf("Error creating enclave: %v", err)
	}

	defer func() {
		logger.Info("Destroying enclave...")
		if err := kCtx.DestroyEnclave(ctx, "e2e-test-enclave"); err != nil {
			t.Fatalf("Error stopping enclave: %v", err)
		}
	}()

	logger.Info("Running Starlark package...")
	_, cancel, err := enclave.RunStarlarkPackage(
		ctx,
		"../kurtosis",
		starlark_run_config.NewRunStarlarkConfig(),
	)
	defer cancel()

	if err != nil {
		t.Fatalf("Error running Starlark package: %v", err)
	}

	logger.Info("Waiting for services to start...")
	services, err := enclave.GetServices()
	if err != nil {
		t.Fatalf("Error getting services: %v", err)
	}

	for k, v := range services {
		logger.Info("Service started", "service", k, "uuid", v)
	}

	time.Sleep(20 * time.Second)
}
