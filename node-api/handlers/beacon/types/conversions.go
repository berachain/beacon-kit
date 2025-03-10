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
	"github.com/berachain/beacon-kit/consensus-types/deneb"
	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/primitives/encoding/hex"
	"github.com/berachain/beacon-kit/primitives/math"
)

func BeaconBlockHeaderFromConsensus(h *deneb.BeaconBlockHeader) *BeaconBlockHeader {
	return &BeaconBlockHeader{
		Slot:          fmt.Sprintf("%d", h.Slot),
		ProposerIndex: fmt.Sprintf("%d", h.ProposerIndex),
		ParentRoot:    h.ParentBlockRoot.Hex(),
		StateRoot:     h.StateRoot.Hex(),
		BodyRoot:      h.BodyRoot.Hex(),
	}
}

func SignedBeaconBlockHeaderFromConsensus(h *deneb.SignedBeaconBlockHeader) *SignedBeaconBlockHeader {
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

func ValidatorFromConsensus(v *deneb.Validator) *Validator {
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
func ValidatorToConsensus(v *Validator) (*deneb.Validator, error) {
	pk, err := parser.ConvertPubkey(v.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed parsing public key: %w", err)
	}
	wc, err := parser.ConvertWithdrawalCredentials(v.WithdrawalCredentials)
	if err != nil {
		return nil, fmt.Errorf("failed parsing withdrawals: %w", err)
	}
	eb, err := math.U64FromString(v.EffectiveBalance)
	if err != nil {
		return nil, fmt.Errorf("failed parsing effective balance: %w", err)
	}
	aee, err := math.U64FromString(v.ActivationEligibilityEpoch)
	if err != nil {
		return nil, fmt.Errorf("failed parsing activation eligibility epoch: %w", err)
	}
	ae, err := math.U64FromString(v.ActivationEpoch)
	if err != nil {
		return nil, fmt.Errorf("failed parsing activation epoch: %w", err)
	}
	ee, err := math.U64FromString(v.ExitEpoch)
	if err != nil {
		return nil, fmt.Errorf("failed parsing exit epoch: %w", err)
	}
	we, err := math.U64FromString(v.WithdrawableEpoch)
	if err != nil {
		return nil, fmt.Errorf("failed parsing withdrawable epoch: %w", err)
	}
	return &deneb.Validator{
		Pubkey:                     pk,
		WithdrawalCredentials:      wc,
		EffectiveBalance:           eb,
		Slashed:                    v.Slashed,
		ActivationEligibilityEpoch: aee,
		ActivationEpoch:            ae,
		ExitEpoch:                  ee,
		WithdrawableEpoch:          we,
	}, nil
}
