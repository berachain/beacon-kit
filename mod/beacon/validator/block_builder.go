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

package validator

import (
	"context"
	"fmt"
	"time"

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	"golang.org/x/sync/errgroup"
)

// buildBlockAndSidecars builds a new beacon block.
func (s *Service[
	_, BeaconBlockT, _, _, BlobSidecarsT, _, _, _, _, _, _, _, SlotDataT,
]) buildBlockAndSidecars(
	ctx context.Context,
	slotData SlotDataT,
) (BeaconBlockT, BlobSidecarsT, error) {
	var (
		blk      BeaconBlockT
		sidecars BlobSidecarsT
		start    = time.Now()
		g, _     = errgroup.WithContext(ctx)
	)
	defer s.metrics.measureRequestBlockForProposalTime(start)

	st := s.sb.StateFromContext(ctx)

	if _, err := s.stateProcessor.ProcessSlots(st, slotData.GetSlot()); err != nil {
		return blk, sidecars, err
	}

	reveal, err := s.buildRandaoReveal(st, slotData.GetSlot())
	if err != nil {
		return blk, sidecars, err
	}

	blk, err = s.getEmptyBeaconBlockForSlot(st, slotData.GetSlot())
	if err != nil {
		return blk, sidecars, err
	}

	envelope, err := s.retrieveExecutionPayload(ctx, st, blk)
	if err != nil || envelope == nil {
		return blk, sidecars, fmt.Errorf("error retrieving payload: %w", err)
	}

	if err = s.buildBlockBody(ctx, st, blk, reveal, envelope, slotData); err != nil {
		return blk, sidecars, err
	}

	g.Go(func() error {
		sidecars, err = s.blobFactory.BuildSidecars(blk, envelope.GetBlobsBundle())
		return err
	})

	g.Go(func() error {
		return s.computeAndSetStateRoot(ctx, st, blk)
	})

	if err = g.Wait(); err != nil {
		return blk, sidecars, err
	}

	s.logger.Info(
		"Beacon block built successfully",
		"slot", slotData.GetSlot().Base10(),
		"state_root", blk.GetStateRoot(),
		"duration", time.Since(start).String(),
	)

	return blk, sidecars, nil
}

// getEmptyBeaconBlockForSlot creates a new empty block.
func (s *Service[
	_, BeaconBlockT, _, BeaconStateT, _, _, _, _, _, _, _, _, _,
]) getEmptyBeaconBlockForSlot(
	st BeaconStateT, slot math.Slot,
) (BeaconBlockT, error) {
	parentBlockRoot, err := st.GetBlockRootAtIndex((slot.Unwrap() - 1) % s.chainSpec.SlotsPerHistoricalRoot())
	if err != nil {
		return BeaconBlockT{}, err
	}

	proposerIndex, err := st.ValidatorIndexByPubkey(s.signer.PublicKey())
	if err != nil {
		return BeaconBlockT{}, err
	}

	return BeaconBlockT.NewWithVersion(
		slot,
		proposerIndex,
		parentBlockRoot,
		s.chainSpec.ActiveForkVersionForSlot(slot),
	)
}

// buildRandaoReveal builds a randao reveal for the given slot.
func (s *Service[
	_, _, _, BeaconStateT, _, _, _, _, _, _, ForkDataT, _, _,
]) buildRandaoReveal(st BeaconStateT, slot math.Slot) (crypto.BLSSignature, error) {
	epoch := s.chainSpec.SlotToEpoch(slot)
	genesisValidatorsRoot, err := st.GetGenesisValidatorsRoot()
	if err != nil {
		return crypto.BLSSignature{}, err
	}

	signingRoot := ForkDataT.New(
		version.FromUint32[common.Version](s.chainSpec.ActiveForkVersionForEpoch(epoch)),
		genesisValidatorsRoot,
	).ComputeRandaoSigningRoot(s.chainSpec.DomainTypeRandao(), epoch)

	return s.signer.Sign(signingRoot[:])
}

