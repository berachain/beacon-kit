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

package beacondb

// GetLatestExecutionPayloadHeader retrieves the latest execution payload
// header from the BeaconStore.
func (s *StateManager[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT,
]) GetLatestExecutionPayloadHeader() (
	ExecutionPayloadHeaderT, error,
) {
	forkVersion, err := s.latestExecutionPayloadVersion.Get()
	if err != nil {
		var t ExecutionPayloadHeaderT
		return t, err
	}
	s.latestExecutionPayloadCodec.SetActiveForkVersion(forkVersion)
	return s.latestExecutionPayloadHeader.Get()
}

// SetLatestExecutionPayloadHeader sets the latest execution payload header in
// the BeaconStore.
func (s *StateManager[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT,
]) SetLatestExecutionPayloadHeader(
	payloadHeader ExecutionPayloadHeaderT,
) error {
	if err := s.latestExecutionPayloadVersion.Set(
		payloadHeader.Version(),
	); err != nil {
		return err
	}
	s.latestExecutionPayloadCodec.SetActiveForkVersion(payloadHeader.Version())
	return s.latestExecutionPayloadHeader.Set(payloadHeader)
}

// GetEth1DepositIndex retrieves the eth1 deposit index from the beacon state.
func (s *StateManager[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT,
]) GetEth1DepositIndex() (uint64, error) {
	return s.eth1DepositIndex.Get()
}

// SetEth1DepositIndex sets the eth1 deposit index in the beacon state.
func (s *StateManager[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT,
]) SetEth1DepositIndex(
	index uint64,
) error {
	return s.eth1DepositIndex.Set(index)
}

// GetEth1Data retrieves the eth1 data from the beacon state.
func (s *StateManager[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT,
]) GetEth1Data() (Eth1DataT, error) {
	return s.eth1Data.Get()
}

// SetEth1Data sets the eth1 data in the beacon state.
func (s *StateManager[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT,
]) SetEth1Data(
	data Eth1DataT,
) error {
	return s.eth1Data.Set(data)
}
