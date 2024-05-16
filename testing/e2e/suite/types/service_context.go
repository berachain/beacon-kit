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

package types

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"

	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/services"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/starlark_run_config"
)

type WrappedServiceContext struct {
	runStarklarkScript func(
		ctx context.Context,
		serializedScript string,
		runConfig *starlark_run_config.StarlarkRunConfig,
	) (*enclaves.StarlarkRunResult, error)
	*services.ServiceContext
	helpersScript string
}

//nolint:dogsled // no risk from e2e suite
func getHelpersScript() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	path := filepath.Join(dir, "../../../../kurtosis/src/lib/helpers.star")

	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func NewWrappedServiceContext(
	serviceCtx *services.ServiceContext,
	runStarklarkScript func(
		ctx context.Context,
		serializedScript string,
		runConfig *starlark_run_config.StarlarkRunConfig,
	) (*enclaves.StarlarkRunResult, error),
) *WrappedServiceContext {
	return &WrappedServiceContext{
		runStarklarkScript: runStarklarkScript,
		ServiceContext:     serviceCtx,
		helpersScript:      getHelpersScript(),
	}
}

func (s *WrappedServiceContext) Start(
	ctx context.Context,
	enclaveContext *enclaves.EnclaveContext,
) (*enclaves.StarlarkRunResult, error) {
	res, err := s.RunHelper(ctx, "start_service", map[string]interface{}{
		"service_name": s.GetServiceName(),
	})
	if err != nil {
		return nil, err
	}

	replacementSCtx, err := enclaveContext.GetServiceContext(
		string(s.ServiceContext.GetServiceName()),
	)
	if err != nil {
		return nil, err
	}

	s.ServiceContext = replacementSCtx

	return res, nil
}

func (s *WrappedServiceContext) Stop(
	ctx context.Context,
) (*enclaves.StarlarkRunResult, error) {
	return s.RunHelper(ctx, "stop_service", map[string]interface{}{
		"service_name": s.GetServiceName(),
	})
}

func (s *WrappedServiceContext) RunHelper(
	ctx context.Context,
	mainFunctionName string,
	args map[string]interface{},
) (*enclaves.StarlarkRunResult, error) {
	jsonBytes, err := json.Marshal(args)
	if err != nil {
		panic(err)
	}

	res, err := s.runStarklarkScript(
		ctx,
		s.helpersScript,
		starlark_run_config.NewRunStarlarkConfig(
			starlark_run_config.WithMainFunctionName(mainFunctionName),
			starlark_run_config.WithSerializedParams(string(jsonBytes)),
		),
	)
	return res, err
}
