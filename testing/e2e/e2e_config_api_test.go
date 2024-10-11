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

package e2e_test

import (
	"encoding/hex"
	"strconv"

	beaconapi "github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// TestConfigSpec tests the node api for config spec of the node.
func (s *BeaconKitE2ESuite) TestConfigSpec() {
	client := s.initNodeTest()

	spec, err := client.Spec(s.Ctx(),
		&beaconapi.SpecOpts{})
	s.Require().NoError(err)
	s.Require().NotNil(spec)
	specData := spec.Data
	s.Require().NotNil(specData)

	depositContractAddress, exists := spec.Data["DEPOSIT_CONTRACT_ADDRESS"]
	s.Require().True(exists,
		"DEPOSIT_CONTRACT_ADDRESS not found in spec data",
	)

	depositContractAddressHex := "0x" + hex.EncodeToString(
		depositContractAddress.([]byte),
	)
	s.Require().Equal("0x4242424242424242424242424242424242424242",
		depositContractAddressHex,
	)

	depositNetworkID, exists := spec.Data["DEPOSIT_NETWORK_ID"]
	s.Require().True(exists,
		"DEPOSIT_NETWORK_ID not found in spec data",
	)

	networkIDUint64, ok := depositNetworkID.(uint64)
	s.Require().True(ok, "DEPOSIT_NETWORK_ID is not a uint64")
	networkIDString := strconv.FormatUint(networkIDUint64, 10)
	s.Require().Equal("80087", networkIDString)

	domainAggregateAndProof, exists := spec.Data["DOMAIN_AGGREGATE_AND_PROOF"]
	s.Require().True(exists,
		"DOMAIN_AGGREGATE_AND_PROOF not found in spec data",
	)
	expectedDomain := phase0.DomainType{0x6, 0x0, 0x0, 0x0}
	s.Require().Equal(expectedDomain, domainAggregateAndProof)

	inactivityPenaltyQuot, exists := spec.Data["INACTIVITY_PENALTY_QUOTIENT"]
	s.Require().True(exists,
		"INACTIVITY_PENALTY_QUOTIENT not found in spec data",
	)
	s.Require().Equal(uint64(0), inactivityPenaltyQuot)

	inactivityPenaltyQuotAltair, exists :=
		spec.Data["INACTIVITY_PENALTY_QUOTIENT_ALTAIR"]
	s.Require().True(exists,
		"INACTIVITY_PENALTY_QUOTIENT_ALTAIR not found in spec data",
	)
	s.Require().Equal(uint64(0), inactivityPenaltyQuotAltair)
}
