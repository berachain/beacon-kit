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

package engine

import (
	"context"
	"fmt"
	"time"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	engineerrors "github.com/berachain/beacon-kit/engine-primitives/errors"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/execution/client"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/cenkalti/backoff/v5"
	cmtcfg "github.com/cometbft/cometbft/config"
)

// consensusPhaseBudgetFraction is the share of the CometBFT phase timeout
// reserved for the engine retry budget. The remainder is headroom for the
// other phase work that shares the same consensus window — getPayload and
// block assembly for PhaseBuild; blob KZG verification, state checks, and
// vote propagation for PhaseValidate.
const consensusPhaseBudgetFraction = 0.75

// Engine is Beacon-Kit's implementation of the `ExecutionEngine`
// from the Ethereum 2.0 Specification.
type Engine struct {
	// ec is the engine client that the engine will use to
	// interact with the execution layer.
	ec *client.EngineClient
	// logger is the logger for the engine.
	logger log.Logger
	// metrics is the metrics for the engine.
	metrics *engineMetrics
	// buildBudget and validateBudget cap the retry loop in PhaseBuild and
	// PhaseValidate respectively. They are sized as a fraction of the
	// corresponding CometBFT phase timeout so a stuck engine call yields in
	// time for consensus to advance, while leaving headroom for the rest of
	// the phase work. PhaseFinalize and PhaseStartup are unbounded and
	// ignore these.
	//
	// validateBudget is pegged to TimeoutPrevote rather than TimeoutPropose:
	// ProcessProposal runs during the propose phase, but the engine call's
	// real deadline is "vote must be cast before TimeoutPrevote elapses."
	// TimeoutPrevote is the safe upper bound, not the tightest one.
	buildBudget    time.Duration
	validateBudget time.Duration
}

// New creates a new Engine. The PhaseBuild and PhaseValidate retry budgets are
// derived from the matching CometBFT consensus phase timeouts (TimeoutPropose
// and TimeoutPrevote) so they track operator-tuned consensus timeouts.
//
// Panics if either timeout is non-positive: zero would yield a zero budget,
// which collapses to unbounded retry and silently undoes the bounded-phase
// contract this constructor enforces.
func New(
	engineClient *client.EngineClient,
	logger log.Logger,
	telemtrySink TelemetrySink,
	consensusCfg *cmtcfg.ConsensusConfig,
) *Engine {
	if consensusCfg.TimeoutPropose <= 0 || consensusCfg.TimeoutPrevote <= 0 {
		panic(fmt.Sprintf(
			"engine.New: ConsensusConfig timeouts must be positive (TimeoutPropose=%v, TimeoutPrevote=%v)",
			consensusCfg.TimeoutPropose, consensusCfg.TimeoutPrevote,
		))
	}
	return &Engine{
		ec:             engineClient,
		logger:         logger,
		metrics:        newEngineMetrics(telemtrySink, logger),
		buildBudget:    time.Duration(float64(consensusCfg.TimeoutPropose) * consensusPhaseBudgetFraction),
		validateBudget: time.Duration(float64(consensusCfg.TimeoutPrevote) * consensusPhaseBudgetFraction),
	}
}

// GetPayload returns the payload and blobs bundle for the given slot.
func (ee *Engine) GetPayload(
	ctx context.Context,
	req *ctypes.GetPayloadRequest,
) (ctypes.BuiltExecutionPayloadEnv, error) {
	return ee.ec.GetPayload(
		ctx, req.PayloadID,
		req.ForkVersion,
	)
}

