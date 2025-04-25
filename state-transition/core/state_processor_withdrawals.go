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
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package core

import (
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/state-transition/core/state"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/params"
)

// processWithdrawals as per the Ethereum 2.0 specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/capella/beacon-chain.md#new-process_withdrawals
//
// NOTE: Modified from the Ethereum 2.0 specification to support EVM inflation:
// 1. The first withdrawal MUST be a fixed EVM inflation withdrawal
// 2. Subsequent withdrawals (if any) are processed as validator withdrawals
// 3. This modification reduces the maximum validator withdrawals per block by one.
//
//nolint:gocognit,funlen // TODO: Consider refactoring
func (sp *StateProcessor) processWithdrawals(
	st *state.StateDB, blk *ctypes.BeaconBlock,
) error {
	// Dequeue and verify the logs.
	var (
		body               = blk.GetBody()
		payload            = body.GetExecutionPayload()
		payloadWithdrawals = payload.GetWithdrawals()
	)

	// Get the expected withdrawals.
	expectedWithdrawals, processedPartialWithdrawalsCount, err := st.ExpectedWithdrawals(blk.GetTimestamp())
	if err != nil {
		return err
	}

	// Common validations
	if len(expectedWithdrawals) != len(payloadWithdrawals) {
		return errors.Wrapf(
			ErrNumWithdrawalsMismatch,
			"withdrawals do not match expected length %d, got %d",
			len(expectedWithdrawals), len(payloadWithdrawals),
		)
	}

	// Enforce that first withdrawal is EVM inflation
	if len(payloadWithdrawals) == 0 {
		return ErrZeroWithdrawals
	}
	if !payloadWithdrawals[0].Equals(st.EVMInflationWithdrawal(blk.GetTimestamp())) {
		return ErrFirstWithdrawalNotEVMInflation
	}
	numWithdrawals := len(expectedWithdrawals)

	// Process all subsequent validator withdrawals.
	for i := 1; i < numWithdrawals; i++ {
		// Ensure the withdrawals match the local state.
		if !expectedWithdrawals[i].Equals(payloadWithdrawals[i]) {
			return errors.Wrapf(
				ErrWithdrawalMismatch,
				"withdrawal at index %d does not match expected %s, got %s",
				i,
				spew.Sdump(expectedWithdrawals[i]),
				spew.Sdump(payloadWithdrawals[i]),
			)
		}

		if err = st.DecreaseBalance(
			expectedWithdrawals[i].GetValidatorIndex(), expectedWithdrawals[i].GetAmount(),
		); err != nil {
			return err
		}
	}

	// Update pending partial withdrawals [Introduced in Electra:EIP7251]
	// This case can only be hit after electra.
	if processedPartialWithdrawalsCount > 0 {
		ppWithdrawals, getErr := st.GetPendingPartialWithdrawals()
		if getErr != nil {
			return getErr
		}
		updatedWithdrawals := ppWithdrawals[processedPartialWithdrawalsCount:]
		if setErr := st.SetPendingPartialWithdrawals(updatedWithdrawals); setErr != nil {
			return setErr
		}
		sp.logger.Info(
			"pending partial withdrawals found",
			"original_count", processedPartialWithdrawalsCount,
			"updated_count", len(updatedWithdrawals),
		)
	}

	if numWithdrawals > 1 {
		if err = st.SetNextWithdrawalIndex(
			(expectedWithdrawals[numWithdrawals-1].GetIndex() + 1).Unwrap(),
		); err != nil {
			return err
		}
	}

	totalValidators, err := st.GetTotalValidators()
	if err != nil {
		return err
	}

	// Update the next validator index to start the next withdrawal sweep.
	var nextValidatorIndex math.ValidatorIndex
	if uint64(numWithdrawals) == sp.cs.MaxWithdrawalsPerPayload() {
		// Next sweep starts after the latest withdrawal's validator index.
		nextValidatorIndex = (expectedWithdrawals[numWithdrawals-1].GetValidatorIndex() + 1) % totalValidators
	} else {
		// Advance sweep by the max length of the sweep if there was not a full set of withdrawals.
		nextValidatorIndex, err = st.GetNextWithdrawalValidatorIndex()
		if err != nil {
			return err
		}
		nextValidatorIndex += sp.cs.MaxValidatorsPerWithdrawalsSweep()
		nextValidatorIndex %= totalValidators
	}

	if err = st.SetNextWithdrawalValidatorIndex(nextValidatorIndex); err != nil {
		return err
	}

	sp.logger.Info(
		"Processed withdrawals",
		"num_withdrawals", numWithdrawals,
		"evm_inflation", float64(payloadWithdrawals[0].GetAmount().Unwrap())/params.GWei,
	)

	return nil
}

