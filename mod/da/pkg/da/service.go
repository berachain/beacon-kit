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

	"github.com/berachain/beacon-kit/mod/async/pkg/event"
	asynctypes "github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/events"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// ServiceName is the name of the service.
const ServiceName = "da"

type Service[
	AvailabilityStoreT AvailabilityStore[BeaconBlockBodyT, BlobSidecarsT],
	BeaconBlockBodyT BeaconBlockBody[ExecutionPayloadT],
	BlobSidecarsT,
	ExecutionPayloadT any,
] struct {
	avs AvailabilityStoreT
	bp  BlobProcessor[
		AvailabilityStoreT, BeaconBlockBodyT,
		BlobSidecarsT, ExecutionPayloadT,
	]
	feed   *event.FeedOf[asynctypes.EventID, *asynctypes.Event[BlobSidecarsT]]
	logger log.Logger[any]
}

// New returns a new DA service.
func NewService[
	AvailabilityStoreT AvailabilityStore[BeaconBlockBodyT, BlobSidecarsT],
	BeaconBlockBodyT BeaconBlockBody[ExecutionPayloadT],
	BlobSidecarsT,
	ExecutionPayloadT any,
](
	avs AvailabilityStoreT,
	bp BlobProcessor[
		AvailabilityStoreT, BeaconBlockBodyT,
		BlobSidecarsT, ExecutionPayloadT,
	],
	feed *event.FeedOf[asynctypes.EventID, *asynctypes.Event[BlobSidecarsT]],
	logger log.Logger[any],
) *Service[
	AvailabilityStoreT, BeaconBlockBodyT, BlobSidecarsT, ExecutionPayloadT,
] {
	return &Service[
		AvailabilityStoreT, BeaconBlockBodyT, BlobSidecarsT, ExecutionPayloadT,
	]{
		avs:  avs,
		bp:   bp,
		feed: feed,
	}
}

// Name returns the name of the service.
func (s *Service[_, _, _, _]) Name() string {
	return ServiceName
}

// Start starts the service.
func (s *Service[_, _, _, _]) Start(ctx context.Context) error {
	go s.start(ctx)
	return nil
}

// start starts the service.
func (s *Service[_, _, BlobSidecarsT, _]) start(ctx context.Context) error {
	ch := make(chan *asynctypes.Event[BlobSidecarsT])
	sub := s.feed.Subscribe(ch)
	defer sub.Unsubscribe()
	for {
		select {
		case <-ctx.Done():
			return nil
		case e := <-ch:
			if e.Type() == events.BlobSidecarsReceived {
				if err := s.processBlobSidecars(ctx, 0, e.Data()); err != nil {
					s.logger.Error(
						"failed to process blob sidecars",
						"error",
						err,
					)
				}
			}
		}
	}
}

// ProcessBlobSidecars processes the blob sidecars.
// TODO: Deprecate this publically and move to event based system.
func (s *Service[_, _, BlobSidecarsT, _]) ProcessBlobSidecars(
	ctx context.Context,
	slot math.Slot,
	sidecars BlobSidecarsT,
) error {
	return s.processBlobSidecars(ctx, slot, sidecars)
}

// processBlobSidecars processes the blob sidecars.
func (s *Service[_, _, BlobSidecarsT, _]) processBlobSidecars(
	_ context.Context,
	slot math.Slot,
	sidecars BlobSidecarsT,
) error {
	// startTime := time.Now()
	// defer s.metrics.measureBlobProcessingDuration(startTime)
	return s.bp.ProcessBlobs(
		slot,
		s.avs,
		sidecars,
	)
}
