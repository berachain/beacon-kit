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

package constants

const (
	// BlobCommitmentVersion is the version of the blob commitment.
	// It is the Version byte for the point evaluation precompile as
	// defined in EIP-4844.
	//
	// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-4844.md
	BlobCommitmentVersion uint8 = 0x01

	// MaxBlobCommitmentsPerBlock is the hardfork-independent fixed
	// theoretical limit same as TARGET_BLOB_GAS_PER_BLOCK (see EIP 4844).
	//
	// https://ethereum.github.io/consensus-specs/specs/deneb/beacon-chain/#execution
	MaxBlobCommitmentsPerBlock = 4096

	// MaxBlobSidecarsPerBlock is the maximum number of blob sidecars that can
	// be included in a block.
	MaxBlobSidecarsPerBlock = 6
)
