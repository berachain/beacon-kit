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

package vm

import (
	"fmt"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/hashing"
	"github.com/berachain/beacon-kit/mod/consensus/pkg/miniavalanche/encoding"
)

type Validator struct {
	NodeID ids.NodeID
	Weight uint64

	// note no nonce, not validatorTxID to distinguish restaked validators
	// TODO: mind this once dynamic validators set is introduced

	bytes []byte
	id    ids.ID
}

func NewValidator(nodeID ids.NodeID, weight uint64) (*Validator, error) {
	val := &Validator{
		NodeID: nodeID,
		Weight: weight,
	}
	return val, val.initValID()
}

func ParseValidator(valBytes []byte) (*Validator, error) {
	val := &Validator{}
	if err := encoding.Decode(valBytes, &val); err != nil {
		return nil, fmt.Errorf("unable to parse validator: %w", err)
	}

	return val, val.initValID()
}

func (v *Validator) initValID() error {
	bytes, err := encoding.Encode(v)
	if err != nil {
		return fmt.Errorf("failed encoding validator %v: %w", v, err)
	}
	v.bytes = bytes
	v.id = hashing.ComputeHash256Array(v.bytes)
	return nil
}
