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

package types

import (
	"fmt"
	"strconv"

	"github.com/berachain/beacon-kit/cli/utils/parser"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/primitives/encoding/hex"
	"github.com/berachain/beacon-kit/primitives/math"
)

func BeaconBlockHeaderFromConsensus(h *ctypes.BeaconBlockHeader) *BeaconBlockHeader {
	return &BeaconBlockHeader{
		Slot:          fmt.Sprintf("%d", h.Slot),
		ProposerIndex: fmt.Sprintf("%d", h.ProposerIndex),
		ParentRoot:    h.ParentBlockRoot.Hex(),
		StateRoot:     h.StateRoot.Hex(),
		BodyRoot:      h.BodyRoot.Hex(),
	}
}

func SignedBeaconBlockHeaderFromConsensus(h *ctypes.SignedBeaconBlockHeader) *SignedBeaconBlockHeader {
	return &SignedBeaconBlockHeader{
		Message:   BeaconBlockHeaderFromConsensus(h.Header),
		Signature: h.Signature.String(),
	}
}

func SidecarFromConsensus(sc *datypes.BlobSidecar) *Sidecar {
	proofs := make([]string, len(sc.InclusionProof))
	for i := range sc.InclusionProof {
		proofs[i] = hex.EncodeBytes(sc.InclusionProof[i][:])
	}
	return &Sidecar{
		Index:                       strconv.FormatUint(sc.Index, 10),
		Blob:                        hex.EncodeBytes(sc.Blob[:]),
		KZGCommitment:               hex.EncodeBytes(sc.KzgCommitment[:]),
		SignedBlockHeader:           SignedBeaconBlockHeaderFromConsensus(sc.SignedBeaconBlockHeader),
		KZGProof:                    hex.EncodeBytes(sc.KzgProof[:]),
		KZGCommitmentInclusionProof: proofs,
	}
}

func ValidatorFromConsensus(v *ctypes.Validator) *Validator {
	return &Validator{
		PublicKey:                  v.GetPubkey().String(),
		WithdrawalCredentials:      v.GetWithdrawalCredentials().String(),
		EffectiveBalance:           v.GetEffectiveBalance().Base10(),
		Slashed:                    v.IsSlashed(),
		ActivationEligibilityEpoch: v.GetActivationEligibilityEpoch().Base10(),
		ActivationEpoch:            v.GetActivationEpoch().Base10(),
		ExitEpoch:                  v.GetExitEpoch().Base10(),
		WithdrawableEpoch:          v.GetWithdrawableEpoch().Base10(),
	}
}

// useful in UTs
func ValidatorToConsensus(v *Validator) (*ctypes.Validator, error) {
	pk, err := parser.ConvertPubkey(v.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed parsing public key: %w", err)
	}
	wc, err := parser.ConvertWithdrawalCredentials(v.WithdrawalCredentials)
	if err != nil {
		return nil, fmt.Errorf("failed parsing withdrawals: %w", err)
	}
	eb, err := strconv.ParseUint(v.EffectiveBalance, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed parsing effective balance: %w", err)
	}
	aee, err := strconv.ParseUint(v.ActivationEligibilityEpoch, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed parsing activation eligibility epoch: %w", err)
	}
	ae, err := strconv.ParseUint(v.ActivationEpoch, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed parsing activation epoch: %w", err)
	}
	ee, err := strconv.ParseUint(v.ExitEpoch, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed parsing exit epoch: %w", err)
	}
	we, err := strconv.ParseUint(v.WithdrawableEpoch, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed parsing withdrawable epoch: %w", err)
	}
	return &ctypes.Validator{
		Pubkey:                     pk,
		WithdrawalCredentials:      wc,
		EffectiveBalance:           math.Gwei(eb),
		Slashed:                    v.Slashed,
		ActivationEligibilityEpoch: math.Epoch(aee),
		ActivationEpoch:            math.Epoch(ae),
		ExitEpoch:                  math.Epoch(ee),
		WithdrawableEpoch:          math.Epoch(we),
	}, nil
}