// NotifyForkchoiceUpdate notifies the execution client of a forkchoice update.
//
// Retry policy is selected by phase:
//   - PhaseBuild and PhaseValidate are bounded so a stuck EL can't trap
//     consensus.
//   - PhaseFinalize and PhaseStartup are unbounded *for transient signals*
//     (transport errors, 5xx, SYNCING) so a brief EL outage doesn't drop a
//     block that consensus has already agreed on.
//
// Fatal errors (HTTP 4xx, pre-defined JSON-RPC errors like parse / invalid
// request) are Permanent in *every* phase including Finalize. They encode
// "this request will never succeed against this EL" — retrying forever turns
// a misconfigured JWT or wrong chain ID into a silent node hang. Brief
// outages produce IsNonFatalError signals, which the case above handles.
func (ee *Engine) NotifyForkchoiceUpdate(
	ctx context.Context,
	req *ctypes.ForkchoiceUpdateRequest,
	phase engineprimitives.EnginePhase,
) (*engineprimitives.PayloadID, error) {
	hasPayloadAttributes := req.PayloadAttributes != nil

	return backoff.Retry(
		ctx,
		func() (*engineprimitives.PayloadID, error) {
			ee.metrics.markNotifyForkchoiceUpdateCalled(hasPayloadAttributes)
			payloadID, err := ee.ec.ForkchoiceUpdated(
				ctx, req.State, req.PayloadAttributes, req.ForkVersion,
			)

			switch {
			case err == nil:
				ee.metrics.markForkchoiceUpdateValid(req.State, hasPayloadAttributes, payloadID)
				if payloadID == nil && hasPayloadAttributes {
					ee.logger.Warn(
						"Received nil payload ID on VALID engine response",
						"head_eth1_hash", req.State.HeadBlockHash,
						"safe_eth1_hash", req.State.SafeBlockHash,
						"finalized_eth1_hash", req.State.FinalizedBlockHash,
					)
					return nil, backoff.Permanent(ErrNilPayloadOnValidResponse)
				}
				return payloadID, nil

			case errors.IsAny(err, engineerrors.ErrSyncingPayloadStatus):
				ee.logger.Info("NotifyForkchoiceUpdate: EL syncing. Retrying...", "phase", phase)
				ee.metrics.markForkchoiceUpdateSyncing(req.State, err)
				return nil, err

			case client.IsNonFatalError(err):
				ee.logger.Info(
					"NotifyForkchoiceUpdate: EL returns non fatal error. Retrying...",
					"phase", phase,
					"err", err,
				)
				ee.metrics.markForkchoiceUpdateNonFatalError(err)
				return nil, err

			case errors.Is(err, engineerrors.ErrInvalidPayloadStatus):
				ee.logger.Error("NotifyForkchoiceUpdate: EL returned invalid payload.", "phase", phase)
				ee.metrics.markForkchoiceUpdateInvalid(req.State, err)
				return nil, backoff.Permanent(err)

			case client.IsFatalError(err):
				ee.logger.Info(
					"NotifyForkchoiceUpdate: EL returns fatal error.",
					"phase", phase,
					"err", err,
				)
				ee.metrics.markForkchoiceUpdateFatalError(err)
				return nil, backoff.Permanent(err)

			default:
				ee.logger.Info(
					"NotifyForkchoiceUpdate: EL returns unknown error.",
					"phase", phase,
					"err", err,
				)
				ee.metrics.markForkchoiceUpdateUndefinedError(err)
				return nil, backoff.Permanent(err)
			}
		},
		backoff.WithBackOff(ee.newBackoff()),
		backoff.WithMaxTries(0),
		backoff.WithMaxElapsedTime(ee.phaseMaxElapsedTime(phase)),
	)
}