// processWithdrawalRequest is the equivalent of process_withdrawal_request as defined in the spec.
// It should only be called after the electra hard fork.
// For invalid withdrawal requests, we return nil, and only return error for system errors.
func (sp *StateProcessor) processWithdrawalRequest(
	st *state.StateDB, withdrawalRequest *ctypes.WithdrawalRequest,
) error {
	// If the amount is 0, it's a full exit.
	isFullExitRequest := withdrawalRequest.Amount == constants.FullExitRequestAmount
	pendingPartialWithdrawals, err := st.GetPendingPartialWithdrawals()
	if err != nil {
		return err
	}

	// If partial withdrawal queue is full, only full exits are processed
	if len(pendingPartialWithdrawals) == constants.PendingPartialWithdrawalsLimit && !isFullExitRequest {
		sp.logger.Warn(
			"skipping processing of withdrawal request as partial withdrawal queue is full",
			withdrawalFields(withdrawalRequest, nil)...,
		)
		return nil
	}

	index, validator, err := validateWithdrawal(st, withdrawalRequest)
	if err != nil {
		sp.logger.Info("Failed to validate withdrawal", withdrawalFields(withdrawalRequest, err)...)
		// Note that we do not return error on invalid requests as it's a user error and invalid withdrawal requests are simply skipped.
		return nil
	}
	if validator == nil {
		// This should never occur as we check in validateWithdrawal for nil, but this is needed to appease nilaway.
		sp.logger.Warn("processWithdrawalRequest: nil validator", withdrawalFields(withdrawalRequest, nil)...)
		return errors.New("processWithdrawalRequest: unexpected nil validator")
	}

	if err = verifyWithdrawalConditions(st, validator); err != nil {
		// Note that we do not return error on invalid requests as it's a user error and invalid withdrawal requests are simply skipped.
		sp.logger.Info("Failed to verify withdrawal conditions", withdrawalFields(withdrawalRequest, err)...)
		return nil
	}

	// Process full exit or partial withdrawal.
	if isFullExitRequest {
		sp.logger.Info("Processing full exit request", withdrawalFields(withdrawalRequest, nil)...)
		return sp.processFullExit(st, index, pendingPartialWithdrawals)
	}
	sp.logger.Info("Processing partial withdrawal request", withdrawalFields(withdrawalRequest, nil)...)
	return sp.processPartialWithdrawal(st, withdrawalRequest, validator, index, pendingPartialWithdrawals)
}

// processFullExit processes the full exit request is not a pending partial withdrawal
// and has passed validation of `processWithdrawalRequest`.
func (sp *StateProcessor) processFullExit(
	st *state.StateDB,
	index math.ValidatorIndex,
	pendingPartialWithdrawals ctypes.PendingPartialWithdrawals,
) error {
	pendingBalance := pendingPartialWithdrawals.PendingBalanceToWithdraw(index)
	if pendingBalance == 0 {
		// Only exit validator if it has no pending withdrawals in the queue
		return sp.InitiateValidatorExit(st, index)
	}
	sp.logger.Info("validator has pending balance and cannot full exit",
		"validator_index", index,
		"pending_balance", pendingBalance,
	)
	return nil
}

