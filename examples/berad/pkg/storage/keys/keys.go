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

package keys

const (
	WithdrawalQueuePrefix byte = iota
	RandaoMixPrefix
	ValidatorIndexPrefix
	BlockRootsPrefix
	StateRootsPrefix
	ValidatorByIndexPrefix
	ValidatorPubkeyToIndexPrefix
	ValidatorConsAddrToIndexPrefix
	ValidatorEffectiveBalanceToIndexPrefix
	LatestBeaconBlockHeaderPrefix
	SlotPrefix
	BalancesPrefix
	Eth1BlockHashPrefix
	Eth1DepositIndexPrefix
	LatestExecutionPayloadHeaderPrefix
	LatestExecutionPayloadVersionPrefix
	GenesisValidatorsRootPrefix
	NextWithdrawalIndexPrefix
	NextWithdrawalValidatorIndexPrefix
	ForkPrefix
)

//nolint:lll
const (
	WithdrawalQueuePrefixHumanReadable                  = "WithdrawalQueuePrefix"
	RandaoMixPrefixHumanReadable                        = "RandaoMixPrefix"
	ValidatorIndexPrefixHumanReadable                   = "ValidatorIndexPrefix"
	BlockRootsPrefixHumanReadable                       = "BlockRootsPrefix"
	StateRootsPrefixHumanReadable                       = "StateRootsPrefix"
	ValidatorByIndexPrefixHumanReadable                 = "ValidatorByIndexPrefix"
	ValidatorPubkeyToIndexPrefixHumanReadable           = "ValidatorPubkeyToIndexPrefix"
	ValidatorConsAddrToIndexPrefixHumanReadable         = "ValidatorConsAddrToIndexPrefix"
	ValidatorEffectiveBalanceToIndexPrefixHumanReadable = "ValidatorEffectiveBalanceToIndexPrefix"
	LatestBeaconBlockHeaderPrefixHumanReadable          = "LatestBeaconBlockHeaderPrefix"
	SlotPrefixHumanReadable                             = "SlotPrefix"
	BalancesPrefixHumanReadable                         = "BalancesPrefix"
	Eth1BlockHashPrefixHumanReadable                    = "Eth1BlockHashPrefix"
	Eth1DepositIndexPrefixHumanReadable                 = "Eth1DepositIndexPrefix"
	LatestExecutionPayloadHeaderPrefixHumanReadable     = "LatestExecutionPayloadHeaderPrefix"
	LatestExecutionPayloadVersionPrefixHumanReadable    = "LatestExecutionPayloadVersionPrefix"
	GenesisValidatorsRootPrefixHumanReadable            = "GenesisValidatorsRootPrefix"
	NextWithdrawalIndexPrefixHumanReadable              = "NextWithdrawalIndexPrefix"
	NextWithdrawalValidatorIndexPrefixHumanReadable     = "NextWithdrawalValidatorIndexPrefix"
	ForkPrefixHumanReadable                             = "ForkPrefix"
)
