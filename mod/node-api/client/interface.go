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

package client

import (
	client "github.com/attestantio/go-eth2-client"
)

// BeaconKitNode is a wrapper around the client.Service interface to add
// additional methods specific to a beacon-kit node's API.
type BeaconKitNode interface {
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
