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

package version

import (
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/common"
)

const (
	// phase0 is the first version of the Beacon Chain.
	phase0 uint32 = 0
	// altair is the first hardfork of the Beacon Chain.
	altair uint32 = 1
	// bellatrix is the second hardfork of the Beacon Chain.
	bellatrix uint32 = 2
	// capella is the third hardfork of the Beacon Chain.
	capella uint32 = 3
	// deneb is the first version of Deneb, used for genesis of Berachain mainnet.
	deneb uint32 = 4
	// deneb1 is the first hardfork of Deneb on Berachain mainnet. LittleEndian of {4, 1, 0, 0}.
	deneb1 uint32 = 260
	// electra is the first version of Electra on Berachain mainnet.
	electra uint32 = 5
)

// Genesis returns the fork version for the genesis of Berachain mainnet, which is Deneb.
func Genesis() common.Version {
	return Deneb()
}

// Phase0 returns phase0 as a common.Version.
func Phase0() common.Version {
	return bytes.FromUint32(phase0)
}

// Altair returns altair as a common.Version.
func Altair() common.Version {
	return bytes.FromUint32(altair)
}

// Bellatrix returns bellatrix as a common.Version.
func Bellatrix() common.Version {
	return bytes.FromUint32(bellatrix)
}

// Capella returns capella as a common.Version.
func Capella() common.Version {
	return bytes.FromUint32(capella)
}

// Deneb returns deneb as a common.Version.
func Deneb() common.Version {
	return bytes.FromUint32(deneb)
}

// Deneb1 returns deneb1 as a common.Version.
func Deneb1() common.Version {
	return bytes.FromUint32(deneb1)
}

// Electra returns electra as a common.Version.
func Electra() common.Version {
	return bytes.FromUint32(electra)
}
