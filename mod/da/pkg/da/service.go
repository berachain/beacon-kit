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

	async "github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/events"
)

// The Data Availability service is responsible for verifying and processing
// incoming blob sidecars.
//

type Service[
	AvailabilityStoreT AvailabilityStore[BeaconBlockBodyT, BlobSidecarsT],
	BeaconBlockBodyT any,
	BlobSidecarsT BlobSidecar,

	ExecutionPayloadT any,
] struct {
	avs AvailabilityStoreT
	bp  BlobProcessor[
		AvailabilityStoreT, BeaconBlockBodyT,
		BlobSidecarsT, ExecutionPayloadT,
	]
	dispatcher           async.EventDispatcher
	logger               log.Logger[any]
	subSidecarsReceived  async.Subscription[async.Event[BlobSidecarsT]]
	subFinalBlobSidecars async.Subscription[async.Event[BlobSidecarsT]]
}

// NewService returns a new DA service.
func NewService[
	AvailabilityStoreT AvailabilityStore[
		BeaconBlockBodyT, BlobSidecarsT,
	],
	BeaconBlockBodyT any,
	BlobSidecarsT BlobSidecar,

	ExecutionPayloadT any,
](
	avs AvailabilityStoreT,
	bp BlobProcessor[
		AvailabilityStoreT, BeaconBlockBodyT,
		BlobSidecarsT, ExecutionPayloadT,
	],
	dispatcher async.EventDispatcher,
	logger log.Logger[any],
) *Service[
	AvailabilityStoreT, BeaconBlockBodyT,
	BlobSidecarsT, ExecutionPayloadT,
] {
	return &Service[
		AvailabilityStoreT, BeaconBlockBodyT,
		BlobSidecarsT, ExecutionPayloadT,
	]{
		avs:                  avs,
		bp:                   bp,
		dispatcher:           dispatcher,
		logger:               logger,
		subSidecarsReceived:  async.NewSubscription[async.Event[BlobSidecarsT]](),
		subFinalBlobSidecars: async.NewSubscription[async.Event[BlobSidecarsT]](),
	}
}

// Name returns the name of the service.
func (s *Service[_, _, _, _]) Name() string {
	return "da"
}

// Start registers this service as the recipient of ProcessSidecars and
// VerifySidecars messages, and begins listening for these requests.
func (s *Service[_, _, BlobSidecarsT, _]) Start(ctx context.Context) error {
	var err error

	// subscribe to SidecarsReceived events
	if err = s.dispatcher.Subscribe(
		events.SidecarsReceived, s.subSidecarsReceived,
	); err != nil {
		return err
	}

	// subscribe to FinalSidecarsReceived events
	if err = s.dispatcher.Subscribe(
		events.FinalSidecarsReceived, s.subFinalBlobSidecars,
	); err != nil {
		return err
	}

	// listen for events and handle accordingly
	s.subSidecarsReceived.Listen(ctx, s.handleSidecarsVerifyRequest)
	s.subFinalBlobSidecars.Listen(ctx, s.handleBlobSidecarsProcessRequest)
	return nil
}

/* -------------------------------------------------------------------------- */
/*                               Event Handlers                             */
/* -------------------------------------------------------------------------- */

// handleBlobSidecarsProcessRequest handles the BlobSidecarsProcessRequest
// event.
// It processes the sidecars and publishes a BlobSidecarsProcessed event.
func (s *Service[_, _, BlobSidecarsT, _]) handleBlobSidecarsProcessRequest(
	msg async.Event[BlobSidecarsT],
) {
	if err := s.processSidecars(msg.Context(), msg.Data()); err != nil {
		s.logger.Error(
			"Failed to process blob sidecars",
			"error",
			err,
		)
	}
}

// handleSidecarsVerifyRequest handles the SidecarsVerifyRequest event.
// It verifies the sidecars and publishes a SidecarsVerified event.
func (s *Service[_, _, BlobSidecarsT, _]) handleSidecarsVerifyRequest(
	msg async.Event[BlobSidecarsT],
) {
	var sidecarsErr error
	// verify the sidecars.
	if sidecarsErr = s.verifySidecars(msg.Data()); sidecarsErr != nil {
		s.logger.Error(
			"Failed to receive blob sidecars",
			"error",
			sidecarsErr,
		)
	}

	// emit the sidecars verification event with error from verifySidecars
	if err := s.dispatcher.PublishEvent(
		async.NewEvent(
			msg.Context(), events.SidecarsVerified, msg.Data(), sidecarsErr,
		),
	); err != nil {
		s.logger.Error("failed to publish event", "err", err)
	}
}

/* -------------------------------------------------------------------------- */
/*                                   helpers                                  */
/* -------------------------------------------------------------------------- */

// ProcessSidecars processes the blob sidecars.
func (s *Service[_, _, BlobSidecarsT, _]) processSidecars(
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
func (s *Service[_, _, BlobSidecarsT, _]) verifySidecars(
	sidecars BlobSidecarsT,
) error {
	// If there are no blobs to verify, return early.
	if sidecars.IsNil() || sidecars.Len() == 0 {
		return nil
	}

	s.logger.Info(
		"Received incoming blob sidecars",
	)

	// Verify the blobs and ensure they match the local state.
	if err := s.bp.VerifySidecars(sidecars); err != nil {
		s.logger.Error(
			"rejecting incoming blob sidecars",
			"reason", err,
		)
		return err
	}

	s.logger.Info(
		"Blob sidecars verification succeeded - accepting incoming blob sidecars",
		"num_blobs",
		sidecars.Len(),
	)

	return nil
}
