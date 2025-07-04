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

package types

//go:generate sszgen -path sync_aggregate.go -objs SyncAggregate -output sync_aggregate_sszgen.go -include ../../primitives/common,../../primitives/bytes,../../primitives/crypto

import (
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/berachain/beacon-kit/primitives/crypto"
)

var (
	_ constraints.SSZMarshallableRootable = (*SyncAggregate)(nil)
	_ common.UnusedEnforcer               = (*SyncAggregate)(nil)
)

const (
	syncCommitteeSize       = 512
	syncCommitteeBitsLength = syncCommitteeSize / 8
)

type SyncAggregate struct {
	SyncCommitteeBits      [64]byte            `ssz-size:"64"`
	SyncCommitteeSignature crypto.BLSSignature `ssz-size:"96"`
}

func (s *SyncAggregate) ValidateAfterDecodingSSZ() error { return s.EnforceUnused() }

// EnforceUnused return true if the SyncAggregate contains all zero values.
// As long as this type remains unused and unvalidated by consensus,
// we must enforce that it contains no data.
func (s *SyncAggregate) EnforceUnused() error {
	if (s != nil && *s != SyncAggregate{}) {
		return errors.New("SyncAggregate must be unused")
	}
	return nil
}

