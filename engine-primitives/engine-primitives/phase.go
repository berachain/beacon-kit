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

package engineprimitives

// EnginePhase tags an engine API call with the consensus phase that issued it,
// so the retry wrapper can pick a budget that fits the phase's safety/liveness
// requirements:
//
//   - PhaseBuild and PhaseValidate are bounded: a stuck call must return so
//     CometBFT can move to a new round/proposer. This is what closes the
//     malicious-payload retry-loop class of bugs.
//   - PhaseFinalize is unbounded: the block is already agreed by >=2/3 of
//     validators, so the node must eventually apply it or fall out of
//     consensus. A brief EL outage is absorbed by retrying; the loop logs
//     loudly so operator alerting can catch a persistently stuck node.
//   - PhaseStartup has an unbounded time budget so a slow cold-starting EL
//     (re-importing chain state, snap-syncing) is acceptable. Fatal errors
//     still propagate immediately — a misconfigured boot (wrong JWT, wrong
//     chain ID) fails fast rather than hot-looping before serving consensus.
type EnginePhase int

const (
	// PhaseBuild is used from PrepareProposal when this node is the proposer
	// and is asking the EL to start building a payload.
	PhaseBuild EnginePhase = iota
	// PhaseValidate is used from ProcessProposal when validating a block
	// proposed by another validator.
	PhaseValidate
	// PhaseFinalize is used from FinalizeBlock and the post-block FCU when
	// applying a block that consensus has already agreed on.
	PhaseFinalize
	// PhaseStartup is used from one-shot startup paths (forceSyncUponProcess
	// / forceSyncUponFinalize).
	PhaseStartup
)

func (p EnginePhase) String() string {
	switch p {
	case PhaseBuild:
		return "build"
	case PhaseValidate:
		return "validate"
	case PhaseFinalize:
		return "finalize"
	case PhaseStartup:
		return "startup"
	default:
		return "unknown"
	}
}
