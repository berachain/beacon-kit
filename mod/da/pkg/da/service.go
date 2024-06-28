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

package da

import (
	"context"

	asynctypes "github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/events"
)

type Service[
	AvailabilityStoreT AvailabilityStore[BeaconBlockBodyT, BlobSidecarsT],
	BeaconBlockBodyT any,
	BlobSidecarsT BlobSidecar,
	//nolint:lll // formatter.
	EventPublisherSubscriberT EventPublisherSubscriber[*asynctypes.Event[BlobSidecarsT]],
	ExecutionPayloadT any,
] struct {
	avs AvailabilityStoreT
	bp  BlobProcessor[
		AvailabilityStoreT, BeaconBlockBodyT,
		BlobSidecarsT, ExecutionPayloadT,
	]
	sidecarsBroker EventPublisherSubscriberT
	logger         log.Logger[any]
}

// NewService returns a new DA service.
func NewService[
	AvailabilityStoreT AvailabilityStore[
		BeaconBlockBodyT, BlobSidecarsT,
	],
	BeaconBlockBodyT any,
	BlobSidecarsT BlobSidecar,
	//nolint:lll // formatter.
	EventPublisherSubscriberT EventPublisherSubscriber[*asynctypes.Event[BlobSidecarsT]],
	ExecutionPayloadT any,
](
	avs AvailabilityStoreT,
	bp BlobProcessor[
		AvailabilityStoreT, BeaconBlockBodyT,
		BlobSidecarsT, ExecutionPayloadT,
	],
	sidecarsBroker EventPublisherSubscriberT,
	logger log.Logger[any],
) *Service[
	AvailabilityStoreT, BeaconBlockBodyT,
	BlobSidecarsT, EventPublisherSubscriberT, ExecutionPayloadT,
] {
	return &Service[
		AvailabilityStoreT, BeaconBlockBodyT,
		BlobSidecarsT, EventPublisherSubscriberT, ExecutionPayloadT,
	]{
		avs:            avs,
		bp:             bp,
		sidecarsBroker: sidecarsBroker,
		logger:         logger,
	}
}

// Name returns the name of the service.
func (s *Service[_, _, _, _, _]) Name() string {
	return "da"
}

// Start starts the service.
func (s *Service[_, _, _, _, _]) Start(ctx context.Context) error {
	subSidecarsCh, err := s.sidecarsBroker.Subscribe()
	if err != nil {
		return err
	}
	go s.start(ctx, subSidecarsCh)
	return nil
}

// start starts the service.
func (s *Service[_, _, BlobSidecarsT, _, _]) start(
	ctx context.Context,
	sidecarsCh chan *asynctypes.Event[BlobSidecarsT],
) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-sidecarsCh:
			switch msg.Type() {
			case events.BlobSidecarsProcessRequest:
				s.handleBlobSidecarsProcessRequest(msg)
			case events.BlobSidecarsReceived:
				s.handleBlobSidecarsReceived(msg)
			}
		}
	}
}

// handleBlobSidecarsProcessRequest handles the BlobSidecarsProcessRequest
// event.
// It processes the sidecars and publishes a BlobSidecarsProcessed event.
func (s *Service[_, _, BlobSidecarsT, _, _]) handleBlobSidecarsProcessRequest(
	msg *asynctypes.Event[BlobSidecarsT],
) {
	err := s.processSidecars(msg.Context(), msg.Data())
	if err != nil {
		s.logger.Error(
			"Failed to process blob sidecars",
			"error",
			err,
		)
	}

	if err = s.sidecarsBroker.Publish(
		msg.Context(),
		asynctypes.NewEvent(
			msg.Context(), events.BlobSidecarsProcessed, msg.Data(), err,
		)); err != nil {
		s.logger.Error(
			"Failed to publish blob sidecars processed event",
			"error",
			err,
		)
	}
}

// handleBlobSidecarsReceived handles the BlobSidecarsReceived event.
// It receives the sidecars and publishes a BlobSidecarsProcessed event.
func (s *Service[_, _, BlobSidecarsT, _, _]) handleBlobSidecarsReceived(
	msg *asynctypes.Event[BlobSidecarsT],
) {
	err := s.receiveSidecars(msg.Data())
	if err != nil {
		s.logger.Error(
			"Failed to receive blob sidecars",
			"error",
			err,
		)
	}

	if err = s.sidecarsBroker.Publish(
		msg.Context(),
		asynctypes.NewEvent(
			msg.Context(), events.BlobSidecarsProcessed, msg.Data(), err,
		)); err != nil {
		s.logger.Error(
			"Failed to publish blob sidecars processed event",
			"error",
			err,
		)
	}
}

// ProcessSidecars processes the blob sidecars.
func (s *Service[_, _, BlobSidecarsT, _, _]) processSidecars(
	_ context.Context,
	sidecars BlobSidecarsT,
) error {
	// startTime := time.Now()
	// defer s.metrics.measureBlobProcessingDuration(startTime)
	return s.bp.ProcessSidecars(
		s.avs,
		sidecars,
	)
}

// VerifyIncomingBlobs receives blobs from the network and processes them.
func (s *Service[_, _, BlobSidecarsT, _, _]) receiveSidecars(
	sidecars BlobSidecarsT,
) error {
	// If there are no blobs to verify, return early.
	if sidecars.IsNil() || sidecars.Len() == 0 {
		return nil
	}

	s.logger.Info(
		"Received incoming blob sidecars 🚔",
	)

	// Verify the blobs and ensure they match the local state.
	if err := s.bp.VerifySidecars(sidecars); err != nil {
		s.logger.Error(
			"rejecting incoming blob sidecars ❌",
			"reason", err,
		)
		return err
	}

	s.logger.Info(
		"Blob sidecars verification succeeded - accepting incoming blob sidecars 💦",
		"num_blobs",
		sidecars.Len(),
	)

	return nil
}