// retrieveExecutionPayload retrieves the execution payload for the block.
func (s *Service[
	_, BeaconBlockT, _, BeaconStateT, _, _, _, _, ExecutionPayloadT, _, _, _, _,
]) retrieveExecutionPayload(
	ctx context.Context, st BeaconStateT, blk BeaconBlockT,
) (engineprimitives.BuiltExecutionPayloadEnv[ExecutionPayloadT], error) {
	envelope, err := s.localPayloadBuilder.RetrievePayload(ctx, blk.GetSlot(), blk.GetParentBlockRoot())
	if err != nil {
		s.metrics.failedToRetrievePayload(blk.GetSlot(), err)
		lph, err := st.GetLatestExecutionPayloadHeader()
		if err != nil {
			return nil, err
		}
		return s.localPayloadBuilder.RequestPayloadSync(
			ctx, st, blk.GetSlot(),
			uint64(time.Now().Unix()+1), lph.GetBlockHash(),
			lph.GetParentHash(),
		)
	}
	return envelope, nil
}

// buildBlockBody assembles the block body with necessary components.
func (s *Service[
	_, BeaconBlockT, _, BeaconStateT, _, _, _, Eth1DataT, ExecutionPayloadT, _, _, _, SlotDataT,
]) buildBlockBody(
	_ context.Context,
	st BeaconStateT,
	blk BeaconBlockT,
	reveal crypto.BLSSignature,
	envelope engineprimitives.BuiltExecutionPayloadEnv[ExecutionPayloadT],
	slotData SlotDataT,
) error {
	body := blk.GetBody()
	if body.IsNil() {
		return ErrNilBlkBody
	}

	body.SetRandaoReveal(reveal)

	blobsBundle := envelope.GetBlobsBundle()
	if blobsBundle == nil {
		return ErrNilBlobsBundle
	}

	body.SetBlobKzgCommitments(blobsBundle.GetCommitments())

	deposits, err := s.sb.DepositStore().GetDepositsByIndex(
		st.GetEth1DepositIndex(), s.chainSpec.MaxDepositsPerBlock(),
	)
	if err != nil {
		return err
	}

	body.SetDeposits(deposits)
	body.SetEth1Data(Eth1DataT{}.New(common.Root{}, 0, common.ExecutionHash{}))

	graffiti, err := bytes.ToBytes32(bytes.ExtendToSize([]byte(s.cfg.Graffiti), bytes.B32Size))
	if err != nil {
		return fmt.Errorf("failed to process graffiti: %w", err)
	}
	body.SetGraffiti(graffiti)

	if s.chainSpec.ActiveForkVersionForEpoch(s.chainSpec.SlotToEpoch(blk.GetSlot())) >= version.DenebPlus {
		body.SetAttestations(slotData.GetAttestationData())
		body.SetSlashingInfo(slotData.GetSlashingInfo())
	}

	body.SetExecutionPayload(envelope.GetExecutionPayload())
	return nil
}

// computeAndSetStateRoot computes the state root of an outgoing block and sets it in the block.
func (s *Service[
	_, BeaconBlockT, _, BeaconStateT, _, _, _, _, _, _, _, _, _,
]) computeAndSetStateRoot(ctx context.Context, st BeaconStateT, blk BeaconBlockT) error {
	stateRoot, err := s.computeStateRoot(ctx, st, blk)
	if err != nil {
		s.logger.Error("failed to compute state root", "slot", blk.GetSlot().Base10(), "error", err)
		return err
	}
	blk.SetStateRoot(stateRoot)
	return nil
}

// computeStateRoot computes the state root of an outgoing block.
func (s *Service[
	_, BeaconBlockT, _, BeaconStateT, _, _, _, _, _, _, _, _, _,
]) computeStateRoot(ctx context.Context, st BeaconStateT, blk BeaconBlockT) (common.Root, error) {
	defer s.metrics.measureStateRootComputationTime(time.Now())
	if _, err := s.stateProcessor.Transition(&transition.Context{
		Context:                 ctx,
		OptimisticEngine:        true,
		SkipPayloadVerification: true,
	}, st, blk); err != nil {
		return common.Root{}, err
	}

	return st.HashTreeRoot(), nil
}
