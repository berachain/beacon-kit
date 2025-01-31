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

package compare_test

import (
	"bytes"
	"slices"
	"testing"
	"testing/quick"
	"unsafe"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/primitives/math"
	zcommon "github.com/protolambda/zrnt/eth2/beacon/common"
	zdeneb "github.com/protolambda/zrnt/eth2/beacon/deneb"
	zspec "github.com/protolambda/zrnt/eth2/configs"
	ztree "github.com/protolambda/ztyp/tree"
	zview "github.com/protolambda/ztyp/view"
	pprim "github.com/prysmaticlabs/prysm/v5/consensus-types/primitives"
	pethpb "github.com/prysmaticlabs/prysm/v5/proto/prysm/v1alpha1"
)

var (
	hFn  = ztree.GetHashFn()
	spec = zspec.Mainnet
)

func TestExecutionPayloadHashTreeRootZrnt(t *testing.T) {
	f := func(payload *ctypes.ExecutionPayload, logsBloom [256]byte) bool {
		// skip these cases lest we trigger a
		// nil-pointer dereference in fastssz
		if payload == nil ||
			payload.Withdrawals == nil ||
			slices.Contains(payload.Withdrawals, nil) ||
			payload.Transactions == nil ||
			slices.ContainsFunc(payload.Transactions, func(e []byte) bool {
				return e == nil
			}) {
			return true
		}

		payload.LogsBloom = logsBloom
		payload.BaseFeePerGas = math.NewU256(123)
		typeRoot := payload.HashTreeRoot()

		baseFeePerGas := zview.Uint256View{}
		baseFeePerGas.SetFromBig(payload.BaseFeePerGas.ToBig())
		zpayload := zdeneb.ExecutionPayload{
			ParentHash:    ztree.Root(payload.ParentHash),
			FeeRecipient:  zcommon.Eth1Address(payload.FeeRecipient),
			StateRoot:     ztree.Root(payload.StateRoot),
			ReceiptsRoot:  ztree.Root(payload.ReceiptsRoot),
			LogsBloom:     zcommon.LogsBloom(payload.LogsBloom),
			PrevRandao:    ztree.Root(payload.Random),
			BlockNumber:   zview.Uint64View(payload.Number),
			GasLimit:      zview.Uint64View(payload.GasLimit),
			GasUsed:       zview.Uint64View(payload.GasUsed),
			Timestamp:     zcommon.Timestamp(payload.Timestamp),
			ExtraData:     []byte(payload.ExtraData),
			BaseFeePerGas: baseFeePerGas,
			BlockHash:     ztree.Root(payload.BlockHash),
			Transactions: *(*zcommon.PayloadTransactions)(
				unsafe.Pointer(&payload.Transactions)),
			Withdrawals:   *(*zcommon.Withdrawals)(unsafe.Pointer(&payload.Withdrawals)),
			BlobGasUsed:   zview.Uint64View(payload.BlobGasUsed.Unwrap()),
			ExcessBlobGas: zview.Uint64View(payload.ExcessBlobGas.Unwrap()),
		}

		zRoot := zpayload.HashTreeRoot(spec, hFn)
		containerRoot := payload.HashTreeRoot()

		return bytes.Equal(typeRoot[:], containerRoot[:]) &&
			bytes.Equal(typeRoot[:], zRoot[:])
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 5_000}); err != nil {
		t.Error(err)
	}
}

func TestBlobSidecarTreeRootPrysm(t *testing.T) {
	f := func(sidecar *datypes.BlobSidecar) bool {
		// skip these cases lest we trigger a
		// nil-pointer dereference in fastssz
		if sidecar == nil ||
			sidecar.InclusionProof == nil ||
			sidecar.SignedBeaconBlockHeader == nil ||
			sidecar.SignedBeaconBlockHeader.Header == nil ||

			// prysm allows only sidecars whose InclusionProof has
			// length 17, while beaconKit allows different length.
			// We only keep 17 long Inclusion proofs for proper comparison
			len(sidecar.InclusionProof) != 17 {
			return true
		}

		sBlkHeader := sidecar.SignedBeaconBlockHeader
		blkHeader := sBlkHeader.Header

		pBlobSidecar := &pethpb.BlobSidecar{
			Index:         sidecar.Index,
			Blob:          sidecar.Blob[:],
			KzgCommitment: sidecar.KzgCommitment[:],
			KzgProof:      sidecar.KzgProof[:],
			SignedBlockHeader: &pethpb.SignedBeaconBlockHeader{
				Header: &pethpb.BeaconBlockHeader{
					Slot:          pprim.Slot(blkHeader.Slot),
					ProposerIndex: pprim.ValidatorIndex(blkHeader.ProposerIndex),
					ParentRoot:    blkHeader.ParentBlockRoot[:],
					StateRoot:     blkHeader.StateRoot[:],
					BodyRoot:      blkHeader.BodyRoot[:],
				},
				Signature: sBlkHeader.Signature[:],
			},
		}

		// Setup inclusion proofs
		inclusionProofs := sidecar.InclusionProof
		pBlobSidecar.CommitmentInclusionProof = make([][]byte, len(inclusionProofs))
		for i, proof := range inclusionProofs {
			pBlobSidecar.CommitmentInclusionProof[i] = proof[:]
		}

		beaconRoot := sidecar.HashTreeRoot()
		prysmRoot, err := pBlobSidecar.HashTreeRoot()
		if err != nil {
			t.Error(err)
		}

		return bytes.Equal(prysmRoot[:], beaconRoot[:])
	}

	if err := quick.Check(f, &quick.Config{MaxCount: 2_500}); err != nil {
		t.Error(err)
	}
}
