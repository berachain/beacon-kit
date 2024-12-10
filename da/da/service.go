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

	asynctypes "github.com/berachain/beacon-kit/async/types"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/primitives/async"
)

// The Data Availability service is responsible for verifying and processing
// incoming blob sidecars.
//

type Service[
	AvailabilityStoreT any,
	ConsensusSidecarsT ConsensusSidecars[BlobSidecarsT, BeaconBlockHeaderT],
	BlobSidecarsT BlobSidecar,
	BeaconBlockHeaderT any,
] struct {
	avs AvailabilityStoreT
	bp  BlobProcessor[
		AvailabilityStoreT,
		ConsensusSidecarsT, BlobSidecarsT,
	]
	dispatcher asynctypes.EventDispatcher
	logger     log.Logger
	// subFinalBlobSidecars is a channel holding FinalSidecarsReceived events.
	subFinalBlobSidecars chan async.Event[BlobSidecarsT]
}

// NewService returns a new DA service.
func NewService[
	AvailabilityStoreT any,
	ConsensusSidecarsT ConsensusSidecars[BlobSidecarsT, BeaconBlockHeaderT],
	BlobSidecarsT BlobSidecar,
	BeaconBlockHeaderT any,
](
	avs AvailabilityStoreT,
	bp BlobProcessor[
		AvailabilityStoreT,
		ConsensusSidecarsT, BlobSidecarsT,
	],
	dispatcher asynctypes.EventDispatcher,
	logger log.Logger,
) *Service[
	AvailabilityStoreT, ConsensusSidecarsT, BlobSidecarsT, BeaconBlockHeaderT,
] {
	return &Service[
		AvailabilityStoreT,
		ConsensusSidecarsT, BlobSidecarsT, BeaconBlockHeaderT,
	]{
		avs:                  avs,
		bp:                   bp,
		dispatcher:           dispatcher,
		logger:               logger,
		subFinalBlobSidecars: make(chan async.Event[BlobSidecarsT]),
	}
}

// Name returns the name of the service.
func (s *Service[_, _, _, _]) Name() string {
	return "da"
}

// Start subscribes the DA service to SidecarsReceived and FinalSidecarsReceived
// events and begins the main event loop to handle them accordingly.
func (s *Service[_, _, _, _]) Start(ctx context.Context) error {
	var err error
	// subscribe to FinalSidecarsReceived events
	if err = s.dispatcher.Subscribe(
		async.FinalSidecarsReceived, s.subFinalBlobSidecars,
	); err != nil {
		return err
	}

	// start the main event loop to listen and handle events.
	go s.eventLoop(ctx)
	return nil
}

// eventLoop listens and handles SidecarsReceived and FinalSidecarsReceived
// events.
func (s *Service[_, _, _, _]) eventLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-s.subFinalBlobSidecars:
			s.handleFinalSidecarsReceived(event)
		}
	}
}

/* -------------------------------------------------------------------------- */
/*                               Event Handlers                             */
/* -------------------------------------------------------------------------- */

// handleFinalSidecarsReceived handles the BlobSidecarsProcessRequest
// event.
// It processes the sidecars and publishes a BlobSidecarsProcessed event.
func (s *Service[_, _, BlobSidecarsT, _]) handleFinalSidecarsReceived(
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
