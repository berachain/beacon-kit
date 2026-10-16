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

//nolint:testpackage // exercises unexported helpers.
package engine

import (
	"testing"
	"time"

	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/stretchr/testify/require"
)

func TestPhaseMaxElapsedTime(t *testing.T) {
	t.Parallel()

	const (
		buildBudget    = 2 * time.Second
		validateBudget = 3 * time.Second
	)
	ee := &Engine{buildBudget: buildBudget, validateBudget: validateBudget}

	cases := []struct {
		phase engineprimitives.EnginePhase
		want  time.Duration
	}{
		{engineprimitives.PhaseBuild, buildBudget},
		{engineprimitives.PhaseValidate, validateBudget},
		// PhaseFinalize and PhaseStartup must be unbounded; never replace 0
		// here without re-reading the rationale in the package doc on
		// EnginePhase.
		{engineprimitives.PhaseFinalize, 0},
		{engineprimitives.PhaseStartup, 0},
	}
	for _, tc := range cases {
		t.Run(tc.phase.String(), func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.want, ee.phaseMaxElapsedTime(tc.phase))
		})
	}
}

func TestPhaseMaxElapsedTime_UnknownDefaultsBounded(t *testing.T) {
	t.Parallel()
	// A wiring bug must not silently produce infinite retry. Unknown phases
	// fall back to the validate budget.
	const validateBudget = 1 * time.Second
	ee := &Engine{validateBudget: validateBudget}
	require.Equal(t, validateBudget, ee.phaseMaxElapsedTime(engineprimitives.EnginePhase(99)))
}
