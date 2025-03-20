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

package sszconstructors

import (
	"fmt"

	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/karalabe/ssz"
)

// SSZUnmarshal is the way we build objects from byte formatted as ssz
func SSZUnmarshal[T interface {
	ssz.Object
	constraints.SSZUnmarshaler
}](buf []byte, v T) error {
	if v.IsUnusedFromSZZ() {
		// we special case construction of unused types, for efficiency
		if len(buf) != 1 {
			return fmt.Errorf("expected 1 byte got %d", len(buf))
		}
		//#nosec:G701 // UnusedType is uint8 and byte is uint8.
		tmp := types.UnusedType(buf[0])
		v, _ = any(tmp).(T) // TODO: any way this could be avoided?
		return nil
	}

	if err := ssz.DecodeFromBytes(buf, v); err != nil {
		return err
	}
	return v.EnsureSyntaxFromSSZ()
}
