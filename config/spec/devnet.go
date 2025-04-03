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

package spec

import (
	"github.com/berachain/beacon-kit/chain"
)

const (
	// devnetDeneb1ForkTime is the timestamp at which the Deneb1 fork occurs.
	// A value of 64 is set for the fork to occur 64 seconds into the test,
	// which is approximately 1 epoch.
	// TODO(fork): Make devnet fork time take place during each devnet test.
	//  devnetDeneb1ForkTime is the UNIX timestamp of when the deneb1 fork occurs. If this value
	//  is 64, it means that the deneb1 fork occurred 64 seconds into 1970 at the start of UNIX
	//  time. Since e2e tests start at `time.Now()`, we would be significantly past the fork time.
	//  To make this fork time happen during an e2e test, we must either:
	//     a. find a way to set the fork time as an offset from genesis time
	//     b. trick `time.Now()` into thinking it's at the UNIX Epoch.
	devnetDeneb1ForkTime = 1 * defaultSlotsPerEpoch * defaultTargetSecondsPerEth1Block

	// devnetElectraForkTime is the timestamp at which the Electra fork occurs.
	devnetElectraForkTime = defaultElectraForkTime
)

// DevnetChainSpecData is the chain.SpecData for a devnet.
func DevnetChainSpecData() *chain.SpecData {
	specData := MainnetChainSpecData()
	specData.DepositEth1ChainID = DevnetEth1ChainID

	// Fork timings are set to facilitate local testing across fork versions.
	specData.Deneb1ForkTime = devnetDeneb1ForkTime
	specData.ElectraForkTime = devnetElectraForkTime

	return specData
}

// DevnetChainSpec is the chain.Spec for a devnet. Used by `make start` and unit tests.
func DevnetChainSpec() (chain.Spec, error) {
	return chain.NewSpec(DevnetChainSpecData())
}
