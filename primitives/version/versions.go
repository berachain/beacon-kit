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

package version

const (
	// Phase0 is the first version of the Beacon Chain.
	Phase0 uint32 = 0
	// Altair is the first hardfork of the Beacon Chain.
	Altair uint32 = 1
	// Bellatrix is the second hardfork of the Beacon Chain.
	Bellatrix uint32 = 2
	// Capella is the third hardfork of the Beacon Chain.
	Capella uint32 = 3
	// Deneb is the first version of Deneb, used for genesis of Berachain mainnet.
	Deneb uint32 = 4
	// Deneb1 is the first hardfork of Deneb on Berachain mainnet (TBD if used).
	// There may also be Deneb2, Deneb3, etc. hardforks.
	Deneb1 uint32 = 260
	// Electra is the first version of Electra on Berachain mainnet.
	Electra uint32 = 5
	// Electra1 is the first hardfork of Electra on Berachain mainnet (TBD if used).
	// There may also be Electra2, Electra3, etc. hardforks.
	Electra1 uint32 = 261
)
