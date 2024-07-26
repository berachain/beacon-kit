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

package types

import (
	"context"
	"os"
	"path/filepath"
	"runtime"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/json"

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
