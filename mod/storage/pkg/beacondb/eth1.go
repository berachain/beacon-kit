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

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/berachain/beacon-kit/mod/storage/pkg/encoding"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetLatestExecutionPayloadHeader retrieves the latest execution payload
// header from the BeaconStore.
func (kv *KVStore[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
]) GetLatestExecutionPayloadHeader() (
	ExecutionPayloadHeaderT, error,
) {
	fmt.Printf(
		"****** GetLatestExecutionPayloadHeader stage=%d\n",
		debugStageID(kv.ctx),
	)
	forkVersion, err := kv.latestExecutionPayloadVersion.Get(kv.ctx)
	if err != nil {
		var t ExecutionPayloadHeaderT
		return t, err
	}
	kv.latestExecutionPayloadCodec.SetActiveForkVersion(forkVersion)

	header, err := kv.latestExecutionPayloadHeader.Get(kv.ctx)
	if err != nil {
		return header, err
	}
	headerSSZ, err := kv.sszDB.GetPath(
		kv.ctx,
		"latest_execution_payload_header",
	)
	if err != nil {
		return header, err
	}

	headerBytes, err := header.MarshalSSZ()
	if err != nil {
		return header, err
	}
	if !bytes.Equal(headerBytes, headerSSZ) {
		h := sha256.Sum256(headerBytes)
		h2 := sha256.Sum256(headerSSZ)
		return header, fmt.Errorf(
			"latest execution payload header SSZ does not match DB; headerSSZ=%x, header=%x",
			h2,
			h,
		)
	}
	return header, nil
}

func debugStageID(ctx context.Context) uint8 {
	const contextlessContext = 77
	sdkCtx, ok := sdk.TryUnwrapSDKContext(ctx)
	if !ok {
		return contextlessContext
	}
	return uint8(sdkCtx.ExecMode())
}

// SetLatestExecutionPayloadHeader sets the latest execution payload header in
// the BeaconStore.
func (kv *KVStore[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
]) SetLatestExecutionPayloadHeader(
	payloadHeader ExecutionPayloadHeaderT,
) error {
	bz, err := payloadHeader.MarshalSSZ()
	if err != nil {
		return err
	}
	fmt.Printf(
		"****** SetLatestExecutionPayloadHeader: stage=%d hash=%x\n",
		debugStageID(kv.ctx),
		sha256.Sum256(bz),
	)
	if err := kv.latestExecutionPayloadVersion.Set(
		kv.ctx, payloadHeader.Version(),
	); err != nil {
		return err
	}
	kv.latestExecutionPayloadCodec.SetActiveForkVersion(payloadHeader.Version())
	err = kv.sszDB.SetLatestExecutionPayloadHeader(kv.ctx, payloadHeader)
	if err != nil {
		return err
	}
	return kv.latestExecutionPayloadHeader.Set(kv.ctx, payloadHeader)
}

// GetEth1DepositIndex retrieves the eth1 deposit index from the beacon state.
func (kv *KVStore[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
]) GetEth1DepositIndex() (uint64, error) {
	return kv.eth1DepositIndex.Get(kv.ctx)
}

// SetEth1DepositIndex sets the eth1 deposit index in the beacon state.
func (kv *KVStore[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
]) SetEth1DepositIndex(
	index uint64,
) error {
	return kv.eth1DepositIndex.Set(kv.ctx, index)
}

// GetEth1Data retrieves the eth1 data from the beacon state.
func (kv *KVStore[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
]) GetEth1Data() (Eth1DataT, error) {
	eth1Data, err := kv.eth1Data.Get(kv.ctx)
	if err != nil {
		return eth1Data, err
	}
	eth1DataSSZ, err := kv.sszDB.GetPath(kv.ctx, "eth1_data")
	if err != nil {
		return eth1Data, err
	}
	bz, err := eth1Data.MarshalSSZ()
	if err != nil {
		return eth1Data, err
	}
	if !bytes.Equal(bz, eth1DataSSZ) {
		return eth1Data, errors.New("eth1 data SSZ does not match DB")
	}
	// TODO: Don't alloc codec, some codec exists already buried in sdk.Collections
	codec := encoding.SSZValueCodec[Eth1DataT]{}
	res, err := codec.Decode(eth1DataSSZ)
	if err != nil {
		return eth1Data, err
	}
	return res, nil
}

// SetEth1Data sets the eth1 data in the beacon state.
func (kv *KVStore[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
]) SetEth1Data(
	data Eth1DataT,
) error {
	if err := kv.sszDB.SetObject(kv.ctx, "eth1_data", data); err != nil {
		return err
	}
	return kv.eth1Data.Set(kv.ctx, data)
}
