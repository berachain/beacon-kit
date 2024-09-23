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
	"encoding/base64"
	"fmt"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/berachain/beacon-kit/mod/consensus/pkg/miniavalanche/encoding"
)

// Base64Genesis is the genesis representation
// that can be fed to the node via cli flags.
type Base64Genesis struct {
	Validators []Base64GenesisValidator
	EthGenesis string
}

type Base64GenesisValidator struct {
	NodeID string
	Weight uint64
}

// ParseBase64StringToBytes is used while parsing cli flags.
func ParseBase64StringToBytes(genesisStr string) ([]byte, error) {
	// Step 1: base64 string to Base64Genesis
	base64Bytes, err := base64.StdEncoding.DecodeString(genesisStr)
	if err != nil {
		return nil, fmt.Errorf("unable to decode base64 genesis content: %w", err)
	}

	base64Gen := &Base64Genesis{}
	if err = encoding.Decode(base64Bytes, base64Gen); err != nil {
		return nil, fmt.Errorf("unable to decode base64 genesis content: %w", err)
	}

	// Step 2: Base64Genesis to InMemoryGenesis
	inMemGen := &Genesis{
		EthGenesis: []byte(base64Gen.EthGenesis),
	}
	for i, v := range base64Gen.Validators {
		var nodeID ids.NodeID
		nodeID, err = ids.NodeIDFromString(v.NodeID)
		if err != nil {
			return nil, fmt.Errorf(
				"unable to turn string %v, pos %d to ids.ID: %w",
				v.NodeID,
				i,
				err,
			)
		}

		var val *Validator
		val, err = NewValidator(nodeID, v.Weight)
		if err != nil {
			return nil, fmt.Errorf("failed building validator: %w", err)
		}

		inMemGen.Validators = append(inMemGen.Validators, val)
	}

	// Step 3: InMemoryGenesis to in memory bytes
	bytes, err := encoding.Encode(inMemGen)
	if err != nil {
		return nil, fmt.Errorf("failed encoding genesis data: %w", err)
	}

	return bytes, nil
}

// BuildBase64GenesisString is used in tools to build a genesis
// (to be fed to a node via cli).
func BuildBase64GenesisString(base64Gen *Base64Genesis) (string, error) {
	bytes, err := encoding.Encode(base64Gen)
	if err != nil {
		return "", fmt.Errorf("failed encoding base64 genesis: %w", err)
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}
