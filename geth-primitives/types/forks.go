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

//nolint:all // This package contains types copied/adapted from go-ethereum (geth).
package types

// Fork is a numerical identifier of specific network upgrades (forks).
type Fork int

const (
	Frontier Fork = iota
	FrontierThawing
	Homestead
	DAO
	TangerineWhistle // a.k.a. the EIP150 fork
	SpuriousDragon   // a.k.a. the EIP155 fork
	Byzantium
	Constantinople
	Petersburg
	Istanbul
	MuirGlacier
	Berlin
	London
	ArrowGlacier
	GrayGlacier
	Paris
	Shanghai
	Cancun
	Prague
	Prague1
	Prague2
	Prague3
	Prague4
	Osaka
	BPO1
	BPO2
	BPO3
	BPO4
	BPO5
	Amsterdam
)