// NotifyNewPayload notifies the execution client of the new payload.
func (ee *Engine) NotifyNewPayload(
	ctx context.Context,
	req ctypes.NewPayloadRequest,
	phase engineprimitives.EnginePhase,
) error {
	var (
		payloadHash       = req.GetExecutionPayload().GetBlockHash()
		payloadParentHash = req.GetExecutionPayload().GetParentHash()
	)

	_, err := backoff.Retry(
		ctx,
		func() (*common.ExecutionHash, error) {
			ee.metrics.markNewPayloadCalled(payloadHash, payloadParentHash)
			lastValidHash, err := ee.ec.NewPayload(ctx, req)

			switch {
			case err == nil:
				ee.metrics.markNewPayloadValid(payloadHash, payloadParentHash)
				return lastValidHash, nil

			case errors.IsAny(err, engineerrors.ErrSyncingPayloadStatus, engineerrors.ErrAcceptedPayloadStatus):
				ee.metrics.markNewPayloadAcceptedSyncingPayloadStatus(err, payloadHash, payloadParentHash)
				// Finalize/Startup accept SYNCING because the FCU that follows
				// drives the EL to sync; the verifying phases must retry until
				// the EL is caught up so they can verify the block.
				if phase == engineprimitives.PhaseFinalize || phase == engineprimitives.PhaseStartup {
					ee.logger.Warn(
						"NotifyNewPayload: pushed new payload to SYNCING/ACCEPTED node.",
						"phase", phase,
						"err", err,
						"blockNum", req.GetExecutionPayload().GetNumber(),
						"blockHash", payloadHash,
					)
					return &common.ExecutionHash{}, nil
				}
				ee.logger.Info(
					"NotifyNewPayload: EL returns non valid status. Retrying...",
					"phase", phase,
					"err", err,
				)
				return nil, err

			case client.IsNonFatalError(err):
				ee.logger.Info(
					"NotifyNewPayload: EL returns non fatal error. Retrying...",
					"phase", phase,
					"err", err,
				)
				if lastValidHash == nil {
					lastValidHash = &common.ExecutionHash{}
				}
				ee.metrics.markNewPayloadNonFatalError(payloadHash, *lastValidHash, err)
				return nil, err

			case errors.Is(err, engineerrors.ErrInvalidPayloadStatus):
				ee.logger.Error("NotifyNewPayload: EL returned invalid payload.", "phase", phase)
				ee.metrics.markNewPayloadInvalidPayloadStatus(payloadHash)
				return nil, backoff.Permanent(err)

			case client.IsFatalError(err):
				ee.logger.Error(
					"NotifyNewPayload: EL returns fatal error.",
					"phase", phase,
					"err", err,
				)
				if lastValidHash == nil {
					lastValidHash = &common.ExecutionHash{}
				}
				ee.metrics.markNewPayloadFatalError(payloadHash, *lastValidHash, err)
				return nil, backoff.Permanent(err)

			default:
				ee.logger.Error(
					"NotifyNewPayload: EL returns unknown error.",
					"phase", phase,
					"err", err,
				)
				ee.metrics.markNewPayloadUndefinedError(payloadHash, err)
				return nil, backoff.Permanent(err)
			}
		},
		backoff.WithBackOff(ee.newBackoff()),
		backoff.WithMaxTries(0),
		backoff.WithMaxElapsedTime(ee.phaseMaxElapsedTime(phase)),
	)
	return err
}

// phaseMaxElapsedTime returns the backoff MaxElapsedTime for a given phase.
// 0 means unbounded.
func (ee *Engine) phaseMaxElapsedTime(phase engineprimitives.EnginePhase) time.Duration {
	switch phase {
	case engineprimitives.PhaseBuild:
		return ee.buildBudget
	case engineprimitives.PhaseValidate:
		return ee.validateBudget
	case engineprimitives.PhaseFinalize, engineprimitives.PhaseStartup:
		return 0
	default:
		// Unknown phase: bound conservatively so a wiring bug can't introduce
		// an infinite retry by accident.
		return ee.validateBudget
	}
}

func (ee *Engine) newBackoff() *backoff.ExponentialBackOff {
	// Configure backoff. Between each retry it waits RPCRetryInterval, growing
	// exponentially up to RPCMaxRetryInterval. MaxElapsedTime is set per-call
	// via phaseMaxElapsedTime.
	engineAPIBackoff := backoff.NewExponentialBackOff()
	engineAPIBackoff.InitialInterval = ee.ec.GetRPCRetryInterval()
	engineAPIBackoff.MaxInterval = ee.ec.GetRPCMaxRetryInterval()
	return engineAPIBackoff
}
