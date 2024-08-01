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

	asynctypes "github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/events"
)

// Service is the validator service.
// It is responsible for building blocks and sidecars on new slots.
type Service[
	AttestationDataT any,
	BeaconBlockT BeaconBlock[
		AttestationDataT, BeaconBlockT, BeaconBlockBodyT, DepositT,
		Eth1DataT, ExecutionPayloadT, SlashingInfoT,
	],
	BeaconBlockBodyT BeaconBlockBody[
		AttestationDataT, DepositT, Eth1DataT, ExecutionPayloadT, SlashingInfoT,
	],
	BlobSidecarsT,
	DepositT any,
	Eth1DataT Eth1Data[Eth1DataT],
	ExecutionPayloadT any,
	SlashingInfoT any,
	SlotDataT SlotData[AttestationDataT, SlashingInfoT],
] struct {
	processor Processor[
		AttestationDataT,
		BeaconBlockT,
		BeaconBlockBodyT,
		BlobSidecarsT,
		DepositT,
		Eth1DataT,
		ExecutionPayloadT,
		SlashingInfoT,
		SlotDataT,
	]
	logger log.Logger[any]
	// blkBroker is a publisher for blocks.
	blkBroker EventPublisher[*asynctypes.Event[BeaconBlockT]]
	// sidecarBroker is a publisher for sidecars.
	sidecarBroker EventPublisher[*asynctypes.Event[BlobSidecarsT]]
	// newSlotSub is a feed for slots.
	slotBroker EventFeed[*asynctypes.Event[SlotDataT]]
}

// NewService creates a new validator service.
func NewService[
	AttestationDataT any,
	BeaconBlockT BeaconBlock[
		AttestationDataT, BeaconBlockT, BeaconBlockBodyT, DepositT,
		Eth1DataT, ExecutionPayloadT, SlashingInfoT,
	],
	BeaconBlockBodyT BeaconBlockBody[
		AttestationDataT, DepositT, Eth1DataT, ExecutionPayloadT, SlashingInfoT,
	],
	BlobSidecarsT,
	DepositT any,
	Eth1DataT Eth1Data[Eth1DataT],
	ExecutionPayloadT any,
	SlashingInfoT any,
	SlotDataT SlotData[AttestationDataT, SlashingInfoT],
](
	processor Processor[
		AttestationDataT, BeaconBlockT, BeaconBlockBodyT, BlobSidecarsT,
		DepositT, Eth1DataT, ExecutionPayloadT, SlashingInfoT, SlotDataT,
	],
	blkBroker EventPublisher[*asynctypes.Event[BeaconBlockT]],
	sidecarBroker EventPublisher[*asynctypes.Event[BlobSidecarsT]],
	slotBroker EventFeed[*asynctypes.Event[SlotDataT]],
) *Service[
	AttestationDataT, BeaconBlockT, BeaconBlockBodyT,
	BlobSidecarsT, DepositT, Eth1DataT, ExecutionPayloadT,
	SlashingInfoT, SlotDataT,
] {
	return &Service[
		AttestationDataT, BeaconBlockT, BeaconBlockBodyT,
		BlobSidecarsT, DepositT, Eth1DataT, ExecutionPayloadT,
		SlashingInfoT, SlotDataT,
	]{
		processor:     processor,
		blkBroker:     blkBroker,
		sidecarBroker: sidecarBroker,
		slotBroker:    slotBroker,
	}
}

// Name returns the name of the service.
func (s *Service[
	_, _, _, _, _, _, _, _, _,
]) Name() string {
	return "validator"
}

// Start initializes the main event loop and listens for new slots.
func (s *Service[
	_, _, _, _, _, _, _, _, _,
]) Start(
	ctx context.Context,
) error {
	subSlotCh, err := s.slotBroker.Subscribe()
	if err != nil {
		return err
	}
	go s.start(ctx, subSlotCh)
	return nil
}

// start listens for new slots and handles them accordingly.
func (s *Service[
	_, _, _, _, _, _, _, _, SlotDataT,
]) start(
	ctx context.Context,
	subSlotCh chan *asynctypes.Event[SlotDataT],
) {
	for {
		select {
		case <-ctx.Done():
			return
		case req := <-subSlotCh:
			if req.Type() == events.NewSlot {
				s.handleNewSlot(req)
			}
		}
	}
}

// handleNewSlot synchronously builds a new block and sidecars for the given
// slot, and publishes the built block and sidecars.
func (s *Service[
	_, _, _, _, _, _, _, _, SlotDataT,
]) handleNewSlot(msg *asynctypes.Event[SlotDataT]) {
	blk, sidecars, err := s.processor.BuildBlockAndSidecars(
		msg.Context(), msg.Data(),
	)
	if err != nil {
		s.logger.Error("failed to build block", "err", err)
	}

	// Publish our built block to the broker.
	if blkErr := s.blkBroker.Publish(
		msg.Context(),
		asynctypes.NewEvent(
			msg.Context(), events.BeaconBlockBuilt, blk, err,
		)); blkErr != nil {
		// Propagate the error from buildBlockAndSidecars
		s.logger.Error("failed to publish block", "err", err)
	}

	// Publish our built blobs to the broker.
	if sidecarsErr := s.sidecarBroker.Publish(
		msg.Context(),
		asynctypes.NewEvent(
			// Propagate the error from buildBlockAndSidecars
			msg.Context(), events.BlobSidecarsBuilt, sidecars, err,
		),
	); sidecarsErr != nil {
		s.logger.Error("failed to publish sidecars", "err", err)
	}
}
