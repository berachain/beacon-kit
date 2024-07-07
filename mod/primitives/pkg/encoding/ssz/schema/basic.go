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

package schema

import (
	"errors"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/types/types"
)

// basic represents a basic SSZ type.
type basic uint64

// ID returns the type ID of the basic type.
func (b basic) ID() types.Type { return types.Basic }

// ItemLength returns the size of the basic type in bytes.
func (b basic) ItemLength() uint64 { return uint64(b) }

// Chunks returns the number of 32-byte chunks required to represent the basic
// type.
func (b basic) HashChunkCount() uint64 { return 1 }

// child returns the basic type itself, as it has no children.
func (b basic) child(_ string) SSZType { return b }

// position always returns an error for basic types, as they have no children.
func (b basic) position(_ string) (uint64, uint8, error) {
	return 0, 0, errors.New("basic type has no children")
}