// processPartialWithdrawal handles the partial withdrawal processing and called after
// request has passed validation of `processWithdrawalRequest`
func (sp *StateProcessor) processPartialWithdrawal(
	st *state.StateDB,
	req *ctypes.WithdrawalRequest,
	validator *ctypes.Validator,
	index math.ValidatorIndex,
	pendingWithdrawals ctypes.PendingPartialWithdrawals,
) error {
	minActivationBalance := sp.cs.MinActivationBalance()
	hasSufficient := validator.GetEffectiveBalance() >= minActivationBalance

	balance, err := st.GetBalance(index)
	if err != nil {
		return err
	}
	pendingBalanceToWithdraw := pendingWithdrawals.PendingBalanceToWithdraw(index)
	hasExcess := balance > minActivationBalance+pendingBalanceToWithdraw

	if validator.HasCompoundingWithdrawalCredential() && hasSufficient && hasExcess {
		toWithdraw := min(balance-minActivationBalance-pendingBalanceToWithdraw, req.Amount)
		// As long as `processPartialWithdrawal` is called after `processSlots`, this will always
		// return the correct slot.
		currentEpoch, getErr := st.GetEpoch()
		if getErr != nil {
			return getErr
		}
		// Note that we do not need to set the ExitEpoch anywhere here as Partial withdrawals do not
		// exit the validator as we enforce that the `toWithdraw` amount is always above the
		// `minActivationBalance`. For example, if a validator's balance is already at
		// `minActivationBalance` (250K BERA on mainnet) and they request a partial withdrawal, the
		// `toWithdraw` amount will always be zero.
		nextEpoch := currentEpoch + 1
		withdrawableEpoch := nextEpoch + sp.cs.MinValidatorWithdrawabilityDelay()
		ppWithdrawal := &ctypes.PendingPartialWithdrawal{
			ValidatorIndex:    index,
			Amount:            toWithdraw,
			WithdrawableEpoch: withdrawableEpoch,
		}
		pendingWithdrawals = append(pendingWithdrawals, ppWithdrawal)
		return st.SetPendingPartialWithdrawals(pendingWithdrawals)
	}
	sp.logger.Info("validator cannot withdraw partial balance",
		"validator_index", index,
		"validator_pubkey", validator.GetPubkey().String(),
		"balance", balance,
		"pending_balance_to_withdraw", pendingBalanceToWithdraw,
		"has_sufficient", hasSufficient,
		"has_excess", hasExcess,
	)
	return nil
}

// validateWithdrawal checks that the validator exists and that the withdrawal credentials match.
func validateWithdrawal(
	st *state.StateDB, withdrawalRequest *ctypes.WithdrawalRequest,
) (math.ValidatorIndex, *ctypes.Validator, error) {
	// Verify pubkey exists
	index, err := st.ValidatorIndexByPubkey(withdrawalRequest.ValidatorPubKey)
	if err != nil {
		return 0, nil, err
	}

	validator, err := st.ValidatorByIndex(index)
	if err != nil {
		return 0, nil, err
	}
	if validator == nil {
		// This should never occur as we expect ErrNotFound if the validator is not found.
		return 0, nil, errors.New("validateWithdrawal: validator does not exist")
	}

	// Verify withdrawal credentials
	if !validator.HasExecutionWithdrawalCredential() {
		return 0, nil, errors.New("validator does not have execution withdrawal credentials")
	}
	correctCred, err := validator.WithdrawalCredentials.ToExecutionAddress()
	if err != nil {
		return 0, nil, err
	}
	if !withdrawalRequest.SourceAddress.Equals(correctCred) {
		return 0, nil, errors.New("source address does not match execution withdrawal credential")
	}
	return index, validator, nil
}

// verifyWithdrawalConditions checks additional conditions like active status, exit not initiated,
// and minimal activation period.
func verifyWithdrawalConditions(st *state.StateDB, validator *ctypes.Validator) error {
	currentEpoch, err := st.GetEpoch()
	if err != nil {
		return err
	}
	// Verify the validator is active
	if !validator.IsActive(currentEpoch) {
		return errors.New("validator is not active")
	}
	// Verify exit has not been initiated
	if validator.GetExitEpoch() != constants.FarFutureEpoch {
		return errors.New("withdrawal already initiated")
	}
	// Verify the validator has been active long enough
	// In the spec, config.SHARD_COMMITTEE_PERIOD is added as well, but we ignore this since
	// it's related to ETH data shards which is not relevant for us.
	if currentEpoch < validator.ActivationEpoch {
		return errors.New("validator not active long enough")
	}
	return nil
}

// withdrawalFields returns the structured fields for logging any WithdrawalRequest.
// error is optional
func withdrawalFields(req *ctypes.WithdrawalRequest, err error) []interface{} {
	logFields := []interface{}{
		"source_address", req.SourceAddress.String(),
		"validator_pubkey", req.ValidatorPubKey.String(),
		"amount", req.Amount,
	}
	if err != nil {
		logFields = append(logFields, "error", err)
	}
	return logFields
}
