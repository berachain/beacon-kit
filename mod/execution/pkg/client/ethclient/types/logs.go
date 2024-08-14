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

package types

import (
	"errors"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
)

// FilterArgs is a map of string to any, used to configure a filter query.
type FilterArgs map[string]any

// New creates a new FilterArgs with the given addresses and topics.
func (f FilterArgs) New(
	addrs []common.ExecutionAddress, topics []common.ExecutionHash,
) FilterArgs {
	return FilterArgs{
		"address": addrs,
		"topics":  topics,
	}
}

// SetFromBlock sets the "fromBlock" field in the FilterArgs, ensuring that
// BlockHash is not also set.
func (f FilterArgs) SetFromBlock(block BlockNumber) error {
	var err error
	if f["fromBlock"], err = block.MarshalText(); err != nil {
		return err
	}
	if f["blockHash"] != nil {
		return errors.New("cannot specify both BlockHash and FromBlock/ToBlock")
	}
	return nil
}

// SetToBlock sets the "toBlock" field in the FilterArgs, ensuring that
// BlockHash is not also set.
func (f FilterArgs) SetToBlock(block BlockNumber) error {
	var err error
	if f["toBlock"], err = block.MarshalText(); err != nil {
		return err
	}
	if f["blockHash"] != nil {
		return errors.New("cannot specify both BlockHash and FromBlock/ToBlock")
	}
	return nil
}

// SetBlockHash sets the "blockHash" field in the FilterArgs, ensuring that
// FromBlock and ToBlock are not also set.
func (f FilterArgs) SetBlockHash(hash common.ExecutionHash) error {
	if f["fromBlock"] != nil || f["toBlock"] != nil {
		return errors.New("cannot specify both BlockHash and FromBlock/ToBlock")
	}
	f["blockHash"] = hash.Hex()
	return nil
}
