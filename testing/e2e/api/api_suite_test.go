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

package api_test

import (
	"context"

	eth2client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/berachain/beacon-kit/testing/e2e/suite"
	"github.com/rs/zerolog"
)

// BeaconAPISuite is a suite of beacon node-api tests with full simulation of a
// beacon-kit network.
type BeaconAPISuite struct {
	suite.KurtosisE2ESuite

	beaconClient eth2client.Service
}

// NewBeaconAPISuite creates a new BeaconAPISuite.
//
// TODO: pass in values from the config.
func NewBeaconAPISuite() *BeaconAPISuite {
	beaconClient, err := http.New(
		context.Background(),
		http.WithAddress("http://localhost:3500"),
		http.WithLogLevel(zerolog.DebugLevel),
	)
	if err != nil {
		panic(err)
	}

	return &BeaconAPISuite{
		beaconClient: beaconClient,
	}
}
