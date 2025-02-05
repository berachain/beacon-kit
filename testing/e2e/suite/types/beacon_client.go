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
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	client "github.com/attestantio/go-eth2-client"
	beaconapi "github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	beaconhttp "github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// BeaconKitNodeClient is a wrapper around the client.Service interface to add
// additional methods specific to a beacon-kit node's API.
type BeaconKitNodeClient interface {
	client.Service

	client.FarFutureEpochProvider
	client.SignedBeaconBlockProvider
	client.BlobSidecarsProvider
	client.BeaconCommitteesProvider
	client.SyncCommitteesProvider
	client.AggregateAttestationProvider
	client.AggregateAttestationsSubmitter
	client.AttestationDataProvider
	client.AttestationPoolProvider
	client.AttestationsSubmitter
	client.AttesterSlashingSubmitter
	client.AttesterDutiesProvider
	client.SyncCommitteeDutiesProvider
	client.SyncCommitteeMessagesSubmitter
	client.SyncCommitteeSubscriptionsSubmitter
	client.SyncCommitteeContributionProvider
	client.SyncCommitteeContributionsSubmitter
	client.BLSToExecutionChangesSubmitter
	client.BeaconBlockHeadersProvider
	client.ProposalProvider
	client.ProposalSlashingSubmitter
	client.BeaconBlockRootProvider
	client.ProposalSubmitter
	client.BeaconCommitteeSubscriptionsSubmitter
	client.BeaconStateProvider
	client.BeaconStateRandaoProvider
	client.BeaconStateRootProvider
	client.BlindedProposalSubmitter
	client.ValidatorRegistrationsSubmitter
	client.EventsProvider
	client.FinalityProvider
	client.ForkChoiceProvider
	client.ForkProvider
	client.ForkScheduleProvider
	client.GenesisProvider
	client.NodePeersProvider
	client.NodeSyncingProvider
	client.NodeVersionProvider
	client.ProposalPreparationsSubmitter
	client.ProposerDutiesProvider
	client.SpecProvider
	client.ValidatorBalancesProvider
	client.ValidatorsProvider
	client.VoluntaryExitSubmitter
	client.VoluntaryExitPoolProvider
	client.DomainProvider
	client.NodeClientProvider

	// Other beacon-kit node-api methods here...
}

// CustomBeaconClient is a custom implementation of the BeaconKitNodeClient interface
// that overrides the Validators method to handle deneb1.
type CustomBeaconClient struct {
	*beaconhttp.Service
	address string
	client  *http.Client
}

// NewBeaconKitNodeClient creates a new beacon-kit node-api client instance
// with the given cancel context.
func NewBeaconKitNodeClient(
	cancelCtx context.Context,
	params ...beaconhttp.Parameter,
) (BeaconKitNodeClient, error) {
	service, err := beaconhttp.New(
		cancelCtx,
		params...,
	)
	if err != nil {
		return nil, err
	}
	if service == nil {
		return nil, errors.New("service is nil")
	}

	address := service.Address()
	if address == "" {
		return nil, errors.New("no address specified")
	}

	// Type assert service to beaconhttp.Service
	httpService, ok := service.(*beaconhttp.Service)
	if !ok {
		return nil, errors.New("failed to cast service to beaconhttp.Service")
	}

	return &CustomBeaconClient{
		Service: httpService,
		address: address,
		client:  &http.Client{},
	}, nil
}

// Validators implements a custom validator query that handles deneb1.
func (c *CustomBeaconClient) Validators(
	ctx context.Context,
	opts *beaconapi.ValidatorsOpts,
) (*beaconapi.Response[map[phase0.ValidatorIndex]*apiv1.Validator], error) {
	if opts == nil {
		return nil, errors.New("no options specified")
	}
	if opts.State == "" {
		return nil, errors.New("no state specified")
	}

	// Construct the URL
	url := fmt.Sprintf("%s/eth/v1/beacon/states/%s/validators", c.address, opts.State)

	// Add indices as query parameters if specified
	if len(opts.Indices) > 0 {
		indices := make([]string, len(opts.Indices))
		for i, idx := range opts.Indices {
			indices[i] = fmt.Sprintf("id=%d", idx)
		}
		url = fmt.Sprintf("%s?%s", url, strings.Join(indices, "&"))
	}

	// Make GET request for empty indices
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	if resp == nil {
		return nil, errors.New("received nil response")
	}
	defer resp.Body.Close()

	// Parse response
	var result struct {
		Data []*apiv1.Validator `json:"data"`
		Meta struct {
			Count     int    `json:"count"`
			NextToken string `json:"next_token,omitempty"`
			TotalSize int    `json:"total_size"`
		} `json:"meta"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert array to map
	validators := make(map[phase0.ValidatorIndex]*apiv1.Validator)
	for _, v := range result.Data {
		validators[v.Index] = v
	}

	return &beaconapi.Response[map[phase0.ValidatorIndex]*apiv1.Validator]{
		Data: validators,
		Metadata: map[string]any{
			"count":      result.Meta.Count,
			"next_token": result.Meta.NextToken,
			"total_size": result.Meta.TotalSize,
		},
	}, nil
}
